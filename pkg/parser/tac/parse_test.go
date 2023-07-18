package tac_test

import (
	"bytes"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"
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

// Test that the parser correctly handles the test file and reports back at least
// one correct Transportation Accounting Code
func (suite *TacParserSuite) TestParsing() {
	reader := bytes.NewReader(suite.txtContent)

	// Parse the text file content
	codes, err := tac.Parse(reader)
	suite.NoError(err)

	// Assuming the txt file has at least one record
	suite.NotEmpty(codes)

	// Do a hard coded check to the first line of data to ensure a 1:1 match to what is expected.
	firstCode := codes[0]
	suite.Equal("0003", firstCode.TAC)
	suite.Equal("FIRST LINE", firstCode.BillingAddressFirstLine)
	suite.Equal("SECOND LINE", firstCode.BillingAddressSecondLine)
	suite.Equal("THIRD LINE", firstCode.BillingAddressThirdLine)
	suite.Equal("FOURTH LINE", firstCode.BillingAddressFourthLine)
	suite.Equal("FOR MOVEMENT TEST 1", firstCode.Transaction)
	suite.Equal("2021-10-01 00:00:00", firstCode.EffectiveDate)
	suite.Equal("2022-09-30 00:00:00", firstCode.ExpirationDate)
	suite.Equal("2022", firstCode.FiscalYear)
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
