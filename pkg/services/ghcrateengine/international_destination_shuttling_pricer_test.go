package ghcrateengine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	idshutTestServiceSchedule      = 2
	idshutTestBasePriceCents       = unit.Cents(353)
	idshutTestEscalationCompounded = 1.125
	idshutTestWeight               = unit.Pound(4000)
	idshutTestPriceCents           = unit.Cents(15880)
)

var idshutTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestInternationalDestinationShuttlingPricer() {
	pricer := NewInternationalDestinationShuttlingPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIDSHUT, idshutTestServiceSchedule, idshutTestBasePriceCents, testdatagen.DefaultContractCode, idshutTestEscalationCompounded)

		paymentServiceItem := suite.setupInternationalDestinationShuttlingServiceItem()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(idshutTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(idshutTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(idshutTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIDSHUT, idshutTestServiceSchedule, idshutTestBasePriceCents, testdatagen.DefaultContractCode, idshutTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, idshutTestRequestedPickupDate, idshutTestWeight, idshutTestServiceSchedule)
		suite.NoError(err)
		suite.Equal(idshutTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIDSHUT, idshutTestServiceSchedule, idshutTestBasePriceCents, testdatagen.DefaultContractCode, idshutTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid weight", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIDSHUT, idshutTestServiceSchedule, idshutTestBasePriceCents, testdatagen.DefaultContractCode, idshutTestEscalationCompounded)
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, idshutTestRequestedPickupDate, badWeight, idshutTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of 500")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIDSHUT, idshutTestServiceSchedule, idshutTestBasePriceCents, testdatagen.DefaultContractCode, idshutTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", idshutTestRequestedPickupDate, idshutTestWeight, idshutTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup International Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIDSHUT, idshutTestServiceSchedule, idshutTestBasePriceCents, testdatagen.DefaultContractCode, idshutTestEscalationCompounded)
		twoYearsLaterPickupDate := idshutTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, idshutTestWeight, idshutTestServiceSchedule)
		suite.Error(err)

		suite.Contains(err.Error(), "could not calculate escalated price")
	})
}

func (suite *GHCRateEngineServiceSuite) setupInternationalDestinationShuttlingServiceItem() models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeIDSHUT,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   idshutTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(idshutTestServiceSchedule),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(idshutTestWeight)),
			},
		}, nil, nil,
	)
}
