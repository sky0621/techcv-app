package repository

import (
	"context"

	"github.com/sky0621/techcv-app/backend/internal/domain"
)

type JobHistoryRepository interface {
	ListJobHistories(ctx context.Context) ([]domain.JobHistory, error)
	CreateJobHistory(ctx context.Context, input domain.JobHistoryInput) (*domain.JobHistory, error)
	UpdateJobHistory(ctx context.Context, id string, input domain.JobHistoryInput) (*domain.JobHistory, error)
	DeleteJobHistory(ctx context.Context, id string) error
	ListJobHistoryOptions(ctx context.Context) (*domain.JobHistoryOptions, error)
	CreateJobEmploymentType(ctx context.Context, input domain.JobEmploymentTypeInput) (*domain.JobEmploymentType, error)
	UpdateJobEmploymentType(ctx context.Context, id string, input domain.JobEmploymentTypeInput) (*domain.JobEmploymentType, error)
}
