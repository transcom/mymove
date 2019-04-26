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

type placeIntoSITStorageInTransit struct {
	db *pop.Connection
}

// PlaceIntoSITStorageInTransit sets the status of a Storage In Transit to IN SIT, saves its ActualStartDate, and returns the updated object.
func (p *placeIntoSITStorageInTransit) PlaceIntoSITStorageInTransit(payload apimessages.StorageInTransitInSitPayload, shipmentID uuid.UUID, session *auth.Session, storageInTransitID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	// Only TSP users are authorized to do this.
	isUserAuthorized, err := authorizeStorageInTransitHTTPRequest(p.db, session, shipmentID, false)

	if err != nil {
		return nil, returnVerrs, err
	}

	if !isUserAuthorized {
		return nil, returnVerrs, err
	}

	storageInTransit, err := models.FetchStorageInTransitByID(p.db, storageInTransitID)

	if err != nil {
		return nil, returnVerrs, err
	}

	// Verify that the shipment we're getting matches what's in the storage in transit
	if shipmentID != storageInTransit.ShipmentID {
		return nil, returnVerrs, models.ErrFetchForbidden
	}

	if !(storageInTransit.Status == models.StorageInTransitStatusAPPROVED) {
		return nil, returnVerrs, models.ErrWriteConflict
	}

	payloadActualStartDate := (time.Time)(payload.ActualStartDate)

	storageInTransit.Status = models.StorageInTransitStatusINSIT
	storageInTransit.ActualStartDate = &payloadActualStartDate

	if verrs, err := p.db.ValidateAndSave(storageInTransit); verrs.HasAny() || err != nil {
		returnVerrs.Append(verrs)
		return nil, returnVerrs, err
	}

	return storageInTransit, returnVerrs, nil

}

// NewStorageInTransitInSITPlacer is the public constructor for a `NewStorageInTransitInSITPlacer`
// using Pop
func NewStorageInTransitInSITPlacer(db *pop.Connection) services.StorageInTransitInSITPlacer {
	return &placeIntoSITStorageInTransit{db}
}
