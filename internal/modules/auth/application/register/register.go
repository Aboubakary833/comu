package register

import (
	"comu/internal/modules/auth/domain"
	"context"
)

type RegisterUC struct {
	userService         domain.UserService
	passwordService     domain.PasswordService
	otpCodeRepository   domain.OtpCodesRepository
	notificationService domain.NotificationService
}

func NewRegisterUseCase(
	userService domain.UserService,
	passwordService domain.PasswordService,
	otpCodeRepository domain.OtpCodesRepository,
	notificationService domain.NotificationService,

) *RegisterUC {
	return &RegisterUC{
		userService:         userService,
		passwordService:     passwordService,
		otpCodeRepository:   otpCodeRepository,
		notificationService: notificationService,
	}
}

func (useCase *RegisterUC) Execute(ctx context.Context, name, email, password string) error {
	hashedPassword, err := useCase.passwordService.Hash(password)

	if err != nil {
		return err
	}

	_, err = useCase.userService.CreateNewUser(ctx, name, email, hashedPassword)

	if err != nil {
		return err
	}

	otpCode, err := useCase.otpCodeRepository.CreateWithUserEmail(ctx, domain.RegisterOTP, email)

	if err != nil {
		return err
	}

	useCase.notificationService.SendOtpCodeMessage(otpCode)

	return nil
}
