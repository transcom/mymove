package ghcapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"

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

//ListMTOAgentsHandler is a struct for the handler.
type ListMTOAgentsHandler struct {
	handlers.HandlerContext
	services.ListFetcher
}

//Handle handles the handling for listing MTO Agents.
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
	var mtoAgents models.MTOAgents
	err = h.FetchRecordList(&mtoAgents, queryFilters, query.NewQueryAssociations([]services.QueryAssociation{}), nil, nil)
	// return errors
	if err != nil {
		if err.Error() == "FETCH_NOT_FOUND" {
			logger.Error(fmt.Sprintf("Error while fetching mto agents. Could not find record with mto shipment with id: %s", mtoShipmentID.String()), zap.Error(err))
			return mtoagentop.NewFetchMTOAgentListNotFound()
		}
		logger.Error(fmt.Sprintf("Error fetching mto agents for mto shipment with id: %s", mtoShipmentID.String()), zap.Error(err))
		return mtoagentop.NewFetchMTOAgentListInternalServerError()
	}

	returnPayload := payloads.MTOAgents(&mtoAgents)
	return mtoagentop.NewFetchMTOAgentListOK().WithPayload(*returnPayload)
}
