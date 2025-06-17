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
	islhTestContractYearName     = testdatagen.DefaultContractYearName
	islhTestPerUnitCents         = unit.Cents(15000)
	islhTestTotalCost            = unit.Cents(349650)
	islhTestIsPeakPeriod         = true
	islhTestEscalationCompounded = 1.11000
	islhTestWeight               = unit.Pound(2100)
)

var islhTestRequestedPickupDate = time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestIntlShippingAndLinehaulPricer() {
	pricer := NewIntlShippingAndLinehaulPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem := suite.setupIntlShippingAndLinehaulServiceItem()

		totalCost, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(islhTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: islhTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(islhTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(islhTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(islhTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid parameters to PriceUsingParams", func() {
		paymentServiceItem := suite.setupIntlShippingAndLinehaulServiceItem()

		paymentServiceItem.PaymentServiceItemParams[3].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to an int", models.ServiceItemParamNameWeightBilled))

		paymentServiceItem.PaymentServiceItemParams[2].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to an int", models.ServiceItemParamNameWeightBilled))

		paymentServiceItem.PaymentServiceItemParams[1].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to a time", models.ServiceItemParamNameReferenceDate))

		paymentServiceItem.PaymentServiceItemParams[0].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to a string", models.ServiceItemParamNameContractCode))
	})

	suite.Run("Price validation errors", func() {

		// No contract code
		_, _, err := pricer.Price(suite.AppContextForTest(), "", islhTestRequestedPickupDate, islhTestWeight, islhTestPerUnitCents.Int())
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())

		// No reference date
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, time.Time{}, islhTestWeight, islhTestPerUnitCents.Int())
		suite.Error(err)
		suite.Equal("referenceDate is required", err.Error())

		// No weight
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, islhTestRequestedPickupDate, 0, islhTestPerUnitCents.Int())
		suite.Error(err)
		suite.Equal(fmt.Sprintf("weight must be at least %d", minIntlWeightHHG), err.Error())

		// No per unit cents
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, islhTestRequestedPickupDate, islhTestWeight, 0)
		suite.Error(err)
		suite.Equal("PerUnitCents is required", err.Error())

	})
}

func (suite *GHCRateEngineServiceSuite) setupIntlShippingAndLinehaulServiceItem() models.PaymentServiceItem {
	contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
	startDate := time.Date(2018, time.January, 1, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(2018, time.December, 31, 12, 0, 0, 0, time.UTC)
	testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			Contract:             contract,
			ContractID:           contract.ID,
			StartDate:            startDate,
			EndDate:              endDate,
			Escalation:           1.0,
			EscalationCompounded: 1.0,
		},
	})
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeISLH,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   contract.Code,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   islhTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNamePerUnitCents,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(islhTestPerUnitCents)),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(islhTestWeight.Int()),
			},
		}, nil, nil,
	)
}
