package factory

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func BuildRank(db *pop.Connection, customs []Customization, traits []Trait) models.Rank {
	customs = setupCustomizations(customs, traits)

	var cRank models.Rank
	if result := findValidCustomization(customs, Rank); result != nil {
		cRank = result.Model.(models.Rank)
		if result.LinkOnly {
			return cRank
		}
	}

	var rank models.Rank

	if db != nil {
		var existingPayGrade models.PayGrade
		err := db.Where("grade = ?", models.ServiceMemberGradeE4).First(&existingPayGrade)
		if err == nil {
			// PayGrade exists
			rank.PayGradeID = existingPayGrade.ID
		} else {
			log.Panic(fmt.Errorf("database is not configured properly and is missing static hhg allowance and pay grade data. pay grade: %s err: %w", models.ServiceMemberGradeE4, err))
		}

		rank.Affiliation = string(models.AffiliationAIRFORCE)
		rank.RankAbbv = "SrA"
		rank.RankName = "Senior Airman"

		testdatagen.MergeModels(&rank, cRank)

		mustCreate(db, &rank)
	}

	return rank
}

// lookup a rank by rank type, if it doesn't exist make it
func FetchOrBuildRankByPayGradeAndAffiliation(db *pop.Connection, payGrade string, affiliation string) models.Rank {
	var rank models.Rank
	var pg models.PayGrade

	// doesn't work if you try to fetch rank and join pay grade for some reason
	err := db.Q().Where("grade = $1", payGrade).First(&pg)
	if err != nil {
		log.Panic(fmt.Errorf("database is not configured properly and is missing static pay grade data. pay grade: %s err: %w", payGrade, err))
	}

	err = db.Q().Where("pay_grade_id = $1", pg.ID).First(&rank)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(fmt.Errorf("database is not configured properly and is missing static pay grade data. pay grade: %s err: %w", payGrade, err))
	} else if err == nil {
		return rank
	}

	return BuildRank(db, nil, nil)
}
