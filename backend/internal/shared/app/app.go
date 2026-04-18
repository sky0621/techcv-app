package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	profilerepository "github.com/sky0621/techcv-app/backend/internal/profile/repository"
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
	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		return nil, fmt.Errorf("MYSQL_DSN is required")
	}

	profileRepository, err := profilerepository.NewMySQLRepository(mysqlDSN)
	if err != nil {
		return nil, err
	}

	if err := profileRepository.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &App{
		Handler: httpserver.NewRouter(profileRepository),
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
