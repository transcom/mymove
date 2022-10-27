package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ProgearCreator creates a Progear that is associated with a PPMShipment
//
//go:generate mockery --name ProgearCreator --disable-version-string
type ProgearCreator interface {
	CreateProgear(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.ProgearWeightTicket, error)
}

// // ProgearUpdater updates a Progear
// //
// //go:generate mockery --name ProgearUpdater --disable-version-string
type ProgearUpdater interface {
	UpdateProgear(appCtx appcontext.AppContext, progear models.ProgearWeightTicket, eTag string) (*models.ProgearWeightTicket, error)
}
