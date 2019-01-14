package database

import (
	"errors"
	"net/url"

	"github.com/dhax/go-base/auth/jwt"
	"github.com/dhax/go-base/auth/pwdless"
	"github.com/dhax/go-base/models"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/urlvalues"
)

var (
	// ErrUniqueEmailConstraint provides error message for already registered email address.
	ErrUniqueEmailConstraint = errors.New("email already registered")
	// ErrBadParams could not parse params to filter
	ErrBadParams = errors.New("bad parameters")
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

// AccountFilter provides pagination and filtering options on accounts.
type AccountFilter struct {
	Pager  *urlvalues.Pager
	Filter *urlvalues.Filter
	Order  []string
}

// NewAccountFilter returns an AccountFilter with options parsed from request url values.
func NewAccountFilter(params interface{}) (*AccountFilter, error) {
	v, ok := params.(url.Values)
	if !ok {
		return nil, ErrBadParams
	}
	p := urlvalues.Values(v)
	f := &AccountFilter{
		Pager:  urlvalues.NewPager(p),
		Filter: urlvalues.NewFilter(p),
		Order:  p["order"],
	}
	return f, nil
}

// Apply applies an AccountFilter on an orm.Query.
func (f *AccountFilter) Apply(q *orm.Query) (*orm.Query, error) {
	q = q.Apply(f.Pager.Pagination)
	q = q.Apply(f.Filter.Filters)
	q = q.Order(f.Order...)
	return q, nil
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *AdmAccountStore) List(f *AccountFilter) ([]pwdless.Account, int, error) {
	a := []pwdless.Account{}
	count, err := s.db.Model(&a).
		Apply(f.Apply).
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// Create creates a new account.
func (s *AdmAccountStore) Create(a *pwdless.Account) error {
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
func (s *AdmAccountStore) Get(id int) (*pwdless.Account, error) {
	a := pwdless.Account{ID: id}
	err := s.db.Select(&a)
	return &a, err
}

// Update account.
func (s *AdmAccountStore) Update(a *pwdless.Account) error {
	err := s.db.Update(a)
	return err
}

// Delete account.
func (s *AdmAccountStore) Delete(a *pwdless.Account) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.Model(&jwt.Token{}).
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
