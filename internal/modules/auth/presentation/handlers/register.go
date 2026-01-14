package handlers

import (
	"comu/internal/modules/auth/application/otp"
	"comu/internal/modules/auth/application/register"
	"comu/internal/modules/auth/application/tokens"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/presentation/validation"
	"errors"

	"github.com/labstack/echo/v4"
)

type RegisterHandlers struct {
	verifyOtp            otp.VerifyOtpUC
	resendOtpUC          otp.ResendOtpUC
	registerUC           register.RegisterUC
	genAuthTokenUC       tokens.GenerateAuthTokensUC
	markUserAsVerifiedUC register.MarkUserAsVerifiedUC
}

func NewRegisterHandlers(
	verifyOtp otp.VerifyOtpUC,
	resendOtpUC otp.ResendOtpUC,
	registerUC register.RegisterUC,
	genAuthTokenUC tokens.GenerateAuthTokensUC,
	markUserAsVerifiedUC register.MarkUserAsVerifiedUC,
) *RegisterHandlers {
	return &RegisterHandlers{
		verifyOtp:            verifyOtp,
		resendOtpUC:          resendOtpUC,
		registerUC:           registerUC,
		genAuthTokenUC:       genAuthTokenUC,
		markUserAsVerifiedUC: markUserAsVerifiedUC,
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

		return jsonInternalErrorResponse(ctx)
	}

	return jsonSuccessMessageResponse(ctx, "A verification code has been sent to your mail.")
}

