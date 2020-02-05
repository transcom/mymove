package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// ShowPPMIncentiveHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMIncentiveHandler struct {
	handlers.HandlerContext
}

// Handle calculates a PPM reimbursement range.
func (h ShowPPMIncentiveHandler) Handle(params ppmop.ShowPPMIncentiveParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	if !session.IsOfficeUser() {
		return ppmop.NewShowPPMIncentiveForbidden()
	}

	ordersID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	move, err := models.FetchMoveByOrderID(h.DB(), ordersID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	engine := rateengine.NewRateEngine(h.DB(), logger, move)

	destinationZip, err := GetDestinationDutyStationPostalCode(h.DB(), ordersID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	distanceMilesFromOriginPickupZip, err := h.Planner().Zip5TransitDistance(params.OriginZip, destinationZip)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	distanceMilesFromOriginDutyStationZip, err := h.Planner().Zip5TransitDistance(params.OriginDutyStationZip, destinationZip)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	cost, err := engine.ComputeLowestCostPPMMove(
		unit.Pound(params.Weight),
		params.OriginZip,
		params.OriginDutyStationZip,
		destinationZip,
		distanceMilesFromOriginPickupZip,
		distanceMilesFromOriginDutyStationZip,
		time.Time(params.OriginalMoveDate),
		0, // We don't want any SIT charges
	)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	gcc := cost.GCC
	incentivePercentage := cost.GCC.MultiplyFloat64(0.95)

	ppmObligation := internalmessages.PPMIncentive{
		Gcc:                 swag.Int64(gcc.Int64()),
		IncentivePercentage: swag.Int64(incentivePercentage.Int64()),
	}
	return ppmop.NewShowPPMIncentiveOK().WithPayload(&ppmObligation)
}
