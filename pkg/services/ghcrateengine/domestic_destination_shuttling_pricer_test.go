package ghcrateengine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	testServiceSchedule      = 2
	testBasePriceCents       = unit.Cents(353)
	testContractYearName     = "Test Year 1"
	testEscalationCompounded = 1.125
	testWeight               = unit.Pound(4000)
	testPriceCents           = unit.Cents(15885) // testBasePriceCents * (testWeight / 100) * testEscalationCompounded
)

var testRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationShuttlingPricer() {
	suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, testServiceSchedule, testBasePriceCents, testdatagen.DefaultContractCode, testEscalationCompounded)

	paymentServiceItem := suite.setupDomesticDestinationShuttlingServiceItem()
	pricer := NewDomesticDestinationShuttlingPricer(suite.DB())

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

func (suite *GHCRateEngineServiceSuite) setupDomesticDestinationShuttlingServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDDSHUT,
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
				Key:     models.ServiceItemParamNameServicesScheduleDest,
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
