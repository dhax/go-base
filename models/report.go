// Package models contains application specific entities.
package models

import (
	"time"
	"database/sql"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-pg/pg/orm"
	//"github.com/pkg/errors"

)

// Profile holds specific application settings linked to an Account.
type Report struct {
	// Reports represents public.reports
	Id              int   `sql:",pk"`          // id
	AccountID       int   // account_id
	UpdatedAt       time.Time      // date
	Temperature     sql.NullFloat64 // temperature
	Cough           bool            // cough
	RunningNose     bool            // running_nose
	SoreThroat      bool            // sore_throat
	DifficultBreath bool            // difficult_breath
	Headache        bool            // headache
	Diarrhea        bool            // diarrhea
	Nausea          bool            // nausea
	VitA            bool            // vit_a
	VitE            bool            // vit_e
	VitD            bool            // vit_d
	VitC            bool            // vit_c
	Sunbathe        bool            // sunbathe
	Exercise        bool            // exercise
	Veg             bool            // veg
	Fruit           bool            // fruit
	SleepEarly      bool            // sleep_early
	Mask            bool            // mask
	Handwash        bool            // handwash
	Complaint       sql.NullString  // complaint
	Medicine        sql.NullString  // medicine
}

// BeforeInsert hook executed before database insert operation.
func (p *Report) BeforeInsert(db orm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook executed before database update operation.
func (p *Report) BeforeUpdate(db orm.DB) error {
	p.UpdatedAt = time.Now()
	return p.Validate()
}

// Validate validates Profile struct and returns validation errors.
func (p *Report) Validate() error {

	return validation.ValidateStruct(p,
		validation.Field(&p.AccountID, validation.Required),
	)
}
