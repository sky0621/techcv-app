package usecase

import (
	"context"

	"github.com/sky0621/techcv-app/backend/internal/domain"
	"github.com/sky0621/techcv-app/backend/internal/repository"
)

type ProfileInput struct {
	FullName           string
	Nickname           string
	Location           string
	Email              string
	Phone              string
	Summary            string
	GitHubURL          string
	ZennURL            string
	QiitaURL           string
	WebsiteURL         string
	PreferredWorkStyle string
	VisibilitySettings map[string]any
}

type ProfileUseCase struct {
	repository repository.ProfileRepository
}

func NewProfileUseCase(repository repository.ProfileRepository) *ProfileUseCase {
	return &ProfileUseCase{repository: repository}
}

func (u *ProfileUseCase) Get(ctx context.Context) (*domain.Profile, error) {
	return u.repository.Get(ctx)
}

func (u *ProfileUseCase) Update(ctx context.Context, input ProfileInput) (*domain.Profile, error) {
	current, err := u.repository.Get(ctx)
	if err != nil {
		return nil, err
	}

	current.FullName = input.FullName
	current.Nickname = input.Nickname
	current.Location = input.Location
	current.Email = input.Email
	current.Phone = input.Phone
	current.Summary = input.Summary
	current.GitHubURL = input.GitHubURL
	current.ZennURL = input.ZennURL
	current.QiitaURL = input.QiitaURL
	current.WebsiteURL = input.WebsiteURL
	current.PreferredWorkStyle = input.PreferredWorkStyle
	if input.VisibilitySettings != nil {
		current.VisibilitySettings = input.VisibilitySettings
	}

	return u.repository.Save(ctx, current)
}
