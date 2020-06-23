package ghcrateengine

import (
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
)


var dlhRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticLinehaul() {
	suite.setupDomesticLinehaulData()
	paymentServiceItem := suite.setupDomesticLinehaulServiceItem()
	linehaulServicePricer := NewDomesticLinehaulPricer(suite.DB())

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := linehaulServicePricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(csPriceCents, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := linehaulServicePricer.Price(testdatagen.DefaultContractCode, dlhRequestedPickupDate, true, int(dlhTestDistance), int(dlhTestWeight), dlhTestServiceArea)
		//contractCode, requestedPickupDate, isPeakPeriod, distanceZip3, weightBilledActual, serviceAreaOrigin
		suite.NoError(err)
		suite.Equal(csPriceCents, priceCents)
	})

	suite.T().Run("sending PaymentServiceItemParams without expected param", func(t *testing.T) {
		_, err := linehaulServicePricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, err := linehaulServicePricer.Price("BOGUS",  dlhRequestedPickupDate, true, int(dlhTestDistance), int(dlhTestWeight), dlhTestServiceArea)
		suite.Error(err)
	})
}

/*
func (suite *GHCRateEngineServiceSuite) TestPriceDomesticLinehaulDELETE() {
	suite.setupDomesticLinehaulData()
	domesticLinehaulPricer := NewDomesticLinehaulPricer(suite.DB())

	suite.T().Run("success within peak period", func(t *testing.T) {
		totalPrice, err := domesticLinehaulPricer.PriceDomesticLinehaul(
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dlhTestDistance, dlhTestWeight, dlhTestServiceArea)

		suite.NoError(err)
		suite.Equal(249770, totalPrice.Int())
	})

	suite.T().Run("success within non-peak period", func(t *testing.T) {
		nonPeakDate := peakStart.addDate(0, -1)
		totalPrice, err := domesticLinehaulPricer.PriceDomesticLinehaul(
			time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
			dlhTestDistance, dlhTestWeight, dlhTestServiceArea)

		suite.NoError(err)
		suite.Equal(224793, totalPrice.Int())
	})

	suite.T().Run("weight below minimum", func(t *testing.T) {
		totalPrice, err := domesticLinehaulPricer.PriceDomesticLinehaul(
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dlhTestDistance, minDomesticWeight/2, dlhTestServiceArea)

		suite.NoError(err)
		suite.Equal(31221, totalPrice.Int())
	})

	suite.T().Run("date outside of valid contract year", func(t *testing.T) {
		_, err := domesticLinehaulPricer.PriceDomesticLinehaul(
			time.Date(testdatagen.TestYear-1, time.January, 1, 0, 0, 0, 0, time.UTC),
			dlhTestDistance, dlhTestWeight, dlhTestServiceArea)

		suite.Error(err)
	})

	suite.T().Run("validation errors", func(t *testing.T) {
		moveDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// No move date
		_, err := domesticLinehaulPricer.PriceDomesticLinehaul(time.Time{}, dlhTestDistance, dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("MoveDate is required", err.Error())

		// No distance
		_, err = domesticLinehaulPricer.PriceDomesticLinehaul(moveDate, 0, dlhTestWeight, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("Distance must be greater than 0", err.Error())

		// No weight
		_, err = domesticLinehaulPricer.PriceDomesticLinehaul(moveDate, dlhTestDistance, 0, dlhTestServiceArea)
		suite.Error(err)
		suite.Equal("Weight must be greater than 0", err.Error())

		// No service area
		_, err = domesticLinehaulPricer.PriceDomesticLinehaul(moveDate, dlhTestDistance, dlhTestWeight, "")
		suite.Error(err)
		suite.Equal("ServiceArea is required", err.Error())
	})
}

 */

func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulData() {

	contractYear := testdatagen.MakeDefaultReContractYear(suite.DB())

	linehaulServiceItem := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDLH,
			},
		})

	taskOrderFee := models.ReDomesticLinehaulPrice{

	}
	suite.MustSave(&taskOrderFee)

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

	linehaulPriceNonPeak := baseLinehaulPrice
	linehaulPriceNonPeak.IsPeakPeriod = false
	linehaulPriceNonPeak.PriceMillicents = 4500 // 0.045
	suite.MustSave(&linehaulPriceNonPeak)
}


func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulServiceItem() models.PaymentServiceItem {
	return suite.setupPaymentServiceItemWithParams(
		models.ReServiceCodeDLH,
		[]createParams{
			{
				models.ServiceItemParamNameContractCode,
				models.ServiceItemParamTypeString,
				testdatagen.DefaultContractCode,
			},
			{
				models.ServiceItemParamNameRequestedPickupDate,
				models.ServiceItemParamTypeTimestamp,
				dlhRequestedPickupDate.Format(TimestampParamFormat),
			},
			{
				models.ServiceItemParamNameDistanceZip3,
				models.ServiceItemParamTypeInteger,
				"1700",
			},
			{
				models.ServiceItemParamNameZipPickupAddress,
				models.ServiceItemParamTypeString,
				"90210",
			},
			{
				models.ServiceItemParamNameZipDestAddress,
				models.ServiceItemParamTypeString,
				"94535",
			},
			{
				models.ServiceItemParamNameWeightBilledActual,
				models.ServiceItemParamTypeInteger,
				"1400",
			},
			{
				models.ServiceItemParamNameWeightActual,
				models.ServiceItemParamTypeInteger,
				"1400",
			},
			{
				models.ServiceItemParamNameWeightEstimated,
				models.ServiceItemParamTypeInteger,
				"1400",
			},
		},
	)
}