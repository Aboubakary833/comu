package handlers

import (
	"comu/internal/modules/auth/application"
	"comu/internal/shared/logger"

	"github.com/labstack/echo/v4"
)

type Handlers interface {
	RegisterRoutes(*echo.Echo)
}

func GetHandlers(ucs application.UseCases, logger *logger.Log) []Handlers {
	otpHandlers := newOtpHandlers(ucs.VerifyOtpUC, ucs.ResendOtpUC, logger)
	loginHandlers := newLoginHandlers(
		ucs.LoginUC, ucs.GenAuthTokenUC,
		ucs.GenAccessTokenFromRefresh, otpHandlers, logger,
	)
	registerHandlers := newRegisterHandlers(
		ucs.RegisterUC, ucs.GenAuthTokenUC,
		ucs.MarkUserAsVerifiedUC, otpHandlers, logger,
	)
	resetPasswordHandlers := newResetPasswordHandlers(
		ucs.NewPasswordUC, ucs.GenResetTokenUC,
		ucs.ResetPasswordUC, otpHandlers, logger,
	)

	return []Handlers{
		loginHandlers,
		registerHandlers,
		resetPasswordHandlers,
	}
}
