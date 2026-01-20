package handlers

import (
	"comu/internal/modules/auth"
	"comu/internal/modules/post/application/posts"
	"comu/internal/modules/post/domain"
	"comu/internal/modules/post/presentation/validation"
	"comu/internal/shared/logger"
	echoRes "comu/internal/shared/utils/echo_res"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	unauthorized echoRes.ErrorResponseType = "unauthorized"
)

type postHandlers struct {
	listPostsUC  *posts.ListPostsUC
	readPostUC   *posts.ReadPostUC
	createPostUC *posts.CreatePostUC
	updatePostUC *posts.UpdatePostUC
	deletePostUC *posts.DeletePostUC

	logger *logger.Log
}

func newPostHandlers(
	listPostsUC *posts.ListPostsUC,
	readPostUC *posts.ReadPostUC,
	createPostUC *posts.CreatePostUC,
	updatePostUC *posts.UpdatePostUC,
	deletePostUC *posts.DeletePostUC,

	logger *logger.Log,
) *postHandlers {
	return &postHandlers{
		listPostsUC:  listPostsUC,
		readPostUC:   readPostUC,
		createPostUC: createPostUC,
		updatePostUC: updatePostUC,
		deletePostUC: deletePostUC,

		logger: logger,
	}
}

func (h *postHandlers) RegisterRoutes(echo *echo.Echo, m ...echo.MiddlewareFunc) {
	group := echo.Group("/posts", m...)

	group.POST("/", h.create)
	group.GET("/", h.list)
	group.GET("/:slug", h.read)
	group.PUT("/:post_id", h.update)
	group.DELETE("/:post_id", h.delete)
}

type postFormData struct {
	Title   string `form:"title" json:"title"`
	Content string `form:"content" json:"content"`
}

func (h *postHandlers) list(ctx echo.Context) error {
	paginator := getPaginatorFromCtx(ctx)

	posts, next, err := h.listPostsUC.Execute(
		ctx.Request().Context(),
		paginator,
	)

	if err != nil {
		h.logger.Error.Println(err)
		return echoRes.JsonInternalErrorResponse(ctx)
	}
	cursor, err := next.ToBase64()

	if err != nil {
		h.logger.Error.Println(err)
		return echoRes.JsonInternalErrorResponse(ctx)
	}

	return echoRes.JsonSuccessWithDataResponse(ctx, map[string]any{
		"posts":  posts,
		"cursor": cursor,
	})
}

func (h *postHandlers) read(ctx echo.Context) error {
	slug := ctx.Param("slug")

	if slug == "" {
		return echoRes.JsonNotFoundResponse(
			ctx, domain.ErrPostNotFound.Error(),
		)
	}

	post, err := h.readPostUC.Execute(ctx.Request().Context(), slug)

	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return echoRes.JsonNotFoundResponse(ctx, err.Error())
		}

		h.logger.Error.Println(err)
		return echoRes.JsonInternalErrorResponse(ctx)
	}

	return echoRes.JsonSuccessWithDataResponse(ctx, map[string]any{
		"post": *post,
	})
}

func (h *postHandlers) create(ctx echo.Context) error {
	handler := postPreHandler(func(validated postFormData, userID uuid.UUID) error {
		post, err := h.createPostUC.Execute(
			ctx.Request().Context(),
			posts.CreatePostInput{
				UserID:  userID,
				Title:   validated.Title,
				Content: validated.Content,
			},
		)

		if err != nil {
			h.logger.Error.Println(err)
			return echoRes.JsonInternalErrorResponse(ctx)
		}

		return echoRes.JsonSuccessWithDataResponse(ctx, map[string]any{
			"post": *post,
		})
	})

	return handler(ctx)
}

func (h *postHandlers) update(ctx echo.Context) error {
	handler := postPreHandler(func(validated postFormData, userID uuid.UUID) error {
		postID, err := uuid.Parse(ctx.Param("post_id"))

		if err != nil {
			return echoRes.JsonNotFoundResponse(
				ctx, domain.ErrPostNotFound.Error(),
			)
		}

		slug, err := h.updatePostUC.Execute(
			ctx.Request().Context(),
			posts.UpdatePostInput{
				PostID:   postID,
				AuthorID: userID,
				Title:    validated.Title,
				Content:  validated.Content,
			},
		)

		if err != nil {
			switch {
			case errors.Is(err, domain.ErrPostNotFound):
				return echoRes.JsonNotFoundResponse(ctx, err.Error())

			case errors.Is(err, domain.ErrUnauthorized):
				return echoRes.JsonForbiddenResponse(ctx, err.Error())

			default:
				h.logger.Error.Println(err)
				return echoRes.JsonInternalErrorResponse(ctx)
			}
		}

		return echoRes.JsonSuccessWithDataResponse(ctx, map[string]string{
			"slug": slug,
		})
	})

	return handler(ctx)
}

func (h *postHandlers) delete(ctx echo.Context) error {
	id := ctx.Get(auth.AuthUserIdCtxKey).(string)
	userID, err := uuid.Parse(id)

	if err != nil {
		return echoRes.JsonUnauthorizedResponse(
			ctx, unauthorized,
			domain.ErrUnauthorized.Error(),
		)
	}

	postID, err := uuid.Parse(ctx.Param("post_id"))

	if err != nil {
		return echoRes.JsonNotFoundResponse(
			ctx, domain.ErrPostNotFound.Error(),
		)
	}

	if err := h.deletePostUC.Execute(
		ctx.Request().Context(),
		postID, userID,
	); err != nil {
		switch {
		case errors.Is(err, domain.ErrPostNotFound):
			return echoRes.JsonNotFoundResponse(ctx, err.Error())

		case errors.Is(err, domain.ErrUnauthorized):
			return echoRes.JsonForbiddenResponse(ctx, err.Error())

		default:
			h.logger.Error.Println(err)
			return echoRes.JsonInternalErrorResponse(ctx)
		}
	}

	return echoRes.JsonSuccessMessageResponse(ctx, "Your post has successfully been deleted.")
}

func postPreHandler(afterFunc func(validated postFormData, userID uuid.UUID) error) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var data postFormData

		if err := ctx.Bind(data); err != nil {
			return echoRes.JsonInvalidRequestResponse(ctx)
		}

		if errList := validation.PostValidator.Validate(&data); errList != nil {
			return echoRes.JsonValidationErrorResponse(ctx, errList)
		}
		id := ctx.Get(auth.AuthUserIdCtxKey).(string)
		userID, err := uuid.Parse(id)

		if err != nil {
			return echoRes.JsonUnauthorizedResponse(
				ctx, unauthorized,
				domain.ErrUnauthorized.Error(),
			)
		}

		return afterFunc(data, userID)
	}
}

func getPaginatorFromCtx(ctx echo.Context) domain.Paginator {

	paginator := domain.Paginator{
		Limit: domain.DefaultPaginatorLimit,
	}
	limitParam := ctx.QueryParam("limit")
	cursorParam := ctx.QueryParam("cursor")

	if value, err := strconv.Atoi(limitParam); err == nil {
		paginator.Limit = value
	}

	if cursorParam == "" {
		return paginator
	}

	rawBytes, err := base64.RawStdEncoding.DecodeString(cursorParam)

	if err != nil {
		return paginator
	}

	cursor := domain.Cursor{}

	if err = json.Unmarshal(rawBytes, &cursor); err != nil {
		return paginator
	}
	paginator.After = &cursor

	return paginator
}
