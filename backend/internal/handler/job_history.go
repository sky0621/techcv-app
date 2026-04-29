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

type JobHistoryHandler struct {
	usecase *usecase.JobHistoryUseCase
}

func NewJobHistoryHandler(usecase *usecase.JobHistoryUseCase) *JobHistoryHandler {
	return &JobHistoryHandler{usecase: usecase}
}

func (h *JobHistoryHandler) ListJobHistories(w http.ResponseWriter, r *http.Request) {
	jobHistories, err := h.usecase.List(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, jobHistoriesResponse{
		JobHistories: toJobHistoryPayloads(jobHistories),
	})
}

func (h *JobHistoryHandler) ListJobHistoryOptions(w http.ResponseWriter, r *http.Request) {
	options, err := h.usecase.ListOptions(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, jobHistoryOptionsResponse{
		EmploymentTypes: toJobEmploymentTypePayloads(options.EmploymentTypes),
	})
}

func (h *JobHistoryHandler) CreateJobHistory(w http.ResponseWriter, r *http.Request) {
	var request jobHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toJobHistoryInput(request)
	if !isValidJobHistoryInput(input) {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "company, displayName, startYear, startMonth and employmentTypeId are required")
		return
	}

	jobHistory, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, jobHistoryResponse{
		JobHistory: toJobHistoryPayload(*jobHistory),
	})
}

func (h *JobHistoryHandler) UpdateJobHistory(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	var request jobHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toJobHistoryInput(request)
	if !isValidJobHistoryInput(input) {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "company, displayName, startYear, startMonth and employmentTypeId are required")
		return
	}

	jobHistory, err := h.usecase.Update(r.Context(), id, input)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			writeJSONError(w, http.StatusNotFound, "not_found", "job history not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, jobHistoryResponse{
		JobHistory: toJobHistoryPayload(*jobHistory),
	})
}

func (h *JobHistoryHandler) DeleteJobHistory(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			writeJSONError(w, http.StatusNotFound, "not_found", "job history not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *JobHistoryHandler) CreateJobEmploymentType(w http.ResponseWriter, r *http.Request) {
	var request jobEmploymentTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toJobEmploymentTypeInput(request, true)
	if input.ID == "" || input.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id and name are required")
		return
	}

	employmentType, err := h.usecase.CreateEmploymentType(r.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, jobEmploymentTypeResponse{
		EmploymentType: toJobEmploymentTypePayload(*employmentType),
	})
}

func (h *JobHistoryHandler) UpdateJobEmploymentType(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "id is required")
		return
	}

	var request jobEmploymentTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	input := toJobEmploymentTypeInput(request, false)
	if input.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request", "name is required")
		return
	}

	employmentType, err := h.usecase.UpdateEmploymentType(r.Context(), id, input)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			writeJSONError(w, http.StatusNotFound, "not_found", "job employment type not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, jobEmploymentTypeResponse{
		EmploymentType: toJobEmploymentTypePayload(*employmentType),
	})
}

type jobHistoriesResponse struct {
	JobHistories []jobHistoryPayload `json:"jobHistories"`
}

type jobHistoryOptionsResponse struct {
	EmploymentTypes []jobEmploymentTypePayload `json:"employmentTypes"`
}

type jobHistoryResponse struct {
	JobHistory jobHistoryPayload `json:"jobHistory"`
}

type jobEmploymentTypeResponse struct {
	EmploymentType jobEmploymentTypePayload `json:"employmentType"`
}

type jobHistoryRequest struct {
	Company          string `json:"company"`
	DisplayName      string `json:"displayName"`
	StartYear        int64  `json:"startYear"`
	StartMonth       int64  `json:"startMonth"`
	EndYear          *int64 `json:"endYear"`
	EndMonth         *int64 `json:"endMonth"`
	EmploymentTypeID string `json:"employmentTypeId"`
}

type jobHistoryPayload struct {
	ID               string `json:"id"`
	Company          string `json:"company"`
	DisplayName      string `json:"displayName"`
	StartYear        int64  `json:"startYear"`
	StartMonth       int64  `json:"startMonth"`
	EndYear          *int64 `json:"endYear"`
	EndMonth         *int64 `json:"endMonth"`
	EmploymentTypeID string `json:"employmentTypeId"`
	EmploymentType   string `json:"employmentType"`
	ProjectCount     int64  `json:"projectCount"`
	SortOrder        int64  `json:"sortOrder"`
}

type jobEmploymentTypeRequest struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type jobEmploymentTypePayload struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SortOrder int64  `json:"sortOrder"`
}

func toJobHistoryPayloads(values []domain.JobHistory) []jobHistoryPayload {
	result := make([]jobHistoryPayload, 0, len(values))
	for _, value := range values {
		result = append(result, toJobHistoryPayload(value))
	}

	return result
}

func toJobHistoryPayload(value domain.JobHistory) jobHistoryPayload {
	return jobHistoryPayload{
		ID:               value.ID,
		Company:          value.Company,
		DisplayName:      value.DisplayName,
		StartYear:        value.StartYear,
		StartMonth:       value.StartMonth,
		EndYear:          value.EndYear,
		EndMonth:         value.EndMonth,
		EmploymentTypeID: value.EmploymentTypeID,
		EmploymentType:   value.EmploymentType,
		ProjectCount:     value.ProjectCount,
		SortOrder:        value.SortOrder,
	}
}

func toJobHistoryInput(request jobHistoryRequest) domain.JobHistoryInput {
	return domain.JobHistoryInput{
		Company:          strings.TrimSpace(request.Company),
		DisplayName:      strings.TrimSpace(request.DisplayName),
		StartYear:        request.StartYear,
		StartMonth:       request.StartMonth,
		EndYear:          request.EndYear,
		EndMonth:         request.EndMonth,
		EmploymentTypeID: strings.TrimSpace(request.EmploymentTypeID),
	}
}

func isValidJobHistoryInput(input domain.JobHistoryInput) bool {
	if input.Company == "" || input.DisplayName == "" || input.StartYear == 0 || !isValidMonth(input.StartMonth) || input.EmploymentTypeID == "" {
		return false
	}
	if input.EndYear == nil || input.EndMonth == nil {
		return input.EndYear == nil && input.EndMonth == nil
	}

	return *input.EndYear != 0 && isValidMonth(*input.EndMonth)
}

func isValidMonth(value int64) bool {
	return value >= 1 && value <= 12
}

func toJobEmploymentTypePayloads(values []domain.JobEmploymentType) []jobEmploymentTypePayload {
	result := make([]jobEmploymentTypePayload, 0, len(values))
	for _, value := range values {
		result = append(result, toJobEmploymentTypePayload(value))
	}

	return result
}

func toJobEmploymentTypePayload(value domain.JobEmploymentType) jobEmploymentTypePayload {
	return jobEmploymentTypePayload{
		ID:        value.ID,
		Name:      value.Name,
		SortOrder: value.SortOrder,
	}
}

func toJobEmploymentTypeInput(request jobEmploymentTypeRequest, includeID bool) domain.JobEmploymentTypeInput {
	input := domain.JobEmploymentTypeInput{
		Name: strings.TrimSpace(request.Name),
	}
	if includeID {
		input.ID = strings.TrimSpace(request.ID)
	}

	return input
}
