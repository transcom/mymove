package rateengine

import (
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) Test_Scenario1() {
	t := suite.T()
	scenario.RunRateEngineScenario1(suite.db)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("could not create a development logger: %+v", err)
	}
	planner := route.NewTestingPlanner(362)
	engine := NewRateEngine(suite.db, logger, planner)

	weight := unit.Pound(4000)
	originZip5 := "32168"
	destinationZip5 := "29429"
	date := time.Date(2018, time.June, 18, 0, 0, 0, 0, time.UTC)
	inverseDiscount := 0.33

	cost, err := engine.computePPM(weight, originZip5, destinationZip5, date, inverseDiscount)
	if err != nil {
		t.Fatalf("could not compute PPM: %+v", err)
	}

	expected := CostComputation{
		LinehaulChargeTotal:   unit.Cents(163434),
		OriginServiceFee:      unit.Cents(4765),
		DestinationServiceFee: unit.Cents(5689),
		FullPackUnpackFee:     unit.Cents(89412),
		GCC:                   unit.Cents(263300),
	}

	if cost.LinehaulChargeTotal != expected.LinehaulChargeTotal {
		t.Errorf("wrong LinehaulChargeTotal: expected %s, got: %s", expected.LinehaulChargeTotal, cost.LinehaulChargeTotal)
	}

	if cost.OriginServiceFee != expected.OriginServiceFee {
		t.Errorf("wrong OriginServiceFee: expected %s, got: %s", expected.OriginServiceFee, cost.OriginServiceFee)
	}

	if cost.DestinationServiceFee != expected.DestinationServiceFee {
		t.Errorf("wrong DestinationServiceFee: expected %s, got: %s", expected.DestinationServiceFee, cost.DestinationServiceFee)
	}

	if cost.FullPackUnpackFee != expected.FullPackUnpackFee {
		t.Errorf("wrong FullPackUnpackFee: expected %d, got: %d", expected.FullPackUnpackFee, cost.FullPackUnpackFee)
	}

	if cost.GCC != expected.GCC {
		t.Errorf("wrong GCC: expected %s, got: %s", expected.GCC, cost.GCC)
	}
}
