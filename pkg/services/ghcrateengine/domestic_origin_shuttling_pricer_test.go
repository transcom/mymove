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
	doshutTestServiceSchedule      = 2
	doshutTestBasePriceCents       = unit.Cents(353)
	doshutTestEscalationCompounded = 1.125
	doshutTestWeight               = unit.Pound(4000)
	doshutTestPriceCents           = unit.Cents(15885) // doshutTestBasePriceCents * (doshutTestWeight / 100) * doshutTestEscalationCompounded
)

var doshutTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginShuttlingPricer() {
	pricer := NewDomesticOriginShuttlingPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, testdatagen.DefaultContractCode, doshutTestEscalationCompounded)

		paymentServiceItem := suite.setupDomesticOriginShuttlingServiceItem()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(doshutTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(doshutTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(doshutTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, testdatagen.DefaultContractCode, doshutTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, doshutTestRequestedPickupDate, doshutTestWeight, doshutTestServiceSchedule)
		suite.NoError(err)
		suite.Equal(doshutTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, testdatagen.DefaultContractCode, doshutTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, testdatagen.DefaultContractCode, doshutTestEscalationCompounded)
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, doshutTestRequestedPickupDate, badWeight, doshutTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of 500")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, testdatagen.DefaultContractCode, doshutTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", doshutTestRequestedPickupDate, doshutTestWeight, doshutTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "Could not lookup Domestic Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, testdatagen.DefaultContractCode, doshutTestEscalationCompounded)
		twoYearsLaterPickupDate := doshutTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, doshutTestWeight, doshutTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "Could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticOriginShuttlingServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDOSHUT,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   doshutTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(doshutTestServiceSchedule),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(doshutTestWeight)),
			},
		},
	)
}
