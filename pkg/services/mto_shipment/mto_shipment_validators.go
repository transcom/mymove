package mtoshipment

import (
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func checkUBShipmentForOconusAddress(shipment models.MTOShipment) error {
	if !*shipment.DestinationAddress.IsOconus && !*shipment.PickupAddress.IsOconus {
		return apperror.NewUnprocessableEntityError("At least one address for a UB shipment must be OCONUS")
	}
	return nil
}
