package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Rank represents a customer's rank (Including civilian)
type Rank struct {
	ID          uuid.UUID `json:"id" db:"id"`
	PayGradeID  uuid.UUID `json:"pay_grade_id" db:"pay_grade_id"`
	Affiliation string    `json:"affiliation" db:"affiliation"`
	RankAbbv    string    `json:"rank_abbv" db:"rank_abbv"`
	RankName    string    `json:"rank_name" db:"rank_name"`
	RankOrder   int       `json:"rank_order" db:"rank_order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" method
func (pgr Rank) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Name: "Affiliation", Field: pgr.Affiliation},
		&validators.StringIsPresent{Name: "RankAbbv", Field: pgr.Affiliation},
		&validators.StringIsPresent{Name: "RankName", Field: pgr.Affiliation},
	), nil
}

// Ranks is a slice of Rank
type Ranks []Rank

// TableName overrides the table name used by Pop.
func (p Rank) TableName() string {
	return "ranks"
}
