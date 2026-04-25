package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sky0621/techcv-app/backend/internal/domain"
	"github.com/sky0621/techcv-app/backend/internal/usecase"
)

type ProfileHandler struct {
	usecase *usecase.ProfileUseCase
}

func NewProfileHandler(usecase *usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{usecase: usecase}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.usecase.Get(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, profileResponse{
		Profile: toProfileResponse(profile),
	})
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var request profileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	profile, err := h.usecase.Update(r.Context(), usecase.ProfileInput{
		FullName:           stringValue(request.DisplayName),
		Location:           stringValue(request.Location),
		Email:              stringValue(request.Email),
		Summary:            stringValue(request.Bio),
		GitHubURL:          stringValue(request.GithubUrl),
		ZennURL:            stringValue(request.ZennUrl),
		QiitaURL:           stringValue(request.QiitaUrl),
		WebsiteURL:         stringValue(request.WebsiteUrl),
		PreferredWorkStyle: stringValue(request.WorkStyle),
		VisibilitySettings: toUseCaseVisibilitySettings(request.VisibilitySettings),
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, profileResponse{
		Profile: toProfileResponse(profile),
	})
}

type profileResponse struct {
	Profile profilePayload `json:"profile"`
}

type profilePayload struct {
	DisplayName        *string         `json:"displayName,omitempty"`
	Location           *string         `json:"location,omitempty"`
	Email              *string         `json:"email,omitempty"`
	Bio                *string         `json:"bio,omitempty"`
	GithubUrl          *string         `json:"githubUrl,omitempty"`
	ZennUrl            *string         `json:"zennUrl,omitempty"`
	QiitaUrl           *string         `json:"qiitaUrl,omitempty"`
	WebsiteUrl         *string         `json:"websiteUrl,omitempty"`
	WorkStyle          *string         `json:"workStyle,omitempty"`
	VisibilitySettings map[string]bool `json:"visibilitySettings"`
}

type profileUpdateRequest struct {
	DisplayName        *string          `json:"displayName,omitempty"`
	Location           *string          `json:"location,omitempty"`
	Email              *string          `json:"email,omitempty"`
	Bio                *string          `json:"bio,omitempty"`
	GithubUrl          *string          `json:"githubUrl,omitempty"`
	ZennUrl            *string          `json:"zennUrl,omitempty"`
	QiitaUrl           *string          `json:"qiitaUrl,omitempty"`
	WebsiteUrl         *string          `json:"websiteUrl,omitempty"`
	WorkStyle          *string          `json:"workStyle,omitempty"`
	VisibilitySettings *map[string]bool `json:"visibilitySettings,omitempty"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func toProfileResponse(profile *domain.Profile) profilePayload {
	return profilePayload{
		DisplayName:        stringPointer(profile.FullName),
		Location:           stringPointer(profile.Location),
		Email:              stringPointer(profile.Email),
		Bio:                stringPointer(profile.Summary),
		GithubUrl:          stringPointer(profile.GitHubURL),
		ZennUrl:            stringPointer(profile.ZennURL),
		QiitaUrl:           stringPointer(profile.QiitaURL),
		WebsiteUrl:         stringPointer(profile.WebsiteURL),
		WorkStyle:          stringPointer(profile.PreferredWorkStyle),
		VisibilitySettings: toVisibilitySettingsResponse(profile.VisibilitySettings),
	}
}

func toVisibilitySettingsResponse(values map[string]any) map[string]bool {
	result := map[string]bool{}
	for key, value := range values {
		booleanValue, ok := value.(bool)
		if ok {
			result[key] = booleanValue
		}
	}

	return result
}

func toUseCaseVisibilitySettings(values *map[string]bool) map[string]any {
	if values == nil {
		return nil
	}

	result := make(map[string]any, len(*values))
	for key, value := range *values {
		result[key] = value
	}

	return result
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func stringPointer(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{
		Code:    code,
		Message: message,
	})
}
