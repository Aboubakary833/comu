package memory

import (
	"comu/internal/modules/auth/domain"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestResetTokensRepositoryStoreMethod(t *testing.T) {

	t.Run("it should successfully store the reset token", func(t *testing.T) {
		repo := NewInMemoryResetTokensRepository(nil)
		token := domain.NewResetToken(uuid.New(), "marcdoe@gmail.com", domain.DefaultResetTokenTTL)

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

func TestInMemoryResetTokensRepositoryFindMethod(t *testing.T) {

	t.Run("it should successfully retrieve a given reset token from the repository", func(t *testing.T) {
		repo := NewInMemoryResetTokensRepository(nil)
		token := domain.NewResetToken(uuid.New(), "johndoe@gmail.com", domain.DefaultResetTokenTTL)
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

	t.Run("it should fail to retrieved a given reset token from the repository", func(t *testing.T) {
		repo := NewInMemoryResetTokensRepository(nil)
		token := domain.NewResetToken(uuid.New(), "lilidoe@gmail.com", domain.DefaultResetTokenTTL)

		retrievedToken, err := repo.Find(context.Background(), token.Token)

		assert.Nil(t, retrievedToken)
		assert.ErrorIs(t, err, domain.ErrTokenNotFound)
	})
}

func TestInMemoryResetTokensRepositoryDeleteMethod(t *testing.T) {

	t.Run("it should successfully delete a given reset token", func(t *testing.T) {
		repo := NewInMemoryResetTokensRepository(nil)
		token := domain.NewResetToken(uuid.New(), "johndoe@gmail.com", domain.DefaultResetTokenTTL)
		ctx := context.Background()

		repo.Store(ctx, token)

		err := repo.Delete(ctx, token.Token)

		assert.NoError(t, err)

		retrievedToken, err := repo.Find(ctx, token.Token)

		if assert.Nil(t, retrievedToken) {
			assert.ErrorIs(t, err, domain.ErrTokenNotFound)
		}
	})

	t.Run("it should fail delete a given token", func(t *testing.T) {
		repo := NewInMemoryResetTokensRepository(nil)
		token := domain.NewResetToken(uuid.New(), "marcdoe@gmail.com", domain.DefaultResetTokenTTL)

		err := repo.Delete(context.Background(), token.Token)

		assert.ErrorIs(t, err, domain.ErrTokenNotFound)
	})

}
