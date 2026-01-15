package auth

import (
	"comu/config"
	"comu/internal/modules/auth/application/login"
	"comu/internal/modules/auth/application/otp"
	"comu/internal/modules/auth/application/register"
	resetPassword "comu/internal/modules/auth/application/reset_password"
	"comu/internal/modules/auth/application/tokens"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/infra/mysql"
	"comu/internal/modules/auth/infra/service"
	"comu/internal/modules/auth/presentation/handlers"
	"comu/internal/modules/users"
	"comu/internal/shared/logger"
	"database/sql"

	"github.com/labstack/echo/v4"
)

type PublicApi interface {
	AuthMiddleware() echo.HandlerFunc
	GuestMiddleware() echo.HandlerFunc
}

type useCases struct {
	loginUC              *login.LoginUC
	registerUC           *register.RegisterUC
	markUserAsVerifiedUC *register.MarkUserAsVerifiedUC

	resetPasswordUC *resetPassword.ResetPasswordUC
	newPasswordUC   *resetPassword.SetNewPasswordUC

	verifyOtpUC *otp.VerifyOtpUC
	resendOtpUC *otp.ResendOtpUC

	genAuthTokenUC            *tokens.GenerateAuthTokensUC
	genResetTokenUC           *tokens.GenerateResetTokenUC
	genAccessTokenFromRefresh *tokens.GenAccessTokenFromRefreshUC
}

type authHandler interface {
	RegisterRoutes(*echo.Echo)
}

type authModule struct {
	api PublicApi
	handlers []authHandler
}

func NewModule(
	db *sql.DB, config *config.Config,
	usersApi users.PublicApi, logger *logger.Log,
) *authModule {
	otpCodesRepo := mysql.NewOtpCodesRepository(db)
	resetTokensRepo := mysql.NewResetTokensRepository(db)
	refreshTokensRepo := mysql.NewRefreshTokensRepository(db)
	resendRequestsRepo := mysql.NewResendOtpRequestsRepository(db)

	jwtService := service.NewJwtService(config.AppKey, domain.DefaultAccessTokenTTL, logger)
	userService := service.NewUserService(usersApi, logger)
	passwordService := service.NewPasswordService(logger)
	notificationService, err := service.NewSmtpNotificationService(
		config.MailHost, config.MailPort, config.MailFrom,
		service.SmtpNotificationAuth{
			Username: config.MailUserName,
			Password: config.MailPassword,
		},
	)

	if err != nil {
		logger.Error.Fatalln(err)
	}

	ucs := newUseCases(
		otpCodesRepo,
		resetTokensRepo,
		refreshTokensRepo,
		resendRequestsRepo,
		jwtService,
		userService,
		passwordService,
		notificationService,
	)
	handlers := getAuthHandlers(ucs, logger)

	return &authModule{
		handlers: handlers,
	}
}

func (module *authModule) RegisterRoutes(echo *echo.Echo) {
	for _, h := range module.handlers {
		h.RegisterRoutes(echo)
	}
}

func (module *authModule) GetPublicApi() PublicApi {
	return module.api
}

func getAuthHandlers(ucs useCases, logger *logger.Log) []authHandler {
	otpHandlers := handlers.NewOtpHandlers(ucs.verifyOtpUC, ucs.resendOtpUC, logger)
	loginHandlers := handlers.NewLoginHandlers(
		ucs.loginUC, ucs.genAuthTokenUC,
		ucs.genAccessTokenFromRefresh, otpHandlers, logger,
	)
	registerHandlers := handlers.NewRegisterHandlers(
		ucs.registerUC, ucs.genAuthTokenUC,
		ucs.markUserAsVerifiedUC, otpHandlers, logger,
	)
	resetPasswordHandlers := handlers.NewResetPasswordHandlers(
		ucs.newPasswordUC, ucs.genResetTokenUC,
		ucs.resetPasswordUC, otpHandlers, logger,
	)

	return []authHandler{
		loginHandlers,
		registerHandlers,
		resetPasswordHandlers,
	}
}

func newUseCases(
	otpCodesRepo domain.OtpCodesRepository,
	resetTokensRepo domain.ResetTokensRepository,
	refreshTokensRepo domain.RefreshTokensRepository,
	resendRequestsRepo domain.ResendOtpRequestsRepository,

	jwtService domain.JwtService,
	userService domain.UserService,
	passwordService domain.PasswordService,
	notificationService domain.NotificationService,
) useCases {
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
	genAccessFromTokenRefreshUC := tokens.NewGenAccessTokenFromRefreshUseCase(jwtService, userService, refreshTokensRepo)

	return useCases{
		loginUC:                   loginUC,
		registerUC:                registerUC,
		markUserAsVerifiedUC:      markUserAsVerifiedUC,
		resetPasswordUC:           resetPasswordUC,
		newPasswordUC:             newPasswordUC,
		verifyOtpUC:               verifyOtpUC,
		resendOtpUC:               resendOtpUC,
		genAuthTokenUC:            genAuthTokenUC,
		genResetTokenUC:           genResetTokenUC,
		genAccessTokenFromRefresh: genAccessFromTokenRefreshUC,
	}
}
