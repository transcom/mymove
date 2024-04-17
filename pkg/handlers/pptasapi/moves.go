package pptasapi

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
	moveop "github.com/transcom/mymove/pkg/gen/pptasoperations/moves"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// IndexMovesHandler returns a list of moves/MTOs via GET /moves
type IndexMovesHandler struct {
	handlers.HandlerConfig
	services.MoveListFetcher
	services.NewQueryFilter
	services.NewPagination
}

var locatorFilterConverters = map[string]func(string) []services.QueryFilter{
	"locator": func(content string) []services.QueryFilter {
		return []services.QueryFilter{query.NewQueryFilter("locator", "=", content)}
	},
}

// Handle retrieves a list of moves
func (h IndexMovesHandler) Handle(params moveop.MovesSinceParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			pagination := h.NewPagination(params.Page, params.PerPage)
			queryFilters := generateQueryFilters(appCtx.Logger(), params.Filter, locatorFilterConverters)
			queryAssociations := []services.QueryAssociation{
				query.NewQueryAssociation("Orders.ServiceMember"),
			}
			associations := query.NewQueryAssociationsPreload(queryAssociations)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			moveList, err := h.FetchMoveList(appCtx, queryFilters, associations, pagination, ordering)
			log.Output(1, moveList[0].Locator)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			payload := pptasmessages.GetMovesSinceResponse{}

			return moveop.NewIndexMovesOK().WithPayload(payload), nil
		})
}
