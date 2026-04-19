package httpserver

import (
	"encoding/json"
	"net/http"

	sharedopenapi "github.com/sky0621/techcv-app/backend/internal/shared/openapi"
)

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(sharedopenapi.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
