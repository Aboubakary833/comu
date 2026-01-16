package application

import (
	"comu/internal/modules/auth/application/login"
	"comu/internal/modules/auth/application/otp"
	"comu/internal/modules/auth/application/register"
	resetPassword "comu/internal/modules/auth/application/reset_password"
	"comu/internal/modules/auth/application/tokens"
	"comu/internal/modules/auth/domain"
)

type UseCases struct {
	LoginUC                   *login.LoginUC
	RegisterUC                *register.RegisterUC
	MarkUserAsVerifiedUC      *register.MarkUserAsVerifiedUC
	ResetPasswordUC           *resetPassword.ResetPasswordUC
	NewPasswordUC             *resetPassword.SetNewPasswordUC
	VerifyOtpUC               *otp.VerifyOtpUC
	ResendOtpUC               *otp.ResendOtpUC
	GenAuthTokenUC            *tokens.GenerateAuthTokensUC
	GenResetTokenUC           *tokens.GenerateResetTokenUC
	VerifyAccessToken		  *tokens.VerifyAccessTokenUC
	GenAccessTokenFromRefresh *tokens.GenAccessTokenFromRefreshUC
}

func InitUseCases(
	otpCodesRepo domain.OtpCodesRepository,
	resetTokensRepo domain.ResetTokensRepository,
	refreshTokensRepo domain.RefreshTokensRepository,
	resendRequestsRepo domain.ResendOtpRequestsRepository,

	jwtService domain.JwtService,
	userService domain.UserService,
	passwordService domain.PasswordService,
	notificationService domain.NotificationService,
) UseCases {
	loginUC := login.NewUseCase(
		userService,
		passwordService,
		otpCodesRepo,
		notificationService,
		resendRequestsRepo,
	)
	registerUC := register.NewRegisterUseCase(
		userService,
		passwordService,
		otpCodesRepo,
		notificationService,
		resendRequestsRepo,
	)

	markUserAsVerifiedUC := register.NewMarkUserAsVerifiedUseCase(userService)

	resetPasswordUC := resetPassword.NewResetPasswordUseCase(
		userService,
		otpCodesRepo,
		notificationService,
		resendRequestsRepo,
	)

	newPasswordUC := resetPassword.NewSetNewPasswordUseCase(
		userService,
		passwordService,
		notificationService,
		resetTokensRepo,
	)

	verifyOtpUC := otp.NewVerifyOtpUseCase(otpCodesRepo)
	resendOtpUC := otp.NewResendOtpUseCase(
		otpCodesRepo,
		notificationService,
		resendRequestsRepo,
	)

	genResetTokenUC := tokens.NewGenResetTokenUseCase(userService, resetTokensRepo)
	genAuthTokenUC := tokens.NewGenAuthTokensUseCase(jwtService, userService, refreshTokensRepo)
	verifyAccessTokenUC := tokens.NewVerifyAccessTokenUseCase(jwtService, userService)
	genAccessFromTokenRefreshUC := tokens.NewGenAccessTokenFromRefreshUseCase(jwtService, userService, refreshTokensRepo)

	return UseCases{
		LoginUC:                   loginUC,
		RegisterUC:                registerUC,
		MarkUserAsVerifiedUC:      markUserAsVerifiedUC,
		ResetPasswordUC:           resetPasswordUC,
		NewPasswordUC:             newPasswordUC,
		VerifyOtpUC:               verifyOtpUC,
		ResendOtpUC:               resendOtpUC,
		GenAuthTokenUC:            genAuthTokenUC,
		GenResetTokenUC:           genResetTokenUC,
		VerifyAccessToken: 		   verifyAccessTokenUC,
		GenAccessTokenFromRefresh: genAccessFromTokenRefreshUC,
	}
}
