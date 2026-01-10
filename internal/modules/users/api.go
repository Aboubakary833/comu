package users

import (
	"comu/internal/modules/users/application"
	"comu/internal/modules/users/domain"
	"context"

	"github.com/google/uuid"
)

type publicApi struct {
	repo domain.Repository
}

func newApi(repo domain.Repository) *publicApi {
	return &publicApi{
		repo: repo,
	}
}

func (api *publicApi) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	useCase := application.NewCreateUserUseCase(api.repo)
	user, err := useCase.Execute(
		ctx, application.CreateUserInput{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		},
	)

	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (api *publicApi) GetUserByEmail(ctx context.Context, email string) (*GetUserResponse, error) {
	useCase := application.NewGetUserByEmailUseCase(api.repo)
	user, err := useCase.Execute(ctx, email)

	if err != nil {
		return nil, err
	}

	return api.newGetUserResponse(user), nil
}

func (api *publicApi) GetUserByID(ctx context.Context, ID uuid.UUID) (*GetUserResponse, error) {
	useCase := application.NewGetUserByIdUseCase(api.repo)
	user, err := useCase.Execute(ctx, ID)

	if err != nil {
		return nil, err
	}

	return api.newGetUserResponse(user), nil
}

func (api *publicApi) MarkEmailAsVerified(ctx context.Context, email string) error {
	useCase := application.NewMarkUserEmailAsVerifiedUseCase(api.repo)
	return useCase.Execute(ctx, email)
}

func (api *publicApi) UpdateUserPassword(ctx context.Context, req UpdateUserPasswordRequest) error {
	useCase := application.NewUpdateUserPasswordUseCase(api.repo)
	return useCase.Execute(ctx, req.ID, req.NewPassword)
}

func (api *publicApi) newGetUserResponse(user *domain.User) *GetUserResponse {
	return &GetUserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Avatar:          user.Avatar,
		Active:          user.Active,
		Password:        user.Password,
		CreatedAt:       user.CreatedAt,
		DeletedAt:       user.DeletedAt,
	}
}
