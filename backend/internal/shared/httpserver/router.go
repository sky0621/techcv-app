package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sky0621/techcv-app/backend/internal/handler"
	"github.com/sky0621/techcv-app/backend/internal/repository"
	"github.com/sky0621/techcv-app/backend/internal/usecase"
)

func NewRouter(profileRepository repository.ProfileRepository) http.Handler {
	r := chi.NewRouter()
	profileServer := handler.NewProfileHandler(usecase.NewProfileUseCase(profileRepository))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	api := chi.NewRouter()
	api.Get("/profile", profileServer.GetProfile)
	api.Put("/profile", profileServer.UpdateProfile)
	r.Mount("/api", api)

	return r
}
