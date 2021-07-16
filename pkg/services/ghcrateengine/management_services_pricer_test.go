package ghcrateengine

import (
	"time"

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
		managementServicesPricer := NewManagementServicesPricer(suite.DB())

		priceCents, displayParams, err := managementServicesPricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
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
		managementServicesPricer := NewManagementServicesPricer(suite.DB())

		priceCents, _, err := managementServicesPricer.Price(testdatagen.DefaultContractCode, msAvailableToPrimeAt)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, msPriceCents)
		managementServicesPricer := NewManagementServicesPricer(suite.DB())

		_, _, err := managementServicesPricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("not finding a rate record", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, msPriceCents)
		managementServicesPricer := NewManagementServicesPricer(suite.DB())

		_, _, err := managementServicesPricer.Price("BOGUS", msAvailableToPrimeAt)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupManagementServicesItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeMS,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameMTOAvailableToPrimeAt,
				KeyType: models.ServiceItemParamTypeTimestamp,
				Value:   msAvailableToPrimeAt.Format(TimestampParamFormat),
			},
		},
	)
}
