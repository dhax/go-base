package mock

import "github.com/dhax/go-base/email"

type Mailer struct {
	LoginTokenFn      func(name, email string, c email.ContentLoginToken) error
	LoginTokenInvoked bool
}

func (s *Mailer) LoginToken(n, e string, c email.ContentLoginToken) error {
	s.LoginTokenInvoked = true
	return s.LoginTokenFn(n, e, c)
}
