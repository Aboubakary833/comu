package mysql

import (
	"comu/internal/modules/auth/domain"
	"context"
	"database/sql"
	"errors"
)

type resetTokensRepository struct {
	db *sql.DB
}

func NewResetTokensRepository(db *sql.DB) *resetTokensRepository {
	return &resetTokensRepository{
		db: db,
	}
}

func (repo *resetTokensRepository) Find(ctx context.Context, tokenString string) (*domain.ResetToken, error) {
	query := "SELECT * FROM reset_tokens WHERE token = ?"
	token := &domain.ResetToken{}

	err := repo.db.QueryRowContext(ctx, query, tokenString).Scan(
		&token.UserID, &token.Token, &token.UserEmail, &token.ExpiredAt, &token.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTokenNotFound
		}

		return nil, err
	}

	return token, nil
}

func (repo *resetTokensRepository) Store(ctx context.Context, token *domain.ResetToken) error {
		query := `
		INSERT INTO reset_tokens (user_id, user_email, token, expired_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := repo.db.ExecContext(
		ctx, query, token.UserID, token.UserEmail,
		token.Token, token.ExpiredAt, token.CreatedAt,
	)

	return err
}

func (repo *resetTokensRepository) Delete(ctx context.Context, token *domain.ResetToken) error {
	query := "DELETE FROM reset_tokens WHERE token = ?"
	_, err := repo.db.ExecContext(ctx, query, token.Token)

	return err
}
