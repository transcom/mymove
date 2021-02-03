package ghcrateengine

import (
	"fmt"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySITSameZip3s() {
	dshZipDest := "30907"
	dshZipSITDest := "30901" // same zip3
	dshDistance := unit.Miles(15)

	suite.T().Run("destination golden path for same zip3s", func(t *testing.T) {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDSH, dddsitTestServiceArea, dddsitTestIsPeakPeriod, dddsitTestDomesticServiceAreaBasePriceCents, dddsitTestEscalationCompounded)
		priceCents, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.NoError(err)
		expectedPrice := unit.Cents(53187) // dddsitTestDomesticServiceAreaBasePriceCents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("invalid service code", func(t *testing.T) {
		_, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pickup/delivery SIT code")
	})

	suite.T().Run("invalid weight", func(t *testing.T) {
		badWeight := unit.Pound(250)
		_, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, badWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		expectedError := fmt.Sprintf("weight of %d less than the minimum", badWeight)
		suite.Contains(err.Error(), expectedError)
	})

	suite.T().Run("bad destination zip", func(t *testing.T) {
		_, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, "309", dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid destination postal code")
	})

	suite.T().Run("bad SIT destination zip", func(t *testing.T) {
		_, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, "1234", dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT destination postal code")
	})

	suite.T().Run("error from shorthaul pricer", func(t *testing.T) {
		_, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price shorthaul")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySIT50PlusMilesDiffZip3s() {
	dlhZipDest := "30907"
	dlhZipSITDest := "36106"       // different zip3
	dlhDistance := unit.Miles(305) // > 50 miles

	suite.T().Run("destination golden path for > 50 miles with different zip3s", func(t *testing.T) {
		suite.setupDomesticLinehaulPrice(dddsitTestServiceArea, dddsitTestIsPeakPeriod, dddsitTestWeightLower, dddsitTestWeightUpper, dddsitTestMilesLower, dddsitTestMilesUpper, dddsitTestDomesticLinehaulBasePriceMillicents, dddsitTestEscalationCompounded)
		priceCents, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dlhZipDest, dlhZipSITDest, dlhDistance)
		suite.NoError(err)
		expectedPriceMillicents := unit.Millicents(45944438) // dddsitTestDomesticLinehaulBasePriceMillicents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded
		expectedPrice := expectedPriceMillicents.ToCents()
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("error from linehaul pricer", func(t *testing.T) {
		_, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dlhZipDest, dlhZipSITDest, dlhDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price linehaul")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySIT50MilesOrLessDiffZip3s() {
	domOtherZipDest := "30907"
	domOtherZipSITDest := "29801"      // different zip3
	domOtherDistance := unit.Miles(37) // <= 50 miles

	suite.T().Run("destination golden path for <= 50 miles with different zip3s", func(t *testing.T) {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestEscalationCompounded)
		priceCents, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.NoError(err)
		expectedPrice := unit.Cents(58355) // dddsitTestDomesticOtherBasePriceCents * (dddsitTestWeight / 100) * dddsitTestEscalationCompounded
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination SIT delivery rate")
	})

	suite.T().Run("not finding a contract year record", func(t *testing.T) {
		twoYearsLaterPickupDate := dddsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}
