package shipment

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"go.uber.org/zap"
)

// ProcessRecalculateShipment is a service object to recalculate a Shipment's Line Items
type ProcessRecalculateShipment struct {
	DB     *pop.Connection
	Logger *zap.Logger
}

/*
	Recalculate a shipment's line items is temporary functionality that will be used when it has been
    determined there is some shipment that requires recalculation. A shipment does not contain line items
    until it has reached the DELIVERED state.

    Some of the reasons a recalculation can happen are:
      - missing required pre-approved line items
      - line items that have been priced incorrectly (this is detected by saying that a line item could have been
        priced erroneously if priced in this specified date range)

    The database table shipment_recalculate contains the date range for shipments that need to be assessed.

    The database table shipment_recalculate_log contains a record entry for each shipment that was recalculated.

    Currently to recalculate a shipment we are looking for the following:
      - Shipment is in DELIVERED or COMPLETED state
      - If the Shipment was created within the specified recalculation date range
      - Shipment has a line item that has been updated before the date shipment_recalculate.updated_before
      - If there is an approved accessorial, this line item must be preserved to maintain the approved status
      - If a Shipment does not have all of the Base Shipment Line Items it will be re-calculated
      - Shipment does not have any line item that has an InvoiceID

    The API that will call this method is GET /shipments/{shipmentId}/accessorials
*/

// Call recalculates a Shipment's Line Items
func (r ProcessRecalculateShipment) Call(shipment *models.Shipment, lineItems models.ShipmentLineItems, planner route.Planner) (bool, error) {

	// If there is an active recalculate date range then continue
	recalculateDates, err := models.FetchShipmentRecalculateDates(r.DB)
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
		errorString := ""
		if verrs.HasAny() {
			errorString = "verrs: " + verrs.String()
		}
		if err != nil {
			errorString = errorString + " err: " + err.Error()
		}
		recalculateError := fmt.Errorf("Error saving shipment for RecalculateShipment %s", errorString)
		// return true for update so that the caller can refresh line items and shipment
		return true, recalculateError
	}

	return true, nil
}

func (r ProcessRecalculateShipment) hasAllBaseLineItems(lineItems models.ShipmentLineItems) bool {
	err := models.VerifyBaseShipmentLineItems(lineItems)
	if err != nil {
		return false
	}
	return true
}

func (r ProcessRecalculateShipment) shipmentLineItemsUpdatedInDateRange(lineItems models.ShipmentLineItems, recalculateDates *models.ShipmentRecalculate) bool {
	for _, item := range lineItems {
		if r.updatedInDateRange(item.UpdatedAt, recalculateDates) {
			return true
		}
	}
	return false
}

func (r ProcessRecalculateShipment) updatedInDateRange(update time.Time, recalculateDates *models.ShipmentRecalculate) bool {
	if update.After(recalculateDates.ShipmentUpdatedAfter) && update.Before(recalculateDates.ShipmentUpdatedBefore) {
		return true
	}
	return false
}
