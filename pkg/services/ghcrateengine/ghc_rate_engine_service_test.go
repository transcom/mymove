package ghcrateengine

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

type GHCRateEngineServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *GHCRateEngineServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestGHCRateEngineServiceSuite(t *testing.T) {
	ts := &GHCRateEngineServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *GHCRateEngineServiceSuite) setupTaskOrderFeeData(code models.ReServiceCode, priceCents unit.Cents) {
	contractYear := testdatagen.MakeDefaultReContractYear(suite.DB())

	counselingService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: code,
			},
		})

	taskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      counselingService.ID,
		PriceCents:     priceCents,
	}
	suite.MustSave(&taskOrderFee)
}

func (suite *GHCRateEngineServiceSuite) setUpDomesticPackAndUnpackData(code models.ReServiceCode) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Escalation:           1.0197,
				EscalationCompounded: 1.0407,
				Name:                 "Base Period Year 1",
			},
		})

	domesticPackUnpackService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: code,
			},
		})

	domesticPackUnpackPrice := models.ReDomesticOtherPrice{
		ContractID:   contractYear.Contract.ID,
		Schedule:     1,
		IsPeakPeriod: true,
		ServiceID:    domesticPackUnpackService.ID,
	}

	domesticPackUnpackPeakPrice := domesticPackUnpackPrice
	domesticPackUnpackPeakPrice.PriceCents = 146
	suite.MustSave(&domesticPackUnpackPeakPrice)

	domesticPackUnpackNonpeakPrice := domesticPackUnpackPrice
	domesticPackUnpackNonpeakPrice.IsPeakPeriod = false
	domesticPackUnpackNonpeakPrice.PriceCents = 127
	suite.MustSave(&domesticPackUnpackNonpeakPrice)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticOtherPrice(code models.ReServiceCode, schedule int, isPeakPeriod bool, priceCents unit.Cents, escalationCompounded float64) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				EscalationCompounded: escalationCompounded,
			},
		})

	service := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: code,
			},
		})

	otherPrice := models.ReDomesticOtherPrice{
		ContractID:   contractYear.Contract.ID,
		ServiceID:    service.ID,
		IsPeakPeriod: isPeakPeriod,
		Schedule:     schedule,
		PriceCents:   priceCents,
	}

	suite.MustSave(&otherPrice)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticServiceAreaPrice(code models.ReServiceCode, serviceAreaCode string, isPeakPeriod bool, priceCents unit.Cents, escalationCompounded float64) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				EscalationCompounded: escalationCompounded,
			},
		})

	service := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: code,
			},
		})

	serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
		testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:    contractYear.Contract,
				ServiceArea: serviceAreaCode,
			},
		})

	serviceAreaPrice := models.ReDomesticServiceAreaPrice{
		ContractID:            contractYear.Contract.ID,
		ServiceID:             service.ID,
		IsPeakPeriod:          isPeakPeriod,
		DomesticServiceAreaID: serviceArea.ID,
		PriceCents:            priceCents,
	}

	suite.MustSave(&serviceAreaPrice)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulPrice(serviceAreaCode string, isPeakPeriod bool, weightLower unit.Pound, weightUpper unit.Pound, milesLower int, milesUpper int, priceMillicents unit.Millicents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
			},
		})

	serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
		testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:    contractYear.Contract,
				ServiceArea: serviceAreaCode,
			},
		})

	baseLinehaulPrice := models.ReDomesticLinehaulPrice{
		ContractID:            contractYear.Contract.ID,
		WeightLower:           weightLower,
		WeightUpper:           weightUpper,
		MilesLower:            milesLower,
		MilesUpper:            milesUpper,
		IsPeakPeriod:          isPeakPeriod,
		DomesticServiceAreaID: serviceArea.ID,
		PriceMillicents:       priceMillicents,
	}

	suite.MustSave(&baseLinehaulPrice)
}

func (suite *GHCRateEngineServiceSuite) HasDisplayParam(displayParams services.PricingDisplayParams, key models.ServiceItemParamName, value string) bool {
	for _, displayParam := range displayParams {
		if displayParam.Key == key {
			return suite.Equal(value, displayParam.Value, "%s param actual value did not match expected", key.String())
		}
	}

	return suite.Failf("Could not find display param", "key=<%s> value=<%s>", key.String(), value)
}
