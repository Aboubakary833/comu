package auth

import (
	"comu/internal/modules/auth/application/tokens"
	echoRes "comu/internal/shared/utils/echo_res"
	"github.com/labstack/echo/v4"
)

var (
	msgUnauthenticatedUser = "User is not authenticated"
	msgAuthenticatedUserFound = "Authenticated user found"
	msgUserIsNotVerified = "User email address is not verified"
)

var (
	notVerified		echoRes.ErrorResponseType = "unverified"
	authenticated	echoRes.ErrorResponseType = "authenticated"
	unauthenticated echoRes.ErrorResponseType = "unauthenticated"
)

type publicApi struct {
	verifyTokenUC *tokens.VerifyAccessTokenUC
}

func newApi(verifyTokenUC *tokens.VerifyAccessTokenUC) *publicApi {
	return &publicApi{
		verifyTokenUC: verifyTokenUC,
	}
}

func (api *publicApi) getAuthToken(ctx echo.Context) string {
	return ctx.Request().Header.Get("Authorization")
}

func (api *publicApi) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := api.getAuthToken(ctx)

		if token == "" {
			return echoRes.JsonUnauthorizedResponse(
				ctx, unauthenticated,
				msgUnauthenticatedUser,
			)
		}

		user, err := api.verifyTokenUC.Execute(ctx.Request().Context(), token)

		if err != nil {
			return echoRes.JsonUnauthorizedResponse(ctx, unauthenticated, err.Error())
		}
		user.Password = ""
		ctx.Set(AuthUserIdCtxKey, user.ID.String())
		ctx.Set(AuthIsUserVerifiedCtxKey, user.EmailVerifiedAt != nil)

		return next(ctx)
	}
}

func (api *publicApi) GuestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := api.getAuthToken(ctx)

		if token == "" {
			return next(ctx)
		}

		_, err := api.verifyTokenUC.Execute(ctx.Request().Context(), token)

		if err != nil {
			return next(ctx)
		}

		return echoRes.JsonUnauthorizedResponse(
			ctx, authenticated,
			msgAuthenticatedUserFound,
		)
	}
}

// This middleware should always come before the AuthMiddleware
func (api *publicApi) VerifiedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if verified, ok := ctx.Get(AuthIsUserVerifiedCtxKey).(bool); !ok || !verified {
			return echoRes.JsonUnauthorizedResponse(ctx, notVerified, msgUserIsNotVerified)
		}

		return next(ctx)
	}
}
