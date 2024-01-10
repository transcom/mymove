package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type PPMCloseout interface {
	GetPPMCloseout(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) *models.PPMCloseout
}
