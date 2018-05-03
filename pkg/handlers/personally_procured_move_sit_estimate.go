package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

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
	sitZip3 := rateengine.Zip5ToZip3(params.DestinationZip)
	cwtWeight := unit.Pound(params.WeightEstimate).ToCWT()
	plannedMoveDateTime := time.Time(params.PlannedMoveDate)
	sitTotal, err := engine.SitCharge(cwtWeight, int(params.DaysInStorage), sitZip3, plannedMoveDateTime, true)
	if err != nil {
		return responseForError(h.logger, err)
	}

	// Most PPMs use COS D, but when there is no COS D rate, the calculation is based on Code 2
	_, sitDiscount, err := models.FetchDiscountRates(h.db, params.OriginZip, params.DestinationZip, "D", plannedMoveDateTime)
	if err != nil {
		if err != models.ErrFetchNotFound {
			return responseForError(h.logger, err)
		}
		h.logger.Info("Couldn't find SIT Discount for COS D, trying with COS 2.", zap.String("origin_zip", params.OriginZip), zap.String("destination_zip", params.DestinationZip), zap.Time("move_date", plannedMoveDateTime), zap.Error(err))
		_, sitDiscount, err = models.FetchDiscountRates(h.db, params.OriginZip, params.DestinationZip, "2", time.Time(params.PlannedMoveDate))
		if err != nil {
			h.logger.Info("Couldn't find SIT Discount for COS 2.", zap.String("origin_zip", params.OriginZip), zap.String("destination_zip", params.DestinationZip), zap.Time("move_date", plannedMoveDateTime), zap.Error(err))
			return responseForError(h.logger, err)
		}
		h.logger.Info("Found SIT Discount for TDL with COS 2.", zap.String("origin_zip", params.OriginZip), zap.String("destination_zip", params.DestinationZip), zap.Time("move_date", plannedMoveDateTime))
	} else {
		h.logger.Info("Found SIT Discount for TDL with COS D.", zap.String("origin_zip", params.OriginZip), zap.String("destination_zip", params.DestinationZip), zap.Time("move_date", plannedMoveDateTime))
	}

	inverseDiscount := (100.00 - sitDiscount) / 100
	// Swagger returns int64 when using the integer type
	sitCharge := int64(sitTotal.MultiplyFloat64(inverseDiscount))
	ppmSitEstimate := internalmessages.PPMSitEstimate{
		Estimate: &sitCharge,
	}
	return ppmop.NewShowPPMSitEstimateOK().WithPayload(&ppmSitEstimate)
}
