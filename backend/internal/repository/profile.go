package repository

import (
	"context"

	"github.com/sky0621/techcv-app/backend/internal/domain"
)

type ProfileRepository interface {
	Get(ctx context.Context) (*domain.Profile, error)
	Save(ctx context.Context, profile *domain.Profile) (*domain.Profile, error)
	ListProfileLinkMasters(ctx context.Context) ([]domain.ProfileLinkMaster, error)
	CreateProfileLinkMaster(ctx context.Context, input domain.ProfileLinkMasterInput) (*domain.ProfileLinkMaster, error)
	UpdateProfileLinkMaster(ctx context.Context, id string, input domain.ProfileLinkMasterInput) (*domain.ProfileLinkMaster, error)
}
