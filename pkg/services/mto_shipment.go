package services

import (
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/models"
)

// MTOShipmentStatusUpdater is the exported interface for updating an MTO shipment status
//go:generate mockery -name MTOShipmentStatusUpdater
type MTOShipmentStatusUpdater interface {
	UpdateMTOShipmentStatus(payload mtoshipmentops.PatchMTOShipmentStatusParams) (*models.MTOShipment, error)
}
