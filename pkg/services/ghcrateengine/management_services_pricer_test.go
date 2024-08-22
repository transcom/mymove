package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	msPriceCents = unit.Cents(12303)
)

var msAvailableToPrimeAt = time.Date(testdatagen.TestYear, time.June, 3, 12, 57, 33, 123, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceManagementServices() {
	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, msPriceCents)
		paymentServiceItem := suite.setupManagementServicesItem()
		managementServicesPricer := NewManagementServicesPricer()

		priceCents, displayParams, err := managementServicesPricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)

		// Check that the PricingDisplayParams were successfully set and returned
		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(msPriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, msPriceCents)
		managementServicesPricer := NewManagementServicesPricer()

		priceCents, _, err := managementServicesPricer.Price(suite.AppContextForTest(), mtoServiceItem)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, msPriceCents)
		managementServicesPricer := NewManagementServicesPricer()

		_, _, err := managementServicesPricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("not finding a rate record", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, msPriceCents)
		managementServicesPricer := NewManagementServicesPricer()

		_, _, err := managementServicesPricer.Price(suite.AppContextForTest(), mtoServiceItem)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupManagementServicesItem() models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeMS,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameMTOAvailableToPrimeAt,
				KeyType: models.ServiceItemParamTypeTimestamp,
				Value:   msAvailableToPrimeAt.Format(TimestampParamFormat),
			},
		}, nil, nil,
	)
}
