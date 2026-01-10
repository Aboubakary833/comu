package mockRepository

import (
	"comu/internal/modules/auth/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

type otpCodesRepositoryMock struct {
	mock.Mock
}

func NewOtpCodesRepositoryMock() *otpCodesRepositoryMock {
	return new(otpCodesRepositoryMock)
}

func (repoMock *otpCodesRepositoryMock) Exists(ctx context.Context, value string) bool {
	args := repoMock.Called(ctx, value)
	return args.Bool(0)
}

func (repoMock *otpCodesRepositoryMock) Find(ctx context.Context, value string) (*domain.OtpCode, error) {
	args := repoMock.Called(ctx, value)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.OtpCode), nil
}

func (repoMock *otpCodesRepositoryMock) FindByUserEmail(ctx context.Context, userEmail string) (*domain.OtpCode, error) {
	args := repoMock.Called(ctx, userEmail)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.OtpCode), nil
}

func (repoMock *otpCodesRepositoryMock) Store(ctx context.Context, otpCode *domain.OtpCode) error {
	args := repoMock.Called(ctx, otpCode)
	return args.Error(0)
}

func (repoMock *otpCodesRepositoryMock) CreateWithUserEmail(ctx context.Context, otpType domain.OtpType, userEmail string) (*domain.OtpCode, error) {
	args := repoMock.Called(ctx, otpType, userEmail)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.OtpCode), nil
}

func (repoMock *otpCodesRepositoryMock) Delete(ctx context.Context, otpCode *domain.OtpCode) error {
	args := repoMock.Called(ctx, otpCode)
	return args.Error(0)
}
