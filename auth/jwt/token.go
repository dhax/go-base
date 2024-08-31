package jwt

import (
	"time"

	"github.com/uptrace/bun"
)

// Token holds refresh jwt information.
type Token struct {
	ID        int       `bun:"id,pk,autoincrement" json:"id,omitempty"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	AccountID int       `bun:"account_id,notnull" json:"-"`

	Token      string    `bun:"token,notnull" json:"-"`
	Expiry     time.Time `bun:"expiry,notnull" json:"-"`
	Mobile     bool      `bun:"mobile,notnull" json:"mobile"`
	Identifier string    `bun:"identifier" json:"identifier,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (t *Token) BeforeInsert(db *bun.DB) error {
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
		t.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate hook executed before database update operation.
func (t *Token) BeforeUpdate(db *bun.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}

// Claims returns the token claims to be signed
func (t *Token) Claims() RefreshClaims {
	return RefreshClaims{
		ID:    t.ID,
		Token: t.Token,
	}
}
