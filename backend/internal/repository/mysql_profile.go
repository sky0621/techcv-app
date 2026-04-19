package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sky0621/techcv-app/backend/internal/domain"
	dbgen "github.com/sky0621/techcv-app/backend/internal/repository/db"
)

type MySQLProfileRepository struct {
	db      *sql.DB
	queries *dbgen.Queries
}

func NewMySQLProfileRepository(dsn string) (*MySQLProfileRepository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return &MySQLProfileRepository{
		db:      db,
		queries: dbgen.New(db),
	}, nil
}

func (r *MySQLProfileRepository) Close() error {
	if r.db == nil {
		return nil
	}

	return r.db.Close()
}

func (r *MySQLProfileRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *MySQLProfileRepository) Get(ctx context.Context) (*domain.Profile, error) {
	row, err := r.queries.GetProfileByUserID(ctx, "user_01")
	if err == nil {
		return toDomainProfile(row)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("query profile: %w", err)
	}

	now := time.Now().UTC()
	profile := domain.Profile{
		ID:                 "profile_01",
		UserID:             "user_01",
		VisibilitySettings: map[string]any{"email": false, "phone": false},
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return r.Save(ctx, &profile)
}

func (r *MySQLProfileRepository) Save(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	visibilitySettings := profile.VisibilitySettings
	if visibilitySettings == nil {
		visibilitySettings = map[string]any{}
	}

	visibilityBytes, err := json.Marshal(visibilitySettings)
	if err != nil {
		return nil, fmt.Errorf("encode visibility settings: %w", err)
	}

	now := time.Now().UTC()
	createdAt := profile.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}

	err = r.queries.UpsertProfile(ctx, dbgen.UpsertProfileParams{
		ID:                 profile.ID,
		UserID:             profile.UserID,
		FullName:           profile.FullName,
		Nickname:           profile.Nickname,
		Location:           profile.Location,
		Email:              profile.Email,
		Phone:              profile.Phone,
		Summary:            profile.Summary,
		GithubUrl:          profile.GitHubURL,
		ZennUrl:            profile.ZennURL,
		QiitaUrl:           profile.QiitaURL,
		WebsiteUrl:         profile.WebsiteURL,
		PreferredWorkStyle: profile.PreferredWorkStyle,
		VisibilitySettings: visibilityBytes,
		CreatedAt:          createdAt,
		UpdatedAt:          now,
	})
	if err != nil {
		return nil, fmt.Errorf("save profile: %w", err)
	}

	row, err := r.queries.GetProfileByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("reload profile: %w", err)
	}

	return toDomainProfile(row)
}

func toDomainProfile(row dbgen.Profile) (*domain.Profile, error) {
	profile := &domain.Profile{
		ID:                 row.ID,
		UserID:             row.UserID,
		FullName:           row.FullName,
		Nickname:           row.Nickname,
		Location:           row.Location,
		Email:              row.Email,
		Phone:              row.Phone,
		Summary:            row.Summary,
		GitHubURL:          row.GithubUrl,
		ZennURL:            row.ZennUrl,
		QiitaURL:           row.QiitaUrl,
		WebsiteURL:         row.WebsiteUrl,
		PreferredWorkStyle: row.PreferredWorkStyle,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}

	if len(row.VisibilitySettings) > 0 {
		if err := json.Unmarshal(row.VisibilitySettings, &profile.VisibilitySettings); err != nil {
			return nil, fmt.Errorf("decode visibility settings: %w", err)
		}
	}

	if profile.VisibilitySettings == nil {
		profile.VisibilitySettings = map[string]any{}
	}

	return profile, nil
}
