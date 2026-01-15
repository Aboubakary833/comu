package handlers

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/shared/utils"
	"comu/internal/shared/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type errorType string

var (
	badRequest    errorType = "bad_request"
	internalError errorType = "internal_error"
)

type errorResponse struct {
	Code    int                              `json:"code"`
	Type    errorType                        `json:"type,omitempty"`
	Status  string                           `json:"status"`
	Message string                           `json:"message"`
	Errors  validator.SchemaValidationErrors `json:"errors,omitempty"`
}

type successResponse[T any] struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func jsonInvalidRequestResponse(ctx echo.Context) error {
	return jsonErrorMessageResponse(ctx, http.StatusBadRequest, badRequest, "Invalid request body")
}

func jsonInternalErrorResponse(ctx echo.Context) error {
	return jsonErrorMessageResponse(ctx, http.StatusInternalServerError, "internal_error", domain.ErrInternal.Error())
}

func jsonUnauthorizedResponse(ctx echo.Context, errType errorType, message string) error {
	return jsonErrorMessageResponse(ctx, http.StatusUnauthorized, errType, message)
}

func jsonErrorMessageResponse(ctx echo.Context, statusCode int, errType errorType, message string) error {
	return ctx.JSON(
		statusCode,
		errorResponse{
			Code:    statusCode,
			Type:    errType,
			Status:  http.StatusText(statusCode),
			Message: utils.UcFirst(message),
		},
	)
}

func jsonValidationErrorResponse(ctx echo.Context, errList validator.SchemaValidationErrors) error {
	return ctx.JSON(
		http.StatusUnprocessableEntity,
		errorResponse{
			Code:    http.StatusUnprocessableEntity,
			Type:    "validation_errors",
			Status:  http.StatusText(http.StatusUnprocessableEntity),
			Message: "Validation failed",
			Errors:  errList,
		},
	)
}

func jsonSuccessResponse[T any](ctx echo.Context, message string, data T) error {
	return ctx.JSON(
		http.StatusOK,
		successResponse[T]{
			Code:    http.StatusOK,
			Status:  "success",
			Message: message,
			Data:    data,
		},
	)
}

func jsonSuccessMessageResponse(ctx echo.Context, message string) error {
	return jsonSuccessResponse(ctx, message, "")
}

func jsonSuccessWithDataResponse[T any](ctx echo.Context, data T) error {
	return jsonSuccessResponse(ctx, "", data)
}
