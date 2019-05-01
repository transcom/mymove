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

type releaseStorageInTransit struct {
	db *pop.Connection
}

// ReleaseStorageInTransit sets the status of a Storage In Transit to released, saves its released on date, and returns the updated object.
func (r *releaseStorageInTransit) ReleaseStorageInTransit(payload apimessages.StorageInTransitReleasePayload, shipmentID uuid.UUID, session *auth.Session, storageInTransitID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	// Only TSPs are authorized to do this and they should only be able to on their own shipments
	isAuthorized, err := authorizeStorageInTransitHTTPRequest(r.db, session, shipmentID, false)

	if err != nil {
		return nil, returnVerrs, err
	}

	if !isAuthorized {
		return nil, returnVerrs, models.ErrFetchForbidden
	}

	if err != nil {
		return nil, returnVerrs, models.ErrFetchForbidden
	}

	storageInTransit, err := models.FetchStorageInTransitByID(r.db, storageInTransitID)

	if err != nil {
		return nil, returnVerrs, err
	}

	// Make sure we're not releasing something that wasn't in SIT or in delivered status.
	// The latter is there so that we can 'undo' a mistaken deliver action.
	if !(storageInTransit.Status == models.StorageInTransitStatusINSIT) &&
		!(storageInTransit.Status == models.StorageInTransitStatusDELIVERED) {
		return nil, returnVerrs, models.ErrWriteConflict
	}
	storageInTransit.Status = models.StorageInTransitStatusRELEASED
	storageInTransit.OutDate = (*time.Time)(&payload.ReleasedOn)

	if verrs, err := r.db.ValidateAndSave(storageInTransit); verrs.HasAny() || err != nil {
		returnVerrs.Append(verrs)
		return nil, returnVerrs, err
	}

	return storageInTransit, returnVerrs, nil
}

// NewStorageInTransitInReleaser is the public constructor for a `NewStorageInTransitInReleaser`
// using Pop
func NewStorageInTransitInReleaser(db *pop.Connection) services.StorageInTransitReleaser {
	return &releaseStorageInTransit{db}
}
