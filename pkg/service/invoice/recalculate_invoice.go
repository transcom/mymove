package invoice

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	shipmentop "github.com/transcom/mymove/pkg/service/shipment"
	"go.uber.org/zap"
)

// TODO: Create model tariff400ng_recalcuate
// TODO: This will have the start and end dates for when to recalculate a shipment's invoice

// Tariff400ngRecalculate struct
type Tariff400ngRecalculate struct {
	ShipmentUpdatedBefore time.Time
	ShipmentCreatedAfter  time.Time
}

var recalculateInfo = Tariff400ngRecalculate{
	ShipmentUpdatedBefore: time.Date(2019, time.January, 01, 00, 00, 00, 0, time.UTC),
	ShipmentCreatedAfter:  time.Date(1970, time.January, 01, 00, 00, 00, 0, time.UTC),
}

// TODO: Create model tariff400ng_recalculate_log
// TODO: When a shipment is recalculated then store the shipmentID, old_price, new_price, date_recalculated

// Tariff400ngRecalculateLog struct
type Tariff400ngRecalculateLog struct {
	DateRecalculated time.Time
	ShipmentID       uuid.UUID
	Code             string
	BeforePrice      float64
	AfterPrice       float64
}

// RecalculateInvoice is a service object to recalculate a Shipment's Invoice
type RecalculateInvoice struct {
	DB     *pop.Connection
	Logger *zap.Logger
}

/*
	Recalculate a shipment's invoice is temporary functionality that will be used when it has been
    determined there is is some shipment that requires recalcuation. A shipment does not contain an invoice
    (or line items) until it has reached the DELIVERED state.

    Some of the reasons a recalculation can happen are:
      - missing required pre-approved line items
      - line items that have been priced incorrectly

    The database table tariff400ng_recalculate contains the date range for shipments that need to be assessed.

    The database table tariff400ng_recalculate_log contains a record entry for each shipment that was recalculated.

    Currently to recalculate a shipment we are looking for the following:
      - Shipment is in DELIVERED or COMPLETED state
      - Shipment's line item has been updated before the date tariff400ng_recalcuate.updated_before
      - If there is an approved accessorial, this line item must be preserved to maintain the approved status

    The API that will call this method is GET /shipments/{shipmentId}/accessorials
*/

// Call recalculates a Shipment's Invoice
func (r RecalculateInvoice) Call(shipment *models.Shipment,
	lineItems models.ShipmentLineItems,
	planner route.Planner) (bool, error) {

	// If the Shipment is in the DELIVERED or COMPLETED state continue
	shipmentStatus := shipment.Status
	if shipmentStatus != models.ShipmentStatusDELIVERED && shipmentStatus != models.ShipmentStatusCOMPLETED {
		return false, nil
	}

	// If the Shipment was created before "ShipmentUpdatedBefore" date then continue
	if !r.updatedInDateRange(shipment.CreatedAt) {
		return false, nil
	}

	// If there are any line items that have been updated in the date range continue
	if !r.shipmentLineItemsUpdatedInDateRange(lineItems) {
		return false, nil
	}

	// Re-calculate the Shipment!
	engine := rateengine.NewRateEngine(r.DB, r.Logger, planner)
	verrs, err := shipmentop.RecalculateShipment{DB: r.DB, Engine: engine}.Call(shipment)
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

func (r RecalculateInvoice) shipmentLineItemsUpdatedInDateRange(lineItems models.ShipmentLineItems) bool {
	for _, item := range lineItems {
		if r.updatedInDateRange(item.UpdatedAt) {
			return true
		}
	}
	return false
}

func (r RecalculateInvoice) updatedInDateRange(update time.Time) bool {
	if update.After(recalculateInfo.ShipmentCreatedAfter) && update.Before(recalculateInfo.ShipmentUpdatedBefore) {
		return true
	}
	return false
}

func (r RecalculateInvoice) getRecalculateDateRange() /*updated_after, updated_before*/ {

}
