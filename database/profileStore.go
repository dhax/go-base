package database

import (
	"context"
	"database/sql"

	"github.com/dhax/go-base/models"
	"github.com/uptrace/bun"
)

// ProfileStore implements database operations for profile management.
type ProfileStore struct {
	db *bun.DB
}

// NewProfileStore returns a ProfileStore implementation.
func NewProfileStore(db *bun.DB) *ProfileStore {
	return &ProfileStore{
		db: db,
	}
}

// Get gets an profile by account ID.
func (s *ProfileStore) Get(accountID int) (*models.Profile, error) {
	p := &models.Profile{AccountID: accountID}
	err := s.db.NewSelect().
		Model(p).
		Where("account_id = ?", accountID).
		Scan(context.Background())

	if err == sql.ErrNoRows {
		_, err = s.db.NewInsert().
			Model(p).
			Exec(context.Background())
	}

	return p, err
}

// Update updates profile.
func (s *ProfileStore) Update(p *models.Profile) error {
	_, err := s.db.NewUpdate().
		Model(p).
		WherePK().
		Exec(context.Background())
	return err
}
