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
	jobHistoryRepository repository.JobHistoryRepository,
) http.Handler {
	r := chi.NewRouter()
	profileServer := handler.NewProfileHandler(usecase.NewProfileUseCase(profileRepository))
	skillOptionsServer := handler.NewSkillOptionsHandler(
		usecase.NewSkillOptionsUseCase(skillOptionsRepository),
	)
	jobHistoryServer := handler.NewJobHistoryHandler(
		usecase.NewJobHistoryUseCase(jobHistoryRepository),
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
	api.Get("/job-histories/options", jobHistoryServer.ListJobHistoryOptions)
	api.Get("/job-histories", jobHistoryServer.ListJobHistories)
	api.Post("/job-histories", jobHistoryServer.CreateJobHistory)
	api.Put("/job-histories/{id}", jobHistoryServer.UpdateJobHistory)
	api.Delete("/job-histories/{id}", jobHistoryServer.DeleteJobHistory)
	api.Post("/job-employment-types", jobHistoryServer.CreateJobEmploymentType)
	api.Put("/job-employment-types/{id}", jobHistoryServer.UpdateJobEmploymentType)
	r.Mount("/api", api)

	return r
}
