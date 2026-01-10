package memory

import (
	"comu/internal/modules/auth/domain"
	"context"
	"sync"
)

type refreshTokenStore map[string]domain.RefreshToken

type inMemoryRefreshTokensRepository struct {
	tokens refreshTokenStore
	sync.Mutex
}

func NewInMemoryRefreshTokensRepository(initialStore refreshTokenStore) *inMemoryRefreshTokensRepository {
	if initialStore == nil {
		initialStore = make(refreshTokenStore)
	}

	return &inMemoryRefreshTokensRepository{
		tokens: initialStore,
	}
}

func (repo *inMemoryRefreshTokensRepository) Find(ctx context.Context, tokenString string) (*domain.RefreshToken, error) {
	repo.Lock()
	defer repo.Unlock()

	token, ok := repo.tokens[tokenString]

	if !ok {
		return nil, domain.ErrTokenNotFound
	}

	return &token, nil
}

func (repo *inMemoryRefreshTokensRepository) Store(ctx context.Context, token *domain.RefreshToken) error {
	repo.Lock()
	defer repo.Unlock()

	repo.tokens[token.Token] = *token

	return nil
}

func (repo *inMemoryRefreshTokensRepository) Update(ctx context.Context, token *domain.RefreshToken) error {
	repo.Lock()
	defer repo.Unlock()

	repo.tokens[token.Token] = *token

	return nil
}

func (repo *inMemoryRefreshTokensRepository) Revoke(ctx context.Context, tokenString string) error {
	token, err := repo.Find(ctx, tokenString)

	if err != nil {
		return err
	}
	repo.Lock()
	defer repo.Unlock()

	token.Revoked = true
	repo.tokens[tokenString] = *token

	return nil
}
