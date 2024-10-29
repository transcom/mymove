package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// FetchOrMakeReService returns the ReService for a given service code, or returns a default (DLH).
func FetchOrMakeReService(db *pop.Connection, assertions Assertions) models.ReService {
	var existingReServices models.ReServices
	code := DefaultServiceCode
	if assertions.ReService.Code != "" {
		code = string(assertions.ReService.Code)
	}
	err := db.Where("code = ?", code).All(&existingReServices)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}

	if len(existingReServices) == 0 {
		return MakeDefaultReService(db)
	}

	return existingReServices[0]
}

// MakeDefaultReService makes a single ReService with default values
func MakeDefaultReService(db *pop.Connection) models.ReService {
	reService := models.ReService{
		Code: "DLH",
		Name: "Domestic linehaul",
	}
	return FetchOrMakeReService(db, Assertions{
		ReService: reService,
	})
}
