package database

import (
	"context"
	"time"

	"github.com/dhax/go-base/auth/jwt"
	"github.com/dhax/go-base/auth/pwdless"
	"github.com/uptrace/bun"
)

// AuthStore implements database operations for account pwdlessentication.
type AuthStore struct {
	db *bun.DB
}

// NewAuthStore return an AuthStore.
func NewAuthStore(db *bun.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

// GetAccount returns an account by ID.
func (s *AuthStore) GetAccount(id int) (*pwdless.Account, error) {
	a := &pwdless.Account{ID: id}
	err := s.db.NewSelect().
		Model(a).
		Where("id = ?", id).
		Scan(context.Background())
	return a, err
}

// GetAccountByEmail returns an account by email.
func (s *AuthStore) GetAccountByEmail(e string) (*pwdless.Account, error) {
	a := &pwdless.Account{Email: e}
	err := s.db.NewSelect().
		Model(a).
		Column("id", "active", "email", "name").
		Where("email = ?", e).
		Scan(context.Background())
	return a, err
}

// UpdateAccount upates account data related to pwdlessentication.
func (s *AuthStore) UpdateAccount(a *pwdless.Account) error {
	_, err := s.db.NewUpdate().
		Model(a).
		Column("last_login").
		WherePK().
		Exec(context.Background())
	return err
}

// GetToken returns refresh token by token identifier.
func (s *AuthStore) GetToken(t string) (*jwt.Token, error) {
	token := &jwt.Token{Token: t}
	err := s.db.NewSelect().
		Model(token).
		Where("token = ?", t).
		Scan(context.Background())
	return token, err
}

// CreateOrUpdateToken creates or updates an existing refresh token.
func (s *AuthStore) CreateOrUpdateToken(t *jwt.Token) error {
	if t.ID == 0 {
		_, err := s.db.NewInsert().
			Model(t).
			Exec(context.Background())
		return err
	}
	_, err := s.db.NewUpdate().
		Model(t).
		WherePK().
		Exec(context.Background())
	return err
}

// DeleteToken deletes a refresh token.
func (s *AuthStore) DeleteToken(t *jwt.Token) error {
	_, err := s.db.NewDelete().
		Model(t).
		WherePK().
		Exec(context.Background())
	return err
}

// PurgeExpiredToken deletes expired refresh token.
func (s *AuthStore) PurgeExpiredToken() error {
	_, err := s.db.NewDelete().
		Model((*jwt.Token)(nil)).
		Where("expiry < ?", time.Now()).
		Exec(context.Background())
	return err
}
