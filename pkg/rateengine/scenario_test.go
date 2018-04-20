package rateengine

import (
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) Test_Scenario1() {
	scenario.RunRateEngineScenario1(suite.db)

	logger, err := zap.NewDevelopment()
	suite.Assertions.Nil(err, "could not create a development logger")

	planner := route.NewTestingPlanner(362)
	engine := NewRateEngine(suite.db, logger, planner)

	weight := unit.Pound(4000)
	originZip5 := "32168"
	destinationZip5 := "29429"
	date := time.Date(2018, time.June, 18, 0, 0, 0, 0, time.UTC)
	inverseDiscount := 0.33

	cost, err := engine.computePPM(weight, originZip5, destinationZip5, date, inverseDiscount)
	suite.Assertions.Nil(err, "could not compute PPM")

	suite.Assertions.Equal(unit.Cents(163434), cost.LinehaulChargeTotal)
	suite.Assertions.Equal(unit.Cents(4765), cost.OriginServiceFee)
	suite.Assertions.Equal(unit.Cents(5689), cost.DestinationServiceFee)
	suite.Assertions.Equal(unit.Cents(89412), cost.FullPackUnpackFee)
	suite.Assertions.Equal(unit.Cents(263300), cost.GCC)
}
