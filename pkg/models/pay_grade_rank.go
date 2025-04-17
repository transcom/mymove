package models

import "github.com/gofrs/uuid"

type PayGradeRank struct {
	ID         uuid.UUID `db:"id" json:"id,omitempty" rw:"r"`
	PayGradeID uuid.UUID `db:"pay_grade_id" json:"payGradeId,omitempty" rw:"r"`
	// PayGrade      PayGrade  `has_one:"pay_grades" belong_to:"pay_grades" fk:"pay_grade_id" json:"grade,omitempty"`
	Affiliation   *string `db:"affiliation" json:"affiliation,omitempty" rw:"r"`
	RankShortName *string `db:"rank_short_name" json:"rankShortName,omitempty" rw:"r"`
	RankName      *string `db:"rank_name" json:"rankName,omitempty" rw:"r"`
	RankOrder     *int64  `db:"rank_order" json:"rankOrder,omitempty" rw:"r"`
}

func (o PayGradeRank) TableName() string {
	return "pay_grade_ranks"
}
