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

	"github.com/sky0621/techcv-app/backend/internal/profile/domain"
	sharedopenapi "github.com/sky0621/techcv-app/backend/internal/shared/openapi"
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
		Profile sharedopenapi.Profile `json:"profile"`
	}
	if err := json.Unmarshal(getRec.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to decode get response: %v", err)
	}

	if getResp.Profile.Id != "profile_01" {
		t.Fatalf("expected profile id profile_01, got %v", getResp.Profile.Id)
	}

	body := []byte(`{
		"fullName":"Sky Sample",
		"nickname":"sky0621",
		"location":"Tokyo",
		"email":"me@example.com",
		"phone":"090-0000-0000",
		"summary":"Backend engineer",
		"githubUrl":"https://github.com/sky0621",
		"preferredWorkStyle":"Full remote",
		"visibilitySettings":{"email":false,"phone":false}
	}`)

	putReq := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(body))
	putRec := httptest.NewRecorder()
	router.ServeHTTP(putRec, putReq)

	if putRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", putRec.Code)
	}

	var putResp struct {
		Profile sharedopenapi.Profile `json:"profile"`
	}
	if err := json.Unmarshal(putRec.Body.Bytes(), &putResp); err != nil {
		t.Fatalf("failed to decode put response: %v", err)
	}

	if putResp.Profile.FullName == nil || *putResp.Profile.FullName != "Sky Sample" {
		t.Fatalf("expected updated fullName, got %v", putResp.Profile.FullName)
	}

	getUpdatedReq := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	getUpdatedRec := httptest.NewRecorder()
	router.ServeHTTP(getUpdatedRec, getUpdatedReq)

	var getUpdatedResp struct {
		Profile sharedopenapi.Profile `json:"profile"`
	}
	if err := json.Unmarshal(getUpdatedRec.Body.Bytes(), &getUpdatedResp); err != nil {
		t.Fatalf("failed to decode updated get response: %v", err)
	}

	if getUpdatedResp.Profile.FullName == nil || *getUpdatedResp.Profile.FullName != "Sky Sample" {
		t.Fatalf("expected persisted fullName, got %v", getUpdatedResp.Profile.FullName)
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
			VisibilitySettings: map[string]any{"email": false, "phone": false},
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
