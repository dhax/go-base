package database

import (
	"github.com/dhax/go-base/models"
	"github.com/go-pg/pg"
)

// ProfileStore implements database operations for profile management.
type ProfileStore struct {
	db *pg.DB
}

// NewProfileStore returns a ProfileStore implementation.
func NewProfileStore(db *pg.DB) *ProfileStore {
	return &ProfileStore{
		db: db,
	}
}

// Get gets an profile by account ID.
func (s *ProfileStore) Get(accountID int) (*models.Profile, error) {
	p := models.Profile{AccountID: accountID}
	_, err := s.db.Model(&p).
		Where("account_id = ?", accountID).
		SelectOrInsert()

	return &p, err
}

// Update updates profile.
func (s *ProfileStore) Update(p *models.Profile) error {
	err := s.db.Update(p)
	return err
}
