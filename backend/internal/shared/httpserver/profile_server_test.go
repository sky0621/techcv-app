package httpserver

import (
	"context"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sky0621/techcv-app/backend/internal/profile/domain"
	"github.com/sky0621/techcv-app/backend/internal/profile/usecase"
)

func TestProfileServerGetProfileMapsDomainToOpenAPI(t *testing.T) {
	now := time.Date(2026, 4, 19, 12, 0, 0, 0, time.UTC)
	repo := &profileRepositoryStub{
		profile: &domain.Profile{
			ID:                 "profile_01",
			UserID:             "user_01",
			FullName:           "Sky Sample",
			Email:              "me@example.com",
			VisibilitySettings: map[string]any{"email": false, "phone": true, "nickname": "public"},
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	server := NewProfileServer(usecase.New(repo))
	resp, err := server.GetProfile(context.Background(), GetProfileRequestObject{})
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}

	okResp, ok := resp.(GetProfile200JSONResponse)
	if !ok {
		t.Fatalf("expected GetProfile200JSONResponse, got %T", resp)
	}

	if okResp.Profile.Id != "profile_01" || okResp.Profile.UserId != "user_01" {
		t.Fatalf("unexpected IDs: %+v", okResp.Profile)
	}
	if okResp.Profile.FullName == nil || *okResp.Profile.FullName != "Sky Sample" {
		t.Fatalf("unexpected FullName: %+v", okResp.Profile.FullName)
	}
	if okResp.Profile.Email == nil || string(*okResp.Profile.Email) != "me@example.com" {
		t.Fatalf("unexpected Email: %+v", okResp.Profile.Email)
	}
	if len(okResp.Profile.VisibilitySettings) != 2 {
		t.Fatalf("expected only boolean visibility settings, got %#v", okResp.Profile.VisibilitySettings)
	}
	if okResp.Profile.VisibilitySettings["email"] != false || okResp.Profile.VisibilitySettings["phone"] != true {
		t.Fatalf("unexpected visibility settings: %#v", okResp.Profile.VisibilitySettings)
	}
}

func TestProfileServerUpdateProfileHandlesNilBody(t *testing.T) {
	server := NewProfileServer(usecase.New(&profileRepositoryStub{}))

	resp, err := server.UpdateProfile(context.Background(), UpdateProfileRequestObject{})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}

	badReqResp, ok := resp.(UpdateProfile400JSONResponse)
	if !ok {
		t.Fatalf("expected UpdateProfile400JSONResponse, got %T", resp)
	}

	if badReqResp.Code != "bad_request" || badReqResp.Message != "invalid request body" {
		t.Fatalf("unexpected bad request response: %+v", badReqResp)
	}
}

func TestProfileServerUpdateProfileMapsOpenAPIInputToUseCase(t *testing.T) {
	now := time.Date(2026, 4, 19, 12, 30, 0, 0, time.UTC)
	repo := &profileRepositoryStub{
		profile: &domain.Profile{
			ID:                 "profile_01",
			UserID:             "user_01",
			VisibilitySettings: map[string]any{"email": true},
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	server := NewProfileServer(usecase.New(repo))

	fullName := "Sky Sample"
	email := openapi_types.Email("me@example.com")
	visibility := VisibilitySettings{"email": false, "phone": true}
	resp, err := server.UpdateProfile(context.Background(), UpdateProfileRequestObject{
		Body: &UpdateProfileJSONRequestBody{
			FullName:           &fullName,
			Email:              &email,
			PreferredWorkStyle: stringRef("Full remote"),
			VisibilitySettings: &visibility,
		},
	})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}

	okResp, ok := resp.(UpdateProfile200JSONResponse)
	if !ok {
		t.Fatalf("expected UpdateProfile200JSONResponse, got %T", resp)
	}

	if okResp.Profile.FullName == nil || *okResp.Profile.FullName != "Sky Sample" {
		t.Fatalf("unexpected FullName: %+v", okResp.Profile.FullName)
	}
	if okResp.Profile.Email == nil || string(*okResp.Profile.Email) != "me@example.com" {
		t.Fatalf("unexpected Email: %+v", okResp.Profile.Email)
	}
	if okResp.Profile.PreferredWorkStyle == nil || *okResp.Profile.PreferredWorkStyle != "Full remote" {
		t.Fatalf("unexpected PreferredWorkStyle: %+v", okResp.Profile.PreferredWorkStyle)
	}
	if okResp.Profile.VisibilitySettings["email"] != false || okResp.Profile.VisibilitySettings["phone"] != true {
		t.Fatalf("unexpected visibility settings: %#v", okResp.Profile.VisibilitySettings)
	}
}

type profileRepositoryStub struct {
	profile *domain.Profile
}

func (r *profileRepositoryStub) Get(context.Context) (*domain.Profile, error) {
	if r.profile == nil {
		r.profile = &domain.Profile{}
	}

	profileCopy := *r.profile
	if r.profile.VisibilitySettings != nil {
		profileCopy.VisibilitySettings = cloneVisibilitySettings(r.profile.VisibilitySettings)
	}

	return &profileCopy, nil
}

func (r *profileRepositoryStub) Save(_ context.Context, profile *domain.Profile) (*domain.Profile, error) {
	profileCopy := *profile
	if profile.VisibilitySettings != nil {
		profileCopy.VisibilitySettings = cloneVisibilitySettings(profile.VisibilitySettings)
	}
	if profileCopy.UpdatedAt.IsZero() {
		profileCopy.UpdatedAt = time.Now().UTC()
	}

	r.profile = &profileCopy
	return &profileCopy, nil
}

func cloneVisibilitySettings(values map[string]any) map[string]any {
	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}

func stringRef(value string) *string {
	return &value
}
