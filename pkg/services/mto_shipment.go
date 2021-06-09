package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

//MTOShipmentUpdater is the service object interface for UpdateMTOShipment
//go:generate mockery --name MTOShipmentUpdater
type MTOShipmentUpdater interface {
	CheckIfMTOShipmentCanBeUpdated(mtoShipment *models.MTOShipment, session *auth.Session) (bool, error)
	MTOShipmentsMTOAvailableToPrime(mtoShipmentID uuid.UUID) (bool, error)
	RetrieveMTOShipment(mtoShipmentID uuid.UUID) (*models.MTOShipment, error)
	UpdateMTOShipment(mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error)
}

//ShipmentDeleter is the service object interface for deleting a shipment
//go:generate mockery --name ShipmentDeleter
type ShipmentDeleter interface {
	DeleteShipment(shipmentID uuid.UUID) (uuid.UUID, error)
}

//ShipmentApprover is the service object interface for approving a shipment
//go:generate mockery --name ShipmentApprover
type ShipmentApprover interface {
	ApproveShipment(shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

//ShipmentDiversionApprover is the service object interface for approving a shipment diversion
//go:generate mockery --name ShipmentDiversionApprover
type ShipmentDiversionApprover interface {
	ApproveShipmentDiversion(shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

// MTOShipmentStatusUpdater is the exported interface for updating an MTO shipment status
//go:generate mockery --name MTOShipmentStatusUpdater
type MTOShipmentStatusUpdater interface {
	UpdateMTOShipmentStatus(shipmentID uuid.UUID, status models.MTOShipmentStatus, rejectionReason *string, eTag string) (*models.MTOShipment, error)
}

// MTOShipmentCreator is the exported interface for creating a payment request
//go:generate mockery --name MTOShipmentCreator
type MTOShipmentCreator interface {
	CreateMTOShipment(MTOShipment *models.MTOShipment, MTOServiceItems models.MTOServiceItems) (*models.MTOShipment, error)
}

// MTOShipmentAddressUpdater is the exported interface for updating an address on an MTO Shipment
type MTOShipmentAddressUpdater interface {
	UpdateMTOShipmentAddress(newAddress *models.Address, mtoShipmentID uuid.UUID, eTag string, mustBeAvailableToPrime bool) (*models.Address, error)
}

// ShipmentRouter is used for setting the status on shipments at different stages
//go:generate mockery --name ShipmentRouter
type ShipmentRouter interface {
	Submit(shipment *models.MTOShipment) error
	Approve(shipment *models.MTOShipment) error
	RequestCancellation(shipment *models.MTOShipment) error
	Cancel(shipment *models.MTOShipment) error
	Reject(shipment *models.MTOShipment, rejectionReason *string) error
	RequestDiversion(shipment *models.MTOShipment) error
	ApproveDiversion(shipment *models.MTOShipment) error
}
