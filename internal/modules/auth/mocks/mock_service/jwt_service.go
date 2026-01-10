package mockService

import (
	"comu/internal/modules/auth/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type jwtServiceMock struct {
	mock.Mock
}

func NewJwtServiceMock() *jwtServiceMock {
	return new(jwtServiceMock)
}

func (serviceMock *jwtServiceMock) GenerateToken(user *domain.AuthUser) (string, error) {
	args := serviceMock.Called(user)
	return args.String(0), args.Error(1)
}

func (serviceMock *jwtServiceMock) ValidateToken(tokenString string) (jwt.Claims, error) {
	args := serviceMock.Called(tokenString)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(jwt.Claims), nil
}
