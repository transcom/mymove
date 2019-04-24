package paperwork

import (
	"errors"
	"time"

	"github.com/transcom/mymove/pkg/route"

	"github.com/transcom/mymove/pkg/rateengine"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"
)

type ppmComputer interface {
	ComputePPMIncludingLHDiscount(weight unit.Pound, originZip5 string, destinationZip5 string, distanceMiles int, date time.Time, daysInSIT int) (cost rateengine.CostComputation, err error)
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
	distanceMiles, err := planner.Zip5TransitDistance(*firstPPM.PickupPostalCode, *firstPPM.DestinationPostalCode)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating distance")
	}
	maxCost, err := sswPpmComputer.ComputePPMIncludingLHDiscount(
		ssfd.WeightAllotment.TotalWeight,
		*firstPPM.PickupPostalCode,
		*firstPPM.DestinationPostalCode,
		distanceMiles,
		*firstPPM.ActualMoveDate,
		0,
	)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating PPM max obligations")
	}
	actualCost, err := sswPpmComputer.ComputePPMIncludingLHDiscount(
		ssfd.PPMRemainingEntitlement,
		*firstPPM.PickupPostalCode,
		*firstPPM.DestinationPostalCode,
		distanceMiles,
		*firstPPM.ActualMoveDate,
		0,
	)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating PPM actual obligations")
	}
	var actualSIT unit.Cents
	if firstPPM.TotalSITCost != nil {
		actualSIT = *firstPPM.TotalSITCost
	}
	if actualSIT > maxCost.SITMax {
		actualSIT = maxCost.SITMax
	}
	maxObligation := models.Obligation{Gcc: maxCost.GCC, SIT: maxCost.SITMax}
	actualObligation := models.Obligation{Gcc: actualCost.GCC, SIT: actualSIT}
	obligations := models.Obligations{MaxObligation: maxObligation, ActualObligation: actualObligation}
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
	if firstPPM.ActualMoveDate == nil {
		return models.PersonallyProcuredMove{}, errors.New("missing required actual move date parameter")
	}
	return firstPPM, nil
}
