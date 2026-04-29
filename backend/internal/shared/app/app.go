package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sky0621/techcv-app/backend/internal/repository"
	"github.com/sky0621/techcv-app/backend/internal/shared/httpserver"
)

type closable interface {
	Close() error
}

type App struct {
	Handler http.Handler
	closers []closable
}

func New(ctx context.Context) (*App, error) {
	sqliteDSN := os.Getenv("SQLITE_DSN")
	if sqliteDSN == "" {
		sqliteDSN = filepath.Join("data", "techcv.db")
	}

	profileRepository, err := repository.NewSQLiteProfileRepository(ctx, sqliteDSN, schemaPath())
	if err != nil {
		return nil, fmt.Errorf("create profile repository: %w", err)
	}

	if err := profileRepository.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &App{
		Handler: httpserver.NewRouter(profileRepository, profileRepository, profileRepository),
		closers: []closable{profileRepository},
	}, nil
}

func (a *App) Close() error {
	for _, closer := range a.closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}

	return nil
}

func schemaPath() string {
	if value := os.Getenv("SQLITE_SCHEMA_PATH"); value != "" {
		return value
	}

	candidates := []string{
		filepath.Join("migrations", "schema.sql"),
		filepath.Join("backend", "migrations", "schema.sql"),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return candidates[0]
}
