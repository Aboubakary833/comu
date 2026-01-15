package tokens

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/modules/auth/infra/memory"
	mockService "comu/internal/modules/auth/mocks/mock_service"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenAccessTokenFromRefreshUseCase(t *testing.T) {

	t.Run("it should fail and return ErrTokenNotFound", func(t *testing.T) {
		repository := memory.NewInMemoryRefreshTokensRepository(nil)
		jwtService := mockService.NewJwtServiceMock()
		userService := mockService.NewUserServiceMock()

		tokenString := "eC9FIPQgybcC6tCItpKMxZyPrW2qNKP8vxoeWE8Vw/s="

		useCase := NewGenAccessTokenFromRefreshUseCase(jwtService, userService, repository)

		_, err := useCase.Execute(context.Background(), tokenString)
		assert.ErrorIs(t, err, domain.ErrTokenNotFound)
		userService.AssertNotCalled(t, "GetUserByID")
		jwtService.AssertNotCalled(t, "GenerateToken")
	})

	t.Run("it should fail and return ErrExpiredToken", func(t *testing.T) {
		repository := memory.NewInMemoryRefreshTokensRepository(nil)
		jwtService := mockService.NewJwtServiceMock()
		userService := mockService.NewUserServiceMock()
		ctx := context.Background()

		token := domain.NewRefreshToken(uuid.New(), -2*time.Hour)
		repository.Store(ctx, token)

		useCase := NewGenAccessTokenFromRefreshUseCase(jwtService, userService, repository)

		_, err := useCase.Execute(ctx, token.Token)
		assert.ErrorIs(t, err, domain.ErrExpiredToken)
		userService.AssertNotCalled(t, "GetUserByID")
		jwtService.AssertNotCalled(t, "GenerateToken")
	})

	t.Run("it should fail and return ErrUserNotFound", func(t *testing.T) {
		repository := memory.NewInMemoryRefreshTokensRepository(nil)
		jwtService := mockService.NewJwtServiceMock()
		userService := mockService.NewUserServiceMock()
		ctx := context.Background()

		token := domain.NewRefreshToken(uuid.New(), domain.DefaultRefreshTokenTTL)
		repository.Store(ctx, token)

		userService.On("GetUserByID", ctx, token.UserID).Return(nil, domain.ErrUserNotFound).Once()

		useCase := NewGenAccessTokenFromRefreshUseCase(jwtService, userService, repository)

		_, err := useCase.Execute(ctx, token.Token)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		userService.AssertExpectations(t)
		jwtService.AssertNotCalled(t, "GenerateToken")
	})

	t.Run("it should succeed, update the refresh token and return a new access token", func(t *testing.T) {
		repository := memory.NewInMemoryRefreshTokensRepository(nil)
		jwtService := mockService.NewJwtServiceMock()
		userService := mockService.NewUserServiceMock()
		ctx := context.Background()

		accessToken := "9eVRLumOWfl+VW2WVFwz5iW3WYLHhvo0ALP2vk8B4uc="
		user := &domain.AuthUser{
			ID:       uuid.New(),
			Name:     "John Doe",
			Email:    "johndoe@gmail.com",
			Password: "secret#pass1234",
		}
		token := domain.NewRefreshToken(user.ID, time.Hour*22)
		repository.Store(ctx, token)

		jwtService.On("GenerateToken", user).Return(accessToken, nil).Once()
		userService.On("GetUserByID", ctx, token.UserID).Return(user, nil).Once()

		useCase := NewGenAccessTokenFromRefreshUseCase(jwtService, userService, repository)

		generatedToken, err := useCase.Execute(ctx, token.Token)
		_assert := assert.New(t)

		if _assert.NoError(err) {
			jwtService.AssertExpectations(t)
			userService.AssertExpectations(t)
			_assert.Equal(accessToken, generatedToken)

			retrievedRefreshToken, err := repository.Find(ctx, token.Token)

			if _assert.NoError(err) {
				_assert.True(
					time.Now().Add(
						domain.DefaultRefreshTokenTTL - 10*time.Hour,
					).Before(retrievedRefreshToken.ExpiredAt),
				)
			}
		}
	})
}
