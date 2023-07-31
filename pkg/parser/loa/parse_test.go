package loa_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/parser/loa"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type LoaParserSuite struct {
	*testingsuite.PopTestSuite
	txtFilename string
	txtContent  []byte
}

func TestLoaParserSuite(t *testing.T) {
	hs := &LoaParserSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		txtFilename:  "./fixtures/Line Of Accounting.txt",
	}

	var err error
	hs.txtContent, err = os.ReadFile(hs.txtFilename)
	if err != nil {
		hs.Logger().Panic("could not read text file", zap.Error(err))
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

// Test that the parser correctly handles the test file and reports back at least
// one correct Line of Accounting
func (suite *LoaParserSuite) TestParsing() {
	reader := bytes.NewReader(suite.txtContent)

	// Parse the text file content
	codes, err := loa.Parse(reader)
	suite.NoError(err)

	// Assuming the txt file has at least one record
	suite.NotEmpty(codes)

	// Do a hard coded check to the first line of data to ensure a 1:1 match to what is expected.
	firstCode := codes[0]
	suite.Equal("124641", firstCode.LOA)
	suite.Equal("97", firstCode.DepartmentID)
	suite.Equal("", firstCode.TransferDepartmentName)
	suite.Equal("4930", firstCode.BasicAppropriationFundID)
	suite.Equal("AA37", firstCode.TreasurySuffixText)
	suite.Equal("", firstCode.MajorClaimantName)
	suite.Equal("6D", firstCode.OperatingAgencyID)
	suite.Equal("0000", firstCode.AllotmentSerialNumberID)
	suite.Equal("MZZF0000", firstCode.ProgramElementID)
	suite.Equal("", firstCode.TaskBudgetSublineText)
	suite.Equal("", firstCode.DefenseAgencyAllocationRecipientID)
	suite.Equal("", firstCode.JobOrderName)
	suite.Equal("", firstCode.SubAllotmentRecipientId)
	suite.Equal("D0000", firstCode.WorkCenterRecipientName)
	suite.Equal("", firstCode.MajorReimbursementSourceID)
	suite.Equal("", firstCode.DetailReimbursementSourceID)
	suite.Equal("", firstCode.CustomerName)
	suite.Equal("22NL", firstCode.ObjectClassID)
	suite.Equal("", firstCode.ServiceSourceID)
	suite.Equal("", firstCode.SpecialInterestID)
	suite.Equal("MA1MN4", firstCode.BudgetAccountClassificationName)
	suite.Equal("FRT1MNUSAF7790", firstCode.DocumentID)
	suite.Equal("", firstCode.ClassReferenceID)
	suite.Equal("011000", firstCode.InstallationAccountingActivityID)
	suite.Equal("", firstCode.LocalInstallationID)
	suite.Equal("", firstCode.FMSTransactionID)
	suite.Equal("CONUS ENTRY", firstCode.DescriptionText)
	suite.Equal(time.Date(2006, 10, 1, 0, 0, 0, 0, time.UTC), firstCode.BeginningDate)
	suite.Equal(time.Date(2007, 9, 30, 0, 0, 0, 0, time.UTC), firstCode.EndDate)
	suite.Equal("", firstCode.FunctionalPersonName)
	suite.Equal("U", firstCode.StatusCode)
	suite.Equal("", firstCode.HistoryStatusCode)
	suite.Equal("", firstCode.HouseholdGoodsCode)
	suite.Equal("DZ", firstCode.OrganizationGroupDefenseFinanceAccountingServiceCode)
	suite.Equal("", firstCode.UnitIdentificationCode)
	suite.Equal("", firstCode.TransactionID)
	suite.Equal("", firstCode.SubordinateAccountID)
	suite.Equal("", firstCode.BusinessEventTypeCode)
	suite.Equal("", firstCode.FundTypeFlagCode)
	suite.Equal("", firstCode.BudgetLineItemID)
	suite.Equal("", firstCode.SecurityCooperationImplementingAgencyCode)
	suite.Equal("", firstCode.SecurityCooperationDesignatorID)
	suite.Equal("", firstCode.SecurityCooperationLineItemID)
	suite.Equal("", firstCode.AgencyDisbursingCode)
	suite.Equal("", firstCode.AgencyAccountingCode)
	suite.Equal("", firstCode.FundCenterID)
	suite.Equal("", firstCode.CostCenterID)
	suite.Equal("", firstCode.ProjectTaskID)
	suite.Equal("", firstCode.ActivityID)
	suite.Equal("", firstCode.CostCode)
	suite.Equal("", firstCode.WorkOrderID)
	suite.Equal("", firstCode.FunctionalAreaID)
	suite.Equal("", firstCode.SecurityCooperationCustomerCode)
	suite.Equal(0, firstCode.EndingFiscalYear)
	suite.Equal(0, firstCode.BeginningFiscalYear)
	suite.Equal("", firstCode.BudgetRestrictionCode)
	suite.Equal("", firstCode.BudgetSubActivityCode)
}

// This test will ensure that the parse function errors on an empty file.
func (suite *LoaParserSuite) TestEmptyFileContent() {
	reader := bytes.NewReader([]byte(""))

	// Attempt to parse an empty file
	_, err := loa.Parse(reader)
	suite.Error(err)
}

// There are 57 expected values per line entry. This test will make sure
// an error is reported if it is not met.
func (suite *LoaParserSuite) TestIncorrectNumberOfValuesInLine() {
	// !Warning, do not touch the format of the byte
	content := []byte(`Unclassified
TAC_SYS_ID|LOA_SYS_ID|TRNSPRTN_ACNT_CD
1234567884061|12345678
Unclassified`)
	reader := bytes.NewReader(content)

	// Attempt to parse the malformed file
	_, err := loa.Parse(reader)
	suite.Error(err)
}

// Test for good data, but bad column headers. Aka, check that the expected
// fields are received from the .txt file.
// This test adds a blank || column header
func (suite *LoaParserSuite) TestColumnHeadersDoNotMatch() {
	// !Warning, do not touch the format of the byte
	content := []byte(`Unclassified
LOA_SYS_ID|LOA_DPT_ID|LOA_TNSFR_DPT_NM|LOA_BAF_ID|LOA_TRSY_SFX_TX|LOA_MAJ_CLM_NM|LOA_OP_AGNCY_ID|LOA_ALLT_SN_ID|LOA_PGM_ELMNT_ID|LOA_TSK_BDGT_SBLN_TX|LOA_DF_AGNCY_ALCTN_RCPNT_ID|LOA_JB_ORD_NM|LOA_SBALTMT_RCPNT_ID|LOA_WK_CNTR_RCPNT_NM|LOA_MAJ_RMBSMT_SRC_ID|LOA_DTL_RMBSMT_SRC_ID|LOA_CUST_NM|LOA_OBJ_CLS_ID|LOA_SRV_SRC_ID|LOA_SPCL_INTR_ID|LOA_BDGT_ACNT_CLS_NM|LOA_DOC_ID|LOA_CLS_REF_ID|LOA_INSTL_ACNTG_ACT_ID|LOA_LCL_INSTL_ID|LOA_FMS_TRNSACTN_ID|LOA_DSC_TX|LOA_BGN_DT|LOA_END_DT|LOA_FNCT_PRS_NM|LOA_STAT_CD|LOA_HIST_STAT_CD|LOA_HS_GDS_CD|ORG_GRP_DFAS_CD|LOA_UIC|LOA_TRNSN_ID|LOA_SUB_ACNT_ID|LOA_BET_CD|LOA_FND_TY_FG_CD|LOA_BGT_LN_ITM_ID|LOA_SCRTY_COOP_IMPL_AGNC_CD|LOA_SCRTY_COOP_DSGNTR_CD|LOA_SCRTY_COOP_LN_ITM_ID|LOA_AGNC_DSBR_CD|LOA_AGNC_ACNTNG_CD|LOA_FND_CNTR_ID|LOA_CST_CNTR_ID|LOA_PRJ_ID|LOA_ACTVTY_ID|LOA_CST_CD|LOA_WRK_ORD_ID|LOA_FNCL_AR_ID|LOA_SCRTY_COOP_CUST_CD|LOA_END_FY_TX|LOA_BG_FY_TX|LOA_BGT_RSTR_CD|LOA_BGT_SUB_ACT_CD||
124641|97||4930|AA37||6D|0000|MZZF0000|||||D0000||||22NL|||MA1MN4|FRT1MNUSAF7790||011000|||CONUS ENTRY|2006-10-01 00:00:00|2007-09-30 00:00:00||U|||DZ|||||||||||||||||||||||
124642|97||4930|AA37||6D|0000|MZZF0000|||||D0000||||22N2|||MA1MN4|FRT1MNUSAF8790||011000|||OCONUS ENTRY|2006-10-01 00:00:00|2007-09-30 00:00:00||U|||DZ|||||||||||||||||||||||
124643|97||4930|AA37||6D|0000|MZZF0000|||||D0000||||22NL|||P50MD6|FRT0MDUSAF9790||011000|||ENTRY TYPE A|2006-10-01 00:00:00|2007-09-30 00:00:00||U|||DZ|||||||||||||||||||||||
Unclassified`)
	reader := bytes.NewReader(content)

	// Attempt to parse the malformed file
	_, err := loa.Parse(reader)

	suite.Error(err)
}

// This function will test the pruning of all expired TACs when called.
// TODO: Add factory
func (suite *LoaParserSuite) TestExpiredTACs() {

	// Create a loa with an empty household good code
	emptyHhgLoa := models.LineOfAccountingDesiredFromTRDM{
		LOA:                                "124641",
		DepartmentID:                       "97",
		TransferDepartmentName:             "",
		BasicAppropriationFundID:           "4930",
		TreasurySuffixText:                 "AA37",
		MajorClaimantName:                  "",
		OperatingAgencyID:                  "6D",
		AllotmentSerialNumberID:            "0000",
		ProgramElementID:                   "MZZF0000",
		TaskBudgetSublineText:              "",
		DefenseAgencyAllocationRecipientID: "",
		JobOrderName:                       "",
		SubAllotmentRecipientId:            "",
		WorkCenterRecipientName:            "D0000",
		MajorReimbursementSourceID:         "",
		DetailReimbursementSourceID:        "",
		CustomerName:                       "",
		ObjectClassID:                      "22NL",
		ServiceSourceID:                    "",
		SpecialInterestID:                  "",
		BudgetAccountClassificationName:    "MA1MN4",
		DocumentID:                         "FRT1MNUSAF7790",
		ClassReferenceID:                   "",
		InstallationAccountingActivityID:   "011000",
		LocalInstallationID:                "",
		FMSTransactionID:                   "",
		DescriptionText:                    "CONUS ENTRY",
		BeginningDate:                      time.Date(2006, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:                            time.Date(2007, 9, 30, 0, 0, 0, 0, time.UTC),
		FunctionalPersonName:               "",
		StatusCode:                         "U",
		HistoryStatusCode:                  "",
		HouseholdGoodsCode:                 "",
		OrganizationGroupDefenseFinanceAccountingServiceCode: "DZ",
		UnitIdentificationCode:                               "",
		TransactionID:                                        "",
		SubordinateAccountID:                                 "",
		BusinessEventTypeCode:                                "",
		FundTypeFlagCode:                                     "",
		BudgetLineItemID:                                     "",
		SecurityCooperationImplementingAgencyCode:            "",
		SecurityCooperationDesignatorID:                      "",
		SecurityCooperationLineItemID:                        "",
		AgencyDisbursingCode:                                 "",
		AgencyAccountingCode:                                 "",
		FundCenterID:                                         "",
		CostCenterID:                                         "",
		ProjectTaskID:                                        "",
		ActivityID:                                           "",
		CostCode:                                             "",
		WorkOrderID:                                          "",
		FunctionalAreaID:                                     "",
		SecurityCooperationCustomerCode:                      "",
		EndingFiscalYear:                                     0,
		BeginningFiscalYear:                                  0,
		BudgetRestrictionCode:                                "",
		BudgetSubActivityCode:                                "",
	}

	// Attempt to prune all expired TACs
	parsedLOAs := []models.LineOfAccountingDesiredFromTRDM{emptyHhgLoa}
	prunedLOAs := loa.PruneEmptyHhgCodes(parsedLOAs)

	// Check that the expired LOA was properly removed
	suite.NotContains(prunedLOAs, emptyHhgLoa)
}
