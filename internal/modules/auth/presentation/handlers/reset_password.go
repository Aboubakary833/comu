package handlers

import (
	resetPassword "comu/internal/modules/auth/application/reset_password"
	"comu/internal/modules/auth/application/tokens"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/presentation/validation"
	"comu/internal/shared/logger"
	"errors"

	"github.com/labstack/echo/v4"
)

type ResetPasswordHandlers struct {
	newPasswordUC   resetPassword.SetNewPasswordUC
	genResetTokenUC tokens.GenerateResetTokenUC
	resetPasswordUC resetPassword.ResetPasswordUC

	otpHandlers *OtpHandlers
	logger      *logger.Log
}

func NewResetPasswordHandlers(
	newPasswordUC resetPassword.SetNewPasswordUC,
	genResetTokenUC tokens.GenerateResetTokenUC,
	resetPasswordUC resetPassword.ResetPasswordUC,

	otpHandlers *OtpHandlers,
	logger *logger.Log,
) *ResetPasswordHandlers {
	return &ResetPasswordHandlers{
		newPasswordUC:   newPasswordUC,
		genResetTokenUC: genResetTokenUC,
		resetPasswordUC: resetPasswordUC,

		otpHandlers: otpHandlers,
		logger:      logger,
	}
}

type resetPasswordFormData struct {
	Email string `json:"email"`
}

type newPasswordFormData struct {
	ResetToken           string `json:"reset_token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (h *ResetPasswordHandlers) Reset(ctx echo.Context) error {
	var data, validated resetPasswordFormData

	if err := ctx.Bind(&data); err != nil {
		return jsonInvalidRequestResponse(ctx)
	}
	errList := validation.ResetPasswordValidator.Validate(data, &validated)

	if errList != nil {
		return jsonValidationErrorResponse(ctx, errList)
	}

	if err := h.resetPasswordUC.Execute(
		ctx.Request().Context(), validated.Email,
	); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return jsonSuccessMessageResponse(ctx, verificationSentMessage)
		}

		h.logger.Error.Println(err)
		return jsonInternalErrorResponse(ctx)
	}

	return jsonSuccessMessageResponse(ctx, verificationSentMessage)
}

func (h *ResetPasswordHandlers) VerifyOtp(ctx echo.Context) error {
	handler := h.otpHandlers.Verify(domain.LoginOTP, func(validated verifyOtpFormData) error {
		token, err := h.genResetTokenUC.Execute(ctx.Request().Context(), validated.Email)

		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				return jsonUnauthorizedResponse(ctx, invalidOtp, domain.ErrInvalidOtp.Error())
			}

			h.logger.Error.Println(err)
			return jsonInternalErrorResponse(ctx)
		}

		return jsonSuccessWithDataResponse(
			ctx, map[string]string{
				"reset_token": token,
			},
		)
	})

	return handler(ctx)
}

func (h *ResetPasswordHandlers) ResendOtp(ctx echo.Context) error {
	handler := h.otpHandlers.Resend(domain.ResetPasswordOTP)
	return handler(ctx)
}

func (h *ResetPasswordHandlers) NewPassword(ctx echo.Context) error {
	var data, validated newPasswordFormData

	if err := ctx.Bind(&data); err != nil {
		return jsonInvalidRequestResponse(ctx)
	}
	errList := validation.NewPasswordValidator.Validate(data, &validated)

	if errList != nil {
		return jsonValidationErrorResponse(ctx, errList)
	}

	if err := h.newPasswordUC.Execute(
		ctx.Request().Context(),
		validated.ResetToken,
		validated.Password,
	); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidToken):
			return jsonUnauthorizedResponse(ctx, invalidToken, err.Error())
		case errors.Is(err, domain.ErrExpiredToken):
			return jsonUnauthorizedResponse(ctx, expiredToken, err.Error())

		case errors.Is(err, domain.ErrUserNotFound):
			return jsonUnauthorizedResponse(ctx, invalidToken, domain.ErrInvalidToken.Error())

		default:
			h.logger.Error.Println(err)
			return jsonInternalErrorResponse(ctx)
		}
	}

	return jsonSuccessMessageResponse(ctx, "Your password has been successfully updated.")
}
