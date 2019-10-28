package ghcapi

import (
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime/middleware"

	//TODO why is this being named move_task_order in generated code. maybe just rename in import?
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateMoveTaskOrderStatusHandlerFunc fetches an access code associated with a service member
type UpdateMoveTaskOrderStatusHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderStatusUpdater services.MoveTaskOrderStatusUpdater
}

// UpdateMoveTaskOrderStatusHandlerFunc updates the status of a MoveTaskOrder
func (h UpdateMoveTaskOrderStatusHandlerFunc) Handle(params move_task_order.UpdateMoveTaskOrderStatusParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	// TODO how are we going to handle auth in new api? Do we need some sort of placeholder to remind us to
	// TODO to revist?
	moveTaskOrderID, status := requestToModels(params)
	mto, err := h.moveTaskOrderStatusUpdater.UpdateMoveTaskOrderStatus(moveTaskOrderID, status)
	if err != nil {
		logger.Error("ghciap.MoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorder.ErrNotFound:
			return move_task_order.NewUpdateMoveTaskOrderStatusNotFound()
		case movetaskorder.ErrInvalidInput:
			return move_task_order.NewUpdateMoveTaskOrderStatusBadRequest()
		default:
			return move_task_order.NewUpdateMoveTaskOrderStatusInternalServerError()
		}
	}
	moveTaskOrderPayload := payloadForMoveTaskOrder(*mto)
	return move_task_order.NewUpdateMoveTaskOrderStatusOK().WithPayload(moveTaskOrderPayload)
}

func requestToModels(params move_task_order.UpdateMoveTaskOrderStatusParams) (uuid.UUID, models.MoveTaskOrderStatus) {
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	status := models.MoveTaskOrderStatus(params.Body.Status)
	return moveTaskOrderID, status
}

// TODO probably could write some some tests for these mappers.
func payloadForMoveTaskOrder(moveTaskOrder models.MoveTaskOrder) *ghcmessages.MoveTaskOrder {
	serviceItems := payloadForServiceItems(moveTaskOrder)
	destinationAddress := payloadForAddress(&moveTaskOrder.DestinationAddress)
	pickupAddress := payloadForAddress(&moveTaskOrder.PickupAddress)
	payload := &ghcmessages.MoveTaskOrder{
		CustomerID:             strfmt.UUID(moveTaskOrder.CustomerID.String()),
		DestinationAddress:     destinationAddress,
		DestinationDutyStation: strfmt.UUID(moveTaskOrder.DestinationDutyStation.ID.String()),
		// TODO the pivotal ticket seems somewhat incomplete compared to the
		// TODO api spec double check that's right
		Entitlements: &ghcmessages.Entitlements{
			StorageInTransit:      moveTaskOrder.SitEntitlement,
			TotalWeightSelf:       moveTaskOrder.WeightEntitlement,
			PrivatelyOwnedVehicle: &moveTaskOrder.POVEntitlement,
			NonTemporaryStorage:   &moveTaskOrder.NTSEntitlement,
		},
		ID:       strfmt.UUID(moveTaskOrder.ID.String()),
		MoveDate: strfmt.Date(moveTaskOrder.RequestedPickupDates),
		MoveID:   strfmt.UUID(moveTaskOrder.MoveID.String()),
		// TODO is UUID in api should it be?
		OriginDutyStation:   strfmt.UUID(moveTaskOrder.OriginDutyStationID.String()),
		PickupAddress:       pickupAddress,
		Remarks:             moveTaskOrder.CustomerRemarks,
		RequestedPickupDate: strfmt.Date(moveTaskOrder.RequestedPickupDates),
		ServiceItems:        serviceItems,
		Status:              string(moveTaskOrder.Status),
		UpdatedAt:           strfmt.Date(moveTaskOrder.UpdatedAt),
	}

	return payload
}

func payloadForAddress(a *models.Address) *ghcmessages.Address {
	if a == nil {
		return nil
	}
	return &ghcmessages.Address{
		ID:             strfmt.UUID(a.ID.String()),
		StreetAddress1: swag.String(a.StreetAddress1),
		StreetAddress2: a.StreetAddress2,
		StreetAddress3: a.StreetAddress3,
		City:           swag.String(a.City),
		State:          swag.String(a.State),
		PostalCode:     swag.String(a.PostalCode),
		Country:        a.Country,
	}
}

func payloadForServiceItems(moveTaskOrder models.MoveTaskOrder) []*ghcmessages.ServiceItem {
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
