package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// PPMCloseoutFetcher fetches all of the necessary calculations needed for display when the SC is reviewing a closeout
//
//go:generate mockery --name PPMCloseoutFetcher
type PPMCloseoutFetcher interface {
	GetPPMCloseout(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMCloseout, error)
	GetActualWeight(ppmShipment *models.PPMShipment) unit.Pound
}
