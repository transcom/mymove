package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	dddsitTestSchedule                            = 1
	dddsitTestServiceArea                         = "888"
	dddsitTestIsPeakPeriod                        = false
	dddsitTestContractYearName                    = testdatagen.DefaultContractYearName
	dddsitTestEscalationCompounded                = 1.11000
	dddsitTestWeight                              = unit.Pound(2250)
	dddsitTestWeightLower                         = unit.Pound(500)
	dddsitTestWeightUpper                         = unit.Pound(4999)
	dddsitTestMilesLower                          = 251
	dddsitTestMilesUpper                          = 500
	dddsitTestDomesticOtherBasePriceCents         = unit.Cents(21796)
	dddsitTestDomesticLinehaulBasePriceMillicents = unit.Millicents(237900)
	dddsitTestDomesticServiceAreaBasePriceCents   = unit.Cents(153)
)

var dddsitTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.December, 10, 10, 22, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationSITDeliveryPricerSameZip3s() {
	zipDest := "30907"
	zipSITDest := "30901" // same zip3
	distance := unit.Miles(37)

	pricer := NewDomesticDestinationSITDeliveryPricer()
	expectedPrice := unit.Cents(544365) // dddsitTestDomesticServiceAreaBasePriceCents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		paymentServiceItem := suite.setupDomesticDestinationSITDeliveryServiceItem(zipDest, zipSITDest, distance)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dddsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dddsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dddsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dddsitTestDomesticOtherBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("PriceUsingParams throw error without contract code", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		badWeight := unit.Pound(250)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, badWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		expectedError := fmt.Sprintf("weight of %d less than the minimum", badWeight)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("bad destination zip", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, "123", zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT original destination postal code")
	})

	suite.Run("bad SIT final destination zip", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, "456", distance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT final destination postal code")
	})

	suite.Run("error from shorthaul pricer", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination SIT delivery rate")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationSITDeliveryPricer50PlusMilesDiffZip3s() {
	zipDest := "30907"
	zipSITDest := "36106"       // different zip3
	distance := unit.Miles(305) // > 50 miles

	pricer := NewDomesticDestinationSITDeliveryPricer()
	expectedPrice := unit.Cents(1812386)

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticLinehaulPrice(dddsitTestServiceArea, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		paymentServiceItem := suite.setupDomesticDestinationSITDeliveryServiceItem(zipDest, zipSITDest, distance)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dddsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dddsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dddsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatFloat(dddsitTestDomesticLinehaulBasePriceMillicents.ToDollarFloatNoRound(), 3)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupDomesticLinehaulPrice(dddsitTestServiceArea, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.Run("error from linehaul pricer", func() {
		suite.setupDomesticLinehaulPrice(dddsitTestServiceArea, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price linehaul")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticDestinationSITDeliveryPricer50MilesOrLessDiffZip3s() {
	zipDest := "30907"
	zipSITDest := "29801"      // different zip3
	distance := unit.Miles(37) // <= 50 miles

	pricer := NewDomesticDestinationSITDeliveryPricer()
	expectedPrice := unit.Cents(544365)

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		paymentServiceItem := suite.setupDomesticDestinationSITDeliveryServiceItem(zipDest, zipSITDest, distance)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dddsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dddsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dddsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dddsitTestDomesticOtherBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination SIT delivery rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestContractYearName, dddsitTestEscalationCompounded)
		twoYearsLaterPickupDate := dddsitTestRequestedPickupDate.AddDate(10, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, zipDest, zipSITDest, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticDestinationSITDeliveryServiceItem(zipDest string, zipSITDest string, distance unit.Miles) models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDDDSIT,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameDistanceZipSITDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(distance)),
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dddsitTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameSITServiceAreaDest,
				KeyType: models.ServiceItemParamTypeString,
				Value:   dddsitTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameSITScheduleDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", dddsitTestSchedule),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dddsitTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameZipSITDestHHGOriginalAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   zipDest,
			},
			{
				Key:     models.ServiceItemParamNameZipSITDestHHGFinalAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   zipSITDest,
			},
		}, nil, nil,
	)
}
