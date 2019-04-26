package storageintransit

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type deliverStorageInTransit struct {
	db *pop.Connection
}

// DeliverStorageInTransit sets the status of a Storage In Transit to delivered and returns the updated object.
func (d *deliverStorageInTransit) DeliverStorageInTransit(shipmentID uuid.UUID, session *auth.Session, storageInTransitID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	// Only TSP users are authorized to do this.
	isUserAuthorized, err := authorizeStorageInTransitHTTPRequest(d.db, session, shipmentID, false)

	if err != nil || !isUserAuthorized {
		return nil, returnVerrs, err
	}

	storageInTransit, err := models.FetchStorageInTransitByID(d.db, storageInTransitID)

	if err != nil {
		return nil, returnVerrs, err
	}

	// Verify that the shipment we're getting matches what's in the storage in transit
	if shipmentID != storageInTransit.ShipmentID {
		return nil, returnVerrs, models.ErrFetchForbidden
	}

	// Make sure we're not trying to set delivered for something that isn't either IN SIT or RELEASED
	if !(storageInTransit.Status == models.StorageInTransitStatusINSIT) &&
		!(storageInTransit.Status == models.StorageInTransitStatusRELEASED) {
		return nil, returnVerrs, models.ErrWriteConflict
	}

	storageInTransit.Status = models.StorageInTransitStatusDELIVERED

	if verrs, err := d.db.ValidateAndSave(storageInTransit); verrs.HasAny() || err != nil {
		returnVerrs.Append(verrs)
		return nil, returnVerrs, err
	}

	return storageInTransit, returnVerrs, nil
}

// NewStorageInTransitInDeliverer is the public constructor for a `NewStorageInTransitInDeliverer`
// using Pop
func NewStorageInTransitInDeliverer(db *pop.Connection) services.StorageInTransitDeliverer {
	return &deliverStorageInTransit{db}
}
