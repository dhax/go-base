package database

import (
	"context"
	"database/sql"
	"errors"
	"net/url"

	"github.com/dhax/go-base/auth/jwt"
	"github.com/dhax/go-base/auth/pwdless"
	"github.com/dhax/go-base/models"
	"github.com/uptrace/bun"
)

var (
	// ErrUniqueEmailConstraint provides error message for already registered email address.
	ErrUniqueEmailConstraint = errors.New("email already registered")
	// ErrBadParams could not parse params to filter
	ErrBadParams = errors.New("bad parameters")
)

// AdmAccountStore implements database operations for account management by admin.
type AdmAccountStore struct {
	db *bun.DB
}

// NewAdmAccountStore returns an AccountStore.
func NewAdmAccountStore(db *bun.DB) *AdmAccountStore {
	return &AdmAccountStore{
		db: db,
	}
}

// AccountFilter provides pagination and filtering options on accounts.
type AccountFilter struct {
	Limit  int
	Offset int
	Filter map[string]interface{}
	Order  []string
}

// NewAccountFilter returns an AccountFilter with options parsed from request url values.
func NewAccountFilter(params interface{}) (*AccountFilter, error) {
	v, ok := params.(url.Values)
	if !ok {
		return nil, ErrBadParams
	}
	f := &AccountFilter{
		Limit:  10, // Default limit
		Offset: 0,  // Default offset
		Filter: make(map[string]interface{}),
		Order:  v["order"],
	}
	// Parse limit and offset
	if limit := v.Get("limit"); limit != "" {
		f.Limit = int(limit[0] - '0')
	}
	if offset := v.Get("offset"); offset != "" {
		f.Offset = int(offset[0] - '0')
	}
	// Parse filters
	for key, values := range v {
		if key != "limit" && key != "offset" && key != "order" {
			f.Filter[key] = values[0]
		}
	}
	return f, nil
}

// Apply applies an AccountFilter on a bun.SelectQuery.
func (f *AccountFilter) Apply(q *bun.SelectQuery) *bun.SelectQuery {
	q = q.Limit(f.Limit).Offset(f.Offset)
	for key, value := range f.Filter {
		q = q.Where("? = ?", bun.Ident(key), value)
	}
	for _, order := range f.Order {
		q = q.Order(order)
	}
	return q
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *AdmAccountStore) List(f *AccountFilter) ([]pwdless.Account, int, error) {
	var a []pwdless.Account
	count, err := s.db.NewSelect().
		Model(&a).
		Apply(f.Apply).
		ScanAndCount(context.Background())
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// Create creates a new account.
func (s *AdmAccountStore) Create(a *pwdless.Account) error {
	exists, err := s.db.NewSelect().
		Model((*pwdless.Account)(nil)).
		Where("email = ?", a.Email).
		Exists(context.Background())
	if err != nil {
		return err
	}

	if exists {
		return ErrUniqueEmailConstraint
	}
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if _, err := tx.NewInsert().
		Model(a).
		Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}
	p := &models.Profile{
		AccountID: a.ID,
	}
	if _, err := tx.NewInsert().
		Model(p).
		Exec(ctx); err != nil {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

// Get account by ID.
func (s *AdmAccountStore) Get(id int) (*pwdless.Account, error) {
	a := &pwdless.Account{ID: id}
	err := s.db.NewSelect().
		Model(a).
		WherePK().
		Scan(context.Background())
	return a, err
}

// Update account.
func (s *AdmAccountStore) Update(a *pwdless.Account) error {
	_, err := s.db.NewUpdate().
		Model(a).
		WherePK().
		Exec(context.Background())
	return err
}

// Delete account.
func (s *AdmAccountStore) Delete(a *pwdless.Account) error {
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
		return err
	}
	tx.Commit()
	return nil
}
