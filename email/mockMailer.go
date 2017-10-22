package email

// MockMailer is a mock Mailer
type MockMailer struct {
	LoginTokenFn      func(name, email string, c ContentLoginToken) error
	LoginTokenInvoked bool
}

// LoginToken is a mock for LoginToken
func (s *MockMailer) LoginToken(n, e string, c ContentLoginToken) error {
	s.LoginTokenInvoked = true
	return s.LoginTokenFn(n, e, c)
}
