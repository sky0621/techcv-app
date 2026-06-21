package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/sky0621/techcv-app/backend/internal/domain"
)

func TestUpdateReplacesFieldsAndPreservesVisibilityWhenInputIsNil(t *testing.T) {
	t.Helper()

	now := time.Date(2026, 4, 19, 10, 0, 0, 0, time.UTC)
	repo := &stubRepository{
		profile: &domain.Profile{
			ID:                 "1",
			FamilyName:         "Before",
			GivenName:          "Name",
			Nickname:           "before",
			Location:           "Osaka",
			Email:              "before@example.com",
			PR:                 "old pr",
			Occupation:         "Before Occupation",
			EmploymentType:     "Before Employment",
			PreferredWorkStyle: "Hybrid",
			VisibilitySettings: map[string]any{"email": true},
			Qualifications: []domain.Qualification{
				{ID: "qualification_before", Name: "Before Cert"},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	uc := NewProfileUseCase(repo)
	got, err := uc.Update(context.Background(), ProfileInput{
		FamilyName:         "After",
		GivenName:          "Name",
		Nickname:           "after",
		Location:           "Tokyo",
		Email:              "after@example.com",
		PR:                 "new pr",
		Occupation:         "Software Engineer",
		EmploymentType:     "Freelance",
		PreferredWorkStyle: "Remote",
		Qualifications: []domain.Qualification{
			{
				ID:           "qualification_after",
				Name:         "After Cert",
				AcquiredDate: "2026-04-26",
				Organization: "Cert Org",
				URL:          "https://example.com/after-cert",
				Memo:         "memo",
			},
		},
		VisibilitySettings: nil,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if repo.getCalls != 1 {
		t.Fatalf("expected Get to be called once, got %d", repo.getCalls)
	}
	if repo.saveCalls != 1 {
		t.Fatalf("expected Save to be called once, got %d", repo.saveCalls)
	}

	if got.FamilyName != "After" || got.GivenName != "Name" || got.Nickname != "after" || got.Location != "Tokyo" {
		t.Fatalf("unexpected basic fields: %+v", got)
	}
	if got.Email != "after@example.com" || got.PR != "new pr" {
		t.Fatalf("unexpected contact fields: %+v", got)
	}
	if got.Occupation != "Software Engineer" || got.EmploymentType != "Freelance" {
		t.Fatalf("unexpected work fields: %+v", got)
	}
	if got.PreferredWorkStyle != "Remote" {
		t.Fatalf("expected PreferredWorkStyle to be updated, got %q", got.PreferredWorkStyle)
	}
	if len(got.Qualifications) != 1 || got.Qualifications[0].Name != "After Cert" {
		t.Fatalf("expected Qualifications to be updated, got %#v", got.Qualifications)
	}
	if got.Qualifications[0].URL != "https://example.com/after-cert" {
		t.Fatalf("expected Qualification URL to be updated, got %#v", got.Qualifications)
	}

	if emailVisible, ok := got.VisibilitySettings["email"]; !ok || emailVisible != true {
		t.Fatalf("expected visibility settings to be preserved, got %#v", got.VisibilitySettings)
	}
}

func TestUpdateReplacesVisibilityWhenInputIsProvided(t *testing.T) {
	repo := &stubRepository{
		profile: &domain.Profile{
			ID:                 "1",
			VisibilitySettings: map[string]any{"email": true},
		},
	}

	uc := NewProfileUseCase(repo)
	got, err := uc.Update(context.Background(), ProfileInput{
		VisibilitySettings: map[string]any{"email": false, "github": true},
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if len(got.VisibilitySettings) != 2 {
		t.Fatalf("expected replaced visibility settings, got %#v", got.VisibilitySettings)
	}
	if got.VisibilitySettings["email"] != false || got.VisibilitySettings["github"] != true {
		t.Fatalf("unexpected visibility settings: %#v", got.VisibilitySettings)
	}
}

type stubRepository struct {
	profile   *domain.Profile
	getCalls  int
	saveCalls int
}

func (r *stubRepository) Get(context.Context) (*domain.Profile, error) {
	r.getCalls++

	profileCopy := *r.profile
	if r.profile.VisibilitySettings != nil {
		profileCopy.VisibilitySettings = cloneMap(r.profile.VisibilitySettings)
	}
	if r.profile.Qualifications != nil {
		profileCopy.Qualifications = cloneQualifications(r.profile.Qualifications)
	}

	return &profileCopy, nil
}

func (r *stubRepository) Save(_ context.Context, profile *domain.Profile) (*domain.Profile, error) {
	r.saveCalls++

	profileCopy := *profile
	if profile.VisibilitySettings != nil {
		profileCopy.VisibilitySettings = cloneMap(profile.VisibilitySettings)
	}
	if profile.Qualifications != nil {
		profileCopy.Qualifications = cloneQualifications(profile.Qualifications)
	}

	r.profile = &profileCopy
	return &profileCopy, nil
}

func (r *stubRepository) ListProfileLinkMasters(context.Context) ([]domain.ProfileLinkMaster, error) {
	return nil, nil
}

func (r *stubRepository) CreateProfileLinkMaster(_ context.Context, input domain.ProfileLinkMasterInput) (*domain.ProfileLinkMaster, error) {
	return &domain.ProfileLinkMaster{
		ID:          input.ID,
		Key:         input.Key,
		Name:        input.Name,
		Icon:        input.Icon,
		Placeholder: input.Placeholder,
		SortOrder:   input.SortOrder,
	}, nil
}

func (r *stubRepository) UpdateProfileLinkMaster(_ context.Context, id string, input domain.ProfileLinkMasterInput) (*domain.ProfileLinkMaster, error) {
	return &domain.ProfileLinkMaster{
		ID:          id,
		Key:         input.Key,
		Name:        input.Name,
		Icon:        input.Icon,
		Placeholder: input.Placeholder,
		SortOrder:   input.SortOrder,
	}, nil
}

func cloneMap(values map[string]any) map[string]any {
	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}
