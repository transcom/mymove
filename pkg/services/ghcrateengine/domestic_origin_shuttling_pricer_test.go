package ghcrateengine

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginShuttlingPricer() {
	suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, testServiceSchedule, testBasePriceCents, testdatagen.DefaultContractCode, testEscalationCompounded)

	paymentServiceItem := suite.setupDomesticOriginShuttlingServiceItem()
	pricer := NewDomesticOriginShuttlingPricer(suite.DB())

	suite.Run("success using PaymentServiceItemParams", func() {
		priceCents, displayParams, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(testPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(testEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(testBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		priceCents, _, err := pricer.Price(testdatagen.DefaultContractCode, testRequestedPickupDate, testWeight, testServiceSchedule)
		suite.NoError(err)
		suite.Equal(testPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		_, _, err := pricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid weight", func() {
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, testRequestedPickupDate, badWeight, testServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of 500")
	})

	suite.Run("not finding a rate record", func() {
		_, _, err := pricer.Price("BOGUS", testRequestedPickupDate, testWeight, testServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "Could not lookup Domestic Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		twoYearsLaterPickupDate := testRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, twoYearsLaterPickupDate, testWeight, testServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "Could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticOriginShuttlingServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDOSHUT,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   testRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(testServiceSchedule),
			},
			{
				Key:     models.ServiceItemParamNameWeightActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "1400",
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(testWeight)),
			},
			{
				Key:     models.ServiceItemParamNameWeightEstimated,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "1400",
			},
		},
	)
}
