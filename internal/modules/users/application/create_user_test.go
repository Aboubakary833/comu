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

func TestCreateUserUseCase(t *testing.T) {

	t.Run("it should successfully create a new user to the repository", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		useCase := application.NewCreateUserUseCase(repo)

		result, _ := useCase.Execute(
			context.Background(),
			application.CreateUserInput{
				Name:     "John Doe",
				Email:    "johndoe@gmail.com",
				Password: "7ySavUthqq1QeQ7XvghiWC4CtV",
			},
		)

		assert.NotNil(t, result)
		assert.NotEqual(t, uuid.Nil.String(), result.ID)
	})

	t.Run("it should failed and return ErrEmailUserTaken", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		useCase := application.NewCreateUserUseCase(repo)
		ctx := context.Background()

		repo.Store(
			ctx,
			domain.NewUser("John Doe", "johndoe@gmail.com", "7ySavUthqq1QeQ7XvghiWC4CtV"),
		)

		result, err := useCase.Execute(
			ctx,
			application.CreateUserInput{
				Name:     "Johnathan Doe",
				Email:    "johndoe@gmail.com",
				Password: "DcyUdK7mfA4F0Xy9LRpcm9",
			},
		)

		_assert := assert.New(t)

		if _assert.Nil(result) {
			_assert.NotNil(err)
			_assert.Equal(domain.ErrUserEmailTaken, err)
		}

	})
}
