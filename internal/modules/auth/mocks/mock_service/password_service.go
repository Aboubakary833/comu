package mockService

import "github.com/stretchr/testify/mock"

type PasswordServiceMock struct {
	mock.Mock
}

func NewPasswordServiceMock() *PasswordServiceMock {
	return new(PasswordServiceMock)
}

func (serviceMock *PasswordServiceMock) Compare(hashedPassword, password string) error {
	args := serviceMock.Called(hashedPassword, password)
	return args.Error(0)
}

func (serviceMock *PasswordServiceMock) Hash(password string) (string, error) {
	args := serviceMock.Called(password)
	return args.String(0), args.Error(1)
}
