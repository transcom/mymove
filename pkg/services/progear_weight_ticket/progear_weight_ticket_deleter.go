package progearweightticket

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	progearops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type progearWeightTicketDeleter struct {
}

func NewProgearWeightTicketDeleter() services.ProgearWeightTicketDeleter {
	return &progearWeightTicketDeleter{}
}

func (d *progearWeightTicketDeleter) DeleteProgearWeightTicket(appCtx appcontext.AppContext, ppmID uuid.UUID, progearWeightTicketID uuid.UUID) (middleware.Responder, error) {
	var ppmShipment models.PPMShipment
	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment.MoveTaskOrder.Orders",
			"ProgearWeightTickets",
		).
		Find(&ppmShipment, ppmID)

	if ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID != appCtx.Session().ServiceMemberID {
		wrongServiceMemberIDErr := apperror.NewSessionError("Attempted delete by wrong service member")
		appCtx.Logger().Error("internalapi.DeleteProgearWeightTicketHandler", zap.Error(wrongServiceMemberIDErr))
		return progearops.NewDeleteProGearWeightTicketForbidden(), wrongServiceMemberIDErr
	}

	progearWeightTicket, err := FetchProgearWeightTicketByIDExcludeDeletedUploads(appCtx, progearWeightTicketID)
	if err != nil {
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
