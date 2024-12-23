package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testhelpers"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPackUnpack() {
	featureFlagValues := testhelpers.MakeMobileHomeFFMap()
	suite.Run("golden path with DNPK", func() {
		suite.setupDomesticNTSPackPrices(dnpkTestServicesScheduleOrigin, dnpkTestIsPeakPeriod, dnpkTestBasePriceCents, models.MarketConus, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)
		isPPM := false
		isMobileHome := false
		priceCents, displayParams, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, dnpkTestRequestedPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.NoError(err)
		suite.Equal(dnpkTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dnpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dnpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dnpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dnpkTestBasePriceCents)},
			{Key: models.ServiceItemParamNameNTSPackingFactor, Value: FormatFloat(dnpkTestFactor, 2)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {
		isPPM := false
		isMobileHome := false
		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, dnpkTestRequestedPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pack/unpack code")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, "", dnpkTestRequestedPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, time.Time{}, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, dnpkTestRequestedPickupDate, unit.Pound(250), dnpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, dnpkTestRequestedPickupDate, dnpkTestWeight, 0, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "Services schedule is required")
	})

	suite.Run("not finding domestic other price", func() {
		suite.setupDomesticNTSPackPrices(dnpkTestServicesScheduleOrigin, dnpkTestIsPeakPeriod, dnpkTestBasePriceCents, models.MarketConus, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)
		badContractCode := "BOGUS"
		isPPM := false
		isMobileHome := false
		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, badContractCode, dnpkTestRequestedPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup domestic other price")
	})

	suite.Run("not finding contract year", func() {
		suite.setupDomesticNTSPackPrices(dnpkTestServicesScheduleOrigin, dnpkTestIsPeakPeriod, dnpkTestBasePriceCents, models.MarketConus, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)
		twoYearsLaterPickupDate := dnpkTestRequestedPickupDate.AddDate(2, 0, 0)
		isPPM := false
		isMobileHome := false
		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})

	suite.Run("not finding shipment type price", func() {
		badMarket := models.MarketOconus
		isPPM := false
		isMobileHome := false
		suite.setupDomesticNTSPackPrices(dnpkTestServicesScheduleOrigin, dnpkTestIsPeakPeriod, dnpkTestBasePriceCents, badMarket, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)
		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, dnpkTestRequestedPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup shipment type price")
	})

}
func (suite *GHCRateEngineServiceSuite) Test_domesticPackAndUnpackWithMobileHome() {
	featureFlagValues := testhelpers.MakeMobileHomeFFMap()
	suite.Run("golden path with DPK", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod, dpkTestBasePriceCents, dpkTestContractYearName, dpkTestEscalationCompounded)

		isPPM := false
		isMobileHome := true
		priceCents, displayParams, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, testdatagen.DefaultContractCode, dpkTestRequestedPickupDate, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.NoError(err)
		suite.Equal(unit.Cents(roundToPrecision(float64(dpkTestPriceCents)*mobileHomeFactor, 2)), priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dpkTestBasePriceCents)},
			{Key: models.ServiceItemParamNameMobileHomeFactor, Value: FormatFloat(mobileHomeFactor, 2)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {
		isPPM := false
		isMobileHome := true
		suite.setupDomesticOtherPrice(models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod, dpkTestBasePriceCents, dpkTestContractYearName, dpkTestEscalationCompounded)

		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, dpkTestRequestedPickupDate, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pack/unpack code")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, "", dnpkTestRequestedPickupDate, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, testdatagen.DefaultContractCode, time.Time{}, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, testdatagen.DefaultContractCode, dpkTestRequestedPickupDate, dpkTestWeight, 0, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "Services schedule is required")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_domesticPackAndUnpackWithPPM() {
	featureFlagValues := testhelpers.MakeMobileHomeFFMap()
	suite.Run("golden path with DPK", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod, dpkTestBasePriceCents, dpkTestContractYearName, dpkTestEscalationCompounded)

		isPPM := true
		isMobileHome := false
		priceCents, displayParams, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, testdatagen.DefaultContractCode, dpkTestRequestedPickupDate, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.NoError(err)
		suite.Equal(dpkTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dpkTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {
		isPPM := true
		isMobileHome := false
		suite.setupDomesticOtherPrice(models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod, dpkTestBasePriceCents, dpkTestContractYearName, dpkTestEscalationCompounded)

		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, dpkTestRequestedPickupDate, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pack/unpack code")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, "", dnpkTestRequestedPickupDate, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, testdatagen.DefaultContractCode, time.Time{}, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, testdatagen.DefaultContractCode, dpkTestRequestedPickupDate, dpkTestWeight, 0, isPPM, isMobileHome, featureFlagValues)
		suite.Error(err)
		suite.Contains(err.Error(), "Services schedule is required")
	})
}
func (suite *GHCRateEngineServiceSuite) Test_priceDomesticFirstDaySIT() {
	suite.Run("destination golden path", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		priceCents, displayParams, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, testdatagen.DefaultContractCode, ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea, false)
		suite.NoError(err)
		suite.Equal(ddfsitTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ddfsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ddfsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ddfsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ddfsitTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		_, _, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea, false)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported first day sit code")
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		badWeight := unit.Pound(250)
		_, _, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, testdatagen.DefaultContractCode, ddfsitTestRequestedPickupDate, badWeight, ddfsitTestServiceArea, false)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.Run("no error when weight minimum is overridden", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		weight := unit.Pound(250)
		_, _, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, testdatagen.DefaultContractCode, ddfsitTestRequestedPickupDate, weight, ddfsitTestServiceArea, true)
		suite.NoError(err)
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		_, _, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, "BOGUS", ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea, false)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination first day SIT rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		twoYearsLaterPickupDate := ddfsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, ddfsitTestWeight, ddfsitTestServiceArea, false)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticAdditionalDaysSIT() {

	suite.Run("destination golden path", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		priceCents, displayParams, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDASIT, testdatagen.DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT, false)
		suite.NoError(err)
		suite.Equal(ddasitTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ddasitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ddasitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ddasitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ddasitTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		_, _, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, testdatagen.DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT, false)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported additional day sit code")
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		badWeight := unit.Pound(499)
		_, _, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDASIT, testdatagen.DefaultContractCode, ddasitTestRequestedPickupDate, badWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT, false)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 499 less than the minimum")
	})

	suite.Run("no error when weight minimum is overridden", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		weight := unit.Pound(499)
		_, _, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDASIT, testdatagen.DefaultContractCode, ddasitTestRequestedPickupDate, weight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT, true)
		suite.NoError(err)
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		_, _, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDASIT, "BOGUS", ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT, false)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination additional days SIT rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		twoYearsLaterPickupDate := ddasitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDASIT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT, false)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySITSameZip3s() {
	dshZipDest := "30907"
	dshZipSITDest := "30901" // same zip3
	dshDistance := unit.Miles(15)
	dshContractName := "dshTestYear"

	suite.Run("destination golden path for same zip3s", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dshContractName, dddsitTestEscalationCompounded)
		priceCents, displayParams, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.NoError(err)
		expectedPrice := unit.Cents(58365) // dddsitTestDomesticServiceAreaBasePriceCents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dshContractName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dddsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dddsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dddsitTestDomesticOtherBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pickup/delivery SIT code")
	})

	suite.Run("invalid weight", func() {
		badWeight := unit.Pound(250)
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, badWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		expectedError := fmt.Sprintf("weight of %d less than the minimum", badWeight)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("bad destination zip", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, "309", dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT original destination postal code")
	})

	suite.Run("bad SIT final destination zip", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, "1234", dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT final destination postal code")
	})

	suite.Run("error fetching domestic destination SIT delivery rate", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination SIT delivery rate")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySIT50PlusMilesDiffZip3s() {
	dlhZipDest := "30907"
	dlhZipSITDest := "36106"       // different zip3
	dlhDistance := unit.Miles(305) // > 50 miles
	dlhContractName := "dhlTestYear"

	suite.Run("destination golden path for > 50 miles with different zip3s", func() {
		suite.setupDomesticLinehaulPrice(dddsitTestServiceArea, dddsitTestIsPeakPeriod, dddsitTestWeightLower, dddsitTestWeightUpper, dddsitTestMilesLower, dddsitTestMilesUpper, dddsitTestDomesticLinehaulBasePriceMillicents, dlhContractName, dddsitTestEscalationCompounded)
		priceCents, displayParams, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dlhZipDest, dlhZipSITDest, dlhDistance)
		suite.NoError(err)
		expectedPrice := unit.Cents(45979)

		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dlhContractName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dddsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dddsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatFloat(dddsitTestDomesticLinehaulBasePriceMillicents.ToDollarFloatNoRound(), 3)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("error from linehaul pricer", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dlhZipDest, dlhZipSITDest, dlhDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price linehaul")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySIT50MilesOrLessDiffZip3s() {
	domOtherZipDest := "30907"
	domOtherZipSITDest := "29801"      // different zip3
	domOtherDistance := unit.Miles(37) // <= 50 miles
	domContractName := "domTestYear"

	suite.Run("destination golden path for <= 50 miles with different zip3s", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, domContractName, dddsitTestEscalationCompounded)
		priceCents, displayParams, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.NoError(err)
		expectedPrice := unit.Cents(58365)
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: domContractName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dddsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dddsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dddsitTestDomesticOtherBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("not finding a rate record", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination SIT delivery rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, domContractName, dddsitTestEscalationCompounded)
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.NoError(err)

		twoYearsLaterPickupDate := dddsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err = priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySIT50MilesOrLessSameZip3s() {
	domOtherZipDest := "30907"
	domOtherZipSITDest := "30910"      // same zip3
	domOtherDistance := unit.Miles(37) // <= 50 miles
	domContractName := "domTestYear"

	suite.Run("destination golden path for <= 50 miles with same zip3s", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, domContractName, dddsitTestEscalationCompounded)
		priceCents, displayParams, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.NoError(err)
		expectedPrice := unit.Cents(58365)
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: domContractName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dddsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dddsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dddsitTestDomesticOtherBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("not finding a rate record", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination SIT delivery rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, domContractName, dddsitTestEscalationCompounded)
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.NoError(err)

		twoYearsLaterPickupDate := dddsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err = priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

type pricerParamsSubtestData struct {
	params             services.PricingDisplayParams
	paymentServiceItem models.PaymentServiceItem
}

func (suite *GHCRateEngineServiceSuite) makePricerParamsSubtestData() (subtestData *pricerParamsSubtestData) {
	subtestData = &pricerParamsSubtestData{}
	subtestData.params = services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: "4000.90",
		}, {
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: "1.06",
		}, {
			Key:   models.ServiceItemParamNameIsPeak,
			Value: "True",
		}, {
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "TRUSS_TEST",
		},
	}

	factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNamePriceRateOrFactor,
				Description: "Price Rate Or Factor",
				Type:        models.ServiceItemParamTypeDecimal,
				Origin:      models.ServiceItemParamOriginPricer,
			},
		},
	}, nil)
	factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameEscalationCompounded,
				Description: "Escalation compounded",
				Type:        models.ServiceItemParamTypeDecimal,
				Origin:      models.ServiceItemParamOriginPricer,
			},
		},
	}, nil)
	factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameIsPeak,
				Description: "Is peak",
				Type:        models.ServiceItemParamTypeBoolean,
				Origin:      models.ServiceItemParamOriginPricer,
			},
		},
	}, nil)
	factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameContractYearName,
				Description: "Contract year name",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginPricer,
			},
		},
	}, nil)

	subtestData.paymentServiceItem = factory.BuildPaymentServiceItem(suite.DB(),
		nil, nil)

	return subtestData
}

func (suite *GHCRateEngineServiceSuite) Test_createPricerGeneratedParams() {
	suite.Run("payment service item params created for the pricer", func() {
		subtestData := suite.makePricerParamsSubtestData()
		paymentServiceItemParams, err := createPricerGeneratedParams(suite.AppContextForTest(), subtestData.paymentServiceItem.ID, subtestData.params)
		suite.NoError(err)
		expectedValues := [4]string{"4000.90", "1.06", "True", "TRUSS_TEST"}
		for _, paymentServiceItemParam := range paymentServiceItemParams {
			switch paymentServiceItemParam.ServiceItemParamKey.Key {
			case models.ServiceItemParamNamePriceRateOrFactor:
				suite.Equal(expectedValues[0], paymentServiceItemParam.Value)
			case models.ServiceItemParamNameEscalationCompounded:
				suite.Equal(expectedValues[1], paymentServiceItemParam.Value)
			case models.ServiceItemParamNameIsPeak:
				suite.Equal(expectedValues[2], paymentServiceItemParam.Value)
			case models.ServiceItemParamNameContractYearName:
				suite.Equal(expectedValues[3], paymentServiceItemParam.Value)
			}
		}
	})

	suite.Run("errors if PaymentServiceItemID is invalid", func() {
		subtestData := suite.makePricerParamsSubtestData()
		invalidID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")

		_, err := createPricerGeneratedParams(suite.AppContextForTest(), invalidID, subtestData.params)
		suite.Error(err)
		suite.Contains(err.Error(), "validation error with creating payment service item param")
	})

	suite.Run("errors if PricingParm points to a serviceItem that doesnt originate from the Pricer", func() {
		subtestData := suite.makePricerParamsSubtestData()
		invalidParam := services.PricingDisplayParams{
			{
				Key:   models.ServiceItemParamNameServiceAreaOrigin,
				Value: "40000.9",
			},
		}

		factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameServiceAreaOrigin,
					Description: "service area actual",
					Type:        models.ServiceItemParamTypeString,
					Origin:      models.ServiceItemParamOriginPrime,
				},
			},
		}, nil)

		_, err := createPricerGeneratedParams(suite.AppContextForTest(), subtestData.paymentServiceItem.ID, invalidParam)
		suite.Error(err)
		suite.Contains(err.Error(), "Service item param key is not a pricer param")
	})

	suite.Run("errors if no PricingParms passed from the Pricer", func() {
		subtestData := suite.makePricerParamsSubtestData()
		emptyParams := services.PricingDisplayParams{}

		_, err := createPricerGeneratedParams(suite.AppContextForTest(), subtestData.paymentServiceItem.ID, emptyParams)
		suite.Error(err)
		suite.Contains(err.Error(), "PricingDisplayParams must not be empty")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticShuttling() {
	suite.Run("destination golden path", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, testdatagen.DefaultContractCode, doshutTestEscalationCompounded)

		priceCents, displayParams, err := priceDomesticShuttling(suite.AppContextForTest(), models.ReServiceCodeDOSHUT, testdatagen.DefaultContractCode, doshutTestRequestedPickupDate, doshutTestWeight, doshutTestServiceSchedule)
		suite.NoError(err)
		suite.Equal(doshutTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(doshutTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(doshutTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, testdatagen.DefaultContractCode, doshutTestEscalationCompounded)
		_, _, err := priceDomesticShuttling(suite.AppContextForTest(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, doshutTestRequestedPickupDate, doshutTestWeight, doshutTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "unsupported domestic shuttling code")
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, testdatagen.DefaultContractCode, doshutTestEscalationCompounded)

		badWeight := unit.Pound(250)
		_, _, err := priceDomesticShuttling(suite.AppContextForTest(), models.ReServiceCodeDOSHUT, testdatagen.DefaultContractCode, doshutTestRequestedPickupDate, badWeight, doshutTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of 500")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, testdatagen.DefaultContractCode, doshutTestEscalationCompounded)

		_, _, err := priceDomesticShuttling(suite.AppContextForTest(), models.ReServiceCodeDOSHUT, "BOGUS", doshutTestRequestedPickupDate, doshutTestWeight, doshutTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup Domestic Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule, ddshutTestBasePriceCents, testdatagen.DefaultContractCode, ddshutTestEscalationCompounded)

		twoYearsLaterPickupDate := doshutTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticShuttling(suite.AppContextForTest(), models.ReServiceCodeDDSHUT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, ddshutTestWeight, ddshutTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "could not calculate escalated price: could not lookup contract year")
	})
}
func (suite *GHCRateEngineServiceSuite) Test_priceDomesticCrating() {
	suite.Run("crating golden path", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)

		priceCents, displayParams, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeDCRT, testdatagen.DefaultContractCode, dcrtTestRequestedPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule, dcrtTestStandaloneCrate, dcrtTestStandaloneCrateCap)
		suite.NoError(err)
		suite.Equal(dcrtTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dcrtTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dcrtTestBasePriceCents)},
			{Key: models.ServiceItemParamNameUncappedRequestTotal, Value: FormatCents(dcrtTestUncappedRequestTotal)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("crating golden path with truncation", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)
		priceCents, displayParams, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeDCRT, testdatagen.DefaultContractCode, dcrtTestRequestedPickupDate, unit.CubicFeet(8.90625), dcrtTestServiceSchedule, dcrtTestStandaloneCrate, dcrtTestStandaloneCrateCap)
		suite.NoError(err)
		suite.Equal(unit.Cents(23049), priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dcrtTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dcrtTestBasePriceCents)},
			{Key: models.ServiceItemParamNameUncappedRequestTotal, Value: FormatCents(23049)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)
		_, _, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, dcrtTestRequestedPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule, dcrtTestStandaloneCrate, dcrtTestStandaloneCrateCap)

		suite.Error(err)
		suite.Contains(err.Error(), "unsupported domestic crating code")
	})

	suite.Run("invalid crate size", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)

		badSize := unit.CubicFeet(1.0)
		_, _, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeDCRT, testdatagen.DefaultContractCode, dcrtTestRequestedPickupDate, badSize, dcrtTestServiceSchedule, dcrtTestStandaloneCrate, dcrtTestStandaloneCrateCap)

		suite.Error(err)
		suite.Contains(err.Error(), "crate must be billed for a minimum of 4 cubic feet")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)

		_, _, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeDCRT, "BOGUS", dcrtTestRequestedPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule, dcrtTestStandaloneCrate, dcrtTestStandaloneCrateCap)

		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup Domestic Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)

		twoYearsLaterPickupDate := dcrtTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeDCRT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule, dcrtTestStandaloneCrate, dcrtTestStandaloneCrateCap)

		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_escalatePriceForContractYear() {
	suite.Run("escalated price is rounded to the nearest cent for non-linehaul pricing", func() {
		escalationCompounded := 1.04071
		cy := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
				},
			})
		isLinehaul := false
		basePrice := 5117.0

		escalatedPrice, contractYear, err := escalatePriceForContractYear(suite.AppContextForTest(), cy.ContractID, cy.StartDate.AddDate(0, 0, 1), isLinehaul, basePrice)

		suite.Nil(err)
		suite.Equal(cy.ID, contractYear.ID)
		suite.Equal(5325.0, escalatedPrice)
	})

	suite.Run("escalated price is rounded to the nearest tenth-cent for linehaul pricing", func() {
		escalationCompounded := 1.04071
		cy := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
				},
			})
		isLinehaul := true
		basePrice := 5117.0

		escalatedPrice, contractYear, err := escalatePriceForContractYear(suite.AppContextForTest(), cy.ContractID, cy.StartDate.AddDate(0, 0, 1), isLinehaul, basePrice)

		suite.Nil(err)
		suite.Equal(cy.ID, contractYear.ID)
		suite.Equal(5325.3, escalatedPrice)
	})

	suite.Run("not finding contract year", func() {
		isLinehaul := true
		basePrice := 5117.0

		_, _, err := escalatePriceForContractYear(suite.AppContextForTest(), uuid.Nil, time.Time{}, isLinehaul, basePrice)

		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_escalatedPriceForContractYearCompounded() {

	setUpContracts := func() map[string]models.ReContractYear {
		escalationCompounded := 1.04071
		basePeriodYear1 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.BasePeriodYear1,
				},
			})
		basePeriodYear2 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.BasePeriodYear2,
					StartDate:            basePeriodYear1.StartDate.AddDate(1, 0, 0),
					EndDate:              basePeriodYear1.EndDate.AddDate(1, 0, 0),
				},
			})
		basePeriodYear3 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.BasePeriodYear3,
					StartDate:            basePeriodYear2.StartDate.AddDate(1, 0, 0),
					EndDate:              basePeriodYear2.EndDate.AddDate(1, 0, 0),
				},
			})
		optionPeriod1 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.OptionPeriod1,
					StartDate:            basePeriodYear3.StartDate.AddDate(1, 0, 0),
					EndDate:              basePeriodYear3.EndDate.AddDate(1, 0, 0),
				},
			})
		optionPeriod2 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.OptionPeriod2,
					StartDate:            optionPeriod1.StartDate.AddDate(1, 0, 0),
					EndDate:              optionPeriod1.EndDate.AddDate(1, 0, 0),
				},
			})
		awardTerm1 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.AwardTerm1,
					StartDate:            optionPeriod2.StartDate.AddDate(1, 0, 0),
					EndDate:              optionPeriod2.EndDate.AddDate(1, 0, 0),
				},
			})
		awardTerm2 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.AwardTerm2,
					StartDate:            awardTerm1.StartDate.AddDate(1, 0, 0),
					EndDate:              awardTerm1.EndDate.AddDate(1, 0, 0),
				},
			})

		optionPeriod3 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.OptionPeriod3,
					StartDate:            awardTerm2.StartDate.AddDate(1, 0, 0),
					EndDate:              awardTerm2.EndDate.AddDate(1, 0, 0),
				},
			})

		contractsYearsMap := make(map[string]models.ReContractYear)
		contractsYearsMap[optionPeriod3.Name] = optionPeriod3
		contractsYearsMap[awardTerm2.Name] = awardTerm2
		contractsYearsMap[awardTerm1.Name] = awardTerm1
		contractsYearsMap[optionPeriod2.Name] = optionPeriod2
		contractsYearsMap[optionPeriod1.Name] = optionPeriod1
		contractsYearsMap[basePeriodYear3.Name] = basePeriodYear3
		contractsYearsMap[basePeriodYear2.Name] = basePeriodYear2
		contractsYearsMap[basePeriodYear1.Name] = basePeriodYear1
		return contractsYearsMap
	}

	setUpContractsWithMissingContracts := func() map[string]models.ReContractYear {
		escalationCompounded := 1.04071
		basePeriodYear1 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.BasePeriodYear1,
				},
			})
		basePeriodYear2 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.BasePeriodYear2,
					StartDate:            basePeriodYear1.StartDate.AddDate(1, 0, 0),
					EndDate:              basePeriodYear1.EndDate.AddDate(1, 0, 0),
				},
			})
		basePeriodYear3 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.BasePeriodYear3,
					StartDate:            basePeriodYear2.StartDate.AddDate(1, 0, 0),
					EndDate:              basePeriodYear2.EndDate.AddDate(1, 0, 0),
				},
			})
		optionPeriod2 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.OptionPeriod2,
					StartDate:            basePeriodYear2.StartDate.AddDate(2, 0, 0),
					EndDate:              basePeriodYear2.EndDate.AddDate(2, 0, 0),
				},
			})
		awardTerm1 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.AwardTerm1,
					StartDate:            optionPeriod2.StartDate.AddDate(1, 0, 0),
					EndDate:              optionPeriod2.EndDate.AddDate(1, 0, 0),
				},
			})
		awardTerm2 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.AwardTerm2,
					StartDate:            awardTerm1.StartDate.AddDate(1, 0, 0),
					EndDate:              awardTerm1.EndDate.AddDate(1, 0, 0),
				},
			})

		optionPeriod3 := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: escalationCompounded,
					Name:                 models.OptionPeriod3,
					StartDate:            awardTerm2.StartDate.AddDate(1, 0, 0),
					EndDate:              awardTerm2.EndDate.AddDate(1, 0, 0),
				},
			})

		contractsYearsMap := make(map[string]models.ReContractYear)
		contractsYearsMap[optionPeriod3.Name] = optionPeriod3
		contractsYearsMap[awardTerm2.Name] = awardTerm2
		contractsYearsMap[awardTerm1.Name] = awardTerm1
		contractsYearsMap[optionPeriod2.Name] = optionPeriod2
		contractsYearsMap[basePeriodYear3.Name] = basePeriodYear3
		contractsYearsMap[basePeriodYear2.Name] = basePeriodYear2
		contractsYearsMap[basePeriodYear1.Name] = basePeriodYear1
		return contractsYearsMap
	}

	suite.Run("should correctly calculate escalated price based on the escalation factors of each contract year", func() {
		contracts := setUpContracts()

		isLinehaul := false
		basePrice := 5117.0

		contract := contracts[models.OptionPeriod3]

		escalatedPrice, contractYear, err := escalatePriceForContractYear(suite.AppContextForTest(), contract.ContractID, contract.StartDate.AddDate(0, 0, 1), isLinehaul, basePrice)

		suite.Nil(err)
		suite.Equal(contract.ID, contractYear.ID)
		suite.Equal(5981.0, escalatedPrice)
	})
	suite.Run("should error if an expected contract needed for the escalation price calculation is not found", func() {
		contracts := setUpContractsWithMissingContracts()
		isLinehaul := false
		basePrice := 5117.0

		contract := contracts[models.OptionPeriod3]
		escalatedPrice, contractYear, err := escalatePriceForContractYear(suite.AppContextForTest(), contract.ContractID, contract.StartDate.AddDate(0, 0, 1), isLinehaul, basePrice)

		suite.Error(err)
		suite.Equal("expected contract Option Period 1 not found", err.Error())
		suite.Equal(contract.ID, contractYear.ID)
		suite.NotNil(escalatedPrice)
	})
}
