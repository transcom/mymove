package tac_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/parser/tac"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type TacParserSuite struct {
	*testingsuite.PopTestSuite
	txtFilename string
	txtContent  []byte
}

func TestTacParserSuite(t *testing.T) {
	hs := &TacParserSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		txtFilename:  "./fixtures/Transportation Account.txt",
	}

	var err error
	hs.txtContent, err = os.ReadFile(hs.txtFilename)
	if err != nil {
		hs.Logger().Panic("could not read text file", zap.Error(err))
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *TacParserSuite) TestParsing() {
	reader := bytes.NewReader(suite.txtContent)

	// Parse the text file content
	codes, err := tac.Parse(reader)
	suite.NoError(err)

	// Assuming the txt file has at least one record
	suite.NotEmpty(codes)
	var bgnDt = time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)
	var endDt = time.Date(2022, 9, 30, 0, 0, 0, 0, time.UTC)

	// Create expected TransportationAccountingCode
	expected := models.TransportationAccountingCode{
		TacSysID:           models.IntPointer(1234567884061),
		LoaSysID:           models.IntPointer(12345678),
		TAC:                "0003",
		TacFyTxt:           models.IntPointer(2022),
		TacFnBlModCd:       models.StringPointer("3"),
		OrgGrpDfasCd:       models.StringPointer("DF"),
		TacMvtDsgID:        models.StringPointer(""),
		TacTyCd:            models.StringPointer("O"),
		TacUseCd:           models.StringPointer("O"),
		TacMajClmtID:       models.StringPointer("USTC"),
		TacBillActTxt:      models.StringPointer(""),
		TacCostCtrNm:       models.StringPointer("G31M32"),
		Buic:               models.StringPointer(""),
		TacHistCd:          models.StringPointer(""),
		TacStatCd:          models.StringPointer("I"),
		TrnsprtnAcntTx:     models.StringPointer("FOR MOVEMENT TEST 1"),
		TrnsprtnAcntBgnDt:  models.TimePointer(bgnDt),
		TrnsprtnAcntEndDt:  models.TimePointer(endDt),
		DdActvtyAdrsID:     models.StringPointer("F55555"),
		TacBlldAddFrstLnTx: models.StringPointer("FIRST LINE"),
		TacBlldAddScndLnTx: models.StringPointer("SECOND LINE"),
		TacBlldAddThrdLnTx: models.StringPointer("THIRD LINE"),
		TacBlldAddFrthLnTx: models.StringPointer("FOURTH LINE"),
		TacFnctPocNm:       models.StringPointer("Contact Person Here"),
	}

	// Do a hard coded check to the first line of data to ensure a 1:1 match to what is expected.
	firstCode := codes[0]
	suite.Equal(expected, firstCode)
}

// This test will ensure that the parse function errors on an empty file.
func (suite *TacParserSuite) TestEmptyFileContent() {
	reader := bytes.NewReader([]byte(""))

	// Attempt to parse an empty file
	_, err := tac.Parse(reader)
	suite.Error(err)
}

// There are 23 expected values per line entry. This test will make sure
// an error is reported if it is not met.
func (suite *TacParserSuite) TestIncorrectNumberOfValuesInLine() {
	// !Warning, do not touch the format of the byte
	content := []byte(`Unclassified
TAC_SYS_ID|LOA_SYS_ID|TRNSPRTN_ACNT_CD
1234567884061|12345678
Unclassified`)
	reader := bytes.NewReader(content)

	// Attempt to parse the malformed file
	_, err := tac.Parse(reader)
	suite.Error(err)
}

// Test for good data, but bad column headers. Aka, check that the expected
// fields are received from the .txt file.
// This test adds a blank column header "||"
func (suite *TacParserSuite) TestColumnHeadersDoNotMatch() {
	// !Warning, do not touch the format of the byte
	content := []byte(`Unclassified
TAC_SYS_ID|LOA_SYS_ID|TRNSPRTN_ACNT_CD|TAC_FY_TXT|TAC_FN_BL_MOD_CD|ORG_GRP_DFAS_CD|TAC_MVT_DSG_ID|TAC_TY_CD|TAC_USE_CD|TAC_MAJ_CLMT_ID|TAC_BILL_ACT_TXT|TAC_COST_CTR_NM|BUIC|TAC_HIST_CD|TAC_STAT_CD|TRNSPRTN_ACNT_TX|TRNSPRTN_ACNT_BGN_DT|TRNSPRTN_ACNT_END_DT|DD_ACTVTY_ADRS_ID|TAC_BLLD_ADD_FRST_LN_TX|TAC_BLLD_ADD_SCND_LN_TX|TAC_BLLD_ADD_THRD_LN_TX|TAC_BLLD_ADD_FRTH_LN_TX|TAC_FNCT_POC_NM|
1234567884061|12345678|0003|2022|3|DF||O|O|USTC||G31M32|||I|FOR MOVEMENT TEST 1|2021-10-01 00:00:00|2022-09-30 00:00:00|F55555|FIRST LINE|SECOND LINE|THIRD LINE|FOURTH LINE|Contact Person Here
3456789|34567890|ZZQE|2022|M|HS||O|N|018301|051800|018301|||I|FOR MOVEMENT TEST 2|2021-10-01 00:00:00|2022-09-30 00:00:00|Z55555|FIRST LINE|SECOND LINE||TAMPA FL 33621|NotARealPerson@USCG.MIL
0000000|00000000|ZZQE|2022|W|HS||O|N|018301|051800|018301|||I|FOR MOVEMENT TEST 2|2021-10-01 00:00:00|2022-09-30 00:00:00|Z55555|FIRST LINE|SECOND LINE||TAMPA FL 33621|NotARealPerson@USCG.MIL
Unclassified`)
	reader := bytes.NewReader(content)

	// Attempt to parse the malformed file
	_, err := tac.Parse(reader)

	suite.Error(err)
}

// This function will test the pruning of all expired TACs when called.
func (suite *TacParserSuite) TestExpiredTACs() {

	// Create a TAC
	expiredTac := factory.BuildFullTransportationAccountingCode(suite.DB())

	// Make it expired
	*expiredTac.TrnsprtnAcntBgnDt = time.Now().AddDate(-1, 0, 0) // A year ago
	*expiredTac.TrnsprtnAcntEndDt = time.Now().AddDate(0, 0, -1) // A day ago

	// Attempt to prune all expired TACs
	parsedTACs := []models.TransportationAccountingCode{expiredTac}
	prunedTACs := tac.PruneExpiredTACs(parsedTACs)

	// Check that the expired TAC was properly removed
	suite.NotContains(prunedTACs, expiredTac)
}

// This function will test the conslidation of two TACs with matching "TAC" and "ExpirationDate" values, but that have a difference in other values.
// It is expected to combine their transaction descriptions and preserve the first code found in the array
func (suite *TacParserSuite) TestDuplicateTACsWithDifferentValuesAndEquivalentExpirationDates() {

	// Create duplicate TACs
	tac1 := factory.BuildFullTransportationAccountingCode(suite.DB())
	tac2 := tac1
	// Set the second TAC to have a different transaction description
	tac2.TrnsprtnAcntTx = models.StringPointer("Different")

	// Create the expected TAC value for comparison
	expectedConsolidatedTAC := tac1
	*expectedConsolidatedTAC.TrnsprtnAcntTx = *tac1.TrnsprtnAcntTx + *tac2.TrnsprtnAcntTx

	parsedTACs := []models.TransportationAccountingCode{tac1, tac2}
	consolidatedTACs := tac.ConsolidateDuplicateTACsDesiredFromTRDM(parsedTACs)

	suite.Contains(consolidatedTACs, expectedConsolidatedTAC)
}

// This function will test the conslidation of two TACs with matching "TAC" values, but that have a difference in other values.
// It is expected to combine their transaction descriptions and preserve the code with the expiration date further in the future.
// The expiration dates will be different.
func (suite *TacParserSuite) TestDuplicateTACsWithDifferentValuesAndDifferentExpirationDates() {
	oneYearAhead := time.Now().AddDate(1, 0, 0)    // A year from now
	twoYearsAhead := oneYearAhead.AddDate(1, 0, 0) // Two years from now

	// Create duplicate TACs
	tac1 := factory.BuildFullTransportationAccountingCode(suite.DB())
	tac2 := tac1

	// Set the first TAC to have a new expiration date
	tac1.TrnsprtnAcntEndDt = &oneYearAhead

	// Set the second TAC to have a different transaction description and expiration date
	tac2.TrnsprtnAcntTx = models.StringPointer("Different")
	tac2.TrnsprtnAcntEndDt = &twoYearsAhead

	parsedTACs := []models.TransportationAccountingCode{tac1, tac2}
	consolidatedTACs := tac.ConsolidateDuplicateTACsDesiredFromTRDM(parsedTACs)

	// Create the expected TAC value for comparison
	// Tac 2 expires a year after tac 1
	expectedConsolidatedTAC := tac2
	*expectedConsolidatedTAC.TrnsprtnAcntTx = *tac1.TrnsprtnAcntTx + *tac2.TrnsprtnAcntTx

	suite.Contains(consolidatedTACs, expectedConsolidatedTAC)
}
