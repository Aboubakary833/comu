package register

import (
	"comu/internal/modules/auth/domain"
	mockRepository "comu/internal/modules/auth/mocks/mock_repository"
	mockService "comu/internal/modules/auth/mocks/mock_service"
	"comu/internal/modules/users"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	t.Run("it should result into success", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		passwordService := mockService.NewPasswordServiceMock()
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendOtpRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userName := "John Doe"
		userEmail := "johndoe@gmail.com"
		userPassword := "BhVmqUnb6m1upSh"
		hashedPassword := "ixReNPXoBPxP9bIBQ6FziHj/9UG5wwzLbxP3vwpSZGo="

		otpCode := domain.NewOtpCode(domain.RegisterOTP, userEmail, domain.DefaultOtpCodeTTL)

		passwordService.On("Hash", userPassword).Return(hashedPassword, nil).Once()
		userService.On("CreateNewUser", ctx, userName, userEmail, hashedPassword).Return(uuid.New(), nil).Once()
		otpCodesRepository.On("CreateWithUserEmail", ctx, domain.RegisterOTP, userEmail).Return(otpCode, nil).Once()
		notificationService.On("SendOtpCodeMessage", otpCode).Return(nil).Once()
		resendOtpRequestsRepository.On("CreateNew", ctx, userEmail).Return(nil).Once()

		useCase := NewRegisterUseCase(
			userService,
			passwordService,
			otpCodesRepository,
			notificationService,
			resendOtpRequestsRepository,
		)

		err := useCase.Execute(ctx, userName, userEmail, userPassword)

		assert.NoError(t, err)
		passwordService.AssertExpectations(t)
		userService.AssertExpectations(t)
		otpCodesRepository.AssertExpectations(t)
		notificationService.AssertExpectations(t)
		resendOtpRequestsRepository.AssertExpectations(t)
	})

	t.Run("it should fail and return email taken error", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		passwordService := mockService.NewPasswordServiceMock()
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendOtpRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userName := "John Doe"
		userEmail := "johndoe@gmail.com"
		userPassword := "BhVmqUnb6m1upSh"
		hashedPassword := "ixReNPXoBPxP9bIBQ6FziHj/9UG5wwzLbxP3vwpSZGo="

		passwordService.On("Hash", userPassword).Return(hashedPassword, nil).Once()
		userService.On("CreateNewUser", ctx, userName, userEmail, hashedPassword).Return(nil, users.ErrUserEmailTaken).Once()

		useCase := NewRegisterUseCase(
			userService,
			passwordService,
			otpCodesRepository,
			notificationService,
			resendOtpRequestsRepository,
		)

		err := useCase.Execute(ctx, userName, userEmail, userPassword)

		assert.ErrorIs(t, err, users.ErrUserEmailTaken)
		passwordService.AssertExpectations(t)
		userService.AssertExpectations(t)
		otpCodesRepository.AssertNotCalled(t, "CreateWithUserEmail")
		notificationService.AssertNotCalled(t, "SendOtpCodeMessage")
		resendOtpRequestsRepository.AssertNotCalled(t, "CreateNew")
	})
}
