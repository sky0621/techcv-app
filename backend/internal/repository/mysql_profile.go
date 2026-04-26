package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sky0621/techcv-app/backend/internal/domain"
	dbgen "github.com/sky0621/techcv-app/backend/internal/repository/db"
)

type SQLiteProfileRepository struct {
	db      *sql.DB
	queries *dbgen.Queries
}

func NewSQLiteProfileRepository(ctx context.Context, dsn string, schemaPath string) (*SQLiteProfileRepository, error) {
	if err := ensureSQLiteDirectory(dsn); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(1)

	repository := &SQLiteProfileRepository{
		db:      db,
		queries: dbgen.New(db),
	}

	if err := repository.configure(ctx); err != nil {
		_ = repository.Close()
		return nil, err
	}

	if err := repository.applySchema(ctx, schemaPath); err != nil {
		_ = repository.Close()
		return nil, err
	}

	return repository, nil
}

func (r *SQLiteProfileRepository) Close() error {
	if r.db == nil {
		return nil
	}

	return r.db.Close()
}

func (r *SQLiteProfileRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *SQLiteProfileRepository) Get(ctx context.Context) (*domain.Profile, error) {
	row, err := r.queries.GetProfileByUserID(ctx, "user_01")
	if err == nil {
		qualifications, err := r.queries.ListQualificationsByProfileID(ctx, row.ID)
		if err != nil {
			return nil, fmt.Errorf("query profile qualifications: %w", err)
		}

		return toDomainProfile(row, qualifications)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("query profile: %w", err)
	}

	now := time.Now().UTC()
	profile := domain.Profile{
		ID:                 "profile_01",
		UserID:             "user_01",
		VisibilitySettings: map[string]any{"email": true, "location": true},
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return r.Save(ctx, &profile)
}

func (r *SQLiteProfileRepository) Save(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	visibilitySettings := sanitizeVisibilitySettings(profile.VisibilitySettings)

	visibilityBytes, err := json.Marshal(visibilitySettings)
	if err != nil {
		return nil, fmt.Errorf("encode visibility settings: %w", err)
	}

	now := time.Now().UTC()
	createdAt := profile.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin profile transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	queries := r.queries.WithTx(tx)

	err = queries.UpsertProfile(ctx, dbgen.UpsertProfileParams{
		ID:                 profile.ID,
		UserID:             profile.UserID,
		FullName:           profile.FullName,
		Nickname:           profile.Nickname,
		Location:           profile.Location,
		Email:              profile.Email,
		Summary:            profile.Summary,
		GithubUrl:          profile.GitHubURL,
		ZennUrl:            profile.ZennURL,
		QiitaUrl:           profile.QiitaURL,
		WebsiteUrl:         profile.WebsiteURL,
		Occupation:         profile.Occupation,
		EmploymentType:     profile.EmploymentType,
		PreferredWorkStyle: profile.PreferredWorkStyle,
		VisibilitySettings: string(visibilityBytes),
		CreatedAt:          createdAt,
		UpdatedAt:          now,
	})
	if err != nil {
		return nil, fmt.Errorf("save profile: %w", err)
	}

	if err := queries.DeleteQualificationsByProfileID(ctx, profile.ID); err != nil {
		return nil, fmt.Errorf("delete profile qualifications: %w", err)
	}

	for index, qualification := range profile.Qualifications {
		if err := queries.InsertQualification(ctx, dbgen.InsertQualificationParams{
			ID:           qualificationID(qualification, index, now),
			ProfileID:    profile.ID,
			Name:         qualification.Name,
			AcquiredDate: qualification.AcquiredDate,
			Organization: qualification.Organization,
			Url:          qualification.URL,
			Memo:         qualification.Memo,
			SortOrder:    int64(index),
			CreatedAt:    now,
			UpdatedAt:    now,
		}); err != nil {
			return nil, fmt.Errorf("insert profile qualification: %w", err)
		}
	}

	row, err := queries.GetProfileByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("reload profile: %w", err)
	}

	qualifications, err := queries.ListQualificationsByProfileID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload profile qualifications: %w", err)
	}

	result, err := toDomainProfile(row, qualifications)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit profile transaction: %w", err)
	}

	return result, nil
}

func (r *SQLiteProfileRepository) configure(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, "PRAGMA busy_timeout = 5000"); err != nil {
		return fmt.Errorf("set busy timeout: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) applySchema(ctx context.Context, schemaPath string) error {
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	statements := strings.Split(string(schemaBytes), ";")
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("apply schema statement %q: %w", statement, err)
		}
	}

	if err := r.dropProfilePhoneColumn(ctx); err != nil {
		return err
	}
	if err := r.addProfileTextColumn(ctx, "occupation"); err != nil {
		return err
	}
	if err := r.addProfileTextColumn(ctx, "employment_type"); err != nil {
		return err
	}

	return r.addQualificationURLColumn(ctx)
}

func ensureSQLiteDirectory(dsn string) error {
	if dsn == ":memory:" || strings.HasPrefix(dsn, "file:") {
		return nil
	}

	dir := filepath.Dir(dsn)
	if dir == "." {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create sqlite data directory: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) dropProfilePhoneColumn(ctx context.Context) error {
	hasPhoneColumn, err := r.profileTableHasColumn(ctx, "phone")
	if err != nil {
		return err
	}
	if !hasPhoneColumn {
		return nil
	}

	if _, err := r.db.ExecContext(ctx, "ALTER TABLE profiles DROP COLUMN phone"); err != nil {
		return fmt.Errorf("drop profiles.phone column: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) addProfileTextColumn(ctx context.Context, columnName string) error {
	hasColumn, err := r.profileTableHasColumn(ctx, columnName)
	if err != nil {
		return err
	}
	if hasColumn {
		return nil
	}

	if _, err := r.db.ExecContext(ctx, "ALTER TABLE profiles ADD COLUMN "+columnName+" TEXT NOT NULL DEFAULT ''"); err != nil {
		return fmt.Errorf("add profiles.%s column: %w", columnName, err)
	}

	return nil
}

func (r *SQLiteProfileRepository) profileTableHasColumn(ctx context.Context, columnName string) (bool, error) {
	return r.tableHasColumn(ctx, "profiles", columnName)
}

func (r *SQLiteProfileRepository) qualificationTableHasColumn(ctx context.Context, columnName string) (bool, error) {
	return r.tableHasColumn(ctx, "profile_qualifications", columnName)
}

func (r *SQLiteProfileRepository) addQualificationURLColumn(ctx context.Context) error {
	hasURLColumn, err := r.qualificationTableHasColumn(ctx, "url")
	if err != nil {
		return err
	}
	if hasURLColumn {
		return nil
	}

	if _, err := r.db.ExecContext(ctx, "ALTER TABLE profile_qualifications ADD COLUMN url TEXT NOT NULL DEFAULT ''"); err != nil {
		return fmt.Errorf("add profile_qualifications.url column: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) tableHasColumn(ctx context.Context, tableName string, columnName string) (bool, error) {
	rows, err := r.db.QueryContext(ctx, "PRAGMA table_info("+tableName+")")
	if err != nil {
		return false, fmt.Errorf("inspect %s columns: %w", tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid          int
			name         string
			columnType   string
			notNull      int
			defaultValue sql.NullString
			primaryKey   int
		)
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, fmt.Errorf("scan %s column: %w", tableName, err)
		}
		if name == columnName {
			return true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("iterate %s columns: %w", tableName, err)
	}

	return false, nil
}

func toDomainProfile(row dbgen.Profile, qualificationRows []dbgen.ProfileQualification) (*domain.Profile, error) {
	profile := &domain.Profile{
		ID:                 row.ID,
		UserID:             row.UserID,
		FullName:           row.FullName,
		Nickname:           row.Nickname,
		Location:           row.Location,
		Email:              row.Email,
		Summary:            row.Summary,
		GitHubURL:          row.GithubUrl,
		ZennURL:            row.ZennUrl,
		QiitaURL:           row.QiitaUrl,
		WebsiteURL:         row.WebsiteUrl,
		Occupation:         row.Occupation,
		EmploymentType:     row.EmploymentType,
		PreferredWorkStyle: row.PreferredWorkStyle,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}

	if row.VisibilitySettings != "" {
		if err := json.Unmarshal([]byte(row.VisibilitySettings), &profile.VisibilitySettings); err != nil {
			return nil, fmt.Errorf("decode visibility settings: %w", err)
		}
	}

	profile.VisibilitySettings = sanitizeVisibilitySettings(profile.VisibilitySettings)
	profile.Qualifications = toDomainQualifications(qualificationRows)

	return profile, nil
}

func toDomainQualifications(rows []dbgen.ProfileQualification) []domain.Qualification {
	qualifications := make([]domain.Qualification, 0, len(rows))
	for _, row := range rows {
		qualifications = append(qualifications, domain.Qualification{
			ID:           row.ID,
			Name:         row.Name,
			AcquiredDate: row.AcquiredDate,
			Organization: row.Organization,
			URL:          row.Url,
			Memo:         row.Memo,
		})
	}

	return qualifications
}

func qualificationID(qualification domain.Qualification, index int, now time.Time) string {
	if qualification.ID != "" {
		return qualification.ID
	}

	return fmt.Sprintf("qualification_%d_%d", now.UnixNano(), index+1)
}

func sanitizeVisibilitySettings(values map[string]any) map[string]any {
	result := map[string]any{}
	for key, value := range values {
		if key == "phone" {
			continue
		}
		result[key] = value
	}

	return result
}
