package models

import "github.com/gofrs/uuid"

type PayGradeRank struct {
	ID            uuid.UUID `db:"id" json:"id,omitempty" rw:"r"`
	PayGradeID    uuid.UUID `db:"pay_grade_id" rw:"r" json:"payGradeId,omitempty"`
	Affiliation   *string   `db:"affiliation" rw:"r" json:"affiliation,omitempty"`
	RankShortName *string   `db:"rank_short_name" rw:"r" json:"rankShortName,omitempty"`
	RankName      *string   `db:"rank_name" rw:"r" json:"rankName,omitempty"`
	RankOrder     *int64    `db:"rank_order" rw:"r" json:"rankOrder,omitempty"`
}

func (o PayGradeRank) TableName() string {
	return "pay_grade_ranks"
}
