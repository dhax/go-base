// Package models contains application specific entities.
package models

import (
	"database/sql"
	"github.com/go-ozzo/ozzo-validation"
	"time"
	//"github.com/pkg/errors"

)

// Profile holds specific application settings linked to an Account.
type Report struct {
	// Reports represents public.reports
	Id              int   `sql:",pk,unique"`          // id
	AccountID       int   // account_id
	UpdatedAt       time.Time      // date
	Complaint       sql.NullString  // complaint
}

// Validate validates Profile struct and returns validation errors.
func (p *Report) Validate() error {

	return validation.ValidateStruct(p,
		validation.Field(&p.AccountID, validation.Required),
	)
}
