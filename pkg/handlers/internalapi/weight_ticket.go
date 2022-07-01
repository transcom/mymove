package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	weightticketops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// CreateWeightTicketHandler
type CreateWeightTicketHandler struct {
	handlers.HandlerConfig
	weightTicketCreator services.WeightTicketCreator
}

// Handle creates a weight ticket
// Depending on the SO, may need to change the document params to weight ticket params
func (h CreateWeightTicketHandler) Handle(params weightticketops.CreateWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// NO NEED FOR payload_to_model, will need for Update
			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				appCtx.Logger().Error("internalapi.CreateWeightTicketHandler", zap.Error(err))
			}
			// ADD AN ERROR CHECK HERE for ppmShipmentID
			var weightTicket *models.WeightTicket
			weightTicket, err = h.weightTicketCreator.CreateWeightTicket(appCtx, ppmShipmentID)

			if err != nil {
				appCtx.Logger().Error("internalapi.CreateWeightTicketHandler", zap.Error(err))
				// Can get a status error
				// Can get an DB error - does the weight ticket, doc create?
				// Can get an error for whether the PPM exist
				// ADD SWITCH STATEMENT
			}
			returnPayload := payloads.CreateWeightTicket(weightTicket)
			return weightticketops.NewCreateWeightTicketOK().WithPayload(returnPayload), nil
		})
}
