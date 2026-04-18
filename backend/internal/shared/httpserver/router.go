package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sky0621/techcv-app/backend/internal/profile/handler"
	profilerepository "github.com/sky0621/techcv-app/backend/internal/profile/repository"
	profileusecase "github.com/sky0621/techcv-app/backend/internal/profile/usecase"
)

func NewRouter(profileRepository profilerepository.Repository) http.Handler {
	r := chi.NewRouter()
	profileHandler := handler.New(
		profileusecase.New(
			profileRepository,
		),
	)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/api", func(api chi.Router) {
		api.Get("/profile", profileHandler.GetProfile)
		api.Put("/profile", profileHandler.UpdateProfile)
	})

	return r
}
