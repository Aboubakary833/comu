package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound   = errors.New("no user is found in the records")
	ErrUserEmailTaken = errors.New("the provided email is already taken")
)

type User struct {
	ID              uuid.UUID
	Name            string
	Email           string
	EmailVerifiedAt *time.Time
	Avatar          string
	Active          bool
	Password        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

func (user *User) SetID(id uuid.UUID) {
	if user == nil {
		user = &User{}
	}
	user.ID = id
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:            name,
		Email:           email,
		EmailVerifiedAt: nil,
		Avatar:          "",
		Active:          true,
		Password:        password,
	}
}

func (user *User) EmailIsVerified() bool {
	return user.EmailVerifiedAt != nil
}

type Repository interface {
	FindByID(ctx context.Context, ID uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Store(context.Context, *User) error
	Update(context.Context, *User) error
	Delete(context.Context, *User) error
}
