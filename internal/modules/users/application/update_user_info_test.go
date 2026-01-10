package application_test

import (
	"comu/internal/modules/users/application"
	"comu/internal/modules/users/domain"
	"comu/internal/modules/users/infra/memory"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserInfoUseCase(t *testing.T) {

	t.Run("it should successfully update a given user data", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		useCase := application.NewUpdateUserInfoUseCase(repo)
		ctx := context.Background()

		user := domain.NewUser("John Doe", "johndoe@gmail.com", "7ySavUthqq1QeQ7XvghiWC4CtV")
		now := time.Now()
		user.EmailVerifiedAt = &now

		repo.Store(ctx, user)

		err := useCase.Execute(
			ctx,
			application.UpdateUserInfoInput{
				ID:        user.ID,
				NewName:   "Johnathan Doe",
				NewEmail:  "johnathandoe@gmail.com",
				NewAvatar: "new-avatar.png",
			},
		)

		retrievedUser, _ := repo.FindByID(ctx, user.ID)
		_assert := assert.New(t)

		if _assert.Nil(err) {
			_assert.Equal("Johnathan Doe", retrievedUser.Name)
			_assert.Equal("johnathandoe@gmail.com", retrievedUser.Email)
			_assert.Equal("new-avatar.png", retrievedUser.Avatar)
			_assert.Nil(retrievedUser.EmailVerifiedAt)
		}
	})

	t.Run("it should only update user name", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		useCase := application.NewUpdateUserInfoUseCase(repo)
		ctx := context.Background()

		user := domain.NewUser("John Doe", "johndoe@gmail.com", "7ySavUthqq1QeQ7XvghiWC4CtV")
		now := time.Now()
		user.EmailVerifiedAt = &now

		repo.Store(ctx, user)

		err := useCase.Execute(
			ctx, application.UpdateUserInfoInput{
				ID:       user.ID,
				NewName:  "Johnathan Doe",
				NewEmail: user.Email,
			},
		)

		retrievedUser, _ := repo.FindByID(ctx, user.ID)
		_assert := assert.New(t)

		if _assert.Nil(err) {
			_assert.Equal("Johnathan Doe", retrievedUser.Name)
			_assert.Equal(user.Email, retrievedUser.Email)
			_assert.NotNil(retrievedUser.EmailVerifiedAt)
		}
	})

	t.Run("it should fail and return ErrUserNotFound", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		useCase := application.NewUpdateUserInfoUseCase(repo)

		err := useCase.Execute(
			context.Background(),
			application.UpdateUserInfoInput{
				ID:       uuid.New(),
				NewName:  "Johnathan Doe",
				NewEmail: "johndoe@gmail.com",
			},
		)

		assert.Equal(t, domain.ErrUserNotFound, err)
	})

	t.Run("it should fail and return ErrUserEmailTaken", func(t *testing.T) {
		repo := memory.NewInMemoryRepository(nil)
		useCase := application.NewUpdateUserInfoUseCase(repo)
		ctx := context.Background()

		now := time.Now()
		subjectUser := domain.NewUser("Johnathan Doe", "johndoe78@gmail.com", "i7BbberVE7aSOcByNBa")
		subjectUser.EmailVerifiedAt = &now

		repo.Store(ctx, domain.NewUser("John Doe", "johndoe@gmail.com", "g1orRgWLIIBKkxF34dra8"))
		repo.Store(ctx, subjectUser)

		err := useCase.Execute(
			ctx, application.UpdateUserInfoInput{
				ID:       subjectUser.ID,
				NewEmail: "johndoe@gmail.com",
			},
		)

		assert.Equal(t, domain.ErrUserEmailTaken, err)
	})
}
