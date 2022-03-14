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

// ShowPPMIncentiveHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMIncentiveHandler struct {
	handlers.HandlerContext
}

// Handle calculates a PPM reimbursement range.
func (h ShowPPMIncentiveHandler) Handle(params ppmop.ShowPPMIncentiveParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {
			if !appCtx.Session().IsOfficeUser() {
				return ppmop.NewShowPPMIncentiveForbidden()
			}

			ordersID, err := uuid.FromString(params.OrdersID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			move, err := models.FetchMoveByOrderID(appCtx.DB(), ordersID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			engine := rateengine.NewRateEngine(move)

			destinationZip, err := GetDestinationDutyLocationPostalCode(appCtx, ordersID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			distanceMilesFromOriginPickupZip, err := h.Planner().Zip5TransitDistanceLineHaul(appCtx, params.OriginZip, destinationZip)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			distanceMilesFromOriginDutyLocationZip, err := h.Planner().Zip5TransitDistanceLineHaul(appCtx, params.OriginDutyLocationZip, destinationZip)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			costDetails, err := engine.ComputePPMMoveCosts(
				appCtx,
				unit.Pound(params.Weight),
				params.OriginZip,
				params.OriginDutyLocationZip,
				destinationZip,
				distanceMilesFromOriginPickupZip,
				distanceMilesFromOriginDutyLocationZip,
				time.Time(params.OriginalMoveDate),
				0, // We don't want any SIT charges
			)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			cost := rateengine.GetWinningCostMove(costDetails)

			gcc := cost.GCC
			incentivePercentage := cost.GCC.MultiplyFloat64(0.95)

			ppmObligation := internalmessages.PPMIncentive{
				Gcc:                 swag.Int64(gcc.Int64()),
				IncentivePercentage: swag.Int64(incentivePercentage.Int64()),
			}
			return ppmop.NewShowPPMIncentiveOK().WithPayload(&ppmObligation)
		})
}
