package payloads

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func AddressModel(address *primemessages.Address) *models.Address {
	if address == nil {
		return nil
	}
	return &models.Address{
		ID:             uuid.FromStringOrNil(address.ID.String()),
		StreetAddress1: *address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           *address.City,
		State:          *address.State,
		PostalCode:     *address.PostalCode,
		Country:        address.Country,
	}
}

func MTOShipmentModel(mtoShipment *primemessages.MTOShipment) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	requestedPickupDate := time.Time(mtoShipment.RequestedPickupDate)
	scheduledPickupDate := time.Time(mtoShipment.ScheduledPickupDate)
	firstAvailableDeliveryDate := time.Time(mtoShipment.FirstAvailableDeliveryDate)

	model := &models.MTOShipment{
		ID:                         uuid.FromStringOrNil(mtoShipment.ID.String()),
		ShipmentType:               models.MTOShipmentType(mtoShipment.ShipmentType),
		RequestedPickupDate:        &requestedPickupDate,
		ScheduledPickupDate:        &scheduledPickupDate,
		FirstAvailableDeliveryDate: &firstAvailableDeliveryDate,
		PickupAddress:              *AddressModel(mtoShipment.PickupAddress),
		PickupAddressID:            uuid.FromStringOrNil(mtoShipment.PickupAddress.ID.String()),
		DestinationAddress:         *AddressModel(mtoShipment.DestinationAddress),
		DestinationAddressID:       uuid.FromStringOrNil(mtoShipment.DestinationAddress.ID.String()),
	}

	if mtoShipment.PrimeActualWeight > 0 {
		actualWeight := unit.Pound(mtoShipment.PrimeActualWeight)
		model.PrimeActualWeight = &actualWeight
	}

	if mtoShipment.PrimeEstimatedWeight > 0 {
		estimatedWeight := unit.Pound(mtoShipment.PrimeEstimatedWeight)
		model.PrimeEstimatedWeight = &estimatedWeight
	}

	if mtoShipment.SecondaryPickupAddress != nil {
		model.SecondaryPickupAddress = AddressModel(mtoShipment.SecondaryPickupAddress)
		secondaryPickupAddressID := uuid.FromStringOrNil(mtoShipment.SecondaryPickupAddress.ID.String())
		model.SecondaryPickupAddressID = &secondaryPickupAddressID
	}

	if mtoShipment.SecondaryDeliveryAddress != nil {
		model.SecondaryDeliveryAddress = AddressModel(mtoShipment.SecondaryDeliveryAddress)
		secondaryDeliveryAddressID := uuid.FromStringOrNil(mtoShipment.SecondaryDeliveryAddress.ID.String())
		model.SecondaryDeliveryAddressID = &secondaryDeliveryAddressID
	}

	return model
}
