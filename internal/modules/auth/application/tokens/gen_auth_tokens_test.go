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

func TestGenerateAuthTokensUseCase(t *testing.T) {

	t.Run("it should fail and return ErrUserNotFound when no user was found with the provided email", func(t *testing.T) {
		jwtService := mockService.NewJwtServiceMock()
		userService := mockService.NewUserServiceMock()
		refreshTokensRepository := memory.NewInMemoryRefreshTokensRepository(nil)
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"

		userService.On("GetUserByEmail", ctx, userEmail).Return(nil, domain.ErrUserNotFound).Once()

		useCase := NewGenAuthTokensUseCase(jwtService, userService, refreshTokensRepository)

		accessToken, refreshToken, err := useCase.Execute(ctx, userEmail)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		userService.AssertExpectations(t)
		jwtService.AssertNotCalled(t, "GenerateToken")
	})

	t.Run("it should succeed and return access and refresh tokens", func(t *testing.T) {
		jwtService := mockService.NewJwtServiceMock()
		userService := mockService.NewUserServiceMock()
		refreshTokensRepository := memory.NewInMemoryRefreshTokensRepository(nil)
		ctx := context.Background()

		userEmail := "johndoe@gmail.com"
		generatedAccessToken := "cyb613GDg42lqkRzP2dY6pzuMhApH2NvaWRjwhbIkBA="

		user := &domain.AuthUser{
			ID:       uuid.New(),
			Name:     "John Doe",
			Email:    userEmail,
			Password: "secret#pass1234",
		}

		userService.On("GetUserByEmail", ctx, userEmail).Return(user, nil).Once()
		jwtService.On("GenerateToken", user).Return(generatedAccessToken, nil).Once()

		useCase := NewGenAuthTokensUseCase(jwtService, userService, refreshTokensRepository)

		accessToken, refreshToken, err := useCase.Execute(ctx, userEmail)
		_assert := assert.New(t)

		if _assert.NoError(err) && _assert.NotEmpty([]string{accessToken, refreshToken}) {

			refreshTokenStruct, err := refreshTokensRepository.Find(ctx, refreshToken)

			if _assert.NoError(err) {
				_assert.Equal(user.ID, refreshTokenStruct.UserID)
				_assert.False(refreshTokenStruct.Expired())
			}
		}
		jwtService.AssertExpectations(t)
		userService.AssertExpectations(t)
	})
}
