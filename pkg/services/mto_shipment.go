package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/transcom/mymove/pkg/models"
)

// MTOShipmentFetcher is the exported interface for FetchMTOShipment
//go:generate mockery -name MTOShipmentFetcher
type MTOShipmentFetcher interface {
	FetchMTOShipment(mtoShipmentID uuid.UUID) (*models.MTOShipment, error)
}

//MTOShipmentUpdater is the service object interface for UpdateMTOShipment
//go:generate mockery -name MTOShipmentUpdater
type MTOShipmentUpdater interface {
	UpdateMTOShipment(mtoShipmentPayload *primemessages.MTOShipment) (*models.MTOShipment, error)
}
