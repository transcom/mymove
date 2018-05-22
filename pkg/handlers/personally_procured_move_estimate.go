package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"time"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// ShowPPMEstimateHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMEstimateHandler HandlerContext

// Handle calculates a PPM reimbursement range.
func (h ShowPPMEstimateHandler) Handle(params ppmop.ShowPPMEstimateParams) middleware.Responder {
	engine := rateengine.NewRateEngine(h.db, h.logger, h.planner)

	lhDiscount, _, err := PPMDiscountFetch(h.db,
		h.logger,
		params.OriginZip,
		params.DestinationZip,
		time.Time(params.PlannedMoveDate),
	)
	if err != nil {
		return responseForError(h.logger, err)
	}

	weightEstimate := params.WeightEstimate
	weightProrate := float64(1)
	// If weight is less than 1000, prorate the rate for 1000 lbs
	if weightEstimate < 1000 {
		weightProrate = float64(weightEstimate) / 1000.0
		weightEstimate = 1000
	}

	cost, err := engine.ComputePPM(unit.Pound(weightEstimate),
		params.OriginZip,
		params.DestinationZip,
		time.Time(params.PlannedMoveDate),
		0, // We don't want any SIT charges
		lhDiscount,
		0.0,
	)

	if err != nil {
		return responseForError(h.logger, err)
	}

	min := cost.GCC.MultiplyFloat64(0.95).MultiplyFloat64(weightProrate)
	max := cost.GCC.MultiplyFloat64(1.05).MultiplyFloat64(weightProrate)

	ppmEstimate := internalmessages.PPMEstimateRange{
		RangeMin: swag.Int64(min.Int64()),
		RangeMax: swag.Int64(max.Int64()),
	}
	return ppmop.NewShowPPMEstimateOK().WithPayload(&ppmEstimate)
}
