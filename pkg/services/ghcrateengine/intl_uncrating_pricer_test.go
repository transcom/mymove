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
	iucrtTestMarket               = models.Market("O")
	iucrtTestBasePriceCents       = unit.Cents(654)
	iucrtTestEscalationCompounded = 1.11000
	iucrtTestBilledCubicFeet      = 10
	iucrtTestPriceCents           = unit.Cents(7260)
	iucrtTestUncappedRequestTotal = unit.Cents(7260)
)

var iucrtTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestIntlUncratingPricer() {
	pricer := NewIntlUncratingPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIUCRT, iucrtTestMarket, iucrtTestBasePriceCents, testdatagen.DefaultContractCode, iucrtTestEscalationCompounded)

		paymentServiceItem := suite.setupIntlUncratingServiceItem()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(iucrtTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(iucrtTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(iucrtTestBasePriceCents)},
			{Key: models.ServiceItemParamNameUncappedRequestTotal, Value: FormatCents(iucrtTestUncappedRequestTotal)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIUCRT, iucrtTestMarket, iucrtTestBasePriceCents, testdatagen.DefaultContractCode, iucrtTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, iucrtTestRequestedPickupDate, iucrtTestBilledCubicFeet, iucrtTestMarket)
		suite.NoError(err)
		suite.Equal(iucrtTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIUCRT, iucrtTestMarket, iucrtTestBasePriceCents, testdatagen.DefaultContractCode, iucrtTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("not finding a rate record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIUCRT, iucrtTestMarket, iucrtTestBasePriceCents, testdatagen.DefaultContractCode, iucrtTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", iucrtTestRequestedPickupDate, iucrtTestBilledCubicFeet, iucrtTestMarket)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup International Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIUCRT, iucrtTestMarket, iucrtTestBasePriceCents, testdatagen.DefaultContractCode, iucrtTestEscalationCompounded)
		twoYearsLaterPickupDate := iucrtTestRequestedPickupDate.AddDate(10, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, iucrtTestBilledCubicFeet, iucrtTestMarket)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupIntlUncratingServiceItem() models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeIUCRT,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameCubicFeetBilled,
				KeyType: models.ServiceItemParamTypeDecimal,
				Value:   fmt.Sprintf("%d", int(iucrtTestBilledCubicFeet)),
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   iucrtTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameMarketDest,
				KeyType: models.ServiceItemParamTypeString,
				Value:   iucrtTestMarket.String(),
			},
		}, nil, nil,
	)
}
