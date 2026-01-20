package otp

import (
	"comu/internal/modules/auth/domain"
	"context"
)

type GenResendOtpRequestUC struct {
	repository domain.ResendOtpRequestsRepository
}

func NewGenResendRequestUseCase(resendRequestsRepository domain.ResendOtpRequestsRepository) *GenResendOtpRequestUC {
	return &GenResendOtpRequestUC{repository: resendRequestsRepository}
}

func (useCase *GenResendOtpRequestUC) Execute(ctx context.Context, userEmail string) (*domain.ResendOtpRequest, error) {
	resendRequest := domain.NewResendOtpRequest(userEmail)
	err := useCase.repository.Store(ctx, resendRequest)

	if err != nil {
		return nil, err
	}

	return resendRequest, nil
}
