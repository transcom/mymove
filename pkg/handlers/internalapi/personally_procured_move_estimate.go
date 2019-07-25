package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/gofrs/uuid"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// ShowPPMEstimateHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMEstimateHandler struct {
	handlers.HandlerContext
}

// Handle calculates a PPM reimbursement range.
func (h ShowPPMEstimateHandler) Handle(params ppmop.ShowPPMEstimateParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	ctx := params.HTTPRequest.Context()

	engine := rateengine.NewRateEngine(h.DB(), logger)

	serviceMemberID, _ := uuid.FromString(session.ServiceMemberID.String())
	serviceMember, err := models.FetchServiceMemberForUser(ctx, h.DB(), session, serviceMemberID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	dutyStationZip := serviceMember.DutyStation.Address.PostalCode

	lhDiscountPickupZip, _, err := models.PPMDiscountFetch(h.DB(),
		logger,
		params.OriginZip,
		params.DestinationZip,
		time.Time(params.OriginalMoveDate),
	)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	lhDiscountDutyStationZip, _, err := models.PPMDiscountFetch(h.DB(),
		logger,
		dutyStationZip,
		params.DestinationZip,
		time.Time(params.OriginalMoveDate),
	)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	distanceMilesFromPickupZip, err := h.Planner().Zip5TransitDistance(params.OriginZip, params.DestinationZip)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	distanceMilesFromDutyStationZip, err := h.Planner().Zip5TransitDistance(dutyStationZip, params.DestinationZip)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	cost, err := engine.ComputePPM(
		unit.Pound(params.WeightEstimate),
		params.OriginZip,
		dutyStationZip,
		params.DestinationZip,
		distanceMilesFromPickupZip,
		distanceMilesFromDutyStationZip,
		time.Time(params.OriginalMoveDate),
		0, // We don't want any SIT charges
		lhDiscountPickupZip,
		lhDiscountDutyStationZip,
		0.0,
	)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	min := cost.GCC.MultiplyFloat64(0.95)
	max := cost.GCC.MultiplyFloat64(1.05)

	ppmEstimate := internalmessages.PPMEstimateRange{
		RangeMin: swag.Int64(min.Int64()),
		RangeMax: swag.Int64(max.Int64()),
	}
	return ppmop.NewShowPPMEstimateOK().WithPayload(&ppmEstimate)
}
