package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func MakeReDomesticAccessorialPrice(db *pop.Connection, assertions Assertions) models.ReDomesticAccessorialPrice {
	contract := assertions.ReContract
	if assertions.ReDomesticAccessorialPrice.Contract.ID != uuid.Nil {
		contract = assertions.ReDomesticAccessorialPrice.Contract
	}

	if contract.ID == uuid.Nil {
		contract = FetchOrMakeReContract(db, assertions)
	}

	reService := assertions.ReService
	if assertions.ReDomesticAccessorialPrice.Service.ID != uuid.Nil {
		reService = assertions.ReDomesticAccessorialPrice.Service
	}

	if reService.ID == uuid.Nil {
		reService = FetchReService(db, assertions)
	}

	reDomesticAccessorialPrice := models.ReDomesticAccessorialPrice{
		ContractID:       contract.ID,
		Contract:         contract,
		Service:          reService,
		ServiceID:        reService.ID,
		PerUnitCents:     unit.Cents(832),
		ServicesSchedule: 1,
	}

	mergeModels(&reDomesticAccessorialPrice, assertions.ReDomesticAccessorialPrice)

	mustCreate(db, &reDomesticAccessorialPrice, assertions.Stub)

	return reDomesticAccessorialPrice
}

func FetchOrMakeReDomesticAccessorialPrice(db *pop.Connection, assertions Assertions) models.ReDomesticAccessorialPrice {
	domesticAccessorialPrice := assertions.ReDomesticAccessorialPrice
	var existingReDomesticAccessorialPrice models.ReDomesticAccessorialPrice
	if domesticAccessorialPrice.ContractID != uuid.Nil && domesticAccessorialPrice.ServiceID != uuid.Nil {
		err := db.Where("contract_id = ? AND service_id = ? and services_schedule = ?",
			domesticAccessorialPrice.ContractID, domesticAccessorialPrice.ServiceID, domesticAccessorialPrice.ServicesSchedule).First(&existingReDomesticAccessorialPrice)
		if err != nil && err != sql.ErrNoRows {
			log.Panic("unexpected query error looking for existing ReDomesticAccessorialPrice", err)
		}

		if existingReDomesticAccessorialPrice.ID != uuid.Nil {
			return existingReDomesticAccessorialPrice
		}
	}

	return MakeReDomesticAccessorialPrice(db, assertions)
}
