package service

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/shared"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secret         string
	accessTokenTTL time.Duration
	logger         *shared.Log
}

func NewJwtService(secret string, accessTokenTTL time.Duration, logger *shared.Log) *JwtService {
	return &JwtService{
		secret:         secret,
		accessTokenTTL: accessTokenTTL,
		logger:         logger,
	}
}

func (service *JwtService) GenerateToken(user *domain.AuthUser) (string, error) {
	expirationTime := time.Now().Add(service.accessTokenTTL)

	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(service.secret)

	if err != nil {
		service.logger.Error.Println(err)
		return "", domain.ErrInternal
	}

	return tokenString, nil
}

func (service *JwtService) ValidateToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}

		return service.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrExpiredToken
		}
		return nil, domain.ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, domain.ErrExpiredToken
}
