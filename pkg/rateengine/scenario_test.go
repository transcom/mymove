package rateengine

import (
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) Test_Scenario1() {
	err := scenario.RunRateEngineScenario1(suite.db)
	suite.Nil(err, "failed to run scenario 1")

	logger, err := zap.NewDevelopment()
	suite.Assertions.Nil(err, "could not create a development logger")

	planner := route.NewTestingPlanner(362)
	engine := NewRateEngine(suite.db, logger, planner)

	weight := unit.Pound(4000)
	originZip5 := "32168"
	destinationZip5 := "29429"
	date := time.Date(2018, time.June, 18, 0, 0, 0, 0, time.UTC)
	lhDiscount := unit.DiscountRate(0.67)

	cost, err := engine.ComputePPM(weight, originZip5, destinationZip5, date, 0, lhDiscount, 0)
	suite.Assertions.Nil(err, "could not compute PPM")

	suite.Equal(unit.Cents(163434), cost.LinehaulChargeTotal)
	suite.Equal(unit.Cents(4765), cost.OriginServiceFee)
	suite.Equal(unit.Cents(5689), cost.DestinationServiceFee)
	suite.Equal(unit.Cents(263300), cost.GCC)
}

func (suite *RateEngineSuite) Test_Scenario2() {
	if err := scenario.RunRateEngineScenario2(suite.db); err != nil {
		suite.NotNil(err, "failed to run scenario 2")
	}

	logger, err := zap.NewDevelopment()
	suite.Assertions.Nil(err, "could not create a development logger")

	planner := route.NewTestingPlanner(1693)
	engine := NewRateEngine(suite.db, logger, planner)

	weight := unit.Pound(7500)
	originZip5 := "94540"
	destinationZip5 := "78626"
	date := time.Date(2018, time.December, 5, 0, 0, 0, 0, time.UTC)
	lhDiscount := unit.DiscountRate(0.67)

	cost, err := engine.ComputePPM(weight, originZip5, destinationZip5, date, 0, lhDiscount, 0)
	suite.Assertions.Nil(err, "could not compute PPM")

	suite.Equal(unit.Cents(430147), cost.LinehaulChargeTotal)
	suite.Equal(unit.Cents(12103), cost.OriginServiceFee)
	suite.Equal(unit.Cents(11187), cost.DestinationServiceFee)
	suite.Equal(unit.Cents(637056), cost.GCC)
}
