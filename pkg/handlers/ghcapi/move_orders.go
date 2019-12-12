package ghcapi

import (
	"database/sql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	moveorderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
)

// GetCustomerInfoHandler fetches the information of a specific customer
type GetMoveOrdersHandler struct {
	handlers.HandlerContext
}

// Handle getting the information of a specific customer
func (h GetMoveOrdersHandler) Handle(params moveorderop.GetMoveOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveOrderID, _ := uuid.FromString(params.MoveOrderID.String())
	moveOrder := &models.MoveOrder{}
	err := h.DB().Eager("DestinationDutyStation.Address", "OriginDutyStation.Address", "Entitlement").Find(moveOrder, moveOrderID)
	if err != nil {
		logger.Error("fetching move order", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return moveorderop.NewGetMoveOrderNotFound()
		default:
			return moveorderop.NewGetMoveOrderInternalServerError()
		}
	}
	moveOrderPayload := payloads.MoveOrders(moveOrder)
	return moveorderop.NewGetMoveOrderOK().WithPayload(moveOrderPayload)
}
