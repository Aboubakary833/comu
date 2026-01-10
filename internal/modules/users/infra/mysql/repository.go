package mysql

import (
	"comu/internal/modules/users/domain"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type mySqlRepository struct {
	db *sql.DB
}

func NewMysqlRepository(db *sql.DB) *mySqlRepository {
	return &mySqlRepository{
		db: db,
	}
}

func (repo *mySqlRepository) FindByID(ctx context.Context, ID uuid.UUID) (*domain.User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	user := &domain.User{}

	err := repo.db.QueryRowContext(ctx, query, ID).Scan(
		&user.ID, &user.Name, &user.Email, &user.EmailVerifiedAt, &user.Avatar,
		&user.Active, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (repo *mySqlRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT 
			id, name, email, email_verified_at, avatar,
			active, password, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ?
	`

	user := &domain.User{}

	err := repo.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.EmailVerifiedAt, &user.Avatar,
		&user.Active, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (repo *mySqlRepository) Store(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (
		id, name, email, email_verified_at, avatar, active,
		password, created_at, updated_at, deleted_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	id, err := uuid.NewV7()

	if err != nil {
		return err
	}
	user.ID = id

	_, err = repo.db.ExecContext(
		ctx, query, user.ID, user.Name, user.Email, user.EmailVerifiedAt,
		user.Avatar, user.Active, user.Password, user.CreatedAt,
		user.UpdatedAt, user.DeletedAt,
	)

	return err
}

func (repo *mySqlRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET name = ?, email = ?, email_verified_at = ?,
	avatar = ?, active = ?, password = ?, updated_at = ?,
	deleted_at = ? WHERE id = ?`

	_, err := repo.db.ExecContext(
		ctx, query, user.Name, user.Email, user.EmailVerifiedAt,
		user.Avatar, user.Active, user.Password, user.CreatedAt,
		user.UpdatedAt, user.DeletedAt, user.ID,
	)

	return err
}

func (repo *mySqlRepository) Delete(ctx context.Context, user *domain.User) error {
	deleteTime := time.Now()
	query := "UPDATE users SET deleted_at = ? WHERE id = ?"
	
	user.DeletedAt = &deleteTime
	_, err := repo.db.ExecContext(ctx, query, deleteTime, user.ID)

	return err
}
