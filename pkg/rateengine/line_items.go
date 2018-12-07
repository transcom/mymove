package rateengine

import (
	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
	"time"
)

// CreateBaseShipmentLineItems will create and return the models for the base shipment line items that every
// shipment should contain
func CreateBaseShipmentLineItems(db *pop.Connection, costByShipment CostByShipment) ([]models.ShipmentLineItem, error) {
	shipment := costByShipment.Shipment
	cost := costByShipment.Cost

	var lineItems []models.ShipmentLineItem

	bqNetWeight := unit.BaseQuantityFromInt(shipment.NetWeight.Int())
	now := time.Now()

	// Linehaul charges ("LHS")
	linehaulItem, err := models.FetchTariff400ngItemByCode(db, "LHS")
	if err != nil {
		return nil, err
	}

	lhAppliedRate := unit.Millicents(0) // total linehaul does not have a rate, and value should be 0
	linehaul := models.ShipmentLineItem{
		ShipmentID:        shipment.ID,
		Tariff400ngItemID: linehaulItem.ID,
		Tariff400ngItem:   linehaulItem,
		Location:          models.ShipmentLineItemLocationNEITHER,
		Quantity1:         bqNetWeight,
		Quantity2:         unit.BaseQuantityFromInt(cost.LinehaulCostComputation.Mileage),
		Status:            models.ShipmentLineItemStatusSUBMITTED,
		AmountCents:       &cost.LinehaulCostComputation.LinehaulChargeTotal,
		AppliedRate:       &lhAppliedRate,
		SubmittedDate:     now,
	}
	lineItems = append(lineItems, linehaul)

	// Origin service fee ("135A")
	originServiceFeeItem, err := models.FetchTariff400ngItemByCode(db, "135A")
	if err != nil {
		return nil, err
	}
	originService := models.ShipmentLineItem{
		ShipmentID:        shipment.ID,
		Tariff400ngItemID: originServiceFeeItem.ID,
		Tariff400ngItem:   originServiceFeeItem,
		Location:          models.ShipmentLineItemLocationORIGIN,
		Quantity1:         bqNetWeight,
		Status:            models.ShipmentLineItemStatusSUBMITTED,
		AmountCents:       &cost.NonLinehaulCostComputation.OriginService.Fee,
		AppliedRate:       &cost.NonLinehaulCostComputation.OriginService.Rate,
		SubmittedDate:     now,
	}
	lineItems = append(lineItems, originService)

	// Destination service fee ("135B")
	destinationServiceFeeItem, err := models.FetchTariff400ngItemByCode(db, "135B")
	if err != nil {
		return nil, err
	}
	destinationService := models.ShipmentLineItem{
		ShipmentID:        shipment.ID,
		Tariff400ngItemID: destinationServiceFeeItem.ID,
		Tariff400ngItem:   destinationServiceFeeItem,
		Location:          models.ShipmentLineItemLocationDESTINATION,
		Quantity1:         bqNetWeight,
		Status:            models.ShipmentLineItemStatusSUBMITTED,
		AmountCents:       &cost.NonLinehaulCostComputation.DestinationService.Fee,
		AppliedRate:       &cost.NonLinehaulCostComputation.DestinationService.Rate,
		SubmittedDate:     now,
	}
	lineItems = append(lineItems, destinationService)

	// Pack fee ("105A")
	fullPackItem, err := models.FetchTariff400ngItemByCode(db, "105A")
	if err != nil {
		return nil, err
	}
	packFee := cost.NonLinehaulCostComputation.Pack.Fee
	packRate := cost.NonLinehaulCostComputation.Pack.Rate
	fullPack := models.ShipmentLineItem{
		ShipmentID:        shipment.ID,
		Tariff400ngItemID: fullPackItem.ID,
		Tariff400ngItem:   fullPackItem,
		Location:          models.ShipmentLineItemLocationORIGIN,
		Quantity1:         bqNetWeight,
		Status:            models.ShipmentLineItemStatusSUBMITTED,
		AmountCents:       &packFee,
		AppliedRate:       &packRate,
		SubmittedDate:     now,
	}
	lineItems = append(lineItems, fullPack)

	// Unpack fee ("105C")
	fullUnpackItem, err := models.FetchTariff400ngItemByCode(db, "105C")
	if err != nil {
		return nil, err
	}
	unpackFee := cost.NonLinehaulCostComputation.Unpack.Fee
	unpackRate := cost.NonLinehaulCostComputation.Unpack.Rate
	fullUnpack := models.ShipmentLineItem{
		ShipmentID:        shipment.ID,
		Tariff400ngItemID: fullUnpackItem.ID,
		Tariff400ngItem:   fullUnpackItem,
		Location:          models.ShipmentLineItemLocationDESTINATION,
		Quantity1:         bqNetWeight,
		Status:            models.ShipmentLineItemStatusSUBMITTED,
		AmountCents:       &unpackFee,
		AppliedRate:       &unpackRate,
		SubmittedDate:     now,
	}
	lineItems = append(lineItems, fullUnpack)

	return lineItems, nil
}
