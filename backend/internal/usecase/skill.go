package usecase

import (
	"context"

	"github.com/sky0621/techcv-app/backend/internal/domain"
	"github.com/sky0621/techcv-app/backend/internal/repository"
)

type SkillOptionsUseCase struct {
	repository repository.SkillOptionsRepository
}

func NewSkillOptionsUseCase(repository repository.SkillOptionsRepository) *SkillOptionsUseCase {
	return &SkillOptionsUseCase{repository: repository}
}

func (u *SkillOptionsUseCase) List(ctx context.Context) (*domain.SkillOptions, error) {
	return u.repository.ListSkillOptions(ctx)
}

func (u *SkillOptionsUseCase) CreateCategory(ctx context.Context, input domain.SkillCategoryInput) (*domain.SkillOption, error) {
	return u.repository.CreateSkillCategory(ctx, input)
}

func (u *SkillOptionsUseCase) UpdateCategory(ctx context.Context, id string, input domain.SkillCategoryInput) (*domain.SkillOption, error) {
	return u.repository.UpdateSkillCategory(ctx, id, input)
}

func (u *SkillOptionsUseCase) ListSkills(ctx context.Context) ([]domain.Skill, error) {
	return u.repository.ListSkills(ctx)
}

func (u *SkillOptionsUseCase) CreateSkill(ctx context.Context, input domain.SkillInput) (*domain.Skill, error) {
	return u.repository.CreateSkill(ctx, input)
}

func (u *SkillOptionsUseCase) UpdateSkill(ctx context.Context, id string, input domain.SkillInput) (*domain.Skill, error) {
	return u.repository.UpdateSkill(ctx, id, input)
}

func (u *SkillOptionsUseCase) DeleteSkill(ctx context.Context, id string) error {
	return u.repository.DeleteSkill(ctx, id)
}
