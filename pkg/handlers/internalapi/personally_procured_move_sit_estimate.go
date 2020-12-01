package internalapi

import (
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
)

// ShowPPMSitEstimateHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMSitEstimateHandler struct {
	handlers.HandlerContext
	services.EstimateCalculator
}

// Handle calculates SIT charge and retrieves SIT discount rate.
// It returns the discount rate applied to relevant SIT charge.
func (h ShowPPMSitEstimateHandler) Handle(params ppmop.ShowPPMSitEstimateParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	ppmID, err := uuid.FromString(params.PersonallyProcuredMoveID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	ppm, err := models.FetchPersonallyProcuredMove(h.DB(), session, ppmID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	ordersID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	move, err := models.FetchMoveByOrderID(h.DB(), ordersID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	var originalMoveDate = time.Time(params.OriginalMoveDate)
	if originalMoveDate.IsZero() {
		logger.Error("original move date invalid")
		return ppmop.NewShowPPMSitEstimateUnprocessableEntity()
	}

	if params.WeightEstimate == 0 {
		logger.Error("weight estimate required")
		return ppmop.NewShowPPMSitEstimateUnprocessableEntity()
	}
	weightEstimate := unit.Pound(params.WeightEstimate)

	// construct ppm with estimated values
	estimatedPPM := *ppm
	estimatedPPM.OriginalMoveDate = &originalMoveDate
	estimatedPPM.PickupPostalCode = &params.OriginZip
	estimatedPPM.DaysInStorage = &params.DaysInStorage
	estimatedPPM.WeightEstimate = &weightEstimate

	sitCharge, _, err := h.CalculateEstimates(&estimatedPPM, move.ID, logger)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	ppmSitEstimate := internalmessages.PPMSitEstimate{
		Estimate: &sitCharge,
	}
	return ppmop.NewShowPPMSitEstimateOK().WithPayload(&ppmSitEstimate)
}
