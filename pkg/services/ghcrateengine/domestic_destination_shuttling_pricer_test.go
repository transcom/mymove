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
	ddshutTestServiceSchedule      = 2
	ddshutTestBasePriceCents       = unit.Cents(353)
	ddshutTestEscalationCompounded = 1.125
	ddshutTestWeight               = unit.Pound(4000)
	ddshutTestPriceCents           = unit.Cents(15885) // ddshutTestBasePriceCents * (ddshutTestWeight / 100) * ddshutTestEscalationCompounded
)

var ddshutTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationShuttlingPricer() {
	pricer := NewDomesticDestinationShuttlingPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule, ddshutTestBasePriceCents, testdatagen.DefaultContractCode, ddshutTestEscalationCompounded)

		paymentServiceItem := suite.setupDomesticDestinationShuttlingServiceItem()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(ddshutTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ddshutTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ddshutTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule, ddshutTestBasePriceCents, testdatagen.DefaultContractCode, ddshutTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, ddshutTestRequestedPickupDate, ddshutTestWeight, ddshutTestServiceSchedule)
		suite.NoError(err)
		suite.Equal(ddshutTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule, ddshutTestBasePriceCents, testdatagen.DefaultContractCode, ddshutTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule, ddshutTestBasePriceCents, testdatagen.DefaultContractCode, ddshutTestEscalationCompounded)
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, ddshutTestRequestedPickupDate, badWeight, ddshutTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of 500")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule, ddshutTestBasePriceCents, testdatagen.DefaultContractCode, ddshutTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", ddshutTestRequestedPickupDate, ddshutTestWeight, ddshutTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "Could not lookup Domestic Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule, ddshutTestBasePriceCents, testdatagen.DefaultContractCode, ddshutTestEscalationCompounded)
		twoYearsLaterPickupDate := ddshutTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, ddshutTestWeight, ddshutTestServiceSchedule)
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
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   ddshutTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(ddshutTestServiceSchedule),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(ddshutTestWeight)),
			},
		},
	)
}
