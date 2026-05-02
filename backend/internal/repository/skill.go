package repository

import (
	"context"

	"github.com/sky0621/techcv-app/backend/internal/domain"
)

type SkillOptionsRepository interface {
	ListSkillOptions(ctx context.Context) (*domain.SkillOptions, error)
	CreateSkillCategory(ctx context.Context, input domain.SkillCategoryInput) (*domain.SkillOption, error)
	UpdateSkillCategory(ctx context.Context, id string, input domain.SkillCategoryInput) (*domain.SkillOption, error)
	UpdateSkillProficiencyLevel(ctx context.Context, id string, input domain.SkillProficiencyLevelInput) (*domain.SkillOption, error)
	CreateSkillMaster(ctx context.Context, input domain.SkillMasterInput) (*domain.SkillMaster, error)
	UpdateSkillMaster(ctx context.Context, id string, input domain.SkillMasterInput) (*domain.SkillMaster, error)
	ListSkills(ctx context.Context) ([]domain.Skill, error)
	CreateSkill(ctx context.Context, input domain.SkillInput) (*domain.Skill, error)
	UpdateSkill(ctx context.Context, id string, input domain.SkillInput) (*domain.Skill, error)
	DeleteSkill(ctx context.Context, id string) error
}
