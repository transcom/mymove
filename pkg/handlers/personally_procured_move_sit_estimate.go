package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// ShowPpmSitEstimateHandler returns PPM SIT estimate for a weight, move date,
type ShowPpmSitEstimateHandler HandlerContext

// Handle calculates SIT charge and retrieves SIT discount rate.
// It returns the discount rate applied to relevant SIT charge.
func (h ShowPpmSitEstimateHandler) Handle(params ppmop.ShowPpmSitEstimateParams) middleware.Responder {
	engine := rateengine.NewRateEngine(h.db, h.logger, h.planner)
	_, sitZip3 := engine.Zip5ToZip3(params.OriginZip, params.DestinationZip)
	cwtWeight := unit.Pound(params.WeightEstimate).ToCWT()
	sitTotal, err := engine.SitCharge(cwtWeight, int(params.DaysInStorage), sitZip3, timeFromDateTime(params.PlannedMoveDate), true)
	if err != nil {
		return responseForError(h.logger, err)
	}

	// TODO: Determine COS for PPMs
	_, sitDiscount, err := models.FetchDiscountRates(h.db, params.OriginZip, params.DestinationZip, "D", timeFromDateTime(params.PlannedMoveDate))
	if err != nil {
		return responseForError(h.logger, err)
	}

	inverseDiscount := (100.00 - sitDiscount) / 100
	sitCharge := int64(sitTotal.MultiplyFloat64(inverseDiscount))
	ppmSitEstimate := internalmessages.PpmSitEstimate{
		Estimate: &sitCharge,
	}
	return ppmop.NewShowPpmSitEstimateOK().WithPayload(&ppmSitEstimate)
}
