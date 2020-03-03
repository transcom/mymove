package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

//MTOShipmentUpdater is the service object interface for UpdateMTOShipment
//go:generate mockery -name MTOShipmentUpdater
type MTOShipmentUpdater interface {
	UpdateMTOShipment(mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error)
}

// MTOShipmentStatusUpdater is the exported interface for updating an MTO shipment status
//go:generate mockery -name MTOShipmentStatusUpdater
type MTOShipmentStatusUpdater interface {
	UpdateMTOShipmentStatus(shipmentID uuid.UUID, status models.MTOShipmentStatus, rejectionReason *string, eTag string) (*models.MTOShipment, error)
}
