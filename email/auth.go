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
	msg := &Mail{
		from:     NewEmail(m.fromName, m.from),
		to:       NewEmail(name, address),
		subject:  "Login Token",
		template: "loginToken",
		content:  content,
	}

	err := m.Send(msg)
	return err
}
