package storageintransit

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type denyStorageInTransit struct {
	db *pop.Connection
}

// DenyStorageInTransit sets the status of a Storage In Transit to denied, saves its Authorization Notes, and returns the updated object.
func (d *denyStorageInTransit) DenyStorageInTransit(payload apimessages.StorageInTransitDenialPayload, shipmentID uuid.UUID, session *auth.Session, storageInTransitID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	// Only office users are authorized to do this.
	if !session.IsOfficeUser() {
		return nil, returnVerrs, models.ErrFetchForbidden
	}

	storageInTransit, err := models.FetchStorageInTransitByID(d.db, storageInTransitID)

	if err != nil {
		return nil, returnVerrs, models.ErrFetchNotFound
	}

	// Verify that the shipment we're getting matches what's in the storage in transit
	if shipmentID != storageInTransit.ShipmentID {
		return nil, returnVerrs, models.ErrFetchForbidden
	}

	if storageInTransit.Status == models.StorageInTransitStatusDELIVERED {
		return nil, returnVerrs, models.ErrWriteConflict
	}

	storageInTransit.Status = models.StorageInTransitStatusDENIED
	storageInTransit.AuthorizationNotes = &payload.AuthorizationNotes

	if verrs, err := d.db.ValidateAndSave(storageInTransit); verrs.HasAny() || err != nil {
		returnVerrs.Append(verrs)
		return nil, returnVerrs, err
	}

	return storageInTransit, returnVerrs, nil

}

// NewStorageInTransitDenier is the public constructor for a `StorageInTransitDenier`
// using Pop
func NewStorageInTransitDenier(db *pop.Connection) services.StorageInTransitDenier {
	return &denyStorageInTransit{db}
}
