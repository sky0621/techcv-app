package handler

import (
	"context"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sky0621/techcv-app/backend/internal/profile/domain"
	"github.com/sky0621/techcv-app/backend/internal/profile/usecase"
	sharedopenapi "github.com/sky0621/techcv-app/backend/internal/shared/openapi"
)

type Handler struct {
	usecase *usecase.UseCase
}

func New(usecase *usecase.UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) GetProfile(ctx context.Context, _ sharedopenapi.GetProfileRequestObject) (sharedopenapi.GetProfileResponseObject, error) {
	profile, err := h.usecase.Get(ctx)
	if err != nil {
		return nil, err
	}

	return sharedopenapi.GetProfile200JSONResponse{
		Profile: toOpenAPIProfile(profile),
	}, nil
}

func (h *Handler) UpdateProfile(ctx context.Context, request sharedopenapi.UpdateProfileRequestObject) (sharedopenapi.UpdateProfileResponseObject, error) {
	if request.Body == nil {
		return sharedopenapi.UpdateProfile400JSONResponse{
			BadRequestJSONResponse: sharedopenapi.BadRequestJSONResponse{
				Code:    "bad_request",
				Message: "invalid request body",
			},
		}, nil
	}

	profile, err := h.usecase.Update(ctx, usecase.ProfileInput{
		FullName:           stringValue(request.Body.FullName),
		Nickname:           stringValue(request.Body.Nickname),
		Location:           stringValue(request.Body.Location),
		Email:              emailValue(request.Body.Email),
		Phone:              stringValue(request.Body.Phone),
		Summary:            stringValue(request.Body.Summary),
		GitHubURL:          stringValue(request.Body.GithubUrl),
		ZennURL:            stringValue(request.Body.ZennUrl),
		QiitaURL:           stringValue(request.Body.QiitaUrl),
		WebsiteURL:         stringValue(request.Body.WebsiteUrl),
		PreferredWorkStyle: stringValue(request.Body.PreferredWorkStyle),
		VisibilitySettings: toUseCaseVisibilitySettings(request.Body.VisibilitySettings),
	})
	if err != nil {
		return nil, err
	}

	return sharedopenapi.UpdateProfile200JSONResponse{
		Profile: toOpenAPIProfile(profile),
	}, nil
}

func toOpenAPIProfile(profile *domain.Profile) sharedopenapi.Profile {
	return sharedopenapi.Profile{
		Id:                 profile.ID,
		UserId:             profile.UserID,
		FullName:           stringPointer(profile.FullName),
		Nickname:           stringPointer(profile.Nickname),
		Location:           stringPointer(profile.Location),
		Email:              emailPointer(profile.Email),
		Phone:              stringPointer(profile.Phone),
		Summary:            stringPointer(profile.Summary),
		GithubUrl:          stringPointer(profile.GitHubURL),
		ZennUrl:            stringPointer(profile.ZennURL),
		QiitaUrl:           stringPointer(profile.QiitaURL),
		WebsiteUrl:         stringPointer(profile.WebsiteURL),
		PreferredWorkStyle: stringPointer(profile.PreferredWorkStyle),
		VisibilitySettings: toOpenAPIVisibilitySettings(profile.VisibilitySettings),
		CreatedAt:          profile.CreatedAt,
		UpdatedAt:          profile.UpdatedAt,
	}
}

func toOpenAPIVisibilitySettings(values map[string]any) sharedopenapi.VisibilitySettings {
	result := sharedopenapi.VisibilitySettings{}
	for key, value := range values {
		booleanValue, ok := value.(bool)
		if ok {
			result[key] = booleanValue
		}
	}

	return result
}

func toUseCaseVisibilitySettings(values *sharedopenapi.VisibilitySettings) map[string]any {
	if values == nil {
		return nil
	}

	result := make(map[string]any, len(*values))
	for key, value := range *values {
		result[key] = value
	}

	return result
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func stringPointer(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func emailValue(value *openapi_types.Email) string {
	if value == nil {
		return ""
	}

	return string(*value)
}

func emailPointer(value string) *openapi_types.Email {
	if value == "" {
		return nil
	}

	email := openapi_types.Email(value)
	return &email
}
