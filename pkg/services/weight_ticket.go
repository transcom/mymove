package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// WeightTicketCreator creates a WeightTicket that is associated with a PPMShipment
//
//go:generate mockery --name WeightTicketCreator
type WeightTicketCreator interface {
	CreateWeightTicket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.WeightTicket, error)
}

// WeightTicketFetcher fetches a WeightTicket that is associated with a PPMShipment
//
//go:generate mockery --name WeightTicketFetcher
type WeightTicketFetcher interface {
	GetWeightTicket(appCtx appcontext.AppContext, weightTicketID uuid.UUID) (*models.WeightTicket, error)
}

// WeightTicketUpdater updates a WeightTicket
//
//go:generate mockery --name WeightTicketUpdater
type WeightTicketUpdater interface {
	UpdateWeightTicket(appCtx appcontext.AppContext, weightTicket models.WeightTicket, eTag string) (*models.WeightTicket, error)
}

// WeightTicketDeleter deletes a WeightTicket
//
//go:generate mockery --name WeightTicketDeleter
type WeightTicketDeleter interface {
	DeleteWeightTicket(appCtx appcontext.AppContext, ppmID uuid.UUID, weightTicketID uuid.UUID) error
}
