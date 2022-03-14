package paperwork

import (
	"errors"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/route"

	"github.com/transcom/mymove/pkg/rateengine"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"
)

type ppmComputer interface {
	ComputePPMMoveCosts(appCtx appcontext.AppContext, weight unit.Pound, originPickupZip5 string, originDutyStationZip5 string, destinationZip5 string, distanceMilesFromOriginPickupZip int, distanceMilesFromOriginDutyStationZip int, date time.Time, daysInSit int) (cost rateengine.CostDetails, err error)
}

//SSWPPMComputer a rate engine wrapper with helper functions to simplify ppm cost calculations specific to shipment summary worksheet
type SSWPPMComputer struct {
	ppmComputer
}

//NewSSWPPMComputer creates a SSWPPMComputer
func NewSSWPPMComputer(PPMComputer ppmComputer) *SSWPPMComputer {
	return &SSWPPMComputer{ppmComputer: PPMComputer}
}

//ObligationType type corresponding to obligation sections of shipment summary worksheet
type ObligationType int

//ComputeObligations is helper function for computing the obligations section of the shipment summary worksheet
func (sswPpmComputer *SSWPPMComputer) ComputeObligations(appCtx appcontext.AppContext, ssfd models.ShipmentSummaryFormData, planner route.Planner) (obligation models.Obligations, err error) {
	firstPPM, err := sswPpmComputer.nilCheckPPM(ssfd)
	if err != nil {
		return models.Obligations{}, err
	}

	originDutyStationZip := ssfd.CurrentDutyLocation.Address.PostalCode
	destDutyStationZip := ssfd.Order.NewDutyLocation.Address.PostalCode

	distanceMilesFromPickupZip, err := planner.Zip5TransitDistanceLineHaul(appCtx, *firstPPM.PickupPostalCode, destDutyStationZip)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating distance")
	}

	distanceMilesFromDutyStationZip, err := planner.Zip5TransitDistanceLineHaul(appCtx, originDutyStationZip, destDutyStationZip)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating distance")
	}

	actualCosts, err := sswPpmComputer.ComputePPMMoveCosts(
		appCtx,
		ssfd.PPMRemainingEntitlement,
		*firstPPM.PickupPostalCode,
		originDutyStationZip,
		destDutyStationZip,
		distanceMilesFromPickupZip,
		distanceMilesFromDutyStationZip,
		*firstPPM.OriginalMoveDate,
		0,
	)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating PPM actual obligations")
	}

	maxCosts, err := sswPpmComputer.ComputePPMMoveCosts(
		appCtx,
		ssfd.WeightAllotment.TotalWeight,
		*firstPPM.PickupPostalCode,
		originDutyStationZip,
		destDutyStationZip,
		distanceMilesFromPickupZip,
		distanceMilesFromDutyStationZip,
		*firstPPM.OriginalMoveDate,
		0,
	)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating PPM max obligations")
	}

	actualCost := rateengine.GetWinningCostMove(actualCosts)
	maxCost := rateengine.GetWinningCostMove(maxCosts)
	nonWinningActualCost := rateengine.GetNonWinningCostMove(actualCosts)
	nonWinningMaxCost := rateengine.GetNonWinningCostMove(maxCosts)

	var actualSIT unit.Cents
	if firstPPM.TotalSITCost != nil {
		actualSIT = *firstPPM.TotalSITCost
	}

	if actualSIT > maxCost.SITMax {
		actualSIT = maxCost.SITMax
	}

	obligations := models.Obligations{
		ActualObligation:           models.Obligation{Gcc: actualCost.GCC, SIT: actualSIT, Miles: unit.Miles(actualCost.Mileage)},
		MaxObligation:              models.Obligation{Gcc: maxCost.GCC, SIT: actualSIT, Miles: unit.Miles(actualCost.Mileage)},
		NonWinningActualObligation: models.Obligation{Gcc: nonWinningActualCost.GCC, SIT: actualSIT, Miles: unit.Miles(nonWinningActualCost.Mileage)},
		NonWinningMaxObligation:    models.Obligation{Gcc: nonWinningMaxCost.GCC, SIT: actualSIT, Miles: unit.Miles(nonWinningActualCost.Mileage)},
	}
	return obligations, nil
}

func (sswPpmComputer *SSWPPMComputer) nilCheckPPM(ssfd models.ShipmentSummaryFormData) (models.PersonallyProcuredMove, error) {
	if len(ssfd.PersonallyProcuredMoves) == 0 {
		return models.PersonallyProcuredMove{}, errors.New("missing ppm")
	}
	firstPPM := ssfd.PersonallyProcuredMoves[0]
	if firstPPM.PickupPostalCode == nil || firstPPM.DestinationPostalCode == nil {
		return models.PersonallyProcuredMove{}, errors.New("missing required address parameter")
	}
	if firstPPM.OriginalMoveDate == nil {
		return models.PersonallyProcuredMove{}, errors.New("missing required original move date parameter")
	}
	return firstPPM, nil
}
