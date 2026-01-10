package otp

import (
	"comu/internal/modules/auth/domain"
	"context"
	"errors"
)

type VerifyOtpInput struct {
	UserEmail string
	OtpCodeType domain.OtpType
	OtpCodeValue string
}

type verifyOtpUC struct {
	otpCodesRepository domain.OtpCodesRepository
}

func NewVerifyOtpUseCase(otpCodesRepository domain.OtpCodesRepository) *verifyOtpUC {
	return &verifyOtpUC{
		otpCodesRepository: otpCodesRepository,
	}
}

// Execute try to retrieve the otp code with the provided otpCodeValue, check if it match the provided
// params and if it's not expired. If everything work fine, the otp code will be deleted from the data source
// and nil will be returned as a success value.
func (useCase *verifyOtpUC) Execute(ctx context.Context, input VerifyOtpInput) error {
	otpCode, err := useCase.otpCodesRepository.Find(ctx, input.OtpCodeValue)

	if err != nil {
		if errors.Is(err, domain.ErrOtpNotFound) {
			return domain.ErrInvalidOtp
		}

		return err
	}

	if otpCode.Type != input.OtpCodeType || otpCode.UserEmail != input.UserEmail {
		return domain.ErrInvalidOtp
	}

	if otpCode.Expired() {
		useCase.otpCodesRepository.Delete(ctx, otpCode)
		return domain.ErrExpiredOtp
	}
	useCase.otpCodesRepository.Delete(ctx, otpCode)

	return nil
}
