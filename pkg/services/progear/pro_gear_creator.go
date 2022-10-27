package weightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type progearCreator struct {
}

// NewCustomerProgearCreator creates a new weightTicketCreator struct with the basic checks
func NewCustomerProgearCreator() services.ProgearCreator {
	return &progearCreator{}
}

func (f *progearCreator) CreateProgear(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.ProgearWeightTicket, error) {
	progear := models.ProgearWeightTicket{}
	return &progear, nil
}
