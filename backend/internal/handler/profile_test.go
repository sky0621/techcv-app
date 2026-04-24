package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sky0621/techcv-app/backend/internal/domain"
	"github.com/sky0621/techcv-app/backend/internal/usecase"
)

func TestGetProfileMapsDomainToResponse(t *testing.T) {
	now := time.Date(2026, 4, 19, 12, 0, 0, 0, time.UTC)
	repo := &profileRepositoryStub{
		profile: &domain.Profile{
			ID:                 "profile_01",
			UserID:             "user_01",
			FullName:           "Sky Sample",
			Email:              "me@example.com",
			VisibilitySettings: map[string]any{"email": false, "github": true, "nickname": "public"},
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	handler := NewProfileHandler(usecase.NewProfileUseCase(repo))
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	rec := httptest.NewRecorder()
	handler.GetProfile(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp profileResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Profile.ID != "profile_01" || resp.Profile.UserID != "user_01" {
		t.Fatalf("unexpected IDs: %+v", resp.Profile)
	}
	if resp.Profile.FullName == nil || *resp.Profile.FullName != "Sky Sample" {
		t.Fatalf("unexpected FullName: %+v", resp.Profile.FullName)
	}
	if resp.Profile.Email == nil || *resp.Profile.Email != "me@example.com" {
		t.Fatalf("unexpected Email: %+v", resp.Profile.Email)
	}
	if len(resp.Profile.VisibilitySettings) != 2 {
		t.Fatalf("expected only boolean visibility settings, got %#v", resp.Profile.VisibilitySettings)
	}
	if resp.Profile.VisibilitySettings["email"] != false || resp.Profile.VisibilitySettings["github"] != true {
		t.Fatalf("unexpected visibility settings: %#v", resp.Profile.VisibilitySettings)
	}
}

func TestUpdateProfileHandlesInvalidBody(t *testing.T) {
	handler := NewProfileHandler(usecase.NewProfileUseCase(&profileRepositoryStub{}))

	req := httptest.NewRequest(http.MethodPut, "/api/profile", nil)
	rec := httptest.NewRecorder()
	handler.UpdateProfile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Code != "bad_request" || resp.Message != "invalid request body" {
		t.Fatalf("unexpected bad request response: %+v", resp)
	}
}

func TestUpdateProfileMapsRequestToUseCase(t *testing.T) {
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

	handler := NewProfileHandler(usecase.NewProfileUseCase(repo))

	fullName := "Sky Sample"
	email := "me@example.com"
	visibility := map[string]bool{"email": false, "github": true}
	body, err := json.Marshal(profileUpdateRequest{
		FullName:           &fullName,
		Email:              &email,
		PreferredWorkStyle: stringRef("Full remote"),
		VisibilitySettings: &visibility,
	})
	if err != nil {
		t.Fatalf("failed to encode request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.UpdateProfile(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp profileResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Profile.FullName == nil || *resp.Profile.FullName != "Sky Sample" {
		t.Fatalf("unexpected FullName: %+v", resp.Profile.FullName)
	}
	if resp.Profile.Email == nil || *resp.Profile.Email != "me@example.com" {
		t.Fatalf("unexpected Email: %+v", resp.Profile.Email)
	}
	if resp.Profile.PreferredWorkStyle == nil || *resp.Profile.PreferredWorkStyle != "Full remote" {
		t.Fatalf("unexpected PreferredWorkStyle: %+v", resp.Profile.PreferredWorkStyle)
	}
	if resp.Profile.VisibilitySettings["email"] != false || resp.Profile.VisibilitySettings["github"] != true {
		t.Fatalf("unexpected visibility settings: %#v", resp.Profile.VisibilitySettings)
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
