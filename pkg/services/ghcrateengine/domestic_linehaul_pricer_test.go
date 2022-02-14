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
	dlhTestServiceArea          = "004"
	dlhTestIsPeakPeriod         = true
	dlhTestWeightLower          = unit.Pound(500)
	dlhTestWeightUpper          = unit.Pound(4999)
	dlhTestMilesLower           = 1001
	dlhTestMilesUpper           = 1500
	dlhTestBasePriceMillicents  = unit.Millicents(5100)
	dlhTestContractYearName     = "DLH Test Year"
	dlhTestEscalationCompounded = 1.04071
	dlhTestDistance             = unit.Miles(1200)
	dlhTestWeight               = unit.Pound(4000)
	dlhPriceCents               = unit.Cents(254766)
)

var dlhRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticLinehaul() {
	linehaulServicePricer := NewDomesticLinehaulPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticLinehaulPrice(dlhTestServiceArea, dlhTestIsPeakPeriod, dlhTestWeightLower, dlhTestWeightUpper, dlhTestMilesLower, dlhTestMilesUpper, dlhTestBasePriceMillicents, dlhTestContractYearName, dlhTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticLinehaulServiceItem()
		priceCents, displayParams, err := linehaulServicePricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(dlhPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dlhTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dlhTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dlhTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatFloat(dlhTestBasePriceMillicents.ToDollarFloatNoRound(), 3)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupDomesticLinehaulPrice(dlhTestServiceArea, dlhTestIsPeakPeriod, dlhTestWeightLower, dlhTestWeightUpper, dlhTestMilesLower, dlhTestMilesUpper, dlhTestBasePriceMillicents, dlhTestContractYearName, dlhTestEscalationCompounded)
		//paymentServiceItem := suite.setupDomesticLinehaulServiceItem()
		priceCents, _, err := linehaulServicePricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dlhRequestedPickupDate, dlhTestDistance, dlhTestWeight, dlhTestServiceArea)
		suite.NoError(err)
		suite.Equal(dlhPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		_, _, err := linehaulServicePricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilled", func() {
		suite.setupDomesticLinehaulPrice(dlhTestServiceArea, dlhTestIsPeakPeriod, dlhTestWeightLower, dlhTestWeightUpper, dlhTestMilesLower, dlhTestMilesUpper, dlhTestBasePriceMillicents, dlhTestContractYearName, dlhTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticLinehaulServiceItem()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 4
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[5].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"

		priceCents, _, err := linehaulServicePricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBelowMinimumWeight)
		suite.Error(err)
		suite.Equal("Weight must be at least 500", err.Error())
		suite.Equal(unit.Cents(0), priceCents)
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticLinehaulPrice(dlhTestServiceArea, dlhTestIsPeakPeriod, dlhTestWeightLower, dlhTestWeightUpper, dlhTestMilesLower, dlhTestMilesUpper, dlhTestBasePriceMillicents, dlhTestContractYearName, dlhTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticLinehaulServiceItem()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 4
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[5].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"

		_, _, err := linehaulServicePricer.Price(suite.AppContextForTest(), "BOGUS", dlhRequestedPickupDate, dlhTestDistance, dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
	})

	suite.Run("validation errors", func() {
		suite.setupDomesticLinehaulPrice(dlhTestServiceArea, dlhTestIsPeakPeriod, dlhTestWeightLower, dlhTestWeightUpper, dlhTestMilesLower, dlhTestMilesUpper, dlhTestBasePriceMillicents, dlhTestContractYearName, dlhTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticLinehaulServiceItem()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 4
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[5].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"

		// No contract code
		_, _, err := linehaulServicePricer.Price(suite.AppContextForTest(), "", dlhRequestedPickupDate, dlhTestDistance, dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())

		// No reference date
		_, _, err = linehaulServicePricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, time.Time{}, dlhTestDistance, dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("ReferenceDate is required", err.Error())

		// No distance
		_, _, err = linehaulServicePricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dlhRequestedPickupDate, unit.Miles(0), dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("Distance must be at least 50", err.Error())

		// Short haul distance
		_, _, err = linehaulServicePricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dlhRequestedPickupDate, unit.Miles(49), dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("Distance must be at least 50", err.Error())

		// No weight
		_, _, err = linehaulServicePricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dlhRequestedPickupDate, dlhTestDistance, unit.Pound(0), dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("Weight must be at least 500", err.Error())

		// No service area
		_, _, err = linehaulServicePricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dlhRequestedPickupDate, dlhTestDistance, dlhTestWeight, "")
		suite.Error(err)
		suite.Equal("ServiceArea is required", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDLH,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip3,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dlhTestDistance)),
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dlhRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   dlhTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dlhTestWeight)),
			},
		},
	)
}
