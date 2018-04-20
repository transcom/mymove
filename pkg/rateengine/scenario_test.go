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

	cost, err := engine.ComputePPM(weight, originZip5, destinationZip5, date, inverseDiscount)
	suite.Assertions.Nil(err, "could not compute PPM")

	suite.Equal(unit.Cents(163434), cost.LinehaulChargeTotal)
	suite.Equal(unit.Cents(4765), cost.OriginServiceFee)
	suite.Equal(unit.Cents(5689), cost.DestinationServiceFee)
	suite.Equal(unit.Cents(89412), cost.FullPackUnpackFee)
	suite.Equal(unit.Cents(263300), cost.GCC)
}

func (suite *RateEngineSuite) Test_Scenario2() {
	scenario.RunRateEngineScenario2(suite.db)

	logger, err := zap.NewDevelopment()
	suite.Assertions.Nil(err, "could not create a development logger")

	planner := route.NewTestingPlanner(1693)
	engine := NewRateEngine(suite.db, logger, planner)

	weight := unit.Pound(7500)
	originZip5 := "94540"
	destinationZip5 := "78626"
	date := time.Date(2018, time.December, 5, 0, 0, 0, 0, time.UTC)
	inverseDiscount := 0.33

	cost, err := engine.ComputePPM(weight, originZip5, destinationZip5, date, inverseDiscount)
	suite.Assertions.Nil(err, "could not compute PPM")

	suite.Equal(unit.Cents(421707), cost.LinehaulChargeTotal)
	suite.Equal(unit.Cents(12103), cost.OriginServiceFee)
	suite.Equal(unit.Cents(11187), cost.DestinationServiceFee)
	suite.Equal(unit.Cents(183620), cost.FullPackUnpackFee)
	suite.Equal(unit.Cents(637056), cost.GCC)
}
