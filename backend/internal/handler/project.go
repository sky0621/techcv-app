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

type ProjectHandler struct {
	usecase *usecase.ProjectUseCase
}

func NewProjectHandler(usecase *usecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{usecase: usecase}
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.usecase.List(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, projectsResponse{
		Projects: toProjectPayloads(projects),
	})
}

func (h *ProjectHandler) ListProjectOptions(w http.ResponseWriter, r *http.Request) {
	options, err := h.usecase.ListOptions(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, projectOptionsResponse{
		Phases: toProjectPhasePayloads(options.Phases),
	})
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var request projectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toProjectInput(request)
	if !isValidProjectInput(input) {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "name, company, startYear and startMonth are required")
		return
	}

	project, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, projectResponse{
		Project: toProjectPayload(*project),
	})
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	var request projectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toProjectInput(request)
	if !isValidProjectInput(input) {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "name, company, startYear and startMonth are required")
		return
	}

	project, err := h.usecase.Update(r.Context(), id, input)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			writeJSONError(w, http.StatusNotFound, "not_found", "project not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, projectResponse{
		Project: toProjectPayload(*project),
	})
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			writeJSONError(w, http.StatusNotFound, "not_found", "project not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) CreateProjectPhase(w http.ResponseWriter, r *http.Request) {
	var request projectPhaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toProjectPhaseInput(request, true)
	if input.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "name is required")
		return
	}

	phase, err := h.usecase.CreatePhase(r.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, projectPhaseResponse{
		Phase: toProjectPhasePayload(*phase),
	})
}

func (h *ProjectHandler) UpdateProjectPhase(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	var request projectPhaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toProjectPhaseInput(request, false)
	if input.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "name is required")
		return
	}

	phase, err := h.usecase.UpdatePhase(r.Context(), id, input)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			writeJSONError(w, http.StatusNotFound, "not_found", "project phase not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, projectPhaseResponse{
		Phase: toProjectPhasePayload(*phase),
	})
}

type projectsResponse struct {
	Projects []projectPayload `json:"projects"`
}

type projectOptionsResponse struct {
	Phases []projectPhasePayload `json:"phases"`
}

type projectResponse struct {
	Project projectPayload `json:"project"`
}

type projectPhaseResponse struct {
	Phase projectPhasePayload `json:"phase"`
}

type projectRequest struct {
	Name         string   `json:"name"`
	Company      string   `json:"company"`
	StartYear    int64    `json:"startYear"`
	StartMonth   int64    `json:"startMonth"`
	EndYear      *int64   `json:"endYear"`
	EndMonth     *int64   `json:"endMonth"`
	Description  string   `json:"description"`
	Role         string   `json:"role"`
	TeamSize     string   `json:"teamSize"`
	Technologies []string `json:"technologies"`
	Phases       []string `json:"phases"`
	Achievements string   `json:"achievements"`
	IsDraft      bool     `json:"isDraft"`
}

type projectPayload struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Company      string   `json:"company"`
	StartYear    int64    `json:"startYear"`
	StartMonth   int64    `json:"startMonth"`
	EndYear      *int64   `json:"endYear"`
	EndMonth     *int64   `json:"endMonth"`
	Description  string   `json:"description"`
	Role         string   `json:"role"`
	TeamSize     string   `json:"teamSize"`
	Technologies []string `json:"technologies"`
	Phases       []string `json:"phases"`
	Achievements string   `json:"achievements"`
	IsDraft      bool     `json:"isDraft"`
	SortOrder    int64    `json:"sortOrder"`
}

type projectPhaseRequest struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	SortOrder int64  `json:"sortOrder"`
}

type projectPhasePayload struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SortOrder int64  `json:"sortOrder"`
}

func toProjectPayloads(values []domain.Project) []projectPayload {
	result := make([]projectPayload, 0, len(values))
	for _, value := range values {
		result = append(result, toProjectPayload(value))
	}

	return result
}

func toProjectPayload(value domain.Project) projectPayload {
	return projectPayload{
		ID:           value.ID,
		Name:         value.Name,
		Company:      value.Company,
		StartYear:    value.StartYear,
		StartMonth:   value.StartMonth,
		EndYear:      value.EndYear,
		EndMonth:     value.EndMonth,
		Description:  value.Description,
		Role:         value.Role,
		TeamSize:     value.TeamSize,
		Technologies: value.Technologies,
		Phases:       value.Phases,
		Achievements: value.Achievements,
		IsDraft:      value.IsDraft,
		SortOrder:    value.SortOrder,
	}
}

func toProjectInput(request projectRequest) domain.ProjectInput {
	return domain.ProjectInput{
		Name:         strings.TrimSpace(request.Name),
		Company:      strings.TrimSpace(request.Company),
		StartYear:    request.StartYear,
		StartMonth:   request.StartMonth,
		EndYear:      request.EndYear,
		EndMonth:     request.EndMonth,
		Description:  strings.TrimSpace(request.Description),
		Role:         strings.TrimSpace(request.Role),
		TeamSize:     strings.TrimSpace(request.TeamSize),
		Technologies: trimStringSlice(request.Technologies),
		Phases:       trimStringSlice(request.Phases),
		Achievements: strings.TrimSpace(request.Achievements),
		IsDraft:      request.IsDraft,
	}
}

func toProjectPhasePayloads(values []domain.ProjectPhase) []projectPhasePayload {
	result := make([]projectPhasePayload, 0, len(values))
	for _, value := range values {
		result = append(result, toProjectPhasePayload(value))
	}

	return result
}

func toProjectPhasePayload(value domain.ProjectPhase) projectPhasePayload {
	return projectPhasePayload{
		ID:        value.ID,
		Name:      value.Name,
		SortOrder: value.SortOrder,
	}
}

func toProjectPhaseInput(request projectPhaseRequest, includeID bool) domain.ProjectPhaseInput {
	input := domain.ProjectPhaseInput{
		Name:      strings.TrimSpace(request.Name),
		SortOrder: request.SortOrder,
	}
	if includeID {
		input.ID = strings.TrimSpace(request.ID)
	}

	return input
}

func isValidProjectInput(input domain.ProjectInput) bool {
	if input.Name == "" || input.Company == "" || input.StartYear == 0 || !isValidMonth(input.StartMonth) {
		return false
	}
	if input.EndYear == nil || input.EndMonth == nil {
		return input.EndYear == nil && input.EndMonth == nil
	}

	return *input.EndYear != 0 && isValidMonth(*input.EndMonth)
}

func trimStringSlice(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		result = append(result, value)
	}

	return result
}
