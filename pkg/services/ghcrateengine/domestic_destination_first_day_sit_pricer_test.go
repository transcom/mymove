package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	ddfsitTestServiceArea          = "456"
	ddfsitTestIsPeakPeriod         = false
	ddfsitTestBasePriceCents       = unit.Cents(525)
	ddfsitTestContractYearName     = "DDFSIT Test Year"
	ddfsitTestEscalationCompounded = 1.052
	ddfsitTestWeight               = unit.Pound(3300)
	ddfsitTestPriceCents           = unit.Cents(18226) // ddfsitTestBasePriceCents * (ddfsitTestWeight / 100) * ddfsitTestEscalationCompounded
)

var ddfsitTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.January, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationFirstDaySITPricer() {
	pricer := NewDomesticDestinationFirstDaySITPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticDestinationFirstDaySITServiceItem()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(ddfsitTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ddfsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ddfsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ddfsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ddfsitTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)
		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea, false)
		suite.NoError(err)
		suite.Equal(ddfsitTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, ddfsitTestRequestedPickupDate, badWeight, ddfsitTestServiceArea, false)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.Run("valid weight when minimum is disabled", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)
		weight := unit.Pound(250)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, ddfsitTestRequestedPickupDate, weight, ddfsitTestServiceArea, true)
		suite.NoError(err)
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea, false)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination first day SIT rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)
		twoYearsLaterPickupDate := ddfsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, ddfsitTestWeight, ddfsitTestServiceArea, false)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticDestinationFirstDaySITServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDDFSIT,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   ddfsitTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaDest,
				KeyType: models.ServiceItemParamTypeString,
				Value:   ddfsitTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(ddfsitTestWeight)),
			},
		},
	)
}
