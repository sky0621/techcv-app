package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sky0621/techcv-app/backend/internal/domain"
)

func TestMySQLProfileRepositoryGetCreatesDefaultProfile(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	testDB := newMySQLProfileTestDatabase(t)

	repo, err := NewMySQLProfileRepository(testDB.dsn())
	if err != nil {
		t.Fatalf("NewMySQLProfileRepository() error = %v", err)
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
	if got.VisibilitySettings["email"] != false || got.VisibilitySettings["phone"] != false {
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

func TestMySQLProfileRepositorySavePersistsProfile(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	testDB := newMySQLProfileTestDatabase(t)

	repo, err := NewMySQLProfileRepository(testDB.dsn())
	if err != nil {
		t.Fatalf("NewMySQLProfileRepository() error = %v", err)
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
	if reloaded.VisibilitySettings["github"] != true || reloaded.VisibilitySettings["email"] != false {
		t.Fatalf("unexpected reloaded visibility settings: %#v", reloaded.VisibilitySettings)
	}
}

type mysqlProfileTestDatabase struct {
	name     string
	host     string
	port     string
	user     string
	password string
	adminDB  *sql.DB
}

func newMySQLProfileTestDatabase(t *testing.T) *mysqlProfileTestDatabase {
	t.Helper()

	cfg := mysqlTestConfigFromEnv()
	adminDB, err := sql.Open("mysql", cfg.adminDSN())
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	if err := adminDB.PingContext(ctx); err != nil {
		_ = adminDB.Close()
		t.Skipf("skipping MySQL integration test: cannot connect to MySQL: %v", err)
	}

	db := &mysqlProfileTestDatabase{
		name:     fmt.Sprintf("techcv_app_test_%d", time.Now().UnixNano()),
		host:     cfg.host,
		port:     cfg.port,
		user:     cfg.user,
		password: cfg.password,
		adminDB:  adminDB,
	}

	if _, err := adminDB.ExecContext(ctx, "CREATE DATABASE "+db.name); err != nil {
		_ = adminDB.Close()
		t.Fatalf("CREATE DATABASE error = %v", err)
	}

	t.Cleanup(func() {
		dropCtx, dropCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer dropCancel()
		_, _ = adminDB.ExecContext(dropCtx, "DROP DATABASE IF EXISTS "+db.name)
		_ = adminDB.Close()
	})

	if err := db.applySchema(ctx); err != nil {
		t.Fatalf("applySchema() error = %v", err)
	}

	return db
}

func (d *mysqlProfileTestDatabase) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", d.user, d.password, d.host, d.port, d.name)
}

func (d *mysqlProfileTestDatabase) applySchema(ctx context.Context) error {
	schemaPath := testSchemaPath()
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	testDB, err := sql.Open("mysql", d.dsn())
	if err != nil {
		return fmt.Errorf("open test database: %w", err)
	}
	defer testDB.Close()

	statements := strings.Split(string(schemaBytes), ";")
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		if _, err := testDB.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("apply statement %q: %w", statement, err)
		}
	}

	return nil
}

func testSchemaPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("..", "..", "migrations", "schema.sql")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "migrations", "schema.sql"))
}

type mysqlTestConfig struct {
	host     string
	port     string
	user     string
	password string
}

func mysqlTestConfigFromEnv() mysqlTestConfig {
	return mysqlTestConfig{
		host:     envOrDefault("MYSQL_HOST", "127.0.0.1"),
		port:     envOrDefault("MYSQL_PORT", "3306"),
		user:     envOrDefault("MYSQL_USER", "root"),
		password: envOrDefault("MYSQL_PASSWORD", "password"),
	}
}

func (c mysqlTestConfig) adminDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true", c.user, c.password, c.host, c.port)
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func profileFixture(createdAt time.Time) domain.Profile {
	return domain.Profile{
		ID:                 "profile_01",
		UserID:             "user_01",
		FullName:           "Sky Sample",
		Nickname:           "sky0621",
		Location:           "Tokyo",
		Email:              "me@example.com",
		Phone:              "090-0000-0000",
		Summary:            "Backend engineer",
		GitHubURL:          "https://github.com/sky0621",
		ZennURL:            "https://zenn.dev/sky0621",
		QiitaURL:           "https://qiita.com/sky0621",
		WebsiteURL:         "https://example.com",
		PreferredWorkStyle: "Full remote",
		VisibilitySettings: map[string]any{
			"email":  false,
			"github": true,
		},
		CreatedAt: createdAt,
	}
}
