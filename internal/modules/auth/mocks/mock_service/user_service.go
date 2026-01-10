package mockService

import (
	"comu/internal/modules/auth/domain"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type userServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *userServiceMock {
	return new(userServiceMock)
}

func (serviceMock *userServiceMock) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.AuthUser, error) {
	args := serviceMock.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.AuthUser), args.Error(1)
}

func (serviceMock *userServiceMock) GetUserByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	args := serviceMock.Called(ctx, email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.AuthUser), args.Error(1)
}

func (serviceMock *userServiceMock) CreateNewUser(ctx context.Context, name, email, password string) (uuid.UUID, error) {
	args := serviceMock.Called(ctx, name, email, password)

	if args.Get(0) == nil {
		return uuid.UUID{}, args.Error(1)
	}

	return args.Get(0).(uuid.UUID), nil
}

func (serviceMock *userServiceMock) MarkUserEmailAsVerified(ctx context.Context, userEmail string) error {
	args := serviceMock.Called(ctx, userEmail)
	return args.Error(0)
}

func (serviceMock *userServiceMock) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	args := serviceMock.Called(ctx, userID, newPassword)
	return args.Error(0)
}
