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

func TestUpdateUserPasswordUseCase(t *testing.T) {

	t.Run("it should fail and return ErrUserNotFound", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		useCase := application.NewUpdateUserPasswordUseCase(repo)

		newPassword := "newSecret#pass1234"

		err := useCase.Execute(context.Background(), uuid.New(), newPassword)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("it should succeed and update the user password", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		ctx := context.Background()

		user := domain.NewUser("John Doe", "johndoe@gmail.com", "Jf9bCEDh6QBFvYsonX0Vc4TGFHQJYXx5X8ro=")
		repo.Store(ctx, user)

		useCase := application.NewUpdateUserPasswordUseCase(repo)

		newPassword := "j3oQQv3RSVJygVeYAHmrtY++2lxkDpVzFXRyM="

		err := useCase.Execute(ctx, user.ID, newPassword)
		_assert := assert.New(t)

		if _assert.NoError(err) {
			retrievedUser, err := repo.FindByID(ctx, user.ID)
			if _assert.NoError(err) && _assert.NotNil(retrievedUser) {
				_assert.Equal(newPassword, retrievedUser.Password)
				_assert.Equal(user.Name, retrievedUser.Name)
				_assert.Equal(user.Email, retrievedUser.Email)
			}
		}
	})
}
