package register

import (
	"comu/internal/modules/auth/domain"
	"context"
)

type markUserAsVerifiedUC struct {
	userService domain.UserService
}

func NewMarkUserAsVerifiedUseCase(userService domain.UserService) *markUserAsVerifiedUC {
	return &markUserAsVerifiedUC{
		userService: userService,
	}
}

func (useCase *markUserAsVerifiedUC) Execute(ctx context.Context, userEmail string) error {
	return useCase.userService.MarkUserEmailAsVerified(ctx, userEmail)
}
