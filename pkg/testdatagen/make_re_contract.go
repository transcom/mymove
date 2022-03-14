package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReContract creates a single ReContract
func MakeReContract(db *pop.Connection, assertions Assertions) models.ReContract {
	reContract := models.ReContract{
		Code: DefaultContractCode,
		Name: "Test Contract",
	}

	// Overwrite values with those from assertions
	mergeModels(&reContract, assertions.ReContract)

	mustCreate(db, &reContract, assertions.Stub)

	return reContract
}

func FetchOrMakeReContract(db *pop.Connection, assertions Assertions) models.ReContract {
	if assertions.ReContract.Code == "" {
		assertions.ReContract.Code = DefaultContractCode
	}

	var reContract models.ReContract
	err := db.Where("re_contracts.code = ?", assertions.ReContract.Code).First(&reContract)

	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}

	if reContract.ID == uuid.Nil {
		return MakeReContract(db, assertions)
	}

	return reContract
}

// MakeDefaultReContract makes a single ReContract with default values
func MakeDefaultReContract(db *pop.Connection) models.ReContract {
	return MakeReContract(db, Assertions{})
}
