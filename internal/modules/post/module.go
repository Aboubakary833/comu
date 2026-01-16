package post

import "github.com/labstack/echo/v4"

type Handlers interface {
	RegisterRoutes(*echo.Echo, ...echo.MiddlewareFunc)
}
