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
	setupTest := func() (time.Time, time.Time, time.Time, string) {
		ordersIssueDate := time.Now()
		startDate := ordersIssueDate.AddDate(-1, 0, 0)
		endDate := ordersIssueDate.AddDate(1, 0, 0)
		tacCode := "CACI"
		return ordersIssueDate, startDate, endDate, tacCode
	}
	suite.Run("successfully fetches LOAs", func() {
		appCtx := suite.AppContextForTest()
		ordersIssueDate, startDate, endDate, tacCode := setupTest()
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
		factory.BuildTransportationAccountingCodeWithoutAttachedLoa(appCtx.DB(), []factory.Customization{
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
	suite.Run("Successfully sorts according to FMBC", func() {
		// Create 4 standalone TACs each with a different fbmc
		fbmcs := []string{
			"M",
			"1",
			"5",
			"3",
		}

		appCtx := suite.AppContextForTest()
		suite.NoError(appCtx.DB().TruncateAll())
		ordersIssueDate, startDate, endDate, tacCode := setupTest()

		// Setup TGET data for tests
		for fbmcMemoryIteration := range fbmcs {
			// Use a unique LoaSysID for each fbmc to avoid duplicates pulling from SQL
			loaSysId := factory.MakeRandomString(20)
			factory.BuildLineOfAccounting(appCtx.DB(), []factory.Customization{
				{
					Model: models.LineOfAccounting{
						LoaBgnDt:   &startDate,
						LoaEndDt:   &endDate,
						LoaSysID:   &loaSysId,
						LoaHsGdsCd: models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
						LoaDscTx:   &fbmcs[fbmcMemoryIteration],
					},
				},
			}, nil)
			factory.BuildTransportationAccountingCodeWithoutAttachedLoa(appCtx.DB(), []factory.Customization{
				{
					Model: models.TransportationAccountingCode{
						TAC:               tacCode,
						TrnsprtnAcntBgnDt: &startDate,
						TrnsprtnAcntEndDt: &endDate,
						TacFnBlModCd:      &fbmcs[fbmcMemoryIteration],
						LoaSysID:          &loaSysId,
					},
				},
			}, nil)
		}

		// Get the TACs that we will extract LOAs from
		tacs, err := suite.tacFetcher.FetchOrderTransportationAccountingCodes(models.AffiliationARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		suite.NotEmpty(tacs)
		// Ensure we got 4 tacs
		suite.Len(tacs, 4)
		// Ensure LOA isn't nil
		suite.NotNil(tacs[0].LineOfAccounting)
		for _, tac := range tacs {
			suite.NotNil(tac.LineOfAccounting)
		}

		// Now that we have our 4 tacs and loas, let's fetch and see if they're in order
		// Extract LOAs
		loas, err := suite.loaFetcher.FetchLongLinesOfAccounting(models.AffiliationARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)

		// Check the order of LOAs
		expectedOrder := []string{"1", "3", "5", "M"}
		// Ensure we got the correct number of LOAs
		suite.Len(loas, len(expectedOrder))
		for i, loa := range loas {
			suite.Equal(expectedOrder[i], *loa.LoaDscTx)
		}
	})
	suite.Run("Filters out FMBC of value P", func() {
		fbmcs := []string{
			"1",
			"P", // Make note of P. Our queries filter this out
		}
		ordersIssueDate, startDate, endDate, tacCode := setupTest()
		appCtx := suite.AppContextForTest()

		// Setup TGET data for tests
		for fbmcMemoryIteration := range fbmcs {
			// Use a unique LoaSysID for each fbmc to avoid duplicates
			loaSysId := factory.MakeRandomString(20)
			factory.BuildLineOfAccounting(appCtx.DB(), []factory.Customization{
				{
					Model: models.LineOfAccounting{
						LoaBgnDt:   &startDate,
						LoaEndDt:   &endDate,
						LoaSysID:   &loaSysId,
						LoaHsGdsCd: models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
						LoaDscTx:   &fbmcs[fbmcMemoryIteration],
					},
				},
			}, nil)
			factory.BuildTransportationAccountingCodeWithoutAttachedLoa(appCtx.DB(), []factory.Customization{
				{
					Model: models.TransportationAccountingCode{
						TAC:               tacCode,
						TrnsprtnAcntBgnDt: &startDate,
						TrnsprtnAcntEndDt: &endDate,
						TacFnBlModCd:      &fbmcs[fbmcMemoryIteration],
						LoaSysID:          &loaSysId,
					},
				},
			}, nil)
		}

		// Get the TACs that we will extract LOAs from
		tacs, err := suite.tacFetcher.FetchOrderTransportationAccountingCodes(models.AffiliationARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		suite.NotEmpty(tacs)
		// Ensure we got 1 tac, FMBC P should've been filtered out
		suite.Len(tacs, 1)
		// Ensure LOA isn't nil
		suite.NotNil(tacs[0].LineOfAccounting)
		for _, tac := range tacs {
			suite.NotNil(tac.LineOfAccounting)
		}
		// Extract LOAs
		loas, err := suite.loaFetcher.FetchLongLinesOfAccounting(models.AffiliationARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		// Ensure we got 1 loa, FBMC P should've been filtered out
		suite.Len(loas, 1)
	})
}
