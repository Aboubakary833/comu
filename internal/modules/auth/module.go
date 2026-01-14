package auth

import "github.com/labstack/echo/v4"

type PublicApi interface {
	AuthMiddleware() echo.HandlerFunc
	GuestMiddleware() echo.HandlerFunc
}

type AuthModule struct {
	api PublicApi
}

func (module *AuthModule) GetPublicApi() PublicApi {
	return module.api
}
