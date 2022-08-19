package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func MakeReDomesticServiceAreaPrice(db *pop.Connection, assertions Assertions) models.ReDomesticServiceAreaPrice {

	contract := assertions.ReContract
	if assertions.ReDomesticServiceAreaPrice.Contract.ID != uuid.Nil {
		contract = assertions.ReDomesticServiceAreaPrice.Contract
	}

	if contract.ID == uuid.Nil {
		contract = FetchOrMakeReContract(db, assertions)
		// this may be used by the following MakeReDomesticServiceArea factory
		assertions.ReContract = contract
	}

	domesticServiceArea := assertions.ReDomesticServiceArea
	if assertions.ReDomesticServiceAreaPrice.DomesticServiceArea.ID != uuid.Nil {
		domesticServiceArea = assertions.ReDomesticServiceAreaPrice.DomesticServiceArea
	}

	if domesticServiceArea.ID == uuid.Nil {
		domesticServiceArea = FetchOrMakeReDomesticServiceArea(db, assertions)
	}

	reService := assertions.ReService
	if assertions.ReDomesticServiceAreaPrice.Service.ID != uuid.Nil {
		reService = assertions.ReDomesticServiceAreaPrice.Service
	}

	if reService.ID == uuid.Nil {
		reService = FetchOrMakeReService(db, assertions)
	}

	reDomesticServiceAreaPrice := models.ReDomesticServiceAreaPrice{
		ContractID:            contract.ID,
		Contract:              contract,
		DomesticServiceAreaID: domesticServiceArea.ID,
		DomesticServiceArea:   domesticServiceArea,
		Service:               reService,
		ServiceID:             reService.ID,
		IsPeakPeriod:          false,
		PriceCents:            unit.Cents(832),
	}

	mergeModels(&reDomesticServiceAreaPrice, assertions.ReDomesticServiceAreaPrice)

	mustCreate(db, &reDomesticServiceAreaPrice, assertions.Stub)

	return reDomesticServiceAreaPrice
}

func FetchOrMakeReDomesticServiceAreaPrice(db *pop.Connection, assertions Assertions) models.ReDomesticServiceAreaPrice {
	var existingServiceAreaPrice models.ReDomesticServiceAreaPrice
	serviceAreaPrice := assertions.ReDomesticServiceAreaPrice
	if serviceAreaPrice.ContractID != uuid.Nil && serviceAreaPrice.DomesticServiceAreaID != uuid.Nil && serviceAreaPrice.ServiceID != uuid.Nil {
		err := db.Where("contract_id = ? AND domestic_service_area_id = ? AND service_id = ? AND is_peak_period = ?", serviceAreaPrice.ContractID, serviceAreaPrice.DomesticServiceAreaID, serviceAreaPrice.ServiceID, serviceAreaPrice.IsPeakPeriod).First(&existingServiceAreaPrice)

		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		}

		if existingServiceAreaPrice.ID == uuid.Nil {
			return MakeReDomesticServiceAreaPrice(db, assertions)
		}
	} else {
		return MakeReDomesticServiceAreaPrice(db, assertions)
	}

	return existingServiceAreaPrice
}
