package register

import (
	"comu/internal/modules/auth/domain"
	"context"
)

type MarkUserAsVerifiedUC struct {
	userService domain.UserService
}

func NewMarkUserAsVerifiedUseCase(userService domain.UserService) *MarkUserAsVerifiedUC {
	return &MarkUserAsVerifiedUC{
		userService: userService,
	}
}

func (useCase *MarkUserAsVerifiedUC) Execute(ctx context.Context, userEmail string) error {
	return useCase.userService.MarkUserEmailAsVerified(ctx, userEmail)
}
