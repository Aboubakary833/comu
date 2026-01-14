package handlers

import (
	"comu/internal/modules/auth/application/login"
	"comu/internal/modules/auth/application/otp"
	"comu/internal/modules/auth/application/tokens"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/presentation/validation"
	"errors"

	"github.com/labstack/echo/v4"
)

type LoginHandlers struct {
	loginUC                     login.LoginUC
	verifyOtpUC                 otp.VerifyOtpUC
	resendOtpUC                 otp.ResendOtpUC
	genAuthTokenUC              tokens.GenerateAuthTokensUC
	genAccessTokenFromRefreshUC tokens.GenAccessTokenFromRefreshUC
}

func NewLoginHandlers(
	loginUC login.LoginUC,
	verifyOtpUC otp.VerifyOtpUC,
	resendOtpUC otp.ResendOtpUC,
	genAuthTokenUC tokens.GenerateAuthTokensUC,
	genAccessTokenFromRefreshUC tokens.GenAccessTokenFromRefreshUC,
) *LoginHandlers {
	return &LoginHandlers{
		loginUC:                     loginUC,
		verifyOtpUC:                 verifyOtpUC,
		genAuthTokenUC:              genAuthTokenUC,
		genAccessTokenFromRefreshUC: genAccessTokenFromRefreshUC,
	}
}

type loginFormData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type otpFormData struct {
	Email string `json:"email"`
	Code  string `json:"code"`
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
			return jsonInternalErrorResponse(ctx)
		}
	}

	return jsonSuccessMessageResponse(ctx, "A verification code has been sent to your mail.")
}

func (h *LoginHandlers) VerifyOtp(ctx echo.Context) error {
	var data, validated otpFormData

	if err := ctx.Bind(&data); err != nil {
		return jsonInvalidRequestResponse(ctx)
	}

	if errList := validation.OtpCodeValidator.Validate(data, &validated); errList != nil {
		return jsonValidationErrorResponse(ctx, errList)
	}

	if err := h.verifyOtpUC.Execute(
		ctx.Request().Context(), otp.VerifyOtpInput{
			UserEmail:    validated.Email,
			OtpCodeType:  domain.LoginOTP,
			OtpCodeValue: validated.Code,
		},
	); err != nil {

		switch {
		case errors.Is(err, domain.ErrInvalidOtp) || errors.Is(err, domain.ErrExpiredOtp):
			return jsonUnauthorizedResponse(ctx, err.Error())
		default:
			return jsonInternalErrorResponse(ctx)
		}
	}

	access, refresh, err := h.genAuthTokenUC.Execute(ctx.Request().Context(), validated.Email)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return jsonUnauthorizedResponse(ctx, domain.ErrInternal.Error())
		}

		return jsonInternalErrorResponse(ctx)
	}

	return jsonSuccessWithDataResponse(ctx, map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
