package users

import (
	"comu/internal/modules/users/domain"
	"comu/internal/modules/users/infra/mysql"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound   = domain.ErrUserNotFound
	ErrUserEmailTaken = domain.ErrUserEmailTaken
)

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
}

type CreateUserResponse struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

type GetUserResponse struct {
	ID              uuid.UUID
	Name            string
	Email           string
	EmailVerifiedAt *time.Time
	Active          bool
	Avatar          string
	Password        string
	CreatedAt       time.Time
	DeletedAt       *time.Time
}

type UpdateUserPasswordRequest struct {
	ID 			uuid.UUID
	NewPassword string
}

type PublicApi interface {
	CreateUser(context.Context, CreateUserRequest) (*CreateUserResponse, error)
	GetUserByID(context.Context, uuid.UUID) (*GetUserResponse, error)
	GetUserByEmail(context.Context, string) (*GetUserResponse, error)
	MarkEmailAsVerified(context.Context, string) error
	UpdateUserPassword(context.Context, UpdateUserPasswordRequest) error
}

type UserModule struct {
	api PublicApi
}

func NewModule(db *sql.DB) *UserModule {
	repo := mysql.NewMysqlRepository(db)
	api := newApi(repo)

	return &UserModule{
		api: api,
	}
}

func (module *UserModule) GetPublicApi() PublicApi {
	return module.api
}
