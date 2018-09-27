package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
)

// MakeAccessorial creates a single accessorial record
func MakeAccessorial(db *pop.Connection, assertions Assertions) models.Accessorial {
	shipment := MakeShipment(db, assertions)

	//filled in dummy data
	accessorial := models.Accessorial{
		ShipmentID:    shipment.ID,
		Code:          "105B",
		Item:          "Pack Reg Crate",
		Location:      models.AccessorialLocationDESTINATION,
		Notes:         "Mounted deer head measures 23\" x 34\" x 27\"; crate will be 16.7 cu ft",
		Quantity:      1670,
		Status:        models.AccessorialStatusSUBMITTED,
		SubmittedDate: time.Now(),
		ApprovedDate:  time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Overwrite values with those from assertions
	mergeModels(&accessorial, assertions.Accessorial)

	mustCreate(db, &accessorial)

	return accessorial
}

// MakeDefaultAccessorial makes an Accessorial with default values
func MakeDefaultAccessorial(db *pop.Connection) models.Accessorial {
	return MakeAccessorial(db, Assertions{})
}
