package ghcrateengine

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	ddasitTestServiceArea          = "789"
	ddasitTestIsPeakPeriod         = false
	ddasitTestBasePriceCents       = unit.Cents(747)
	ddasitTestEscalationCompounded = 1.042
	ddasitTestWeight               = unit.Pound(4200)
	ddasitTestNumberOfDaysInSIT    = 29
	ddasitTestPriceCents           = unit.Cents(948060) // ddasitTestBasePriceCents * (ddasitTestWeight / 100) * ddasitTestEscalationCompounded * ddasitTestNumberOfDaysInSIT
)

var ddasitTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.January, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationAdditionalDaysSITPricer() {
	suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestEscalationCompounded)
	paymentServiceItem := suite.setupDomesticDestinationAdditionalDaysSITServiceItem()
	pricer := NewDomesticDestinationAdditionalDaysSITPricer(suite.DB())

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(ddasitTestPriceCents, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := pricer.Price(testdatagen.DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestIsPeakPeriod, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.NoError(err)
		suite.Equal(ddasitTestPriceCents, priceCents)
	})

	suite.T().Run("PriceUsingParams but sending empty params", func(t *testing.T) {
		_, err := pricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
		// this is the first param checked for, otherwise error doesn't matter
		suite.Equal("could not find param with key ContractCode", err.Error())
	})

	suite.T().Run("invalid weight", func(t *testing.T) {
		badWeight := unit.Pound(250)
		_, err := pricer.Price(testdatagen.DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestIsPeakPeriod, badWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, err := pricer.Price("BOGUS", ddasitTestRequestedPickupDate, ddasitTestIsPeakPeriod, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination additional days SIT rate")
	})

	suite.T().Run("not finding a contract year record", func(t *testing.T) {
		twoYearsLaterPickupDate := ddasitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, err := pricer.Price(testdatagen.DefaultContractCode, twoYearsLaterPickupDate, ddasitTestIsPeakPeriod, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationAdditionalDaysSITPricerMissingParams() {
	pricer := NewDomesticDestinationAdditionalDaysSITPricer(suite.DB())
	testData := []struct {
		testDescription string
		psiParams       []testdatagen.CreatePaymentServiceItemParams
		expectedError   string
	}{
		{
			testDescription: "not finding number of days in SIT",
			expectedError:   "could not find param with key NumberDaysSIT",
			psiParams: []testdatagen.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   testdatagen.DefaultContractCode,
				},
				{
					Key:     models.ServiceItemParamNameRequestedPickupDate,
					KeyType: models.ServiceItemParamTypeTimestamp,
					Value:   ddasitTestRequestedPickupDate.Format(TimestampParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameServiceAreaDest,
					KeyType: models.ServiceItemParamTypeString,
					Value:   ddasitTestServiceArea,
				},
				{
					Key:     models.ServiceItemParamNameWeightActual,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   "2700",
				},
				{
					Key:     models.ServiceItemParamNameWeightBilledActual,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(ddasitTestWeight)),
				},
				{
					Key:     models.ServiceItemParamNameWeightEstimated,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   "2500",
				},
				{
					Key:     models.ServiceItemParamNameZipDestAddress,
					KeyType: models.ServiceItemParamTypeString,
					Value:   "30907",
				},
			},
		},
		{
			testDescription: "not finding service area dest",
			expectedError:   "could not find param with key ServiceAreaDest",
			psiParams: []testdatagen.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   testdatagen.DefaultContractCode,
				},
				{
					Key:     models.ServiceItemParamNameRequestedPickupDate,
					KeyType: models.ServiceItemParamTypeTimestamp,
					Value:   ddasitTestRequestedPickupDate.Format(TimestampParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameWeightActual,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   "2700",
				},
				{
					Key:     models.ServiceItemParamNameWeightBilledActual,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(ddasitTestWeight)),
				},
				{
					Key:     models.ServiceItemParamNameWeightEstimated,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   "2500",
				},
			},
		},
		{
			testDescription: "not finding weight billed actual",
			expectedError:   "could not find param with key WeightBilledActual",
			psiParams: []testdatagen.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   testdatagen.DefaultContractCode,
				},
				{
					Key:     models.ServiceItemParamNameRequestedPickupDate,
					KeyType: models.ServiceItemParamTypeTimestamp,
					Value:   ddasitTestRequestedPickupDate.Format(TimestampParamFormat),
				},
			},
		},
		{
			testDescription: "not finding requested pickup date",
			expectedError:   "could not find param with key RequestedPickupDate",
			psiParams: []testdatagen.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   testdatagen.DefaultContractCode,
				},
			},
		},
	}

	for _, data := range testData {
		paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDASIT,
			data.psiParams,
		)

		suite.T().Run(data.testDescription, func(t *testing.T) {
			_, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
			suite.Error(err)
			suite.Contains(err.Error(), data.expectedError)
		})
	}
}

func (suite *GHCRateEngineServiceSuite) setupDomesticDestinationAdditionalDaysSITServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDDASIT,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeTimestamp,
				Value:   ddasitTestRequestedPickupDate.Format(TimestampParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaDest,
				KeyType: models.ServiceItemParamTypeString,
				Value:   ddasitTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameWeightActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "2700",
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(ddasitTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameWeightEstimated,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "2500",
			},
			{
				Key:     models.ServiceItemParamNameZipDestAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   "30907",
			},
			{
				Key:     models.ServiceItemParamNameNumberDaysSIT,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(ddasitTestNumberOfDaysInSIT),
			},
		},
	)
}
