package mtoshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func checkUBShipmentForOconusAddress(shipment models.MTOShipment) error {
	if !*shipment.DestinationAddress.IsOconus && !*shipment.PickupAddress.IsOconus {
		return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "At least one address must be an OCONUS address")
	}
	return nil
}
