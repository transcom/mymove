package primeapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

type ListMoveTaskOrdersHandler struct {
	handlers.HandlerContext
}

func (h ListMoveTaskOrdersHandler) Handle(params movetaskorderops.ListMoveTaskOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	var mtos models.MoveTaskOrders

	query := h.DB().Q()
	if params.Since != nil {
		since := time.Unix(*params.Since, 0)
		query = query.Where("updated_at > ?", since)
	}

	err := query.All(&mtos)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return movetaskorderops.NewListMoveTaskOrdersInternalServerError()
	}

	payload := make(primemessages.MoveTaskOrders, len(mtos))

	for i, m := range mtos {
		payload[i] = payloadForMoveTaskOrder(m)
	}

	return movetaskorderops.NewListMoveTaskOrdersOK().WithPayload(payload)
}

func payloadForMoveTaskOrder(moveTaskOrder models.MoveTaskOrder) *primemessages.MoveTaskOrder {
	destinationAddress := payloadForAddress(&moveTaskOrder.DestinationAddress)
	pickupAddress := payloadForAddress(&moveTaskOrder.PickupAddress)
	entitlements := payloadForEntitlements(&moveTaskOrder.Entitlements)
	payload := &primemessages.MoveTaskOrder{
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
		Status:                 string(moveTaskOrder.Status),
		UpdatedAt:              strfmt.Date(moveTaskOrder.UpdatedAt),
	}
	return payload
}

func payloadForAddress(a *models.Address) *primemessages.Address {
	if a == nil {
		return nil
	}
	return &primemessages.Address{
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

func payloadForEntitlements(entitlement *models.GHCEntitlement) *primemessages.Entitlements {
	if entitlement == nil {
		return nil
	}
	return &primemessages.Entitlements{
		DependentsAuthorized:  entitlement.DependentsAuthorized,
		NonTemporaryStorage:   handlers.FmtBool(entitlement.NonTemporaryStorage),
		PrivatelyOwnedVehicle: handlers.FmtBool(entitlement.PrivatelyOwnedVehicle),
		ProGearWeight:         int64(entitlement.ProGearWeight),
		ProGearWeightSpouse:   int64(entitlement.ProGearWeightSpouse),
		StorageInTransit:      int64(entitlement.StorageInTransit),
		TotalDependents:       int64(entitlement.TotalDependents),
	}
}
