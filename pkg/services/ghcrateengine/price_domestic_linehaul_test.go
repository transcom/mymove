package ghcrateengine

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	testContractCode = "TEST"
	testContractYear = 2019
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticLinehaul() {
	suite.setupDomesticLinehaulData()
	domesticLinehaulPricer := NewDomesticLinehaulPricer(suite.DB(), suite.logger, testContractCode)

	suite.T().Run("success within peak period", func(t *testing.T) {
		pricingData := services.DomesticServicePricingData{
			MoveDate:    time.Date(testContractYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
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
			MoveDate:    time.Date(testContractYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
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
			MoveDate:    time.Date(testContractYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			Distance:    unit.Miles(1200),
			Weight:      minWeight / 2,
			ServiceArea: "004",
		}
		totalPrice, err := domesticLinehaulPricer.PriceDomesticLinehaul(pricingData)

		suite.NoError(err)
		suite.Equal(31221, totalPrice.Int())
	})

	suite.T().Run("date outside of valid contract year", func(t *testing.T) {
		pricingData := services.DomesticServicePricingData{
			MoveDate:    time.Date(testContractYear-1, time.January, 1, 0, 0, 0, 0, time.UTC),
			Distance:    unit.Miles(1200),
			Weight:      unit.Pound(4000),
			ServiceArea: "004",
		}
		_, err := domesticLinehaulPricer.PriceDomesticLinehaul(pricingData)

		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulData() {
	contract := models.ReContract{
		Code: testContractCode,
		Name: "Test Contract",
	}
	suite.MustSave(&contract)

	contractYear := models.ReContractYear{
		ContractID:           contract.ID,
		Name:                 "Base Period Year 3",
		StartDate:            time.Date(testContractYear, time.January, 1, 0, 0, 0, 0, time.UTC),
		EndDate:              time.Date(testContractYear, time.December, 31, 0, 0, 0, 0, time.UTC),
		Escalation:           1.0197,
		EscalationCompounded: 1.04071,
	}
	suite.MustSave(&contractYear)

	serviceArea := models.ReDomesticServiceArea{
		BasePointCity:    "Birmingham",
		State:            "AL",
		ServiceArea:      "004",
		ServicesSchedule: 2,
		SITPDSchedule:    2,
	}
	suite.MustSave(&serviceArea)

	linehaulPricePeak := models.ReDomesticLinehaulPrice{
		ContractID:            contract.ID,
		WeightLower:           500,
		WeightUpper:           4999,
		MilesLower:            1001,
		MilesUpper:            1500,
		IsPeakPeriod:          true,
		DomesticServiceAreaID: serviceArea.ID,
		PriceMillicents:       5000, // 0.050
	}
	suite.MustSave(&linehaulPricePeak)

	linehaulPriceNonPeak := models.ReDomesticLinehaulPrice{
		ContractID:            contract.ID,
		WeightLower:           500,
		WeightUpper:           4999,
		MilesLower:            1001,
		MilesUpper:            1500,
		IsPeakPeriod:          false,
		DomesticServiceAreaID: serviceArea.ID,
		PriceMillicents:       4500, // 0.045
	}
	suite.MustSave(&linehaulPriceNonPeak)
}
