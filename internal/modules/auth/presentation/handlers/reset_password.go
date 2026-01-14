package handlers

import (
	resetPassword "comu/internal/modules/auth/application/reset_password"
	"comu/internal/modules/auth/application/tokens"
)

type ResetPasswordHandlers struct {
	newPasswordUC resetPassword.SetNewPasswordUC
	genResetTokenUC tokens.GenerateResetTokenUC
	resetPasswordUC resetPassword.ResetPasswordUC
}

