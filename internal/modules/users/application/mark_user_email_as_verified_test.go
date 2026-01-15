package application_test

import (
	"comu/internal/modules/users/application"
	"comu/internal/modules/users/domain"
	"comu/internal/modules/users/infra/memory"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkUserEmailAsVerifiedUseCase(t *testing.T) {
	t.Run("it should fail and return ErrUserNotFound", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"

		useCase := application.NewMarkUserEmailAsVerifiedUseCase(repo)

		err := useCase.Execute(ctx, userEmail)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("it should succeed and update the user password", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		user := domain.NewUser("John Doe", userEmail, "secret#pass1234")
		repo.Store(ctx, user)

		useCase := application.NewMarkUserEmailAsVerifiedUseCase(repo)

		err := useCase.Execute(ctx, userEmail)
		_assert := assert.New(t)

		if _assert.NoError(err) {
			retrievedUser, err := repo.FindByID(ctx, user.ID)
			if _assert.NoError(err) && _assert.NotNil(retrievedUser) {
				_assert.Equal(user.Name, retrievedUser.Name)
				_assert.Equal(user.Email, retrievedUser.Email)
				_assert.True(retrievedUser.EmailIsVerified())
			}
		}
	})
}
