package mysql

import (
	"comu/internal/modules/users/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (repo *repository) findQuery(ctx context.Context, column, value string) (*domain.User, error) {
	queryVal := "?"

	if column == "id" {
		queryVal = "UUID_TO_BIN(?)"
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE %s = %s", column, queryVal)
	user := &domain.User{}

	err := repo.db.QueryRowContext(ctx, query, value).Scan(
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

func (repo *repository) FindByID(ctx context.Context, ID uuid.UUID) (*domain.User, error) {
	return repo.findQuery(ctx, "id", ID.String())
}

func (repo *repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return repo.findQuery(ctx, "email", email)
}

func (repo *repository) Store(ctx context.Context, user *domain.User) error {
	if repo.emailIsTaken(ctx, user.Email) {
		return domain.ErrUserEmailTaken
	}

	query := `
	INSERT INTO users (
		id, name, email, email_verified_at, avatar, active,
		password, created_at, updated_at, deleted_at
	) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	id, err := uuid.NewV7()

	if err != nil {
		return err
	}
	user.ID = id
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = repo.db.ExecContext(
		ctx, query, user.ID, user.Name, user.Email, user.EmailVerifiedAt,
		user.Avatar, user.Active, user.Password, user.CreatedAt,
		user.UpdatedAt, user.DeletedAt,
	)

	return err
}

func (repo *repository) Update(ctx context.Context, user *domain.User) error {
	savedUser, err := repo.FindByID(ctx, user.ID)

	if err != nil {
		return err
	}

	if savedUser.Email != user.Email && repo.emailIsTaken(ctx, user.Email) {
		return domain.ErrUserEmailTaken
	}

	query := `UPDATE users SET name = ?, email = ?, email_verified_at = ?,
	avatar = ?, active = ?, password = ?, updated_at = ?,
	deleted_at = ? WHERE id = UUID_TO_BIN(?)`

	user.UpdatedAt = time.Now()

	_, err = repo.db.ExecContext(
		ctx, query, user.Name, user.Email, user.EmailVerifiedAt,
		user.Avatar, user.Active, user.Password,
		user.UpdatedAt, user.DeletedAt, user.ID,
	)

	return err
}

func (repo *repository) Delete(ctx context.Context, user *domain.User) error {
	deleteTime := time.Now()
	query := "UPDATE users SET deleted_at = ? WHERE id = UUID_TO_BIN(?)"

	user.DeletedAt = &deleteTime
	_, err := repo.db.ExecContext(ctx, query, deleteTime, user.ID)

	return err
}

func (repo *repository) emailIsTaken(ctx context.Context, email string) bool {
	user, err := repo.FindByEmail(ctx, email)

	if err != nil && errors.Is(err, domain.ErrUserNotFound) {
		return false
	}

	return user != nil
}
