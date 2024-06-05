package transportationaccountingcode

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type TransportationAccountingCodeServiceSuite struct {
	*testingsuite.PopTestSuite
	tacFetcher services.TransportationAccountingCodeFetcher
}

func TestTransportationAccountingCodeServiceSuite(t *testing.T) {
	ts := &TransportationAccountingCodeServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction(),
		),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *TransportationAccountingCodeServiceSuite) SetupTest() {
	suite.tacFetcher = NewTransportationAccountingCodeFetcher()
}

func (suite *TransportationAccountingCodeServiceSuite) TestFetchOrderTransportationAccountingCodes() {
	suite.Run("successfully fetches TACs by affiliation", func() {
		ordersIssueDate := time.Now()
		endDate := ordersIssueDate.AddDate(1, 0, 0)
		tacCode := "CACI"
		loa := factory.BuildLineOfAccounting(suite.AppContextForTest().DB(), []factory.Customization{
			{
				Model: models.LineOfAccounting{
					LoaBgnDt:   &ordersIssueDate,
					LoaEndDt:   &endDate,
					LoaSysID:   models.StringPointer("1234567890"),
					LoaHsGdsCd: models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
				},
			},
		}, nil)
		factory.BuildTransportationAccountingCode(suite.AppContextForTest().DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               tacCode,
					TrnsprtnAcntBgnDt: &ordersIssueDate,
					TrnsprtnAcntEndDt: &endDate,
					TacFnBlModCd:      models.StringPointer("1"),
					LoaSysID:          loa.LoaSysID,
				},
			},
		}, nil)

		testCases := []struct {
			affiliation models.ServiceMemberAffiliation
			shouldError bool
		}{
			{models.AffiliationCOASTGUARD, false},
			{models.AffiliationARMY, false},
		}
		for _, testCase := range testCases {
			tacs, err := suite.tacFetcher.FetchOrderTransportationAccountingCodes(testCase.affiliation, ordersIssueDate, tacCode, suite.AppContextForTest())
			if testCase.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
			suite.NotEmpty(tacs)
			suite.Equal(tacCode, tacs[0].TAC)
			// Assert that the TAC came back with a LOA. This is important as
			// the line of accounting service object will need these.
			// The LOA service object uses this service as the line of accounting
			// is attached to a transportation accounting code.
			suite.NotNil(tacs[0].LineOfAccounting)
		}
	})
}
