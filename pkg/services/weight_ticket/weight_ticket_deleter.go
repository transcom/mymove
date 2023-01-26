package weightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
)

type weightTicketDeleter struct {
}

func NewWeightTicketDeleter() services.WeightTicketDeleter {
	return &weightTicketDeleter{}
}

func (d *weightTicketDeleter) DeleteWeightTicket(appCtx appcontext.AppContext, weightTicketID uuid.UUID) error {
	return nil
}
