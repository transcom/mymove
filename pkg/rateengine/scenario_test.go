package rateengine

import (
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *RateEngineSuite) Test_Scenario1() {
	scenario.RunRateEngineScenario1(suite.db)

	// engine := NewRateEngine(suite.db, logger, etc...)
}
