package weightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type progearCreator struct {
}

func (f *progearCreator) CreateProgear(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.ProgearWeightTicket, error) {
	progear := models.ProgearWeightTicket{}
	return &progear, nil
}
