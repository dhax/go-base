package database

import (
	"github.com/dhax/go-base/auth/jwt"
	"github.com/dhax/go-base/auth/pwdless"
	"github.com/dhax/go-base/models"
	"github.com/go-pg/pg"
)

// AccountStore implements database operations for account management by user.
type AccountStore struct {
	db *pg.DB
}

// NewAccountStore returns an AccountStore.
func NewAccountStore(db *pg.DB) *AccountStore {
	return &AccountStore{
		db: db,
	}
}

// Get an account by ID.
func (s *AccountStore) Get(id int) (*pwdless.Account, error) {
	a := pwdless.Account{ID: id}
	err := s.db.Model(&a).
		Where("account.id = ?id").
		Column("account.*", "Token").
		First()
	return &a, err
}

// Update an account.
func (s *AccountStore) Update(a *pwdless.Account) error {
	_, err := s.db.Model(a).
		Column("email", "name").
		WherePK().
		Update()
	return err
}

// Delete an account.
func (s *AccountStore) Delete(a *pwdless.Account) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.Model(&jwt.Token{}).
			Where("account_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(&models.Profile{}).
			Where("account_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		return tx.Delete(a)
	})
	return err
}

// UpdateToken updates a jwt refresh token.
func (s *AccountStore) UpdateToken(t *jwt.Token) error {
	_, err := s.db.Model(t).
		Column("identifier").
		WherePK().
		Update()
	return err
}

// DeleteToken deletes a jwt refresh token.
func (s *AccountStore) DeleteToken(t *jwt.Token) error {
	err := s.db.Delete(t)
	return err
}
