package ghcrateengine

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
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
