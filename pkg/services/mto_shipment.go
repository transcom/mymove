package services

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

//MTOShipmentUpdater is the service object interface for UpdateMTOShipment
//go:generate mockery --name MTOShipmentUpdater --disable-version-string
type MTOShipmentUpdater interface {
	CheckIfMTOShipmentCanBeUpdated(mtoShipment *models.MTOShipment, session *auth.Session) (bool, error)
	MTOShipmentsMTOAvailableToPrime(mtoShipmentID uuid.UUID) (bool, error)
	RetrieveMTOShipment(mtoShipmentID uuid.UUID) (*models.MTOShipment, error)
	UpdateMTOShipmentOffice(ctx context.Context, mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error)
	UpdateMTOShipmentCustomer(ctx context.Context, mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error)
	UpdateMTOShipmentPrime(ctx context.Context, mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error)
}

//ShipmentDeleter is the service object interface for deleting a shipment
//go:generate mockery --name ShipmentDeleter --disable-version-string
type ShipmentDeleter interface {
	DeleteShipment(shipmentID uuid.UUID) (uuid.UUID, error)
}

//ShipmentApprover is the service object interface for approving a shipment
//go:generate mockery --name ShipmentApprover --disable-version-string
type ShipmentApprover interface {
	ApproveShipment(shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

//ShipmentDiversionRequester is the service object interface for approving a shipment diversion
//go:generate mockery --name ShipmentDiversionRequester --disable-version-string
type ShipmentDiversionRequester interface {
	RequestShipmentDiversion(shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

//ShipmentDiversionApprover is the service object interface for approving a shipment diversion
//go:generate mockery --name ShipmentDiversionApprover --disable-version-string
type ShipmentDiversionApprover interface {
	ApproveShipmentDiversion(shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

//ShipmentRejecter is the service object interface for approving a shipment
//go:generate mockery --name ShipmentRejecter --disable-version-string
type ShipmentRejecter interface {
	RejectShipment(shipmentID uuid.UUID, eTag string, reason *string) (*models.MTOShipment, error)
}

//ShipmentCancellationRequester is the service object interface for approving a shipment diversion
//go:generate mockery --name ShipmentCancellationRequester --disable-version-string
type ShipmentCancellationRequester interface {
	RequestShipmentCancellation(shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

//ShipmentReweighRequester is the service object interface for approving a shipment diversion
//go:generate mockery --name ShipmentReweighRequester --disable-version-string
type ShipmentReweighRequester interface {
	RequestShipmentReweigh(ctx context.Context, shipmentID uuid.UUID) (*models.Reweigh, error)
}

// MTOShipmentStatusUpdater is the exported interface for updating an MTO shipment status
//go:generate mockery --name MTOShipmentStatusUpdater --disable-version-string
type MTOShipmentStatusUpdater interface {
	UpdateMTOShipmentStatus(shipmentID uuid.UUID, status models.MTOShipmentStatus, rejectionReason *string, eTag string) (*models.MTOShipment, error)
}

// MTOShipmentCreator is the exported interface for creating a payment request
//go:generate mockery --name MTOShipmentCreator --disable-version-string
type MTOShipmentCreator interface {
	CreateMTOShipment(MTOShipment *models.MTOShipment, MTOServiceItems models.MTOServiceItems) (*models.MTOShipment, error)
}

// MTOShipmentAddressUpdater is the exported interface for updating an address on an MTO Shipment
type MTOShipmentAddressUpdater interface {
	UpdateMTOShipmentAddress(newAddress *models.Address, mtoShipmentID uuid.UUID, eTag string, mustBeAvailableToPrime bool) (*models.Address, error)
}

// ShipmentRouter is used for setting the status on shipments at different stages
//go:generate mockery --name ShipmentRouter --disable-version-string
type ShipmentRouter interface {
	Submit(shipment *models.MTOShipment) error
	Approve(shipment *models.MTOShipment) error
	RequestCancellation(shipment *models.MTOShipment) error
	Cancel(shipment *models.MTOShipment) error
	Reject(shipment *models.MTOShipment, rejectionReason *string) error
	RequestDiversion(shipment *models.MTOShipment) error
	ApproveDiversion(shipment *models.MTOShipment) error
}
