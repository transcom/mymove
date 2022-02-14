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
	pricer := NewDomesticOriginAdditionalDaysSITPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOASIT, doasitTestServiceArea, doasitTestIsPeakPeriod, doasitTestBasePriceCents, doasitTestContractYearName, doasitTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticOriginAdditionalDaysSITServiceItem()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
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
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOASIT, doasitTestServiceArea, doasitTestIsPeakPeriod, doasitTestBasePriceCents, doasitTestContractYearName, doasitTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, doasitTestRequestedPickupDate, doasitTestWeight, doasitTestServiceArea, doasitTestNumberOfDaysInSIT)
		suite.NoError(err)
		suite.Equal(doasitTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOASIT, doasitTestServiceArea, doasitTestIsPeakPeriod, doasitTestBasePriceCents, doasitTestContractYearName, doasitTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
		// this is the first param checked for, otherwise error doesn't matter
		suite.Equal("could not find param with key ContractCode", err.Error())
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOASIT, doasitTestServiceArea, doasitTestIsPeakPeriod, doasitTestBasePriceCents, doasitTestContractYearName, doasitTestEscalationCompounded)
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, doasitTestRequestedPickupDate, badWeight, doasitTestServiceArea, doasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOASIT, doasitTestServiceArea, doasitTestIsPeakPeriod, doasitTestBasePriceCents, doasitTestContractYearName, doasitTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", doasitTestRequestedPickupDate, doasitTestWeight, doasitTestServiceArea, doasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic origin additional days SIT rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOASIT, doasitTestServiceArea, doasitTestIsPeakPeriod, doasitTestBasePriceCents, doasitTestContractYearName, doasitTestEscalationCompounded)
		twoYearsLaterPickupDate := doasitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, doasitTestWeight, doasitTestServiceArea, doasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginAdditionalDaysSITPricerMissingParams() {
	pricer := NewDomesticOriginAdditionalDaysSITPricer()
	testData := []struct {
		testDescription string
		psiParams       []testdatagen.CreatePaymentServiceItemParams
		expectedError   string
	}{
		{
			testDescription: "not finding weight billed",
			expectedError:   "could not find param with key WeightBilled",
			psiParams: []testdatagen.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   testdatagen.DefaultContractCode,
				},
				{
					Key:     models.ServiceItemParamNameNumberDaysSIT,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(doasitTestNumberOfDaysInSIT)),
				},
				{
					Key:     models.ServiceItemParamNameReferenceDate,
					KeyType: models.ServiceItemParamTypeDate,
					Value:   doasitTestRequestedPickupDate.Format(DateParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameServiceAreaOrigin,
					KeyType: models.ServiceItemParamTypeString,
					Value:   doasitTestServiceArea,
				},
			},
		},
		{
			testDescription: "not finding service area origin",
			expectedError:   "could not find param with key ServiceAreaOrigin",
			psiParams: []testdatagen.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   testdatagen.DefaultContractCode,
				},
				{
					Key:     models.ServiceItemParamNameNumberDaysSIT,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(doasitTestNumberOfDaysInSIT)),
				},
				{
					Key:     models.ServiceItemParamNameReferenceDate,
					KeyType: models.ServiceItemParamTypeDate,
					Value:   doasitTestRequestedPickupDate.Format(DateParamFormat),
				},
			},
		},
		{
			testDescription: "not finding reference date",
			expectedError:   "could not find param with key ReferenceDate",
			psiParams: []testdatagen.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   testdatagen.DefaultContractCode,
				},
				{
					Key:     models.ServiceItemParamNameNumberDaysSIT,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(doasitTestNumberOfDaysInSIT)),
				},
			},
		},
		{
			testDescription: "not finding number of days in SIT",
			expectedError:   "could not find param with key NumberDaysSIT",
			psiParams: []testdatagen.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   testdatagen.DefaultContractCode,
				},
			},
		},
		{
			testDescription: "not finding contract code",
			expectedError:   "could not find param with key ContractCode",
			psiParams:       []testdatagen.CreatePaymentServiceItemParams{},
		},
	}

	for _, data := range testData {
		suite.Run(data.testDescription, func() {
			paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItemWithParams(
				suite.DB(),
				models.ReServiceCodeDOASIT,
				data.psiParams,
			)

			_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
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
				Key:     models.ServiceItemParamNameNumberDaysSIT,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(doasitTestNumberOfDaysInSIT),
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   doasitTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   doasitTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(doasitTestWeight)),
			},
		},
	)
}
