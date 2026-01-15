package mockRepository

import (
	"comu/internal/modules/auth/domain"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type resendOtpRequestsRepository struct {
	mock.Mock
}

func NewResendOtpRequestsRepositoryMock() *resendOtpRequestsRepository {
	return new(resendOtpRequestsRepository)
}

func (repoMock *resendOtpRequestsRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ResendOtpRequest, error) {
	args := repoMock.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.ResendOtpRequest), nil
}

func (repoMock *resendOtpRequestsRepository) FindByUserEmail(ctx context.Context, userEmail string) (*domain.ResendOtpRequest, error) {
	args := repoMock.Called(ctx, userEmail)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.ResendOtpRequest), nil
}

func (repoMock *resendOtpRequestsRepository) IncrementCount(ctx context.Context, req *domain.ResendOtpRequest) error {
	args := repoMock.Called(ctx, req)
	req.Count += 1

	return args.Error(0)
}

func (repoMock *resendOtpRequestsRepository) CreateNew(ctx context.Context, userEmail string) error {
	args := repoMock.Called(ctx, userEmail)
	return args.Error(0)
}

func (repoMock *resendOtpRequestsRepository) Store(ctx context.Context, req *domain.ResendOtpRequest) error {
	args := repoMock.Called(ctx, req)

	return args.Error(0)
}

func (repoMock *resendOtpRequestsRepository) Delete(ctx context.Context, req *domain.ResendOtpRequest) error {
	args := repoMock.Called(ctx, req)

	return args.Error(0)
}
