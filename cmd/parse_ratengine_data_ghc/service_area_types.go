package main

import (
	"log"
	"strings"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

type domesticServiceArea struct {
	BasePointCity     string
	State             string
	ServiceAreaNumber string
	Zip3s             []string
}

func (dsa *domesticServiceArea) csvHeader() []string {
	header := []string{
		"Base Point City",
		"State",
		"Service Area Number",
		"Zip3's",
	}

	return header
}

func (dsa *domesticServiceArea) toSlice() []string {
	var values []string

	values = append(values, dsa.BasePointCity)
	values = append(values, dsa.State)
	values = append(values, dsa.ServiceAreaNumber)
	values = append(values, strings.Join(dsa.Zip3s, ","))

	return values
}

func (dsa *domesticServiceArea) saveToDatabase(db *pop.Connection) {
	// need to turn dsa into re_zip3 and re_domestic_service_area
	rdsa := models.ReDomesticServiceArea{
		BasePointCity:    dsa.BasePointCity,
		State:            dsa.State,
		ServiceArea:      dsa.ServiceAreaNumber,
		ServicesSchedule: 2, // TODO Need to look up or parse out the ServicesSchedule
		SITPDSchedule:    2, // TODO Need to look up or parse out the SITPDSchedule
	}
	verrs, err := db.ValidateAndSave(&rdsa)
	if err != nil || verrs.HasAny() {
		var dbError string
		if err != nil {
			dbError = err.Error()
		}
		if verrs.HasAny() {
			dbError = dbError + verrs.Error()
		}
		log.Fatalf("Failed to save Service Area: %v\n  with error: %v\n", rdsa, dbError)
	}
	for _, zip3 := range dsa.Zip3s {
		rz3 := models.ReZip3{
			Zip3:                  zip3,
			DomesticServiceAreaID: rdsa.ID,
		}
		verrs, err = db.ValidateAndSave(&rz3)
		if err != nil || verrs.HasAny() {
			var dbError string
			if err != nil {
				dbError = err.Error()
			}
			if verrs.HasAny() {
				dbError = dbError + verrs.Error()
			}
			log.Fatalf("Failed to save Zip3: %v\n  with error: %v\n", rz3, dbError)
		}
	}
}
