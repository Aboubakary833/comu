package mysql

import (
	"comu/internal/modules/auth/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type resendOtpRequestsRepository struct {
	db *sql.DB
}

func NewResendOtpRequestsRepository(db *sql.DB) *resendOtpRequestsRepository {
	return &resendOtpRequestsRepository{
		db: db,
	}
}

func (repo *resendOtpRequestsRepository) findQuery(ctx context.Context, column, value string) (*domain.ResendOtpRequest, error) {
	queryVal := "?"

	if column == "id" {
		queryVal = "UUID_TO_BIN(?)"
	}

	query := fmt.Sprintf("SELECT * FROM resend_otp_requests WHERE %s = %s", column, queryVal)
	req := &domain.ResendOtpRequest{}

	err := repo.db.QueryRowContext(ctx, query, value).Scan(
		&req.ID, &req.UserEmail, &req.Count,
		&req.LastSendAt, &req.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResendRequestNotFound
		}

		return nil, err
	}

	return req, nil
}

func (repo *resendOtpRequestsRepository) FindByID(ctx context.Context, ID uuid.UUID) (*domain.ResendOtpRequest, error) {
	return repo.findQuery(ctx, "id", ID.String())
}

func (repo *resendOtpRequestsRepository) FindByUserEmail(ctx context.Context, userEmail string) (*domain.ResendOtpRequest, error) {
	return repo.findQuery(ctx, "user_email", userEmail)
}

func (repo *resendOtpRequestsRepository) Store(ctx context.Context, req *domain.ResendOtpRequest) error {
	query := `
		INSERT INTO resend_otp_requests (id, user_email, count, last_sent_at, created_at)
		VALUES (UUID_TO_BIN(?), ?, ?, ?, ?)
	`

	_, err := repo.db.ExecContext(
		ctx, query, req.ID.String(), req.UserEmail,
		req.Count, req.LastSendAt, req.CreatedAt,
	)

	return err
}

func (repo *resendOtpRequestsRepository) CreateNew(ctx context.Context, userEmail string) error {
	req := domain.NewResendOtpRequest(userEmail)
	return repo.Store(ctx, req)
}

func (repo *resendOtpRequestsRepository) Delete(ctx context.Context, req *domain.ResendOtpRequest) error {
	query := "DELETE FROM resend_otp_requests WHERE id = UUID_TO_BIN(?)"
	_, err := repo.db.ExecContext(ctx, query, req.ID.String())

	return err
}

func (repo *resendOtpRequestsRepository) IncrementCount(ctx context.Context, req *domain.ResendOtpRequest) error {
	query := "UPDATE resend_otp_requests SET count = ? WHERE id = UUID_TO_BIN(?)"
	req.Count += 1

	_, err := repo.db.ExecContext(ctx, query, req.Count, req.ID.String())

	return err
}
