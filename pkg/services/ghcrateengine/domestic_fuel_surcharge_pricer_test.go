package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticFuelSurcharge() {
	moveDate := time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC)
	planner := &mocks.Planner{}
	sourceZip := "00001"
	destinationZip := "90210"
	weight := unit.Pound(3000)

	//Returns an error for now. Once the PriceDomesticFuelSurcharge implementation is complete, this test will check for the correct fuel surcharge value.
	suite.Run("PriceDomesticFuelSurcharge when weight is less than 5000", func() {
		domesticFuelSurchargePricer := NewDomesticFuelSurchargePricer(suite.DB(), suite.logger, testdatagen.DefaultContractCode)
		fuelSurcharge, err := domesticFuelSurchargePricer.PriceDomesticFuelSurcharge(
			moveDate,
			planner,
			weight,
			sourceZip,
			destinationZip,
		)
		suite.Error(err)
		suite.Equal(err.Error(), "Error calculating fuel surcharge")
		suite.Equal(fuelSurcharge, unit.Cents(0))
	})
}
