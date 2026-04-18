package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sky0621/techcv-app/backend/internal/profile/domain"
	profileusecase "github.com/sky0621/techcv-app/backend/internal/profile/usecase"
)

type ProfileServer struct {
	usecase *profileusecase.UseCase
}

func NewProfileServer(usecase *profileusecase.UseCase) *ProfileServer {
	return &ProfileServer{usecase: usecase}
}

func (s *ProfileServer) GetProfile(ctx context.Context, _ GetProfileRequestObject) (GetProfileResponseObject, error) {
	profile, err := s.usecase.Get(ctx)
	if err != nil {
		return nil, err
	}

	return GetProfile200JSONResponse{
		Profile: toOpenAPIProfile(profile),
	}, nil
}

func (s *ProfileServer) UpdateProfile(ctx context.Context, request UpdateProfileRequestObject) (UpdateProfileResponseObject, error) {
	if request.Body == nil {
		return UpdateProfile400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{
				Code:    "bad_request",
				Message: "invalid request body",
			},
		}, nil
	}

	profile, err := s.usecase.Update(ctx, profileusecase.ProfileInput{
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

	return UpdateProfile200JSONResponse{
		Profile: toOpenAPIProfile(profile),
	}, nil
}

func toOpenAPIProfile(profile *domain.Profile) Profile {
	return Profile{
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

func toOpenAPIVisibilitySettings(values map[string]any) VisibilitySettings {
	result := VisibilitySettings{}
	for key, value := range values {
		booleanValue, ok := value.(bool)
		if ok {
			result[key] = booleanValue
		}
	}

	return result
}

func toUseCaseVisibilitySettings(values *VisibilitySettings) map[string]any {
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

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	})
}
