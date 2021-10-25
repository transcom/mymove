package ppmservices

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type estimateCalculator struct {
	db *pop.Connection
	//TODO: add logger back once we are able to pass a logger in through internalapi (see https://dp3.atlassian.net/browse/MB-2352)
	//logger  Logger
	planner route.Planner
}

// NewEstimateCalculator returns a new estimateCalculator
func NewEstimateCalculator(db *pop.Connection, planner route.Planner) services.EstimateCalculator {
	return &estimateCalculator{db: db, planner: planner}
}

func (e *estimateCalculator) CalculateEstimates(ppm *models.PersonallyProcuredMove, moveID uuid.UUID, logger *zap.Logger) (int64, rateengine.CostComputation, error) {
	// temporarily passing in logger here until fix listed above in service struct
	var sitCharge int64
	cost := rateengine.CostComputation{}
	move, err := models.FetchMoveByMoveID(e.db, moveID)
	if err != nil {
		switch err {
		case models.ErrFetchNotFound:
			return sitCharge, cost, apperror.NewNotFoundError(moveID, "Unable to calculate estimate")
		default:
			return sitCharge, cost, apperror.NewQueryError("Move", err, fmt.Sprintf("error calculating estimate: unable to fetch move with ID %s", moveID))
		}
	}

	re := rateengine.NewRateEngine(e.db, logger, move)
	daysInSIT := 0
	if ppm.HasSit != nil && *ppm.HasSit && ppm.DaysInStorage != nil {
		daysInSIT = int(*ppm.DaysInStorage)
	}

	originDutyStationZip := ppm.Move.Orders.ServiceMember.DutyStation.Address.PostalCode
	destinationDutyStationZip := ppm.Move.Orders.NewDutyStation.Address.PostalCode

	distanceMilesFromOriginPickupZip, err := e.planner.Zip5TransitDistanceLineHaul(*ppm.PickupPostalCode, destinationDutyStationZip)
	if err != nil {
		return sitCharge, cost, fmt.Errorf("error calculating estimate: cannot get distance from origin pickup to destination: %w", err)
	}

	distanceMilesFromOriginDutyStationZip, err := e.planner.Zip5TransitDistanceLineHaul(originDutyStationZip, destinationDutyStationZip)
	if err != nil {
		return sitCharge, cost, fmt.Errorf("error calculating estimate: cannot get distance from origin duty station to destination: %w", err)
	}

	costDetails, err := re.ComputePPMMoveCosts(
		*ppm.WeightEstimate,
		*ppm.PickupPostalCode,
		originDutyStationZip,
		destinationDutyStationZip,
		distanceMilesFromOriginPickupZip,
		distanceMilesFromOriginDutyStationZip,
		*ppm.OriginalMoveDate,
		daysInSIT,
	)
	if err != nil {
		return sitCharge, cost, fmt.Errorf("error calculating estimate: cannot compute PPM move costs: %w", err)
	}

	// get the SIT charge
	cost = rateengine.GetWinningCostMove(costDetails)
	cwtWeight := unit.Pound(*ppm.WeightEstimate).ToCWT()
	sitZip3 := rateengine.Zip5ToZip3(destinationDutyStationZip)
	if !*ppm.HasSit {
		return sitCharge, cost, nil
	}
	sitComputation, sitChargeErr := re.SitCharge(cwtWeight, daysInSIT, sitZip3, *ppm.OriginalMoveDate, true)
	if sitChargeErr != nil {
		return sitCharge, cost, sitChargeErr
	}
	sitCharge = int64(sitComputation.ApplyDiscount(cost.LHDiscount, cost.SITDiscount))

	return sitCharge, cost, nil
}
