package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ShipmentAddressUpdateRequester Interface for the service object that creates an approved SIT Address Update
//
//go:generate mockery --name ShipmentAddressUpdateRequester
type ShipmentAddressUpdateRequester interface {
	RequestShipmentDeliveryAddressUpdate(appCtx appcontext.AppContext, shipmentID uuid.UUID, newAddress models.Address, contractorRemarks string) (*models.ShipmentAddressUpdate, error)
}
