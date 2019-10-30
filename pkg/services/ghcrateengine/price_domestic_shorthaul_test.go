package ghcrateengine

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const (
	dshTestServiceArea = "005"
	dshWeight          = 3500
	dshMileage         = 1200
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticShorthaul() {
	suite.setUpDomesticShorthaulData()

	pricer := NewDomesticShorthaulPricer(suite.DB(), suite.logger, testdatagen.DefaultContractCode)

	suite.T().Run("success shorthaul cost within peak period", func(t *testing.T) {
		cost, err := pricer.PriceDomesticShorthaul(
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dshMileage,
			dshWeight,
			dshTestServiceArea,
		)
		expectedCost := 23569
		suite.NoError(err)
		suite.Equal(expectedCost, cost)

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
	suite.MustSave(&domesticShorthaulPrice)

	domesticShorthaulNonpeakPrice := domesticShorthaulPrice
	domesticShorthaulNonpeakPrice.IsPeakPeriod = false
	domesticShorthaulNonpeakPrice.PriceCents = 127
	suite.MustSave(&domesticShorthaulPrice)
}