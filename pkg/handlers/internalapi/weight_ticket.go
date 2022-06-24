package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	weightticketops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"go.uber.org/zap"
)

// Create weightTicketHandler
type CreateWeightTicketHandler struct {
	handlers.HandlerConfig
	weightTicketCreator services.WeightTicketCreator
}

// Handle creates a weight ticket
// Depending on the SO, may need to change the doument params to weight ticket params
func (h CreateWeightTicketHandler) Handle(params weightticketops.CreateWightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// Get payload
			payload := params.Body
			if payload == nil {
				missingBodyErr := apperror.NotFoundError{}
				appCtx.logger().Error(missingBodyErr.Error())
				return weightticketops.NewCreateWightTicketBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage, "Wight Ticket body cannot be empty", h.GetTraceIDFromRequest(params.HTTPRequest))), missingBodyErr
			}
			weightTicket := payloads.WeightTicketModelFromCreate(payload)
			var err error

			weightTicket, err = h.weightTicketCreator.CreateWightTicket(appCtx, weightTicket)

			if err != nil {
				appCtx.Logger().Error("internalapi.CreateWeightTicketHandler", zap.Error(err))
				//TODO: maybe add a switch statement here?
			}
			returnPayload := payloads.WeightTicket(weightTicket)
			return weightticketops.NewCreateWeightTicketOK().WithPayload(returnPayload), nil
		})
}
