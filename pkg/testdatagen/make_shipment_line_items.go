package testdatagen

import (
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
	var rate unit.Cents
	rate = 2354
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
