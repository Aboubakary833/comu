package application

import (
	"comu/internal/modules/users/domain"
	"context"
	"time"
)

type markUserEmailAsVerifiedUC struct {
	repo domain.Repository
}

func NewMarkUserEmailAsVerifiedUseCase(repo domain.Repository) *markUserEmailAsVerifiedUC {
	return &markUserEmailAsVerifiedUC{
		repo: repo,
	}
}

func (useCase *markUserEmailAsVerifiedUC) Execute(ctx context.Context, userEmail string) error {
	user, err := useCase.repo.FindByEmail(ctx, userEmail)

	if err != nil {
		return err
	}
	now := time.Now()
	user.EmailVerifiedAt = &now

	return useCase.repo.Update(ctx, user)
}
