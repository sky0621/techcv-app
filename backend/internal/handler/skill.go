package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sky0621/techcv-app/backend/internal/domain"
	"github.com/sky0621/techcv-app/backend/internal/usecase"
)

type SkillOptionsHandler struct {
	usecase *usecase.SkillOptionsUseCase
}

func NewSkillOptionsHandler(usecase *usecase.SkillOptionsUseCase) *SkillOptionsHandler {
	return &SkillOptionsHandler{usecase: usecase}
}

func (h *SkillOptionsHandler) ListSkillOptions(w http.ResponseWriter, r *http.Request) {
	options, err := h.usecase.List(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, skillOptionsResponse{
		Categories:        toSkillOptionPayloads(options.Categories),
		ProficiencyLevels: toSkillOptionPayloads(options.ProficiencyLevels),
	})
}

func (h *SkillOptionsHandler) CreateSkillCategory(w http.ResponseWriter, r *http.Request) {
	var request skillCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toSkillCategoryInput(request, true)
	if input.ID == "" || input.Name == "" || input.Icon == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id, name and icon are required")
		return
	}

	category, err := h.usecase.CreateCategory(r.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, skillCategoryResponse{
		Category: toSkillOptionPayload(*category),
	})
}

func (h *SkillOptionsHandler) UpdateSkillCategory(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	var request skillCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toSkillCategoryInput(request, false)
	if input.Name == "" || input.Icon == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "name and icon are required")
		return
	}

	category, err := h.usecase.UpdateCategory(r.Context(), id, input)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			writeJSONError(w, http.StatusNotFound, "not_found", "skill category not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, skillCategoryResponse{
		Category: toSkillOptionPayload(*category),
	})
}

func (h *SkillOptionsHandler) ListSkills(w http.ResponseWriter, r *http.Request) {
	skills, err := h.usecase.ListSkills(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, skillsResponse{
		Skills: toSkillPayloads(skills),
	})
}

func (h *SkillOptionsHandler) CreateSkill(w http.ResponseWriter, r *http.Request) {
	var request skillRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toSkillInput(request)
	if input.Name == "" || input.CategoryID == "" || input.ProficiencyLevelID == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "name, categoryId and proficiencyLevelId are required")
		return
	}

	skill, err := h.usecase.CreateSkill(r.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, skillResponse{
		Skill: toSkillPayload(*skill),
	})
}

func (h *SkillOptionsHandler) UpdateSkill(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	var request skillRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toSkillInput(request)
	if input.Name == "" || input.CategoryID == "" || input.ProficiencyLevelID == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "name, categoryId and proficiencyLevelId are required")
		return
	}

	skill, err := h.usecase.UpdateSkill(r.Context(), id, input)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			writeJSONError(w, http.StatusNotFound, "not_found", "skill not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, skillResponse{
		Skill: toSkillPayload(*skill),
	})
}

func (h *SkillOptionsHandler) DeleteSkill(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	if err := h.usecase.DeleteSkill(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			writeJSONError(w, http.StatusNotFound, "not_found", "skill not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type skillOptionsResponse struct {
	Categories        []skillOptionPayload `json:"categories"`
	ProficiencyLevels []skillOptionPayload `json:"proficiencyLevels"`
}

type skillCategoryResponse struct {
	Category skillOptionPayload `json:"category"`
}

type skillsResponse struct {
	Skills []skillPayload `json:"skills"`
}

type skillResponse struct {
	Skill skillPayload `json:"skill"`
}

type skillCategoryRequest struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type skillRequest struct {
	Name               string `json:"name"`
	CategoryID         string `json:"categoryId"`
	Experience         string `json:"experience"`
	ProficiencyLevelID string `json:"proficiencyLevelId"`
}

type skillOptionPayload struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon,omitempty"`
	SortOrder int64  `json:"sortOrder"`
}

type skillPayload struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	CategoryID         string `json:"categoryId"`
	Category           string `json:"category"`
	Experience         string `json:"experience"`
	ProficiencyLevelID string `json:"proficiencyLevelId"`
	Proficiency        string `json:"proficiency"`
	SortOrder          int64  `json:"sortOrder"`
}

func toSkillOptionPayloads(values []domain.SkillOption) []skillOptionPayload {
	result := make([]skillOptionPayload, 0, len(values))
	for _, value := range values {
		result = append(result, toSkillOptionPayload(value))
	}

	return result
}

func toSkillOptionPayload(value domain.SkillOption) skillOptionPayload {
	return skillOptionPayload{
		ID:        value.ID,
		Name:      value.Name,
		Icon:      value.Icon,
		SortOrder: value.SortOrder,
	}
}

func toSkillPayloads(values []domain.Skill) []skillPayload {
	result := make([]skillPayload, 0, len(values))
	for _, value := range values {
		result = append(result, toSkillPayload(value))
	}

	return result
}

func toSkillPayload(value domain.Skill) skillPayload {
	return skillPayload{
		ID:                 value.ID,
		Name:               value.Name,
		CategoryID:         value.CategoryID,
		Category:           value.CategoryName,
		Experience:         value.Experience,
		ProficiencyLevelID: value.ProficiencyLevelID,
		Proficiency:        value.ProficiencyName,
		SortOrder:          value.SortOrder,
	}
}

func toSkillCategoryInput(request skillCategoryRequest, includeID bool) domain.SkillCategoryInput {
	input := domain.SkillCategoryInput{
		Name: strings.TrimSpace(request.Name),
		Icon: strings.TrimSpace(request.Icon),
	}
	if includeID {
		input.ID = strings.TrimSpace(request.ID)
	}

	return input
}

func toSkillInput(request skillRequest) domain.SkillInput {
	return domain.SkillInput{
		Name:               strings.TrimSpace(request.Name),
		CategoryID:         strings.TrimSpace(request.CategoryID),
		Experience:         strings.TrimSpace(request.Experience),
		ProficiencyLevelID: strings.TrimSpace(request.ProficiencyLevelID),
	}
}
