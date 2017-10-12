package database

import (
	"github.com/dhax/go-base/auth"
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
func (s *AccountStore) Get(id int) (*auth.Account, error) {
	a := auth.Account{ID: id}
	err := s.db.Model(&a).
		Where("account.id = ?id").
		Column("account.*", "Token").
		First()
	return &a, err
}

// Update an account.
func (s *AccountStore) Update(a *auth.Account) error {
	_, err := s.db.Model(a).
		Column("email", "name").
		Update()
	return err
}

// Delete an account.
func (s *AccountStore) Delete(a *auth.Account) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.Model(&auth.Token{}).
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
func (s *AccountStore) UpdateToken(t *auth.Token) error {
	_, err := s.db.Model(t).
		Column("identifier").
		Update()
	return err
}

// DeleteToken deletes a jwt refresh token.
func (s *AccountStore) DeleteToken(t *auth.Token) error {
	err := s.db.Delete(t)
	return err
}
