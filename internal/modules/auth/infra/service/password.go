package service

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/shared"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	logger *shared.Log
}

func NewPasswordService(logger *shared.Log) *PasswordService {
	return &PasswordService{
		logger: logger,
	}
}

func (service *PasswordService) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}

func (service *PasswordService) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		service.logger.Error.Println(err)
		return "", domain.ErrInternal
	}

	return string(hashBytes), nil
}
