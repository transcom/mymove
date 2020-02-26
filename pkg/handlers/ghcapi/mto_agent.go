package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"

	mtoagentop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_agent"
)

type ListMTOAgentsHandler struct {
	handlers.HandlerContext
	services.MTOAgentListFetcher
}

func (h ListMTOAgentsHandler) Handle(params mtoagentop.FetchMTOAgentListParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	mtoShipmentID, err := uuid.FromString(params.ShipmentID.String())
	// Return parsing sadness
	if err != nil {
		parsingError := fmt.Errorf("UUID Parsing for %s: %w", "MTOShipmentID", err).Error()
		logger.Error(parsingError)
		payload := payloadForValidationError("UUID(s) parsing error", parsingError, h.GetTraceID(), validate.NewErrors())
		return mtoagentop.NewFetchMTOAgentListUnprocessableEntity().WithPayload(payload)
	}

	// Let's set up our filter for the service object call
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("mto_shipment_id", "=", mtoShipmentID.String()),
	}

	mtoAgents, err := h.FetchMTOAgentList(queryFilters)
	// return errors
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching mto agents for mto shipment with id: %s", mtoShipmentID.String()), zap.Error(err))
		return mtoagentop.NewFetchMTOAgentListInternalServerError()
	}

	if mtoAgents == nil {
		logger.Error(fmt.Sprintf("Found 0 mto agents for mto shipment id: %s", mtoShipmentID.String()))
		return mtoagentop.NewFetchMTOAgentListNotFound()
	}

	returnPayload := payloads.MTOAgents(mtoAgents)
	return mtoagentop.NewFetchMTOAgentListOK().WithPayload(*returnPayload)
}