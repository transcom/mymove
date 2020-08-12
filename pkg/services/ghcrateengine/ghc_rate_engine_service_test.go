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
	suite.DB().TruncateAll()
}

func TestGHCRateEngineServiceSuite(t *testing.T) {
	ts := &GHCRateEngineServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

type createParams struct {
	key     models.ServiceItemParamName
	keyType models.ServiceItemParamType
	value   string
}

func (suite *GHCRateEngineServiceSuite) setupPaymentServiceItemWithParams(serviceCode models.ReServiceCode, paramsToCreate []createParams) models.PaymentServiceItem {
	var params models.PaymentServiceItemParams

	paymentServiceItem := testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: serviceCode,
		},
	})

	for _, param := range paramsToCreate {
		serviceItemParamKey := testdatagen.MakeServiceItemParamKey(suite.DB(),
			testdatagen.Assertions{
				ServiceItemParamKey: models.ServiceItemParamKey{
					Key:  param.key,
					Type: param.keyType,
				},
			})

		serviceItemParam := testdatagen.MakePaymentServiceItemParam(suite.DB(),
			testdatagen.Assertions{
				PaymentServiceItem:  paymentServiceItem,
				ServiceItemParamKey: serviceItemParamKey,
				PaymentServiceItemParam: models.PaymentServiceItemParam{
					Value: param.value,
				},
			})
		params = append(params, serviceItemParam)
	}

	paymentServiceItem.PaymentServiceItemParams = params

	return paymentServiceItem
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

func (suite *GHCRateEngineServiceSuite) setUpDomesticPackData(code models.ReServiceCode) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Escalation:           1.0197,
				EscalationCompounded: 1.0407,
			},
		})

	domesticPackService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: code,
			},
		})

	domesticPackPrice := models.ReDomesticOtherPrice{
		ContractID:   contractYear.Contract.ID,
		Schedule:     servicesScheduleOrigin,
		IsPeakPeriod: true,
		ServiceID:    domesticPackService.ID,
	}

	domesticPackPeakPrice := domesticPackPrice
	domesticPackPeakPrice.PriceCents = 146
	suite.MustSave(&domesticPackPeakPrice)

	domesticPackNonpeakPrice := domesticPackPrice
	domesticPackNonpeakPrice.IsPeakPeriod = false
	domesticPackNonpeakPrice.PriceCents = 127
	suite.MustSave(&domesticPackNonpeakPrice)
}
