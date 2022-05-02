package testdatagen

import (
	"github.com/gobuffalo/pop/v5"
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
