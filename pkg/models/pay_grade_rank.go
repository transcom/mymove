package models

import "github.com/gofrs/uuid"

type PayGradeRank struct {
	// revert once ready
	// ID            uuid.UUID `db:"id" fk_id:"pay_grade_rank_id" references:"id"`
	// PayGradeID    uuid.UUID `db:"pay_grade_id"`
	// Affiliation   *string   `db:"affiliation"`
	// RankShortName *string   `db:"rank_short_name"`
	// RankName      *string   `db:"rank_name"`
	// RankOrder     *int64    `db:"rank_order"`
	ID            uuid.UUID `db:"-" json:"id,omitempty"`
	PayGradeID    uuid.UUID `db:"-" json:"payGradeId,omitempty"`
	Affiliation   *string   `db:"-" json:"affiliation,omitempty"`
	RankShortName *string   `db:"-" json:"rankShortName,omitempty"`
	RankName      *string   `db:"-" json:"rankName,omitempty"`
	RankOrder     *int64    `db:"-" json:"rankOrder,omitempty"`
}

func (o PayGradeRank) TableName() string {
	return "pay_grade_ranks"
}
