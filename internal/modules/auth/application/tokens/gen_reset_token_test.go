package tokens

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/infra/memory"
	mockService "comu/internal/modules/auth/mocks/mock_service"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateResetTokenUseCase(t *testing.T) {

	t.Run("it should fail and return ErrUserNotFound", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		resetTokensRepository := memory.NewInMemoryResetTokensRepository(nil)
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"

		userService.On("GetUserByEmail", ctx, userEmail).Return(nil, domain.ErrUserNotFound).Once()

		useCase := NewGenResetTokenUseCase(userService, resetTokensRepository)

		_, err := useCase.Execute(ctx, userEmail)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		userService.AssertExpectations(t)
	})

	t.Run("it should succeed and return the reset token string", func(t *testing.T) {
		userService := mockService.NewUserServiceMock()
		resetTokensRepository := memory.NewInMemoryResetTokensRepository(nil)
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		user := &domain.AuthUser{
			ID:       uuid.New(),
			Name:     "John Doe",
			Email:    userEmail,
			Password: "secret#pass1234",
		}

		userService.On("GetUserByEmail", ctx, userEmail).Return(user, nil).Once()

		useCase := NewGenResetTokenUseCase(userService, resetTokensRepository)

		tokenString, err := useCase.Execute(ctx, userEmail)
		_assert := assert.New(t)

		if _assert.NoError(err) && _assert.NotEmpty(tokenString) {
			token, err := resetTokensRepository.Find(ctx, tokenString)

			if _assert.NoError(err) && _assert.NotNil(token) {
				_assert.Equal(user.ID, token.UserID)
				_assert.False(token.Expired())
			}
		}
	})
}
