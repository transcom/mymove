package pptasapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/appcontext"
	moveop "github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations/moves"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
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

			payload := pptasmessages.GetMovesSinceResponse{}

			return moveop.NewMovesSinceOK().WithPayload(&payload), nil
		})
}
