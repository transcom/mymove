package ghcapi

import (
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime/middleware"

	//TODO why is this being named move_task_order
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForAccessCodeModel(moveTaskOrder models.MoveTaskOrder) *ghcmessages.MoveTaskOrder {
	payload := &ghcmessages.MoveTaskOrder{
		Customer:               moveTaskOrder.Customer,
		DestinationDutyStation: strfmt.UUID(moveTaskOrder.DestinationDutyStation.ID.String()),
		Entitlements: &ghcmessages.Entitlements{
			DependentsAuthorized:  false,
			NonTemporaryStorage:   false,
			PrivatelyOwnedVehicle: false,
			ProGearWeight:         0,
			ProGearWeightSpouse:   0,
			StorageInTransit:      0,
			TotalDependents:       0,
			TotalWeightSelf:       0,
		},
		ID:                  "",
		MoveDate:            strfmt.Date{},
		MoveID:              "",
		MoveTaskOrdersType:  "",
		OriginDutyStation:   "",
		OriginPPSO:          "",
		Remarks:             "",
		RequestedPickupDate: strfmt.Date{},
		ServiceItems:        nil,
		Status:              "",
		UpdatedAt:           strfmt.Date{},
	}

	return payload
}

// FetchAccessCodeHandler fetches an access code associated with a service member
type UpdateMoveTaskOrderHandler struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// NewGhcAPIHandler returns a handler for the GHC API
func (h UpdateMoveTaskOrderHandler) Handle(params move_task_order.UpdateMoveTaskOrderParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if session == nil {
		return move_task_order.NewDeleteMoveTaskOrderForbidden()
	}

	// Fetch access code
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(session.ServiceMemberID)
	if err != nil {
		logger.Error("ghciap.MoveTaskOrderHandler error", zap.Error(err))
	}
	moveTaskOrderPayload := payloadForAccessCodeModel(*mto)
	return move_task_order.NewUpdateMoveTaskOrderOK().WithPayload(moveTaskOrderPayload)
}
