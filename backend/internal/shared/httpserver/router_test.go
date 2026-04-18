package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	profilerepository "github.com/sky0621/techcv-app/backend/internal/profile/repository"
)

func TestProfileRoutes(t *testing.T) {
	router := NewRouter(profilerepository.NewMemoryRepository())

	getReq := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", getRec.Code)
	}

	var getResp struct {
		Profile Profile `json:"profile"`
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
		Profile Profile `json:"profile"`
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
		Profile Profile `json:"profile"`
	}
	if err := json.Unmarshal(getUpdatedRec.Body.Bytes(), &getUpdatedResp); err != nil {
		t.Fatalf("failed to decode updated get response: %v", err)
	}

	if getUpdatedResp.Profile.FullName == nil || *getUpdatedResp.Profile.FullName != "Sky Sample" {
		t.Fatalf("expected persisted fullName, got %v", getUpdatedResp.Profile.FullName)
	}
}
