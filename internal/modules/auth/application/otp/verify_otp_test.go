package otp

import (
	"comu/internal/modules/auth/domain"
	mockRepository "comu/internal/modules/auth/mocks/mock_repository"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVerifyOtpUseCase(t *testing.T) {

	t.Run("it should return ErrInvalidOtp when otpCode not found", func(t *testing.T) {
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		otpCodeValue := "012345"
		userEmail := "johndoe@gmail.com"

		otpCodesRepository.On("Find", ctx, otpCodeValue).Return(nil, domain.ErrOtpNotFound).Once()

		useCase := NewVerifyOtpUseCase(otpCodesRepository, resendRequestsRepository)

		err := useCase.Execute(
			ctx, VerifyOtpInput{
				UserEmail:    userEmail,
				OtpCodeType:  domain.LoginOTP,
				OtpCodeValue: otpCodeValue,
			},
		)

		assert.ErrorIs(t, err, domain.ErrInvalidOtp)
		otpCodesRepository.AssertExpectations(t)
		resendRequestsRepository.AssertNotCalled(t, "FindByUserEmail")
		resendRequestsRepository.AssertNotCalled(t, "Delete")
	})

	t.Run("it should return ErrInvalidOtp when provided email does'nt match retrieved otp email", func(t *testing.T) {
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		otpCode := domain.NewOtpCode(domain.LoginOTP, "jeannettedoe@gmail.com", domain.DefaultOtpCodeTTL)
		userEmail := "johndoe@gmail.com"

		otpCodesRepository.On("Find", ctx, otpCode.Value).Return(otpCode, nil).Once()

		useCase := NewVerifyOtpUseCase(otpCodesRepository, resendRequestsRepository)

		err := useCase.Execute(
			ctx, VerifyOtpInput{
				UserEmail:    userEmail,
				OtpCodeType:  domain.LoginOTP,
				OtpCodeValue: otpCode.Value,
			},
		)

		assert.ErrorIs(t, err, domain.ErrInvalidOtp)
		otpCodesRepository.AssertExpectations(t)
		resendRequestsRepository.AssertNotCalled(t, "FindByUserEmail")
		resendRequestsRepository.AssertNotCalled(t, "Delete")
	})

	t.Run("it should return ErrInvalidOtp when provided type does'nt match retrieved otp type", func(t *testing.T) {
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		otpCode := domain.NewOtpCode(domain.RegisterOTP, userEmail, domain.DefaultOtpCodeTTL)

		otpCodesRepository.On("Find", ctx, otpCode.Value).Return(otpCode, nil).Once()

		useCase := NewVerifyOtpUseCase(otpCodesRepository, resendRequestsRepository)

		err := useCase.Execute(
			ctx, VerifyOtpInput{
				UserEmail:    userEmail,
				OtpCodeType:  domain.LoginOTP,
				OtpCodeValue: otpCode.Value,
			},
		)

		assert.ErrorIs(t, err, domain.ErrInvalidOtp)
		otpCodesRepository.AssertExpectations(t)
		resendRequestsRepository.AssertNotCalled(t, "FindByUserEmail")
		resendRequestsRepository.AssertNotCalled(t, "Delete")
	})

	t.Run("it should return ErrExpiredOtp when retrieved otp expired", func(t *testing.T) {
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		otpCode := domain.NewOtpCode(domain.RegisterOTP, userEmail, -2*time.Minute)

		otpCodesRepository.On("Find", ctx, otpCode.Value).Return(otpCode, nil).Once()
		otpCodesRepository.On("Delete", ctx, otpCode).Return(nil).Once()

		useCase := NewVerifyOtpUseCase(otpCodesRepository, resendRequestsRepository)

		err := useCase.Execute(
			ctx, VerifyOtpInput{
				UserEmail:    userEmail,
				OtpCodeType:  domain.RegisterOTP,
				OtpCodeValue: otpCode.Value,
			},
		)

		assert.ErrorIs(t, err, domain.ErrExpiredOtp)
		otpCodesRepository.AssertExpectations(t)
		resendRequestsRepository.AssertNotCalled(t, "FindByUserEmail")
		resendRequestsRepository.AssertNotCalled(t, "Delete")
	})

	t.Run("it should succeed and delete otp and resend request entries from repos", func(t *testing.T) {
		otpCodesRepository := mockRepository.NewOtpCodesRepositoryMock()
		resendRequestsRepository := mockRepository.NewResendOtpRequestsRepositoryMock()
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		otpCode := domain.NewOtpCode(domain.RegisterOTP, userEmail, domain.DefaultOtpCodeTTL)
		resendReq := domain.NewResendOtpRequest(userEmail)

		otpCodesRepository.On("Find", ctx, otpCode.Value).Return(otpCode, nil).Once()
		otpCodesRepository.On("Delete", ctx, otpCode).Return(nil).Once()
		resendRequestsRepository.On("FindByUserEmail", ctx, userEmail).Return(resendReq, nil).Once()
		resendRequestsRepository.On("Delete", ctx, resendReq).Return(nil)

		useCase := NewVerifyOtpUseCase(otpCodesRepository, resendRequestsRepository)

		err := useCase.Execute(
			ctx, VerifyOtpInput{
				UserEmail:    userEmail,
				OtpCodeType:  domain.RegisterOTP,
				OtpCodeValue: otpCode.Value,
			},
		)

		assert.NoError(t, err)
		otpCodesRepository.AssertExpectations(t)
		resendRequestsRepository.AssertExpectations(t)
	})
}
