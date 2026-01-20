package otp

import (
	"comu/internal/modules/auth/domain"
	"context"
	"errors"
)

type VerifyOtpInput struct {
	UserEmail    string
	OtpCodeType  domain.OtpType
	OtpCodeValue string
}

type VerifyOtpUC struct {
	otpCodesRepository domain.OtpCodesRepository
	resendRequestsRepository domain.ResendOtpRequestsRepository
}

func NewVerifyOtpUseCase(
	otpCodesRepository domain.OtpCodesRepository,
	resendRequestsRepository domain.ResendOtpRequestsRepository,
	) *VerifyOtpUC {
	return &VerifyOtpUC{
		otpCodesRepository: otpCodesRepository,
		resendRequestsRepository: resendRequestsRepository,
	}
}

// Execute try to retrieve the otp code with the provided otpCodeValue, check if it match the provided
// params and if it's not expired. If everything work fine, the otp code will be deleted from the data source
// and nil will be returned as a success value.
func (useCase *VerifyOtpUC) Execute(ctx context.Context, input VerifyOtpInput) error {
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

	resendReq, _ := useCase.resendRequestsRepository.FindByUserEmail(ctx, otpCode.UserEmail)
	if resendReq != nil {
		useCase.resendRequestsRepository.Delete(ctx, resendReq)
	}
	useCase.otpCodesRepository.Delete(ctx, otpCode)

	return nil
}
