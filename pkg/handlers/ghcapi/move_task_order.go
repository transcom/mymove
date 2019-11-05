package ghcapi

import (
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	movetaskorderservice "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime/middleware"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateMoveTaskOrderStatusHandlerFunc updates the status of a Move Task Order
type UpdateMoveTaskOrderStatusHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderStatusUpdater services.MoveTaskOrderStatusUpdater
}

// UpdateMoveTaskOrderStatusHandlerFunc updates the status of a MoveTaskOrder
func (h UpdateMoveTaskOrderStatusHandlerFunc) Handle(params movetaskorderops.UpdateMoveTaskOrderStatusParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	// TODO how are we going to handle auth in new api? Do we need some sort of placeholder to remind us to
	// TODO to revisit?
	moveTaskOrderID, status := requestToModels(params)
	mto, err := h.moveTaskOrderStatusUpdater.UpdateMoveTaskOrderStatus(moveTaskOrderID, status)
	if err != nil {
		logger.Error("ghciap.MoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusNotFound()
		case movetaskorderservice.ErrInvalidInput:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusBadRequest()
		default:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError()
		}
	}
	moveTaskOrderPayload := payloadForMoveTaskOrder(*mto)
	return movetaskorderops.NewUpdateMoveTaskOrderStatusOK().WithPayload(moveTaskOrderPayload)
}

func requestToModels(params movetaskorderops.UpdateMoveTaskOrderStatusParams) (uuid.UUID, models.MoveTaskOrderStatus) {
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	status := models.MoveTaskOrderStatus(params.Body.Status)
	return moveTaskOrderID, status
}

// TODO it might make sense to create a package for these various mappers, rather than making them all private
// TODO since some can be reused. Might also make sense to have some mapper specfic tests
func payloadForMoveTaskOrder(moveTaskOrder models.MoveTaskOrder) *ghcmessages.MoveTaskOrder {
	serviceItems := payloadForServiceItems(&moveTaskOrder.ServiceItems)
	destinationAddress := payloadForAddress(&moveTaskOrder.DestinationAddress)
	pickupAddress := payloadForAddress(&moveTaskOrder.PickupAddress)
	entitlements := payloadForEntitlements(&moveTaskOrder.Entitlements)
	payload := &ghcmessages.MoveTaskOrder{
		CustomerID:             strfmt.UUID(moveTaskOrder.CustomerID.String()),
		DestinationAddress:     destinationAddress,
		DestinationDutyStation: strfmt.UUID(moveTaskOrder.DestinationDutyStation.ID.String()),
		Entitlements:           entitlements,
		ID:                     strfmt.UUID(moveTaskOrder.ID.String()),
		MoveDate:               strfmt.Date(moveTaskOrder.RequestedPickupDate),
		MoveID:                 strfmt.UUID(moveTaskOrder.MoveID.String()),
		OriginDutyStation:      strfmt.UUID(moveTaskOrder.OriginDutyStationID.String()),
		PickupAddress:          pickupAddress,
		Remarks:                moveTaskOrder.CustomerRemarks,
		RequestedPickupDate:    strfmt.Date(moveTaskOrder.RequestedPickupDate),
		ServiceItems:           serviceItems,
		Status:                 string(moveTaskOrder.Status),
		UpdatedAt:              strfmt.Date(moveTaskOrder.UpdatedAt),
	}

	if moveTaskOrder.ActualWeight != nil {
		payload.ActualWeight = int64(*moveTaskOrder.ActualWeight)
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

func payloadForServiceItems(serviceItems *models.ServiceItems) []*ghcmessages.ServiceItem {
	var sis models.ServiceItems
	if serviceItems == nil {
		return nil
	}
	var serviceItemsPayload []*ghcmessages.ServiceItem
	sis = *serviceItems
	for _, si := range sis {
		serviceItemsPayload = append(serviceItemsPayload, &ghcmessages.ServiceItem{
			MoveTaskOrderID: strfmt.UUID(si.MoveTaskOrderID.String()),
			CreatedAt:       strfmt.Date(si.CreatedAt),
			UpdatedAt:       strfmt.Date(si.UpdatedAt),
		})
	}
	return serviceItemsPayload
}

func payloadForEntitlements(entitlement *models.GHCEntitlement) *ghcmessages.Entitlements {
	if entitlement == nil {
		return nil
	}
	return &ghcmessages.Entitlements{
		DependentsAuthorized:  entitlement.DependentsAuthorized,
		NonTemporaryStorage:   handlers.FmtBool(entitlement.NonTemporaryStorage),
		PrivatelyOwnedVehicle: handlers.FmtBool(entitlement.PrivatelyOwnedVehicle),
		ProGearWeight:         int64(entitlement.ProGearWeight),
		ProGearWeightSpouse:   int64(entitlement.ProGearWeightSpouse),
		StorageInTransit:      int64(entitlement.StorageInTransit),
		TotalDependents:       int64(entitlement.TotalDependents),
	}
}
