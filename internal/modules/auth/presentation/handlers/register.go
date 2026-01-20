package handlers

import (
	"comu/internal/modules/auth/application/otp"
	"comu/internal/modules/auth/application/register"
	"comu/internal/modules/auth/application/tokens"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/presentation/validation"
	"comu/internal/shared/logger"
	echoRes "comu/internal/shared/utils/echo_res"
	"errors"

	"github.com/labstack/echo/v4"
)

var (
	userEmailTaken echoRes.ErrorResponseType = "user_email_taken"
)

type registerHandlers struct {
	registerUC           *register.RegisterUC
	genAuthTokenUC       *tokens.GenerateAuthTokensUC
	markUserAsVerifiedUC *register.MarkUserAsVerifiedUC
	genResendRequestUC   *otp.GenResendOtpRequestUC

	otpHandlers *otpHandlers
	logger      *logger.Log
}

func newRegisterHandlers(
	registerUC *register.RegisterUC,
	genAuthTokenUC *tokens.GenerateAuthTokensUC,
	markUserAsVerifiedUC *register.MarkUserAsVerifiedUC,
	genResendRequestUC *otp.GenResendOtpRequestUC,

	otpHandler *otpHandlers,
	logger *logger.Log,
) *registerHandlers {
	return &registerHandlers{
		registerUC:           registerUC,
		genAuthTokenUC:       genAuthTokenUC,
		markUserAsVerifiedUC: markUserAsVerifiedUC,
		genResendRequestUC:   genResendRequestUC,

		otpHandlers: otpHandler,
		logger:      logger,
	}
}

type registerFormData struct {
	Name     string `form:"name" json:"name"`
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (h *registerHandlers) register(ctx echo.Context) error {
	var data registerFormData

	if err := ctx.Bind(&data); err != nil {
		return echoRes.JsonInvalidRequestResponse(ctx)
	}
	errList := validation.RegisterValidator.Validate(&data)

	if errList != nil {
		return echoRes.JsonValidationErrorResponse(ctx, errList)
	}

	if err := h.registerUC.Execute(
		ctx.Request().Context(),
		data.Name, data.Email, data.Password,
	); err != nil {

		if errors.Is(err, domain.ErrUserEmailTaken) {
			return echoRes.JsonUnauthorizedResponse(ctx, userEmailTaken, err.Error())
		}

		h.logger.Error.Println(err)
		return echoRes.JsonInternalErrorResponse(ctx)
	}

	resendRequest, _ := h.genResendRequestUC.Execute(ctx.Request().Context(), data.Email)

	return echoRes.JsonSuccessResponse(ctx, verificationSentMessage, map[string]string{
		"resend_token": resendRequest.ID.String(),
	})
}

func (h *registerHandlers) verifyOtp(ctx echo.Context) error {
	handler := h.otpHandlers.verify(domain.RegisterOTP, func(validated verifyOtpFormData) error {

		if err := h.markUserAsVerifiedUC.Execute(ctx.Request().Context(), validated.Email); err != nil {
			if errors.Is(err, domain.ErrUserEmailTaken) {
				return echoRes.JsonUnauthorizedResponse(ctx, invalidOtp, domain.ErrInvalidOtp.Error())
			}

			h.logger.Error.Println(err)
			return echoRes.JsonInternalErrorResponse(ctx)
		}

		access, refresh, err := h.genAuthTokenUC.Execute(ctx.Request().Context(), validated.Email)

		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				return echoRes.JsonUnauthorizedResponse(ctx, invalidOtp, domain.ErrInvalidOtp.Error())
			}

			h.logger.Error.Println(err)
			return echoRes.JsonInternalErrorResponse(ctx)
		}

		return echoRes.JsonSuccessWithDataResponse(ctx, map[string]string{
			"access_token":  access,
			"refresh_token": refresh,
		})
	})

	return handler(ctx)
}

func (h *registerHandlers) resendOtp(ctx echo.Context) error {
	handler := h.otpHandlers.resend(domain.RegisterOTP)
	return handler(ctx)
}

func (h *registerHandlers) RegisterRoutes(echo *echo.Echo, m ...echo.MiddlewareFunc) {
	groupRouter := echo.Group("/register", m...)

	groupRouter.POST("", h.register)
	groupRouter.POST("/verify", h.verifyOtp)
	groupRouter.POST("/resend_otp", h.resendOtp)
}
