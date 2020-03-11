package paperwork

import (
	"fmt"
	"errors"
	"time"

	"github.com/transcom/mymove/pkg/route"

	"github.com/transcom/mymove/pkg/rateengine"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"
)

type ppmComputer interface {
	ComputePPMMoveCosts(weight unit.Pound, originPickupZip5 string, originDutyStationZip5 string, destinationZip5 string, distanceMilesFromOriginPickupZip int, distanceMilesFromOriginDutyStationZip int, date time.Time, daysInSit int) (cost rateengine.CostDetails, err error)
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
func (sswPpmComputer *SSWPPMComputer) ComputeObligations(ssfd models.ShipmentSummaryFormData, planner route.Planner) (obligation models.Obligations, err error) {
	firstPPM, err := sswPpmComputer.nilCheckPPM(ssfd)
	if err != nil {
		return models.Obligations{}, err
	}

	originDutyStationZip := ssfd.CurrentDutyStation.Address.PostalCode
	destDutyStationZip := ssfd.Order.NewDutyStation.Address.PostalCode

	distanceMilesFromPickupZip, err := planner.Zip5TransitDistance(*firstPPM.PickupPostalCode, destDutyStationZip)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating distance")
	}

	distanceMilesFromDutyStationZip, err := planner.Zip5TransitDistance(originDutyStationZip, destDutyStationZip)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating distance")
	}

	actualCost, err := sswPpmComputer.ComputePPMMoveCosts(
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

	// mileageWon := unit.Miles(actualCost.Mileage)

	maxCost, err := sswPpmComputer.ComputePPMMoveCosts(
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

	var actualSIT unit.Cents
	if firstPPM.TotalSITCost != nil {
		actualSIT = *firstPPM.TotalSITCost
	}
	// This logic needs to be put back in!!
	// if actualSIT > maxCost.SITMax {
	// 	actualSIT = maxCost.SITMax
	// }

	var lowestActualObligation models.Obligation
	var actualObligation models.Obligation
	var maxObligation models.Obligation
	var lowestMaxObligation models.Obligation
	if actualCost["pickupLocation"].IsLowest {
		lowestActualObligation = models.Obligation{Gcc: actualCost["pickupLocation"].Cost.GCC, SIT: actualSIT, Miles: unit.Miles(actualCost["pickupLocation"].Cost.Mileage)}
		actualObligation = models.Obligation{Gcc: actualCost["originDutyStation"].Cost.GCC, SIT: actualSIT, Miles: unit.Miles(actualCost["originDutyStation"].Cost.Mileage)}
	} else {
		actualObligation = models.Obligation{Gcc: actualCost["pickupLocation"].Cost.GCC, SIT: actualSIT, Miles: unit.Miles(actualCost["pickupLocation"].Cost.Mileage)}
		lowestActualObligation = models.Obligation{Gcc: actualCost["originDutyStation"].Cost.GCC, SIT: actualSIT, Miles: unit.Miles(actualCost["originDutyStation"].Cost.Mileage)}
	}

	if maxCost["pickupLocation"].IsLowest {
		lowestMaxObligation = models.Obligation{Gcc: maxCost["pickupLocation"].Cost.GCC, SIT: actualSIT, Miles: unit.Miles(maxCost["pickupLocation"].Cost.Mileage)}
		maxObligation = models.Obligation{Gcc: maxCost["originDutyStation"].Cost.GCC, SIT: actualSIT, Miles: unit.Miles(maxCost["originDutyStation"].Cost.Mileage)}
	} else {
		maxObligation = models.Obligation{Gcc: maxCost["pickupLocation"].Cost.GCC, SIT: actualSIT, Miles: unit.Miles(maxCost["pickupLocation"].Cost.Mileage)}
		lowestMaxObligation = models.Obligation{Gcc: maxCost["originDutyStation"].Cost.GCC, SIT: actualSIT, Miles: unit.Miles(maxCost["originDutyStation"].Cost.Mileage)}
	}

	obligations := models.Obligations{LowestMaxObligation: lowestMaxObligation, LowestActualObligation: lowestActualObligation, MaxObligation: maxObligation, ActualObligation: actualObligation}
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
