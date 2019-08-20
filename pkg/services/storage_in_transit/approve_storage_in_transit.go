package storageintransit

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type approveStorageInTransit struct {
	db *pop.Connection
}

// ApproveStorageInTransit sets the status of a Storage In Transit to approved, saves its Authorization Notes, saves its ActualDate, and returns the updated object.
func (a *approveStorageInTransit) ApproveStorageInTransit(payload apimessages.StorageInTransitApprovalPayload, shipmentID uuid.UUID, session *auth.Session, storageInTransitID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	// Only office users are authorized to do this.
	if !session.IsOfficeUser() {
		return nil, returnVerrs, models.ErrFetchForbidden
	}

	storageInTransit, err := models.FetchStorageInTransitByID(a.db, storageInTransitID)

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

	storageInTransit.Status = models.StorageInTransitStatusAPPROVED
	storageInTransit.AuthorizationNotes = payload.AuthorizationNotes
	storageInTransit.AuthorizedStartDate = (*time.Time)(&payload.AuthorizedStartDate)

	if verrs, err := a.db.ValidateAndSave(storageInTransit); verrs.HasAny() || err != nil {
		returnVerrs.Append(verrs)
		return nil, returnVerrs, err
	}

	return storageInTransit, returnVerrs, nil
}

// NewStorageInTransitApprover is the public constructor for a `StorageInTransitApprover`
// using Pop
func NewStorageInTransitApprover(db *pop.Connection) services.StorageInTransitApprover {
	return &approveStorageInTransit{db}
}
