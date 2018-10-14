package database

import (
	"time"

	"github.com/dhax/go-base/auth/jwt"
	"github.com/dhax/go-base/auth/pwdless"
	"github.com/go-pg/pg"
)

// AuthStore implements database operations for account pwdlessentication.
type AuthStore struct {
	db *pg.DB
}

// NewAuthStore return an AuthStore.
func NewAuthStore(db *pg.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

// GetAccount returns an account by ID.
func (s *AuthStore) GetAccount(id int) (*pwdless.Account, error) {
	a := pwdless.Account{ID: id}
	err := s.db.Model(&a).
		Column("account.*").
		Where("id = ?id").
		First()
	return &a, err
}

// GetAccountByEmail returns an account by email.
func (s *AuthStore) GetAccountByEmail(e string) (*pwdless.Account, error) {
	a := pwdless.Account{Email: e}
	err := s.db.Model(&a).
		Column("id", "active", "email", "name").
		Where("email = ?email").
		First()
	return &a, err
}

// UpdateAccount upates account data related to pwdlessentication.
func (s *AuthStore) UpdateAccount(a *pwdless.Account) error {
	_, err := s.db.Model(a).
		Column("last_login").
		WherePK().
		Update()
	return err
}

// GetToken returns refresh token by token identifier.
func (s *AuthStore) GetToken(t string) (*jwt.Token, error) {
	token := jwt.Token{Token: t}
	err := s.db.Model(&token).
		Where("token = ?token").
		First()

	return &token, err
}

// CreateOrUpdateToken creates or updates an existing refresh token.
func (s *AuthStore) CreateOrUpdateToken(t *jwt.Token) error {
	var err error
	if t.ID == 0 {
		err = s.db.Insert(t)
	} else {
		err = s.db.Update(t)
	}
	return err
}

// DeleteToken deletes a refresh token.
func (s *AuthStore) DeleteToken(t *jwt.Token) error {
	err := s.db.Delete(t)
	return err
}

// PurgeExpiredToken deletes expired refresh token.
func (s *AuthStore) PurgeExpiredToken() error {
	_, err := s.db.Model(&jwt.Token{}).
		Where("expiry < ?", time.Now()).
		Delete()

	return err
}
