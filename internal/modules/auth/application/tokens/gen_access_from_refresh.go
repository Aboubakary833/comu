package tokens

import (
	"comu/internal/modules/auth/domain"
	"context"
)

type genAccessTokenFromRefreshUC struct {
	jwtService domain.JwtService
	userService domain.UserService
	refreshTokensRepository domain.RefreshTokensRepository
}

func NewGenAccessTokenFromRefreshUseCase(
	jwtService domain.JwtService,
	userService domain.UserService,
	tokensRepository domain.RefreshTokensRepository,
	) *genAccessTokenFromRefreshUC {
	return &genAccessTokenFromRefreshUC{
		jwtService: jwtService,
		userService: userService,
		refreshTokensRepository: tokensRepository,
	}
}

func (useCase *genAccessTokenFromRefreshUC) Execute(ctx context.Context, tokenString string) (string, error) {
	token, err := useCase.refreshTokensRepository.Find(ctx, tokenString)

	if err != nil {
		return "", err
	}

	if token.Expired() {
		return "", domain.ErrExpiredToken
	}
	user, err := useCase.userService.GetUserByID(ctx, token.UserID)

	if err != nil {
		return "", err
	}
	jwtToken, err := useCase.jwtService.GenerateToken(user)

	if err != nil {
		return "", err
	}
	// if the token expire in the next 24H, we give it another week of expiration time.
	// NOTE: Handling the update is not necessary here.
	if token.ExpireInNext24H() {
		token.ExpiredAt = token.ExpiredAt.Add(domain.DefaultRefreshTokenTTL - 24)
		useCase.refreshTokensRepository.Update(ctx, token)
	}

	return jwtToken, nil
}
