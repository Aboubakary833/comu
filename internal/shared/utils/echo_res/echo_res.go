package echoRes

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/shared/utils"
	"comu/internal/shared/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponseType string

var (
	badRequest    ErrorResponseType = "bad_request"
	internalError ErrorResponseType = "internal_error"
)

type ErrorResponse struct {
	Code    int                              `json:"code"`
	Type    ErrorResponseType                `json:"type,omitempty"`
	Status  string                           `json:"status"`
	Message string                           `json:"message"`
	Errors  validator.SchemaValidationErrors `json:"errors,omitempty"`
}

type SuccessResponse[T any] struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func JsonInvalidRequestResponse(ctx echo.Context) error {
	return JsonErrorMessageResponse(ctx, http.StatusBadRequest, badRequest, "Invalid request body")
}

func JsonInternalErrorResponse(ctx echo.Context) error {
	return JsonErrorMessageResponse(ctx, http.StatusInternalServerError, "internal_error", domain.ErrInternal.Error())
}

func JsonUnauthorizedResponse(ctx echo.Context, errType ErrorResponseType, message string) error {
	return JsonErrorMessageResponse(ctx, http.StatusUnauthorized, errType, message)
}

func JsonErrorMessageResponse(ctx echo.Context, statusCode int, errType ErrorResponseType, message string) error {
	return ctx.JSON(
		statusCode,
		ErrorResponse{
			Code:    statusCode,
			Type:    errType,
			Status:  http.StatusText(statusCode),
			Message: utils.UcFirst(message),
		},
	)
}

func JsonValidationErrorResponse(ctx echo.Context, errList validator.SchemaValidationErrors) error {
	return ctx.JSON(
		http.StatusUnprocessableEntity,
		ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Type:    "validation_errors",
			Status:  http.StatusText(http.StatusUnprocessableEntity),
			Message: "Validation failed",
			Errors:  errList,
		},
	)
}

func JsonSuccessResponse[T any](ctx echo.Context, message string, data T) error {
	return ctx.JSON(
		http.StatusOK,
		SuccessResponse[T]{
			Code:    http.StatusOK,
			Status:  "success",
			Message: message,
			Data:    data,
		},
	)
}

func JsonSuccessMessageResponse(ctx echo.Context, message string) error {
	return JsonSuccessResponse(ctx, message, "")
}

func JsonSuccessWithDataResponse[T any](ctx echo.Context, data T) error {
	return JsonSuccessResponse(ctx, "", data)
}
