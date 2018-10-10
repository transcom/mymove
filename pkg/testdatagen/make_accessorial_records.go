package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeDummy400ngItem creates a hardcoded accessorial model
// This should be deprecated quickly once we get the real codes into the db
func MakeDummy400ngItem(db *pop.Connection) models.Tariff400ngItem {
	item := models.Tariff400ngItem{
		Code:             "105B",
		Item:             "Pack Reg Crate",
		DiscountType:     models.Tariff400ngItemDiscountTypeNONE,
		AllowedLocation:  models.Tariff400ngItemAllowedLocationEITHER,
		MeasurementUnit1: models.Tariff400ngItemMeasurementUnitEACH,
		MeasurementUnit2: models.Tariff400ngItemMeasurementUnitNONE,
		RateRefCode:      models.Tariff400ngItemRateRefCodeNONE,
	}

	mustCreate(db, &item)

	return item
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
		accessorial = MakeDummy400ngItem(db)
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
