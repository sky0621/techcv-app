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

type skillOptionsResponse struct {
	Categories        []skillOptionPayload `json:"categories"`
	ProficiencyLevels []skillOptionPayload `json:"proficiencyLevels"`
}

type skillCategoryResponse struct {
	Category skillOptionPayload `json:"category"`
}

type skillCategoryRequest struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type skillOptionPayload struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon,omitempty"`
	SortOrder int64  `json:"sortOrder"`
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
