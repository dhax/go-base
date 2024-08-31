package pwdless

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/uptrace/bun"

	"github.com/dhax/go-base/auth/jwt"
)

// Account represents an authenticated application user
type Account struct {
	ID        int       `bun:"id,pk,autoincrement" json:"id"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	LastLogin time.Time `bun:"last_login" json:"last_login,omitempty"`

	Email  string   `bun:"email,notnull" json:"email"`
	Name   string   `bun:"name,notnull" json:"name"`
	Active bool     `bun:"active,notnull" json:"active"`
	Roles  []string `bun:"roles,array" json:"roles,omitempty"`

	Token []jwt.Token `bun:"rel:has-many" json:"token,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (a *Account) BeforeInsert(db *bun.DB) error {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
		a.UpdatedAt = now
	}
	return a.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (a *Account) BeforeUpdate(db *bun.DB) error {
	a.UpdatedAt = time.Now()
	return a.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (a *Account) BeforeDelete(db *bun.DB) error {
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

// CanLogin returns true if user is allowed to login.
func (a *Account) CanLogin() bool {
	return a.Active
}

// Claims returns the account's claims to be signed
func (a *Account) Claims() jwt.AppClaims {
	return jwt.AppClaims{
		ID:    a.ID,
		Sub:   a.Name,
		Roles: a.Roles,
	}
}
