package handlers

import (
	"comu/internal/modules/auth/application/otp"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/presentation/validation"
	"comu/internal/shared/logger"
	echoRes "comu/internal/shared/utils/echo_res"
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	invalidOtp                 echoRes.ErrorResponseType = "invalid_otp"
	expiredOtp                 echoRes.ErrorResponseType = "expired_otp"
	invalidResendRequest       echoRes.ErrorResponseType = "invalid_request"
	unprocessableResendRequest echoRes.ErrorResponseType = "unprocessable_resend_request"
	exceededResendRequestCount echoRes.ErrorResponseType = "exceeded_resend_request_count"
)

type otpHandlers struct {
	verifyOtpUC *otp.VerifyOtpUC
	resendOtpUC *otp.ResendOtpUC

	logger *logger.Log
}

func newOtpHandlers(
	verifyOtpUC *otp.VerifyOtpUC,
	resendOtpUC *otp.ResendOtpUC,
	logger *logger.Log,
) *otpHandlers {
	return &otpHandlers{
		verifyOtpUC: verifyOtpUC,
		resendOtpUC: resendOtpUC,

		logger: logger,
	}
}

type verifyOtpFormData struct {
	Email string `form:"email" json:"email"`
	Code  string `form:"code" json:"code"`
}

type resendOtpFormData struct {
	Email       string `form:"email" json:"email"`
	ResendToken string `form:"resend_token" json:"resend_token"`
}

func (h *otpHandlers) verify(otpType domain.OtpType, afterFunc func(verifyOtpFormData) error) echo.HandlerFunc {

	return func(ctx echo.Context) error {
		var data verifyOtpFormData

		if err := ctx.Bind(&data); err != nil {
			return echoRes.JsonInvalidRequestResponse(ctx)
		}

		if errList := validation.OtpCodeValidator.Validate(&data); errList != nil {
			return echoRes.JsonValidationErrorResponse(ctx, errList)
		}

		if err := h.verifyOtpUC.Execute(
			ctx.Request().Context(), otp.VerifyOtpInput{
				UserEmail:    data.Email,
				OtpCodeType:  otpType,
				OtpCodeValue: data.Code,
			},
		); err != nil {

			switch {
			case errors.Is(err, domain.ErrInvalidOtp):
				return echoRes.JsonUnauthorizedResponse(ctx, invalidOtp, err.Error())
			case errors.Is(err, domain.ErrExpiredOtp):
				return echoRes.JsonUnauthorizedResponse(ctx, expiredOtp, err.Error())
			default:
				h.logger.Error.Println(err)
				return echoRes.JsonInternalErrorResponse(ctx)
			}
		}

		return afterFunc(data)
	}

}

func (h *otpHandlers) resend(otpType domain.OtpType) echo.HandlerFunc {

	return func(ctx echo.Context) error {
		var data resendOtpFormData

		if err := ctx.Bind(&data); err != nil {
			return echoRes.JsonInvalidRequestResponse(ctx)
		}

		h.logger.Info.Println(data)
		if errList := validation.ResendOtpValidator.Validate(&data); errList != nil {
			return echoRes.JsonValidationErrorResponse(ctx, errList)
		}

		id, err := uuid.Parse(data.ResendToken)

		if err != nil {
			return echoRes.JsonUnauthorizedResponse(
				ctx, invalidResendRequest,
				domain.ErrInvalidResendRequest.Error(),
			)
		}

		if err := h.resendOtpUC.Execute(
			ctx.Request().Context(), otp.ResendOtpInput{
				ID:          id,
				UserEmail:   data.Email,
				OtpCodeType: otpType,
			},
		); err != nil {

			switch {

			case errors.Is(err, domain.ErrResendRequestNotFound):
				return echoRes.JsonUnauthorizedResponse(
					ctx, invalidResendRequest,
					domain.ErrInvalidResendRequest.Error(),
				)

			case errors.Is(err, domain.ErrInvalidResendRequest):
				return echoRes.JsonUnauthorizedResponse(ctx, invalidResendRequest, err.Error())

			case errors.Is(err, domain.ErrResendRequestCountExceeded):
				return echoRes.JsonUnauthorizedResponse(ctx, exceededResendRequestCount, err.Error())

			case errors.Is(err, domain.ErrResendRequestCantBeProcessed):
				return echoRes.JsonUnauthorizedResponse(ctx, unprocessableResendRequest, err.Error())

			default:
				h.logger.Error.Println(err)
				return echoRes.JsonInternalErrorResponse(ctx)
			}

		}

		return echoRes.JsonSuccessMessageResponse(ctx, "A new otp code has been sent to your mail.")
	}
}
