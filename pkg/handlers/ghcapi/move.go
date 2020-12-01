package ghcapi

import (
	"fmt"
	"regexp"

	"github.com/go-openapi/runtime/middleware"

	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"

	"go.uber.org/zap"
)

// GetMoveHandler gets a move by locator
type GetMoveHandler struct {
	handlers.HandlerContext
	services.Fetcher
	services.NewQueryFilter
}

// Handle handles the handling
func (h GetMoveHandler) Handle(params moveop.GetMoveParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	locator := params.Locator

	move := &models.Move{}
	err := h.Fetcher.FetchRecord(move,
		[]services.QueryFilter{query.NewQueryFilter("locator", "=", locator)})

	notFoundErr := regexp.MustCompile("Resource not found:")

	if err != nil {
		if notFoundErr.MatchString(err.Error()) {
			logger.Error(fmt.Sprintf("No move found with locator %s", locator), zap.Error(err))
			return moveop.NewGetMoveNotFound()
		}
		logger.Error(fmt.Sprintf("Error fetching move with locator: %s", locator), zap.Error(err))
		return moveop.NewGetMoveInternalServerError()
	}

	payload := payloads.Move(move)
	return moveop.NewGetMoveOK().WithPayload(payload)
}
