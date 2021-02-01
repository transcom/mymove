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

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticAdditionalDaysSit() {
	suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestEscalationCompounded)

	suite.T().Run("destination golden path", func(t *testing.T) {
		priceCents, err := priceDomesticAdditionalDaysSit(suite.DB(), models.ReServiceCodeDDASIT, DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestIsPeakPeriod, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.NoError(err)
		suite.Equal(ddasitTestPriceCents, priceCents)
	})

	suite.T().Run("invalid service code", func(t *testing.T) {
		_, err := priceDomesticAdditionalDaysSit(suite.DB(), models.ReServiceCodeDDFSIT, DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestIsPeakPeriod, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported additional day sit code")
	})

	suite.T().Run("invalid weight", func(t *testing.T) {
		badWeight := unit.Pound(499)
		_, err := priceDomesticAdditionalDaysSit(suite.DB(), models.ReServiceCodeDDASIT, DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestIsPeakPeriod, badWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 499 less than the minimum")
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, err := priceDomesticAdditionalDaysSit(suite.DB(), models.ReServiceCodeDDASIT, "BOGUS", ddasitTestRequestedPickupDate, ddasitTestIsPeakPeriod, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination additional days SIT rate")
	})

	suite.T().Run("not finding a contract year record", func(t *testing.T) {
		twoYearsLaterPickupDate := ddasitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, err := priceDomesticAdditionalDaysSit(suite.DB(), models.ReServiceCodeDDASIT, DefaultContractCode, twoYearsLaterPickupDate, ddasitTestIsPeakPeriod, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}
