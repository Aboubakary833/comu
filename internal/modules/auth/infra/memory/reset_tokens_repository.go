package memory

import (
	"comu/internal/modules/auth/domain"
	"context"
	"sync"
)

type resetTokenStore map[string]domain.ResetToken

type inMemoryResetTokensRepository struct {
	tokens resetTokenStore
	sync.Mutex
}

func NewInMemoryResetTokensRepository(initialStore resetTokenStore) *inMemoryResetTokensRepository {
	if initialStore == nil {
		initialStore = make(resetTokenStore)
	}

	return &inMemoryResetTokensRepository{
		tokens: initialStore,
	}
}

func (repo *inMemoryResetTokensRepository) Find(ctx context.Context, tokenString string) (*domain.ResetToken, error) {
	repo.Lock()
	defer repo.Unlock()

	token, ok := repo.tokens[tokenString]

	if !ok {
		return nil, domain.ErrTokenNotFound
	}

	return &token, nil
}

func (repo *inMemoryResetTokensRepository) Store(ctx context.Context, token *domain.ResetToken) error {
	repo.Lock()
	defer repo.Unlock()

	repo.tokens[token.Token] = *token

	return nil
}

func (repo *inMemoryResetTokensRepository) Delete(ctx context.Context, tokenString string) error {
	if _, err := repo.Find(ctx, tokenString); err != nil {
		return err
	}
	repo.Lock()
	defer repo.Unlock()

	delete(repo.tokens, tokenString)

	return nil
}
