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
	dlhTestServiceArea          = "004"
	dlhTestIsPeakPeriod         = true
	dlhTestWeightLower          = unit.Pound(500)
	dlhTestWeightUpper          = unit.Pound(4999)
	dlhTestMilesLower           = 1001
	dlhTestMilesUpper           = 1500
	dlhTestBasePriceMillicents  = unit.Millicents(5000)
	dlhTestEscalationCompounded = 1.04071
	dlhTestDistance             = unit.Miles(1200)
	dlhTestWeight               = unit.Pound(4000)
	dlhPriceCents               = unit.Cents(249770)
)

var dlhRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticLinehaul() {
	suite.setupDomesticLinehaulPrice(dlhTestServiceArea, dlhTestIsPeakPeriod, dlhTestWeightLower, dlhTestWeightUpper, dlhTestMilesLower, dlhTestMilesUpper, dlhTestBasePriceMillicents, dlhTestEscalationCompounded)
	paymentServiceItem := suite.setupDomesticLinehaulServiceItem()
	linehaulServicePricer := NewDomesticLinehaulPricer(suite.DB())

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, _, err := linehaulServicePricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(dlhPriceCents, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, _, err := linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, dlhTestDistance, dlhTestWeight, dlhTestServiceArea)
		suite.NoError(err)
		suite.Equal(dlhPriceCents, priceCents)
	})

	suite.T().Run("sending PaymentServiceItemParams without expected param", func(t *testing.T) {
		_, _, err := linehaulServicePricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
	weightBilledActualIndex := 5
	if paramsWithBelowMinimumWeight[weightBilledActualIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilledActual {
		suite.T().Fatalf("Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilledActual, paramsWithBelowMinimumWeight[5].ServiceItemParamKey.Key)
	}
	paramsWithBelowMinimumWeight[weightBilledActualIndex].Value = "200"
	suite.T().Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilledActual", func(t *testing.T) {
		priceCents, _, err := linehaulServicePricer.PriceUsingParams(paramsWithBelowMinimumWeight)
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: weight must be at least 500", err.Error())
		suite.Equal(unit.Cents(0), priceCents)
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, _, err := linehaulServicePricer.Price("BOGUS", dlhRequestedPickupDate, dlhTestDistance, dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
	})

	suite.T().Run("validation errors", func(t *testing.T) {
		// No move date
		_, _, err := linehaulServicePricer.Price("BOGUS", time.Time{}, dlhTestDistance, dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: MoveDate is required", err.Error())

		// No distance
		_, _, err = linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, unit.Miles(0), dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: distance must be at least 50", err.Error())

		// Short haul distance
		_, _, err = linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, unit.Miles(49), dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: distance must be at least 50", err.Error())

		// No weight
		_, _, err = linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, dlhTestDistance, unit.Pound(0), dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: weight must be at least 500", err.Error())

		// No service area
		_, _, err = linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, dlhTestDistance, dlhTestWeight, "")
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: ServiceArea is required", err.Error())
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
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dlhRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip3,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dlhTestDistance)),
			},
			{
				Key:     models.ServiceItemParamNameZipPickupAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   "90210",
			},
			{
				Key:     models.ServiceItemParamNameZipDestAddress,
				KeyType: models.ServiceItemParamTypeString,
				Value:   "94535",
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dlhTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameWeightActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(dlhTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameWeightEstimated,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "1400",
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   dlhTestServiceArea,
			},
		},
	)
}
