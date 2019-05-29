package shipment

import (
	"errors"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	//sitservice "github.com/transcom/mymove/pkg/services/storage_in_transit"
)

// DeliverAndPriceShipment is a service object to deliver and price a Shipment
type DeliverAndPriceShipment struct {
	DB      *pop.Connection
	Engine  *rateengine.RateEngine
	Planner route.Planner
}

// Call delivers a Shipment and prices associated line items
func (c DeliverAndPriceShipment) Call(deliveryDate time.Time, shipment *models.Shipment, tspID uuid.UUID) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	var err error
	c.DB.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("rollback")

		err = shipment.Deliver(deliveryDate)
		if err != nil {
			return err
		}

		verrs, err = PriceShipment{DB: db, Engine: c.Engine, Planner: c.Planner}.Call(shipment, ShipmentPriceNEW)
		if err != nil || verrs.HasAny() {
			return transactionError
		}

		//sitDeliverer := sitservice.NewStorageInTransitsDeliverer(c.DB)
		//_, verrs, err = sitDeliverer.DeliverStorageInTransits(shipment.ID, tspID)
		//if err != nil || verrs.HasAny() {
		//	return transactionError
		//}

		c.deliverStorageInTransits(shipment.ID, tspID)
		if err != nil || verrs.HasAny() {
			return transactionError
		}

		return nil
	})

	return verrs, err
}

// DeliverStorageInTransits delivers multiple SITS
func (c *DeliverAndPriceShipment) deliverStorageInTransits(shipmentID uuid.UUID, tspID uuid.UUID) ([]models.StorageInTransit, *validate.Errors, error) {
	// TODO: it looks like from the wireframes for the delivery status change form that this will also need to edit
	//  delivery address(es) and the actual delivery date.
	verrs := validate.NewErrors()

	storageInTransits, err := models.FetchStorageInTransitsOnShipment(c.DB, shipmentID)
	if err != nil {
		return nil, verrs, err
	}
	sitsToReturn := []models.StorageInTransit{}
	for _, sit := range storageInTransits {
		// only deliver DESTINATION Sits that are IN_SIT
		if sit.Status == models.StorageInTransitStatusINSIT &&
			sit.Location == models.StorageInTransitLocationDESTINATION {
			modifiedSit, verrs, err := c.deliverStorageInTransit(shipmentID, sit.ID, tspID)
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

func (c *DeliverAndPriceShipment) deliverStorageInTransit(shipmentID uuid.UUID, storageInTransitID uuid.UUID, tspID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	storageInTransit, err := models.FetchStorageInTransitByID(c.DB, storageInTransitID)
	if err != nil {
		return nil, returnVerrs, err
	}
	// Verify that the shipment we're getting matches what's in the storage in transit
	if shipmentID != storageInTransit.ShipmentID {
		return nil, returnVerrs, models.ErrFetchForbidden
	}

	shipment, err := models.FetchShipmentByTSP(c.DB, tspID, shipmentID)
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
	if verrs, err := c.DB.ValidateAndSave(storageInTransit); verrs.HasAny() || err != nil {
		returnVerrs.Append(verrs)
		return nil, returnVerrs, err
	}

	return storageInTransit, returnVerrs, nil
}
