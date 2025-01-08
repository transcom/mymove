package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func FetchOrMakeDomesticOtherPrice(db *pop.Connection, customs []Customization, traits []Trait) models.ReDomesticOtherPrice {
	customs = setupCustomizations(customs, traits)

	var cReDomesticOtherPrice models.ReDomesticOtherPrice
	if result := findValidCustomization(customs, ReDomesticOtherPrice); result != nil {
		cReDomesticOtherPrice = result.Model.(models.ReDomesticOtherPrice)
		if result.LinkOnly {
			return cReDomesticOtherPrice
		}
	}

	// fetch first before creating
	// the contractID, serviceID, peak, and schedule need to be unique
	var reDomesticOtherPrice models.ReDomesticOtherPrice
	if cReDomesticOtherPrice.ContractID != uuid.Nil && cReDomesticOtherPrice.ServiceID != uuid.Nil && cReDomesticOtherPrice.Schedule != 0 {
		err := db.Where("contract_id = ? AND service_id = ? AND is_peak_period = ? AND schedule = ?",
			cReDomesticOtherPrice.ContractID,
			cReDomesticOtherPrice.ServiceID,
			cReDomesticOtherPrice.IsPeakPeriod,
			cReDomesticOtherPrice.Schedule).
			First(&reDomesticOtherPrice)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		}
		return reDomesticOtherPrice
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&reDomesticOtherPrice, cReDomesticOtherPrice)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &reDomesticOtherPrice)
	}
	return reDomesticOtherPrice
}
