package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"

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
	setIDErr := setUpdateMTOAgentIDs(mtoAgent, agentID, mtoShipmentID)
	if setIDErr != nil {
		return mtoshipmentops.NewUpdateMTOAgentUnprocessableEntity().WithPayload(
			payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), setIDErr.ValidationErrors))
	}

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

// setUpdateMTOAgentIDs sets the ID values from the path on the MTOAgent model
// and also checks that no conflicting values are present
func setUpdateMTOAgentIDs(agent *models.MTOAgent, agentID uuid.UUID, mtoShipmentID uuid.UUID) *services.InvalidInputError {
	verrs := validate.NewErrors()

	if agent.ID != agentID && agent.ID != uuid.Nil {
		verrs.Add("id", "must match the agentID in the path or be omitted from the request")
	}

	if agent.MTOShipmentID != mtoShipmentID && agent.MTOShipmentID != uuid.Nil {
		verrs.Add("mtoShipmentID", "must match the mtoShipmentID in the path or be omitted from the request")
	}

	if verrs.HasAny() {
		err := services.NewInvalidInputError(agentID, nil, verrs, "Invalid input found in the request.")
		return &err
	}

	// Set the values on the model if everything was well:
	agent.ID = agentID
	agent.MTOShipmentID = mtoShipmentID

	return nil
}
