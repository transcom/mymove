package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// GetMoveHandler gets a move by locator
type GetMoveHandler struct {
	handlers.HandlerContext
	services.MoveFetcher
}

// Handle handles the handling
func (h GetMoveHandler) Handle(params moveop.GetMoveParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	locator := params.Locator
	if locator == "" {
		return moveop.NewGetMoveBadRequest()
	}

	move, err := h.FetchMove(locator)

	if err != nil {
		logger.Error("Error retrieving move by locator", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return moveop.NewGetMoveNotFound()
		default:
			return moveop.NewGetMoveInternalServerError()
		}
	}

	payload := payloads.Move(move)
	return moveop.NewGetMoveOK().WithPayload(payload)
}
