package users

import (
	"comu/internal/modules/users/application"
	"comu/internal/modules/users/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
}

type CreateUserResponse struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

type GetUserResponse struct {
	ID              uuid.UUID
	Name            string
	Email           string
	EmailVerifiedAt *time.Time
	Active          bool
	Avatar          string
	Password        string
	CreatedAt       time.Time
	DeletedAt       *time.Time
}

type UpdateUserPasswordRequest struct {
	ID          uuid.UUID
	NewPassword string
}

type publicApi struct {
	createUserUC              *application.CreateUserUC
	getUserByIdUC             *application.GetUserByIdUC
	getUserByEmailUC          *application.GetUserByEmailUC
	updateUserPasswordUC      *application.UpdateUserPasswordUC
	markUserEmailAsVerifiedUC *application.MarkUserEmailAsVerifiedUC
}

func newApi(
	createUserUC *application.CreateUserUC,
	getUserByIdUC *application.GetUserByIdUC,
	getUserByEmailUC *application.GetUserByEmailUC,
	updateUserPasswordUC *application.UpdateUserPasswordUC,
	markUserEmailAsVerifiedUC *application.MarkUserEmailAsVerifiedUC,
) *publicApi {
	return &publicApi{
		createUserUC:              createUserUC,
		getUserByIdUC:             getUserByIdUC,
		getUserByEmailUC:          getUserByEmailUC,
		updateUserPasswordUC:      updateUserPasswordUC,
		markUserEmailAsVerifiedUC: markUserEmailAsVerifiedUC,
	}
}

func (api *publicApi) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	user, err := api.createUserUC.Execute(
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
	user, err := api.getUserByEmailUC.Execute(ctx, email)

	if err != nil {
		return nil, err
	}

	return api.newGetUserResponse(user), nil
}

func (api *publicApi) GetUserByID(ctx context.Context, ID uuid.UUID) (*GetUserResponse, error) {
	user, err := api.getUserByIdUC.Execute(ctx, ID)

	if err != nil {
		return nil, err
	}

	return api.newGetUserResponse(user), nil
}

func (api *publicApi) MarkEmailAsVerified(ctx context.Context, email string) error {
	return api.markUserEmailAsVerifiedUC.Execute(ctx, email)
}

func (api *publicApi) UpdateUserPassword(ctx context.Context, req UpdateUserPasswordRequest) error {
	return api.updateUserPasswordUC.Execute(ctx, req.ID, req.NewPassword)
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
