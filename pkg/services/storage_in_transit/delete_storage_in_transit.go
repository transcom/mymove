package storageintransit

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type deleteStorageInTransit struct {
	db *pop.Connection
}

// DeleteStorageInTransit deletes an existing storage in transit object and returns nil if successful, an error otherwise.
func (d *deleteStorageInTransit) DeleteStorageInTransit(shipmentID uuid.UUID, storageInTransitID uuid.UUID, session *auth.Session) (*models.StorageInTransit, error) {

	// TSPs can delete their own SIT requests
	isAuthorized, err := authorizeStorageInTransitHTTPRequest(d.db, session, shipmentID, false)

	if err != nil {
		return nil, err
	}

	if !isAuthorized {
		return nil, models.ErrFetchForbidden
	}

	storageInTransit, err := models.DeleteStorageInTransit(d.db, storageInTransitID)

	if err != nil {
		return nil, err
	}

	return storageInTransit, nil
}

// NewStorageInTransitDeleter is the public constructor for a `NewStorageInTransitDeleter`
// using Pop
func NewStorageInTransitDeleter(db *pop.Connection) services.StorageInTransitDeleter {
	return &deleteStorageInTransit{db}
}
