package application

import (
	"comu/internal/modules/users/domain"
	"context"
)

type createUserUC struct {
	repo domain.Repository
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

func NewCreateUserUseCase(repo domain.Repository) *createUserUC {
	return &createUserUC{
		repo: repo,
	}
}

func (useCase *createUserUC) Execute(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	newUser := domain.NewUser(input.Name, input.Email, input.Password)
	err := useCase.repo.Store(ctx, newUser)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}
