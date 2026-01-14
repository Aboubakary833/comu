package application

import (
	"comu/internal/modules/users/domain"
	"context"

	"github.com/google/uuid"
)

type UpdateUserInfoUC struct {
	repo domain.Repository
}

type UpdateUserInfoInput struct {
	ID        uuid.UUID
	NewName   string
	NewEmail  string
	NewAvatar string
}

func NewUpdateUserInfoUseCase(repo domain.Repository) *UpdateUserInfoUC {
	return &UpdateUserInfoUC{
		repo: repo,
	}
}

func (useCase *UpdateUserInfoUC) Execute(ctx context.Context, input UpdateUserInfoInput) error {
	user, err := useCase.repo.FindByID(ctx, input.ID)

	if err != nil {
		return err
	}
	if input.NewName != "" {
		user.Name = input.NewName
	}

	if input.NewEmail != "" && user.Email != input.NewEmail {
		user.Email = input.NewEmail
		user.EmailVerifiedAt = nil
	}

	if input.NewAvatar != "" {
		user.Avatar = input.NewAvatar
	}

	err = useCase.repo.Update(ctx, user)

	if err != nil {
		return err
	}

	return err
}
