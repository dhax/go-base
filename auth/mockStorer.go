package auth

// MockStorer mocks Storer interface.
type MockStorer struct {
	GetByIDFn      func(id int) (*Account, error)
	GetByIDInvoked bool

	GetByEmailFn      func(email string) (*Account, error)
	GetByEmailInvoked bool

	GetByRefreshTokenFn      func(token string) (*Account, *Token, error)
	GetByRefreshTokenInvoked bool

	UpdateAccountFn      func(a *Account) error
	UpdateAccountInvoked bool

	SaveRefreshTokenFn      func(t *Token) error
	SaveRefreshTokenInvoked bool

	DeleteRefreshTokenFn      func(t *Token) error
	DeleteRefreshTokenInvoked bool

	PurgeExpiredTokenFn      func() error
	PurgeExpiredTokenInvoked bool
}

// GetByID mock returns an account by ID.
func (s *MockStorer) GetByID(id int) (*Account, error) {
	s.GetByIDInvoked = true
	return s.GetByIDFn(id)
}

// GetByEmail mock returns an account by email.
func (s *MockStorer) GetByEmail(email string) (*Account, error) {
	s.GetByEmailInvoked = true
	return s.GetByEmailFn(email)
}

// GetByRefreshToken mock returns an account and refresh token by token identifier.
func (s *MockStorer) GetByRefreshToken(token string) (*Account, *Token, error) {
	s.GetByRefreshTokenInvoked = true
	return s.GetByRefreshTokenFn(token)
}

// UpdateAccount mock upates account data related to authentication.
func (s *MockStorer) UpdateAccount(a *Account) error {
	s.UpdateAccountInvoked = true
	return s.UpdateAccountFn(a)
}

// SaveRefreshToken mock creates or updates a refresh token.
func (s *MockStorer) SaveRefreshToken(t *Token) error {
	s.SaveRefreshTokenInvoked = true
	return s.SaveRefreshTokenFn(t)
}

// DeleteRefreshToken mock deletes a refresh token.
func (s *MockStorer) DeleteRefreshToken(t *Token) error {
	s.DeleteRefreshTokenInvoked = true
	return s.DeleteRefreshTokenFn(t)
}

// PurgeExpiredToken mock deletes expired refresh token.
func (s *MockStorer) PurgeExpiredToken() error {
	s.PurgeExpiredTokenInvoked = true
	return s.PurgeExpiredTokenFn()
}
