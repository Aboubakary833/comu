package resetPassword

import (
	"comu/internal/modules/auth/domain"
	mockRepository "comu/internal/modules/auth/mocks/mock_repository"
	mockService "comu/internal/modules/auth/mocks/mock_service"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestResetPasswordUseCase(t *testing.T) {

	t.Run("it should fail and return ErrUserNotFound", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		otpCodeRepository := mockRepository.NewOtpCodesRepositoryMock()
		notificationService := mockService.NewNotificationServiceMock()
		resendOtpRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"

		userService.On("GetUserByEmail", ctx, userEmail).Return(nil, domain.ErrUserNotFound).Once()

		useCase := NewResetPasswordUseCase(userService, otpCodeRepository, notificationService, resendOtpRequestsRepository)

		err := useCase.Execute(ctx, userEmail)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		userService.AssertExpectations(t)
		otpCodeRepository.AssertNotCalled(t, "CreateWithUserEmail")
		notificationService.AssertNotCalled(t, "SendOtpCodeMessage")
		resendOtpRequestsRepository.AssertNotCalled(t, "CreateNew")
	})

	t.Run("it should result into success and send reset code", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		otpCodeRepository := mockRepository.NewOtpCodesRepositoryMock()
		notificationService := mockService.NewNotificationServiceMock()
		resendOtpRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"

		user := &domain.AuthUser{
			ID:       uuid.New(),
			Name:     "John Doe",
			Email:    userEmail,
			Password: "secret#pass1234",
		}
		otpCode := domain.NewOtpCode(domain.RegisterOTP, userEmail, domain.DefaultOtpCodeTTL)

		userService.On("GetUserByEmail", ctx, userEmail).Return(user, nil).Once()
		otpCodeRepository.On("CreateWithUserEmail", ctx, domain.ResetPasswordOTP, userEmail).Return(otpCode, nil).Once()
		notificationService.On("SendOtpCodeMessage", otpCode).Return(nil).Once()
		resendOtpRequestsRepository.On("CreateNew", ctx, userEmail).Return(nil).Once()

		useCase := NewResetPasswordUseCase(userService, otpCodeRepository, notificationService, resendOtpRequestsRepository)

		err := useCase.Execute(ctx, userEmail)

		assert.NoError(t, err)
		userService.AssertExpectations(t)
		otpCodeRepository.AssertExpectations(t)
		notificationService.AssertExpectations(t)
	})

}
