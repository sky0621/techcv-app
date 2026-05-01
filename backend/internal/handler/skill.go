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
		SkillMasters:      toSkillMasterPayloads(options.SkillMasters),
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

func (h *SkillOptionsHandler) CreateSkillMaster(w http.ResponseWriter, r *http.Request) {
	var request skillMasterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toSkillMasterInput(request, true)
	if input.ID == "" || input.Name == "" || input.CategoryID == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id, name and categoryId are required")
		return
	}

	skillMaster, err := h.usecase.CreateMaster(r.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, skillMasterResponse{
		SkillMaster: toSkillMasterPayload(*skillMaster),
	})
}

func (h *SkillOptionsHandler) UpdateSkillMaster(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	var request skillMasterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toSkillMasterInput(request, false)
	if input.Name == "" || input.CategoryID == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "name and categoryId are required")
		return
	}

	skillMaster, err := h.usecase.UpdateMaster(r.Context(), id, input)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			writeJSONError(w, http.StatusNotFound, "not_found", "skill master not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, skillMasterResponse{
		SkillMaster: toSkillMasterPayload(*skillMaster),
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
		if isUniqueConstraintError(err) {
			writeJSONError(w, http.StatusConflict, "conflict", "skill already exists")
			return
		}
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
		if isUniqueConstraintError(err) {
			writeJSONError(w, http.StatusConflict, "conflict", "skill already exists")
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

func isUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

type skillOptionsResponse struct {
	Categories        []skillOptionPayload `json:"categories"`
	ProficiencyLevels []skillOptionPayload `json:"proficiencyLevels"`
	SkillMasters      []skillMasterPayload `json:"skillMasters"`
}

type skillCategoryResponse struct {
	Category skillOptionPayload `json:"category"`
}

type skillMasterResponse struct {
	SkillMaster skillMasterPayload `json:"skillMaster"`
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

type skillMasterRequest struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	CategoryID string `json:"categoryId"`
}

type skillRequest struct {
	Name               string `json:"name"`
	CategoryID         string `json:"categoryId"`
	Experience         int64  `json:"experience"`
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
	Experience         int64  `json:"experience"`
	ProficiencyLevelID string `json:"proficiencyLevelId"`
	Proficiency        string `json:"proficiency"`
	SortOrder          int64  `json:"sortOrder"`
}

type skillMasterPayload struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CategoryID string `json:"categoryId"`
	Category   string `json:"category"`
	SortOrder  int64  `json:"sortOrder"`
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

func toSkillMasterPayloads(values []domain.SkillMaster) []skillMasterPayload {
	result := make([]skillMasterPayload, 0, len(values))
	for _, value := range values {
		result = append(result, toSkillMasterPayload(value))
	}

	return result
}

func toSkillMasterPayload(value domain.SkillMaster) skillMasterPayload {
	return skillMasterPayload{
		ID:         value.ID,
		Name:       value.Name,
		CategoryID: value.CategoryID,
		Category:   value.CategoryName,
		SortOrder:  value.SortOrder,
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

func toSkillMasterInput(request skillMasterRequest, includeID bool) domain.SkillMasterInput {
	input := domain.SkillMasterInput{
		Name:       strings.TrimSpace(request.Name),
		CategoryID: strings.TrimSpace(request.CategoryID),
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
		Experience:         request.Experience,
		ProficiencyLevelID: strings.TrimSpace(request.ProficiencyLevelID),
	}
}
