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
	ubpTestContractYearName     = "Base Period Year 1"
	ubpTestPerUnitCents         = unit.Cents(15000)
	ubpTestTotalCost            = unit.Cents(83250)
	ubpTestIsPeakPeriod         = true
	ubpTestEscalationCompounded = 1.11000
	ubpTestWeight               = unit.Pound(500)
)

var ubpTestRequestedPickupDate = time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestUbpPricer() {
	pricer := NewIntlUBPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem := suite.setupIntlUBPricerServiceItem()

		totalCost, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(ubpTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ubpTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ubpTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ubpTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ubpTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid parameters to PriceUsingParams", func() {
		paymentServiceItem := suite.setupIntlUBPricerServiceItem()

		// PerUnitCents
		paymentServiceItem.PaymentServiceItemParams[2].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to an int", models.ServiceItemParamNamePerUnitCents))

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

	suite.Run("Price validation errors", func() {

		// No contract code
		_, _, err := pricer.Price(suite.AppContextForTest(), "", ubpTestRequestedPickupDate, ubpTestWeight, ubpTestPerUnitCents.Int())
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())

		// No reference date
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, time.Time{}, ubpTestWeight, ubpTestPerUnitCents.Int())
		suite.Error(err)
		suite.Equal("ReferenceDate is required", err.Error())

		// No weight
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, ubpTestRequestedPickupDate, 0, ubpTestPerUnitCents.Int())
		suite.Error(err)
		suite.Equal(fmt.Sprintf("Weight must be at least %d pounds", minIntlWeightUB), err.Error())

		// No per unit cents
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, ubpTestRequestedPickupDate, ubpTestWeight, 0)
		suite.Error(err)
		suite.Equal("PerUnitCents is required", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) setupIntlUBPricerServiceItem() models.PaymentServiceItem {
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
		models.ReServiceCodeUBP,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   contract.Code,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   ubpTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNamePerUnitCents,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(ubpTestPerUnitCents)),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(ubpTestWeight.Int()),
			},
		}, nil, nil,
	)
}
