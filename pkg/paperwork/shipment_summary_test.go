package paperwork

import (
	"errors"
	"time"

	"github.com/transcom/mymove/pkg/route"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

type ppmComputerParams struct {
	Weight          unit.Pound
	OriginZip5      string
	DestinationZip5 string
	Miles           int
	Date            time.Time
	DaysInSIT       int
	SitDiscount     unit.DiscountRate
}

type mockPPMComputer struct {
	costComputation   rateengine.CostComputation
	err               error
	ppmComputerParams ppmComputerParams
}

func (mppmc *mockPPMComputer) ComputePPMIncludingLHDiscount(weight unit.Pound, originZip5 string, destinationZip5 string, distanceMiles int, date time.Time, daysInSIT int, sitDiscount unit.DiscountRate) (cost rateengine.CostComputation, err error) {
	mppmc.ppmComputerParams = ppmComputerParams{
		Weight:          weight,
		OriginZip5:      originZip5,
		DestinationZip5: destinationZip5,
		Miles:           distanceMiles,
		Date:            date,
		DaysInSIT:       daysInSIT,
		SitDiscount:     sitDiscount,
	}
	return mppmc.costComputation, mppmc.err
}

func (mppmc *mockPPMComputer) CalledWith() ppmComputerParams {
	return mppmc.ppmComputerParams
}

func (suite *PaperworkSuite) TestComputeObligationsParams() {
	ppmComputer := NewSSWPPMComputer(&mockPPMComputer{})
	pickupPostalCode := "85369"
	destinationPostalCode := "31905"
	ppm := models.PersonallyProcuredMove{
		PickupPostalCode:      &pickupPostalCode,
		DestinationPostalCode: &destinationPostalCode,
	}
	noPPM := models.ShipmentSummaryFormData{PersonallyProcuredMoves: models.PersonallyProcuredMoves{}}
	missingZip := models.ShipmentSummaryFormData{PersonallyProcuredMoves: models.PersonallyProcuredMoves{{}}}
	missingOrigMoveDate := models.ShipmentSummaryFormData{PersonallyProcuredMoves: models.PersonallyProcuredMoves{ppm}}
	missingActualMoveDate := models.ShipmentSummaryFormData{PersonallyProcuredMoves: models.PersonallyProcuredMoves{ppm}}

	_, err1 := ppmComputer.ComputeObligations(noPPM, route.NewTestingPlanner(10), MaxObligation)
	_, err2 := ppmComputer.ComputeObligations(missingZip, route.NewTestingPlanner(10), MaxObligation)
	_, err3 := ppmComputer.ComputeObligations(missingActualMoveDate, route.NewTestingPlanner(10), ActualObligation)
	_, err4 := ppmComputer.ComputeObligations(missingOrigMoveDate, route.NewTestingPlanner(10), MaxObligation)

	suite.NotNil(err1)
	suite.NotNil(err2)
	suite.NotNil(err3)
	suite.NotNil(err4)
}

func (suite *PaperworkSuite) TestTestComputeObligations() {
	var netWeight int64 = 2000
	miles := 100
	maxObligation := 1000
	planner := route.NewTestingPlanner(miles)
	origMoveDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	actualDate := time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC)
	pickupPostalCode := "85369"
	destinationPostalCode := "31905"
	ppm := models.PersonallyProcuredMove{
		NetWeight:             &netWeight,
		OriginalMoveDate:      &origMoveDate,
		ActualMoveDate:        &actualDate,
		PickupPostalCode:      &pickupPostalCode,
		DestinationPostalCode: &destinationPostalCode,
	}
	params := models.ShipmentSummaryFormData{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{ppm},
		TotalWeightAllotment:    maxObligation,
	}
	mockComputer := mockPPMComputer{
		costComputation: rateengine.CostComputation{GCC: 100},
	}
	ppmComputer := NewSSWPPMComputer(&mockComputer)

	suite.Run("TestComputeMaxObligations", func() {
		expectMaxObligationParams := ppmComputerParams{
			Weight:          unit.Pound(maxObligation),
			OriginZip5:      pickupPostalCode,
			DestinationZip5: destinationPostalCode,
			Miles:           miles,
			Date:            origMoveDate,
			DaysInSIT:       0,
			SitDiscount:     0,
		}

		_, err := ppmComputer.ComputeObligations(params, planner, MaxObligation)

		calledWith := mockComputer.CalledWith()
		suite.Equal(calledWith, expectMaxObligationParams)
		suite.Nil(err)
	})

	suite.Run("TestComputeActualObligations", func() {
		expectActualObligationParams := ppmComputerParams{
			Weight:          unit.Pound(netWeight),
			OriginZip5:      pickupPostalCode,
			DestinationZip5: destinationPostalCode,
			Date:            actualDate,
			Miles:           miles,
			DaysInSIT:       0,
			SitDiscount:     0,
		}

		_, err := ppmComputer.ComputeObligations(params, planner, ActualObligation)

		calledWith := mockComputer.CalledWith()
		suite.Equal(calledWith, expectActualObligationParams)
		suite.Nil(err)
	})

	suite.Run("TestCalcError", func() {
		mockComputer := mockPPMComputer{err: errors.New("ERROR")}
		ppmComputer := SSWPPMComputer{&mockComputer}

		_, err := ppmComputer.ComputeObligations(params, planner, ActualObligation)

		suite.NotNil(err)
	})
}
