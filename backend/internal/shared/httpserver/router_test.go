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
		"zennUrl":"https://zenn.dev/sky0621",
		"qiitaUrl":"https://qiita.com/sky0621",
		"websiteUrl":"https://example.com",
		"workStyle":"Full remote",
		"visibilitySettings":{"email":false,"location":false}
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
			Location    string `json:"location"`
			Email       string `json:"email"`
			Bio         string `json:"bio"`
			GithubURL   string `json:"githubUrl"`
			ZennURL     string `json:"zennUrl"`
			QiitaURL    string `json:"qiitaUrl"`
			WebsiteURL  string `json:"websiteUrl"`
			WorkStyle   string `json:"workStyle"`
			Visibility  struct {
				Email    bool `json:"email"`
				Location bool `json:"location"`
			} `json:"visibilitySettings"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(putRec.Body.Bytes(), &putResp); err != nil {
		t.Fatalf("failed to decode put response: %v", err)
	}

	if putResp.Profile.DisplayName != "Sky Sample" {
		t.Fatalf("expected updated displayName, got %v", putResp.Profile.DisplayName)
	}
	if putResp.Profile.Location != "Tokyo" {
		t.Fatalf("expected updated location, got %v", putResp.Profile.Location)
	}
	if putResp.Profile.Email != "me@example.com" {
		t.Fatalf("expected updated email, got %v", putResp.Profile.Email)
	}
	if putResp.Profile.Bio != "Backend engineer" {
		t.Fatalf("expected updated bio, got %v", putResp.Profile.Bio)
	}
	if putResp.Profile.GithubURL != "https://github.com/sky0621" {
		t.Fatalf("expected updated githubUrl, got %v", putResp.Profile.GithubURL)
	}
	if putResp.Profile.ZennURL != "https://zenn.dev/sky0621" {
		t.Fatalf("expected updated zennUrl, got %v", putResp.Profile.ZennURL)
	}
	if putResp.Profile.QiitaURL != "https://qiita.com/sky0621" {
		t.Fatalf("expected updated qiitaUrl, got %v", putResp.Profile.QiitaURL)
	}
	if putResp.Profile.WebsiteURL != "https://example.com" {
		t.Fatalf("expected updated websiteUrl, got %v", putResp.Profile.WebsiteURL)
	}
	if putResp.Profile.WorkStyle != "Full remote" {
		t.Fatalf("expected updated workStyle, got %v", putResp.Profile.WorkStyle)
	}
	if putResp.Profile.Visibility.Email != false || putResp.Profile.Visibility.Location != false {
		t.Fatalf("expected updated visibility settings, got %+v", putResp.Profile.Visibility)
	}

	getUpdatedReq := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	getUpdatedRec := httptest.NewRecorder()
	router.ServeHTTP(getUpdatedRec, getUpdatedReq)

	var getUpdatedResp struct {
		Profile struct {
			DisplayName string `json:"displayName"`
			Location    string `json:"location"`
			Email       string `json:"email"`
			Bio         string `json:"bio"`
			GithubURL   string `json:"githubUrl"`
			ZennURL     string `json:"zennUrl"`
			QiitaURL    string `json:"qiitaUrl"`
			WebsiteURL  string `json:"websiteUrl"`
			WorkStyle   string `json:"workStyle"`
			Visibility  struct {
				Email    bool `json:"email"`
				Location bool `json:"location"`
			} `json:"visibilitySettings"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(getUpdatedRec.Body.Bytes(), &getUpdatedResp); err != nil {
		t.Fatalf("failed to decode updated get response: %v", err)
	}

	if getUpdatedResp.Profile.DisplayName != "Sky Sample" {
		t.Fatalf("expected persisted displayName, got %v", getUpdatedResp.Profile.DisplayName)
	}
	if getUpdatedResp.Profile.Location != "Tokyo" {
		t.Fatalf("expected persisted location, got %v", getUpdatedResp.Profile.Location)
	}
	if getUpdatedResp.Profile.Email != "me@example.com" {
		t.Fatalf("expected persisted email, got %v", getUpdatedResp.Profile.Email)
	}
	if getUpdatedResp.Profile.Bio != "Backend engineer" {
		t.Fatalf("expected persisted bio, got %v", getUpdatedResp.Profile.Bio)
	}
	if getUpdatedResp.Profile.GithubURL != "https://github.com/sky0621" {
		t.Fatalf("expected persisted githubUrl, got %v", getUpdatedResp.Profile.GithubURL)
	}
	if getUpdatedResp.Profile.ZennURL != "https://zenn.dev/sky0621" {
		t.Fatalf("expected persisted zennUrl, got %v", getUpdatedResp.Profile.ZennURL)
	}
	if getUpdatedResp.Profile.QiitaURL != "https://qiita.com/sky0621" {
		t.Fatalf("expected persisted qiitaUrl, got %v", getUpdatedResp.Profile.QiitaURL)
	}
	if getUpdatedResp.Profile.WebsiteURL != "https://example.com" {
		t.Fatalf("expected persisted websiteUrl, got %v", getUpdatedResp.Profile.WebsiteURL)
	}
	if getUpdatedResp.Profile.WorkStyle != "Full remote" {
		t.Fatalf("expected persisted workStyle, got %v", getUpdatedResp.Profile.WorkStyle)
	}
	if getUpdatedResp.Profile.Visibility.Email != false || getUpdatedResp.Profile.Visibility.Location != false {
		t.Fatalf("expected persisted visibility settings, got %+v", getUpdatedResp.Profile.Visibility)
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
