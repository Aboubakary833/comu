package auth

import (
	"context"
	"errors"

	"comu/internal/modules/auth/domain"
	echoRes "comu/internal/shared/utils/echo_res"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrAuthenticatedUserFound = errors.New("authenticated user found")
	ErrUnauthenticatedUser    = errors.New("user is not authenticated")
	ErrAuthUserIDMissing      = errors.New("auth user id missing from context")
	ErrAuthUserIDInvalid      = errors.New("auth user id is invalid")
)

type publicApi struct {
	userService domain.UserService
}

func newApi(userService domain.UserService) *publicApi {
	return &publicApi{
		userService: userService,
	}
}

func authUserID(ctx echo.Context) (uuid.UUID, error) {
	raw := ctx.Get(AuthUserIdCtxKey)
	if raw == nil {
		return uuid.Nil, ErrAuthUserIDMissing
	}

	id, ok := raw.(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.Nil, ErrAuthUserIDInvalid
	}

	return id, nil
}

func (api *publicApi) userExists(ctx context.Context, id uuid.UUID) bool {
	_, err := api.userService.GetUserByID(ctx, id)
	return err == nil
}

func (api *publicApi) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, err := authUserID(ctx)
		if err != nil {
			return echoRes.JsonUnauthorizedResponse(
				ctx,
				"unauthenticated",
				ErrUnauthenticatedUser.Error(),
			)
		}

		if !api.userExists(ctx.Request().Context(), userID) {
			return echoRes.JsonUnauthorizedResponse(
				ctx,
				"unauthenticated",
				ErrUnauthenticatedUser.Error(),
			)
		}

		return next(ctx)
	}
}

func (api *publicApi) GuestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, err := authUserID(ctx)
		if err != nil {
			return next(ctx)
		}

		if api.userExists(ctx.Request().Context(), userID) {
			return echoRes.JsonUnauthorizedResponse(
				ctx,
				"authenticated",
				ErrAuthenticatedUserFound.Error(),
			)
		}

		return next(ctx)
	}
}
