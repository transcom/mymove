package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func FetchOrMakeDomesticServiceAreaPrice(db *pop.Connection, customs []Customization, traits []Trait) models.ReDomesticServiceAreaPrice {
	customs = setupCustomizations(customs, traits)

	var cReDomesticServiceAreaPrice models.ReDomesticServiceAreaPrice
	if result := findValidCustomization(customs, ReDomesticServiceAreaPrice); result != nil {
		cReDomesticServiceAreaPrice = result.Model.(models.ReDomesticServiceAreaPrice)
		if result.LinkOnly {
			return cReDomesticServiceAreaPrice
		}
	}

	// fetch first before creating
	// the contractID, serviceID, peak, and schedule need to be unique
	var reDomesticServiceAreaPrice models.ReDomesticServiceAreaPrice
	if cReDomesticServiceAreaPrice.ContractID != uuid.Nil && cReDomesticServiceAreaPrice.ServiceID != uuid.Nil && cReDomesticServiceAreaPrice.DomesticServiceAreaID != uuid.Nil {
		err := db.Where("contract_id = ? AND service_id = ? AND domestic_service_area_id = ? AND is_peak_period = ?",
			cReDomesticServiceAreaPrice.ContractID,
			cReDomesticServiceAreaPrice.ServiceID,
			cReDomesticServiceAreaPrice.DomesticServiceAreaID,
			cReDomesticServiceAreaPrice.IsPeakPeriod).
			First(&reDomesticServiceAreaPrice)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		} else if err == sql.ErrNoRows {
			// if it isn't found, then we need to create it
			testdatagen.MergeModels(&reDomesticServiceAreaPrice, cReDomesticServiceAreaPrice)
			mustCreate(db, &reDomesticServiceAreaPrice)
		}
		return reDomesticServiceAreaPrice
	}
	return reDomesticServiceAreaPrice
}
