package ghcrateengine

import (
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticFuelSurcharge() {
	planner := route.NewTestingPlanner(1000)
	sourceZip := "00001"
	destinationZip := "90210"
	weight := unit.Pound(3000)

	suite.Run("PriceDomesticFuelSurcharge when weight is less than 5000", func() {
		domesticFuelSurchargePricer := NewDomesticFuelSurchargePricer(suite.DB(), suite.logger, testdatagen.DefaultContractCode)
		fuelSurcharge, err := domesticFuelSurchargePricer.PriceDomesticFuelSurcharge(
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
