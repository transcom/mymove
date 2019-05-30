package shipment

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
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

	storageInTransits, err := models.FetchStorageInTransitsOnShipment(c.DB, shipment.ID)
	if err != nil {
		return verrs, err
	}
	err = c.DB.Transaction(func(db *pop.Connection) error {
		transactionError := shipment.Deliver(deliveryDate)
		if transactionError != nil {
			return transactionError
		}

		verrs, transactionError = PriceShipment{DB: db, Engine: c.Engine, Planner: c.Planner}.Call(shipment, ShipmentPriceNEW)
		if transactionError != nil || verrs.HasAny() {
			return transactionError
		}

		_, verrs, transactionError = deliverStorageInTransits(db, storageInTransits, deliveryDate, tspID)
		if transactionError != nil || verrs.HasAny() {
			return transactionError
		}

		return nil
	})

	return verrs, err
}

// DeliverStorageInTransits delivers multiple SITS
func deliverStorageInTransits(db *pop.Connection, storageInTransits []models.StorageInTransit, deliveryDate time.Time, tspID uuid.UUID) (sitsToReturn []models.StorageInTransit, verrs *validate.Errors, err error) {
	// TODO: it looks like from the wireframes for the delivery status change form that this will also need to edit
	//  delivery address(es) and the actual delivery date.
	returnVerrs := validate.NewErrors()
	for _, sit := range storageInTransits {

		// only deliver DESTINATION Sits that are IN_SIT
		if sit.Status == models.StorageInTransitStatusINSIT &&
			sit.Location == models.StorageInTransitLocationDESTINATION {
			modifiedSit, verrs, err := deliverStorageInTransit(db, deliveryDate, sit, tspID)
			if verrs.HasAny() || err != nil {
				returnVerrs.Append(verrs)
				return nil, verrs, err
			}
			sitsToReturn = append(sitsToReturn, *modifiedSit)
		} else {
			sitsToReturn = append(sitsToReturn, sit)
		}
	}
	return sitsToReturn, returnVerrs, err
}

func deliverStorageInTransit(db *pop.Connection, deliveryDate time.Time, storageInTransit models.StorageInTransit, tspID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	// Make sure we're not trying to set delivered for something that isn't both IN SIT and a DESTINATION SIT
	if !(storageInTransit.Status == models.StorageInTransitStatusINSIT &&
		storageInTransit.Location == models.StorageInTransitLocationDESTINATION) {
		return &storageInTransit, returnVerrs, models.ErrWriteConflict
	}

	storageInTransit.Status = models.StorageInTransitStatusDELIVERED
	storageInTransit.OutDate = &deliveryDate
	if verrs, err := db.ValidateAndSave(&storageInTransit); verrs.HasAny() || err != nil {
		returnVerrs.Append(verrs)
		return nil, returnVerrs, err
	}

	return &storageInTransit, returnVerrs, nil
}
