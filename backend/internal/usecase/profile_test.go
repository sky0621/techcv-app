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
			ID:                 "profile_01",
			UserID:             "user_01",
			FullName:           "Before",
			Nickname:           "before",
			Location:           "Osaka",
			Email:              "before@example.com",
			Phone:              "000",
			Summary:            "old summary",
			GitHubURL:          "https://github.com/before",
			ZennURL:            "https://zenn.dev/before",
			QiitaURL:           "https://qiita.com/before",
			WebsiteURL:         "https://before.example.com",
			PreferredWorkStyle: "Hybrid",
			VisibilitySettings: map[string]any{"email": true, "phone": false},
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	uc := NewProfileUseCase(repo)
	got, err := uc.Update(context.Background(), ProfileInput{
		FullName:           "After",
		Nickname:           "after",
		Location:           "Tokyo",
		Email:              "after@example.com",
		Phone:              "111",
		Summary:            "new summary",
		GitHubURL:          "https://github.com/after",
		ZennURL:            "https://zenn.dev/after",
		QiitaURL:           "https://qiita.com/after",
		WebsiteURL:         "https://after.example.com",
		PreferredWorkStyle: "Remote",
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

	if got.FullName != "After" || got.Nickname != "after" || got.Location != "Tokyo" {
		t.Fatalf("unexpected basic fields: %+v", got)
	}
	if got.Email != "after@example.com" || got.Phone != "111" || got.Summary != "new summary" {
		t.Fatalf("unexpected contact fields: %+v", got)
	}
	if got.GitHubURL != "https://github.com/after" || got.ZennURL != "https://zenn.dev/after" {
		t.Fatalf("unexpected social fields: %+v", got)
	}
	if got.QiitaURL != "https://qiita.com/after" || got.WebsiteURL != "https://after.example.com" {
		t.Fatalf("unexpected site fields: %+v", got)
	}
	if got.PreferredWorkStyle != "Remote" {
		t.Fatalf("expected PreferredWorkStyle to be updated, got %q", got.PreferredWorkStyle)
	}

	if emailVisible, ok := got.VisibilitySettings["email"]; !ok || emailVisible != true {
		t.Fatalf("expected visibility settings to be preserved, got %#v", got.VisibilitySettings)
	}
}

func TestUpdateReplacesVisibilityWhenInputIsProvided(t *testing.T) {
	repo := &stubRepository{
		profile: &domain.Profile{
			ID:                 "profile_01",
			UserID:             "user_01",
			VisibilitySettings: map[string]any{"email": true, "phone": false},
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

	return &profileCopy, nil
}

func (r *stubRepository) Save(_ context.Context, profile *domain.Profile) (*domain.Profile, error) {
	r.saveCalls++

	profileCopy := *profile
	if profile.VisibilitySettings != nil {
		profileCopy.VisibilitySettings = cloneMap(profile.VisibilitySettings)
	}

	r.profile = &profileCopy
	return &profileCopy, nil
}

func cloneMap(values map[string]any) map[string]any {
	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}
