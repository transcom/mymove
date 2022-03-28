package ppmservices

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type estimateCalculator struct {
	planner route.Planner
}

// NewEstimateCalculator returns a new estimateCalculator
func NewEstimateCalculator(planner route.Planner) services.EstimateCalculator {
	return &estimateCalculator{planner: planner}
}

func (e *estimateCalculator) CalculateEstimates(appCtx appcontext.AppContext, ppm *models.PersonallyProcuredMove, moveID uuid.UUID) (int64, rateengine.CostComputation, error) {
	var sitCharge int64
	cost := rateengine.CostComputation{}
	move, err := models.FetchMoveByMoveID(appCtx.DB(), moveID)
	if err != nil {
		switch err {
		case models.ErrFetchNotFound:
			return sitCharge, cost, apperror.NewNotFoundError(moveID, "Unable to calculate estimate")
		default:
			return sitCharge, cost, apperror.NewQueryError("Move", err, fmt.Sprintf("error calculating estimate: unable to fetch move with ID %s", moveID))
		}
	}

	re := rateengine.NewRateEngine(move)
	daysInSIT := 0
	if ppm.HasSit != nil && *ppm.HasSit && ppm.DaysInStorage != nil {
		daysInSIT = int(*ppm.DaysInStorage)
	}

	originDutyLocationZip := ppm.Move.Orders.ServiceMember.DutyLocation.Address.PostalCode
	destinationDutyLocationZip := ppm.Move.Orders.NewDutyLocation.Address.PostalCode

	distanceMilesFromOriginPickupZip, err := e.planner.Zip5TransitDistanceLineHaul(appCtx, *ppm.PickupPostalCode, destinationDutyLocationZip)
	if err != nil {
		return sitCharge, cost, fmt.Errorf("error calculating estimate: cannot get distance from origin pickup to destination: %w", err)
	}

	distanceMilesFromOriginDutyLocationZip, err := e.planner.Zip5TransitDistanceLineHaul(appCtx, originDutyLocationZip, destinationDutyLocationZip)
	if err != nil {
		return sitCharge, cost, fmt.Errorf("error calculating estimate: cannot get distance from origin duty location to destination: %w", err)
	}

	costDetails, err := re.ComputePPMMoveCosts(
		appCtx,
		*ppm.WeightEstimate,
		*ppm.PickupPostalCode,
		originDutyLocationZip,
		destinationDutyLocationZip,
		distanceMilesFromOriginPickupZip,
		distanceMilesFromOriginDutyLocationZip,
		*ppm.OriginalMoveDate,
		daysInSIT,
	)
	if err != nil {
		return sitCharge, cost, fmt.Errorf("error calculating estimate: cannot compute PPM move costs: %w", err)
	}

	// get the SIT charge
	cost = rateengine.GetWinningCostMove(costDetails)
	cwtWeight := unit.Pound(*ppm.WeightEstimate).ToCWT()
	sitZip3 := rateengine.Zip5ToZip3(destinationDutyLocationZip)
	if !*ppm.HasSit {
		return sitCharge, cost, nil
	}
	sitComputation, sitChargeErr := re.SitCharge(appCtx, cwtWeight, daysInSIT, sitZip3, *ppm.OriginalMoveDate, true)
	if sitChargeErr != nil {
		return sitCharge, cost, sitChargeErr
	}
	sitCharge = int64(sitComputation.ApplyDiscount(cost.LHDiscount, cost.SITDiscount))

	return sitCharge, cost, nil
}
