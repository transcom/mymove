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
	ddasitTestServiceArea          = "789"
	ddasitTestIsPeakPeriod         = false
	ddasitTestBasePriceCents       = unit.Cents(747)
	ddasitTestContractYearName     = "DDASIT Test Year"
	ddasitTestEscalationCompounded = 1.042
	ddasitTestWeight               = unit.Pound(4200)
	ddasitTestNumberOfDaysInSIT    = 29
	ddasitTestPriceCents           = unit.Cents(948060) // ddasitTestBasePriceCents * (ddasitTestWeight / 100) * ddasitTestEscalationCompounded * ddasitTestNumberOfDaysInSIT
)

var ddasitTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.January, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationAdditionalDaysSITPricer() {
	pricer := NewDomesticDestinationAdditionalDaysSITPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticDestinationAdditionalDaysSITServiceItem()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.TestAppContext(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(ddasitTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ddasitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ddasitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ddasitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ddasitTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.TestAppContext(), testdatagen.DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.NoError(err)
		suite.Equal(ddasitTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.TestAppContext(), models.PaymentServiceItemParams{})
		suite.Error(err)
		// this is the first param checked for, otherwise error doesn't matter
		suite.Equal("could not find param with key ContractCode", err.Error())
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(suite.TestAppContext(), testdatagen.DefaultContractCode, ddasitTestRequestedPickupDate, badWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)
		_, _, err := pricer.Price(suite.TestAppContext(), "BOGUS", ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination additional days SIT rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)
		twoYearsLaterPickupDate := ddasitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(suite.TestAppContext(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationAdditionalDaysSITPricerMissingParams() {
	pricer := NewDomesticDestinationAdditionalDaysSITPricer()
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
		//Value:   ddasitTestRequestedPickupDate.Format(TimestampParamFormat),
		//},
		//{
		//Key:     models.ServiceItemParamNameServiceAreaDest,
		//KeyType: models.ServiceItemParamTypeString,
		//Value:   ddasitTestServiceArea,
		//},
		//{
		//Key:     models.ServiceItemParamNameWeightOriginal,
		//KeyType: models.ServiceItemParamTypeInteger,
		//Value:   "2700",
		//},
		//{
		//Key:     models.ServiceItemParamNameWeightBilledActual,
		//KeyType: models.ServiceItemParamTypeInteger,
		//Value:   fmt.Sprintf("%d", int(ddasitTestWeight)),
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
					Key:     models.ServiceItemParamNameWeightOriginal,
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
		suite.Run(data.testDescription, func() {
			paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItemWithParams(
				suite.DB(),
				models.ReServiceCodeDDASIT,
				data.psiParams,
			)

			_, _, err := pricer.PriceUsingParams(suite.TestAppContext(), paymentServiceItem.PaymentServiceItemParams)
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
				Key:     models.ServiceItemParamNameWeightOriginal,
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
