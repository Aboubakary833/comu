package handlers

import (
	"comu/internal/modules/post/application"
	"comu/internal/shared/logger"

	"github.com/labstack/echo/v4"
)

type Handlers interface {
	RegisterRoutes(*echo.Echo, ...echo.MiddlewareFunc)
}

func GetHandlers(ucs application.UseCases, logger *logger.Log) []Handlers {
	postsHandlers := newPostHandlers(
		ucs.ListPostsUC, ucs.ReadPostUC, ucs.CreatePostUC,
		ucs.UpdatePostUC, ucs.DeletePostUC, logger,
	)
	commentHandlers := newCommentsHandler(
		ucs.ListCommentUC, ucs.CreateCommentUC,
		ucs.UpdateCommentUC, ucs.DeleteCommentUC, logger,
	)

	return []Handlers{postsHandlers, commentHandlers}
}
