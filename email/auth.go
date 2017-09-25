package email

import "time"

type LoginTokenContent struct {
	Email  string
	Name   string
	URL    string
	Token  string
	Expiry time.Time
}

func (s *EmailService) LoginToken(name, address string, content LoginTokenContent) error {
	msg := &Message{
		from:     NewEmail(s.fromName, s.from),
		to:       NewEmail(name, address),
		subject:  "Login Token",
		template: "loginToken",
		content:  content,
	}

	err := s.send(msg)
	return err
}
