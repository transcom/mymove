package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ProgearWeightTickerCreator creates a ProgearWeightTicket that is associated with a PPMShipment
//
//go:generate mockery --name ProgearWeightTicketCreator --disable-version-string
type ProgearWeightTicketCreator interface {
	CreateProgearWeightTicket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.ProgearWeightTicket, error)
}

// ProgearWeightTicketUpdater updates a ProgearWeightTicket
//
//go:generate mockery --name ProgearWeightTicketUpdater --disable-version-string
type ProgearWeightTicketUpdater interface {
	UpdateProgearWeightTicket(appCtx appcontext.AppContext, progearWeightTicket models.ProgearWeightTicket, eTag string) (*models.ProgearWeightTicket, error)
}

// ProgearWeightTicketDeleter deletes a ProgearWeightTicket
//
//go:generate mockery --name ProgearWeightTicketDeleter --disable-version-string
type ProgearWeightTicketDeleter interface {
	DeleteProgearWeightTicket(appCtx appcontext.AppContext, progearWeightTicketID uuid.UUID) error
}
