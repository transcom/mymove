package ghcrateengine

import (
	"fmt"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	dddsitTestSchedule                    = 1
	dddsitTestServiceArea                 = "888"
	dddsitTestIsPeakPeriod                = false
	dddsitTestDomesticOtherBasePriceCents = unit.Cents(2518)
	dddsitTestEscalationCompounded        = 1.03
	dddsitTestWeight                      = unit.Pound(2250)
	dddsitTestWeightLower                 = unit.Pound(500)
	dddsitTestWeightUpper                 = unit.Pound(4999)
	dddsitTestMilesLower                  = 251
	dddsitTestMilesUpper                  = 500
	dddsitTestBasePriceMillicents         = unit.Millicents(6500)
)

var dddsitTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.December, 10, 10, 22, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationSITDeliveryPricer50MilesOrLess() {
	suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestEscalationCompounded)

	zipDest := "30907"
	zipSITDest := "30901" // same zip3, but different zip3 should work the same
	distance := unit.Miles(15)

	paymentServiceItem := suite.setupDomesticDestinationSITDeliveryServiceItem(zipDest, zipSITDest, distance)
	pricer := NewDomesticDestinationSITDeliveryPricer(suite.DB())
	expectedPrice := unit.Cents(58355) // dddsitTestBasePriceCents * (dddsitTestWeight / 100) * dddsitTestEscalationCompounded

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams, same zip3s", func(t *testing.T) {
		priceCents, err := pricer.Price(testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams, different zip3s", func(t *testing.T) {
		zipSITDestDiffZip3 := "29841" // Should get the same answer with different zip3 and same mileage
		priceCents, err := pricer.Price(testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDestDiffZip3, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("PriceUsingParams but sending empty params", func(t *testing.T) {
		_, err := pricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.T().Run("invalid weight", func(t *testing.T) {
		badWeight := unit.Pound(250)
		_, err := pricer.Price(testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, badWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, err := pricer.Price("BOGUS", dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination SIT delivery rate")
	})

	suite.T().Run("not finding a contract year record", func(t *testing.T) {
		twoYearsLaterPickupDate := dddsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, err := pricer.Price(testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationSITDeliveryPricerMore50PlusMilesDiffZip3s() {
	suite.setupDomesticLinehaulPrice(dddsitTestServiceArea, dddsitTestIsPeakPeriod, dddsitTestWeightLower, dddsitTestWeightUpper, dddsitTestMilesLower, dddsitTestMilesUpper, dddsitTestBasePriceMillicents, dddsitTestEscalationCompounded)

	zipDest := "30907"
	zipSITDest := "36106"       // different zip3
	distance := unit.Miles(305) // more than 50 miles

	paymentServiceItem := suite.setupDomesticDestinationSITDeliveryServiceItem(zipDest, zipSITDest, distance)
	pricer := NewDomesticDestinationSITDeliveryPricer(suite.DB())
	expectedPriceMillicents := unit.Millicents(45944438) // dddsitTestBasePriceMillicents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded
	expectedPrice := expectedPriceMillicents.ToCents()

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := pricer.Price(testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("error from linehaul pricer", func(t *testing.T) {
		_, err := pricer.Price("BOGUS", dddsitTestRequestedPickupDate, dddsitTestIsPeakPeriod, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price linehaul")
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticDestinationSITDeliveryServiceItem(zipDest string, zipSITDest string, distance unit.Miles) models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDDDSIT,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dddsitTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaDest,
				KeyType: models.ServiceItemParamTypeString,
				Value:   dddsitTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameSITScheduleDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", dddsitTestSchedule),
			},
			{
				Key:     models.ServiceItemParamNameWeightActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dddsitTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dddsitTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameWeightEstimated,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "2400",
			},
			{
				Key:     models.ServiceItemParamNameZipDestAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   zipDest,
			},
			{
				Key:     models.ServiceItemParamNameZipSITDestHHGFinalAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   zipSITDest,
			},
			{
				Key:     models.ServiceItemParamNameDistanceZipSITDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(distance)),
			},
		},
	)
}
