package repository

import (
	"context"

	"github.com/sky0621/techcv-app/backend/internal/domain"
)

type SkillOptionsRepository interface {
	ListSkillOptions(ctx context.Context) (*domain.SkillOptions, error)
	CreateSkillCategory(ctx context.Context, input domain.SkillCategoryInput) (*domain.SkillOption, error)
	UpdateSkillCategory(ctx context.Context, id string, input domain.SkillCategoryInput) (*domain.SkillOption, error)
}
