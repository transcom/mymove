package ghcrateengine

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticLinehaul() {
	suite.setupDomesticLinehaulData()
	domesticLinehaulPricer := NewDomesticLinehaulPricer(suite.DB(), suite.logger, testdatagen.DefaultContractCode)

	suite.T().Run("success within peak period", func(t *testing.T) {
		pricingData := services.DomesticServicePricingData{
			MoveDate:    time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			Distance:    unit.Miles(1200),
			Weight:      unit.Pound(4000),
			ServiceArea: "004",
		}
		totalPrice, err := domesticLinehaulPricer.PriceDomesticLinehaul(pricingData)

		suite.NoError(err)
		suite.Equal(249770, totalPrice.Int())
	})

	suite.T().Run("success within non-peak period", func(t *testing.T) {
		nonPeakDate := peakStart.addDate(0, -1)
		pricingData := services.DomesticServicePricingData{
			MoveDate:    time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
			Distance:    unit.Miles(1200),
			Weight:      unit.Pound(4000),
			ServiceArea: "004",
		}
		totalPrice, err := domesticLinehaulPricer.PriceDomesticLinehaul(pricingData)

		suite.NoError(err)
		suite.Equal(224793, totalPrice.Int())
	})

	suite.T().Run("weight below minimum", func(t *testing.T) {
		pricingData := services.DomesticServicePricingData{
			MoveDate:    time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			Distance:    unit.Miles(1200),
			Weight:      minDomesticWeight / 2,
			ServiceArea: "004",
		}
		totalPrice, err := domesticLinehaulPricer.PriceDomesticLinehaul(pricingData)

		suite.NoError(err)
		suite.Equal(31221, totalPrice.Int())
	})

	suite.T().Run("date outside of valid contract year", func(t *testing.T) {
		pricingData := services.DomesticServicePricingData{
			MoveDate:    time.Date(testdatagen.TestYear-1, time.January, 1, 0, 0, 0, 0, time.UTC),
			Distance:    unit.Miles(1200),
			Weight:      unit.Pound(4000),
			ServiceArea: "004",
		}
		_, err := domesticLinehaulPricer.PriceDomesticLinehaul(pricingData)

		suite.Error(err)
	})

	suite.T().Run("validation errors", func(t *testing.T) {
		basePricingData := services.DomesticServicePricingData{
			MoveDate:    time.Date(testdatagen.TestYear-1, time.January, 1, 0, 0, 0, 0, time.UTC),
			Distance:    unit.Miles(1200),
			Weight:      unit.Pound(4000),
			ServiceArea: "004",
		}

		noMoveDate := basePricingData
		noMoveDate.MoveDate = time.Time{}
		_, err := domesticLinehaulPricer.PriceDomesticLinehaul(noMoveDate)
		suite.Error(err)
		suite.Equal("MoveDate is required", err.Error())

		noDistance := basePricingData
		noDistance.Distance = 0
		_, err = domesticLinehaulPricer.PriceDomesticLinehaul(noDistance)
		suite.Error(err)
		suite.Equal("Distance must be greater than 0", err.Error())

		noWeight := basePricingData
		noWeight.Weight = 0
		_, err = domesticLinehaulPricer.PriceDomesticLinehaul(noWeight)
		suite.Error(err)
		suite.Equal("Weight must be greater than 0", err.Error())

		noServiceArea := basePricingData
		noServiceArea.ServiceArea = ""
		_, err = domesticLinehaulPricer.PriceDomesticLinehaul(noServiceArea)
		suite.Error(err)
		suite.Equal("ServiceArea is required", err.Error())
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
				ServiceArea: "004",
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
