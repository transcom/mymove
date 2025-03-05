package ghcrateengine

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *GHCRateEngineServiceSuite) TestIntlNTSHHGPackPricer() {
	ihpkPricer := NewIntlHHGPackPricer()
	pricer := NewIntlNTSHHGPackPricer(ihpkPricer)

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem, contract := suite.setupIntlPackServiceItem(models.ReServiceCodeIHPK)

		totalCost, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)

		// Fetch the INPK market factor from the DB
		inpkReService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeINPK)
		ntsMarketFactor, err := models.FetchMarketFactor(suite.AppContextForTest(), contract.ID, inpkReService.ID, "O")
		suite.FatalNoError(err)

		// Multiply the IHPK price by the NTS markert factor to ensure it math'd properly
		suite.Equal((float64(ihpkTestTotalCost) * ntsMarketFactor), float64(totalCost))

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ihpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ihpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ihpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ihpkTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})
}
