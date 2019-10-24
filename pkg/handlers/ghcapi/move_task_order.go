package ghcapi

import (
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	//TODO why is this being named move_task_order
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type MoveTaskOrder struct {
	ID                     uuid.UUID
	MoveID                 uuid.UUID            `db:"move_id"`
	CreatedAt              time.Time            `db:"created_at"`
	UpdatedAt              time.Time            `db:"updated_at"`
	Customer               models.ServiceMember `has_one:"service_members" fk_id:"id"`
	OriginDutyStation      models.DutyStation   `has_one:"duty_stations" fk_id:"id"`
	DestinationDutyStation models.DutyStation   `has_one:"duty_stations" fk_id:"id"`
	PickupAddress          models.Address       `has_one:"addresses" fk_id:"id"`
	DestinationAddress     models.Address       `has_one:"addresses" fk_id:"id"`
	RequestedPickupDates   time.Time            `db:"request_pickup_dates"`
	CustomerRemarks        string               `db:"customer_remarks"`
	WeightEntitlement      int                  `db:"weight_entitlement"`
	SitEntitlement         int                  `db:"sit_entitlement"`
	POVEntitlement         bool                 `db:"pov_entitlement"`
	NTSEntitlement         bool                 `db:"nts_entitlement"`
}

func payloadForAccessCodeModel(moveTaskOrder MoveTaskOrder) *ghcmessages.MoveTaskOrder {
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
type MoveTaskOrderHandler struct {
	handlers.HandlerContext
	accessCodeFetcher services.AccessCodeFetcher
}

// NewGhcAPIHandler returns a handler for the GHC API
func (h MoveTaskOrderHandler) Handle(params move_task_order.UpdateMoveTaskOrderParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if session == nil {
		return move_task_order.NewDeleteMoveTaskOrderForbidden()
	}

	// Fetch access code
	accessCode, err := h.accessCodeFetcher.FetchAccessCode(session.ServiceMemberID)
	var fetchAccessCodePayload *internalmessages.AccessCode

	fetchAccessCodePayload = payloadForAccessCodeModel(*moveTaskOrder)

	return move_task_order.NewUpdateMoveTaskOrderOK().WithPayload(fetchAccessCodePayload)
}
