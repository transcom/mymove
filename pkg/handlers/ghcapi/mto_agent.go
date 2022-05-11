package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	mtoagentop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_agent"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

//ListMTOAgentsHandler is a struct for the handler.
type ListMTOAgentsHandler struct {
	handlers.HandlerContext
	services.ListFetcher
}

//Handle handles the handling for listing MTO Agents.
func (h ListMTOAgentsHandler) Handle(params mtoagentop.FetchMTOAgentListParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			mtoShipmentID, err := uuid.FromString(params.ShipmentID.String())
			// Return parsing sadness
			if err != nil {
				parsingError := fmt.Errorf("UUID Parsing for %s: %w", "MTOShipmentID", err).Error()
				appCtx.Logger().Error(parsingError)
				payload := payloadForValidationError("UUID(s) parsing error", parsingError, h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
				return mtoagentop.NewFetchMTOAgentListUnprocessableEntity().WithPayload(payload), err
			}

			// Let's set up our filter for the service object call
			queryFilters := []services.QueryFilter{
				query.NewQueryFilter("mto_shipment_id", "=", mtoShipmentID.String()),
			}
			var mtoAgents models.MTOAgents
			err = h.FetchRecordList(appCtx, &mtoAgents, queryFilters, nil, nil, nil)
			// return errors
			if err != nil {
				if err.Error() == "FETCH_NOT_FOUND" {
					appCtx.Logger().Error(fmt.Sprintf("Error while fetching mto agents. Could not find record with mto shipment with id: %s", mtoShipmentID.String()), zap.Error(err))
					return mtoagentop.NewFetchMTOAgentListNotFound(), err
				}
				appCtx.Logger().Error(fmt.Sprintf("Error fetching mto agents for mto shipment with id: %s", mtoShipmentID.String()), zap.Error(err))
				return mtoagentop.NewFetchMTOAgentListInternalServerError(), err
			}

			returnPayload := payloads.MTOAgents(&mtoAgents)
			return mtoagentop.NewFetchMTOAgentListOK().WithPayload(*returnPayload), nil
		})
}
