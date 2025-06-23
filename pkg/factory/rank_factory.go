package factory

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func FetchOrBuildRank(db *pop.Connection, customs []Customization, traits []Trait) models.Rank {
	customs = setupCustomizations(customs, traits)

	var rank models.Rank
	if result := findValidCustomization(customs, Rank); result != nil {
		rank = result.Model.(models.Rank)
		if result.LinkOnly {
			return rank
		}
	}

	// cannot get a rank unless there is a provided pay grade id or it's provided link only
	if rank.PayGradeID == uuid.Nil && db != nil {
		var existingPayGrade models.PayGrade
		err := db.Where("grade = ?", models.ServiceMemberGradeE1).First(&existingPayGrade)
		if err == nil {
			// PayGrade exists
			rank.PayGradeID = existingPayGrade.ID
		} else {
			log.Panic(fmt.Errorf("database is not configured properly and is missing static hhg allowance and pay grade data. pay grade: %s err: %w", models.ServiceMemberGradeE4, err))
		}
	}
	if rank.Affiliation == "" {
		rank.Affiliation = string(models.DepartmentIndicatorARMY)
	}
	if rank.RankAbbv == "" {
		rank.RankAbbv = "SrA"
	}
	if rank.RankName == "" {
		rank.RankName = "Senior Airman"
	}
	// cannot save to DB without a paygrade ID return local copy only
	if rank.PayGradeID != uuid.Nil {
		mustCreate(db, &rank)
	}
	return rank
}
