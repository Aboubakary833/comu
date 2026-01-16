package auth

import (
	"comu/config"
	"comu/internal/modules/auth/application"
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/infra/mysql"
	"comu/internal/modules/auth/infra/service"
	"comu/internal/modules/auth/presentation/handlers"
	"comu/internal/modules/users"
	"comu/internal/shared/logger"
	"database/sql"

	"github.com/labstack/echo/v4"
)

var (
	AuthUserIdCtxKey = "userID"
	AuthIsUserVerifiedCtxKey = "isUserVerified"
)

type PublicApi interface {
	AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	GuestMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	VerifiedMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type authModule struct {
	api      PublicApi
	handlers []handlers.Handlers
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

	useCases := application.InitUseCases(
		otpCodesRepo,
		resetTokensRepo,
		refreshTokensRepo,
		resendRequestsRepo,
		jwtService,
		userService,
		passwordService,
		notificationService,
	)

	api := newApi(useCases.VerifyAccessToken)
	handlers := handlers.GetHandlers(useCases, logger)

	return &authModule{
		api:      api,
		handlers: handlers,
	}
}

func (module *authModule) RegisterRoutes(echo *echo.Echo) {
	for _, h := range module.handlers {
		h.RegisterRoutes(echo, module.api.GuestMiddleware)
	}
}

func (module *authModule) GetPublicApi() PublicApi {
	return module.api
}
