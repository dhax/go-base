package email

import "time"

// ContentLoginToken defines content for login token email template.
type ContentLoginToken struct {
	Email  string
	Name   string
	URL    string
	Token  string
	Expiry time.Time
}

// LoginToken creates and sends a login token email with provided template content.
func (m *Mailer) LoginToken(name, address string, content ContentLoginToken) error {
	msg := &message{
		from:     m.from,
		to:       NewEmail(name, address),
		subject:  "Login Token",
		template: "loginToken",
		content:  content,
	}

	if err := msg.parse(); err != nil {
		return err
	}

	return m.Send(msg)
}
