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
			ID:                 "1",
			UserID:             "user_01",
			FullName:           "Sky Sample",
			Email:              "me@example.com",
			Summary:            "Backend engineer",
			Occupation:         "Software Engineer",
			EmploymentType:     "Freelance",
			PreferredWorkStyle: "Full remote",
			VisibilitySettings: map[string]any{"email": false, "github": true, "nickname": "public"},
			Qualifications: []domain.Qualification{
				{
					ID:           "qualification_01",
					Name:         "AWS Certified Solutions Architect",
					AcquiredDate: "2026-04-26",
					Organization: "Amazon Web Services",
					URL:          "https://aws.amazon.com/certification/certified-solutions-architect-associate/",
					Memo:         "Associate",
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
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

	if resp.Profile.DisplayName != "Sky Sample" {
		t.Fatalf("unexpected DisplayName: %+v", resp.Profile.DisplayName)
	}
	if resp.Profile.Email != "me@example.com" {
		t.Fatalf("unexpected Email: %+v", resp.Profile.Email)
	}
	if resp.Profile.Bio != "Backend engineer" {
		t.Fatalf("unexpected Bio: %+v", resp.Profile.Bio)
	}
	if resp.Profile.Occupation != "Software Engineer" {
		t.Fatalf("unexpected Occupation: %+v", resp.Profile.Occupation)
	}
	if resp.Profile.EmploymentType != "Freelance" {
		t.Fatalf("unexpected EmploymentType: %+v", resp.Profile.EmploymentType)
	}
	if resp.Profile.WorkStyle != "Full remote" {
		t.Fatalf("unexpected WorkStyle: %+v", resp.Profile.WorkStyle)
	}
	if len(resp.Profile.Qualifications) != 1 {
		t.Fatalf("expected one qualification, got %#v", resp.Profile.Qualifications)
	}
	if resp.Profile.Qualifications[0].Name != "AWS Certified Solutions Architect" {
		t.Fatalf("unexpected qualification: %#v", resp.Profile.Qualifications[0])
	}
	if resp.Profile.Qualifications[0].URL != "https://aws.amazon.com/certification/certified-solutions-architect-associate/" {
		t.Fatalf("unexpected qualification URL: %#v", resp.Profile.Qualifications[0])
	}
	if len(resp.Profile.VisibilitySettings) != 3 {
		t.Fatalf("expected boolean visibility settings with defaults, got %#v", resp.Profile.VisibilitySettings)
	}
	if resp.Profile.VisibilitySettings["email"] != false || resp.Profile.VisibilitySettings["github"] != true {
		t.Fatalf("unexpected visibility settings: %#v", resp.Profile.VisibilitySettings)
	}
	if resp.Profile.VisibilitySettings["location"] != true {
		t.Fatalf("expected default location visibility, got %#v", resp.Profile.VisibilitySettings)
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
			ID:                 "1",
			UserID:             "user_01",
			VisibilitySettings: map[string]any{"email": true},
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	handler := NewProfileHandler(usecase.NewProfileUseCase(repo))

	displayName := "Sky Sample"
	email := "me@example.com"
	bio := "Backend engineer"
	occupation := "Software Engineer"
	employmentType := "Freelance"
	workStyle := "Full remote"
	qualifications := []qualificationPayload{
		{
			Name:         "AWS Certified Solutions Architect",
			AcquiredDate: "2026-04-26",
			Organization: "Amazon Web Services",
			URL:          "https://aws.amazon.com/certification/certified-solutions-architect-associate/",
			Memo:         "Associate",
		},
	}
	visibility := map[string]bool{"email": false, "location": false}
	body, err := json.Marshal(profileUpdateRequest{
		DisplayName:        &displayName,
		Email:              &email,
		Bio:                &bio,
		Occupation:         &occupation,
		EmploymentType:     &employmentType,
		WorkStyle:          &workStyle,
		Qualifications:     qualifications,
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

	if resp.Profile.DisplayName != "Sky Sample" {
		t.Fatalf("unexpected DisplayName: %+v", resp.Profile.DisplayName)
	}
	if resp.Profile.Email != "me@example.com" {
		t.Fatalf("unexpected Email: %+v", resp.Profile.Email)
	}
	if resp.Profile.Bio != "Backend engineer" {
		t.Fatalf("unexpected Bio: %+v", resp.Profile.Bio)
	}
	if resp.Profile.Occupation != "Software Engineer" {
		t.Fatalf("unexpected Occupation: %+v", resp.Profile.Occupation)
	}
	if resp.Profile.EmploymentType != "Freelance" {
		t.Fatalf("unexpected EmploymentType: %+v", resp.Profile.EmploymentType)
	}
	if resp.Profile.WorkStyle != "Full remote" {
		t.Fatalf("unexpected WorkStyle: %+v", resp.Profile.WorkStyle)
	}
	if len(resp.Profile.Qualifications) != 1 {
		t.Fatalf("expected one qualification, got %#v", resp.Profile.Qualifications)
	}
	if resp.Profile.Qualifications[0].Name != "AWS Certified Solutions Architect" {
		t.Fatalf("unexpected qualification: %#v", resp.Profile.Qualifications[0])
	}
	if resp.Profile.Qualifications[0].URL != "https://aws.amazon.com/certification/certified-solutions-architect-associate/" {
		t.Fatalf("unexpected qualification URL: %#v", resp.Profile.Qualifications[0])
	}
	if resp.Profile.VisibilitySettings["email"] != false || resp.Profile.VisibilitySettings["location"] != false {
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
	if r.profile.Qualifications != nil {
		profileCopy.Qualifications = cloneQualifications(r.profile.Qualifications)
	}

	return &profileCopy, nil
}

func (r *profileRepositoryStub) Save(_ context.Context, profile *domain.Profile) (*domain.Profile, error) {
	profileCopy := *profile
	if profile.VisibilitySettings != nil {
		profileCopy.VisibilitySettings = cloneVisibilitySettings(profile.VisibilitySettings)
	}
	if profile.Qualifications != nil {
		profileCopy.Qualifications = cloneQualifications(profile.Qualifications)
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

func cloneQualifications(values []domain.Qualification) []domain.Qualification {
	cloned := make([]domain.Qualification, len(values))
	copy(cloned, values)
	return cloned
}

func stringRef(value string) *string {
	return &value
}
