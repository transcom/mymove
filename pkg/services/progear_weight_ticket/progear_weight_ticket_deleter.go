package progearweightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
)

type progearWeightTicketDeleter struct {
}

func NewProgearWeightTicketDeleter() services.ProgearWeightTicketDeleter {
	return &progearWeightTicketDeleter{}
}

func (d *progearWeightTicketDeleter) DeleteProgearWeightTicket(appCtx appcontext.AppContext, progearWeightTicketID uuid.UUID) error {
	return nil
}
