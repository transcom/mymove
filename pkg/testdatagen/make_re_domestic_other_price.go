package testdatagen

import (
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
