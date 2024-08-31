package database

import (
	"context"
	"database/sql"

	"github.com/dhax/go-base/auth/jwt"
	"github.com/dhax/go-base/auth/pwdless"
	"github.com/dhax/go-base/models"
	"github.com/uptrace/bun"
)

// AccountStore implements database operations for account management by user.
type AccountStore struct {
	db *bun.DB
}

// NewAccountStore returns an AccountStore.
func NewAccountStore(db *bun.DB) *AccountStore {
	return &AccountStore{
		db: db,
	}
}

// Get an account by ID.
func (s *AccountStore) Get(id int) (*pwdless.Account, error) {
	a := &pwdless.Account{ID: id}
	err := s.db.NewSelect().
		Model(a).
		Where("id = ?", id).
		Relation("Token").
		Scan(context.Background())
	return a, err
}

// Update an account.
func (s *AccountStore) Update(a *pwdless.Account) error {
	_, err := s.db.NewUpdate().
		Model(a).
		Column("email", "name").
		WherePK().
		Exec(context.Background())
	return err
}

// Delete an account.
func (s *AccountStore) Delete(a *pwdless.Account) error {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if _, err := tx.NewDelete().
		Model((*jwt.Token)(nil)).
		Where("account_id = ?", a.ID).
		Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.NewDelete().
		Model((*models.Profile)(nil)).
		Where("account_id = ?", a.ID).
		Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.NewDelete().
		Model(a).
		WherePK().
		Exec(ctx); err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return nil
}

// UpdateToken updates a jwt refresh token.
func (s *AccountStore) UpdateToken(t *jwt.Token) error {
	_, err := s.db.NewUpdate().
		Model(t).
		Column("identifier").
		WherePK().
		Exec(context.Background())
	return err
}

// DeleteToken deletes a jwt refresh token.
func (s *AccountStore) DeleteToken(t *jwt.Token) error {
	_, err := s.db.NewDelete().
		Model(t).
		WherePK().
		Exec(context.Background())
	return err
}
