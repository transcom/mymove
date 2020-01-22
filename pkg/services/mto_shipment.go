package services

import (
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
)

// MTOShipmentStatusUpdater is the exported interface for updating an MTO shipment status
//go:generate mockery -name MTOShipmentStatusUpdater
type MTOShipmentStatusUpdater interface {
	UpdateMTOShipmentStatus(model *models.MTOShipment, status string) (*validate.Errors, error)
}
