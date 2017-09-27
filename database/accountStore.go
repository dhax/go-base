package database

import (
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
func (s *AccountStore) Get(id int) (*models.Account, error) {
	a := models.Account{ID: id}
	err := s.db.Model(&a).
		Where("account.id = ?id").
		Column("account.*", "Profile", "Token").
		First()
	return &a, err
}

// Update an account.
func (s *AccountStore) Update(a *models.Account) error {
	_, err := s.db.Model(a).
		Column("email", "name").
		Update()
	return err
}

// Delete an account.
func (s *AccountStore) Delete(a *models.Account) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.Model(&models.Token{}).
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
func (s *AccountStore) UpdateToken(t *models.Token) error {
	_, err := s.db.Model(t).
		Column("identifier").
		Update()
	return err
}

// DeleteToken deletes a jwt refresh token.
func (s *AccountStore) DeleteToken(t *models.Token) error {
	err := s.db.Delete(t)
	return err
}

// UpdateProfile updates corresponding account profile.
func (s *AccountStore) UpdateProfile(p *models.Profile) error {
	err := s.db.Update(p)
	return err
}
