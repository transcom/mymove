package weightticket

import (
	"errors"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
)

type weightTicketDeleter struct {
	services.WeightTicketFetcher
	services.PPMEstimator
}

func NewWeightTicketDeleter(fetcher services.WeightTicketFetcher, estimator services.PPMEstimator) services.WeightTicketDeleter {
	return &weightTicketDeleter{
		fetcher,
		estimator,
	}
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
func (d *weightTicketDeleter) DeleteWeightTicket(appCtx appcontext.AppContext, weightTicketID uuid.UUID) error {
	weightTicket, err := d.GetWeightTicket(appCtx, weightTicketID)
	if err != nil {
		return err
	}
	oldPPM, err := ppmshipment.FindPPMShipmentAndWeightTickets(appCtx, weightTicket.PPMShipmentID)
	if err != nil {
		return err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// All weightTicket documents are belongs_to relations, so will not be automatically
		// deleted when we call SoftDestroy on the weight ticket
		err = utilities.SoftDestroy(appCtx.DB(), &weightTicket.EmptyDocument)
		if err != nil {
			return handleSoftDestroyError(err)
		}
		err = utilities.SoftDestroy(appCtx.DB(), &weightTicket.FullDocument)
		if err != nil {
			return handleSoftDestroyError(err)
		}
		err = utilities.SoftDestroy(appCtx.DB(), &weightTicket.ProofOfTrailerOwnershipDocument)
		if err != nil {
			return handleSoftDestroyError(err)
		}
		err = utilities.SoftDestroy(appCtx.DB(), weightTicket)
		if err != nil {
			return handleSoftDestroyError(err)
		}
		newPPM, err := ppmshipment.FindPPMShipmentAndWeightTickets(appCtx, weightTicket.PPMShipmentID)
		if err != nil {
			return err
		}

		finalIncentive, err := d.PPMEstimator.FinalIncentiveWithDefaultChecks(appCtx, *oldPPM, newPPM)
		if err != nil {
			return err
		}
		newPPM.FinalIncentive = finalIncentive
		verrs, err := appCtx.DB().ValidateAndUpdate(newPPM)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
