package ghcrateengine

import (
	"fmt"
	"github.com/transcom/mymove/pkg/unit"
	"testing"
)

const (
	dlhTestServiceArea = "004"
	dlhTestDistance    = unit.Miles(1200)
	dlhTestWeight      = unit.Pound(4000)
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticServiceArea() {
	suite.setUpDomesticServiceAreaData()
	var services []string{
		"Origin/Destination",
		"SIT 1st Day",
		"SIT Add'l Days",
	}
	for serviceName := range services {
		suite.T().Run(fmt.Sprintf("success %s cost within peak period", serviceName, func(t *testing.T) {

		}

		suite.T().Run(fmt.Sprintf("success %s cost within non-peak period", serviceName,, func(t *testing.T) {

		}

		suite.T().Run(fmt.Sprintf("%s cost weight below minimum", serviceName,, func(t *testing.T) {

		}

		suite.T().Run(fmt.Sprintf("%s date outside of valid contract year", func(t *testing.T) {

		}
	}

	suite.T().Run("validation errors", func(t *testing.T) {

	}
}




func (suite *GHCRateEngineServiceSuite) setUpDomesticServiceAreaData() {
	// create contracts, domestic
}
