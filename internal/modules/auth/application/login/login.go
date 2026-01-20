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
	passwordService     domain.PasswordService
}

func NewUseCase(
	userService domain.UserService,
	passwordService domain.PasswordService,
	otpCodeRepository domain.OtpCodesRepository,
	notificationService domain.NotificationService,
) *LoginUC {
	return &LoginUC{
		userService:         userService,
		passwordService:     passwordService,
		otpCodeRepository:   otpCodeRepository,
		notificationService: notificationService,
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

	err = useCase.notificationService.SendOtpCodeMessage(otpCode)

	if err != nil {
		useCase.otpCodeRepository.Delete(ctx, otpCode)
		return err
	}

	return nil
}
