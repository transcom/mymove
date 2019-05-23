package storageintransit

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type deliverStorageInTransit struct {
	db *pop.Connection
}

// DeliverStorageInTransits delivers multiple SITS
func (d *deliverStorageInTransit) DeliverStorageInTransits(shipmentID uuid.UUID, tspID uuid.UUID) ([]models.StorageInTransit, *validate.Errors, error) {
	// TODO: it looks like from the wireframes for the delivery status change form that this will also need to edit
	//  delivery address(es) and the actual delivery date.
	verrs := validate.NewErrors()

	storageInTransits, err := models.FetchStorageInTransitsOnShipment(d.db, shipmentID)
	if err != nil {
		return nil, verrs, err
	}
	sitsToReturn := []models.StorageInTransit{}
	for _, sit := range storageInTransits {
		// only deliver DESTINATION Sits that are IN_SIT
		if sit.Status == models.StorageInTransitStatusINSIT &&
			sit.Location == models.StorageInTransitLocationDESTINATION {
			modifiedSit, verrs, err := d.deliverStorageInTransit(shipmentID, sit.ID, tspID)
			if verrs.HasAny() || err != nil {
				verrs.Append(verrs)
				return nil, verrs, err
			}
			sitsToReturn = append(sitsToReturn, *modifiedSit)
		} else {
			sitsToReturn = append(sitsToReturn, sit)
		}
	}
	return sitsToReturn, verrs, err
}

// DeliverStorageInTransit sets the status of a Storage In Transit to delivered and returns the updated object.
func (d *deliverStorageInTransit) deliverStorageInTransit(shipmentID uuid.UUID, storageInTransitID uuid.UUID, tspID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	storageInTransit, err := models.FetchStorageInTransitByID(d.db, storageInTransitID)
	if err != nil {
		return nil, returnVerrs, err
	}
	// Verify that the shipment we're getting matches what's in the storage in transit
	if shipmentID != storageInTransit.ShipmentID {
		return nil, returnVerrs, models.ErrFetchForbidden
	}

	shipment, err := models.FetchShipmentByTSP(d.db, tspID, shipmentID)
	if err != nil {
		return storageInTransit, returnVerrs, err
	}

	// Make sure we're not trying to set delivered for something that isn't both IN SIT and a DESTINATION SIT
	if !(storageInTransit.Status == models.StorageInTransitStatusINSIT &&
		storageInTransit.Location == models.StorageInTransitLocationDESTINATION) {
		return storageInTransit, returnVerrs, models.ErrWriteConflict
	}

	storageInTransit.Status = models.StorageInTransitStatusDELIVERED
	storageInTransit.OutDate = shipment.ActualDeliveryDate
	if verrs, err := d.db.ValidateAndSave(storageInTransit); verrs.HasAny() || err != nil {
		returnVerrs.Append(verrs)
		return nil, returnVerrs, err
	}

	return storageInTransit, returnVerrs, nil
}

// NewStorageInTransitsDeliverer is the public constructor for a `NewStorageInTransitInDeliverer`
// using Pop
func NewStorageInTransitsDeliverer(db *pop.Connection) services.StorageInTransitDeliverer {
	return &deliverStorageInTransit{db}
}
