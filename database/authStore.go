package database

import (
	"time"

	"github.com/dhax/go-base/models"
	"github.com/go-pg/pg"
)

type AuthStore struct {
	db *pg.DB
}

func NewAuthStore(db *pg.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

func (s *AuthStore) GetByID(id int) (*models.Account, error) {
	a := models.Account{ID: id}
	err := s.db.Model(&a).
		Column("account.*").
		Where("id = ?id").
		First()
	return &a, err
}

func (s *AuthStore) GetByEmail(e string) (*models.Account, error) {
	a := models.Account{Email: e}
	err := s.db.Model(&a).
		Column("id", "active", "email", "name").
		Where("email = ?email").
		First()
	return &a, err
}

func (s *AuthStore) GetByRefreshToken(t string) (*models.Account, *models.Token, error) {
	token := models.Token{Token: t}
	err := s.db.Model(&token).
		Where("token = ?token").
		First()
	if err != nil {
		return nil, nil, err
	}

	a := models.Account{ID: token.AccountID}
	err = s.db.Model(&a).
		Column("account.*").
		Where("id = ?id").
		First()

	return &a, &token, err
}

func (s *AuthStore) UpdateAccount(a *models.Account) error {
	_, err := s.db.Model(a).
		Column("last_login").
		Update()
	return err
}

func (s *AuthStore) SaveRefreshToken(t *models.Token) error {
	var err error
	if t.ID == 0 {
		err = s.db.Insert(t)
	} else {
		err = s.db.Update(t)
	}
	return err
}

func (s *AuthStore) DeleteRefreshToken(t *models.Token) error {
	err := s.db.Delete(t)
	return err
}

func (s *AuthStore) PurgeExpiredToken() error {
	_, err := s.db.Model(&models.Token{}).
		Where("expiry < ?", time.Now()).
		Delete()

	return err
}
