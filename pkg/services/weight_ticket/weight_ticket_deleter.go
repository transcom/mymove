package weightticket

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
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

func (d *weightTicketDeleter) DeleteWeightTicket(appCtx appcontext.AppContext, ppmID uuid.UUID, weightTicketID uuid.UUID) error {
	var ppmShipment models.PPMShipment
	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment.MoveTaskOrder.Orders",
			"WeightTickets",
		).
		Find(&ppmShipment, ppmID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFoundError(weightTicketID, "while looking for WeightTicket")
		}
		return apperror.NewQueryError("WeightTicket fetch original", err, "")
	}

	if appCtx.Session().IsMilApp() {
		if ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID != appCtx.Session().ServiceMemberID && !appCtx.Session().IsOfficeUser() {
			wrongServiceMemberIDErr := apperror.NewForbiddenError("Attempted delete by wrong service member")
			appCtx.Logger().Error("internalapi.DeleteWeightTicketHandler", zap.Error(wrongServiceMemberIDErr))
			return wrongServiceMemberIDErr
		}
	}

	found := false
	for _, lineItem := range ppmShipment.WeightTickets {
		if lineItem.ID == weightTicketID {
			found = true
			break
		}
	}
	if !found {
		mismatchedPPMShipmentAndWeightTicketIDErr := apperror.NewNotFoundError(weightTicketID, "Weight ticket does not exist on ppm shipment")
		appCtx.Logger().Error("internalapi.DeleteWeightTicketHandler", zap.Error(mismatchedPPMShipmentAndWeightTicketIDErr))
		return mismatchedPPMShipmentAndWeightTicketIDErr
	}

	weightTicket, err := d.GetWeightTicket(appCtx, weightTicketID)
	if err != nil {
		return err
	}
	oldPPM, err := ppmshipment.FindPPMShipmentAndWeightTickets(appCtx, weightTicket.PPMShipmentID)
	if err != nil {
		return err
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
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
		newPPM, err := ppmshipment.FindPPMShipmentAndWeightTickets(appCtx, weightTicket.PPMShipmentID)
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
