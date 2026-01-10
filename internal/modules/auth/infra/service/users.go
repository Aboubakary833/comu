package service

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/users"
	"comu/internal/shared"
	"context"
	"errors"

	"github.com/google/uuid"
)

type userService struct {
	api    users.PublicApi
	logger *shared.Log
}

func NewUserService(api users.PublicApi, logger *shared.Log) *userService {
	return &userService{
		api:    api,
		logger: logger,
	}
}

func (service *userService) GetUser(fn func() (*users.GetUserResponse, error)) (*domain.AuthUser, error) {
	response, err := fn()

	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		service.logger.Error.Println(err)
		return nil, domain.ErrInternal
	}

	return service.newAuthUserFromGetUserResponse(response), nil
}

func (service *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.AuthUser, error) {
	return service.GetUser(func() (*users.GetUserResponse, error) {
		return service.api.GetUserByID(ctx, id)
	})
}

func (service *userService) GetUserByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	return service.GetUser(func() (*users.GetUserResponse, error) {
		return service.api.GetUserByEmail(ctx, email)
	})
}

func (service *userService) CreateUser(ctx context.Context, name, email, password string) (uuid.UUID, error) {

	response, err := service.api.CreateUser(ctx, users.CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})

	if err != nil {
		if !errors.Is(err, users.ErrUserEmailTaken) {
			service.logger.Error.Println(err)
			return uuid.UUID{}, domain.ErrInternal
		}

		return uuid.UUID{}, err
	}

	return response.ID, nil
}

func (service *userService) MarkUserAsVerified(ctx context.Context, userEmail string) error {
	err := service.api.MarkEmailAsVerified(ctx, userEmail)

	if err != nil {
		if !errors.Is(err, users.ErrUserNotFound) {
			service.logger.Error.Println(err)
			return domain.ErrInternal
		}

		return domain.ErrUserNotFound
	}

	return nil
}

func (service *userService) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	err := service.api.UpdateUserPassword(
		ctx, users.UpdateUserPasswordRequest{
			ID:          userID,
			NewPassword: newPassword,
		},
	)

	if err != nil {
		if !errors.Is(err, users.ErrUserNotFound) {
			service.logger.Error.Println(err)
			return domain.ErrInternal
		}

		return domain.ErrUserNotFound
	}

	return nil
}

func (service *userService) newAuthUserFromGetUserResponse(response *users.GetUserResponse) *domain.AuthUser {
	return &domain.AuthUser{
		ID:              response.ID,
		Name:            response.Name,
		Email:           response.Email,
		EmailVerifiedAt: response.EmailVerifiedAt,
		Avatar:          response.Avatar,
		Active:          response.Active,
		Password:        response.Password,
		CreatedAt:       response.CreatedAt,
		DeletedAt:       response.DeletedAt,
	}
}
