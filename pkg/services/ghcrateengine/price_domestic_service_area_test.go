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
	dsaTestServiceArea = "005"
	dsaTestWeight      = unit.Pound(3500)
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticServiceArea() {
	suite.setUpDomesticServiceAreaPricesData()
	pricer := NewDomesticServiceAreaPricer(suite.DB(), suite.logger, testdatagen.DefaultContractCode)

	type testDataStruct struct {
		serviceCode         string
		serviceName         string
		expectedPeakCost    int
		expectedNonpeakCost int
		expectedMinPeakCost int
	}

	testCases := []testDataStruct{
		{
			serviceCode:         "DODP",
			serviceName:         "Dom. O/D Price",
			expectedPeakCost:    28848,
			expectedNonpeakCost: 25096,
			expectedMinPeakCost: 4121,
		},
		{
			serviceCode:         "DFSIT",
			serviceName:         "Dom. O/D 1st Day SIT",
			expectedPeakCost:    80898,
			expectedNonpeakCost: 70335,
			expectedMinPeakCost: 11556,
		},
		{
			serviceCode:         "DASIT",
			serviceName:         "Dom. O/D Add'l SIT",
			expectedPeakCost:    2841,
			expectedNonpeakCost: 2476,
			expectedMinPeakCost: 405,
		},
	}

	for _, c := range testCases {
		suite.T().Run(fmt.Sprintf("success %s cost within peak period", c.serviceName), func(t *testing.T) {

			cost, err := pricer.PriceDomesticServiceArea(
				time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
				dsaTestWeight,
				dsaTestServiceArea,
				c.serviceCode)

			suite.NoError(err)
			suite.Equal(c.expectedPeakCost, cost.Int())
		})

		suite.T().Run(fmt.Sprintf("success %s cost within non-peak period", c.serviceName), func(t *testing.T) {
			nonPeakDate := peakStart.addDate(0, -1)
			cost, err := pricer.PriceDomesticServiceArea(
				time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
				dsaTestWeight,
				dsaTestServiceArea,
				c.serviceCode)

			suite.NoError(err)
			suite.Equal(c.expectedNonpeakCost, cost.Int())
		})

		suite.T().Run(fmt.Sprintf("%s cost weight below minimum", c.serviceName), func(t *testing.T) {
			cost, err := pricer.PriceDomesticServiceArea(
				time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
				450,
				dsaTestServiceArea,
				c.serviceCode)

			suite.NoError(err)
			suite.Equal(c.expectedMinPeakCost, cost.Int())
		})

		//suite.T().Run(fmt.Sprintf("%s date outside of valid contract year", serviceName), func(t *testing.T) {
		//
		//})
	}

	//suite.T().Run("validation errors", func(t *testing.T) {
	//
	//}
}

func (suite *GHCRateEngineServiceSuite) setUpDomesticServiceAreaPricesData() {
	// create contractYear, domesticServiceArea, services data
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
				ServiceArea: dsaTestServiceArea,
			},
		})

	originDestinationService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DODP",
				Name: "Dom. O/D Price",
			},
		})

	sit1Service := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DFSIT",
				Name: "Dom. O/D 1st Day SIT",
			},
		})

	addlSITService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DASIT",
				Name: "Dom. O/D Add'l SIT",
			},
		})

	baseDomesticServiceAreaPrice := models.ReDomesticServiceAreaPrice{
		ContractID:            contractYear.Contract.ID,
		DomesticServiceAreaID: serviceArea.ID,
		IsPeakPeriod:          true,
	}

	// Origin/Destination Price
	oDPrice := baseDomesticServiceAreaPrice
	oDPrice.ServiceID = originDestinationService.ID

	oDPeakPrice := oDPrice
	oDPeakPrice.PriceCents = 792
	suite.MustSave(&oDPeakPrice)

	oDNonpeakPrice := oDPrice
	oDNonpeakPrice.IsPeakPeriod = false
	oDNonpeakPrice.PriceCents = 689
	suite.MustSave(&oDNonpeakPrice)

	// SIT Day 1
	sit1Price := baseDomesticServiceAreaPrice
	sit1Price.ServiceID = sit1Service.ID

	sit1PeakPrice := sit1Price
	sit1PeakPrice.PriceCents = 2221
	suite.MustSave(&sit1PeakPrice)

	sit1NonpeakPrice := sit1Price
	sit1NonpeakPrice.IsPeakPeriod = false
	sit1NonpeakPrice.PriceCents = 1931
	suite.MustSave(&sit1NonpeakPrice)

	// SIT Additional Days
	addlSITPrice := baseDomesticServiceAreaPrice
	addlSITPrice.ServiceID = addlSITService.ID

	addlSITPeakPrice := addlSITPrice
	addlSITPeakPrice.PriceCents = 78
	suite.MustSave(&addlSITPeakPrice)

	addlSITNonpeakPrice := addlSITPrice
	addlSITNonpeakPrice.IsPeakPeriod = false
	addlSITNonpeakPrice.PriceCents = 68
	suite.MustSave(&addlSITNonpeakPrice)
}
