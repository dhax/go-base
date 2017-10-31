package pwdless

import "github.com/dhax/go-base/auth/jwt"

// MockAuthStore mocks AuthStorer interface.
type MockAuthStore struct {
	GetAccountFn      func(id int) (*Account, error)
	GetAccountInvoked bool

	GetAccountByEmailFn      func(email string) (*Account, error)
	GetAccountByEmailInvoked bool

	UpdateAccountFn      func(a *Account) error
	UpdateAccountInvoked bool

	GetTokenFn      func(token string) (*jwt.Token, error)
	GetTokenInvoked bool

	CreateOrUpdateTokenFn      func(t *jwt.Token) error
	CreateOrUpdateTokenInvoked bool

	DeleteTokenFn      func(t *jwt.Token) error
	DeleteTokenInvoked bool

	PurgeExpiredTokenFn      func() error
	PurgeExpiredTokenInvoked bool
}

// GetAccount mock returns an account by ID.
func (s *MockAuthStore) GetAccount(id int) (*Account, error) {
	s.GetAccountInvoked = true
	return s.GetAccountFn(id)
}

// GetAccountByEmail mock returns an account by email.
func (s *MockAuthStore) GetAccountByEmail(email string) (*Account, error) {
	s.GetAccountByEmailInvoked = true
	return s.GetAccountByEmailFn(email)
}

// UpdateAccount mock upates account data related to authentication.
func (s *MockAuthStore) UpdateAccount(a *Account) error {
	s.UpdateAccountInvoked = true
	return s.UpdateAccountFn(a)
}

// GetToken mock returns an account and refresh token by token identifier.
func (s *MockAuthStore) GetToken(token string) (*jwt.Token, error) {
	s.GetTokenInvoked = true
	return s.GetTokenFn(token)
}

// CreateOrUpdateToken mock creates or updates a refresh token.
func (s *MockAuthStore) CreateOrUpdateToken(t *jwt.Token) error {
	s.CreateOrUpdateTokenInvoked = true
	return s.CreateOrUpdateTokenFn(t)
}

// DeleteToken mock deletes a refresh token.
func (s *MockAuthStore) DeleteToken(t *jwt.Token) error {
	s.DeleteTokenInvoked = true
	return s.DeleteTokenFn(t)
}

// PurgeExpiredToken mock deletes expired refresh token.
func (s *MockAuthStore) PurgeExpiredToken() error {
	s.PurgeExpiredTokenInvoked = true
	return s.PurgeExpiredTokenFn()
}
