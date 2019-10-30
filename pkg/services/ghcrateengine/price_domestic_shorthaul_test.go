package ghcrateengine

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const (
	dshTestServiceArea = "006"
	dshTestWeight      = 3600
	dshTestMileage     = 1200
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticShorthaul() {
	suite.setUpDomesticShorthaulData()

	pricer := NewDomesticShorthaulPricer(suite.DB(), suite.logger, testdatagen.DefaultContractCode)

	suite.T().Run("success shorthaul cost within peak period", func(t *testing.T) {
		cost, err := pricer.PriceDomesticShorthaul(
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dshTestMileage,
			dshTestWeight,
			dshTestServiceArea,
		)
		expectedCost := unit.Cents(6563903)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)

	})

	suite.T().Run("success shorthaul cost within non-peak period", func(t *testing.T) {
		nonPeakDate := peakStart.addDate(0, -1)
		cost, err := pricer.PriceDomesticShorthaul(
			time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
			dshTestMileage,
			dshTestWeight,
			dshTestServiceArea,
		)
		expectedCost := unit.Cents(5709696)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.T().Run("Failure if move date is outside of contract year", func(t *testing.T) {
		_, err := pricer.PriceDomesticShorthaul(
			time.Date(testdatagen.TestYear+1, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dshTestMileage,
			dshTestWeight,
			dshTestServiceArea,
		)

		suite.Error(err)
	})

	suite.T().Run("weight below minimum", func(t *testing.T) {
		cost, err := pricer.PriceDomesticShorthaul(
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dshTestMileage,
			unit.Pound(499),
			dshTestServiceArea,
		)
		expectedCost := unit.Cents(911653)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)

	})

	suite.T().Run("validation errors", func(t *testing.T) {
		moveDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// No move date
		_, err := pricer.PriceDomesticShorthaul(time.Time{}, dshTestMileage, dshTestWeight, dshTestServiceArea)
		suite.Error(err)
		suite.Equal("MoveDate is required", err.Error())

		// No distance
		_, err = pricer.PriceDomesticShorthaul(moveDate, 0, dshTestWeight, dshTestServiceArea)
		suite.Error(err)
		suite.Equal("Distance must be greater than 0", err.Error())

		// No weight
		_, err = pricer.PriceDomesticShorthaul(moveDate, dshTestMileage, 0, dshTestServiceArea)
		suite.Error(err)
		suite.Equal("Weight must be greater than 0", err.Error())

		// No service area
		_, err = pricer.PriceDomesticShorthaul(moveDate, dshTestMileage, dshTestWeight, "")
		suite.Error(err)
		suite.Equal("ServiceArea is required", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) setUpDomesticShorthaulData() {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Escalation:           1.0197,
				EscalationCompounded: 1.0407,
			},
		})

	serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
		testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea: dshTestServiceArea,
			},
		})

	domesticShorthaulService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DSH",
				Name: "Dom. Shorthaul",
			},
		})

	domesticShorthaulPrice := models.ReDomesticServiceAreaPrice{
		ContractID:            contractYear.Contract.ID,
		DomesticServiceAreaID: serviceArea.ID,
		IsPeakPeriod:          true,
		ServiceID:             domesticShorthaulService.ID,
	}

	domesticShorthaulPeakPrice := domesticShorthaulPrice
	domesticShorthaulPeakPrice.PriceCents = 146
	suite.MustSave(&domesticShorthaulPeakPrice)

	domesticShorthaulNonpeakPrice := domesticShorthaulPrice
	domesticShorthaulNonpeakPrice.IsPeakPeriod = false
	domesticShorthaulNonpeakPrice.PriceCents = 127
	suite.MustSave(&domesticShorthaulNonpeakPrice)
}