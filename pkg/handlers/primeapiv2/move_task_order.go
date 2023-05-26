package primeapiv2

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primev2api/primev2operations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapiv2/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// ListMovesHandler lists moves with the option to filter since a particular date. Optimized ver.
type ListMovesHandler struct {
	handlers.HandlerConfig
	services.MoveTaskOrderFetcherV2
}

// Handle fetches all moves with the option to filter since a particular date. Optimized version.
func (h ListMovesHandler) Handle(params movetaskorderops.ListMovesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			var searchParams services.MoveTaskOrderFetcherParams
			if params.Since != nil {
				since := handlers.FmtDateTimePtrToPop(params.Since)
				searchParams.Since = &since
			}

			mtos, err := h.MoveTaskOrderFetcherV2.ListPrimeMoveTaskOrdersV2(appCtx, &searchParams)

			if err != nil {
				appCtx.Logger().Error("Unexpected error while fetching moves:", zap.Error(err))
				return movetaskorderops.NewListMovesInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			payload := payloads.ListMoves(&mtos)

			return movetaskorderops.NewListMovesOK().WithPayload(payload), nil
		})
}
