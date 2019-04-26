package storageintransit

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type patchStorageInTransit struct {
	db *pop.Connection
}

func patchStorageInTransitWithPayload(storageInTransit *models.StorageInTransit, payload *apimessages.StorageInTransit) {
	if *payload.Location == "ORIGIN" {
		storageInTransit.Location = models.StorageInTransitLocationORIGIN
	} else {
		storageInTransit.Location = models.StorageInTransitLocationDESTINATION
	}

	if payload.EstimatedStartDate != nil {
		storageInTransit.EstimatedStartDate = *(*time.Time)(payload.EstimatedStartDate)
	}

	storageInTransit.Notes = handlers.FmtStringPtrNonEmpty(payload.Notes)

	if payload.WarehouseID != nil {
		storageInTransit.WarehouseID = *payload.WarehouseID
	}

	if payload.WarehouseName != nil {
		storageInTransit.WarehouseName = *payload.WarehouseName
	}

	if payload.WarehouseAddress != nil {
		updateAddressWithPayload(&storageInTransit.WarehouseAddress, payload.WarehouseAddress)
	}

	storageInTransit.WarehousePhone = handlers.FmtStringPtrNonEmpty(payload.WarehousePhone)
	storageInTransit.WarehouseEmail = handlers.FmtStringPtrNonEmpty(payload.WarehouseEmail)
}

// PatchStorageInTransit edits an existing storage in transit and returns the updated object.
func (p *patchStorageInTransit) PatchStorageInTransit(payload apimessages.StorageInTransit, shipmentID uuid.UUID, storageInTransitID uuid.UUID, session *auth.Session) (*models.StorageInTransit, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	// Both TSPs and Office users can do this. TSPs can edit based on whether or not its their shipment.
	isAuthorized, err := authorizeStorageInTransitHTTPRequest(p.db, session, shipmentID, true)

	if err != nil {
		return nil, returnVerrs, err
	}

	if !isAuthorized {
		return nil, returnVerrs, models.ErrFetchForbidden
	}

	storageInTransit, err := models.FetchStorageInTransitByID(p.db, storageInTransitID)

	if err != nil {
		return nil, returnVerrs, err
	}

	patchStorageInTransitWithPayload(storageInTransit, &payload)

	verrs, err := models.SaveStorageInTransitAndAddress(p.db, storageInTransit)
	if err != nil || verrs.HasAny() {
		returnVerrs.Append(verrs)
		return nil, returnVerrs, err
	}

	return storageInTransit, returnVerrs, nil
}

// NewStorageInTransitPatcher is the public constructor for a `NewStorageInTransitPatcher`
// using Pop
func NewStorageInTransitPatcher(db *pop.Connection) services.StorageInTransitPatcher {
	return &patchStorageInTransit{db}
}
