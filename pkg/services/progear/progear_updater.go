package weightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type progearUpdater struct {
}

// NewCustomerProgearUpdater creates a new progearUpdater struct with the basic checks
func NewCustomerProgearUpdater() services.ProgearCreator {
	return &progearCreator{}
}

func (f *progearUpdater) UpdateProgear(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.ProgearWeightTicket, error) {
	progear := models.ProgearWeightTicket{}
	return &progear, nil
}
