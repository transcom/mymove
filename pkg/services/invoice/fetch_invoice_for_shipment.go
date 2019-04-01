package invoice

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// FetchShipmentForInvoice is a service object for fetching a shipment with the fields required for an invoice
// This struct should contain dependencies
type FetchShipmentForInvoice struct {
	DB *pop.Connection
}

// Call queries a shipment for a given ID along with required associations
// Conditions for adding line items are:
// - must be approved or not require preapproval
// - must NOT have an existing invoice association (ie. has been invoiced already)
// - must be associated with the passed shipment ID
func (f FetchShipmentForInvoice) Call(shipmentID uuid.UUID) (models.Shipment, error) {
	var shipment models.Shipment
	err := f.DB.
		Eager(
			"PickupAddress",
			"Move.Orders.NewDutyStation.Address",
			"Move.Orders.NewDutyStation.TransportationOffice",
			"ServiceMember.DutyStation.TransportationOffice",
			"ShipmentOffers.TransportationServiceProviderPerformance.TransportationServiceProvider",
			"ShipmentOffers.TransportationServiceProviderPerformance",
		).
		Find(&shipment, shipmentID)
	if err != nil {
		return shipment, err
	}

	var lineItems models.ShipmentLineItems
	err = f.DB.Q().
		Eager("Tariff400ngItem").
		LeftJoin("tariff400ng_items as ti", "shipment_line_items.tariff400ng_item_id = ti.id").
		Where("(shipment_line_items.status=? OR ti.requires_pre_approval = false)",
			models.ShipmentLineItemStatusAPPROVED).
		Where("shipment_line_items.invoice_id IS NULL").
		Where("shipment_line_items.shipment_id=?", shipmentID).
		All(&lineItems)
	lineItemsFiltered := (filter35AItems(lineItems))
	shipment.ShipmentLineItems = lineItemsFiltered
	return shipment, err
}

// filter35AItems: 35A items are invoiced if they have an `actual_amount_cents` value
func filter35AItems(lineItems models.ShipmentLineItems) models.ShipmentLineItems {
	var lineItemsFiltered models.ShipmentLineItems
	for _, li := range lineItems {
		if li.Tariff400ngItem.Code == "35A" && li.ActualAmountCents == nil {
			continue
		} else {
			lineItemsFiltered = append(lineItemsFiltered, li)
		}
	}
	return lineItemsFiltered
}
