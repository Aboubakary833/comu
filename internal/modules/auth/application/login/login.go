package login

import (
	"comu/internal/modules/auth/domain"
	"context"
	"errors"
)

type LoginUC struct {
	userService         domain.UserService
	otpCodeRepository   domain.OtpCodesRepository
	notificationService domain.NotificationService
	passwordService    domain.PasswordService
	resendOtpRequestsRepository domain.ResendOtpRequestsRepository
}

func NewUseCase(
	userService domain.UserService,
	passwordService domain.PasswordService,
	otpCodeRepository domain.OtpCodesRepository,
	notificationService domain.NotificationService,
	resendOtpRequestsRepository domain.ResendOtpRequestsRepository,
) *LoginUC {
	return &LoginUC{
		userService:         userService,
		passwordService:    passwordService,
		otpCodeRepository:   otpCodeRepository,
		notificationService: notificationService,
		resendOtpRequestsRepository: resendOtpRequestsRepository,
	}
}

func (useCase *LoginUC) Execute(ctx context.Context, email, password string) error {
	user, err := useCase.userService.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(domain.ErrUserNotFound, err) {
			return domain.ErrInvalidCredentials
		}

		return err
	}
	err = useCase.passwordService.Compare(user.Password, password)

	if err != nil {
		return domain.ErrInvalidCredentials
	}
	otpCode, err := useCase.otpCodeRepository.CreateWithUserEmail(ctx, domain.LoginOTP, user.Email)

	if err != nil {
		return err
	}
	err = useCase.resendOtpRequestsRepository.CreateNew(ctx, user.Email)
	
	if err != nil {
		return err
	}
	useCase.notificationService.SendOtpCodeMessage(otpCode)

	return nil
}
