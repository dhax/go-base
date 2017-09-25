package mock

import "github.com/dhax/go-base/email"

type EmailService struct {
	LoginTokenFn      func(name, email string, c email.LoginTokenContent) error
	LoginTokenInvoked bool
}

func (s *EmailService) LoginToken(n, e string, c email.LoginTokenContent) error {
	s.LoginTokenInvoked = true
	return s.LoginTokenFn(n, e, c)
}
