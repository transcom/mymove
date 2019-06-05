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
	DB            *pop.Connection
	Engine        *rateengine.RateEngine
	Planner       route.Planner
	PriceShipment PriceShipment
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

		verrs, transactionError = c.PriceShipment.Call(shipment, ShipmentPriceNEW)
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
	for _, sit := range storageInTransits {
		// only deliver DESTINATION Sits that are IN_SIT
		if sit.Status == models.StorageInTransitStatusINSIT &&
			sit.Location == models.StorageInTransitLocationDESTINATION {
			var modifiedSit *models.StorageInTransit
			modifiedSit, err = sit.Deliver(db, deliveryDate, sit, tspID)
			if err != nil {
				return nil, verrs, err
			}
			sitsToReturn = append(sitsToReturn, *modifiedSit)
		} else {
			sitsToReturn = append(sitsToReturn, sit)
		}
	}
	verrs, err = db.ValidateAndSave(&sitsToReturn)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}

	return sitsToReturn, verrs, err
}
