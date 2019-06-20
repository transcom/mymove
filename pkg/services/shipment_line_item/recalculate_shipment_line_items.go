package shipmentlineitem

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/invoice"
	shipmentop "github.com/transcom/mymove/pkg/services/shipment"
)

type recalculateShipmentLineItems struct {
	db     *pop.Connection
	logger Logger
}

// RecalculateShipmentLineItems returns a collection of ShipmentLineItems that have been recalculated.
func (r *recalculateShipmentLineItems) RecalculateShipmentLineItems(shipmentID uuid.UUID, session *auth.Session, route route.Planner) ([]models.ShipmentLineItem, error) {
	if session.IsTspUser() {
		// Check that the TSP user can access the shipment
		_, _, err := models.FetchShipmentForVerifiedTSPUser(r.db, session.TspUserID, shipmentID)
		if err != nil {
			return nil, err
		}
	} else if !session.IsOfficeUser() {
		return nil, models.ErrFetchForbidden
	}

	shipmentLineItems, err := models.FetchLineItemsByShipmentID(r.db, &shipmentID)

	// If there is a shipment line item with an invoice do not run the recalculate function
	// the system is currently not setup to re-price a shipment with an existing invoice
	// and currently the system does not expect to have multiple invoices per shipment
	for _, item := range shipmentLineItems {
		if item.InvoiceID != nil {
			return nil, nil
		}
	}

	// Need to fetch Shipment to get the Accepted Offer and the ShipmentLineItems
	// Only returning ShipmentLineItems that are approved and have no InvoiceID
	shipment, err := invoice.FetchShipmentForInvoice{DB: r.db}.Call(shipmentID)
	if err != nil {
		return nil, err
	}
	// Run re-calculation process
	update, err := shipmentop.ProcessRecalculateShipment{
		DB:     r.db,
		Logger: r.logger,
	}.Call(&shipment, shipmentLineItems, route)

	if err != nil {
		return nil, err
	}

	if !update {
		return nil, nil
	}

	returnShipmentLineItems, err := models.FetchLineItemsByShipmentID(r.db, &shipmentID)
	if err != nil {
		return nil, err
	}

	return returnShipmentLineItems, nil
}

// NewShipmentLineItemRecalculator is the public constructor for a `NewShipmentLineItemRecalculator`
// using Pop
func NewShipmentLineItemRecalculator(db *pop.Connection, logger Logger) services.ShipmentLineItemRecalculator {
	return &recalculateShipmentLineItems{db, logger}
}
