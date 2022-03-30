package internalapi

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
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
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			ppmID, err := uuid.FromString(params.PersonallyProcuredMoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			ppm, err := models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), ppmID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			ordersID, err := uuid.FromString(params.OrdersID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			move, err := models.FetchMoveByOrderID(appCtx.DB(), ordersID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			var originalMoveDate = time.Time(params.OriginalMoveDate)
			if originalMoveDate.IsZero() {
				errMsg := "original move date invalid"
				appCtx.Logger().Error(errMsg)
				return ppmop.NewShowPPMSitEstimateUnprocessableEntity(), apperror.NewBadDataError(errMsg)
			}

			if params.WeightEstimate == 0 {
				errMsg := "weight estimate required"
				appCtx.Logger().Error(errMsg)
				return ppmop.NewShowPPMSitEstimateUnprocessableEntity(), apperror.NewBadDataError(errMsg)
			}
			weightEstimate := unit.Pound(params.WeightEstimate)

			// construct ppm with estimated values
			estimatedPPM := *ppm
			estimatedPPM.OriginalMoveDate = &originalMoveDate
			estimatedPPM.PickupPostalCode = &params.OriginZip
			estimatedPPM.DaysInStorage = &params.DaysInStorage
			estimatedPPM.WeightEstimate = &weightEstimate

			sitCharge, _, err := h.CalculateEstimates(appCtx, &estimatedPPM, move.ID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			ppmSitEstimate := internalmessages.PPMSitEstimate{
				Estimate: &sitCharge,
			}
			return ppmop.NewShowPPMSitEstimateOK().WithPayload(&ppmSitEstimate), nil
		})
}
