package database

import (
	"github.com/dhax/go-base/models"
	"github.com/go-pg/pg"
)

type AccountStore struct {
	db *pg.DB
}

func NewAccountStore(db *pg.DB) *AccountStore {
	return &AccountStore{
		db: db,
	}
}

func (s *AccountStore) Get(id int) (*models.Account, error) {
	a := models.Account{ID: id}
	err := s.db.Model(&a).
		Where("account.id = ?id").
		Column("account.*", "Profile", "Token").
		First()
	return &a, err
}

func (s *AccountStore) Update(a *models.Account) error {
	_, err := s.db.Model(a).
		Column("email", "name").
		Update()
	return err
}

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

func (s *AccountStore) UpdateToken(t *models.Token) error {
	_, err := s.db.Model(t).
		Column("identifier").
		Update()
	return err
}

func (s *AccountStore) DeleteToken(t *models.Token) error {
	err := s.db.Delete(t)
	return err
}

func (s *AccountStore) UpdateProfile(p *models.Profile) error {
	err := s.db.Update(p)
	return err
}
