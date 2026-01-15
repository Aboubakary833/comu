package handlers

import (
	"comu/internal/modules/auth/application/register"
	"comu/internal/modules/auth/application/tokens"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/presentation/validation"
	"comu/internal/shared/logger"
	"errors"

	"github.com/labstack/echo/v4"
)

type RegisterHandlers struct {
	registerUC           register.RegisterUC
	genAuthTokenUC       tokens.GenerateAuthTokensUC
	markUserAsVerifiedUC register.MarkUserAsVerifiedUC

	otpHandlers *OtpHandlers
	logger      *logger.Log
}

func NewRegisterHandlers(
	registerUC register.RegisterUC,
	genAuthTokenUC tokens.GenerateAuthTokensUC,
	markUserAsVerifiedUC register.MarkUserAsVerifiedUC,

	otpHandler *OtpHandlers,
	logger *logger.Log,
) *RegisterHandlers {
	return &RegisterHandlers{
		registerUC:           registerUC,
		genAuthTokenUC:       genAuthTokenUC,
		markUserAsVerifiedUC: markUserAsVerifiedUC,

		otpHandlers: otpHandler,
		logger:      logger,
	}
}

type registerFormData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *RegisterHandlers) Register(ctx echo.Context) error {
	var data, validated registerFormData

	if err := ctx.Bind(&data); err != nil {
		return jsonInvalidRequestResponse(ctx)
	}
	errList := validation.RegisterValidator.Validate(data, &validated)

	if errList != nil {
		return jsonValidationErrorResponse(ctx, errList)
	}

	if err := h.registerUC.Execute(
		ctx.Request().Context(),
		validated.Name, validated.Email, validated.Password,
	); err != nil {

		if errors.Is(err, domain.ErrUserEmailTaken) {
			return jsonUnauthorizedResponse(ctx, err.Error())
		}

		h.logger.Error.Println(err)
		return jsonInternalErrorResponse(ctx)
	}

	return jsonSuccessMessageResponse(ctx, verificationSentMessage)
}

func (h *RegisterHandlers) VerifyOtp(ctx echo.Context) error {
	handler := h.otpHandlers.Verify(domain.LoginOTP, func(validated verifyOtpFormData) error {

		if err := h.markUserAsVerifiedUC.Execute(ctx.Request().Context(), validated.Email); err != nil {
			if errors.Is(err, domain.ErrUserEmailTaken) {
				return jsonUnauthorizedResponse(ctx, domain.ErrInvalidOtp.Error())
			}

			h.logger.Error.Println(err)
			return jsonInternalErrorResponse(ctx)
		}

		access, refresh, err := h.genAuthTokenUC.Execute(ctx.Request().Context(), validated.Email)

		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				return jsonUnauthorizedResponse(ctx, domain.ErrInternal.Error())
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

func (h *RegisterHandlers) ResendOtp(ctx echo.Context) error {
	handler := h.otpHandlers.Resend(domain.RegisterOTP)
	return handler(ctx)
}
