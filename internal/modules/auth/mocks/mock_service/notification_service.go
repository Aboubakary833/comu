package mockService

import (
	"comu/internal/modules/auth/domain"

	"github.com/stretchr/testify/mock"
)

type notificationServiceMock struct {
	mock.Mock
}

func NewNotificationServiceMock() *notificationServiceMock {
	return new(notificationServiceMock)
}

func (serviceMock *notificationServiceMock) SendOtpCodeMessage(code *domain.OtpCode) error {
	args := serviceMock.Called(code)
	return args.Error(0)
}

func (serviceMock *notificationServiceMock) SendPasswordChangedMessage(userEmail string) error {
	args := serviceMock.Called(userEmail)
	return args.Error(0)
}
