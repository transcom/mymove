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

// ShowPPMEstimateHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMEstimateHandler struct {
	handlers.HandlerContext
}

// Handle calculates a PPM reimbursement range.
func (h ShowPPMEstimateHandler) Handle(params ppmop.ShowPPMEstimateParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	originDutyStationZip := move.Orders.ServiceMember.DutyStation.Address.PostalCode

	distanceMilesFromOriginPickupZip, err := h.Planner().Zip5TransitDistance(params.OriginZip, params.DestinationZip)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	distanceMilesFromOriginDutyStationZip, err := h.Planner().Zip5TransitDistance(originDutyStationZip, params.DestinationZip)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	engine := rateengine.NewRateEngine(h.DB(), logger)

	cost, err := engine.ComputeLowestCostPPMMove(
		unit.Pound(params.WeightEstimate),
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

	min := cost.GCC.MultiplyFloat64(0.95)
	max := cost.GCC.MultiplyFloat64(1.05)

	ppmEstimate := internalmessages.PPMEstimateRange{
		RangeMin: swag.Int64(min.Int64()),
		RangeMax: swag.Int64(max.Int64()),
	}
	return ppmop.NewShowPPMEstimateOK().WithPayload(&ppmEstimate)
}
