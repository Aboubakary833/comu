package handlers

import (
	"comu/internal/modules/auth/application/otp"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/presentation/validation"
	"comu/internal/shared/logger"
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	invalidOtp                 errorType = "invalid_otp"
	expiredOtp                 errorType = "expired_otp"
	invalidResendRequest       errorType = "invalid_request"
	unprocessableResendRequest errorType = "unprocessable_resend_request"
	exceededResendRequestCount errorType = "exceeded_resend_request_count"
)

type OtpHandlers struct {
	verifyOtpUC otp.VerifyOtpUC
	resendOtpUC otp.ResendOtpUC

	logger *logger.Log
}

func NewOtpHandlers(
	verifyOtpUC otp.VerifyOtpUC,
	resendOtpUC otp.ResendOtpUC,
	logger *logger.Log,
) *OtpHandlers {
	return &OtpHandlers{
		verifyOtpUC: verifyOtpUC,
		resendOtpUC: resendOtpUC,

		logger: logger,
	}
}

type verifyOtpFormData struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type resendOtpFormData struct {
	Email       string `json:"email"`
	ResendToken string `json:"resend_token"`
}

func (h *OtpHandlers) Verify(otpType domain.OtpType, afterFunc func(verifyOtpFormData) error) echo.HandlerFunc {

	return func(ctx echo.Context) error {
		var data, validated verifyOtpFormData

		if err := ctx.Bind(&data); err != nil {
			return jsonInvalidRequestResponse(ctx)
		}

		if errList := validation.OtpCodeValidator.Validate(data, &validated); errList != nil {
			return jsonValidationErrorResponse(ctx, errList)
		}

		if err := h.verifyOtpUC.Execute(
			ctx.Request().Context(), otp.VerifyOtpInput{
				UserEmail:    validated.Email,
				OtpCodeType:  otpType,
				OtpCodeValue: validated.Code,
			},
		); err != nil {

			switch {
			case errors.Is(err, domain.ErrInvalidOtp):
				return jsonUnauthorizedResponse(ctx, invalidOtp, err.Error())
			case errors.Is(err, domain.ErrExpiredOtp):
				return jsonUnauthorizedResponse(ctx, expiredOtp, err.Error())
			default:
				h.logger.Error.Println(err)
				return jsonInternalErrorResponse(ctx)
			}
		}

		return afterFunc(validated)
	}

}

func (h *OtpHandlers) Resend(otpType domain.OtpType) echo.HandlerFunc {

	return func(ctx echo.Context) error {
		var data, validated resendOtpFormData

		if err := ctx.Bind(&data); err != nil {
			return jsonInvalidRequestResponse(ctx)
		}

		if errList := validation.ResendOtpValidator.Validate(data, &validated); errList != nil {
			return jsonValidationErrorResponse(ctx, errList)
		}

		id, err := uuid.Parse(validated.ResendToken)

		if err != nil {
			return jsonUnauthorizedResponse(
				ctx, invalidResendRequest,
				domain.ErrInvalidResendRequest.Error(),
			)
		}

		if err := h.resendOtpUC.Execute(
			ctx.Request().Context(), otp.ResendOtpInput{
				ID:          id,
				UserEmail:   validated.Email,
				OtpCodeType: otpType,
			},
		); err != nil {

			switch {
			case errors.Is(err, domain.ErrInvalidResendRequest):
				return jsonUnauthorizedResponse(ctx, invalidResendRequest, err.Error())

			case errors.Is(err, domain.ErrResendRequestCountExceeded):
				return jsonUnauthorizedResponse(ctx, exceededResendRequestCount, err.Error())

			case errors.Is(err, domain.ErrResendRequestCantBeProcessed):
				return jsonUnauthorizedResponse(ctx, unprocessableResendRequest, err.Error())

			default:
				h.logger.Error.Println(err)
				return jsonInternalErrorResponse(ctx)
			}

		}

		return jsonSuccessMessageResponse(ctx, "A new otp code has been sent to your mail.")
	}
}
