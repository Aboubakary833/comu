package mysql

import (
	"comu/internal/modules/auth/domain"
	"context"
	"database/sql"
	"errors"
)

type refreshTokensRepository struct {
	db *sql.DB
}

func NewRefreshTokensRepository(db *sql.DB) *refreshTokensRepository {
	return &refreshTokensRepository{
		db: db,
	}
}

func (repo *refreshTokensRepository) Find(ctx context.Context, tokenString string) (*domain.RefreshToken, error) {
	query := "SELECT * FROM refresh_tokens WHERE token = ?"
	token := &domain.RefreshToken{}

	err := repo.db.QueryRowContext(ctx, query, tokenString).Scan(
		&token.UserID, &token.Token, &token.ExpiredAt,
		&token.CreatedAt, &token.Revoked,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTokenNotFound
		}

		return nil, err
	}

	return token, nil
}

func (repo *refreshTokensRepository) Store(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expired_at, created_at, revoked)
		VALUES (UUID_TO_BIN(?), ?, ?, ?, ?)
	`

	_, err := repo.db.ExecContext(
		ctx, query, token.UserID, token.Token,
		token.ExpiredAt, token.CreatedAt, token.Revoked,
	)

	return err
}

func (repo *refreshTokensRepository) Update(ctx context.Context, token *domain.RefreshToken) error {
	query := "UPDATE refresh_tokens SET expired_at = ? WHERE token = ?"
	_, err := repo.db.ExecContext(ctx, query, token.ExpiredAt, token.Token)

	return err
}

func (repo *refreshTokensRepository) Revoke(ctx context.Context, token *domain.RefreshToken) error {
	query := "UPDATE refresh_tokens SET revoked = ? WHERE token = ?"
	_, err := repo.db.ExecContext(ctx, query, token.Revoked, token.Token)

	return err
}
