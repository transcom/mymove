package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// ShowPPMSitEstimateHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMSitEstimateHandler HandlerContext

// Handle calculates SIT charge and retrieves SIT discount rate.
// It returns the discount rate applied to relevant SIT charge.
func (h ShowPPMSitEstimateHandler) Handle(params ppmop.ShowPPMSitEstimateParams) middleware.Responder {
	engine := rateengine.NewRateEngine(h.db, h.logger, h.planner)
	sitZip3 := engine.Zip5ToZip3(params.DestinationZip)
	cwtWeight := unit.Pound(params.WeightEstimate).ToCWT()
	sitTotal, err := engine.SitCharge(cwtWeight, int(params.DaysInStorage), sitZip3, time.Time(params.PlannedMoveDate), true)
	if err != nil {
		return responseForError(h.logger, err)
	}

	// Most PPMs use COS D, but when there is no COS D rate, the calculation is based on Code 2
	_, sitDiscount, err := models.FetchDiscountRates(h.db, params.OriginZip, params.DestinationZip, "D", time.Time(params.PlannedMoveDate))
	if err != nil {
		if errors.Cause(err).Error() != models.RecordNotFoundErrorString {
			return responseForError(h.logger, err)
		}
		_, sitDiscount, err = models.FetchDiscountRates(h.db, params.OriginZip, params.DestinationZip, "2", time.Time(params.PlannedMoveDate))
		if err != nil {
			return responseForError(h.logger, err)
		}
	}

	inverseDiscount := (100.00 - sitDiscount) / 100
	// Swagger returns int64 when using the integer type
	sitCharge := int64(sitTotal.MultiplyFloat64(inverseDiscount))
	ppmSitEstimate := internalmessages.PPMSitEstimate{
		Estimate: &sitCharge,
	}
	return ppmop.NewShowPPMSitEstimateOK().WithPayload(&ppmSitEstimate)
}
