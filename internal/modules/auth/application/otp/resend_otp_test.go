package otp

import (
	"comu/internal/modules/auth/domain"
	mockRepository "comu/internal/modules/auth/mocks/mock_repository"
	mockService "comu/internal/modules/auth/mocks/mock_service"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestResendOtpUseCase(t *testing.T) {

	t.Run("it should fail and return ErrResendRequestNotFound", func(t *testing.T) {
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendOtpRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		reqID := uuid.New()
		userEmail := "johndoe@gmail.com"

		resendOtpRequestsRepository.On("FindByID", ctx, reqID).Return(nil, domain.ErrResendRequestNotFound).Once()

		useCase := NewResendOtpUseCase(
			otpCodesRepository,
			notificationService,
			resendOtpRequestsRepository,
		)

		err := useCase.Execute(
			ctx, ResendOtpInput{
				ID:          reqID,
				UserEmail:   userEmail,
				OtpCodeType: domain.LoginOTP,
			},
		)

		assert.ErrorIs(t, err, domain.ErrResendRequestNotFound)
		resendOtpRequestsRepository.AssertExpectations(t)
		otpCodesRepository.AssertNotCalled(t, "FindByUserEmail")
		notificationService.AssertNotCalled(t, "SendOtpCodeMessage")
	})

	t.Run("it should fail and return ErrInvalidResendRequest when otp code is'nt found", func(t *testing.T) {
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendOtpRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		req := domain.NewResendOtpRequest(userEmail)

		resendOtpRequestsRepository.On("FindByID", ctx, req.ID).Return(req, nil).Once()
		otpCodesRepository.On("FindByUserEmail", ctx, userEmail).Return(nil, domain.ErrOtpNotFound).Once()

		useCase := NewResendOtpUseCase(
			otpCodesRepository,
			notificationService,
			resendOtpRequestsRepository,
		)

		err := useCase.Execute(
			ctx, ResendOtpInput{
				ID:          req.ID,
				UserEmail:   userEmail,
				OtpCodeType: domain.LoginOTP,
			},
		)

		assert.ErrorIs(t, err, domain.ErrInvalidResendRequest)
		resendOtpRequestsRepository.AssertExpectations(t)
		otpCodesRepository.AssertExpectations(t)
		otpCodesRepository.AssertNotCalled(t, "Delete")
		otpCodesRepository.AssertNotCalled(t, "CreateWithUserEmail")
		resendOtpRequestsRepository.AssertNotCalled(t, "Update")
		notificationService.AssertNotCalled(t, "SendOtpCodeMessage")
	})

	t.Run("it should fail and return ErrInvalidResendRequest when otp code type don't match input type", func(t *testing.T) {
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendOtpRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		otpCode := domain.NewOtpCode(domain.RegisterOTP, userEmail, domain.DefaultOtpCodeTTL)
		req := domain.NewResendOtpRequest(userEmail)

		resendOtpRequestsRepository.On("FindByID", ctx, req.ID).Return(req, nil).Once()
		otpCodesRepository.On("FindByUserEmail", ctx, userEmail).Return(otpCode, nil).Once()

		useCase := NewResendOtpUseCase(
			otpCodesRepository,
			notificationService,
			resendOtpRequestsRepository,
		)

		err := useCase.Execute(
			ctx, ResendOtpInput{
				ID:          req.ID,
				UserEmail:   userEmail,
				OtpCodeType: domain.LoginOTP,
			},
		)

		assert.ErrorIs(t, err, domain.ErrInvalidResendRequest)
		resendOtpRequestsRepository.AssertExpectations(t)
		otpCodesRepository.AssertExpectations(t)
		otpCodesRepository.AssertNotCalled(t, "Delete")
		otpCodesRepository.AssertNotCalled(t, "CreateWithUserEmail")
		resendOtpRequestsRepository.AssertNotCalled(t, "Update")
		notificationService.AssertNotCalled(t, "SendOtpCodeMessage")
	})

	t.Run("it should fail and return ErrResendRequestCantBeProcessed", func(t *testing.T) {
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendOtpRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		otpCode := domain.NewOtpCode(domain.LoginOTP, userEmail, domain.DefaultOtpCodeTTL)
		req := domain.NewResendOtpRequest(userEmail)
		req.LastSendAt = time.Now().Add(time.Minute * 5)

		resendOtpRequestsRepository.On("FindByID", ctx, req.ID).Return(req, nil).Once()
		otpCodesRepository.On("FindByUserEmail", ctx, userEmail).Return(otpCode, nil).Once()

		useCase := NewResendOtpUseCase(
			otpCodesRepository,
			notificationService,
			resendOtpRequestsRepository,
		)

		err := useCase.Execute(
			ctx, ResendOtpInput{
				ID:          req.ID,
				UserEmail:   userEmail,
				OtpCodeType: domain.LoginOTP,
			},
		)

		assert.ErrorIs(t, err, domain.ErrResendRequestCantBeProcessed)
		resendOtpRequestsRepository.AssertExpectations(t)
		otpCodesRepository.AssertExpectations(t)
		otpCodesRepository.AssertNotCalled(t, "Delete")
		otpCodesRepository.AssertNotCalled(t, "CreateWithUserEmail")
		notificationService.AssertNotCalled(t, "SendOtpCodeMessage")
		resendOtpRequestsRepository.AssertNotCalled(t, "Update")
	})

	t.Run("it should return ErrResendRequestCountExceeded", func(t *testing.T) {
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendOtpRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		otpCode := domain.NewOtpCode(domain.LoginOTP, userEmail, domain.DefaultOtpCodeTTL)
		req := domain.NewResendOtpRequest(userEmail)
		req.Count = 5
		req.LastSendAt = time.Now().Add(-5 * time.Minute)

		resendOtpRequestsRepository.On("FindByID", ctx, req.ID).Return(req, nil).Once()
		otpCodesRepository.On("FindByUserEmail", ctx, userEmail).Return(otpCode, nil).Once()

		useCase := NewResendOtpUseCase(
			otpCodesRepository,
			notificationService,
			resendOtpRequestsRepository,
		)

		err := useCase.Execute(
			ctx, ResendOtpInput{
				ID:          req.ID,
				UserEmail:   userEmail,
				OtpCodeType: domain.LoginOTP,
			},
		)

		assert.ErrorIs(t, err, domain.ErrResendRequestCountExceeded)
		resendOtpRequestsRepository.AssertExpectations(t)
		otpCodesRepository.AssertExpectations(t)
		otpCodesRepository.AssertNotCalled(t, "Delete")
		otpCodesRepository.AssertNotCalled(t, "CreateWithUserEmail")
		notificationService.AssertNotCalled(t, "SendOtpCodeMessage")
		resendOtpRequestsRepository.AssertNotCalled(t, "Update")
	})

	t.Run("it should succeed", func(t *testing.T) {
		notificationService := mockService.NewNotificationServiceMock()
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendOtpRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()
		_assert := assert.New(t)

		userEmail := "johndoe@gmail.com"
		otpCode := domain.NewOtpCode(domain.LoginOTP, userEmail, domain.DefaultOtpCodeTTL)
		req := domain.NewResendOtpRequest(userEmail)
		req.LastSendAt = time.Now().Add(-5 * time.Minute)

		resendOtpRequestsRepository.On("FindByID", ctx, req.ID).Return(req, nil).Once()
		resendOtpRequestsRepository.On("Update", ctx, req).Return(nil).Once()
		otpCodesRepository.On("FindByUserEmail", ctx, userEmail).Return(otpCode, nil).Once()
		otpCodesRepository.On("Delete", ctx, otpCode).Return(nil)
		otpCodesRepository.On("CreateWithUserEmail", ctx, otpCode.Type, otpCode.UserEmail).Return(otpCode, nil)
		notificationService.On("SendOtpCodeMessage", otpCode).Return(nil).Once()

		useCase := NewResendOtpUseCase(
			otpCodesRepository,
			notificationService,
			resendOtpRequestsRepository,
		)

		err := useCase.Execute(
			ctx, ResendOtpInput{
				ID:          req.ID,
				UserEmail:   userEmail,
				OtpCodeType: domain.LoginOTP,
			},
		)

		if _assert.NoError(err) {
			resendOtpRequestsRepository.AssertExpectations(t)
			otpCodesRepository.AssertExpectations(t)
			notificationService.AssertExpectations(t)

			_assert.Equal(2, req.Count)
			_assert.Equal(false, req.CanOtpBeSent())
		}
	})
}
