package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func FetchOrMakeAccessorialOtherPrice(db *pop.Connection, customs []Customization, traits []Trait) models.ReDomesticAccessorialPrice {
	customs = setupCustomizations(customs, traits)

	var cReDomesticAccessorialPrice models.ReDomesticAccessorialPrice
	if result := findValidCustomization(customs, ReDomesticAccessorialPrice); result != nil {
		cReDomesticAccessorialPrice = result.Model.(models.ReDomesticAccessorialPrice)
		if result.LinkOnly {
			return cReDomesticAccessorialPrice
		}
	}

	contractYear := testdatagen.FetchOrMakeReContractYear(db, testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			StartDate: testdatagen.ContractStartDate,
			EndDate:   testdatagen.ContractEndDate,
		},
	})

	// fetch first before creating
	// the contractID, serviceID, peak, and schedule need to be unique
	var reDomesticAccessorialPrice models.ReDomesticAccessorialPrice
	if contractYear.ContractID != uuid.Nil && cReDomesticAccessorialPrice.ServiceID != uuid.Nil && cReDomesticAccessorialPrice.ServicesSchedule != 0 {
		err := db.Where("contract_id = ? AND service_id = ? AND services_schedule = ?",
			contractYear.ContractID,
			cReDomesticAccessorialPrice.ServiceID,
			cReDomesticAccessorialPrice.ServicesSchedule).
			First(&reDomesticAccessorialPrice)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		}
		return reDomesticAccessorialPrice
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&reDomesticAccessorialPrice, cReDomesticAccessorialPrice)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &reDomesticAccessorialPrice)
	}
	return reDomesticAccessorialPrice
}
