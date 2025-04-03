package models

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

type PaygradeRank struct {
	ID           strfmt.UUID `db:"id" fk_id:"pay_grade_rank_id" references:"id"`
	PaygradeId   strfmt.UUID `db:"pay_grade_id"`
	Affiliation  *string     `db:"affiliation"`
	RankNameAbbv *string     `db:"rank_short_name"`
	RankName     *string     `db:"rank_name"`
	RankOrder    *int64      `db:"rank_order"`
}

func (o PaygradeRank) TableName() string {
	return "pay_grade_ranks"
}

func (o PaygradeRank) FormatToRankPayload() *internalmessages.Rank {
	var rank = &internalmessages.Rank{}

	rank.Affiliation = o.Affiliation
	rank.ID = o.ID
	rank.PaygradeID = o.PaygradeId
	rank.RankShortName = o.RankNameAbbv
	rank.RankName = o.RankName

	rank.RankOrder = o.RankOrder

	return rank
}
