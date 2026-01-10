package application_test

import (
	"comu/internal/modules/users/application"
	"comu/internal/modules/users/domain"
	"comu/internal/modules/users/infra/memory"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {

	t.Run("it should successfully find and return the user data", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		user := domain.NewUser("John Doe", "johndoe@gmail.com", "7ySavUthqq1QeQ7XvghiWC4CtV")
		ctx := context.Background()
		repo.Store(ctx, user)

		useCase := application.NewGetUserByIdUseCase(repo)
		result, _ := useCase.Execute(ctx, user.ID)
		_assert := assert.New(t)

		if _assert.NotNil(result) {
			_assert.Equal(user.ID.String(), result.ID.String())
			_assert.Equal(user.Name, result.Name)
			_assert.Equal(user.Email, result.Email)
			_assert.True(result.Active)
		}
	})

	t.Run("it should fail and return ErrUserNotFound", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		id := uuid.New()

		useCase := application.NewGetUserByIdUseCase(repo)
		result, err := useCase.Execute(context.Background(), id)
		_assert := assert.New(t)

		if _assert.Nil(result) {
			_assert.NotNil(err)
			_assert.Equal(domain.ErrUserNotFound, err)
		}
	})
}

func TestGetUserByEmailCase(t *testing.T) {

	t.Run("it should successfully find and return the user data", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		user := domain.NewUser("John Doe", "johndoe@gmail.com", "7ySavUthqq1QeQ7XvghiWC4CtV")
		ctx := context.Background()
		repo.Store(ctx, user)

		useCase := application.NewGetUserByEmailUseCase(repo)
		result, _ := useCase.Execute(ctx, user.Email)
		_assert := assert.New(t)

		if _assert.NotNil(result) {
			_assert.Equal(user.ID.String(), result.ID.String())
			_assert.Equal(user.Name, result.Name)
			_assert.Equal(user.Email, result.Email)
			_assert.True(result.Active)
		}
	})

	t.Run("it should fail and return ErrUserNotFound", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)

		useCase := application.NewGetUserByEmailUseCase(repo)
		result, err := useCase.Execute(context.Background(), "johndoe@gmail.com")
		_assert := assert.New(t)

		if _assert.Nil(result) {
			_assert.NotNil(err)
			_assert.Equal(domain.ErrUserNotFound, err)
		}
	})
}
