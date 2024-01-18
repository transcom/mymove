package progearweightticket

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	progearops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/services"
)

type progearWeightTicketDeleter struct {
}

func NewProgearWeightTicketDeleter() services.ProgearWeightTicketDeleter {
	return &progearWeightTicketDeleter{}
}

func (d *progearWeightTicketDeleter) DeleteProgearWeightTicket(appCtx appcontext.AppContext, ppmID uuid.UUID, progearWeightTicketID uuid.UUID) (middleware.Responder, error) {
	progearWeightTicket, err := FetchProgearWeightTicketByIDExcludeDeletedUploads(appCtx, ppmID, progearWeightTicketID)
	if err != nil {
		if err == apperror.NewSessionError("Attempted delete by wrong service member") {
			return progearops.NewDeleteProGearWeightTicketForbidden(), err
		}

		return progearops.NewDeleteProGearWeightTicketNotFound(), err
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
		return progearops.NewDeleteProGearWeightTicketInternalServerError(), transactionError
	}

	return nil, nil
}
