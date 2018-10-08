package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeDummyAccessorial creates a hardcoded accessorial model
// This should be deprecated quickly once we get the real codes into the db
func MakeDummyAccessorial(db *pop.Connection) models.Accessorial {
	accessorial := models.Accessorial{
		Code:             "105B",
		Item:             "Pack Reg Crate",
		DiscountType:     models.AccessorialDiscountTypeNONE,
		AllowedLocation:  models.AccessorialAllowedLocationEITHER,
		MeasurementUnit1: models.AccessorialMeasurementUnitEACH,
		MeasurementUnit2: models.AccessorialMeasurementUnitNONE,
		RateRefCode:      models.AccessorialRateRefCodeNONE,
	}

	mustCreate(db, &accessorial)

	return accessorial
}

// MakeShipmentAccessorial creates a single accessorial record
func MakeShipmentAccessorial(db *pop.Connection, assertions Assertions) models.ShipmentAccessorial {
	shipmentID := assertions.ShipmentAccessorial.ShipmentID
	if isZeroUUID(shipmentID) {
		shipment := MakeShipment(db, assertions)
		shipmentID = shipment.ID
	}

	accessorial := assertions.ShipmentAccessorial.Accessorial
	if isZeroUUID(accessorial.ID) {
		accessorial = MakeDummyAccessorial(db)
	}

	//filled in dummy data
	shipmentAccessorial := models.ShipmentAccessorial{
		ShipmentID:    shipmentID,
		AccessorialID: accessorial.ID,
		Accessorial:   accessorial,
		Location:      models.ShipmentAccessorialLocationDESTINATION,
		Notes:         "Mounted deer head measures 23\" x 34\" x 27\"; crate will be 16.7 cu ft",
		Quantity1:     unit.BaseQuantity(1670),
		Status:        models.ShipmentAccessorialStatusSUBMITTED,
		SubmittedDate: time.Now(),
		ApprovedDate:  time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Overwrite values with those from assertions
	mergeModels(&shipmentAccessorial, assertions.ShipmentAccessorial)

	mustCreate(db, &shipmentAccessorial)

	return shipmentAccessorial
}

// MakeDefaultShipmentAccessorial makes an Accessorial with default values
func MakeDefaultShipmentAccessorial(db *pop.Connection) models.ShipmentAccessorial {
	return MakeShipmentAccessorial(db, Assertions{})
}
