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
	RequestShipmentDeliveryAddressUpdate(appCtx appcontext.AppContext, shipmentID uuid.UUID, newAddress models.ShipmentAddressUpdate, contractorRemarks string, eTag string) (*models.ShipmentAddressUpdate, error)
	ReviewShipmentAddressChange(appCtx appcontext.AppContext, shipmentID uuid.UUID, tooApprovalStatus models.ShipmentAddressUpdateStatus, tooRemarks string) (*models.ShipmentAddressUpdate, error)
}
