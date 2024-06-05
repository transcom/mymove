package lineofaccounting

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	transportationaccountingcode "github.com/transcom/mymove/pkg/services/transportation_accounting_code"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type LineOfAccountingServiceSuite struct {
	*testingsuite.PopTestSuite
	tacFetcher services.TransportationAccountingCodeFetcher
	loaFetcher services.LinesOfAccountingFetcher
}

func TestLineOfAccountingServiceSuite(t *testing.T) {
	ts := &LineOfAccountingServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction(),
		),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *LineOfAccountingServiceSuite) SetupTest() {
	suite.tacFetcher = transportationaccountingcode.NewTransportationAccountingCodeFetcher()
	suite.loaFetcher = NewLinesOfAccountingFetcher(suite.tacFetcher)
}

func (suite *LineOfAccountingServiceSuite) TestFetchOrderLineOfAccountings() {
	suite.Run("successfully fetches LOAs", func() {
		appCtx := suite.AppContextForTest()
		ordersIssueDate := time.Now()
		startDate := ordersIssueDate.AddDate(-1, 0, 0)
		endDate := ordersIssueDate.AddDate(1, 0, 0)
		tacCode := "CACI"
		loa := factory.BuildLineOfAccounting(appCtx.DB(), []factory.Customization{
			{
				Model: models.LineOfAccounting{
					LoaBgnDt:   &startDate,
					LoaEndDt:   &endDate,
					LoaSysID:   models.StringPointer("1234567890"),
					LoaHsGdsCd: models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
				},
			},
		}, nil)
		factory.BuildTransportationAccountingCode(appCtx.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               tacCode,
					TrnsprtnAcntBgnDt: &startDate,
					TrnsprtnAcntEndDt: &endDate,
					TacFnBlModCd:      models.StringPointer("1"),
					LoaSysID:          loa.LoaSysID,
				},
			},
		}, nil)

		// Get the TACs that we will extract LOAs from
		tacs, err := suite.tacFetcher.FetchOrderTransportationAccountingCodes(models.AffiliationARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		suite.NotEmpty(tacs)
		// Ensure LOA isn't nil
		suite.NotNil(tacs[0].LineOfAccounting)
		// Extract LOAs
		loas, err := suite.loaFetcher.FetchLongLinesOfAccounting(models.AffiliationARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		suite.Equal(loa.ID, loas[0].ID)
	})
}
