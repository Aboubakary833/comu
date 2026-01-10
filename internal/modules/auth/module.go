package auth

import (
	"net/http"
)

type PublicApi interface {
	AuthMiddleware() http.Handler
}

type AuthModule struct {
	api PublicApi
}

func (module *AuthModule) GetPublicApi() PublicApi {
	return module.api
}
