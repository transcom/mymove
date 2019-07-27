package internalapi

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
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
	moveID, _ := uuid.FromString(params.MoveID.String())
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	originDutyStationZip := move.Orders.ServiceMember.DutyStation.Address.PostalCode
	if !session.IsOfficeUser() {
		return ppmop.NewShowPPMIncentiveForbidden()
	}

	engine := rateengine.NewRateEngine(h.DB(), logger)

	distanceMilesFromOriginPickupZip, err := h.Planner().Zip5TransitDistance(params.OriginZip, params.DestinationZip)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	distanceMilesFromOriginDutyStationZip, err := h.Planner().Zip5TransitDistance(originDutyStationZip, params.DestinationZip)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	cost, err := engine.ComputeLowestCostPPMMove(
		unit.Pound(params.Weight),
		params.OriginZip,
		originDutyStationZip,
		params.DestinationZip,
		distanceMilesFromOriginPickupZip,
		distanceMilesFromOriginDutyStationZip,
		time.Time(params.OriginalMoveDate),
		0,
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
