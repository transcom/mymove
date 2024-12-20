package ghcrateengine

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	msPriceCents = unit.Cents(12303)
)

func (suite *GHCRateEngineServiceSuite) TestPriceManagementServices() {
	lockedPrice := csPriceCents
	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, msPriceCents)
		paymentServiceItem := suite.setupManagementServicesItem()
		managementServicesPricer := NewManagementServicesPricer()

		priceCents, displayParams, err := managementServicesPricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams, nil)
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

		priceCents, _, err := managementServicesPricer.Price(suite.AppContextForTest(), &lockedPrice)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, msPriceCents)
		managementServicesPricer := NewManagementServicesPricer()

		_, _, err := managementServicesPricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{}, nil)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupManagementServicesItem() models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeMS,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameLockedPriceCents,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   msPriceCents.ToMillicents().ToCents().String(),
			},
		}, nil, nil,
	)
}
