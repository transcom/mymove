package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"

	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// GetMoveHistoryHandler gets a move history by locator
type GetMoveHistoryHandler struct {
	handlers.HandlerContext
	services.MoveHistoryFetcher
}

// Handle handles the getMoveHistory by locator request
func (h GetMoveHistoryHandler) Handle(params moveop.GetMoveHistoryParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			locator := params.Locator
			if locator == "" {
				return moveop.NewGetMoveHistoryBadRequest(), apperror.NewBadDataError("missing required parameter: locator")
			}

			move, err := h.FetchMoveHistory(appCtx, locator)

			if err != nil {
				appCtx.Logger().Error("Error retrieving move history by locator", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return moveop.NewGetMoveHistoryNotFound(), err
				default:
					return moveop.NewGetMoveHistoryInternalServerError(), err
				}
			}

			payload := payloads.MoveHistory(move)
			return moveop.NewGetMoveHistoryOK().WithPayload(payload), nil
		})
}
