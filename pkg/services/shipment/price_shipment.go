package shipment

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	storageintransit "github.com/transcom/mymove/pkg/services/storage_in_transit"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

const (
	// ShipmentPriceNEW captures enum value "NEW"
	ShipmentPriceNEW services.PricingType = "NEW"
	// ShipmentPriceRECALCULATE captures enum value "RECALCULATE"
	ShipmentPriceRECALCULATE services.PricingType = "RECALCULATE"
)

// PriceShipment is a service object to price a Shipment
type shipmentPricer struct {
	db      *pop.Connection
	engine  *rateengine.RateEngine
	planner route.Planner
}

// Call prices a Shipment
func (c *shipmentPricer) PriceShipment(shipment *models.Shipment, price services.PricingType) (*validate.Errors, error) {
	origin := shipment.PickupAddress
	if origin == nil || origin.ID == uuid.Nil {
		return validate.NewErrors(), errors.New("PickupAddress not provided")
	}

	destination := shipment.Move.Orders.NewDutyStation.Address
	if shipment.DestinationAddressOnAcceptance != nil {
		destination = *shipment.DestinationAddressOnAcceptance
	}

	if destination.ID == uuid.Nil {
		return validate.NewErrors(), errors.New("Destination address not provided")
	}

	useZipOnlyForDistance := true
	distanceCalculation, err := models.NewDistanceCalculation(c.planner, *origin, destination, useZipOnlyForDistance)
	if err != nil {
		return validate.NewErrors(), errors.Wrap(err, "Error creating DistanceCalculation model")
	}

	// Delivering a shipment is a trigger to populate several shipment line items in the database.  First
	// calculate charges, then submit the updated shipment record and line items in a DB transaction.
	shipmentCost, err := c.engine.HandleRunOnShipment(*shipment, distanceCalculation)
	if err != nil {
		return validate.NewErrors(), err
	}

	if price == ShipmentPriceRECALCULATE {
		// Delete Base Shipment Line Items for repricing
		err = shipment.DeleteBaseShipmentLineItems(c.db)
		if err != nil {
			return validate.NewErrors(), err
		}

		// Delete Storage in Transit Shipment Line Items for repricing
		err = shipment.DeleteStorageInTransitLineItems(c.db)
		if err != nil {
			return validate.NewErrors(), err
		}

		// Save and validate Shipment after deleting Base Shipment Line Items
		// and Storage in Transit Shipment Line Items
		verrs, saveShipmentErr := models.SaveShipment(c.db, shipment)
		if verrs.HasAny() || saveShipmentErr != nil {
			saveError := errors.Wrap(saveShipmentErr, "Error saving shipment for ShipmentPriceRECALCULATE")
			return verrs, saveError
		}
	}

	// Create Base Shipment Line Items for Shipment
	baseLineItems, err := rateengine.CreateBaseShipmentLineItems(c.db, shipmentCost)
	if err != nil {
		return validate.NewErrors(), err
	}

	// Create Storage in Transit (SIT) line items for Shipment
	createStorageInTransitLineItems := storageintransit.CreateStorageInTransitLineItems{
		DB:      c.db,
		Planner: c.planner,
	}
	storageInTransitLineItems, err := createStorageInTransitLineItems.CreateStorageInTransitLineItems(shipmentCost)
	if err != nil {
		return validate.NewErrors(), err
	}

	// Price existing approved accessorials that require approval (and have been approved) and price storage in transit line items
	additionalLineItems, err := c.engine.PriceAdditionalRequestsForShipment(*shipment, storageInTransitLineItems)
	if err != nil {
		return validate.NewErrors(), err
	}

	verrs, err := shipment.SaveShipmentAndPricingInfo(c.db, baseLineItems, additionalLineItems, distanceCalculation)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	if price == ShipmentPriceRECALCULATE {
		log := models.ShipmentRecalculateLog{ShipmentID: shipment.ID}
		log.SaveShipmentRecalculateLog(c.db)
	}

	return validate.NewErrors(), nil
}

func NewShipmentPricer(db *pop.Connection, engine *rateengine.RateEngine, planner route.Planner) services.ShipmentPricer {
	return &shipmentPricer{db: db, engine: engine, planner: planner}
}
