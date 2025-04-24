package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// PayGradeRank represents a customer's pay grade (Including civilian)
type PayGradeRank struct {
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
func (pgr PayGradeRank) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Name: "Affiliation", Field: pgr.Affiliation},
		&validators.StringIsPresent{Name: "RankAbbv", Field: pgr.Affiliation},
		&validators.StringIsPresent{Name: "RankName", Field: pgr.Affiliation},
	), nil
}

// PayGradeRanks is a slice of PayGradeRank
type PayGradeRanks []PayGradeRank

// TableName overrides the table name used by Pop.
func (p PayGradeRank) TableName() string {
	return "pay_grade_ranks"
}

// get pay grade / rank for orders drop down
func GetPayGradeRankDropdownOptions(db *pop.Connection, affiliation string) ([]string, error) {
	var dropdownOptions []string

	err := db.Q().RawQuery(`
		select pgr.id, pgr.rank_abbv ||' / '|| pg.grade
		from pay_grade_ranks pgr
		join pay_grades pg on pgr.pay_grade_id = pg.id
		where affiliation = $1
		order by pgr.rank_order
	`, affiliation).All(&dropdownOptions)
	if err != nil {
		return nil, err
	}

	return dropdownOptions, nil
}
