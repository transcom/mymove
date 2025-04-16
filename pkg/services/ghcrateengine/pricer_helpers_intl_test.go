package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) Test_priceInternationalShuttling() {
	suite.Run("origin golden path", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)

		priceCents, displayParams, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeIOSHUT, testdatagen.DefaultContractCode, ioshutTestRequestedPickupDate, ioshutTestWeight, ioshutTestMarket)
		suite.NoError(err)
		suite.Equal(ioshutTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ioshutTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ioshutTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)
		_, _, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, ioshutTestRequestedPickupDate, ioshutTestWeight, ioshutTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "unsupported international shuttling code")
	})

	suite.Run("invalid weight", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)

		badWeight := unit.Pound(250)
		_, _, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeIOSHUT, testdatagen.DefaultContractCode, ioshutTestRequestedPickupDate, badWeight, ioshutTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of 500")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)

		_, _, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeIOSHUT, "BOGUS", ioshutTestRequestedPickupDate, ioshutTestWeight, ioshutTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup International Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)

		tenYearsLaterPickupDate := ioshutTestRequestedPickupDate.AddDate(10, 0, 0)
		_, _, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeIOSHUT, testdatagen.DefaultContractCode, tenYearsLaterPickupDate, ioshutTestWeight, ioshutTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "could not calculate escalated price: could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceIntlPackUnpack() {
	suite.Run("success with IHPK", func() {
		suite.setupIntlPackServiceItem()
		totalCost, displayParams, err := priceIntlPackUnpack(suite.AppContextForTest(), models.ReServiceCodeIHPK, testdatagen.DefaultContractCode, ihpkTestRequestedPickupDate, ihpkTestWeight, ihpkTestPerUnitCents.Int())
		suite.NoError(err)
		suite.Equal(ihpkTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ihpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ihpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ihpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ihpkTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {
		suite.setupIntlPackServiceItem()
		_, _, err := priceIntlPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, ihpkTestRequestedPickupDate, ihpkTestWeight, ihpkTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pack/unpack code")

		_, _, err = priceIntlPackUnpack(suite.AppContextForTest(), models.ReServiceCodeIHPK, "", ihpkTestRequestedPickupDate, ihpkTestWeight, ihpkTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceIntlPackUnpack(suite.AppContextForTest(), models.ReServiceCodeIHPK, testdatagen.DefaultContractCode, time.Time{}, ihpkTestWeight, ihpkTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceIntlPackUnpack(suite.AppContextForTest(), models.ReServiceCodeIHPK, testdatagen.DefaultContractCode, ihpkTestRequestedPickupDate, ihpkTestWeight, 0)
		suite.Error(err)
		suite.Contains(err.Error(), "PerUnitCents is required")
	})

}

func (suite *GHCRateEngineServiceSuite) TestPriceIntlFirstDaySIT() {
	suite.Run("success with IDFSIT", func() {
		suite.setupIntlDestinationFirstDayServiceItem()
		totalCost, displayParams, err := priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDFSIT, testdatagen.DefaultContractCode, idfsitTestRequestedPickupDate, idfsitTestWeight, idfsitTestPerUnitCents.Int())
		suite.NoError(err)
		suite.Equal(idfsitTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: idfsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(idfsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(idfsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(idfsitTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success with IOFSIT", func() {
		suite.setupIntlOriginFirstDayServiceItem()
		totalCost, displayParams, err := priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeIOFSIT, testdatagen.DefaultContractCode, iofsitTestRequestedPickupDate, iofsitTestWeight, iofsitTestPerUnitCents.Int())
		suite.NoError(err)
		suite.Equal(iofsitTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: iofsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(iofsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(iofsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(iofsitTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {
		suite.setupIntlDestinationFirstDayServiceItem()
		_, _, err := priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, idfsitTestRequestedPickupDate, idfsitTestWeight, idfsitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported first day SIT code")

		_, _, err = priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDFSIT, "", idfsitTestRequestedPickupDate, idfsitTestWeight, idfsitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDFSIT, testdatagen.DefaultContractCode, time.Time{}, idfsitTestWeight, idfsitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDFSIT, testdatagen.DefaultContractCode, idfsitTestRequestedPickupDate, idfsitTestWeight, 0)
		suite.Error(err)
		suite.Contains(err.Error(), "PerUnitCents is required")
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceIntlAdditionalDaySIT() {
	suite.Run("success with IDASIT", func() {
		suite.setupIntlDestinationAdditionalDayServiceItem()
		totalCost, displayParams, err := priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDASIT, testdatagen.DefaultContractCode, idasitTestRequestedPickupDate, idasitNumerDaysInSIT, idasitTestWeight, idasitTestPerUnitCents.Int())
		suite.NoError(err)
		suite.Equal(idasitTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: idasitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(idasitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(idasitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(idasitTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success with IOASIT", func() {
		suite.setupIntlOriginAdditionalDayServiceItem()
		totalCost, displayParams, err := priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIOASIT, testdatagen.DefaultContractCode, ioasitTestRequestedPickupDate, idasitNumerDaysInSIT, ioasitTestWeight, ioasitTestPerUnitCents.Int())
		suite.NoError(err)
		suite.Equal(ioasitTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ioasitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ioasitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ioasitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ioasitTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {
		suite.setupIntlDestinationAdditionalDayServiceItem()
		_, _, err := priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, idasitTestRequestedPickupDate, idasitNumerDaysInSIT, idasitTestWeight, idasitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported additional day SIT code")

		_, _, err = priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDASIT, "", idasitTestRequestedPickupDate, idasitNumerDaysInSIT, idasitTestWeight, idasitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDASIT, testdatagen.DefaultContractCode, time.Time{}, idasitNumerDaysInSIT, idasitTestWeight, idasitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDASIT, testdatagen.DefaultContractCode, idasitTestRequestedPickupDate, idasitNumerDaysInSIT, idasitTestWeight, 0)
		suite.Error(err)
		suite.Contains(err.Error(), "PerUnitCents is required")

		_, _, err = priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDASIT, testdatagen.DefaultContractCode, idasitTestRequestedPickupDate, 0, idasitTestWeight, 0)
		suite.Error(err)
		suite.Contains(err.Error(), "NumberDaysSIT is required")
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceIntlCratingUncrating() {
	suite.Run("crating golden path", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		priceCents, displayParams, err := priceIntlCratingUncrating(suite.AppContextForTest(), models.ReServiceCodeICRT, testdatagen.DefaultContractCode, icrtTestRequestedPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)
		suite.NoError(err)
		suite.Equal(icrtTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(icrtTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(icrtTestBasePriceCents)},
			{Key: models.ServiceItemParamNameUncappedRequestTotal, Value: FormatCents(dcrtTestUncappedRequestTotal)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)
		_, _, err := priceIntlCratingUncrating(suite.AppContextForTest(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, icrtTestRequestedPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "unsupported international crating/uncrating code")
	})

	suite.Run("invalid crate size - external crate", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		badSize := unit.CubicFeet(1.0)
		_, _, err := priceIntlCratingUncrating(suite.AppContextForTest(), models.ReServiceCodeICRT, testdatagen.DefaultContractCode, icrtTestRequestedPickupDate, badSize, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, true, icrtTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "external crates must be billed for a minimum of 4.00 cubic feet")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		_, _, err := priceIntlCratingUncrating(suite.AppContextForTest(), models.ReServiceCodeICRT, "BOGUS", icrtTestRequestedPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup International Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		twoYearsLaterPickupDate := ioshutTestRequestedPickupDate.AddDate(10, 0, 0)
		_, _, err := priceIntlCratingUncrating(suite.AppContextForTest(), models.ReServiceCodeICRT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "could not calculate escalated price: could not lookup contract year")
	})
}
