package notifier

import "comu/internal/modules/notification/domain"

type mailNotifier struct {
	address  string
	port     int
	user     string
	password string
}

func NewMailNotifier(address string, port int, user, password string) *mailNotifier {
	return &mailNotifier{
		address: address,
		port: port,
		user: user,
		password: password,
	}
}

func (notifier *mailNotifier) Send(headers domain.Headers, body domain.Body) error {
	return nil
}
