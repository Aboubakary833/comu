package newPassword

import (
	"comu/internal/modules/auth/domain"
	"context"
	"errors"
)

type setNewPasswordUC struct {
	userService domain.UserService
	passwordService domain.PasswordService
	notificationService domain.NotificationService
	resetTokensRepository domain.ResetTokensRepository
}

func NewUseCase(
	userService domain.UserService,
	passwordService domain.PasswordService,
	notificationService domain.NotificationService,
	resetTokensRepository domain.ResetTokensRepository,
) *setNewPasswordUC {
	return &setNewPasswordUC{
		userService: userService,
		passwordService: passwordService,
		notificationService: notificationService,
		resetTokensRepository: resetTokensRepository,
	}
}


func (useCase *setNewPasswordUC) Execute(ctx context.Context, tokenString string, newPassword string) error {
	hashedNewPassword, err := useCase.passwordService.Hash(newPassword)
	
	if err != nil {
		return err
	}
	token, err := useCase.resetTokensRepository.Find(ctx, tokenString)

	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return domain.ErrInvalidToken
		}
		return err
	}
	err = useCase.userService.UpdateUserPassword(ctx, token.UserID, hashedNewPassword)

	if err != nil {
		return err
	}
	useCase.resetTokensRepository.Delete(ctx, tokenString)
	useCase.notificationService.SendPasswordChangedMessage(token.UserEmail)

	return nil
}
