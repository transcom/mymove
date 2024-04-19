package pptasapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	pptasop "github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations/moves"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/pptasapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// IndexMovesHandler returns a list of moves/MTOs via GET /moves
type IndexMovesHandler struct {
	handlers.HandlerConfig
	services.MoveListFetcher
	services.MoveSearcher
}

// Handle retrieves a list of moves
func (h IndexMovesHandler) Handle(params pptasop.MovesSinceParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			searchMovesParams := services.SearchMovesParams{
				MoveCreatedDate: handlers.FmtDateTimePtrToPopPtr(&params.Body.MoveSinceDate),
			}

			moves, _, err := h.SearchMoves(appCtx, &searchMovesParams)

			if err != nil {
				appCtx.Logger().Error("Error searching for move", zap.Error(err))
				return pptasop.NewMovesSinceInternalServerError(), err
			}

			movesToReturn := moves[0:params.Body.NumMoves]

			payload := payloads.MovesSince(appCtx, movesToReturn)
			return pptasop.NewMovesSinceOK().WithPayload(payload), nil
		})
}
