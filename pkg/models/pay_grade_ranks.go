package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
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
func GetPayGradeRankDropdownOptions(db *pop.Connection, affiliation string) ([]*internalmessages.Rank, error) {
	var dropdownOptions []*internalmessages.Rank

	err := db.Q().RawQuery(`
		select
			pay_grade_ranks.rank_abbv || ' / ' || pay_grades.grade as RankGradeName,
			pay_grade_ranks.id,
			pay_grade_ranks.pay_grade_id as PaygradeID,
			pay_grade_ranks.rank_order as RankOrder
		from pay_grade_ranks
		join pay_grades on pay_grade_ranks.pay_grade_id = pay_grades.id
		where affiliation = $1
		order by pay_grade_ranks.rank_order DESC
	`, affiliation).All(&dropdownOptions)
	if err != nil {
		return nil, err
	}

	return dropdownOptions, nil
}
