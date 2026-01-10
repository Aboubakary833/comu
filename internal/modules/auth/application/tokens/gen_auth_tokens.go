package tokens

import (
	"comu/internal/modules/auth/domain"
	"context"
)

type generateAuthTokensUC struct {
	jwtService              domain.JwtService
	userService             domain.UserService
	refreshTokensRepository domain.RefreshTokensRepository
}

func NewGenAuthTokensUseCase(
	jwtService domain.JwtService,
	userService domain.UserService,
	refreshTokensRepository domain.RefreshTokensRepository,
) *generateAuthTokensUC {
	return &generateAuthTokensUC{
		jwtService:              jwtService,
		userService:             userService,
		refreshTokensRepository: refreshTokensRepository,
	}
}

func (useCase *generateAuthTokensUC) Execute(ctx context.Context, userEmail string) (accessToken, refreshToken string, err error) {
	user, err := useCase.userService.GetUserByEmail(ctx, userEmail)

	if err != nil {
		return
	}
	accessToken, err = useCase.jwtService.GenerateToken(user)

	if err != nil {
		return
	}

	newRefreshToken := domain.NewRefreshToken(user.ID, domain.DefaultRefreshTokenTTL)
	err = useCase.refreshTokensRepository.Store(ctx, newRefreshToken)

	if err != nil {
		return
	}
	refreshToken = newRefreshToken.Token

	return
}
