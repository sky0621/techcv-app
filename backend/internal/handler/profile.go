package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
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
		Links:              toUseCaseProfileLinks(request.Links),
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, profileResponse{
		Profile: toProfileResponse(profile),
	})
}

func (h *ProfileHandler) ListProfileLinkMasters(w http.ResponseWriter, r *http.Request) {
	masters, err := h.usecase.ListLinkMasters(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, profileLinkMastersResponse{
		LinkMasters: toProfileLinkMasterPayloads(masters),
	})
}

func (h *ProfileHandler) CreateProfileLinkMaster(w http.ResponseWriter, r *http.Request) {
	var request profileLinkMasterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toProfileLinkMasterInput(request, true)
	if input.ID == "" || input.Key == "" || input.Name == "" || input.Icon == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id, key, name and icon are required")
		return
	}

	master, err := h.usecase.CreateLinkMaster(r.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, profileLinkMasterResponse{
		LinkMaster: toProfileLinkMasterPayload(*master),
	})
}

func (h *ProfileHandler) UpdateProfileLinkMaster(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	var request profileLinkMasterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toProfileLinkMasterInput(request, false)
	if input.Key == "" || input.Name == "" || input.Icon == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "key, name and icon are required")
		return
	}

	master, err := h.usecase.UpdateLinkMaster(r.Context(), id, input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, profileLinkMasterResponse{
		LinkMaster: toProfileLinkMasterPayload(*master),
	})
}

type profileResponse struct {
	Profile profilePayload `json:"profile"`
}

type profileLinkMastersResponse struct {
	LinkMasters []profileLinkMasterPayload `json:"linkMasters"`
}

type profileLinkMasterResponse struct {
	LinkMaster profileLinkMasterPayload `json:"linkMaster"`
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
	Links              []profileLinkPayload   `json:"links"`
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
	Links              []profileLinkPayload   `json:"links,omitempty"`
}

type qualificationPayload struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AcquiredDate string `json:"acquiredDate"`
	Organization string `json:"organization"`
	URL          string `json:"url"`
	Memo         string `json:"memo"`
}

type profileLinkPayload struct {
	ID           string `json:"id"`
	LinkMasterID string `json:"linkMasterId"`
	Key          string `json:"key"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	Placeholder  string `json:"placeholder"`
	URL          string `json:"url"`
	SortOrder    int64  `json:"sortOrder"`
}

type profileLinkMasterPayload struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Placeholder string `json:"placeholder"`
	SortOrder   int64  `json:"sortOrder"`
}

type profileLinkMasterRequest struct {
	ID          string `json:"id,omitempty"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Placeholder string `json:"placeholder"`
	SortOrder   int64  `json:"sortOrder"`
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
		Links:              toProfileLinkPayloads(profile.Links),
	}
}

func toProfileLinkMasterPayloads(values []domain.ProfileLinkMaster) []profileLinkMasterPayload {
	result := make([]profileLinkMasterPayload, 0, len(values))
	for _, value := range values {
		result = append(result, toProfileLinkMasterPayload(value))
	}

	return result
}

func toProfileLinkMasterPayload(value domain.ProfileLinkMaster) profileLinkMasterPayload {
	return profileLinkMasterPayload{
		ID:          value.ID,
		Key:         value.Key,
		Name:        value.Name,
		Icon:        value.Icon,
		Placeholder: value.Placeholder,
		SortOrder:   value.SortOrder,
	}
}

func toProfileLinkPayloads(values []domain.ProfileLink) []profileLinkPayload {
	result := make([]profileLinkPayload, 0, len(values))
	for _, value := range values {
		result = append(result, profileLinkPayload{
			ID:           value.ID,
			LinkMasterID: value.LinkMasterID,
			Key:          value.Key,
			Name:         value.Name,
			Icon:         value.Icon,
			Placeholder:  value.Placeholder,
			URL:          value.URL,
			SortOrder:    value.SortOrder,
		})
	}

	return result
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

func toUseCaseProfileLinks(values []profileLinkPayload) []domain.ProfileLink {
	result := make([]domain.ProfileLink, 0, len(values))
	for _, value := range values {
		result = append(result, domain.ProfileLink{
			ID:           value.ID,
			LinkMasterID: value.LinkMasterID,
			Key:          value.Key,
			URL:          value.URL,
			SortOrder:    value.SortOrder,
		})
	}

	return result
}

func toProfileLinkMasterInput(request profileLinkMasterRequest, includeID bool) domain.ProfileLinkMasterInput {
	input := domain.ProfileLinkMasterInput{
		Key:         strings.TrimSpace(request.Key),
		Name:        strings.TrimSpace(request.Name),
		Icon:        strings.TrimSpace(request.Icon),
		Placeholder: strings.TrimSpace(request.Placeholder),
		SortOrder:   request.SortOrder,
	}
	if includeID {
		input.ID = strings.TrimSpace(request.ID)
	}

	return input
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
