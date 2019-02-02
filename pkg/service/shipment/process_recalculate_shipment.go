package shipment

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"go.uber.org/zap"
)

// Tariff400ngRecalculateLog struct
type Tariff400ngRecalculateLog struct {
	DateRecalculated time.Time
	ShipmentID       uuid.UUID
	Code             string
	BeforePrice      float64
	AfterPrice       float64
}

// ProcessRecalculateShipment is a service object to recalculate a Shipment's Line Items
type ProcessRecalculateShipment struct {
	DB     *pop.Connection
	Logger *zap.Logger
}

/*
	Recalculate a shipment's line items is temporary functionality that will be used when it has been
    determined there is is some shipment that requires recalcuation. A shipment does not contain line items
    until it has reached the DELIVERED state.

    Some of the reasons a recalculation can happen are:
      - missing required pre-approved line items
      - line items that have been priced incorrectly (this is detected by saying that a line item could have been
        price erroneously if priced in this specified date range)

    The database table tariff400ng_recalculate contains the date range for shipments that need to be assessed.

    The database table tariff400ng_recalculate_log contains a record entry for each shipment that was recalculated.

    Currently to recalculate a shipment we are looking for the following:
      - Shipment is in DELIVERED or COMPLETED state
      - If the Shipment was created within the specified recalculation date range
      - Shipment's line item has been updated before the date tariff400ng_recalcuate.updated_before
      - If there is an approved accessorial, this line item must be preserved to maintain the approved status
      - If a Shipment does have have all of the Base Shipment Line Items it will be re-calculated

    The API that will call this method is GET /shipments/{shipmentId}/accessorials

    TODO: Potential issue: FetchShipmentForInvoice{}.Call() is used to retrieve a Shipment and it's ShipmentLineItems
    TODO: ths fetch takes care not to return shipment line items that have an Invoice ID attached. However,
    TODO: if re-calculate is triggered on a shipment with lingering shipment line items and the Base ShipmentLineItems
    TODO: have an Invoice ID already attached, then calling rateengine.CreateBaseShipmentLineItems() fails, because
    TODO: it tries to add BaseShipmentLineItems to the Shipment and they are then duplicates. Since no Invoice's have
    TODO: been sent to date, this isn't an issue. But if we use recalculate functionality when there are
    TODO: invoices sent it will cause a problem. This problem is not unique to RecalculateShipment{}.Call() it is
    TODO: it is an issue with the underlying PriceShipment{}.Call().
    TODO:
    TODO: The problem that is unique to  RecalculateShipment{}.Call() is that it is called for shipments that are in the
    TODO: ShipmentStatusDELIVERED or ShipmentStatusCOMPLETED state. And in that state, a Shipment could potentially
    TODO: have a sent invoice and that will cause an issue. Hitting the payment button doesn't change the status of the
    TODO: HHG Shipment.
*/

// Call recalculates a Shipment's Line Items
func (r ProcessRecalculateShipment) Call(shipment *models.Shipment, lineItems models.ShipmentLineItems, planner route.Planner) (bool, error) {

	// If there is an active recalculate date range then continue
	recalculateDates, err := models.FetchTariff400ngRecalculateDates(r.DB)
	if recalculateDates == nil || err != nil {
		return false, nil
	}

	// If the Shipment is in the DELIVERED or COMPLETED state continue
	shipmentStatus := shipment.Status
	if shipmentStatus != models.ShipmentStatusDELIVERED && shipmentStatus != models.ShipmentStatusCOMPLETED {
		return false, nil
	}

	// If the Shipment was created before "ShipmentUpdatedBefore" date then continue
	if !r.updatedInDateRange(shipment.CreatedAt, recalculateDates) {
		return false, nil
	}

	// If Shipment does not have all of the base line items expected or
	// a shipment line item was updated within the recalculate update range then continue
	if r.hasAllBaseLineItems(lineItems) && !r.shipmentLineItemsUpdatedInDateRange(lineItems, recalculateDates) {
		return false, nil
	}

	// Re-calculate the Shipment!
	engine := rateengine.NewRateEngine(r.DB, r.Logger, planner)
	verrs, err := RecalculateShipment{
		DB:     r.DB,
		Logger: r.Logger,
		Engine: engine,
	}.Call(shipment)
	if verrs.HasAny() || err != nil {
		verrsString := ""
		if verrs.HasAny() {
			verrsString = "verrs: " + verrs.String()
		}
		recalculateError := errors.Wrap(err, fmt.Sprintf("Error saving shipment for RecalculateShipment %s", verrsString))
		// return true for update so that the caller can refresh line items and shipment
		return true, recalculateError
	}

	return true, nil
}

func (r ProcessRecalculateShipment) hasAllBaseLineItems(lineItems models.ShipmentLineItems) bool {
	var count int
	count = 0
	for _, item := range lineItems {
		if models.FindBaseShipmentLineItem(item.Tariff400ngItem.Code) {
			count++
		}
	}
	if count == len(models.BaseShipmentLineItems) {
		return true
	}
	return false
}

func (r ProcessRecalculateShipment) shipmentLineItemsUpdatedInDateRange(lineItems models.ShipmentLineItems, recalculateDates *models.Tariff400ngRecalculate) bool {
	for _, item := range lineItems {
		if r.updatedInDateRange(item.UpdatedAt, recalculateDates) {
			return true
		}
	}
	return false
}

func (r ProcessRecalculateShipment) updatedInDateRange(update time.Time, recalculateDates *models.Tariff400ngRecalculate) bool {
	if update.After(recalculateDates.ShipmentUpdatedAfter) && update.Before(recalculateDates.ShipmentUdpatedBefore) {
		return true
	}
	return false
}
