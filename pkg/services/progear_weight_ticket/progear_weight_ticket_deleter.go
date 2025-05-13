package progearweightticket

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type progearWeightTicketDeleter struct {
}

func NewProgearWeightTicketDeleter() services.ProgearWeightTicketDeleter {
	return &progearWeightTicketDeleter{}
}

func (d *progearWeightTicketDeleter) DeleteProgearWeightTicket(appCtx appcontext.AppContext, ppmID uuid.UUID, progearWeightTicketID uuid.UUID) error {
	var ppmShipment models.PPMShipment
	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment.MoveTaskOrder.Orders",
			"ProgearWeightTickets",
		).
		Find(&ppmShipment, ppmID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFoundError(progearWeightTicketID, "while looking for ProgearWeightTicket")
		}
		return apperror.NewQueryError("Progear Weight Ticket fetch original", err, "")
	}

	if appCtx.Session().IsMilApp() {
		if ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID != appCtx.Session().ServiceMemberID && !appCtx.Session().IsOfficeUser() {
			wrongServiceMemberIDErr := apperror.NewForbiddenError("Attempted delete by wrong service member")
			appCtx.Logger().Error("internalapi.DeleteWeightTicketHandler", zap.Error(wrongServiceMemberIDErr))
			return wrongServiceMemberIDErr
		}
	}

	found := false
	for _, lineItem := range ppmShipment.ProgearWeightTickets {
		if lineItem.ID == progearWeightTicketID {
			found = true
			break
		}
	}
	if !found {
		mismatchedPPMShipmentAndProgearWeightTicketIDErr := apperror.NewNotFoundError(progearWeightTicketID, "Pro-gear weight ticket does not exist on ppm shipment")
		appCtx.Logger().Error("internalapi.DeleteProGearWeightTicketHandler", zap.Error(mismatchedPPMShipmentAndProgearWeightTicketIDErr))
		return mismatchedPPMShipmentAndProgearWeightTicketIDErr
	}

	progearWeightTicket, err := FetchProgearWeightTicketByIDExcludeDeletedUploads(appCtx, progearWeightTicketID)
	if err != nil {
		return err
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
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

		// Grab all ProgearWeightTickets from DB except the deleted ones
		var tickets []models.ProgearWeightTicket
		if err := appCtx.DB().
			Q().
			Scope(utilities.ExcludeDeletedScope(models.ProgearWeightTicket{})).
			Where("ppm_shipment_id = ?", ppmID).
			All(&tickets); err != nil {
			return apperror.NewQueryError("fetching ProgearWeightTickets", err, "")
		}

		// Total up the tickets
		totalSelf := 0
		totalSpouse := 0
		for _, t := range tickets {
			if t.BelongsToSelf != nil && *t.BelongsToSelf {
				if t.Weight != nil {
					totalSelf += int(*t.Weight)
				}
			} else {
				if t.Weight != nil {
					totalSpouse += int(*t.Weight)
				}
			}
		}

		// Update actual progear weight and actual spouse progear weight in mto_shipment
		if err := appCtx.DB().RawQuery(
			`UPDATE mto_shipments
				SET actual_pro_gear_weight        = NULLIF(?, 0),
					actual_spouse_pro_gear_weight = NULLIF(?, 0)
			WHERE id = (
				SELECT shipment_id
					From ppm_shipments
				WHERE id = ?
				)`,
			totalSelf, totalSpouse, ppmID,
		).Exec(); err != nil {
			return apperror.NewQueryError("MTOShipment update actual pro-gear weights", err, "")
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
