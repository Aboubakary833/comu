package tokens

import (
	"comu/internal/modules/auth/domain"
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


type VerifyAccessTokenUC struct {
	jwtService domain.JwtService
	userService domain.UserService
}


func NewVerifyAccessTokenUseCase(
	jwtService domain.JwtService,
	userService domain.UserService,
) *VerifyAccessTokenUC {
	return &VerifyAccessTokenUC{
		jwtService: jwtService,
		userService: userService,
	}
}


func (useCase *VerifyAccessTokenUC) Execute(ctx context.Context, token string) (*domain.AuthUser, error) {
	claims, err := useCase.getClaimsFromToken(token)

	if err != nil {
		return nil, err
	}

	ID, err := useCase.getUserIdFromClaims(claims)

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	user, err := useCase.userService.GetUserByID(ctx, ID)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidToken
		}

		return nil, err
	}

	return user, nil
}

func (useCase *VerifyAccessTokenUC) getClaimsFromToken(token string) (jwt.MapClaims, error) {
	if strings.HasPrefix(token, "bearer ") {
		token = strings.Replace(token, "bearer ", "", 1)
	}

	claims, err := useCase.jwtService.ValidateToken(token)

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (useCase *VerifyAccessTokenUC) getUserIdFromClaims(claims jwt.MapClaims) (uuid.UUID, error) {
	idString, err := claims.GetSubject()

	if err != nil {
		return uuid.Nil, err
	}

	ID, err := uuid.Parse(idString)

	if err != nil {
		return uuid.Nil, err
	}

	return ID, nil
}
