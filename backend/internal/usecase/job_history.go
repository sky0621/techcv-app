package usecase

import (
	"context"

	"github.com/sky0621/techcv-app/backend/internal/domain"
	"github.com/sky0621/techcv-app/backend/internal/repository"
)

type JobHistoryUseCase struct {
	repository repository.JobHistoryRepository
}

func NewJobHistoryUseCase(repository repository.JobHistoryRepository) *JobHistoryUseCase {
	return &JobHistoryUseCase{repository: repository}
}

func (u *JobHistoryUseCase) List(ctx context.Context) ([]domain.JobHistory, error) {
	return u.repository.ListJobHistories(ctx)
}

func (u *JobHistoryUseCase) Create(ctx context.Context, input domain.JobHistoryInput) (*domain.JobHistory, error) {
	return u.repository.CreateJobHistory(ctx, input)
}

func (u *JobHistoryUseCase) Update(ctx context.Context, id string, input domain.JobHistoryInput) (*domain.JobHistory, error) {
	return u.repository.UpdateJobHistory(ctx, id, input)
}

func (u *JobHistoryUseCase) Delete(ctx context.Context, id string) error {
	return u.repository.DeleteJobHistory(ctx, id)
}
