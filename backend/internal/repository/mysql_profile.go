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
	if err := r.addSkillCategoryIconColumn(ctx); err != nil {
		return err
	}
	if err := r.seedSkillOptions(ctx); err != nil {
		return err
	}
	if err := r.backfillEmptySkillCategoryIcons(ctx); err != nil {
		return err
	}

	return r.addQualificationURLColumn(ctx)
}

func (r *SQLiteProfileRepository) ListSkillOptions(ctx context.Context) (*domain.SkillOptions, error) {
	categories, err := r.queries.ListSkillCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("query skill categories: %w", err)
	}

	proficiencyLevels, err := r.queries.ListSkillProficiencyLevels(ctx)
	if err != nil {
		return nil, fmt.Errorf("query skill proficiency levels: %w", err)
	}

	return &domain.SkillOptions{
		Categories:        toDomainSkillCategoryOptions(categories),
		ProficiencyLevels: toDomainSkillProficiencyLevelOptions(proficiencyLevels),
	}, nil
}

func (r *SQLiteProfileRepository) CreateSkillCategory(ctx context.Context, input domain.SkillCategoryInput) (*domain.SkillOption, error) {
	row, err := r.queries.InsertSkillCategory(ctx, dbgen.InsertSkillCategoryParams{
		ID:   input.ID,
		Name: input.Name,
		Icon: input.Icon,
	})
	if err != nil {
		return nil, fmt.Errorf("insert skill category: %w", err)
	}

	return toDomainSkillCategoryOption(row), nil
}

func (r *SQLiteProfileRepository) UpdateSkillCategory(ctx context.Context, id string, input domain.SkillCategoryInput) (*domain.SkillOption, error) {
	row, err := r.queries.UpdateSkillCategory(ctx, dbgen.UpdateSkillCategoryParams{
		ID:   id,
		Name: input.Name,
		Icon: input.Icon,
	})
	if err != nil {
		return nil, fmt.Errorf("update skill category: %w", err)
	}

	return toDomainSkillCategoryOption(row), nil
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

func (r *SQLiteProfileRepository) addSkillCategoryIconColumn(ctx context.Context) error {
	hasIconColumn, err := r.tableHasColumn(ctx, "skill_categories", "icon")
	if err != nil {
		return err
	}
	if hasIconColumn {
		return nil
	}

	if _, err := r.db.ExecContext(ctx, "ALTER TABLE skill_categories ADD COLUMN icon TEXT NOT NULL DEFAULT ''"); err != nil {
		return fmt.Errorf("add skill_categories.icon column: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) seedSkillOptions(ctx context.Context) error {
	categories := []domain.SkillOption{
		{ID: "skill_category_language", Name: "言語", Icon: "code", SortOrder: 1},
		{ID: "skill_category_framework", Name: "フレームワーク", Icon: "code", SortOrder: 2},
		{ID: "skill_category_database", Name: "データベース", Icon: "database", SortOrder: 3},
		{ID: "skill_category_infrastructure", Name: "インフラ", Icon: "cloud", SortOrder: 4},
		{ID: "skill_category_tool", Name: "ツール", Icon: "wrench", SortOrder: 5},
		{ID: "skill_category_other", Name: "その他", Icon: "wrench", SortOrder: 6},
	}
	for _, category := range categories {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO skill_categories (id, name, icon, sort_order) VALUES (?, ?, ?, ?)",
			category.ID,
			category.Name,
			category.Icon,
			category.SortOrder,
		); err != nil {
			return fmt.Errorf("seed skill category %s: %w", category.ID, err)
		}
	}

	proficiencyLevels := []domain.SkillOption{
		{ID: "skill_proficiency_beginner", Name: "初級", SortOrder: 1},
		{ID: "skill_proficiency_intermediate", Name: "中級", SortOrder: 2},
		{ID: "skill_proficiency_advanced", Name: "上級", SortOrder: 3},
		{ID: "skill_proficiency_expert", Name: "エキスパート", SortOrder: 4},
	}
	for _, level := range proficiencyLevels {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO skill_proficiency_levels (id, name, sort_order) VALUES (?, ?, ?)",
			level.ID,
			level.Name,
			level.SortOrder,
		); err != nil {
			return fmt.Errorf("seed skill proficiency level %s: %w", level.ID, err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) backfillEmptySkillCategoryIcons(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE skill_categories SET icon = 'wrench', updated_at = CURRENT_TIMESTAMP WHERE icon = ''"); err != nil {
		return fmt.Errorf("backfill empty skill category icons: %w", err)
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

func toDomainSkillCategoryOptions(rows []dbgen.SkillCategory) []domain.SkillOption {
	options := make([]domain.SkillOption, 0, len(rows))
	for _, row := range rows {
		options = append(options, *toDomainSkillCategoryOption(row))
	}

	return options
}

func toDomainSkillCategoryOption(row dbgen.SkillCategory) *domain.SkillOption {
	return &domain.SkillOption{
		ID:        row.ID,
		Name:      row.Name,
		Icon:      row.Icon,
		SortOrder: row.SortOrder,
	}
}

func toDomainSkillProficiencyLevelOptions(rows []dbgen.SkillProficiencyLevel) []domain.SkillOption {
	options := make([]domain.SkillOption, 0, len(rows))
	for _, row := range rows {
		options = append(options, domain.SkillOption{
			ID:        row.ID,
			Name:      row.Name,
			SortOrder: row.SortOrder,
		})
	}

	return options
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
