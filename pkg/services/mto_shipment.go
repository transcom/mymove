package services

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/transcom/mymove/pkg/models"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
)

// MTOShipmentFetcher is the exported interface for FetchMTOShipment
//go:generate mockery -name MTOShipmentFetcher
type MTOShipmentFetcher interface {
	FetchMTOShipment(mtoShipmentID uuid.UUID) (*models.MTOShipment, error)
}

//MTOShipmentUpdater is the service object interface for UpdateMTOShipment
//go:generate mockery -name MTOShipmentUpdater
type MTOShipmentUpdater interface {
	UpdateMTOShipment(unmodifiedSince time.Time, mtoShipmentPayload *primemessages.MTOShipment) (*models.MTOShipment, error)
}

// MTOShipmentStatusUpdater is the exported interface for updating an MTO shipment status
//go:generate mockery -name MTOShipmentStatusUpdater
type MTOShipmentStatusUpdater interface {
	UpdateMTOShipmentStatus(payload mtoshipmentops.PatchMTOShipmentStatusParams) (*models.MTOShipment, error)
}
