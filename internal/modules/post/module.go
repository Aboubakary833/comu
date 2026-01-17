package post

import (
	"comu/internal/modules/auth"
	"comu/internal/modules/post/application"
	"comu/internal/modules/post/infra/mysql"
	"comu/internal/modules/post/presentation/handlers"
	"comu/internal/shared/logger"
	"database/sql"

	"github.com/labstack/echo/v4"
)

type postModule struct {
	authApi  auth.PublicApi
	handlers []handlers.Handlers
}

func NewModule(db *sql.DB, authApi auth.PublicApi, logger *logger.Log) *postModule {
	postsRepo := mysql.NewPostRepository(db)
	commentsRepo := mysql.NewCommentsRepository(db)

	useCases := application.InitUseCases(postsRepo, commentsRepo)
	handlers := handlers.GetHandlers(useCases, logger)

	return &postModule{
		authApi:  authApi,
		handlers: handlers,
	}
}

func (module *postModule) RegisterRoutes(echo *echo.Echo) {
	for _, h := range module.handlers {
		h.RegisterRoutes(
			echo,
			module.authApi.AuthMiddleware,
			module.authApi.VerifiedMiddleware,
		)
	}
}
