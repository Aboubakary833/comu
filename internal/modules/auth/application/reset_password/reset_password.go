package resetPassword

import (
	"comu/internal/modules/auth/domain"
	"context"
)

type ResetPasswordUC struct {
	userService         domain.UserService
	otpCodesRepository  domain.OtpCodesRepository
	notificationService domain.NotificationService
}

func NewResetPasswordUseCase(
	userService domain.UserService,
	otpCodesRepository domain.OtpCodesRepository,
	notificationService domain.NotificationService,
) *ResetPasswordUC {
	return &ResetPasswordUC{
		userService:         userService,
		otpCodesRepository:  otpCodesRepository,
		notificationService: notificationService,
	}
}

func (useCase *ResetPasswordUC) Execute(ctx context.Context, userEmail string) error {
	_, err := useCase.userService.GetUserByEmail(ctx, userEmail)

	if err != nil {
		return err
	}

	otpCode, err := useCase.otpCodesRepository.CreateWithUserEmail(ctx, domain.ResetPasswordOTP, userEmail)

	if err != nil {
		return err
	}

	err = useCase.notificationService.SendOtpCodeMessage(otpCode)

	if err != nil {
		useCase.otpCodesRepository.Delete(ctx, otpCode)
		return err
	}

	return nil
}
