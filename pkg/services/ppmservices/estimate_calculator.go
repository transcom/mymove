package ppmservices

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type estimateCalculator struct {
	db      *pop.Connection
	logger  Logger
	planner route.Planner
}

// NewEstimateCalculator returns a new estimateCalculator
func NewEstimateCalculator(db *pop.Connection, logger Logger, planner route.Planner) services.EstimateCalculator {
	return &estimateCalculator{db: db, logger: logger, planner: planner}
}

func (e *estimateCalculator) CalculateEstimate(ppm *models.PersonallyProcuredMove, moveID uuid.UUID) error {
	move, err := models.FetchMoveByMoveID(e.db, moveID)
	if err != nil {
		if err == models.ErrFetchNotFound {
			return services.NewNotFoundError(moveID, "Unable to calculate estimate")
		}
		return fmt.Errorf("error calculating estimate: unable to fetch move with ID %s: %w", moveID, err)
	}

	re := rateengine.NewRateEngine(e.db, e.logger, move)
	daysInSIT := 0
	if ppm.HasSit != nil && *ppm.HasSit && ppm.DaysInStorage != nil {
		daysInSIT = int(*ppm.DaysInStorage)
	}

	originDutyStationZip := ppm.Move.Orders.ServiceMember.DutyStation.Address.PostalCode
	destinationDutyStationZip := ppm.Move.Orders.NewDutyStation.Address.PostalCode

	distanceMilesFromOriginPickupZip, err := e.planner.Zip5TransitDistance(*ppm.PickupPostalCode, destinationDutyStationZip)
	if err != nil {
		return fmt.Errorf("error calculating estimate: cannot get distance from origin pickup to destination: %w", err)
	}

	distanceMilesFromOriginDutyStationZip, err := e.planner.Zip5TransitDistance(originDutyStationZip, destinationDutyStationZip)
	if err != nil {
		return fmt.Errorf("error calculating estimate: cannot get distance from origin duty station to destination: %w", err)
	}

	costDetails, err := re.ComputePPMMoveCosts(
		unit.Pound(*ppm.WeightEstimate),
		*ppm.PickupPostalCode,
		originDutyStationZip,
		destinationDutyStationZip,
		distanceMilesFromOriginPickupZip,
		distanceMilesFromOriginDutyStationZip,
		time.Time(*ppm.OriginalMoveDate),
		daysInSIT,
	)
	if err != nil {
		return fmt.Errorf("error calculating estimate: cannot compute PPM move costs: %w", err)
	}

	cost := rateengine.GetWinningCostMove(costDetails)

	// Update SIT estimate
	if ppm.HasSit != nil && *ppm.HasSit {
		cwtWeight := unit.Pound(*ppm.WeightEstimate).ToCWT()
		sitZip3 := rateengine.Zip5ToZip3(destinationDutyStationZip)
		sitComputation, sitChargeErr := re.SitCharge(cwtWeight, daysInSIT, sitZip3, *ppm.OriginalMoveDate, true)
		if sitChargeErr != nil {
			return fmt.Errorf("error calculating estimate: cannot calculate SIT charge: %w", sitChargeErr)
		}
		sitCharge := float64(sitComputation.ApplyDiscount(cost.LHDiscount, cost.SITDiscount))
		reimbursementString := fmt.Sprintf("$%.2f", sitCharge/100)
		ppm.EstimatedStorageReimbursement = &reimbursementString
	}

	mileage := int64(cost.LinehaulCostComputation.Mileage)
	ppm.Mileage = &mileage
	ppm.PlannedSITMax = &cost.SITFee
	ppm.SITMax = &cost.SITMax
	min := cost.GCC.MultiplyFloat64(0.95)
	max := cost.GCC.MultiplyFloat64(1.05)
	ppm.IncentiveEstimateMin = &min
	ppm.IncentiveEstimateMax = &max

	return nil
}