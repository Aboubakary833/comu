package mysql

import (
	"comu/internal/modules/auth/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type otpCodesRepository struct {
	db *sql.DB
}

func NewOtpCodesRepository(db *sql.DB) *otpCodesRepository {
	return &otpCodesRepository{
		db: db,
	}
}

func (repo *otpCodesRepository) findQuery(ctx context.Context, column, value string) (*domain.OtpCode, error) {
	query := fmt.Sprintf("SELECT * FROM otp_codes WHERE %s = ?", column)
	otpCode := &domain.OtpCode{}

	err := repo.db.QueryRowContext(ctx, query, value).Scan(
		&otpCode.Type, &otpCode.UserEmail,
		&otpCode.Value, &otpCode.ExpiredAt, &otpCode.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOtpNotFound
		}

		return nil, err
	}

	return otpCode, nil
}

func (repo *otpCodesRepository) Find(ctx context.Context, value string) (*domain.OtpCode, error) {
	return repo.findQuery(ctx, "value", value)
}

func (repo *otpCodesRepository) FindByUserEmail(ctx context.Context, userEmail string) (*domain.OtpCode, error) {
	return repo.findQuery(ctx, "user_email", userEmail)
}

func (repo *otpCodesRepository) Exists(ctx context.Context, value string) bool {
	_, err := repo.Find(ctx, value)
	return err == nil
}

func (repo *otpCodesRepository) Store(ctx context.Context, otpCode *domain.OtpCode) error {
	query := `
		INSERT INTO otp_codes (type, user_email, value, expired_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := repo.db.ExecContext(
		ctx, query, otpCode.Type, otpCode.UserEmail,
		otpCode.Value, otpCode.ExpiredAt, otpCode.CreatedAt,
	)

	return err
}

func (repo *otpCodesRepository) CreateWithUserEmail(ctx context.Context, otpType domain.OtpType, email string) (*domain.OtpCode, error) {
	otpCode := domain.NewOtpCode(otpType, email, domain.DefaultOtpCodeTTL)
	err := repo.Store(ctx, otpCode)

	if err != nil {
		return nil, err
	}

	return otpCode, nil
}

func (repo *otpCodesRepository) Delete(ctx context.Context, otpCode *domain.OtpCode) error {
	query := "DELETE FROM otp_codes WHERE value = ?"
	_, err := repo.db.ExecContext(ctx, query, otpCode.Value)

	return err
}
