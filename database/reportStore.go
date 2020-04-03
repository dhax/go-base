package database

import (
	"github.com/dhax/go-base/models"
	"github.com/go-pg/pg"
	//"database/sql"
)

// ProfileStore implements database operations for profile management.
type ReportStore struct {
	db *pg.DB
}

// NewProfileStore returns a ProfileStore implementation.
func NewReportStore(db *pg.DB) *ReportStore {
	return &ReportStore{
		db: db,
	}
}

// Get gets an profile by account ID.
func (s *ReportStore) Get(accountID int) (*models.Report, error) {
	p := models.Report{AccountID: accountID}
	_, err := s.db.Model(&p).
		Where("account_id = ?", accountID).
		SelectOrInsert()

	return &p, err
}

// Update updates profile.
func (s *ReportStore) Update(p *models.Report) error {
	err := s.db.Update(p)
	return err
}


// Create inserts the Reports to the database.
func (s *ReportStore) Insert(p *models.Report) error {
	err := s.db.Insert(p)

	//_, err := r.db.Exec(
	//	`INSERT INTO reports (id, user_id, date, temperature, cough, running_nose, sore_throat, difficult_breath, headache, diarrhea, nausea, vit_a, vit_e, vit_d, vit_c, sunbathe, excercise, veg, fruit, sleep_early, mask, handwash, complaint, medicine) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)`,
	//	&r.ID, &r.AccountID, &r.Date, &r.Temperature, &r.Cough, &r.RunningNose, &r.SoreThroat, &r.DifficultBreath, &r.Headache, &r.Diarrhea, &r.Nausea, &r.VitA, &r.VitE, &r.VitD, &r.VitC, &r.Sunbathe, &r.Exercise, &r.Veg, &r.Fruit, &r.SleepEarly, &r.Mask, &r.Handwash, &r.Complaint, &r.Medicine)
	if err != nil {
		return err
	}
	return err
}