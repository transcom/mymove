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
	doasitTestServiceArea          = "789"
	doasitTestIsPeakPeriod         = false
	doasitTestBasePriceCents       = unit.Cents(747)
	doasitTestContractYearName     = "DOASIT Test Year"
	doasitTestEscalationCompounded = 1.042
	doasitTestWeight               = unit.Pound(4200)
	doasitTestNumberOfDaysInSIT    = 29
	doasitTestPriceCents           = unit.Cents(948060) // doasitTestBasePriceCents * (doasitTestWeight / 100) * doasitTestEscalationCompounded * doasitTestNumberOfDaysInSIT
)

var doasitTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.January, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginAdditionalDaysSITPricer() {
	suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOASIT, doasitTestServiceArea, doasitTestIsPeakPeriod, doasitTestBasePriceCents, doasitTestContractYearName, doasitTestEscalationCompounded)
	paymentServiceItem := suite.setupDomesticOriginAdditionalDaysSITServiceItem()
	pricer := NewDomesticOriginAdditionalDaysSITPricer(suite.DB())

	suite.Run("success using PaymentServiceItemParams", func() {
		priceCents, displayParams, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(doasitTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: doasitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(doasitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(doasitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(doasitTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		priceCents, _, err := pricer.Price(testdatagen.DefaultContractCode, doasitTestRequestedPickupDate, doasitTestWeight, doasitTestServiceArea, doasitTestNumberOfDaysInSIT)
		suite.NoError(err)
		suite.Equal(doasitTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		_, _, err := pricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
		// this is the first param checked for, otherwise error doesn't matter
		suite.Equal("could not find param with key ContractCode", err.Error())
	})

	suite.Run("invalid weight", func() {
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, doasitTestRequestedPickupDate, badWeight, doasitTestServiceArea, doasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.Run("not finding a rate record", func() {
		_, _, err := pricer.Price("BOGUS", doasitTestRequestedPickupDate, doasitTestWeight, doasitTestServiceArea, doasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic origin additional days SIT rate")
	})

	suite.Run("not finding a contract year record", func() {
		twoYearsLaterPickupDate := doasitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, twoYearsLaterPickupDate, doasitTestWeight, doasitTestServiceArea, doasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginAdditionalDaysSITPricerMissingParams() {
	pricer := NewDomesticOriginAdditionalDaysSITPricer(suite.DB())
	testData := []struct {
		testDescription string
		psiParams       []testdatagen.CreatePaymentServiceItemParams
		expectedError   string
	}{
		// TODO: cannot run this test until MB-1564 is done
		//{
		//testDescription: "not finding number of days in SIT",
		//expectedError:   "could not find param with key NumberDaysSIT",
		//psiParams: []testdatagen.CreatePaymentServiceItemParams{
		//{
		//Key:     models.ServiceItemParamNameContractCode,
		//KeyType: models.ServiceItemParamTypeString,
		//Value:   testdatagen.DefaultContractCode,
		//},
		//{
		//Key:     models.ServiceItemParamNameRequestedPickupDate,
		//KeyType: models.ServiceItemParamTypeTimestamp,
		//Value:   doasitTestRequestedPickupDate.Format(TimestampParamFormat),
		//},
		//{
		//Key:     models.ServiceItemParamNameServiceAreaOrigin,
		//KeyType: models.ServiceItemParamTypeString,
		//Value:   doasitTestServiceArea,
		//},
		//{
		//Key:     models.ServiceItemParamNameWeightActual,
		//KeyType: models.ServiceItemParamTypeInteger,
		//Value:   "2700",
		//},
		//{
		//Key:     models.ServiceItemParamNameWeightBilledActual,
		//KeyType: models.ServiceItemParamTypeInteger,
		//Value:   fmt.Sprintf("%d", int(doasitTestWeight)),
		//},
		//{
		//Key:     models.ServiceItemParamNameWeightEstimated,
		//KeyType: models.ServiceItemParamTypeInteger,
		//Value:   "2500",
		//},
		//{
		//Key:     models.ServiceItemParamNameZipDestAddress,
		//KeyType: models.ServiceItemParamTypeString,
		//Value:   "30907",
		//},
		//},
		//},
		{
			testDescription: "not finding service area dest",
			expectedError:   "could not find param with key ServiceAreaOrigin",
			psiParams: []testdatagen.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   testdatagen.DefaultContractCode,
				},
				{
					Key:     models.ServiceItemParamNameRequestedPickupDate,
					KeyType: models.ServiceItemParamTypeTimestamp,
					Value:   doasitTestRequestedPickupDate.Format(TimestampParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameWeightActual,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   "2700",
				},
				{
					Key:     models.ServiceItemParamNameWeightBilledActual,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(doasitTestWeight)),
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
					Value:   doasitTestRequestedPickupDate.Format(TimestampParamFormat),
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
			models.ReServiceCodeDOASIT,
			data.psiParams,
		)

		suite.Run(data.testDescription, func() {
			_, _, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
			suite.Error(err)
			suite.Contains(err.Error(), data.expectedError)
		})
	}
}

func (suite *GHCRateEngineServiceSuite) setupDomesticOriginAdditionalDaysSITServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDOASIT,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeTimestamp,
				Value:   doasitTestRequestedPickupDate.Format(TimestampParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   doasitTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameWeightActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "2700",
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(doasitTestWeight)),
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
				Value:   strconv.Itoa(doasitTestNumberOfDaysInSIT),
			},
		},
	)
}
