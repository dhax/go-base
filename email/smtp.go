package email

import (
	"github.com/go-mail/mail"
	"github.com/spf13/viper"
)

// SMTPMailer is a SMTP mailer.
type SMTPMailer struct {
	client *mail.Dialer
	from   Email
}

// NewMailer returns a configured SMTP Mailer.
func NewMailer() (Mailer, error) {
	if err := parseTemplates(); err != nil {
		return nil, err
	}

	smtp := struct {
		Host     string
		Port     int
		User     string
		Password string
	}{
		viper.GetString("email_smtp_host"),
		viper.GetInt("email_smtp_port"),
		viper.GetString("email_smtp_user"),
		viper.GetString("email_smtp_password"),
	}

	if smtp.Host == "" {
		return NewMockMailer(), nil
	}

	s := &SMTPMailer{
		client: mail.NewDialer(smtp.Host, smtp.Port, smtp.User, smtp.Password),
		from:   NewEmail(viper.GetString("email_from_name"), viper.GetString("email_from_address")),
	}

	d, err := s.client.Dial()
	if err == nil {
		d.Close()
		return s, nil
	}
	return nil, err
}

// Send sends the mail via smtp.
func (m *SMTPMailer) Send(email Message) error {
	if err := email.parse(); err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetAddressHeader("From", email.From.Address, email.From.Name)
	msg.SetAddressHeader("To", email.To.Address, email.To.Name)
	msg.SetHeader("Subject", email.Subject)
	msg.SetBody("text/plain", email.text)
	msg.AddAlternative("text/html", email.html)

	return m.client.DialAndSend(msg)
}
