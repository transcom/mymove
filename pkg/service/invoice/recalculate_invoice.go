package invoice

import (
	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
)

// TODO: Create model tariff400ng_recalcuate
// TODO: This will have the start and end dates for when to recalculate a shipment's invoice

// TODO: Create model tariff400ng_recalculate_log
// TODO: When a shipment is recalculated then store the shipmentID, old_price, new_price, date_recalculated

// RecalculateInvoice is a service object to recalculate a Shipment's Invoice
type RecalculateInvoice struct {
	DB     *pop.Connection
	Logger *zap.Logger
}

/*
	Recalculate a shipment's invoice is temporary functionality that will be used when it has been
    determined there is is some shipment that requires recalcuation. A shipment does not contain an invoice
    (or line items) until it has reached the delivered state.

    Some of the reasons a recalculation can happen are:
      - missing required pre-approved line items
      - line items that have been priced incorrectly

    The database table tariff400ng_recalcuate contains the date range for shipments that need to be assessed.

    The database table tariff400ng_recalculate_log contains a record entry for each shipment that was recalculated.

    Currently to recalculate a shipment we are looking for the following:
      - Shipment is in DELIVERED or COMPLETED state
      - Shipment's line item has been updated before the date tariff400ng_recalcuate.updated_before
      - If there is an approved accessorial, this line item must be preserved to maintain the approved status
*/

// Call recalculates a Shipment's Invoice
func (r RecalculateInvoice) Call() error {
	return nil
}

func (r RecalculateInvoice) getRecalculateDateRange() /*updated_after, updated_before*/ {

}
