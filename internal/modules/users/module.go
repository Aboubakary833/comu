package users

import (
	"comu/internal/modules/users/application"
	"comu/internal/modules/users/domain"
	"comu/internal/modules/users/infra/mysql"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound   = domain.ErrUserNotFound
	ErrUserEmailTaken = domain.ErrUserEmailTaken
)

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
	repo := mysql.NewRepository(db)

	//usecases
	createUserUC := application.NewCreateUserUseCase(repo)
	getUserByIdUC := application.NewGetUserByIdUseCase(repo)
	getUserByEmailUC := application.NewGetUserByEmailUseCase(repo)
	updateUserPasswordUC := application.NewUpdateUserPasswordUseCase(repo)
	markUserEmailAsVerifiedUC := application.NewMarkUserEmailAsVerifiedUseCase(repo)

	api := newApi(
		createUserUC, getUserByIdUC, getUserByEmailUC,
		updateUserPasswordUC, markUserEmailAsVerifiedUC,
	)

	return &UserModule{
		api: api,
	}
}

func (module *UserModule) GetPublicApi() PublicApi {
	return module.api
}
