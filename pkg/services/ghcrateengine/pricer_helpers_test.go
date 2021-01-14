package ghcrateengine

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticFirstDaySit() {
	suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestEscalationCompounded)

	suite.T().Run("destination golden path", func(t *testing.T) {
		priceCents, err := priceDomesticFirstDaySit(suite.DB(), models.ReServiceCodeDDFSIT, DefaultContractCode, ddfsitTestRequestedPickupDate, ddfsitTestIsPeakPeriod, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.NoError(err)
		suite.Equal(ddfsitTestPriceCents, priceCents)
	})

	suite.T().Run("invalid service code", func(t *testing.T) {
		_, err := priceDomesticFirstDaySit(suite.DB(), models.ReServiceCodeCS, DefaultContractCode, ddfsitTestRequestedPickupDate, ddfsitTestIsPeakPeriod, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported first day sit code")
	})

	suite.T().Run("invalid weight", func(t *testing.T) {
		badWeight := unit.Pound(250)
		_, err := priceDomesticFirstDaySit(suite.DB(), models.ReServiceCodeDDFSIT, DefaultContractCode, ddfsitTestRequestedPickupDate, ddfsitTestIsPeakPeriod, badWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, err := priceDomesticFirstDaySit(suite.DB(), models.ReServiceCodeDDFSIT, "BOGUS", ddfsitTestRequestedPickupDate, ddfsitTestIsPeakPeriod, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination first day SIT rate")
	})

	suite.T().Run("not finding a contract year record", func(t *testing.T) {
		twoYearsLaterPickupDate := ddfsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, err := priceDomesticFirstDaySit(suite.DB(), models.ReServiceCodeDDFSIT, DefaultContractCode, twoYearsLaterPickupDate, ddfsitTestIsPeakPeriod, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}
