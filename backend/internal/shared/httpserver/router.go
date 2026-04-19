package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	profilehandler "github.com/sky0621/techcv-app/backend/internal/profile/handler"
	profilerepository "github.com/sky0621/techcv-app/backend/internal/profile/repository"
	profileusecase "github.com/sky0621/techcv-app/backend/internal/profile/usecase"
	sharedopenapi "github.com/sky0621/techcv-app/backend/internal/shared/openapi"
)

func NewRouter(profileRepository profilerepository.Repository) http.Handler {
	r := chi.NewRouter()
	profileServer := profilehandler.New(profileusecase.New(profileRepository))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	api := chi.NewRouter()
	sharedopenapi.HandlerFromMux(
		sharedopenapi.NewStrictHandlerWithOptions(profileServer, nil, sharedopenapi.StrictHTTPServerOptions{
			RequestErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
				writeJSONError(w, http.StatusBadRequest, "bad_request", err.Error())
			},
			ResponseErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
				writeJSONError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			},
		}),
		api,
	)
	r.Mount("/api", api)

	return r
}
