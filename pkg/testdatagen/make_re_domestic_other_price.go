package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func MakeReDomesticOtherPrice(db *pop.Connection, assertions Assertions) models.ReDomesticOtherPrice {
	contract := assertions.ReContract
	if assertions.ReDomesticOtherPrice.Contract.ID != uuid.Nil {
		contract = assertions.ReDomesticOtherPrice.Contract
	}

	if contract.ID == uuid.Nil {
		contract = FetchOrMakeReContract(db, assertions)
	}

	reService := assertions.ReService
	if assertions.ReDomesticOtherPrice.Service.ID != uuid.Nil {
		reService = assertions.ReDomesticOtherPrice.Service
	}

	if reService.ID == uuid.Nil {
		reService = FetchOrMakeReService(db, assertions)
	}

	reDomesticOtherPrice := models.ReDomesticOtherPrice{
		ContractID:   contract.ID,
		Contract:     contract,
		Service:      reService,
		ServiceID:    reService.ID,
		IsPeakPeriod: false,
		PriceCents:   unit.Cents(832),
		Schedule:     1,
	}

	mergeModels(&reDomesticOtherPrice, assertions.ReDomesticOtherPrice)

	mustCreate(db, &reDomesticOtherPrice, assertions.Stub)

	return reDomesticOtherPrice
}

func FetchOrMakeReDomesticOtherPrice(db *pop.Connection, assertions Assertions) models.ReDomesticOtherPrice {
	domesticOtherPrice := assertions.ReDomesticOtherPrice
	var existingReDomesticOtherPrice models.ReDomesticOtherPrice
	if domesticOtherPrice.ContractID != uuid.Nil && domesticOtherPrice.ServiceID != uuid.Nil {
		err := db.Where("contract_id = ? AND service_id = ? and price_cents = ? and is_peak_period = ? and schedule = ?",
			domesticOtherPrice.ContractID, domesticOtherPrice.ServiceID, domesticOtherPrice.PriceCents, domesticOtherPrice.IsPeakPeriod, domesticOtherPrice.Schedule).First(&existingReDomesticOtherPrice)
		if err != nil && err != sql.ErrNoRows {
			log.Panic("unexpected query error looking for existing ReDomesticOtherPrice", err)
		}

		if existingReDomesticOtherPrice.ID != uuid.Nil {
			return existingReDomesticOtherPrice
		}
	}

	return MakeReDomesticOtherPrice(db, assertions)
}
