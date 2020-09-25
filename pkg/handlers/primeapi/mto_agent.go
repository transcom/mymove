package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	mtoagent "github.com/transcom/mymove/pkg/services/mto_agent"

	"github.com/transcom/mymove/pkg/services"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
)

// UpdateMTOAgentHandler is the handler to update an agent
type UpdateMTOAgentHandler struct {
	handlers.HandlerContext
	MTOAgentUpdater services.MTOAgentUpdater
}

// Handle updates an MTO Agent for a shipment
func (h UpdateMTOAgentHandler) Handle(params mtoshipmentops.UpdateMTOAgentParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	// Get the params and payload
	payload := params.Body
	eTag := params.IfMatch
	mtoShipmentID := uuid.FromStringOrNil(params.MtoShipmentID.String())
	agentID := uuid.FromStringOrNil(params.AgentID.String())

	// Get the new agent model
	mtoAgent := payloads.MTOAgentModel(payload)
	mtoAgent.ID = agentID // TODO set IDs method w/ error
	mtoAgent.MTOShipmentID = mtoShipmentID

	// Call the service object
	updatedAgent, err := h.MTOAgentUpdater.UpdateMTOAgent(mtoAgent, eTag, mtoagent.UpdateMTOAgentPrimeValidator)

	// Convert the errors into error responses to return to caller
	if err != nil {
		logger.Error("primeapi.UpdateMTOAgentHandler", zap.Error(err))

		switch e := err.(type) {
		// PreconditionFailedError -> Precondition Failed Response
		case services.PreconditionFailedError:
			return mtoshipmentops.NewUpdateMTOAgentPreconditionFailed().WithPayload(
				payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		// NotFoundError -> Not Found Response
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateMTOAgentNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		// InvalidInputError -> Unprocessable Entity Response
		case services.InvalidInputError:
			return mtoshipmentops.NewUpdateMTOAgentUnprocessableEntity().WithPayload(
				payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))
		// ConflictError -> Conflict Error Response
		case services.ConflictError:
			return mtoshipmentops.NewUpdateMTOAgentConflict().WithPayload(
				payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
		// QueryError -> Internal Server Error
		case services.QueryError:
			if e.Unwrap() != nil {
				logger.Error("primeapi.UpdateMTOAgentHandler error", zap.Error(e.Unwrap()))
			}
			return mtoshipmentops.NewUpdateMTOAgentInternalServerError().WithPayload(
				payloads.InternalServerError(nil, h.GetTraceID()))
		// Unknown -> Internal Server Error
		default:
			return mtoshipmentops.NewUpdateMTOAgentInternalServerError().
				WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}

	}

	// If no error, create a successful payload to return
	mtoAgentPayload := payloads.MTOAgent(updatedAgent)
	return mtoshipmentops.NewUpdateMTOAgentOK().WithPayload(mtoAgentPayload)
}
