package application

import (
	"comu/internal/modules/users/domain"
	"context"

	"github.com/google/uuid"
)

type GetUserByIdUC struct {
	repo domain.Repository
}

type getUserByEmailUC struct {
	repo domain.Repository
}

func NewGetUserByIdUseCase(repo domain.Repository) *GetUserByIdUC {
	return &GetUserByIdUC{
		repo: repo,
	}
}

func NewGetUserByEmailUseCase(repo domain.Repository) *getUserByEmailUC {
	return &getUserByEmailUC{
		repo: repo,
	}
}

func (useCase *GetUserByIdUC) Execute(ctx context.Context, ID uuid.UUID) (*domain.User, error) {
	user, err := useCase.repo.FindByID(ctx, ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (useCase *getUserByEmailUC) Execute(ctx context.Context, email string) (*domain.User, error) {
	user, err := useCase.repo.FindByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return user, nil
}
