package domain

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mazen160/go-random"
)

type OtpType = int

const (
	LoginOTP OtpType = iota
	RegisterOTP
	ResetPasswordOTP
)

const (
	DefaultOtpCodeTTL      = time.Minute * 10
	DefaultResetTokenTTL   = time.Minute * 15
	DefaultAccessTokenTTL  = time.Minute * 15
	DefaultRefreshTokenTTL = time.Hour * 24 * 7
)

var (
	ErrTokenNotFound                = errors.New("no refresh token was found")
	ErrInvalidToken                 = errors.New("the provided token is invalid")
	ErrExpiredToken                 = errors.New("the provided token has expired")
	ErrOtpNotFound                  = errors.New("no otp code was found")
	ErrInvalidOtp                   = errors.New("the provided otp code is invalid")
	ErrExpiredOtp                   = errors.New("the provided otp code has expired")
	ErrOtpValueTaken                = errors.New("the provided otp is already taken")
	ErrUserNotFound                 = errors.New("no user was found")
	ErrUserEmailTaken               = errors.New("the provided email is already taken")
	ErrInvalidCredentials           = errors.New("invalid credentials. Please try with another credentials")
	ErrResendRequestNotFound        = errors.New("no resend otp request found")
	ErrInvalidResendRequest         = errors.New("the provided opt resend request data is invalid")
	ErrResendRequestCountExceeded   = errors.New("you exceeded the authorized otp resend request limit")
	ErrResendRequestCantBeProcessed = errors.New("you can't request an otp code resend at the moment")
	ErrInternal                     = errors.New("internal error. Please retry or contact our support team")
)

type AuthUser struct {
	ID              uuid.UUID
	Name            string
	Email           string
	EmailVerifiedAt *time.Time
	Avatar          string
	Active          bool
	Password        string
	CreatedAt       time.Time
	DeletedAt       *time.Time
}

type OtpCode struct {
	Type      OtpType
	UserEmail string
	Value     string
	ExpiredAt time.Time
	CreatedAt time.Time
}

type RefreshToken struct {
	UserID    uuid.UUID
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
	Revoked   bool
}

type ResetToken struct {
	UserID    uuid.UUID
	UserEmail string
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
}

type ResendOtpRequest struct {
	ID         uuid.UUID
	UserEmail  string
	Count      int
	LastSendAt time.Time
	CreatedAt  time.Time
}

func NewRefreshToken(userID uuid.UUID, ttl time.Duration) *RefreshToken {
	token, _ := random.String(64)
	expiredAt := time.Now().Add(ttl)

	return &RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
		Revoked:   false,
	}
}

func NewOtpCode(otpType OtpType, userEmail string, ttl time.Duration) *OtpCode {
	code, _ := random.Random(6, random.Digits, true)
	expiredAt := time.Now().Add(ttl)

	return &OtpCode{
		Type:      otpType,
		UserEmail: userEmail,
		Value:     code,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
	}
}

func NewResetToken(userID uuid.UUID, userEmail string, ttl time.Duration) *ResetToken {
	token, _ := random.String(64)
	expiredAt := time.Now().Add(ttl)

	return &ResetToken{
		UserID:    userID,
		UserEmail: userEmail,
		Token:     token,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
	}
}

func NewResendOtpRequest(userEmail string) *ResendOtpRequest {
	return &ResendOtpRequest{
		ID:         uuid.New(),
		UserEmail:  userEmail,
		Count:      0,
		LastSendAt: time.Now(),
		CreatedAt:  time.Now(),
	}
}

func (req *ResendOtpRequest) CanOtpBeSent() bool {
	coolTimeExpiration := req.LastSendAt.Add(time.Minute * 5)
	return time.Now().Before(coolTimeExpiration)
}

func (req *ResendOtpRequest) IsCountExceeded() bool {
	return req.Count >= 5
}

func (token *RefreshToken) Expired() bool {
	return time.Now().After(token.ExpiredAt)
}

func (token *RefreshToken) ExpireInNext24H() bool {
	return time.Until(token.ExpiredAt) <= time.Hour*24
}

func (otpCode *OtpCode) Expired() bool {
	return time.Now().After(otpCode.ExpiredAt)
}

func (token *ResetToken) Expired() bool {
	return time.Now().After(token.ExpiredAt)
}

type OtpCodesRepository interface {
	Find(context.Context, string) (*OtpCode, error)
	FindByUserEmail(context.Context, string) (*OtpCode, error)
	Store(context.Context, *OtpCode) error
	Exists(context.Context, string) bool
	Delete(context.Context, *OtpCode) error
	CreateWithUserEmail(ctx context.Context, otpType OtpType, email string) (*OtpCode, error)
}

type RefreshTokensRepository interface {
	Find(context.Context, string) (*RefreshToken, error)
	Store(context.Context, *RefreshToken) error
	Update(context.Context, *RefreshToken) error
	Revoke(context.Context, string) error
}

type ResetTokensRepository interface {
	Find(context.Context, string) (*ResetToken, error)
	Store(context.Context, *ResetToken) error
	Delete(context.Context, string) error
}

type ResendOtpRequestsRepository interface {
	FindByID(context.Context, uuid.UUID) (*ResendOtpRequest, error)
	FindByUserEmail(context.Context, string) (*ResendOtpRequest, error)
	IncrementCount(context.Context, *ResendOtpRequest) error
	Store(context.Context, *ResendOtpRequest) error
	CreateNew(ctx context.Context, userEmail string) error
	Delete(context.Context, *ResendOtpRequest) error
}

type UserService interface {
	GetUserByID(context.Context, uuid.UUID) (*AuthUser, error)
	GetUserByEmail(context.Context, string) (*AuthUser, error)
	CreateNewUser(ctx context.Context, name, email, password string) (uuid.UUID, error)
	MarkUserEmailAsVerified(ctx context.Context, userEmail string) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
}

type PasswordService interface {
	Compare(hash, password string) error
	Hash(string) (string, error)
}

type JwtService interface {
	GenerateToken(*AuthUser) (string, error)
	ValidateToken(string) (jwt.Claims, error)
}

type NotificationService interface {
	SendOtpCodeMessage(code *OtpCode) error
	SendPasswordChangedMessage(userEmail string) error
}
