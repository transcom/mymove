package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// WeightTicketCreator creates a WeightTicket that is associated with a PPMShipment
type WeightTicketCreator interface {
	CreateWeightTicket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.WeightTicket, error)
}

// WeightTicketUpdater updates a WeightTicket
type WeightTicketUpdater interface {
	UpdateWeightTicket(appCtx appcontext.AppContext, weightTicket models.WeightTicket, eTag string) (*models.WeightTicket, error)
}
