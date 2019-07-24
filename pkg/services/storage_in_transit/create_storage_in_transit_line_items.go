package storageintransit

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

type CreateStorageInTransitLineItems struct {
	DB      *pop.Connection
	Planner route.Planner
}

func (c CreateStorageInTransitLineItems) storageInTransitDistance(storageInTransit models.StorageInTransit, shipment models.Shipment) (*models.DistanceCalculation, error) {

	var origin models.Address
	var destination models.Address

	if storageInTransit.Location == models.StorageInTransitLocationDESTINATION {
		origin = storageInTransit.WarehouseAddress
		destination = shipment.Move.Orders.NewDutyStation.Address
		if shipment.DestinationAddressOnAcceptance != nil {
			destination = *shipment.DestinationAddressOnAcceptance
		}
	} else if storageInTransit.Location == models.StorageInTransitLocationORIGIN {
		if shipment.PickupAddress != nil {
			origin = *shipment.PickupAddress
		} else {
			return nil, errors.New("StorageInTransit PickupAddress not provided")
		}
		destination = storageInTransit.WarehouseAddress
	}

	if origin.ID == uuid.Nil {
		return nil, errors.New("StorageInTransit PickupAddress not provided")
	}

	if destination.ID == uuid.Nil {
		return nil, errors.New("StorageInTransit Destination address not provided")
	}

	useZipOnlyForDistance := false
	distanceCalculation, err := models.NewDistanceCalculation(c.Planner, origin, destination, useZipOnlyForDistance)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating StorageInTransit DistanceCalculation model")
	}

	verrs, err := c.DB.ValidateAndSave(&distanceCalculation)
	if verrs.HasAny() || err != nil {
		saveError := errors.Wrapf(err, "Error saving storage in transit distance %s", verrs.Error())
		return nil, saveError
	}

	return &distanceCalculation, nil
}

func (c CreateStorageInTransitLineItems) shipmentItemLocation(location models.StorageInTransitLocation) models.ShipmentLineItemLocation {
	if location == models.StorageInTransitLocationORIGIN {
		return models.ShipmentLineItemLocationORIGIN
	}

	if location == models.StorageInTransitLocationDESTINATION {
		return models.ShipmentLineItemLocationDESTINATION
	}

	return models.ShipmentLineItemLocationNEITHER
}
func (c CreateStorageInTransitLineItems) saveLineItem(lineItem *models.ShipmentLineItem) error {
	verrs, err := c.DB.ValidateAndCreate(lineItem)

	if err != nil || verrs.HasAny() {
		responseError := errors.Wrapf(err, "Error saving storage in transit line item for shipment %s and item %s with verr %s",
			lineItem.ShipmentID, lineItem.Tariff400ngItemID, verrs.Error())
		return responseError
	}
	return nil
}

func (c CreateStorageInTransitLineItems) createShipmentLineItem(itemCode string, shipment models.Shipment, sit models.StorageInTransit, currentTime time.Time) (models.ShipmentLineItem, error) {

	tariffItem, err := models.FetchTariff400ngItemByCode(c.DB, itemCode)
	if err != nil {
		return models.ShipmentLineItem{}, errors.Wrapf(err, "Error fetching item code %s - CreateStorageInTransitLineItems()", itemCode)
	}
	lineItem := models.ShipmentLineItem{
		ShipmentID:        shipment.ID,
		Shipment:          shipment,
		Tariff400ngItemID: tariffItem.ID,
		Tariff400ngItem:   tariffItem,
		Location:          c.shipmentItemLocation(sit.Location),
		Quantity1:         unit.BaseQuantityFromInt(sit.StorageInTransitDistance.DistanceMiles),
		Status:            models.ShipmentLineItemStatusAPPROVED,
		SubmittedDate:     currentTime,
		Address:           sit.WarehouseAddress,
		AddressID:         &sit.WarehouseAddressID,
	}
	err = c.saveLineItem(&lineItem)
	if err != nil {
		return models.ShipmentLineItem{}, errors.Wrapf(err, "Error saving line item %s - CreateStorageInTransitLineItems()", itemCode)
	}

	return lineItem, nil
}

func (c CreateStorageInTransitLineItems) CreateStorageInTransitLineItems(costByShipment rateengine.CostByShipment) ([]models.ShipmentLineItem, error) {
	var lineItems []models.ShipmentLineItem
	shipment := costByShipment.Shipment
	//cost := costByShipment.Cost needed for 16B Fuelsurcharge
	now := time.Now()
	storageInTransits, err := models.FetchStorageInTransitsOnShipment(c.DB, shipment.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating StorageInTransit line items")
	}

	for _, sit := range storageInTransits {

		// Only add line items for storage in transits that have been released or delivered
		if sit.Status != models.StorageInTransitStatusDELIVERED && sit.Status != models.StorageInTransitStatusRELEASED {
			continue
		}

		// Calculate distance for storage in transit
		distanceCalculation, err := c.storageInTransitDistance(sit, shipment)
		if distanceCalculation == nil || err != nil {
			return nil, errors.Wrap(err, "Error finding the distance calculation for Storage In Transit")
		}
		sit.StorageInTransitDistance = *distanceCalculation
		sit.StorageInTransitDistanceID = &(*distanceCalculation).ID

		/****************************************************************************
		 * Add 210A, 210B, and 210C Shipment Line Items
		 * 210A-E = Additional flat rate costs based on distance to/from the SIT facility. These vary based on geographical schedules.
		 *
		 *
		 * Up to 30 miles: Item 210A --  SIT Pup/Del - 30 or Less Miles
		 * Up to 50 miles: Item 201A & 210B -- SIT Pup/Del 31 - 50 Miles
		 * Over 50 miles : Item 210C (Use the linehaul tables for computation of charges) -- SIT Pup/Del - Over 50 Miles
		 * Over 50 miles (Alaska only) : Item 210F (Use linehaul tables section 7 Intra-AK)
		 ****************************************************************************/

		const distance30Miles = 30
		const distance50Miles = 50
		if sit.StorageInTransitDistance.DistanceMiles > distance50Miles {
			additionalFlatRate, err := c.createShipmentLineItem("210C", shipment, sit, now)
			if err != nil {
				return nil, errors.Wrapf(err, "Error creating shipment line item for 210C - CreateStorageInTransitLineItems()")
			}
			lineItems = append(lineItems, additionalFlatRate)
		} else {
			if sit.StorageInTransitDistance.DistanceMiles > distance30Miles {
				additionalFlatRate, err := c.createShipmentLineItem("210B", shipment, sit, now)
				if err != nil {
					return nil, errors.Wrapf(err, "Error creating shipment line item for 210B - CreateStorageInTransitLineItems()")
				}
				lineItems = append(lineItems, additionalFlatRate)
			}
			additionalFlatRate, err := c.createShipmentLineItem("210A", shipment, sit, now)
			if err != nil {
				return nil, errors.Wrapf(err, "Error creating shipment line item for 210A - CreateStorageInTransitLineItems()")
			}
			lineItems = append(lineItems, additionalFlatRate)
		}

		/* TODO: Failing to load 16B from database
		// https://www.pivotaltracker.com/story/show/166766741
		// Fuel Surcharge (16B) - DEL to/from SIT (Deliver to and from Storage in Transit)
		fuelSurchargeItem, err := models.FetchTariff400ngItemByCode(c.DB, "16B")
		if err != nil {
			return nil, err
		}

		fsAppliedRate := &cost.LinehaulCostComputation.FuelSurcharge.Rate
		fuelSurcharge := models.ShipmentLineItem{
			ShipmentID:        shipment.ID,
			Tariff400ngItemID: fuelSurchargeItem.ID,
			Tariff400ngItem:   fuelSurchargeItem,
			Location:          models.ShipmentLineItemLocation(fuelSurchargeItem.AllowedLocation),
			Quantity1:         unit.BaseQuantityFromInt(shipment.NetWeight.Int()),
			Quantity2:         unit.BaseQuantityFromInt(sit.StorageInTransitDistance.DistanceMiles),
			Status:            models.ShipmentLineItemStatusAPPROVED,
			AmountCents:       &cost.LinehaulCostComputation.FuelSurcharge.Fee,
			AppliedRate:       fsAppliedRate,
			SubmittedDate:     now,
		}
		lineItems = append(lineItems, fuelSurcharge)
		*/

	}
	return lineItems, nil
}
