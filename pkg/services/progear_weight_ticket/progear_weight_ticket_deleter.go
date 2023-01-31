package progearweightticket

import (
	"errors"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/services"
)

type progearWeightTicketDeleter struct {
}

func NewProgearWeightTicketDeleter() services.ProgearWeightTicketDeleter {
	return &progearWeightTicketDeleter{}
}

func handleSoftDestroyError(err error) error {
	if err == nil {
		return nil
	}
	switch err.Error() {
	case "error updating model":
		return apperror.NewUnprocessableEntityError("while updating model")
	case "this model does not have deleted_at field":
		return apperror.NewPreconditionFailedError(uuid.Nil, errors.New("model or sub table missing deleted_at field"))
	default:
		return apperror.NewInternalServerError("failed attempt to soft delete model")
	}
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
			return handleSoftDestroyError(err)
		}
		err = utilities.SoftDestroy(appCtx.DB(), progearWeightTicket)
		if err != nil {
			return handleSoftDestroyError(err)
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
