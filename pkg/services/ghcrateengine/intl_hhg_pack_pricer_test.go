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
	ihpkTestContractYearName     = testdatagen.DefaultContractYearName
	ihpkTestPerUnitCents         = unit.Cents(15000)
	ihpkTestTotalCost            = unit.Cents(349650)
	ihpkTestIsPeakPeriod         = true
	ihpkTestEscalationCompounded = 1.11000
	ihpkTestWeight               = unit.Pound(2100)
	ihpkTestPriceCents           = unit.Cents(193064)
)

var ihpkTestRequestedPickupDate = time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestIntlHHGPackPricer() {
	pricer := NewIntlHHGPackPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem, _ := suite.setupIntlPackServiceItem(models.ReServiceCodeIHPK)

		totalCost, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(ihpkTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ihpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ihpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ihpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ihpkTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid parameters to PriceUsingParams", func() {
		paymentServiceItem, _ := suite.setupIntlPackServiceItem(models.ReServiceCodeIHPK)

		// WeightBilled
		paymentServiceItem.PaymentServiceItemParams[3].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to an int", models.ServiceItemParamNameWeightBilled))

		// PerUnitCents
		paymentServiceItem.PaymentServiceItemParams[2].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
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
}

func (suite *GHCRateEngineServiceSuite) setupIntlPackServiceItem(code models.ReServiceCode) (models.PaymentServiceItem, models.ReContract) {
	contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
	startDate := time.Date(2018, time.January, 1, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(2018, time.December, 31, 12, 0, 0, 0, time.UTC)
	testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			Contract:             contract,
			ContractID:           contract.ID,
			StartDate:            startDate,
			EndDate:              endDate,
			Escalation:           1.11,
			EscalationCompounded: 1.11,
		},
	})
	availableToPrimeAt := time.Date(2018, time.September, 14, 0, 0, 0, 0, time.UTC)
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		code,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   contract.Code,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   ihpkTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNamePerUnitCents,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(ihpkTestPerUnitCents)),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(ihpkTestWeight.Int()),
			},
			{
				Key:     models.ServiceItemParamNameNTSPackingFactor,
				KeyType: models.ServiceItemParamTypeDecimal,
				// Note for a future dev
				// If you are looking at line, the packing factor is probably breaking something
				// If that is the case, this is supposed to match the non-truncated db value
				// It shouldn't really ever change, but if it did just update it here to match
				Value: strconv.FormatFloat(1.45, 'f', -1, 64),
			},
		}, []factory.Customization{{
			// Available to prime is used to fetch market factors for moves
			// The market factor can only be fetched if it's available to the Prime
			// And if it isn't available to the Prime, then we shouldn't be processing the creation
			// of this payment request
			Model: models.Move{
				AvailableToPrimeAt: &availableToPrimeAt,
			},
		},
		}, nil,
	), contract
}
