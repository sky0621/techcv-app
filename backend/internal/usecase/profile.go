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
	Summary            string
	GitHubURL          string
	ZennURL            string
	QiitaURL           string
	WebsiteURL         string
	Occupation         string
	EmploymentType     string
	PreferredWorkStyle string
	VisibilitySettings map[string]any
	Qualifications     []domain.Qualification
	Links              []domain.ProfileLink
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
	current.Summary = input.Summary
	current.GitHubURL = input.GitHubURL
	current.ZennURL = input.ZennURL
	current.QiitaURL = input.QiitaURL
	current.WebsiteURL = input.WebsiteURL
	current.Occupation = input.Occupation
	current.EmploymentType = input.EmploymentType
	current.PreferredWorkStyle = input.PreferredWorkStyle
	current.Qualifications = cloneQualifications(input.Qualifications)
	current.Links = cloneProfileLinks(input.Links)
	if input.VisibilitySettings != nil {
		current.VisibilitySettings = input.VisibilitySettings
	}

	return u.repository.Save(ctx, current)
}

func (u *ProfileUseCase) ListLinkMasters(ctx context.Context) ([]domain.ProfileLinkMaster, error) {
	return u.repository.ListProfileLinkMasters(ctx)
}

func (u *ProfileUseCase) CreateLinkMaster(ctx context.Context, input domain.ProfileLinkMasterInput) (*domain.ProfileLinkMaster, error) {
	return u.repository.CreateProfileLinkMaster(ctx, input)
}

func (u *ProfileUseCase) UpdateLinkMaster(ctx context.Context, id string, input domain.ProfileLinkMasterInput) (*domain.ProfileLinkMaster, error) {
	return u.repository.UpdateProfileLinkMaster(ctx, id, input)
}

func cloneQualifications(values []domain.Qualification) []domain.Qualification {
	if values == nil {
		return nil
	}

	cloned := make([]domain.Qualification, len(values))
	copy(cloned, values)
	return cloned
}

func cloneProfileLinks(values []domain.ProfileLink) []domain.ProfileLink {
	if values == nil {
		return nil
	}

	cloned := make([]domain.ProfileLink, len(values))
	copy(cloned, values)
	return cloned
}
