package service

import (
	"comu/internal/modules/auth/domain"
	"fmt"

	"github.com/wneessen/go-mail"
)

type SmtpNotificationAuth struct {
	Username string
	Password string
}

type smtpNotificationService struct {
	client *mail.Client
	from   string
}

func NewSmtpNotificationService(host string, port int, mailFrom string, auth SmtpNotificationAuth, enableTLS bool) (*smtpNotificationService, error) {
	mailOptions := []mail.Option{
		mail.WithPort(port),
		mail.WithUsername(auth.Username),
		mail.WithPassword(auth.Password),
	}

	// The default TLSPolicy is TLSMandatory
	if !enableTLS {
		mailOptions = append(mailOptions, mail.WithTLSPolicy(mail.NoTLS))
	}

	client, err := mail.NewClient(host, mailOptions...)

	if err != nil {
		return nil, err
	}

	return &smtpNotificationService{
		client: client,
		from:   mailFrom,
	}, nil
}

func (service *smtpNotificationService) SendOtpCodeMessage(code *domain.OtpCode) error {
	msg, err := service.newMessage(code.UserEmail)

	if err != nil {
		return err
	}

	msg.Subject(service.getOtpCodeSubject(code.Type))
	msg.SetBodyString(
		mail.TypeTextPlain,
		fmt.Sprintf(`
			Your verification code is: %s

			This code is valid for %d minutes.
			If you did not request this code, please ignore this message.

		`, code.Value, int(domain.DefaultOtpCodeTTL.Minutes())),
	)

	return service.client.DialAndSend(msg)
}

func (service *smtpNotificationService) SendPasswordChangedMessage(userEmail string) error {
	msg, err := service.newMessage(userEmail)

	if err != nil {
		return err
	}

	msg.Subject("Your password has been changed")
	msg.SetBodyString(
		mail.TypeTextPlain,
		`
		Your password has been successfully changed.

		If you did not make this change, please contact support immediately.
		`,
	)

	return service.client.DialAndSend(msg)
}

func (service *smtpNotificationService) newMessage(receiverEmail string) (*mail.Msg, error) {
	msg := mail.NewMsg()

	if err := msg.From(service.from); err != nil {
		return nil, err
	}

	if err := msg.To(receiverEmail); err != nil {
		return nil, err
	}

	return msg, nil
}

func (service *smtpNotificationService) getOtpCodeSubject(t domain.OtpType) string {
	if t == domain.LoginOTP {
		return "Login verification"
	}

	if t == domain.RegisterOTP {
		return "Confirm your registration"
	}

	return "Reset password confirmation"
}
