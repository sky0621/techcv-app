package handler

import (
	"net/http"

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

type skillOptionsResponse struct {
	Categories        []skillOptionPayload `json:"categories"`
	ProficiencyLevels []skillOptionPayload `json:"proficiencyLevels"`
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
		result = append(result, skillOptionPayload{
			ID:        value.ID,
			Name:      value.Name,
			Icon:      value.Icon,
			SortOrder: value.SortOrder,
		})
	}

	return result
}
