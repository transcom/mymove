package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	csPriceCents = unit.Cents(8327)
)

var csAvailableToPrimeAt = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceCounselingServices() {
	suite.setupTaskOrderFeeData(models.ReServiceCodeCS, csPriceCents)
	paymentServiceItem := suite.setupCounselingServicesItem()
	counselingServicesPricer := NewCounselingServicesPricer(suite.DB())

	suite.Run("success using PaymentServiceItemParams", func() {
		priceCents, displayParams, err := counselingServicesPricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(csPriceCents, priceCents)

		// Check that PricingDisplayParams have been set and are returned
		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(csPriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		priceCents, _, err := counselingServicesPricer.Price(testdatagen.DefaultContractCode, csAvailableToPrimeAt)
		suite.NoError(err)
		suite.Equal(csPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		_, _, err := counselingServicesPricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("not finding a rate record", func() {
		_, _, err := counselingServicesPricer.Price("BOGUS", csAvailableToPrimeAt)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupCounselingServicesItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeCS,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameMTOAvailableToPrimeAt,
				KeyType: models.ServiceItemParamTypeTimestamp,
				Value:   csAvailableToPrimeAt.Format(TimestampParamFormat),
			},
		},
	)
}
