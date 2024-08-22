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
	csPriceCents = unit.Cents(12303)
)

var csAvailableToPrimeAt = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

var lockedPriceCents = unit.Cents(12303)
var mtoServiceItem = models.MTOServiceItem{
	LockedPriceCents: &lockedPriceCents,
}

var failedMtoServiceItem = models.MTOServiceItem{
	LockedPriceCents: nil,
}

func (suite *GHCRateEngineServiceSuite) TestPriceCounselingServices() {
	counselingServicesPricer := NewCounselingServicesPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem := suite.setupCounselingServicesItem()

		priceCents, displayParams, err := counselingServicesPricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(lockedPriceCents, priceCents)

		// Check that PricingDisplayParams have been set and are returned
		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(csPriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeCS, csPriceCents)

		priceCents, _, err := counselingServicesPricer.Price(suite.AppContextForTest(), mtoServiceItem)
		suite.NoError(err)
		suite.Equal(csPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeCS, csPriceCents)

		_, _, err := counselingServicesPricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("not finding a rate record", func() {
		_, _, err := counselingServicesPricer.Price(suite.AppContextForTest(), failedMtoServiceItem)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupCounselingServicesItem() models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeCS,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameMTOAvailableToPrimeAt,
				KeyType: models.ServiceItemParamTypeTimestamp,
				Value:   csAvailableToPrimeAt.Format(TimestampParamFormat),
			},
		}, nil, nil,
	)
}
