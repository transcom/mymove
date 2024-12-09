package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// PayGrade represents a customer's pay grade (Including civilian)
type PayGrade struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	Grade            string     `json:"grade" db:"grade"`
	GradeDescription *string    `json:"grade_description" db:"grade_description"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (pg PayGrade) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Name: "Grade", Field: pg.Grade},
	), nil
}

// PayGrades is a slice of PayGrade
type PayGrades []PayGrade
