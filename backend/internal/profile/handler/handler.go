package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sky0621/techcv-app/backend/internal/profile/domain"
	"github.com/sky0621/techcv-app/backend/internal/profile/usecase"
)

type Handler struct {
	usecase *usecase.UseCase
}

type updateProfileRequest struct {
	FullName           string         `json:"fullName"`
	Nickname           string         `json:"nickname"`
	Location           string         `json:"location"`
	Email              string         `json:"email"`
	Phone              string         `json:"phone"`
	Summary            string         `json:"summary"`
	GitHubURL          string         `json:"githubUrl"`
	ZennURL            string         `json:"zennUrl"`
	QiitaURL           string         `json:"qiitaUrl"`
	WebsiteURL         string         `json:"websiteUrl"`
	PreferredWorkStyle string         `json:"preferredWorkStyle"`
	VisibilitySettings map[string]any `json:"visibilitySettings"`
}

type profileResponse struct {
	ID                 string         `json:"id"`
	UserID             string         `json:"userId"`
	FullName           string         `json:"fullName"`
	Nickname           string         `json:"nickname"`
	Location           string         `json:"location"`
	Email              string         `json:"email"`
	Phone              string         `json:"phone"`
	Summary            string         `json:"summary"`
	GitHubURL          string         `json:"githubUrl"`
	ZennURL            string         `json:"zennUrl"`
	QiitaURL           string         `json:"qiitaUrl"`
	WebsiteURL         string         `json:"websiteUrl"`
	PreferredWorkStyle string         `json:"preferredWorkStyle"`
	VisibilitySettings map[string]any `json:"visibilitySettings"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
}

type getProfileResponse struct {
	Profile profileResponse `json:"profile"`
}

type updateProfileResponse struct {
	Profile profileResponse `json:"profile"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func New(usecase *usecase.UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.usecase.Get(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_server_error", "failed to load profile")
		return
	}

	writeJSON(w, http.StatusOK, getProfileResponse{
		Profile: newProfileResponse(profile),
	})
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var request updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	profile, err := h.usecase.Update(r.Context(), usecase.ProfileInput{
		FullName:           request.FullName,
		Nickname:           request.Nickname,
		Location:           request.Location,
		Email:              request.Email,
		Phone:              request.Phone,
		Summary:            request.Summary,
		GitHubURL:          request.GitHubURL,
		ZennURL:            request.ZennURL,
		QiitaURL:           request.QiitaURL,
		WebsiteURL:         request.WebsiteURL,
		PreferredWorkStyle: request.PreferredWorkStyle,
		VisibilitySettings: request.VisibilitySettings,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_server_error", "failed to update profile")
		return
	}

	writeJSON(w, http.StatusOK, updateProfileResponse{
		Profile: newProfileResponse(profile),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{
		Code:    code,
		Message: message,
	})
}

func newProfileResponse(profile *domain.Profile) profileResponse {
	return profileResponse{
		ID:                 profile.ID,
		UserID:             profile.UserID,
		FullName:           profile.FullName,
		Nickname:           profile.Nickname,
		Location:           profile.Location,
		Email:              profile.Email,
		Phone:              profile.Phone,
		Summary:            profile.Summary,
		GitHubURL:          profile.GitHubURL,
		ZennURL:            profile.ZennURL,
		QiitaURL:           profile.QiitaURL,
		WebsiteURL:         profile.WebsiteURL,
		PreferredWorkStyle: profile.PreferredWorkStyle,
		VisibilitySettings: profile.VisibilitySettings,
		CreatedAt:          profile.CreatedAt,
		UpdatedAt:          profile.UpdatedAt,
	}
}
