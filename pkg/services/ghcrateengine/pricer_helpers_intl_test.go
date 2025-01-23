package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

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
