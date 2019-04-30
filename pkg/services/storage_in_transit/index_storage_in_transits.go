package storageintransit

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type indexStorageInTransits struct {
	db *pop.Connection
}

// IndexStorageInTransits returns a collection of Storage In Transits that are associated with a specific shipmentID
func (i *indexStorageInTransits) IndexStorageInTransits(shipmentID uuid.UUID, session *auth.Session) ([]models.StorageInTransit, error) {
	isUserAuthorized, err := authorizeStorageInTransitHTTPRequest(i.db, session, shipmentID, true)
	if err != nil {
		return nil, err
	}

	if !isUserAuthorized {
		return nil, models.ErrFetchForbidden
	}

	storageInTransits, err := models.FetchStorageInTransitsOnShipment(i.db, shipmentID)

	if err != nil {
		return nil, err
	}

	return storageInTransits, nil
}

// NewStorageInTransitIndexer is the public constructor for a `StorageInTransitIndexer`
// using Pop
func NewStorageInTransitIndexer(db *pop.Connection) services.StorageInTransitsIndexer {
	return &indexStorageInTransits{db}
}
