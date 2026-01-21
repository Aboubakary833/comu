package handlers

import (
	"comu/internal/modules/auth"
	"comu/internal/modules/post/application/comments"
	"comu/internal/modules/post/domain"
	"comu/internal/modules/post/presentation/validation"
	"comu/internal/shared/logger"
	echoRes "comu/internal/shared/utils/echo_res"
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type commentHandlers struct {
	listCommentsUC  *comments.ListCommentsUC
	createCommentUC *comments.CreateCommentUC
	updateCommentUC *comments.UpdateCommentUC
	deleteCommentUC *comments.DeleteCommentUC

	logger *logger.Log
}

func newCommentsHandler(
	listCommentsUC *comments.ListCommentsUC,
	createCommentUC *comments.CreateCommentUC,
	updateCommentUC *comments.UpdateCommentUC,
	deleteCommentUC *comments.DeleteCommentUC,

	logger *logger.Log,
) *commentHandlers {
	return &commentHandlers{
		listCommentsUC:  listCommentsUC,
		createCommentUC: createCommentUC,
		updateCommentUC: updateCommentUC,
		deleteCommentUC: deleteCommentUC,

		logger: logger,
	}
}

func (h *commentHandlers) RegisterRoutes(echo *echo.Echo, m ...echo.MiddlewareFunc) {
	group := echo.Group("/comments", m...)

	group.GET("/list/:post_id", h.list)
	group.POST("/create", h.create)
	group.PUT("/update/:comment_id", h.list)
	group.DELETE("/delete/:comment_id", h.list)
}

type defaultCommentFormData struct {
	PostID  string `form:"post_id" json:"post_id"`
	Content string `form:"content" json:"content"`
}

func (h *commentHandlers) list(ctx echo.Context) error {
	paginator := getPaginatorFromCtx(ctx)

	postID, err := uuid.Parse(ctx.Param("post_id"))

	if err != nil {
		return echoRes.JsonUnauthorizedResponse(
			ctx, unauthorized,
			domain.ErrUnauthorized.Error(),
		)
	}

	comments, next, err := h.listCommentsUC.Execute(
		ctx.Request().Context(),
		postID,
		paginator,
	)

	if next == nil {
		return echoRes.JsonSuccessWithDataResponse(ctx, map[string]any{
			"comments": comments,
			"cursor":   "",
		})
	}

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
		"comments": comments,
		"cursor":   cursor,
	})
}

func (h *commentHandlers) create(ctx echo.Context) error {
	handler := commentPreHandler(func(authorID uuid.UUID, validated defaultCommentFormData) error {

		postID, err := uuid.Parse(validated.PostID)

		if err != nil {
			return echoRes.JsonUnauthorizedResponse(
				ctx, unauthorized,
				domain.ErrUnauthorized.Error(),
			)
		}

		comment, err := h.createCommentUC.Execute(
			ctx.Request().Context(),
			comments.CreateCommentInput{
				PostID:   postID,
				AuthorID: authorID,
				Content:  validated.Content,
			},
		)

		if err != nil {
			h.logger.Error.Println(err)
			return echoRes.JsonInternalErrorResponse(ctx)
		}

		return echoRes.JsonSuccessWithDataResponse(ctx, map[string]any{
			"comment": *comment,
		})
	})

	return handler(ctx)
}

func (h *commentHandlers) update(ctx echo.Context) error {
	handler := commentPreHandler(func(authorID uuid.UUID, validated defaultCommentFormData) error {
		commentID, err := uuid.Parse(ctx.Param("comment_id"))

		if err != nil {
			return echoRes.JsonNotFoundResponse(
				ctx, domain.ErrCommentNotFound.Error(),
			)
		}

		if err := h.updateCommentUC.Execute(
			ctx.Request().Context(),
			commentID, authorID, validated.Content,
		); err != nil {
			switch {
			case errors.Is(err, domain.ErrCommentNotFound):
				return echoRes.JsonNotFoundResponse(ctx, err.Error())

			case errors.Is(err, domain.ErrUnauthorized):
				return echoRes.JsonForbiddenResponse(ctx, err.Error())

			default:
				h.logger.Error.Println(err)
				return echoRes.JsonInternalErrorResponse(ctx)
			}
		}

		return echoRes.JsonSuccessMessageResponse(ctx, "Your comment has successfully been updated.")
	})

	return handler(ctx)
}

func (h *commentHandlers) delete(ctx echo.Context) error {
	id := ctx.Get(auth.AuthUserIdCtxKey).(string)
	userID, err := uuid.Parse(id)

	if err != nil {
		return echoRes.JsonUnauthorizedResponse(
			ctx, unauthorized,
			domain.ErrUnauthorized.Error(),
		)
	}

	commentID, err := uuid.Parse(ctx.Param("comment_id"))

	if err != nil {
		return echoRes.JsonNotFoundResponse(
			ctx, domain.ErrCommentNotFound.Error(),
		)
	}

	if err := h.deleteCommentUC.Execute(
		ctx.Request().Context(),
		commentID, userID,
	); err != nil {
		switch {
		case errors.Is(err, domain.ErrCommentNotFound):
			return echoRes.JsonNotFoundResponse(ctx, err.Error())

		case errors.Is(err, domain.ErrUnauthorized):
			return echoRes.JsonForbiddenResponse(ctx, err.Error())

		default:
			h.logger.Error.Println(err)
			return echoRes.JsonInternalErrorResponse(ctx)
		}
	}

	return echoRes.JsonSuccessMessageResponse(ctx, "Your comment has successfully been deleted.")
}

func commentPreHandler(afterFunc func(authorID uuid.UUID, validated defaultCommentFormData) error) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var data defaultCommentFormData

		if err := ctx.Bind(data); err != nil {
			return echoRes.JsonInvalidRequestResponse(ctx)
		}

		if errList := validation.CommentValidator.Validate(&data); errList != nil {
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

		return afterFunc(userID, data)
	}
}
