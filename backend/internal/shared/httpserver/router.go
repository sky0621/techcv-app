package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sky0621/techcv-app/backend/internal/handler"
	"github.com/sky0621/techcv-app/backend/internal/repository"
	"github.com/sky0621/techcv-app/backend/internal/usecase"
)

func NewRouter(
	profileRepository repository.ProfileRepository,
	skillOptionsRepository repository.SkillOptionsRepository,
) http.Handler {
	r := chi.NewRouter()
	profileServer := handler.NewProfileHandler(usecase.NewProfileUseCase(profileRepository))
	skillOptionsServer := handler.NewSkillOptionsHandler(
		usecase.NewSkillOptionsUseCase(skillOptionsRepository),
	)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	api := chi.NewRouter()
	api.Get("/profile", profileServer.GetProfile)
	api.Put("/profile", profileServer.UpdateProfile)
	api.Get("/skills/options", skillOptionsServer.ListSkillOptions)
	api.Get("/skills", skillOptionsServer.ListSkills)
	api.Post("/skills", skillOptionsServer.CreateSkill)
	api.Put("/skills/{id}", skillOptionsServer.UpdateSkill)
	api.Delete("/skills/{id}", skillOptionsServer.DeleteSkill)
	api.Post("/skills/categories", skillOptionsServer.CreateSkillCategory)
	api.Put("/skills/categories/{id}", skillOptionsServer.UpdateSkillCategory)
	r.Mount("/api", api)

	return r
}
