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
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractCode},
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
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, idshutTestBasePriceCents, testdatagen.DefaultContractCode, idshutTestEscalationCompounded)

		twoYearsLaterPickupDate := ioshutTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeIDSHUT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, idshutTestWeight, ioshutTestMarket)

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
