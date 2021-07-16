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
	dopsitTestSchedule                            = 3
	dopsitTestServiceArea                         = "012"
	dopsitTestIsPeakPeriod                        = true
	dopsitTestContractYearName                    = "DOPSIT Test Year"
	dopsitTestEscalationCompounded                = 1.0445
	dopsitTestWeight                              = unit.Pound(4555)
	dopsitTestWeightLower                         = unit.Pound(4000)
	dopsitTestWeightUpper                         = unit.Pound(4999)
	dopsitTestMilesLower                          = 51
	dopsitTestMilesUpper                          = 250
	dopsitTestDomesticOtherBasePriceCents         = unit.Cents(2810)
	dopsitTestDomesticLinehaulBasePriceMillicents = unit.Millicents(4455)
	dopsitTestDomesticServiceAreaBasePriceCents   = unit.Cents(223)
)

var dopsitTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.July, 5, 10, 22, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginSITPickupPricerSameZip3s() {
	suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDSH, dopsitTestServiceArea, dopsitTestIsPeakPeriod, dopsitTestDomesticServiceAreaBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)

	zipOriginal := "29201"
	zipActual := "29212" // same zip3
	distance := unit.Miles(12)

	paymentServiceItem := suite.setupDomesticOriginSITPickupServiceItem(zipOriginal, zipActual, distance)
	pricer := NewDomesticOriginSITPickupPricer(suite.DB())
	expectedPrice := unit.Cents(127316) // dopsitTestDomesticServiceAreaBasePriceCents * (dopsitTestWeight / 100) * distance * dopsitTestEscalationCompounded

	suite.Run("success using PaymentServiceItemParams", func() {
		priceCents, displayParams, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dopsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dopsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dopsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dopsitTestDomesticServiceAreaBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		priceCents, displayParams, err := pricer.Price(testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dopsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dopsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dopsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dopsitTestDomesticServiceAreaBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		_, _, err := pricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid weight", func() {
		badWeight := unit.Pound(333)
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, badWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.Error(err)
		expectedError := fmt.Sprintf("weight of %d less than the minimum", badWeight)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("bad original zip", func() {
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, "7891", zipActual, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT origin original postal code")
	})

	suite.Run("bad actual zip", func() {
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, "12", distance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT origin actual postal code")
	})

	suite.Run("error from shorthaul pricer", func() {
		_, _, err := pricer.Price("BOGUS", dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price shorthaul")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginSITPickupPricer50PlusMilesDiffZip3s() {
	suite.setupDomesticLinehaulPrice(dopsitTestServiceArea, dopsitTestIsPeakPeriod, dopsitTestWeightLower, dopsitTestWeightUpper, dopsitTestMilesLower, dopsitTestMilesUpper, dopsitTestDomesticLinehaulBasePriceMillicents, dopsitTestContractYearName, dopsitTestEscalationCompounded)

	zipOriginal := "29201"
	zipActual := "30907"       // different zip3
	distance := unit.Miles(77) // > 50 miles

	paymentServiceItem := suite.setupDomesticOriginSITPickupServiceItem(zipOriginal, zipActual, distance)
	pricer := NewDomesticOriginSITPickupPricer(suite.DB())
	expectedPriceMillicents := unit.Millicents(16320568) // dopsitTestDomesticLinehaulBasePriceMillicents * (dopsitTestWeight / 100) * distance * dopsitTestEscalationCompounded
	expectedPrice := expectedPriceMillicents.ToCents()

	suite.Run("success using PaymentServiceItemParams", func() {
		priceCents, displayParams, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dopsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dopsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dopsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatFloat(dopsitTestDomesticLinehaulBasePriceMillicents.ToDollarFloatNoRound(), 3)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		priceCents, _, err := pricer.Price(testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.Run("error from linehaul pricer", func() {
		_, _, err := pricer.Price("BOGUS", dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price linehaul")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginSITPickupPricer50MilesOrLessDiffZip3s() {
	suite.setupDomesticOtherPrice(models.ReServiceCodeDOPSIT, dopsitTestSchedule, dopsitTestIsPeakPeriod, dopsitTestDomesticOtherBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)

	zipOriginal := "29201"
	zipActual := "29123"       // different zip3
	distance := unit.Miles(23) // <= 50 miles

	paymentServiceItem := suite.setupDomesticOriginSITPickupServiceItem(zipOriginal, zipActual, distance)
	pricer := NewDomesticOriginSITPickupPricer(suite.DB())
	expectedPrice := unit.Cents(133691) // dopsitTestDomesticOtherBasePriceCents * (dopsitTestWeight / 100) * dopsitTestEscalationCompounded

	suite.Run("success using PaymentServiceItemParams", func() {
		priceCents, displayParams, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dopsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dopsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dopsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dopsitTestDomesticOtherBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		priceCents, _, err := pricer.Price(testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.Run("not finding a rate record", func() {
		_, _, err := pricer.Price("BOGUS", dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic origin SIT pickup rate")
	})

	suite.Run("not finding a contract year record", func() {
		twoYearsLaterPickupDate := dopsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticOriginSITPickupServiceItem(zipOriginal string, zipActual string, distance unit.Miles) models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDOPSIT,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dopsitTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   dopsitTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameSITScheduleOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", dopsitTestSchedule),
			},
			{
				Key:     models.ServiceItemParamNameWeightActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dopsitTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dopsitTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameWeightEstimated,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "2400",
			},
			{
				Key:     models.ServiceItemParamNameZipSITOriginHHGOriginalAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   zipOriginal,
			},
			{
				Key:     models.ServiceItemParamNameZipSITOriginHHGActualAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   zipActual,
			},
			{
				Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(distance)),
			},
		},
	)
}
