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
	dupkTestEscalationCompounded = 1.2311
	dupkTestIsPeakPeriod         = false
	dupkTestWeight               = unit.Pound(3600)
	dupkTestServicesScheduleDest = 1
	dupkTestContractYearName     = "DUPK Test Year"
	dupkTestBasePriceCents       = unit.Cents(123)
	dupkTestPriceCents           = unit.Cents(5436)
)

var dupkTestRequestedPickupDate = time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

func (suite *GHCRateEngineServiceSuite) TestDomesticUnpackPricer() {
	pricer := NewDomesticUnpackPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDUPK, dupkTestServicesScheduleDest, dupkTestIsPeakPeriod, dupkTestBasePriceCents, dupkTestContractYearName, dupkTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticUnpackServiceItem()

		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(dupkTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dupkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dupkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dupkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dupkTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid parameters to PriceUsingParams", func() {
		paymentServiceItem := suite.setupDomesticUnpackServiceItem()

		// Setting each param's type to something incorrect should trigger an error.

		// WeightBilled
		paymentServiceItem.PaymentServiceItemParams[3].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to an int", models.ServiceItemParamNameWeightBilled))

		// ServicesScheduleDest
		paymentServiceItem.PaymentServiceItemParams[2].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to an int", models.ServiceItemParamNameServicesScheduleDest))

		// ReferenceDate
		paymentServiceItem.PaymentServiceItemParams[1].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to a time", models.ServiceItemParamNameReferenceDate))

		// ContractCode
		paymentServiceItem.PaymentServiceItemParams[0].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to a string", models.ServiceItemParamNameContractCode))
	})

	suite.Run("successfully finds dom unpack price for ppm with weight < 500 lbs with PriceUsingParams method", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDUPK, dupkTestServicesScheduleDest, dupkTestIsPeakPeriod, dupkTestBasePriceCents, dupkTestContractYearName, dupkTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticUnpackServiceItem()
		params := paymentServiceItem.PaymentServiceItemParams
		params[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType = models.MTOShipmentTypePPM
		weightBilledIndex := 3

		params[weightBilledIndex].Value = "500"
		basePriceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)

		params[weightBilledIndex].Value = "250"
		halfPriceCents, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(basePriceCents/2, halfPriceCents)

		params[weightBilledIndex].Value = "100"
		fifthPriceCents, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(basePriceCents/5, fifthPriceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dupkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dupkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dupkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dupkTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticUnpackServiceItem() models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDUPK,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dupkTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dupkTestServicesScheduleDest),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dupkTestWeight.Int()),
			},
		}, nil, nil,
	)
}
