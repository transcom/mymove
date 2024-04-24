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
	zipOriginal := "29201"
	zipActual := "29212"       // same zip3
	distance := unit.Miles(12) // distance will follow pricer logic for moves under 50 miles

	pricer := NewDomesticOriginSITPickupPricer()
	expectedPrice := unit.Cents(10613) // dopsitTestDomesticServiceAreaBasePriceCents * (dopsitTestWeight / 100) * distance * dopsitTestEscalationCompounded

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDOPSIT, dopsitTestSchedule, dopsitTestIsPeakPeriod, dopsitTestDomesticServiceAreaBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)

		paymentServiceItem := suite.setupDomesticOriginSITPickupServiceItem(zipOriginal, zipActual, distance)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
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
		suite.setupDomesticOtherPrice(models.ReServiceCodeDOPSIT, dopsitTestSchedule, dopsitTestIsPeakPeriod, dopsitTestDomesticServiceAreaBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)

		priceCents, displayParams, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
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
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDSH, dopsitTestServiceArea, dopsitTestIsPeakPeriod, dopsitTestDomesticServiceAreaBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDSH, dopsitTestServiceArea, dopsitTestIsPeakPeriod, dopsitTestDomesticServiceAreaBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)
		badWeight := unit.Pound(333)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, badWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.Error(err)
		expectedError := fmt.Sprintf("weight of %d less than the minimum", badWeight)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("bad original zip", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDSH, dopsitTestServiceArea, dopsitTestIsPeakPeriod, dopsitTestDomesticServiceAreaBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, "7891", zipActual, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT origin original postal code")
	})

	suite.Run("bad actual zip", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDSH, dopsitTestServiceArea, dopsitTestIsPeakPeriod, dopsitTestDomesticServiceAreaBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, "12", distance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT origin actual postal code")
	})

	suite.Run("error from domestic origin SIT pickup pricer", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDSH, dopsitTestServiceArea, dopsitTestIsPeakPeriod, dopsitTestDomesticServiceAreaBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic origin SIT pickup rate")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginSITPickupPricer50PlusMilesDiffZip3s() {
	zipOriginal := "29201"
	zipActual := "30907"       // different zip3
	distance := unit.Miles(77) // > 50 miles

	pricer := NewDomesticOriginSITPickupPricer()
	expectedPrice := unit.Cents(16485)

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticLinehaulPrice(dopsitTestServiceArea, dopsitTestIsPeakPeriod, dopsitTestWeightLower, dopsitTestWeightUpper, dopsitTestMilesLower, dopsitTestMilesUpper, dopsitTestDomesticLinehaulBasePriceMillicents, dopsitTestContractYearName, dopsitTestEscalationCompounded)

		paymentServiceItem := suite.setupDomesticOriginSITPickupServiceItem(zipOriginal, zipActual, distance)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
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
		suite.setupDomesticLinehaulPrice(dopsitTestServiceArea, dopsitTestIsPeakPeriod, dopsitTestWeightLower, dopsitTestWeightUpper, dopsitTestMilesLower, dopsitTestMilesUpper, dopsitTestDomesticLinehaulBasePriceMillicents, dopsitTestContractYearName, dopsitTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.Run("error from linehaul pricer", func() {
		suite.setupDomesticLinehaulPrice(dopsitTestServiceArea, dopsitTestIsPeakPeriod, dopsitTestWeightLower, dopsitTestWeightUpper, dopsitTestMilesLower, dopsitTestMilesUpper, dopsitTestDomesticLinehaulBasePriceMillicents, dopsitTestContractYearName, dopsitTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price linehaul")
	})
}

func (suite *GHCRateEngineServiceSuite) TestDomesticOriginSITPickupPricer50MilesOrLessDiffZip3s() {
	zipOriginal := "29201"
	zipActual := "29123"       // different zip3
	distance := unit.Miles(23) // <= 50 miles

	pricer := NewDomesticOriginSITPickupPricer()
	expectedPrice := unit.Cents(133689)

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDOPSIT, dopsitTestSchedule, dopsitTestIsPeakPeriod, dopsitTestDomesticOtherBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)

		paymentServiceItem := suite.setupDomesticOriginSITPickupServiceItem(zipOriginal, zipActual, distance)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
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
		suite.setupDomesticOtherPrice(models.ReServiceCodeDOPSIT, dopsitTestSchedule, dopsitTestIsPeakPeriod, dopsitTestDomesticOtherBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.NoError(err)
		suite.Equal(expectedPrice, priceCents)
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDOPSIT, dopsitTestSchedule, dopsitTestIsPeakPeriod, dopsitTestDomesticOtherBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", dopsitTestRequestedPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic origin SIT pickup rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDOPSIT, dopsitTestSchedule, dopsitTestIsPeakPeriod, dopsitTestDomesticOtherBasePriceCents, dopsitTestContractYearName, dopsitTestEscalationCompounded)
		twoYearsLaterPickupDate := dopsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dopsitTestWeight, dopsitTestServiceArea, dopsitTestSchedule, zipOriginal, zipActual, distance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticOriginSITPickupServiceItem(zipOriginal string, zipActual string, distance unit.Miles) models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDOPSIT,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(distance)),
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dopsitTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameSITServiceAreaOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   dopsitTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameSITScheduleOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", dopsitTestSchedule),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dopsitTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameZipSITOriginHHGActualAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   zipActual,
			},
			{
				Key:     models.ServiceItemParamNameZipSITOriginHHGOriginalAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   zipOriginal,
			},
		}, nil, nil,
	)
}
