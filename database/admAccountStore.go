package database

import (
	"errors"

	"github.com/dhax/go-base/models"
	"github.com/go-pg/pg"
)

var (
	ErrUniqueEmailConstraint = errors.New("email already registered")
)

type AdmAccountStore struct {
	db *pg.DB
}

func NewAdmAccountStore(db *pg.DB) *AdmAccountStore {
	return &AdmAccountStore{
		db: db,
	}
}

func (s *AdmAccountStore) List(f models.AccountFilter) (*[]models.Account, int, error) {
	var a []models.Account
	count, err := s.db.Model(&a).
		Apply(f.Filter).
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return &a, count, nil
}

func (s *AdmAccountStore) Create(a *models.Account) error {
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

func (s *AdmAccountStore) Get(id int) (*models.Account, error) {
	a := models.Account{ID: id}
	err := s.db.Select(&a)
	return &a, err
}

func (s *AdmAccountStore) Update(a *models.Account) error {
	err := s.db.Update(a)
	return err
}

func (s *AdmAccountStore) Delete(a *models.Account) error {
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
