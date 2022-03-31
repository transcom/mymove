package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// ShowPPMEstimateHandler returns PPM estimate for a weight, move date, origin zip, order id
type ShowPPMEstimateHandler struct {
	handlers.HandlerContext
}

// Handle calculates a PPM reimbursement range.
func (h ShowPPMEstimateHandler) Handle(params ppmop.ShowPPMEstimateParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			ordersID, err := uuid.FromString(params.OrdersID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			move, err := models.FetchMoveByOrderID(appCtx.DB(), ordersID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			engine := rateengine.NewRateEngine(move)

			destinationZip, err := GetDestinationDutyLocationPostalCode(appCtx, ordersID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			distanceMilesFromOriginPickupZip, err := h.Planner().Zip5TransitDistanceLineHaul(appCtx, params.OriginZip, destinationZip)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			distanceMilesFromOriginDutyLocationZip, err := h.Planner().Zip5TransitDistanceLineHaul(appCtx, params.OriginDutyLocationZip, destinationZip)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			costDetails, err := engine.ComputePPMMoveCosts(
				appCtx,
				unit.Pound(params.WeightEstimate),
				params.OriginZip,
				params.OriginDutyLocationZip,
				destinationZip,
				distanceMilesFromOriginPickupZip,
				distanceMilesFromOriginDutyLocationZip,
				time.Time(params.OriginalMoveDate),
				0, // We don't want any SIT charges
			)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			cost := rateengine.GetWinningCostMove(costDetails)

			min := cost.GCC.MultiplyFloat64(0.95)
			max := cost.GCC.MultiplyFloat64(1.05)

			ppmEstimate := internalmessages.PPMEstimateRange{
				RangeMin: swag.Int64(min.Int64()),
				RangeMax: swag.Int64(max.Int64()),
			}
			return ppmop.NewShowPPMEstimateOK().WithPayload(&ppmEstimate), nil
		})
}
