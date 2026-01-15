package otp

import (
	"comu/internal/modules/auth/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

type ResendOtpInput struct {
	ID          uuid.UUID
	UserEmail   string
	OtpCodeType domain.OtpType
}

type ResendOtpUC struct {
	otpCodesRepository          domain.OtpCodesRepository
	notificationService         domain.NotificationService
	resendOtpRequestsRepository domain.ResendOtpRequestsRepository
}

func NewResendOtpUseCase(
	otpCodesRepository domain.OtpCodesRepository,
	notificationService domain.NotificationService,
	resendOtpRequestsRepository domain.ResendOtpRequestsRepository,
) *ResendOtpUC {
	return &ResendOtpUC{
		otpCodesRepository:          otpCodesRepository,
		notificationService:         notificationService,
		resendOtpRequestsRepository: resendOtpRequestsRepository,
	}
}

func (useCase *ResendOtpUC) Execute(ctx context.Context, input ResendOtpInput) error {
	req, err := useCase.resendOtpRequestsRepository.FindByID(ctx, input.ID)

	if err != nil {
		return err
	} else if req.UserEmail != input.UserEmail {
		return domain.ErrInvalidResendRequest
	}

	otpCode, err := useCase.otpCodesRepository.FindByUserEmail(ctx, input.UserEmail)

	if err != nil {
		if errors.Is(err, domain.ErrOtpNotFound) {
			return domain.ErrInvalidResendRequest
		}
		return err
	}

	if otpCode.Type != input.OtpCodeType {
		return domain.ErrInvalidResendRequest
	}

	if !req.CanOtpBeSent() {
		return domain.ErrResendRequestCantBeProcessed
	}

	if req.IsCountExceeded() {
		return domain.ErrResendRequestCountExceeded
	}
	useCase.otpCodesRepository.Delete(ctx, otpCode)
	otpCode, err = useCase.otpCodesRepository.CreateWithUserEmail(ctx, otpCode.Type, otpCode.UserEmail)

	if err != nil {
		return err
	}

	useCase.resendOtpRequestsRepository.IncrementCount(ctx, req)
	return useCase.notificationService.SendOtpCodeMessage(otpCode)
}
