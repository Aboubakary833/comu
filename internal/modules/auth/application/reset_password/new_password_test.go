package resetPassword

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/infra/memory"
	mockService "comu/internal/modules/auth/mocks/mock_service"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPasswordUseCase(t *testing.T) {

	t.Run("it should fail and return ErrInvalidToken", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		passwordService := mockService.NewPasswordServiceMock()
		notificationService := mockService.NewNotificationServiceMock()
		resetTokensRepository := memory.NewInMemoryResetTokensRepository(nil)

		tokenString := "zr7JAzt1zCEmQdNFxH6ukyYX+zk0aocNvLwNX69Qnhs="
		newPassword := "xdAPktpKLjcEy8ncy7Cqall95m4"
		hashedNewPassword := "cZDdc3CmKwYE8AoNMJ+kG4D52C6IYKzcuxWAmOSr1vs"

		passwordService.On("Hash", newPassword).Return(hashedNewPassword, nil)

		useCase := NewSetNewPasswordUseCase(userService, passwordService, notificationService, resetTokensRepository)

		err := useCase.Execute(context.Background(), tokenString, newPassword)

		assert.ErrorIs(t, err, domain.ErrInvalidToken)
		passwordService.AssertExpectations(t)
		userService.AssertNotCalled(t, "UpdateUserPassword")
		notificationService.AssertNotCalled(t, "SendPasswordChangedMessage")
	})

	t.Run("it should fail and return ErrExpiredToken", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		passwordService := mockService.NewPasswordServiceMock()
		notificationService := mockService.NewNotificationServiceMock()
		resetTokensRepository := memory.NewInMemoryResetTokensRepository(nil)
		ctx := context.Background()

		userID := uuid.New()
		userEmail := "johndoe@gmail.com"
		token := domain.NewResetToken(userID, userEmail, -20*time.Minute)
		newPassword := "xdAPktpKLjcEy8ncy7Cqall95m4"
		hashedNewPassword := "cZDdc3CmKwYE8AoNMJ+kG4D52C6IYKzcuxWAmOSr1vs"

		resetTokensRepository.Store(ctx, token)
		passwordService.On("Hash", newPassword).Return(hashedNewPassword, nil)

		useCase := NewSetNewPasswordUseCase(userService, passwordService, notificationService, resetTokensRepository)

		err := useCase.Execute(ctx, token.Token, newPassword)

		assert.ErrorIs(t, err, domain.ErrExpiredToken)
		passwordService.AssertExpectations(t)
		userService.AssertNotCalled(t, "UpdateUserPassword")
		notificationService.AssertNotCalled(t, "SendPasswordChangedMessage")
	})

	t.Run("it should fail and return ErrUserNotFound", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		passwordService := mockService.NewPasswordServiceMock()
		notificationService := mockService.NewNotificationServiceMock()
		resetTokensRepository := memory.NewInMemoryResetTokensRepository(nil)
		ctx := context.Background()

		userID := uuid.New()
		userEmail := "johndoe@gmail.com"
		token := domain.NewResetToken(userID, userEmail, domain.DefaultResetTokenTTL)
		newPassword := "xdAPktpKLjcEy8ncy7Cqall95m4"
		hashedNewPassword := "cZDdc3CmKwYE8AoNMJ+kG4D52C6IYKzcuxWAmOSr1vs"

		resetTokensRepository.Store(ctx, token)
		passwordService.On("Hash", newPassword).Return(hashedNewPassword, nil)
		userService.On("UpdateUserPassword", ctx, userID, hashedNewPassword).Return(domain.ErrUserNotFound)

		useCase := NewSetNewPasswordUseCase(userService, passwordService, notificationService, resetTokensRepository)

		err := useCase.Execute(ctx, token.Token, newPassword)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		passwordService.AssertExpectations(t)
		userService.AssertExpectations(t)
		notificationService.AssertNotCalled(t, "SendPasswordChangedMessage")
	})

	t.Run("it should succeed and update user password", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		passwordService := mockService.NewPasswordServiceMock()
		notificationService := mockService.NewNotificationServiceMock()
		resetTokensRepository := memory.NewInMemoryResetTokensRepository(nil)
		ctx := context.Background()

		userID := uuid.New()
		userEmail := "johndoe@gmail.com"
		token := domain.NewResetToken(userID, userEmail, domain.DefaultResetTokenTTL)
		newPassword := "xdAPktpKLjcEy8ncy7Cqall95m4"
		hashedNewPassword := "cZDdc3CmKwYE8AoNMJ+kG4D52C6IYKzcuxWAmOSr1vs"

		resetTokensRepository.Store(ctx, token)
		passwordService.On("Hash", newPassword).Return(hashedNewPassword, nil)
		userService.On("UpdateUserPassword", ctx, userID, hashedNewPassword).Return(nil)
		notificationService.On("SendPasswordChangedMessage", token.UserEmail).Return(nil)

		useCase := NewSetNewPasswordUseCase(userService, passwordService, notificationService, resetTokensRepository)

		err := useCase.Execute(ctx, token.Token, newPassword)

		assert.NoError(t, err)
		passwordService.AssertExpectations(t)
		userService.AssertExpectations(t)
		notificationService.AssertExpectations(t)
	})
}
