package services

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// MTOShipmentStatusUpdater is the exported interface for updating an MTO shipment status
//go:generate mockery -name MTOShipmentStatusUpdater
type MTOShipmentStatusUpdater interface {
	UpdateMTOShipmentStatus(id uuid.UUID, status string) (*validate.Errors, error)
}
