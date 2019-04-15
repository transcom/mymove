package storageintransit

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type storageInTransitFetcher struct {
	db *pop.Connection
}

func authorizeStorageInTransitRequest(db *pop.Connection, session *auth.Session, shipmentID uuid.UUID, allowOffice bool) (isUserAuthorized bool, err error) {
	if session.IsTspUser() {
		_, _, err := models.FetchShipmentForVerifiedTSPUser(db, session.TspUserID, shipmentID)

		if err != nil {
			return false, err
		}
		return true, nil
	} else if session.IsOfficeUser() {
		if allowOffice {
			return true, nil
		}
	} else {
		return false, models.ErrFetchForbidden
	}
	return false, models.ErrFetchForbidden
}

// FetchStorageInTransitbyID gets a Storage In Transit record by ID
// Authorizes based on session and shipment ID
func (s storageInTransitFetcher) FetchStorageInTransitByID(storageInTransitID uuid.UUID, shipmentID uuid.UUID, session *auth.Session) (*models.StorageInTransit, error) {
	isAuthorized, err := authorizeStorageInTransitRequest(s.db, session, shipmentID, true)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, models.ErrFetchForbidden
	}
	return models.FetchStorageInTransitByID(s.db, storageInTransitID)
}

// NewStorageInTransitByIDFetcher is the public constructor for a `StorageInTransitFetcher`
// using Pop
func NewStorageInTransitByIDFetcher(db *pop.Connection) services.StorageInTransitByIDFetcher {
	return storageInTransitFetcher{db}
}
