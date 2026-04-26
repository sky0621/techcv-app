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
		Occupation:         stringValue(request.Occupation),
		EmploymentType:     stringValue(request.EmploymentType),
		PreferredWorkStyle: stringValue(request.WorkStyle),
		VisibilitySettings: toUseCaseVisibilitySettings(request.VisibilitySettings),
		Qualifications:     toUseCaseQualifications(request.Qualifications),
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
	DisplayName        string                 `json:"displayName"`
	Location           string                 `json:"location"`
	Email              string                 `json:"email"`
	Bio                string                 `json:"bio"`
	GithubUrl          string                 `json:"githubUrl"`
	ZennUrl            string                 `json:"zennUrl"`
	QiitaUrl           string                 `json:"qiitaUrl"`
	WebsiteUrl         string                 `json:"websiteUrl"`
	Occupation         string                 `json:"occupation"`
	EmploymentType     string                 `json:"employmentType"`
	WorkStyle          string                 `json:"workStyle"`
	VisibilitySettings map[string]bool        `json:"visibilitySettings"`
	Qualifications     []qualificationPayload `json:"qualifications"`
}

type profileUpdateRequest struct {
	DisplayName        *string                `json:"displayName,omitempty"`
	Location           *string                `json:"location,omitempty"`
	Email              *string                `json:"email,omitempty"`
	Bio                *string                `json:"bio,omitempty"`
	GithubUrl          *string                `json:"githubUrl,omitempty"`
	ZennUrl            *string                `json:"zennUrl,omitempty"`
	QiitaUrl           *string                `json:"qiitaUrl,omitempty"`
	WebsiteUrl         *string                `json:"websiteUrl,omitempty"`
	Occupation         *string                `json:"occupation,omitempty"`
	EmploymentType     *string                `json:"employmentType,omitempty"`
	WorkStyle          *string                `json:"workStyle,omitempty"`
	VisibilitySettings *map[string]bool       `json:"visibilitySettings,omitempty"`
	Qualifications     []qualificationPayload `json:"qualifications,omitempty"`
}

type qualificationPayload struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AcquiredDate string `json:"acquiredDate"`
	Organization string `json:"organization"`
	URL          string `json:"url"`
	Memo         string `json:"memo"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func toProfileResponse(profile *domain.Profile) profilePayload {
	return profilePayload{
		DisplayName:        profile.FullName,
		Location:           profile.Location,
		Email:              profile.Email,
		Bio:                profile.Summary,
		GithubUrl:          profile.GitHubURL,
		ZennUrl:            profile.ZennURL,
		QiitaUrl:           profile.QiitaURL,
		WebsiteUrl:         profile.WebsiteURL,
		Occupation:         profile.Occupation,
		EmploymentType:     profile.EmploymentType,
		WorkStyle:          profile.PreferredWorkStyle,
		VisibilitySettings: toVisibilitySettingsResponse(profile.VisibilitySettings),
		Qualifications:     toQualificationPayloads(profile.Qualifications),
	}
}

func toQualificationPayloads(values []domain.Qualification) []qualificationPayload {
	result := make([]qualificationPayload, 0, len(values))
	for _, value := range values {
		result = append(result, qualificationPayload{
			ID:           value.ID,
			Name:         value.Name,
			AcquiredDate: value.AcquiredDate,
			Organization: value.Organization,
			URL:          value.URL,
			Memo:         value.Memo,
		})
	}

	return result
}

func toUseCaseQualifications(values []qualificationPayload) []domain.Qualification {
	result := make([]domain.Qualification, 0, len(values))
	for _, value := range values {
		result = append(result, domain.Qualification{
			ID:           value.ID,
			Name:         value.Name,
			AcquiredDate: value.AcquiredDate,
			Organization: value.Organization,
			URL:          value.URL,
			Memo:         value.Memo,
		})
	}

	return result
}

func toVisibilitySettingsResponse(values map[string]any) map[string]bool {
	result := map[string]bool{
		"email":    true,
		"location": true,
	}
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
