package handlers

import (
	"comu/internal/modules/auth/application/otp"
	resetPassword "comu/internal/modules/auth/application/reset_password"
	"comu/internal/modules/auth/application/tokens"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/presentation/validation"
	"comu/internal/shared/logger"
	echoRes "comu/internal/shared/utils/echo_res"
	"errors"

	"github.com/labstack/echo/v4"
)

type resetPasswordHandlers struct {
	newPasswordUC      *resetPassword.SetNewPasswordUC
	genResetTokenUC    *tokens.GenerateResetTokenUC
	resetPasswordUC    *resetPassword.ResetPasswordUC
	genResendRequestUC *otp.GenResendOtpRequestUC

	otpHandlers *otpHandlers
	logger      *logger.Log
}

func newResetPasswordHandlers(
	newPasswordUC *resetPassword.SetNewPasswordUC,
	genResetTokenUC *tokens.GenerateResetTokenUC,
	resetPasswordUC *resetPassword.ResetPasswordUC,
	genResendRequestUC *otp.GenResendOtpRequestUC,

	otpHandlers *otpHandlers,
	logger *logger.Log,
) *resetPasswordHandlers {
	return &resetPasswordHandlers{
		newPasswordUC:      newPasswordUC,
		genResetTokenUC:    genResetTokenUC,
		resetPasswordUC:    resetPasswordUC,
		genResendRequestUC: genResendRequestUC,

		otpHandlers: otpHandlers,
		logger:      logger,
	}
}

var msgPasswordsDoNotMatch = "Password and password confirmation don't match"

type resetPasswordFormData struct {
	Email string `form:"email" json:"email"`
}

type newPasswordFormData struct {
	ResetToken           string `form:"reset_token" json:"reset_token"`
	Password             string `form:"password" json:"password"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation"`
}

func (h *resetPasswordHandlers) reset(ctx echo.Context) error {
	var data resetPasswordFormData

	if err := ctx.Bind(&data); err != nil {
		return echoRes.JsonInvalidRequestResponse(ctx)
	}
	errList := validation.ResetPasswordValidator.Validate(&data)

	if errList != nil {
		return echoRes.JsonValidationErrorResponse(ctx, errList)
	}

	if err := h.resetPasswordUC.Execute(
		ctx.Request().Context(), data.Email,
	); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return echoRes.JsonSuccessMessageResponse(ctx, verificationSentMessage)
		}

		h.logger.Error.Println(err)
		return echoRes.JsonInternalErrorResponse(ctx)
	}

	resendRequest, _ := h.genResendRequestUC.Execute(ctx.Request().Context(), data.Email)

	return echoRes.JsonSuccessResponse(ctx, verificationSentMessage, map[string]string{
		"resend_token": resendRequest.ID.String(),
	})
}

func (h *resetPasswordHandlers) verifyOtp(ctx echo.Context) error {
	handler := h.otpHandlers.verify(domain.ResetPasswordOTP, func(validated verifyOtpFormData) error {
		token, err := h.genResetTokenUC.Execute(ctx.Request().Context(), validated.Email)

		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				return echoRes.JsonUnauthorizedResponse(ctx, invalidOtp, domain.ErrInvalidOtp.Error())
			}

			h.logger.Error.Println(err)
			return echoRes.JsonInternalErrorResponse(ctx)
		}

		return echoRes.JsonSuccessWithDataResponse(
			ctx, map[string]string{
				"reset_token": token,
			},
		)
	})

	return handler(ctx)
}

func (h *resetPasswordHandlers) resendOtp(ctx echo.Context) error {
	handler := h.otpHandlers.resend(domain.ResetPasswordOTP)
	return handler(ctx)
}

func (h *resetPasswordHandlers) newPassword(ctx echo.Context) error {
	var data newPasswordFormData

	if err := ctx.Bind(&data); err != nil {
		return echoRes.JsonInvalidRequestResponse(ctx)
	}
	errList := validation.NewPasswordValidator.Validate(&data)

	if errList != nil {
		return echoRes.JsonValidationErrorResponse(ctx, errList)
	}

	if data.Password != data.PasswordConfirmation {
		return echoRes.JsonValidationErrorResponse(
			ctx, map[string]string{"password": msgPasswordsDoNotMatch},
		)
	}

	if err := h.newPasswordUC.Execute(
		ctx.Request().Context(),
		data.ResetToken,
		data.Password,
	); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidToken):
			return echoRes.JsonUnauthorizedResponse(ctx, invalidToken, err.Error())
		case errors.Is(err, domain.ErrExpiredToken):
			return echoRes.JsonUnauthorizedResponse(ctx, expiredToken, err.Error())

		case errors.Is(err, domain.ErrUserNotFound):
			return echoRes.JsonUnauthorizedResponse(ctx, invalidToken, domain.ErrInvalidToken.Error())

		default:
			h.logger.Error.Println(err)
			return echoRes.JsonInternalErrorResponse(ctx)
		}
	}

	return echoRes.JsonSuccessMessageResponse(ctx, "Your password has been successfully updated.")
}

func (h *resetPasswordHandlers) RegisterRoutes(echo *echo.Echo, m ...echo.MiddlewareFunc) {
	groupRouter := echo.Group("/reset_password", m...)

	groupRouter.POST("", h.reset)
	groupRouter.POST("/verify", h.verifyOtp)
	groupRouter.POST("/resend_otp", h.resendOtp)
	groupRouter.POST("/new_password", h.newPassword)
}
