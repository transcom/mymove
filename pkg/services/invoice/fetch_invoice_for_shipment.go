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
	lineItemsFiltered := (filterRobust35AItems(lineItems))
	shipment.ShipmentLineItems = lineItemsFiltered
	return shipment, err
}

// filterRobust35AItems: robust 35A items are invoiced only if they have an `actual_amount_cents` value or are 35A legacy item
func filterRobust35AItems(lineItems models.ShipmentLineItems) models.ShipmentLineItems {
	var lineItemsFiltered models.ShipmentLineItems
	for _, li := range lineItems {
		if is35AItemMissingActualAmount(li) {
			continue
		} else {
			lineItemsFiltered = append(lineItemsFiltered, li)
		}
	}
	return lineItemsFiltered
}

// is35AItemMissingActualAmount check if 35A item is robust but has no actual_amount
func is35AItemMissingActualAmount(item models.ShipmentLineItem) bool {
	if item.Tariff400ngItem.Code == "35A" {
		return item.Description != nil && item.Reason != nil && item.EstimateAmountCents != nil && item.ActualAmountCents == nil
	}
	return false
}
