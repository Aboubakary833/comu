package service

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/shared/logger"

	"golang.org/x/crypto/bcrypt"
)

type passwordService struct {
	logger *logger.Log
}

func NewPasswordService(logger *logger.Log) *passwordService {
	return &passwordService{
		logger: logger,
	}
}

func (service *passwordService) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}

func (service *passwordService) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		service.logger.Error.Println(err)
		return "", domain.ErrInternal
	}

	return string(hashBytes), nil
}
