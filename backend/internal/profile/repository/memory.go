package repository

import (
	"context"
	"sync"
	"time"

	"github.com/sky0621/techcv-app/backend/internal/profile/domain"
)

type Repository interface {
	Get(ctx context.Context) (*domain.Profile, error)
	Save(ctx context.Context, profile *domain.Profile) (*domain.Profile, error)
}

type MemoryRepository struct {
	mu      sync.RWMutex
	profile *domain.Profile
}

func NewMemoryRepository() *MemoryRepository {
	now := time.Now().UTC()

	return &MemoryRepository{
		profile: &domain.Profile{
			ID:                 "profile_01",
			UserID:             "user_01",
			VisibilitySettings: map[string]any{"email": false, "phone": false},
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}
}

func (r *MemoryRepository) Get(_ context.Context) (*domain.Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	profileCopy := *r.profile
	return &profileCopy, nil
}

func (r *MemoryRepository) Save(_ context.Context, profile *domain.Profile) (*domain.Profile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	profileCopy := *profile
	profileCopy.UpdatedAt = time.Now().UTC()
	if profileCopy.CreatedAt.IsZero() {
		profileCopy.CreatedAt = profileCopy.UpdatedAt
	}

	r.profile = &profileCopy

	savedCopy := *r.profile
	return &savedCopy, nil
}
