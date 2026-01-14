package application

import (
	"comu/internal/modules/users/domain"
	"context"

	"github.com/google/uuid"
)

type UpdateUserPasswordUC struct {
	repo domain.Repository
}

func NewUpdateUserPasswordUseCase(repo domain.Repository) *UpdateUserPasswordUC {
	return &UpdateUserPasswordUC{
		repo: repo,
	}
}

func (useCase *UpdateUserPasswordUC) Execute(ctx context.Context, userID uuid.UUID, newPassword string) error {
	user, err := useCase.repo.FindByID(ctx, userID)

	if err != nil {
		return err
	}
	user.Password = newPassword
	
	return useCase.repo.Update(ctx, user)
}
