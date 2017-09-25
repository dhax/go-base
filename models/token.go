package models

import (
	"time"

	"github.com/go-pg/pg/orm"
)

type Token struct {
	ID        int       `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	AccountID int       `json:"-"`

	Token      string    `json:"-"`
	Expiry     time.Time `json:"-"`
	Mobile     bool      `sql:",notnull" json:"mobile"`
	Identifier string    `json:"identifier,omitempty"`
}

func (t *Token) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
		t.UpdatedAt = now
	}
	return nil
}

func (t *Token) BeforeUpdate(db orm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
