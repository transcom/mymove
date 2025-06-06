package factory

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
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
		err := db.Where("grade = ?", models.ServiceMemberGradeE1).First(&existingPayGrade)
		if err == nil {
			// PayGrade exists
			rank.PayGradeID = existingPayGrade.ID
		} else {
			log.Panic(fmt.Errorf("database is not configured properly and is missing static hhg allowance and pay grade data. pay grade: %s err: %w", models.ServiceMemberGradeE4, err))
		}

		rank.Affiliation = string(models.DepartmentIndicatorARMY)
		rank.RankAbbv = "SrA"
		rank.RankName = "Senior Airman"

		mustCreate(db, &rank)
	}

	return rank
}
