package database

import (
	"errors"

	"github.com/dhax/go-base/auth"
	"github.com/dhax/go-base/models"
	"github.com/go-pg/pg"
)

var (
	// ErrUniqueEmailConstraint provides error message for already registered email address.
	ErrUniqueEmailConstraint = errors.New("email already registered")
)

// AdmAccountStore implements database operations for account management by admin.
type AdmAccountStore struct {
	db *pg.DB
}

// NewAdmAccountStore returns an AccountStore.
func NewAdmAccountStore(db *pg.DB) *AdmAccountStore {
	return &AdmAccountStore{
		db: db,
	}
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *AdmAccountStore) List(f auth.AccountFilter) ([]auth.Account, int, error) {
	a := []auth.Account{}
	count, err := s.db.Model(&a).
		Apply(f.Filter).
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// Create creates a new account.
func (s *AdmAccountStore) Create(a *auth.Account) error {
	count, _ := s.db.Model(a).
		Where("email = ?email").
		Count()

	if count != 0 {
		return ErrUniqueEmailConstraint
	}

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		err := tx.Insert(a)
		if err != nil {
			return err
		}
		p := &models.Profile{
			AccountID: a.ID,
		}
		return tx.Insert(p)
	})

	return err
}

// Get account by ID.
func (s *AdmAccountStore) Get(id int) (*auth.Account, error) {
	a := auth.Account{ID: id}
	err := s.db.Select(&a)
	return &a, err
}

// Update account.
func (s *AdmAccountStore) Update(a *auth.Account) error {
	err := s.db.Update(a)
	return err
}

// Delete account.
func (s *AdmAccountStore) Delete(a *auth.Account) error {
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
