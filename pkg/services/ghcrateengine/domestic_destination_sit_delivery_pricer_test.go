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
	dddsitTestSchedule                            = 1
	dddsitTestServiceArea                         = "888"
	dddsitTestIsPeakPeriod                        = false
	dddsitTestContractYearName                    = "DDDSIT Test Year"
	dddsitTestEscalationCompounded                = 1.03
	dddsitTestWeight                              = unit.Pound(2250)
	dddsitTestWeightLower                         = unit.Pound(500)
	dddsitTestWeightUpper                         = unit.Pound(4999)
	dddsitTestMilesLower                          = 251
	dddsitTestMilesUpper                          = 500
	dddsitTestDomesticOtherBasePriceCents         = unit.Cents(2518)
	dddsitTestDomesticLinehaulBasePriceMillicents = unit.Millicents(6500)
	dddsitTestDomesticServiceAreaBasePriceCents   = unit.Cents(153)
)

var dddsitTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.December, 10, 10, 22, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationSITDeliveryPricerSameZip3s() {
	suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDSH, dddsitTestServiceArea, dddsitTestIsPeakPeriod, dddsitTestDomesticServiceAreaBasePriceCents, dddsitTestEscalationCompounded)

	zipDest := "30907"
	zipSITDest := "30901" // same zip3
	distance := unit.Miles(15)

	paymentServiceItem := suite.setupDomesticDestinationSITDeliveryServiceItem(zipDest, zipSITDest, distance)
	pricer := NewDomesticDestinationSITDeliveryPricer(suite.DB())
	expectedPrice := unit.Cents(53187) // dddsitTestDomesticServiceAreaBasePriceCents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, _, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, _, err := pricer.Price(testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("PriceUsingParams but sending empty params", func(t *testing.T) {
		_, _, err := pricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.T().Run("invalid weight", func(t *testing.T) {
		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, badWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		expectedError := fmt.Sprintf("weight of %d less than the minimum", badWeight)
		suite.Contains(err.Error(), expectedError)
	})

	suite.T().Run("bad destination zip", func(t *testing.T) {
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, "123", zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid destination postal code")
	})

	suite.T().Run("bad SIT final destination zip", func(t *testing.T) {
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, "456", distance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT final destination postal code")
	})

	suite.T().Run("error from shorthaul pricer", func(t *testing.T) {
		_, _, err := pricer.Price("BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price shorthaul")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationSITDeliveryPricer50PlusMilesDiffZip3s() {
	suite.setupDomesticLinehaulPrice(dddsitTestServiceArea, dddsitTestIsPeakPeriod, dddsitTestWeightLower, dddsitTestWeightUpper, dddsitTestMilesLower, dddsitTestMilesUpper, dddsitTestDomesticLinehaulBasePriceMillicents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

	zipDest := "30907"
	zipSITDest := "36106"       // different zip3
	distance := unit.Miles(305) // > 50 miles

	paymentServiceItem := suite.setupDomesticDestinationSITDeliveryServiceItem(zipDest, zipSITDest, distance)
	pricer := NewDomesticDestinationSITDeliveryPricer(suite.DB())
	expectedPriceMillicents := unit.Millicents(45944438) // dddsitTestDomesticLinehaulBasePriceMillicents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded
	expectedPrice := expectedPriceMillicents.ToCents()

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, _, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, _, err := pricer.Price(testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("error from linehaul pricer", func(t *testing.T) {
		_, _, err := pricer.Price("BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price linehaul")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationSITDeliveryPricer50MilesOrLessDiffZip3s() {
	suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestEscalationCompounded)

	zipDest := "30907"
	zipSITDest := "29801"      // different zip3
	distance := unit.Miles(37) // <= 50 miles

	paymentServiceItem := suite.setupDomesticDestinationSITDeliveryServiceItem(zipDest, zipSITDest, distance)
	pricer := NewDomesticDestinationSITDeliveryPricer(suite.DB())
	expectedPrice := unit.Cents(58355) // dddsitTestDomesticOtherBasePriceCents * (dddsitTestWeight / 100) * dddsitTestEscalationCompounded

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, _, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, _, err := pricer.Price(testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, _, err := pricer.Price("BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination SIT delivery rate")
	})

	suite.T().Run("not finding a contract year record", func(t *testing.T) {
		twoYearsLaterPickupDate := dddsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
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
