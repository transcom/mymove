package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
)

// ShowPpmSitEstimateHandler returns PPM SIT estimate for a weight, move date,
type ShowPpmSitEstimateHandler HandlerContext

// Handle retrieves orders in the system belonging to the logged in user given order ID
func (h ShowPpmSitEstimate) Handle(params ppmop.ShowPpmSitEstimateParams) middleware.Responder {
	sitZip3 := rateengine.zip5ToZip3(params.destinationZip)
	sitTotal, err := rateengine.sitCharge(params.WeightEstimate, params.DaysInStorage, sitZip3, params.PlannedMoveDate, true)
	if err != nil {
		return responseForError(h.logger, err)
	}

	sitDiscount := FetchDiscountRates(params.originZip, params.destinationZip, "D", params.PlannedMoveDate)
	inverseDiscount := 1.00 - sitDiscount
	sitCharge := sitTotal * inverseDiscount
	ppmSitEstimate := internalmessages.PpmSitEstimate{
		Estimate: sitCharge,
	}
	return ppmop.NewShowOrdersOK().WithPayload(ppmSitEstimate)
}
