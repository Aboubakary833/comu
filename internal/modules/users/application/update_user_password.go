package application

import (
	"comu/internal/modules/users/domain"
	"context"

	"github.com/google/uuid"
)

type updateUserPasswordUC struct {
	repo domain.Repository
}

func NewUpdateUserPasswordUseCase(repo domain.Repository) *updateUserPasswordUC {
	return &updateUserPasswordUC{
		repo: repo,
	}
}

func (useCase *updateUserPasswordUC) Execute(ctx context.Context, userID uuid.UUID, newPassword string) error {
	user, err := useCase.repo.FindByID(ctx, userID)

	if err != nil {
		return err
	}
	user.Password = newPassword
	
	return useCase.repo.Update(ctx, user)
}
