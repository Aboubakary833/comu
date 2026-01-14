package application

import (
	"comu/internal/modules/users/domain"
	"context"
)

type CreateUserUC struct {
	repo domain.Repository
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

func NewCreateUserUseCase(repo domain.Repository) *CreateUserUC {
	return &CreateUserUC{
		repo: repo,
	}
}

func (useCase *CreateUserUC) Execute(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	newUser := domain.NewUser(input.Name, input.Email, input.Password)
	err := useCase.repo.Store(ctx, newUser)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}
