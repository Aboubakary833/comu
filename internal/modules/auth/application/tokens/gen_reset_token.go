package tokens

import (
	"comu/internal/modules/auth/domain"
	"context"
)

type GenerateResetTokenUC struct {
	userService           domain.UserService
	resetTokensRepository domain.ResetTokensRepository
}

func NewGenResetTokenUseCase(
	userService domain.UserService,
	resetTokensRepository domain.ResetTokensRepository,
) *GenerateResetTokenUC {
	return &GenerateResetTokenUC{
		userService:           userService,
		resetTokensRepository: resetTokensRepository,
	}
}

func (useCase *GenerateResetTokenUC) Execute(ctx context.Context, userEmail string) (tokenString string, err error) {
	user, err := useCase.userService.GetUserByEmail(ctx, userEmail)

	if err != nil {
		return
	}

	token := domain.NewResetToken(user.ID, userEmail, domain.DefaultResetTokenTTL)
	err = useCase.resetTokensRepository.Store(ctx, token)

	if err != nil {
		return
	}

	tokenString = token.Token
	return
}
