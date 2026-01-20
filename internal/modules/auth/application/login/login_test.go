package login

import (
	"comu/internal/modules/auth/domain"
	mockRepository "comu/internal/modules/auth/mocks/mock_repository"
	mockService "comu/internal/modules/auth/mocks/mock_service"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUseCase(t *testing.T) {

	t.Run("it should result into success", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		passwordService := mockService.NewPasswordServiceMock()
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		userPassword := "BhVmqUnb6m1upSh"
		hashedPassword := "ixReNPXoBPxP9bIBQ6FziHj/9UG5wwzLbxP3vwpSZGo="

		user := domain.AuthUser{
			ID:       uuid.New(),
			Name:     "John Doe",
			Email:    userEmail,
			Password: hashedPassword,
		}

		otpCode := domain.NewOtpCode(domain.LoginOTP, userEmail, domain.DefaultOtpCodeTTL)

		userService.On("GetUserByEmail", ctx, userEmail).Return(&user, nil).Once()
		passwordService.On("Compare", hashedPassword, userPassword).Return(nil).Once()
		otpCodesRepository.On("CreateWithUserEmail", ctx, domain.LoginOTP, userEmail).Return(otpCode, nil).Once()
		notificationService.On("SendOtpCodeMessage", otpCode).Return(nil).Once()

		useCase := NewUseCase(
			userService,
			passwordService,
			otpCodesRepository,
			notificationService,
		)

		err := useCase.Execute(ctx, userEmail, userPassword)

		assert.NoError(t, err)
		userService.AssertExpectations(t)
		passwordService.AssertExpectations(t)
		otpCodesRepository.AssertExpectations(t)
		notificationService.AssertExpectations(t)
	})

	t.Run("it should fail and return ErrInvalidCredentials when user not found", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		passwordService := mockService.NewPasswordServiceMock()
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		userPassword := "BhVmqUnb6m1upSh"

		userService.On("GetUserByEmail", ctx, userEmail).Return(nil, domain.ErrUserNotFound).Once()

		useCase := NewUseCase(
			userService,
			passwordService,
			otpCodesRepository,
			notificationService,
		)

		err := useCase.Execute(ctx, userEmail, userPassword)

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		userService.AssertExpectations(t)
		passwordService.AssertNotCalled(t, "Compare")
		otpCodesRepository.AssertNotCalled(t, "CreateWithUserEmail")
		notificationService.AssertNotCalled(t, "SendOtpCodeMessage")
	})

	t.Run("it should fail and return ErrInvalidCredentials when passwords don't match", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		passwordService := mockService.NewPasswordServiceMock()
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		userPassword := "secret#pass1234"
		hashedPassword := "ixReNPXoBPxP9bIBQ6FziHj/9UG5wwzLbxP3vwpSZGo="

		user := domain.AuthUser{
			ID:       uuid.New(),
			Name:     "John Doe",
			Email:    userEmail,
			Password: hashedPassword,
		}

		userService.On("GetUserByEmail", ctx, userEmail).Return(&user, nil).Once()
		passwordService.On("Compare", hashedPassword, userPassword).Return(bcrypt.ErrMismatchedHashAndPassword).Once()

		useCase := NewUseCase(
			userService,
			passwordService,
			otpCodesRepository,
			notificationService,
		)

		err := useCase.Execute(ctx, userEmail, userPassword)

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		userService.AssertExpectations(t)
		passwordService.AssertExpectations(t)
		otpCodesRepository.AssertNotCalled(t, "CreateWithUserEmail")
		notificationService.AssertNotCalled(t, "SendOtpCodeMessage")
	})
}
