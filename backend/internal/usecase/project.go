package usecase

import (
	"context"

	"github.com/sky0621/techcv-app/backend/internal/domain"
	"github.com/sky0621/techcv-app/backend/internal/repository"
)

type ProjectUseCase struct {
	repository repository.ProjectRepository
}

func NewProjectUseCase(repository repository.ProjectRepository) *ProjectUseCase {
	return &ProjectUseCase{repository: repository}
}

func (u *ProjectUseCase) List(ctx context.Context) ([]domain.Project, error) {
	return u.repository.ListProjects(ctx)
}

func (u *ProjectUseCase) Create(ctx context.Context, input domain.ProjectInput) (*domain.Project, error) {
	return u.repository.CreateProject(ctx, input)
}

func (u *ProjectUseCase) Update(ctx context.Context, id string, input domain.ProjectInput) (*domain.Project, error) {
	return u.repository.UpdateProject(ctx, id, input)
}

func (u *ProjectUseCase) Delete(ctx context.Context, id string) error {
	return u.repository.DeleteProject(ctx, id)
}

func (u *ProjectUseCase) ListOptions(ctx context.Context) (*domain.ProjectOptions, error) {
	return u.repository.ListProjectOptions(ctx)
}

func (u *ProjectUseCase) CreatePhase(ctx context.Context, input domain.ProjectPhaseInput) (*domain.ProjectPhase, error) {
	return u.repository.CreateProjectPhase(ctx, input)
}

func (u *ProjectUseCase) UpdatePhase(ctx context.Context, id string, input domain.ProjectPhaseInput) (*domain.ProjectPhase, error) {
	return u.repository.UpdateProjectPhase(ctx, id, input)
}
