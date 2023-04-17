package weightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
)

type weightTicketDeleter struct {
	services.WeightTicketFetcher
	services.PPMShipmentFetcher
	services.PPMEstimator
}

func NewWeightTicketDeleter(weightTicketFetcher services.WeightTicketFetcher, ppmShipmentFetcher services.PPMShipmentFetcher, estimator services.PPMEstimator) services.WeightTicketDeleter {
	return &weightTicketDeleter{
		weightTicketFetcher,
		ppmShipmentFetcher,
		estimator,
	}
}

func (d *weightTicketDeleter) DeleteWeightTicket(appCtx appcontext.AppContext, weightTicketID uuid.UUID) error {
	weightTicket, err := d.GetWeightTicket(appCtx, weightTicketID)
	if err != nil {
		return err
	}

	oldPPM, err := d.GetPPMShipment(
		appCtx,
		weightTicket.PPMShipmentID,
		[]string{ppmshipment.EagerPreloadAssociationShipment, ppmshipment.EagerPreloadAssociationWeightTickets},
		nil,
	)

	if err != nil {
		return err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// All weightTicket documents are belongs_to relations, so will not be automatically
		// deleted when we call SoftDestroy on the weight ticket
		err = utilities.SoftDestroy(appCtx.DB(), &weightTicket.EmptyDocument)
		if err != nil {
			return err
		}
		err = utilities.SoftDestroy(appCtx.DB(), &weightTicket.FullDocument)
		if err != nil {
			return err
		}
		err = utilities.SoftDestroy(appCtx.DB(), &weightTicket.ProofOfTrailerOwnershipDocument)
		if err != nil {
			return err
		}
		err = utilities.SoftDestroy(appCtx.DB(), weightTicket)
		if err != nil {
			return err
		}
		newPPM, err := d.GetPPMShipment(
			appCtx,
			weightTicket.PPMShipmentID,
			[]string{ppmshipment.EagerPreloadAssociationShipment, ppmshipment.EagerPreloadAssociationWeightTickets},
			nil,
		)
		if err != nil {
			return err
		}

		finalIncentive, err := d.PPMEstimator.FinalIncentiveWithDefaultChecks(appCtx, *oldPPM, newPPM)
		if err != nil {
			return err
		}

		// Only update PPM if the incentive has changed
		if finalIncentive != oldPPM.FinalIncentive || (finalIncentive != nil && oldPPM.FinalIncentive != nil && *finalIncentive != *oldPPM.FinalIncentive) {
			newPPM.FinalIncentive = finalIncentive
			verrs, err := appCtx.DB().ValidateAndUpdate(newPPM)
			if err != nil {
				return err
			}
			if verrs.HasAny() {
				return verrs
			}
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
