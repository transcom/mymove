package ghcrateengine

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	csPriceCents = unit.Cents(12303)
)

func (suite *GHCRateEngineServiceSuite) TestPriceCounselingServices() {
	lockedPrice := csPriceCents
	counselingServicesPricer := NewCounselingServicesPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem := suite.setupCounselingServicesItem()

		priceCents, displayParams, err := counselingServicesPricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(csPriceCents, priceCents)

		// Check that PricingDisplayParams have been set and are returned
		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(csPriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeCS, csPriceCents)

		priceCents, _, err := counselingServicesPricer.Price(suite.AppContextForTest(), &lockedPrice)
		suite.NoError(err)
		suite.Equal(csPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeCS, csPriceCents)

		_, _, err := counselingServicesPricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupCounselingServicesItem() models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeCS,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameLockedPriceCents,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   csPriceCents.ToMillicents().ToCents().String(),
			},
		}, nil, nil,
	)
}
