package handlers

import (
	"comu/internal/modules/auth/application/login"
	"comu/internal/modules/auth/application/otp"
	"comu/internal/modules/auth/application/tokens"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/presentation/validation"
	"comu/internal/shared/logger"
	echoRes "comu/internal/shared/utils/echo_res"
	"errors"

	"github.com/labstack/echo/v4"
)

var verificationSentMessage = "A verification code has been sent to your mail."

var (
	invalidCredentials echoRes.ErrorResponseType = "invalid_credentials"
	invalidToken       echoRes.ErrorResponseType = "invalid_token"
	expiredToken       echoRes.ErrorResponseType = "expired_token"
)

type loginHandlers struct {
	loginUC                     *login.LoginUC
	genAuthTokenUC              *tokens.GenerateAuthTokensUC
	genResendRequestUC          *otp.GenResendOtpRequestUC
	genAccessTokenFromRefreshUC *tokens.GenAccessTokenFromRefreshUC

	otpHandlers *otpHandlers
	logger      *logger.Log
}

func newLoginHandlers(
	loginUC *login.LoginUC,
	genAuthTokenUC *tokens.GenerateAuthTokensUC,
	genResendRequestUC *otp.GenResendOtpRequestUC,
	genAccessTokenFromRefreshUC *tokens.GenAccessTokenFromRefreshUC,

	otpHandler *otpHandlers,
	logger *logger.Log,
) *loginHandlers {
	return &loginHandlers{
		loginUC:                     loginUC,
		genAuthTokenUC:              genAuthTokenUC,
		genResendRequestUC:          genResendRequestUC,
		genAccessTokenFromRefreshUC: genAccessTokenFromRefreshUC,

		otpHandlers: otpHandler,
		logger:      logger,
	}
}

type loginFormData struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

type refreshFormData struct {
	Token string `form:"refresh_token" json:"refresh_token"`
}

func (h *loginHandlers) loginAttempt(ctx echo.Context) error {
	var data loginFormData

	if err := ctx.Bind(&data); err != nil {
		return echoRes.JsonInvalidRequestResponse(ctx)
	}

	if errList := validation.LoginValidator.Validate(&data); errList != nil {
		return echoRes.JsonValidationErrorResponse(ctx, errList)
	}

	if err := h.loginUC.Execute(
		ctx.Request().Context(),
		data.Email,
		data.Password,
	); err != nil {

		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			return echoRes.JsonUnauthorizedResponse(ctx, invalidCredentials, err.Error())
		default:
			h.logger.Error.Println(err)
			return echoRes.JsonInternalErrorResponse(ctx)
		}
	}

	resendRequest, _ := h.genResendRequestUC.Execute(ctx.Request().Context(), data.Email)

	return echoRes.JsonSuccessResponse(ctx, verificationSentMessage, map[string]string{
		"resend_token": resendRequest.ID.String(),
	})
}

func (h *loginHandlers) verifyOtp(ctx echo.Context) error {
	handler := h.otpHandlers.verify(domain.LoginOTP, func(validated verifyOtpFormData) error {
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

func (h *loginHandlers) resendOtp(ctx echo.Context) error {
	handler := h.otpHandlers.resend(domain.LoginOTP)
	return handler(ctx)
}

func (h *loginHandlers) refreshToken(ctx echo.Context) error {
	var data refreshFormData

	if err := ctx.Bind(&data); err != nil {
		return echoRes.JsonInvalidRequestResponse(ctx)
	}

	if data.Token == "" {
		return echoRes.JsonUnauthorizedResponse(ctx, invalidToken, domain.ErrInvalidToken.Error())
	}

	token, err := h.genAccessTokenFromRefreshUC.Execute(
		ctx.Request().Context(),
		data.Token,
	)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrExpiredToken):
			return echoRes.JsonUnauthorizedResponse(ctx, expiredToken, err.Error())

		case errors.Is(err, domain.ErrUserEmailTaken):
			return echoRes.JsonUnauthorizedResponse(
				ctx, invalidToken,
				domain.ErrInvalidToken.Error(),
			)

		default:
			h.logger.Error.Println(err)
			return echoRes.JsonInternalErrorResponse(ctx)
		}
	}

	return echoRes.JsonSuccessWithDataResponse(ctx, map[string]string{
		"access_token": token,
	})
}

func (h *loginHandlers) RegisterRoutes(echo *echo.Echo, m ...echo.MiddlewareFunc) {
	groupRouter := echo.Group("/login", m...)

	groupRouter.POST("", h.loginAttempt)
	groupRouter.POST("/verify", h.verifyOtp)
	groupRouter.POST("/resend_otp", h.resendOtp)
	groupRouter.POST("/refresh", h.refreshToken)
}
