package handlers

import (
	"comu/internal/modules/auth/application/login"
	"comu/internal/modules/auth/application/tokens"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/presentation/validation"
	"comu/internal/shared/logger"
	"errors"

	"github.com/labstack/echo/v4"
)

var verificationSentMessage = "A verification code has been sent to your mail."

type LoginHandlers struct {
	loginUC                     login.LoginUC
	genAuthTokenUC              tokens.GenerateAuthTokensUC
	genAccessTokenFromRefreshUC tokens.GenAccessTokenFromRefreshUC

	otpHandlers *OtpHandlers
	logger		*logger.Log
}

func NewLoginHandlers(
	loginUC login.LoginUC,
	genAuthTokenUC tokens.GenerateAuthTokensUC,
	genAccessTokenFromRefreshUC tokens.GenAccessTokenFromRefreshUC,

	otpHandler *OtpHandlers,
	logger *logger.Log,
) *LoginHandlers {
	return &LoginHandlers{
		loginUC:                     loginUC,
		genAuthTokenUC:              genAuthTokenUC,
		genAccessTokenFromRefreshUC: genAccessTokenFromRefreshUC,

		otpHandlers: otpHandler,
		logger: logger,
	}
}

type loginFormData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshFormData struct {
	Token string `json:"refresh_token"`
}

func (h *LoginHandlers) Login(ctx echo.Context) error {
	var data, validated loginFormData

	if err := ctx.Bind(&data); err != nil {
		return jsonInvalidRequestResponse(ctx)
	}

	if errList := validation.LoginValidator.Validate(data, &validated); errList != nil {
		return jsonValidationErrorResponse(ctx, errList)
	}

	if err := h.loginUC.Execute(
		ctx.Request().Context(),
		validated.Email,
		validated.Password,
	); err != nil {

		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			return jsonUnauthorizedResponse(ctx, err.Error())
		default:
			h.logger.Error.Println(err)
			return jsonInternalErrorResponse(ctx)
		}
	}

	return jsonSuccessMessageResponse(ctx, verificationSentMessage)
}

func (h *LoginHandlers) VerifyOtp(ctx echo.Context) error {
	handler := h.otpHandlers.Verify(domain.LoginOTP, func(validated verifyOtpFormData) error {
		access, refresh, err := h.genAuthTokenUC.Execute(ctx.Request().Context(), validated.Email)

		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				return jsonUnauthorizedResponse(ctx, domain.ErrInvalidOtp.Error())
			}

			h.logger.Error.Println(err)
			return jsonInternalErrorResponse(ctx)
		}

		return jsonSuccessWithDataResponse(ctx, map[string]string{
			"access_token":  access,
			"refresh_token": refresh,
		})
	})

	return handler(ctx)
}

func (h *LoginHandlers) ResendOtp(ctx echo.Context) error {
	handler := h.otpHandlers.Resend(domain.LoginOTP)
	return handler(ctx)
}

func (h *LoginHandlers) Refresh(ctx echo.Context) error {
	var data refreshFormData

	if err := ctx.Bind(&data); err != nil {
		return jsonInvalidRequestResponse(ctx)
	}

	if data.Token == "" {
		return jsonUnauthorizedResponse(ctx, domain.ErrInvalidToken.Error())
	}

	token, err := h.genAccessTokenFromRefreshUC.Execute(
		ctx.Request().Context(),
		data.Token,
	);

	if err != nil {
		
	}
}
