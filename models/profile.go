package models

import (
	"time"

	"github.com/go-ozzo/ozzo-validation"

	"github.com/go-pg/pg/orm"
)

type Profile struct {
	ID        int       `json:"id,omitempty"`
	AccountID int       `json:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Theme string `json:"theme,omitempty"`
}

func (p *Profile) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
		p.UpdatedAt = now
	}
	return nil
}

func (p *Profile) BeforeUpdate(db orm.DB) error {
	if err := p.Validate(); err != nil {
		return err
	}
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Profile) Validate() error {

	return validation.ValidateStruct(p,
		validation.Field(&p.Theme, validation.Required, validation.In("default", "dark")),
	)
}
