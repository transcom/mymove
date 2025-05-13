package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	ioshutTestMarket               = "O"
	ioshutTestBasePriceCents       = unit.Cents(15623)
	ioshutTestEscalationCompounded = 1.11000
	ioshutTestWeight               = unit.Pound(4000)
	ioshutTestPriceCents           = unit.Cents(693680)
)

var ioshutTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestInternationalOriginShuttlingPricer() {
	pricer := NewInternationalOriginShuttlingPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)

		paymentServiceItem := suite.setupInternationalOriginShuttlingServiceItem()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(ioshutTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ioshutTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ioshutTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, ioshutTestRequestedPickupDate, ioshutTestWeight, ioshutTestMarket)
		suite.NoError(err)
		suite.Equal(ioshutTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid weight", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, ioshutTestRequestedPickupDate, badWeight, ioshutTestMarket)
		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of 500")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", ioshutTestRequestedPickupDate, ioshutTestWeight, ioshutTestMarket)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup International Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)
		twoYearsLaterPickupDate := ioshutTestRequestedPickupDate.AddDate(10, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, ioshutTestWeight, ioshutTestMarket)
		suite.Error(err)

		suite.Contains(err.Error(), "could not calculate escalated price")

	})
}

func (suite *GHCRateEngineServiceSuite) setupInternationalOriginShuttlingServiceItem() models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeIOSHUT,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   ioshutTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameMarketOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   ioshutTestMarket,
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(ioshutTestWeight)),
			},
		}, nil, nil,
	)
}
