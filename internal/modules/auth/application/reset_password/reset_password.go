package resetPassword

import (
	"comu/internal/modules/auth/domain"
	"context"
)

type resetPasswordUC struct {
	userService         domain.UserService
	otpCodesRepository  domain.OtpCodesRepository
	notificationService domain.NotificationService
	resendOtpRequestsRepository domain.ResendOtpRequestsRepository
}

func NewUseCase(
	userService domain.UserService,
	otpCodesRepository domain.OtpCodesRepository,
	notificationService domain.NotificationService,
	resendOtpRequestsRepository domain.ResendOtpRequestsRepository,
) *resetPasswordUC {
	return &resetPasswordUC{
		userService:         userService,
		otpCodesRepository:  otpCodesRepository,
		notificationService: notificationService,
		resendOtpRequestsRepository: resendOtpRequestsRepository,
	}
}

func (useCase *resetPasswordUC) Execute(ctx context.Context, userEmail string) error {

	_, err := useCase.userService.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	otpCode, err := useCase.otpCodesRepository.CreateWithUserEmail(ctx, domain.ResetPasswordOTP, userEmail)

	if err != nil {
		return err
	}
	err = useCase.resendOtpRequestsRepository.CreateNew(ctx, userEmail)
	
	if err != nil {
		return err
	}

	return useCase.notificationService.SendOtpCodeMessage(otpCode)
}
