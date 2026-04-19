package repository

import (
	"context"

	"github.com/sky0621/techcv-app/backend/internal/profile/domain"
)

type Repository interface {
	Get(ctx context.Context) (*domain.Profile, error)
	Save(ctx context.Context, profile *domain.Profile) (*domain.Profile, error)
}
