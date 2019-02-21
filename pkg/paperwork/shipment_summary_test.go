package paperwork

import (
	"errors"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

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
		Date:            date,
		DaysInSIT:       daysInSIT,
		SitDiscount:     sitDiscount,
	}
	return mppmc.costComputation, mppmc.err
}

func (mppmc *mockPPMComputer) CalledWith() ppmComputerParams {
	return mppmc.ppmComputerParams
}

func (suite *PaperworkSuite) TestComputeObligationsBadPPM() {
	ppmComputer := NewSSWPPMComputer(&mockPPMComputer{})
	badPpm := models.PersonallyProcuredMove{}
	_, err1 := ppmComputer.ComputeObligations(badPpm, 0, MaxObligation)
	_, err2 := ppmComputer.ComputeObligations(badPpm, 0, ActualObligation)
	suite.NotNil(err1)
	suite.NotNil(err2)
}

func (suite *PaperworkSuite) TestTestComputeObligations() {
	mockComputer := mockPPMComputer{
		costComputation: rateengine.CostComputation{GCC: 100},
	}
	ppmComputer := NewSSWPPMComputer(&mockComputer)
	var wtgEstimate int64 = 1000
	var netWeight int64 = 2000
	origMoveDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	actualDate := time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC)
	pickupPostalCode := "85369"
	destinationPostalCode := "31905"
	ppm := models.PersonallyProcuredMove{
		WeightEstimate:        &wtgEstimate,
		NetWeight:             &netWeight,
		OriginalMoveDate:      &origMoveDate,
		ActualMoveDate:        &actualDate,
		PickupPostalCode:      &pickupPostalCode,
		DestinationPostalCode: &destinationPostalCode,
	}
	suite.Run("TestComputeMaxObligations", func() {

		_, err := ppmComputer.ComputeObligations(ppm, 0, MaxObligation)

		calledWith := mockComputer.CalledWith()
		expectMaxObligationParams := ppmComputerParams{
			Weight:          unit.Pound(wtgEstimate),
			OriginZip5:      pickupPostalCode,
			DestinationZip5: destinationPostalCode,
			Date:            origMoveDate,
			DaysInSIT:       0,
			SitDiscount:     0,
		}
		suite.Equal(calledWith, expectMaxObligationParams)
		suite.Nil(err)
	})

	suite.Run("TestComputeActualObligations", func() {

		expectActualObligationParams := ppmComputerParams{
			Weight:          unit.Pound(netWeight),
			OriginZip5:      pickupPostalCode,
			DestinationZip5: destinationPostalCode,
			Date:            actualDate,
			DaysInSIT:       0,
			SitDiscount:     0,
		}
		_, err := ppmComputer.ComputeObligations(ppm, 0, ActualObligation)

		calledWith := mockComputer.CalledWith()
		suite.Equal(calledWith, expectActualObligationParams)
		suite.Nil(err)
	})

	suite.Run("TestCalcError", func() {
		mockComputer := mockPPMComputer{err: errors.New("ERROR")}
		ppmComputer := SSWPPMComputer{&mockComputer}

		_, err := ppmComputer.ComputeObligations(ppm, 0, ActualObligation)

		suite.NotNil(err)
	})
}
