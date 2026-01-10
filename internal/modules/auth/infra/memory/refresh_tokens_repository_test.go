package memory

import (
	"comu/internal/modules/auth/domain"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRefreshTokensRepositoryStoreMethod(t *testing.T) {
	t.Run("it successfully store a new refresh token in the repository", func(t *testing.T) {
		repo := NewInMemoryRefreshTokensRepository(nil)
		token := domain.NewRefreshToken(uuid.New(), domain.DefaultRefreshTokenTTL)

		err := repo.Store(context.Background(), token)
		_assert := assert.New(t)

		_assert.NoError(err)

		retrievedToken, ok := repo.tokens[token.Token]

		if _assert.True(ok) {
			_assert.Equal(token.UserID.String(), retrievedToken.UserID.String())
			_assert.Equal(token.ExpiredAt, retrievedToken.ExpiredAt)
		}
	})
}

func TestInMemoryRefreshTokensRepositoryFindMethod(t *testing.T) {

	t.Run("it should successfully retrieve a given refresh token from the repository", func(t *testing.T) {
		repo := NewInMemoryRefreshTokensRepository(nil)
		token := domain.NewRefreshToken(uuid.New(), domain.DefaultRefreshTokenTTL)
		ctx := context.Background()

		repo.Store(ctx, token)

		retrievedToken, err := repo.Find(ctx, token.Token)
		_assert := assert.New(t)

		if _assert.NoError(err) {
			_assert.Equal(token.UserID.String(), retrievedToken.UserID.String())
			_assert.Equal(token.ExpiredAt, retrievedToken.ExpiredAt)
			_assert.Equal(token.CreatedAt, retrievedToken.CreatedAt)
		}
	})

	t.Run("it should fail to retrieved a given token from the repository", func(t *testing.T) {
		repo := NewInMemoryRefreshTokensRepository(nil)
		token := domain.NewRefreshToken(uuid.New(), domain.DefaultRefreshTokenTTL)

		retrievedToken, err := repo.Find(context.Background(), token.Token)

		assert.Nil(t, retrievedToken)
		assert.ErrorIs(t, err, domain.ErrTokenNotFound)
	})
}

func TestInMemoryRefreshTokensRepositoryRevokeMethod(t *testing.T) {

	t.Run("it should successfully revoke a given token", func(t *testing.T) {
		repo := NewInMemoryRefreshTokensRepository(nil)
		token := domain.NewRefreshToken(uuid.New(), domain.DefaultRefreshTokenTTL)
		ctx := context.Background()

		repo.Store(ctx, token)

		err := repo.Revoke(ctx, token.Token)

		assert.NoError(t, err)

		retrievedToken, err := repo.Find(ctx, token.Token)

		if assert.NoError(t, err) {
			assert.True(t, retrievedToken.Revoked)
		}
	})

	t.Run("it should fail revoking a given token", func(t *testing.T) {
		repo := NewInMemoryRefreshTokensRepository(nil)
		token := domain.NewRefreshToken(uuid.New(), domain.DefaultRefreshTokenTTL)

		err := repo.Revoke(context.Background(), token.Token)

		assert.ErrorIs(t, err, domain.ErrTokenNotFound)
	})

}
