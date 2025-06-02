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
	ID               uuid.UUID `json:"id" db:"id"`
	Grade            string    `json:"grade" db:"grade"`
	GradeDescription *string   `json:"grade_description" db:"grade_description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" method
func (pg PayGrade) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Name: "Grade", Field: pg.Grade},
	), nil
}

// PayGrades is a slice of PayGrade
type PayGrades []PayGrade

func GetPayGradesForAffiliation(db *pop.Connection, affiliation string) (PayGrades, error) {
	var payGrades PayGrades

	err := db.Q().
		Join("ranks", "ranks.pay_grade_id = pay_grades.id").
		Where("ranks.affiliation = $1", affiliation).
		GroupBy("pay_grades.id, pay_grades.grade, pay_grades.created_at, pay_grades.updated_at").
		Order("pay_grades.sort_order").
		All(&payGrades)
	if err != nil {
		return nil, err
	}

	return payGrades, nil
}
