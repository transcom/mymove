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
	dlhTestServiceArea = "004"
	dlhTestDistance    = unit.Miles(1200)
	dlhTestWeight      = unit.Pound(4000)
	dlhPriceCents      = unit.Cents(249770)
)

var dlhRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticLinehaul() {
	suite.setupDomesticLinehaulData()
	paymentServiceItem := suite.setupDomesticLinehaulServiceItem()
	linehaulServicePricer := NewDomesticLinehaulPricer(suite.DB())

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := linehaulServicePricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(dlhPriceCents, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, true, int(dlhTestDistance), int(dlhTestWeight), dlhTestServiceArea)
		//contractCode, requestedPickupDate, isPeakPeriod, distanceZip3, weightBilledActual, serviceAreaOrigin
		suite.NoError(err)
		suite.Equal(dlhPriceCents, priceCents)
	})

	suite.T().Run("sending PaymentServiceItemParams without expected param", func(t *testing.T) {
		_, err := linehaulServicePricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
	weightBilledActualIndex := 5
	if paramsWithBelowMinimumWeight[weightBilledActualIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilledActual {
		suite.T().Fatalf("Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilledActual, paramsWithBelowMinimumWeight[5].ServiceItemParamKey.Key)
	}
	paramsWithBelowMinimumWeight[weightBilledActualIndex].Value = "200"
	suite.T().Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilledActual", func(t *testing.T) {
		priceCents, err := linehaulServicePricer.PriceUsingParams(paramsWithBelowMinimumWeight)
		suite.Equal("could not fetch domestic linehaul rate: weight must be at least 500", err.Error())
		suite.Equal(unit.Cents(0), priceCents)
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, err := linehaulServicePricer.Price("BOGUS", dlhRequestedPickupDate, true, int(dlhTestDistance), int(dlhTestWeight), dlhTestServiceArea)
		suite.Error(err)
	})

	suite.T().Run("validation errors", func(t *testing.T) {
		// No move date
		_, err := linehaulServicePricer.Price("BOGUS", time.Time{}, true, int(dlhTestDistance), int(dlhTestWeight), dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: MoveDate is required", err.Error())

		// No distance
		_, err = linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, true, 0, int(dlhTestWeight), dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: distance must be at least 50", err.Error())

		// Short haul distance
		_, err = linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, true, 49, int(dlhTestWeight), dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: distance must be at least 50", err.Error())

		// No weight
		_, err = linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, true, int(dlhTestDistance), 0, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: weight must be at least 500", err.Error())

		// No service area
		_, err = linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, true, int(dlhTestDistance), int(dlhTestWeight), "")
		suite.Error(err)
		suite.Equal("could not fetch domestic linehaul rate: ServiceArea is required", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulData() {

	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Escalation:           1.0197,
				EscalationCompounded: 1.04071,
			},
		})

	serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
		testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:    contractYear.Contract,
				ServiceArea: dlhTestServiceArea,
			},
		})

	baseLinehaulPrice := models.ReDomesticLinehaulPrice{
		ContractID:            contractYear.Contract.ID,
		WeightLower:           500,
		WeightUpper:           4999,
		MilesLower:            1001,
		MilesUpper:            1500,
		IsPeakPeriod:          true,
		DomesticServiceAreaID: serviceArea.ID,
	}

	linehaulPricePeak := baseLinehaulPrice
	linehaulPricePeak.PriceMillicents = 5000 // 0.050
	suite.MustSave(&linehaulPricePeak)

}

func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulServiceItem() models.PaymentServiceItem {
	return testdatagen.MakePaymentServiceItemWithParams(
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
				KeyType: models.ServiceItemParamTypeTimestamp,
				Value:   dlhRequestedPickupDate.Format(TimestampParamFormat),
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
				Value:   "1400",
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
