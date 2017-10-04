package models

import (
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg/orm"
)

// Account represents an authenticated application user
type Account struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	LastLogin time.Time `json:"last_login,omitempty"`

	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Active bool     `sql:",notnull" json:"active"`
	Roles  []string `pg:",array" json:"roles,omitempty"`

	Profile *Profile `json:"profile,omitempty"`
	Token   []*Token `json:"token,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (a *Account) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
		a.UpdatedAt = now
	}
	if err := a.Validate(); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate hook executed before database update operation.
func (a *Account) BeforeUpdate(db orm.DB) error {
	if err := a.Validate(); err != nil {
		return err
	}
	a.UpdatedAt = time.Now()
	return nil
}

// BeforeDelete hook executed before database delete operation.
func (a *Account) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates Account struct and returns validation errors.
func (a *Account) Validate() error {
	a.Email = strings.TrimSpace(a.Email)
	a.Email = strings.ToLower(a.Email)
	a.Name = strings.TrimSpace(a.Name)

	return validation.ValidateStruct(a,
		validation.Field(&a.Email, validation.Required, is.Email, is.LowerCase),
		validation.Field(&a.Name, validation.Required, is.ASCII),
	)
}

// CanLogin returns true if is user is allowed to login.
func (a *Account) CanLogin() bool {
	return a.Active
}

func (a *Account) Claims() jwtauth.Claims {
	return jwtauth.Claims{
		"id":    a.ID,
		"sub":   a.Name,
		"roles": a.Roles,
	}
}

// AccountFilter provides pagination and filtering options on accounts.
type AccountFilter struct {
	orm.Pager
	Filters url.Values
	Order   []string
}

// Filter applies an AccountFilter on an orm.Query.
func (f *AccountFilter) Filter(q *orm.Query) (*orm.Query, error) {
	q = q.Apply(f.Pager.Paginate)
	q = q.Apply(orm.URLFilters(f.Filters))
	q = q.Order(f.Order...)
	return q, nil
}

// NewAccountFilter returns an AccountFilter with options parsed from request url values.
func NewAccountFilter(v url.Values) AccountFilter {
	var f AccountFilter
	f.SetURLValues(v)
	f.Filters = v
	f.Order = v["order"]
	return f
}
