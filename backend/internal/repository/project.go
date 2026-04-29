package repository

import (
	"context"

	"github.com/sky0621/techcv-app/backend/internal/domain"
)

type ProjectRepository interface {
	ListProjects(ctx context.Context) ([]domain.Project, error)
	CreateProject(ctx context.Context, input domain.ProjectInput) (*domain.Project, error)
	UpdateProject(ctx context.Context, id string, input domain.ProjectInput) (*domain.Project, error)
	DeleteProject(ctx context.Context, id string) error
}
