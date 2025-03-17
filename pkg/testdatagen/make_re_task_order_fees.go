package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func MakeReTaskOrderFees(db *pop.Connection, assertions Assertions) models.ReTaskOrderFee {
	var contractYear models.ReContractYear
	if assertions.ReTaskOrderFee.ContractYearID == uuid.Nil {
		contractYear = MakeReContractYear(db, assertions)
	} else {
		contractYear = assertions.ReContractYear
	}

	var reService = assertions.ReService

	reTaskOrderFees := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      reService.ID,
		PriceCents:     unit.Cents(100),
	}

	mergeModels(&reTaskOrderFees, assertions.ReTaskOrderFee)

	mustCreate(db, &reTaskOrderFees, assertions.Stub)

	return reTaskOrderFees
}

func FetchOrMakeReTaskOrderFees(db *pop.Connection, assertions Assertions) models.ReTaskOrderFee {
	var contractYear models.ReContractYear
	if assertions.ReTaskOrderFee.ContractYearID == uuid.Nil {
		contractYear = MakeReContractYear(db, assertions)
	} else {
		contractYear = assertions.ReContractYear
	}

	var reService = assertions.ReService

	var existingReTaskOrderFee models.ReTaskOrderFee
	err := db.Where("contract_year_id = ? AND service_id = ? ", contractYear.ID, reService.ID).First(&existingReTaskOrderFee)
	if err != nil && err != sql.ErrNoRows {
		log.Panic("unexpected query error looking for existing ReTaskOrderFees", err)
	}

	if existingReTaskOrderFee.ID != uuid.Nil {
		return existingReTaskOrderFee
	}
	return MakeReTaskOrderFees(db, assertions)
}
