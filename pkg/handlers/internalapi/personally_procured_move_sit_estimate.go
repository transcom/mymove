package internalapi

import (
	"time"

	"github.com/transcom/mymove/pkg/models"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// ShowPPMSitEstimateHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMSitEstimateHandler struct {
	handlers.HandlerContext
}

// Handle calculates SIT charge and retrieves SIT discount rate.
// It returns the discount rate applied to relevant SIT charge.
func (h ShowPPMSitEstimateHandler) Handle(params ppmop.ShowPPMSitEstimateParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	engine := rateengine.NewRateEngine(h.DB(), logger)

	ordersID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	destinationZip, err := GetDestionDutyStationPostalCode(h.DB(), ordersID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	sitZip3 := rateengine.Zip5ToZip3(destinationZip)
	cwtWeight := unit.Pound(params.WeightEstimate).ToCWT()
	originalMoveDateTime := time.Time(params.OriginalMoveDate)

	lhDiscount, sitDiscount, err := models.PPMDiscountFetch(h.DB(),
		logger,
		params.OriginZip,
		destinationZip,
		time.Time(params.OriginalMoveDate),
	)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	sitComputation, err := engine.SitCharge(cwtWeight, int(params.DaysInStorage), sitZip3, originalMoveDateTime, true)

	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	// Swagger returns int64 when using the integer type
	sitCharge := int64(sitComputation.ApplyDiscount(lhDiscount, sitDiscount))

	ppmSitEstimate := internalmessages.PPMSitEstimate{
		Estimate: &sitCharge,
	}
	return ppmop.NewShowPPMSitEstimateOK().WithPayload(&ppmSitEstimate)
}
