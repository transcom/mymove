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
	loaFetcher services.LineOfAccountingFetcher
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
		tacs, err := suite.tacFetcher.FetchOrderTransportationAccountingCodes(models.DepartmentIndicatorARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		suite.NotEmpty(tacs)
		// Ensure LOA isn't nil
		suite.NotNil(tacs[0].LineOfAccounting)
		// Extract LOAs
		loas, err := suite.loaFetcher.FetchLongLinesOfAccounting(models.DepartmentIndicatorARMY, ordersIssueDate, tacCode, appCtx)
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
		tacs, err := suite.tacFetcher.FetchOrderTransportationAccountingCodes(models.DepartmentIndicatorARMY, ordersIssueDate, tacCode, appCtx)
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
		loas, err := suite.loaFetcher.FetchLongLinesOfAccounting(models.DepartmentIndicatorARMY, ordersIssueDate, tacCode, appCtx)
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
		tacs, err := suite.tacFetcher.FetchOrderTransportationAccountingCodes(models.DepartmentIndicatorARMY, ordersIssueDate, tacCode, appCtx)
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
		loas, err := suite.loaFetcher.FetchLongLinesOfAccounting(models.DepartmentIndicatorARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		// Ensure we got 1 loa, FBMC P should've been filtered out
		suite.Len(loas, 1)
	})
	suite.Run("Checks for valid HHG program code and a valid LOA for a given TAC", func() {
		appCtx := suite.AppContextForTest()
		ordersIssueDate, startDate, endDate, tacCode := setupTest()
		loaFY := 12
		loa := factory.BuildLineOfAccounting(appCtx.DB(), []factory.Customization{
			{
				Model: models.LineOfAccounting{
					LoaBgnDt:               &startDate,
					LoaEndDt:               &endDate,
					LoaSysID:               models.StringPointer("1234567890"),
					LoaHsGdsCd:             models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
					LoaDptID:               models.StringPointer("1"),
					LoaTnsfrDptNm:          models.StringPointer("1"),
					LoaBafID:               models.StringPointer("1"),
					LoaTrsySfxTx:           models.StringPointer("1"),
					LoaMajClmNm:            models.StringPointer("1"),
					LoaOpAgncyID:           models.StringPointer("1"),
					LoaAlltSnID:            models.StringPointer("1"),
					LoaPgmElmntID:          models.StringPointer("1"),
					LoaTskBdgtSblnTx:       models.StringPointer("1"),
					LoaDfAgncyAlctnRcpntID: models.StringPointer("1"),
					LoaJbOrdNm:             models.StringPointer("1"),
					LoaSbaltmtRcpntID:      models.StringPointer("1"),
					LoaWkCntrRcpntNm:       models.StringPointer("1"),
					LoaMajRmbsmtSrcID:      models.StringPointer("1"),
					LoaDtlRmbsmtSrcID:      models.StringPointer("1"),
					LoaCustNm:              models.StringPointer("1"),
					LoaObjClsID:            models.StringPointer("1"),
					LoaSrvSrcID:            models.StringPointer("1"),
					LoaSpclIntrID:          models.StringPointer("1"),
					LoaBdgtAcntClsNm:       models.StringPointer("1"),
					LoaDocID:               models.StringPointer("1"),
					LoaClsRefID:            models.StringPointer("1"),
					LoaInstlAcntgActID:     models.StringPointer("1"),
					LoaLclInstlID:          models.StringPointer("1"),
					LoaFmsTrnsactnID:       models.StringPointer("1"),
					LoaTrnsnID:             models.StringPointer("1"),
					LoaUic:                 models.StringPointer("1"),
					LoaBgFyTx:              &loaFY,
					LoaEndFyTx:             &loaFY,
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
		tacs, err := suite.tacFetcher.FetchOrderTransportationAccountingCodes(models.DepartmentIndicatorARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		suite.NotEmpty(tacs)
		// Ensure LOA isn't nil
		suite.NotNil(tacs[0].LineOfAccounting)
		// Extract LOAs
		loas, err := suite.loaFetcher.FetchLongLinesOfAccounting(models.DepartmentIndicatorARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		suite.Equal(loa.ID, loas[0].ID)

		ValidHhgProgramCodeForLoaReturnValue := true
		ValidLoaForTacReturnValue := true
		suite.Equal(loas[0].ValidHhgProgramCodeForLoa, &ValidHhgProgramCodeForLoaReturnValue)
		suite.Equal(loas[0].ValidLoaForTac, &ValidLoaForTacReturnValue)
	})
	suite.Run("Checks for invalid HHG program code and an invalid LOA for a given TAC", func() {
		appCtx := suite.AppContextForTest()
		ordersIssueDate, startDate, endDate, tacCode := setupTest()
		loa := factory.BuildLineOfAccounting(appCtx.DB(), []factory.Customization{
			{
				Model: models.LineOfAccounting{
					LoaBgnDt:               &startDate,
					LoaEndDt:               &endDate,
					LoaSysID:               models.StringPointer("1234567890"),
					LoaHsGdsCd:             models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
					LoaDptID:               models.StringPointer("1"),
					LoaTnsfrDptNm:          models.StringPointer("1"),
					LoaBafID:               models.StringPointer("1"),
					LoaTrsySfxTx:           models.StringPointer("1"),
					LoaMajClmNm:            models.StringPointer("1"),
					LoaOpAgncyID:           models.StringPointer("1"),
					LoaAlltSnID:            models.StringPointer("1"),
					LoaPgmElmntID:          models.StringPointer("1"),
					LoaTskBdgtSblnTx:       models.StringPointer("1"),
					LoaDfAgncyAlctnRcpntID: models.StringPointer("1"),
					LoaJbOrdNm:             models.StringPointer("1"),
					LoaSbaltmtRcpntID:      models.StringPointer("1"),
					LoaWkCntrRcpntNm:       models.StringPointer("1"),
					LoaMajRmbsmtSrcID:      models.StringPointer("1"),
					LoaDtlRmbsmtSrcID:      models.StringPointer("1"),
					LoaCustNm:              models.StringPointer("1"),
					LoaObjClsID:            models.StringPointer("1"),
					LoaSrvSrcID:            models.StringPointer("1"),
					LoaSpclIntrID:          models.StringPointer("1"),
					LoaBdgtAcntClsNm:       models.StringPointer("1"),
					LoaDocID:               models.StringPointer("1"),
					LoaClsRefID:            models.StringPointer("1"),
					LoaInstlAcntgActID:     models.StringPointer("1"),
					LoaLclInstlID:          models.StringPointer("1"),
					LoaFmsTrnsactnID:       models.StringPointer("1"),
					LoaDscTx:               models.StringPointer("1"),
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
		tacs, err := suite.tacFetcher.FetchOrderTransportationAccountingCodes(models.DepartmentIndicatorARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		suite.NotEmpty(tacs)
		// Ensure LOA isn't nil
		suite.NotNil(tacs[0].LineOfAccounting)

		// Extract LOAs
		loas, err := suite.loaFetcher.FetchLongLinesOfAccounting(models.DepartmentIndicatorARMY, ordersIssueDate, tacCode, appCtx)
		suite.NoError(err)
		suite.Equal(loa.ID, loas[0].ID)

		// The LoaUic field is missing, so the ValidLoaForTac field should be false
		ValidLoaForTacReturnValue := false
		suite.Equal(loas[0].ValidLoaForTac, &ValidLoaForTacReturnValue)

		// In the event that the LoaHsGdsCd is for some reason nil, the ValidHhgProgramCodeForLoa should be false
		// Forcing LoaHsGdsCd to be nil since it can't be nil earlier in the test
		loas[0].LoaHsGdsCd = nil
		loasWithValidityCheck, err := checkForValidHhgProgramCodeForLoaAndValidLoaForTac(loas, appCtx)
		// This is a "soft warning" and shouldn't prevent the return of the LOA information, so err is nil
		suite.Nil(err)
		ValidHhgProgramCodeForLoaReturnValue := false
		// When LoaHsGdsCd is nil, so the ValidHhgProgramCodeForLoa should be false
		suite.Equal(loasWithValidityCheck[0].ValidHhgProgramCodeForLoa, &ValidHhgProgramCodeForLoaReturnValue)
	})
}
