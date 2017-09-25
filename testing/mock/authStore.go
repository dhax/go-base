package mock

import "github.com/dhax/go-base/models"

type AuthStore struct {
	GetByIDFn      func(id int) (*models.Account, error)
	GetByIDInvoked bool

	GetByEmailFn      func(email string) (*models.Account, error)
	GetByEmailInvoked bool

	GetByRefreshTokenFn      func(token string) (*models.Account, *models.Token, error)
	GetByRefreshTokenInvoked bool

	UpdateAccountFn      func(a *models.Account) error
	UpdateAccountInvoked bool

	SaveRefreshTokenFn      func(u *models.Token) error
	SaveRefreshTokenInvoked bool

	DeleteRefreshTokenFn      func(t *models.Token) error
	DeleteRefreshTokenInvoked bool

	PurgeExpiredTokenFn      func() error
	PurgeExpiredTokenInvoked bool
}

func (s *AuthStore) GetByID(id int) (*models.Account, error) {
	s.GetByIDInvoked = true
	return s.GetByIDFn(id)
}

func (s *AuthStore) GetByEmail(email string) (*models.Account, error) {
	s.GetByEmailInvoked = true
	return s.GetByEmailFn(email)
}

func (s *AuthStore) GetByRefreshToken(token string) (*models.Account, *models.Token, error) {
	s.GetByRefreshTokenInvoked = true
	return s.GetByRefreshTokenFn(token)
}

func (s *AuthStore) UpdateAccount(a *models.Account) error {
	s.UpdateAccountInvoked = true
	return s.UpdateAccountFn(a)
}

func (s *AuthStore) SaveRefreshToken(u *models.Token) error {
	s.SaveRefreshTokenInvoked = true
	return s.SaveRefreshTokenFn(u)
}

func (s *AuthStore) DeleteRefreshToken(t *models.Token) error {
	s.DeleteRefreshTokenInvoked = true
	return s.DeleteRefreshTokenFn(t)
}

func (s *AuthStore) PurgeExpiredToken() error {
	s.PurgeExpiredTokenInvoked = true
	return s.PurgeExpiredTokenFn()
}
