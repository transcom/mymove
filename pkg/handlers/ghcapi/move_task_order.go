package ghcapi

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime/middleware"

	//TODO why is this being named move_task_order
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForMoveTaskOrder(moveTaskOrder models.MoveTaskOrder) *ghcmessages.MoveTaskOrder {
	serviceItems := serviceItemResponseMapper(moveTaskOrder)
	payload := &ghcmessages.MoveTaskOrder{
		Customer:               moveTaskOrder.Customer,
		DestinationDutyStation: strfmt.UUID(moveTaskOrder.DestinationDutyStation.ID.String()),
		// TODO the pivotal ticket seems somewhat incomplete compared to the
		// TODO api spec double check that's right
		Entitlements: &ghcmessages.Entitlements{
			StorageInTransit: moveTaskOrder.SitEntitlement,
			TotalWeightSelf:  moveTaskOrder.WeightEntitlement,
		},
		ID:       strfmt.UUID(moveTaskOrder.ID.String()),
		MoveDate: strfmt.Date(moveTaskOrder.RequestedPickupDates),
		MoveID:   strfmt.UUID(moveTaskOrder.MoveID.String()),
		// TODO is UUID in api should it be?
		OriginDutyStation:   strfmt.UUID(moveTaskOrder.OriginDutyStationID.String()),
		Remarks:             moveTaskOrder.CustomerRemarks,
		RequestedPickupDate: strfmt.Date(moveTaskOrder.RequestedPickupDates),
		ServiceItems:        serviceItems,
		Status:              string(moveTaskOrder.Status),
		UpdatedAt:           strfmt.Date(moveTaskOrder.UpdatedAt),
	}

	return payload
}

func serviceItemResponseMapper(moveTaskOrder models.MoveTaskOrder) []*ghcmessages.ServiceItem {
	var serviceItems []*ghcmessages.ServiceItem
	for _, si := range moveTaskOrder.ServiceItems {
		serviceItems = append(serviceItems, &ghcmessages.ServiceItem{
			MoveTaskOrderID: strfmt.UUID(si.MoveTaskOrderID.String()),
			CreatedAt:       strfmt.Date(si.CreatedAt),
			UpdatedAt:       strfmt.Date(si.UpdatedAt),
		})
	}
	return serviceItems
}

// FetchAccessCodeHandler fetches an access code associated with a service member
type UpdateMoveTaskOrderStatusHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// NewGhcAPIHandler returns a handler for the GHC API
func (h UpdateMoveTaskOrderStatusHandlerFunc) Handle(params move_task_order.UpdateMoveTaskOrderStatusParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	// TODO how are we going to handle auth in new api
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		logger.Error("ghciap.MoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorder.ErrNotFound:
			return move_task_order.NewUpdateMoveTaskOrderStatusNotFound()
		default:
			return move_task_order.NewUpdateMoveTaskOrderStatusInternalServerError()
		}
	}
	moveTaskOrderPayload := payloadForMoveTaskOrder(*mto)
	return move_task_order.NewUpdateMoveTaskOrderStatusOK().WithPayload(moveTaskOrderPayload)
}
