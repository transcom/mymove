package ghcrateengine

import (
	"math"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) TestIntlNTSHHGPackPricer() {
	ihpkPricer := NewIntlHHGPackPricer()
	pricer := NewIntlNTSHHGPackPricer(ihpkPricer)

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem, contract := suite.setupIntlPackServiceItem(models.ReServiceCodeIHPK)

		totalCost, pricerParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.FatalNoError(err)

		// Fetch the INPK market factor from the DB
		inpkReService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeINPK)
		ntsMarketFactor, err := models.FetchMarketFactor(suite.AppContextForTest(), contract.ID, inpkReService.ID, "O")
		suite.FatalNoError(err)
		suite.FatalTrue(suite.NotEmpty(ntsMarketFactor))

		// Multiply the IHPK price by the NTS market factor to ensure it math'd properly
		suite.FatalTrue(suite.Equal(math.Round((float64(ihpkTestTotalCost) * ntsMarketFactor)), float64(totalCost)))

		expectedPricerParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ihpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ihpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ihpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ihpkTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedPricerParams, pricerParams)
	})
	suite.Run("success using PaymentServiceItemParams", func() {
		pickupDate := time.Date(2018, time.September, 14, 12, 0, 0, 0, time.UTC)
		cents := unit.Cents(6000)
		_, contract := suite.setupIntlPackServiceItem(models.ReServiceCodeIHPK)

		// Get the IHPK pricer result
		ihpkCost, _, err := ihpkPricer.Price(suite.AppContextForTest(), contract.Code, pickupDate, 500, int(cents))
		suite.NoError(err)

		// Get the INPK pricer result
		totalCost, displayParams, err := pricer.Price(suite.AppContextForTest(), contract.Code, pickupDate, 500, 6000)
		suite.NoError(err)

		// Fetch the INPK market factor from the DB
		inpkReService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeINPK)
		ntsMarketFactor, err := models.FetchMarketFactor(suite.AppContextForTest(), contract.ID, inpkReService.ID, "O")
		suite.FatalNoError(err)

		// Multiply the IHPK pricer result by the NTS market factor to ensure it math'd properly and matches the INPK pricer
		suite.Equal((float64(ihpkCost) * ntsMarketFactor), float64(totalCost))

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ihpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ihpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ihpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(cents)},
			{Key: models.ServiceItemParamNameNTSPackingFactor, Value: FormatFloat(ntsMarketFactor, -1)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})
}
