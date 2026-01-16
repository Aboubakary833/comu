package tokens

import (
	"comu/internal/modules/auth/domain"
	mockService "comu/internal/modules/auth/mocks/mock_service"
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestVerifyAccessToken(t *testing.T) {

	t.Run("it should return ErrInvalidToken when token parsing failed or token expired", func(t *testing.T) {
		jwtService := mockService.NewJwtServiceMock()
		userService := mockService.NewUserServiceMock()

		tokenString := "/Vd6cOMwVI8ZUv84fwOVcQSH6nd5bwFYdw3roB4+Pmo="
		inputToken := "bearer " + tokenString

		errs := []error{domain.ErrInvalidToken, domain.ErrExpiredToken}

		for _, tErr := range errs {
			jwtService.On("ValidateToken", tokenString).Return(nil, tErr).Once()

			useCase := NewVerifyAccessTokenUseCase(jwtService, userService)

			_, err := useCase.Execute(context.Background(), inputToken)
			assert.ErrorIs(t, err, tErr)
			jwtService.AssertExpectations(t)
			userService.AssertNotCalled(t, "GetUserByID")
		}
	})

	t.Run("it should return ErrInvalidToken when no user with the provided ID is found", func(t *testing.T) {
		jwtService := mockService.NewJwtServiceMock()
		userService := mockService.NewUserServiceMock()
		ctx := context.Background()

		userID := uuid.New()
		userEmail := "johndoe@gmail.com"
		expirationTime := time.Now().Add(time.Minute * 15)

		jwtClaims := jwt.MapClaims{
			"sub":   userID.String(),
			"email": userEmail,
			"exp":   expirationTime.Unix(),
			"iat":   time.Now().Unix(),
		}

		tokenString := "/Vd6cOMwVI8ZUv84fwOVcQSH6nd5bwFYdw3roB4+Pmo="
		inputToken := "bearer " + tokenString

		jwtService.On("ValidateToken", tokenString).Return(jwtClaims, nil).Once()
		userService.On("GetUserByID", ctx, userID).Return(nil, domain.ErrUserNotFound).Once()

		useCase := NewVerifyAccessTokenUseCase(jwtService, userService)

		_, err := useCase.Execute(context.Background(), inputToken)

		assert.ErrorIs(t, err, domain.ErrInvalidToken)
		jwtService.AssertExpectations(t)
		userService.AssertExpectations(t)
	})

	t.Run("it should succeed and return the authenticated user", func(t *testing.T) {
		_assert := assert.New(t)
		ctx := context.Background()
		jwtService := mockService.NewJwtServiceMock()
		userService := mockService.NewUserServiceMock()

		userID := uuid.New()
		verifiedAt := time.Now().Add(2 * time.Hour * 24)
		user := domain.AuthUser{
			ID:              userID,
			Name:            "John Doe",
			Email:           "johndoe@gmail.com",
			EmailVerifiedAt: &verifiedAt,
		}
		expirationTime := time.Now().Add(time.Minute * 15)

		jwtClaims := jwt.MapClaims{
			"sub":   user.ID.String(),
			"email": user.Email,
			"exp":   expirationTime.Unix(),
			"iat":   time.Now().Unix(),
		}

		tokenString := "/Vd6cOMwVI8ZUv84fwOVcQSH6nd5bwFYdw3roB4+Pmo="
		inputToken := "bearer " + tokenString

		jwtService.On("ValidateToken", tokenString).Return(jwtClaims, nil).Once()
		userService.On("GetUserByID", ctx, userID).Return(&user, nil)

		useCase := NewVerifyAccessTokenUseCase(jwtService, userService)

		u, err := useCase.Execute(context.Background(), inputToken)

		if _assert.NoError(err) {
			_assert.Equal(user.ID, u.ID)
			_assert.Equal(user.Name, u.Name)
			_assert.Equal(user.Email, u.Email)
			_assert.Equal(user.EmailVerifiedAt, u.EmailVerifiedAt)
		}
		jwtService.AssertExpectations(t)
		userService.AssertExpectations(t)
	})
}
