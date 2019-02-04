package testdatagen

import (
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeShipmentLineItem creates a single shipment line item record with an associated tariff400ngItem
func MakeShipmentLineItem(db *pop.Connection, assertions Assertions) models.ShipmentLineItem {
	shipment := assertions.ShipmentLineItem.Shipment
	if isZeroUUID(shipment.ID) {
		shipment = MakeShipment(db, assertions)
	}

	tariff400ngItem := assertions.ShipmentLineItem.Tariff400ngItem
	if isZeroUUID(tariff400ngItem.ID) {
		tariff400ngItem = MakeTariff400ngItem(db, assertions)
	}
	var rate unit.Millicents
	rate = 2354000
	//filled in dummy data
	shipmentLineItem := models.ShipmentLineItem{
		ShipmentID:        shipment.ID,
		Shipment:          shipment,
		Tariff400ngItemID: tariff400ngItem.ID,
		Tariff400ngItem:   tariff400ngItem,
		Location:          models.ShipmentLineItemLocationDESTINATION,
		Notes:             "Mounted deer head measures 23\" x 34\" x 27\"; crate will be 16.7 cu ft",
		Quantity1:         unit.BaseQuantity(1670),
		AppliedRate:       &rate,
		Status:            models.ShipmentLineItemStatusSUBMITTED,
		SubmittedDate:     time.Now(),
		ApprovedDate:      time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Overwrite values with those from assertions
	mergeModels(&shipmentLineItem, assertions.ShipmentLineItem)

	mustCreate(db, &shipmentLineItem)

	return shipmentLineItem
}

// MakeDefaultShipmentLineItem makes a shipment line item with default values
func MakeDefaultShipmentLineItem(db *pop.Connection) models.ShipmentLineItem {
	return MakeShipmentLineItem(db, Assertions{})
}

// MakeCompleteShipmentLineItem makes a shipmentLineItem with all dependencies that "just works"
func MakeCompleteShipmentLineItem(db *pop.Connection, assertions Assertions) models.ShipmentLineItem {
	// First we need a shipment that has proper zip3, serviceArea, etc. set up
	shipment := assertions.ShipmentLineItem.Shipment
	if isZeroUUID(shipment.ID) {
		var err error
		shipment, err = MakeShipmentForPricing(db, assertions)
		if err != nil {
			log.Panic(err)
		}
		assertions.ShipmentLineItem.Shipment = shipment
		assertions.ShipmentLineItem.ShipmentID = shipment.ID
	}

	// Then we need a 400ng item
	tariff400ngItem := assertions.ShipmentLineItem.Tariff400ngItem
	if isZeroUUID(tariff400ngItem.ID) {
		tariff400ngItem = MakeTariff400ngItem(db, assertions)

		assertions.ShipmentLineItem.Tariff400ngItem = tariff400ngItem
		assertions.ShipmentLineItem.Tariff400ngItemID = tariff400ngItem.ID
	}

	// And lastly we need a valid rate for the item code
	rateAssertions := assertions
	rateAssertions.Tariff400ngItemRate.Code = tariff400ngItem.Code
	MakeTariff400ngItemRate(db, rateAssertions)

	return MakeShipmentLineItem(db, assertions)
}
