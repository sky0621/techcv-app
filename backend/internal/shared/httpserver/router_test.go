package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/sky0621/techcv-app/backend/internal/domain"
)

func TestProfileRoutes(t *testing.T) {
	router := NewRouter(newTestProfileRepository())

	getReq := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", getRec.Code)
	}

	var getResp struct {
		Profile struct {
			DisplayName string `json:"displayName"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(getRec.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to decode get response: %v", err)
	}

	if getResp.Profile.DisplayName != "Initial User" {
		t.Fatalf("expected initial displayName, got %v", getResp.Profile.DisplayName)
	}

	body := []byte(`{
		"displayName":"Sky Sample",
		"location":"Tokyo",
		"email":"me@example.com",
		"bio":"Backend engineer",
		"githubUrl":"https://github.com/sky0621",
		"workStyle":"Full remote",
		"visibilitySettings":{"email":false}
	}`)

	putReq := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(body))
	putRec := httptest.NewRecorder()
	router.ServeHTTP(putRec, putReq)

	if putRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", putRec.Code)
	}

	var putResp struct {
		Profile struct {
			DisplayName string `json:"displayName"`
			Bio         string `json:"bio"`
			WorkStyle   string `json:"workStyle"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(putRec.Body.Bytes(), &putResp); err != nil {
		t.Fatalf("failed to decode put response: %v", err)
	}

	if putResp.Profile.DisplayName != "Sky Sample" {
		t.Fatalf("expected updated displayName, got %v", putResp.Profile.DisplayName)
	}
	if putResp.Profile.Bio != "Backend engineer" {
		t.Fatalf("expected updated bio, got %v", putResp.Profile.Bio)
	}
	if putResp.Profile.WorkStyle != "Full remote" {
		t.Fatalf("expected updated workStyle, got %v", putResp.Profile.WorkStyle)
	}

	getUpdatedReq := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	getUpdatedRec := httptest.NewRecorder()
	router.ServeHTTP(getUpdatedRec, getUpdatedReq)

	var getUpdatedResp struct {
		Profile struct {
			DisplayName string `json:"displayName"`
			Bio         string `json:"bio"`
			WorkStyle   string `json:"workStyle"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(getUpdatedRec.Body.Bytes(), &getUpdatedResp); err != nil {
		t.Fatalf("failed to decode updated get response: %v", err)
	}

	if getUpdatedResp.Profile.DisplayName != "Sky Sample" {
		t.Fatalf("expected persisted displayName, got %v", getUpdatedResp.Profile.DisplayName)
	}
	if getUpdatedResp.Profile.Bio != "Backend engineer" {
		t.Fatalf("expected persisted bio, got %v", getUpdatedResp.Profile.Bio)
	}
	if getUpdatedResp.Profile.WorkStyle != "Full remote" {
		t.Fatalf("expected persisted workStyle, got %v", getUpdatedResp.Profile.WorkStyle)
	}
}

type testProfileRepository struct {
	mu      sync.RWMutex
	profile *domain.Profile
}

func newTestProfileRepository() *testProfileRepository {
	now := time.Now().UTC()

	return &testProfileRepository{
		profile: &domain.Profile{
			ID:                 "profile_01",
			UserID:             "user_01",
			FullName:           "Initial User",
			VisibilitySettings: map[string]any{"email": false},
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}
}

func (r *testProfileRepository) Get(_ context.Context) (*domain.Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	profileCopy := *r.profile
	return &profileCopy, nil
}

func (r *testProfileRepository) Save(_ context.Context, profile *domain.Profile) (*domain.Profile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	profileCopy := *profile
	profileCopy.UpdatedAt = time.Now().UTC()
	if profileCopy.CreatedAt.IsZero() {
		profileCopy.CreatedAt = profileCopy.UpdatedAt
	}

	r.profile = &profileCopy

	savedCopy := *r.profile
	return &savedCopy, nil
}
