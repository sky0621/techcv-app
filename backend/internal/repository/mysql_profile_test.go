package repository

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/sky0621/techcv-app/backend/internal/domain"
)

func TestSQLiteProfileRepositoryGetCreatesDefaultProfile(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	testDB := newSQLiteProfileTestDatabase(t)

	repo, err := NewSQLiteProfileRepository(ctx, testDB.dsn(), testSchemaPath())
	if err != nil {
		t.Fatalf("NewSQLiteProfileRepository() error = %v", err)
	}
	t.Cleanup(func() {
		_ = repo.Close()
	})

	got, err := repo.Get(ctx)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.ID != "profile_01" {
		t.Fatalf("expected default ID profile_01, got %q", got.ID)
	}
	if got.UserID != "user_01" {
		t.Fatalf("expected default UserID user_01, got %q", got.UserID)
	}
	if got.VisibilitySettings["email"] != true || got.VisibilitySettings["location"] != true {
		t.Fatalf("unexpected default visibility settings: %#v", got.VisibilitySettings)
	}

	reloaded, err := repo.Get(ctx)
	if err != nil {
		t.Fatalf("second Get() error = %v", err)
	}

	if !got.CreatedAt.Equal(reloaded.CreatedAt) {
		t.Fatalf("expected created_at to persist, first=%s second=%s", got.CreatedAt, reloaded.CreatedAt)
	}
}

func TestSQLiteProfileRepositorySavePersistsProfile(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	testDB := newSQLiteProfileTestDatabase(t)

	repo, err := NewSQLiteProfileRepository(ctx, testDB.dsn(), testSchemaPath())
	if err != nil {
		t.Fatalf("NewSQLiteProfileRepository() error = %v", err)
	}
	t.Cleanup(func() {
		_ = repo.Close()
	})

	createdAt := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
	profile := profileFixture(createdAt)
	saved, err := repo.Save(ctx, &profile)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if saved.FullName != "Sky Sample" || saved.Email != "me@example.com" {
		t.Fatalf("unexpected saved profile: %+v", saved)
	}
	if saved.VisibilitySettings["github"] != true {
		t.Fatalf("expected visibility settings to persist, got %#v", saved.VisibilitySettings)
	}
	if !saved.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected CreatedAt %s, got %s", createdAt, saved.CreatedAt)
	}
	if !saved.UpdatedAt.After(createdAt) {
		t.Fatalf("expected UpdatedAt after CreatedAt, created=%s updated=%s", createdAt, saved.UpdatedAt)
	}

	reloaded, err := repo.Get(ctx)
	if err != nil {
		t.Fatalf("Get() after Save() error = %v", err)
	}

	if reloaded.FullName != "Sky Sample" || reloaded.PreferredWorkStyle != "Full remote" {
		t.Fatalf("unexpected reloaded profile: %+v", reloaded)
	}
	if reloaded.Occupation != "Software Engineer" || reloaded.EmploymentType != "Freelance" {
		t.Fatalf("unexpected reloaded work fields: %+v", reloaded)
	}
	if reloaded.VisibilitySettings["github"] != true || reloaded.VisibilitySettings["email"] != false {
		t.Fatalf("unexpected reloaded visibility settings: %#v", reloaded.VisibilitySettings)
	}
	if len(reloaded.Qualifications) != 2 {
		t.Fatalf("expected qualifications to persist, got %#v", reloaded.Qualifications)
	}
	if reloaded.Qualifications[0].ID != "qualification_01" || reloaded.Qualifications[0].Name != "AWS Certified Solutions Architect" {
		t.Fatalf("unexpected first qualification: %#v", reloaded.Qualifications[0])
	}
	if reloaded.Qualifications[0].URL != "https://aws.amazon.com/certification/certified-solutions-architect-associate/" {
		t.Fatalf("unexpected first qualification URL: %#v", reloaded.Qualifications[0])
	}
	if reloaded.Qualifications[1].Name != "基本情報技術者" || reloaded.Qualifications[1].Organization != "IPA" {
		t.Fatalf("unexpected second qualification: %#v", reloaded.Qualifications[1])
	}

	reloaded.Qualifications = []domain.Qualification{
		{
			Name:         "Google Cloud Professional Cloud Architect",
			AcquiredDate: "2025-01-01",
			Organization: "Google Cloud",
			URL:          "https://cloud.google.com/learn/certification/cloud-architect",
			Memo:         "replacement",
		},
	}
	replaced, err := repo.Save(ctx, reloaded)
	if err != nil {
		t.Fatalf("Save() replacement error = %v", err)
	}
	if len(replaced.Qualifications) != 1 {
		t.Fatalf("expected qualifications to be replaced, got %#v", replaced.Qualifications)
	}
	if replaced.Qualifications[0].ID == "" {
		t.Fatalf("expected replacement qualification ID to be generated, got %#v", replaced.Qualifications[0])
	}
	if replaced.Qualifications[0].Name != "Google Cloud Professional Cloud Architect" {
		t.Fatalf("unexpected replacement qualification: %#v", replaced.Qualifications[0])
	}
	if replaced.Qualifications[0].URL != "https://cloud.google.com/learn/certification/cloud-architect" {
		t.Fatalf("unexpected replacement qualification URL: %#v", replaced.Qualifications[0])
	}
}

func TestSQLiteProfileRepositoryMigratesProfilePhoneColumn(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	testDB := newSQLiteProfileTestDatabase(t)

	rawDB, err := sql.Open("sqlite3", testDB.dsn())
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	if _, err := rawDB.ExecContext(ctx, oldProfileSchema); err != nil {
		_ = rawDB.Close()
		t.Fatalf("failed to create old schema: %v", err)
	}
	if err := rawDB.Close(); err != nil {
		t.Fatalf("failed to close old schema database: %v", err)
	}

	repo, err := NewSQLiteProfileRepository(ctx, testDB.dsn(), testSchemaPath())
	if err != nil {
		t.Fatalf("NewSQLiteProfileRepository() error = %v", err)
	}
	t.Cleanup(func() {
		_ = repo.Close()
	})

	if _, err := repo.Get(ctx); err != nil {
		t.Fatalf("Get() after migration error = %v", err)
	}

	hasPhoneColumn, err := repo.profileTableHasColumn(ctx, "phone")
	if err != nil {
		t.Fatalf("profileTableHasColumn() error = %v", err)
	}
	if hasPhoneColumn {
		t.Fatal("expected profiles.phone column to be dropped")
	}
	hasOccupationColumn, err := repo.profileTableHasColumn(ctx, "occupation")
	if err != nil {
		t.Fatalf("profileTableHasColumn(occupation) error = %v", err)
	}
	if !hasOccupationColumn {
		t.Fatal("expected profiles.occupation column to be added")
	}
	hasEmploymentTypeColumn, err := repo.profileTableHasColumn(ctx, "employment_type")
	if err != nil {
		t.Fatalf("profileTableHasColumn(employment_type) error = %v", err)
	}
	if !hasEmploymentTypeColumn {
		t.Fatal("expected profiles.employment_type column to be added")
	}
}

func TestSQLiteProfileRepositoryMigratesQualificationURLColumn(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	testDB := newSQLiteProfileTestDatabase(t)

	rawDB, err := sql.Open("sqlite3", testDB.dsn())
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	if _, err := rawDB.ExecContext(ctx, oldProfileSchema+oldQualificationSchemaWithoutURL); err != nil {
		_ = rawDB.Close()
		t.Fatalf("failed to create old qualification schema: %v", err)
	}
	if err := rawDB.Close(); err != nil {
		t.Fatalf("failed to close old qualification schema database: %v", err)
	}

	repo, err := NewSQLiteProfileRepository(ctx, testDB.dsn(), testSchemaPath())
	if err != nil {
		t.Fatalf("NewSQLiteProfileRepository() error = %v", err)
	}
	t.Cleanup(func() {
		_ = repo.Close()
	})

	hasURLColumn, err := repo.qualificationTableHasColumn(ctx, "url")
	if err != nil {
		t.Fatalf("qualificationTableHasColumn() error = %v", err)
	}
	if !hasURLColumn {
		t.Fatal("expected profile_qualifications.url column to be added")
	}
}

type sqliteProfileTestDatabase struct {
	path string
}

func newSQLiteProfileTestDatabase(t *testing.T) *sqliteProfileTestDatabase {
	t.Helper()

	return &sqliteProfileTestDatabase{
		path: filepath.Join(t.TempDir(), "techcv-test.db"),
	}
}

func (d *sqliteProfileTestDatabase) dsn() string {
	return d.path
}

func testSchemaPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("..", "..", "migrations", "schema.sql")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "migrations", "schema.sql"))
}

func profileFixture(createdAt time.Time) domain.Profile {
	return domain.Profile{
		ID:                 "profile_01",
		UserID:             "user_01",
		FullName:           "Sky Sample",
		Nickname:           "sky0621",
		Location:           "Tokyo",
		Email:              "me@example.com",
		Summary:            "Backend engineer",
		GitHubURL:          "https://github.com/sky0621",
		ZennURL:            "https://zenn.dev/sky0621",
		QiitaURL:           "https://qiita.com/sky0621",
		WebsiteURL:         "https://example.com",
		Occupation:         "Software Engineer",
		EmploymentType:     "Freelance",
		PreferredWorkStyle: "Full remote",
		VisibilitySettings: map[string]any{
			"email":  false,
			"github": true,
		},
		Qualifications: []domain.Qualification{
			{
				ID:           "qualification_01",
				Name:         "AWS Certified Solutions Architect",
				AcquiredDate: "2026-04-26",
				Organization: "Amazon Web Services",
				URL:          "https://aws.amazon.com/certification/certified-solutions-architect-associate/",
				Memo:         "Associate",
			},
			{
				ID:           "qualification_02",
				Name:         "基本情報技術者",
				AcquiredDate: "2020-10-01",
				Organization: "IPA",
				URL:          "https://www.ipa.go.jp/shiken/kubun/fe.html",
				Memo:         "",
			},
		},
		CreatedAt: createdAt,
	}
}

const oldProfileSchema = `
CREATE TABLE profiles (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL,
    nickname TEXT NOT NULL,
    location TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    summary TEXT NOT NULL,
    github_url TEXT NOT NULL,
    zenn_url TEXT NOT NULL,
    qiita_url TEXT NOT NULL,
    website_url TEXT NOT NULL,
    preferred_work_style TEXT NOT NULL,
    visibility_settings TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

const oldQualificationSchemaWithoutURL = `
CREATE TABLE profile_qualifications (
    id TEXT NOT NULL PRIMARY KEY,
    profile_id TEXT NOT NULL,
    name TEXT NOT NULL,
    acquired_date TEXT NOT NULL,
    organization TEXT NOT NULL,
    memo TEXT NOT NULL,
    sort_order INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);`
