package gunsafeweightticket

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

type gunsafeWeightTicketDeleter struct {
}

func NewGunSafeWeightTicketDeleter() services.GunSafeWeightTicketDeleter {
	return &gunsafeWeightTicketDeleter{}
}

func (d *gunsafeWeightTicketDeleter) DeleteGunSafeWeightTicket(appCtx appcontext.AppContext, ppmID uuid.UUID, gunsafeWeightTicketID uuid.UUID) error {
	var ppmShipment models.PPMShipment
	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment.MoveTaskOrder.Orders",
			"GunSafeWeightTickets",
		).
		Find(&ppmShipment, ppmID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFoundError(gunsafeWeightTicketID, "while looking for GunSafeWeightTicket")
		}
		return apperror.NewQueryError("GunSafe Weight Ticket fetch original", err, "")
	}

	if appCtx.Session().IsMilApp() {
		if ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID != appCtx.Session().ServiceMemberID && !appCtx.Session().IsOfficeUser() {
			wrongServiceMemberIDErr := apperror.NewForbiddenError("Attempted delete by wrong service member")
			appCtx.Logger().Error("internalapi.DeleteWeightTicketHandler", zap.Error(wrongServiceMemberIDErr))
			return wrongServiceMemberIDErr
		}
	}

	found := false
	for _, lineItem := range ppmShipment.GunSafeWeightTickets {
		if lineItem.ID == gunsafeWeightTicketID {
			found = true
			break
		}
	}
	if !found {
		mismatchedPPMShipmentAndGunSafeWeightTicketIDErr := apperror.NewNotFoundError(gunsafeWeightTicketID, "Gun safe weight ticket does not exist on ppm shipment")
		appCtx.Logger().Error("internalapi.DeleteGunSafeWeightTicketHandler", zap.Error(mismatchedPPMShipmentAndGunSafeWeightTicketIDErr))
		return mismatchedPPMShipmentAndGunSafeWeightTicketIDErr
	}

	gunsafeWeightTicket, err := FetchGunSafeWeightTicketByIDExcludeDeletedUploads(appCtx, gunsafeWeightTicketID)
	if err != nil {
		return err
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
		// gunsafeWeightTicket.Document is a belongs_to relation, so will not be automatically
		// deleted when we call SoftDestroy on the moving expense
		err = utilities.SoftDestroy(appCtx.DB(), &gunsafeWeightTicket.Document)
		if err != nil {
			return err
		}
		err = utilities.SoftDestroy(appCtx.DB(), gunsafeWeightTicket)
		if err != nil {
			return err
		}

		if err := appCtx.DB().
			RawQuery("SELECT update_actual_gun_safe_weight_totals($1)", ppmID).
			Exec(); err != nil {
			return apperror.NewQueryError("update_actual_gunsafe_weight_totals", err, "")
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
