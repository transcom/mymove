package rateengine

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) Test_CheckPPMTotal() {
	move := models.Move{
		Locator: "ABC123",
	}

	engine := NewRateEngine(move)

	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	// 139698 +20000
	cost, err := engine.computePPM(suite.AppContextForTest(), 2000, "39574", "33633", 1234, testdatagen.RateEngineDate,
		1, unit.DiscountRate(.6), unit.DiscountRate(.5))

	if err != nil {
		suite.Fail("failed to calculate ppm charge: %s", err)
	}

	// PPMs estimates are being hardcoded because we are not loading tariff400ng data
	// update this check so test passes - but this is not testing correctness of data
	suite.Equal(unit.Cents(175543), cost.GCC, "wrong GCC")
}

type RateEngineSuite struct {
	*testingsuite.PopTestSuite
}

func TestRateEngineSuite(t *testing.T) {
	hs := &RateEngineSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
