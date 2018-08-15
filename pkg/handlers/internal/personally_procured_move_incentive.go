package internal

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/auth"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// ShowPPMIncentiveHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMIncentiveHandler utils.HandlerContext

// Handle calculates a PPM reimbursement range.
func (h ShowPPMIncentiveHandler) Handle(params ppmop.ShowPPMIncentiveParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if !session.IsOfficeUser() {
		return ppmop.NewShowPPMIncentiveForbidden()
	}
	engine := rateengine.NewRateEngine(h.Db, h.Logger, h.Planner)

	lhDiscount, _, err := PPMDiscountFetch(h.Db,
		h.Logger,
		params.OriginZip,
		params.DestinationZip,
		time.Time(params.PlannedMoveDate),
	)
	if err != nil {
		return responseForError(h.Logger, err)
	}

	cost, err := engine.ComputePPM(unit.Pound(params.Weight),
		params.OriginZip,
		params.DestinationZip,
		time.Time(params.PlannedMoveDate),
		0, // We don't want any SIT charges
		lhDiscount,
		0.0,
	)

	if err != nil {
		return responseForError(h.Logger, err)
	}

	gcc := cost.GCC
	incentivePercentage := cost.GCC.MultiplyFloat64(0.95)

	ppmObligation := internalmessages.PPMIncentive{
		Gcc:                 swag.Int64(gcc.Int64()),
		IncentivePercentage: swag.Int64(incentivePercentage.Int64()),
	}
	return ppmop.NewShowPPMIncentiveOK().WithPayload(&ppmObligation)
}
