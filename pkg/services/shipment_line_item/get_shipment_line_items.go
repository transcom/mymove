package shipmentlineitem

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type getShipmentLineItems struct {
	db *pop.Connection
}

// IndexStorageInTransits returns a collection of Storage In Transits that are associated with a specific shipmentID
func (i *getShipmentLineItems) GetShipmentLineItemsByShipmentID(shipmentID uuid.UUID, session *auth.Session) ([]models.ShipmentLineItem, error) {

	if session.IsTspUser() {
		// Check that the TSP user can access the shipment
		_, _, err := models.FetchShipmentForVerifiedTSPUser(i.db, session.TspUserID, shipmentID)
		if err != nil {
			return nil, err
		}
	} else if !session.IsOfficeUser() {
		return nil, models.ErrFetchForbidden
	}

	shipmentLineItems, err := models.FetchLineItemsByShipmentID(i.db, &shipmentID)
	if err != nil {
		return nil, err
	}

	return shipmentLineItems, nil
}

// NewShipmentLineItemFetcher is the public constructor for a `ShipmentLineItemFetcher`
// using Pop
func NewShipmentLineItemFetcher(db *pop.Connection) services.ShipmentLineItemFetcher {
	return &getShipmentLineItems{db}
}
