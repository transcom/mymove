package progearweightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/services"
)

type progearWeightTicketDeleter struct {
}

func NewProgearWeightTicketDeleter() services.ProgearWeightTicketDeleter {
	return &progearWeightTicketDeleter{}
}

func (d *progearWeightTicketDeleter) DeleteProgearWeightTicket(appCtx appcontext.AppContext, progearWeightTicketID uuid.UUID) error {
	progearWeightTicket, err := FetchProgearWeightTicketByIDExcludeDeletedUploads(appCtx, progearWeightTicketID)
	if err != nil {
		return err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// progearWeightTicket.Document is a belongs_to relation, so will not be automatically
		// deleted when we call SoftDestroy on the moving expense
		err = utilities.SoftDestroy(appCtx.DB(), &progearWeightTicket.Document)
		if err != nil {
			return err
		}
		err = utilities.SoftDestroy(appCtx.DB(), progearWeightTicket)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
