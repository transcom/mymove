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
	dofsitTestServiceArea          = "123"
	dofsitTestIsPeakPeriod         = true
	dofsitTestBasePriceCents       = unit.Cents(353)
	dofsitTestContractYearName     = "DOFSIT Test Year"
	dofsitTestEscalationCompounded = 1.125
	dofsitTestWeight               = unit.Pound(4000)
	dofsitTestPriceCents           = unit.Cents(15885) // dofsitTestBasePriceCents * (dofsitTestWeight / 100) * dofsitTestEscalationCompounded
)

var dofsitTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginFirstDaySITPricer() {
	pricer := NewDomesticOriginFirstDaySITPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, dofsitTestServiceArea, dofsitTestIsPeakPeriod, dofsitTestBasePriceCents, dofsitTestContractYearName, dofsitTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticOriginFirstDaySITServiceItem()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.TestAppContext(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(dofsitTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dofsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dofsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dofsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dofsitTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, dofsitTestServiceArea, dofsitTestIsPeakPeriod, dofsitTestBasePriceCents, dofsitTestContractYearName, dofsitTestEscalationCompounded)
		priceCents, _, err := pricer.Price(suite.TestAppContext(), testdatagen.DefaultContractCode, dofsitTestRequestedPickupDate, dofsitTestWeight, dofsitTestServiceArea)
		suite.NoError(err)
		suite.Equal(dofsitTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, dofsitTestServiceArea, dofsitTestIsPeakPeriod, dofsitTestBasePriceCents, dofsitTestContractYearName, dofsitTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.TestAppContext(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, dofsitTestServiceArea, dofsitTestIsPeakPeriod, dofsitTestBasePriceCents, dofsitTestContractYearName, dofsitTestEscalationCompounded)
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(suite.TestAppContext(), testdatagen.DefaultContractCode, dofsitTestRequestedPickupDate, badWeight, dofsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, dofsitTestServiceArea, dofsitTestIsPeakPeriod, dofsitTestBasePriceCents, dofsitTestContractYearName, dofsitTestEscalationCompounded)
		_, _, err := pricer.Price(suite.TestAppContext(), "BOGUS", dofsitTestRequestedPickupDate, dofsitTestWeight, dofsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic origin first day SIT rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, dofsitTestServiceArea, dofsitTestIsPeakPeriod, dofsitTestBasePriceCents, dofsitTestContractYearName, dofsitTestEscalationCompounded)
		twoYearsLaterPickupDate := dofsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(suite.TestAppContext(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dofsitTestWeight, dofsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticOriginFirstDaySITServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDOFSIT,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dofsitTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   dofsitTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameWeightActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "1400",
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dofsitTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameWeightEstimated,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "1400",
			},
			{
				Key:     models.ServiceItemParamNameZipPickupAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   "90210",
			},
		},
	)
}
