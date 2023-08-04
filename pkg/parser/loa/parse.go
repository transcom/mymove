package loa

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/models"
)

var expectedColumnNames = []string{"LOA_SYS_ID", "LOA_DPT_ID", "LOA_TNSFR_DPT_NM", "LOA_BAF_ID", "LOA_TRSY_SFX_TX", "LOA_MAJ_CLM_NM", "LOA_OP_AGNCY_ID", "LOA_ALLT_SN_ID", "LOA_PGM_ELMNT_ID", "LOA_TSK_BDGT_SBLN_TX", "LOA_DF_AGNCY_ALCTN_RCPNT_ID", "LOA_JB_ORD_NM", "LOA_SBALTMT_RCPNT_ID", "LOA_WK_CNTR_RCPNT_NM", "LOA_MAJ_RMBSMT_SRC_ID", "LOA_DTL_RMBSMT_SRC_ID", "LOA_CUST_NM", "LOA_OBJ_CLS_ID", "LOA_SRV_SRC_ID", "LOA_SPCL_INTR_ID", "LOA_BDGT_ACNT_CLS_NM", "LOA_DOC_ID", "LOA_CLS_REF_ID", "LOA_INSTL_ACNTG_ACT_ID", "LOA_LCL_INSTL_ID", "LOA_FMS_TRNSACTN_ID", "LOA_DSC_TX", "LOA_BGN_DT", "LOA_END_DT", "LOA_FNCT_PRS_NM", "LOA_STAT_CD", "LOA_HIST_STAT_CD", "LOA_HS_GDS_CD", "ORG_GRP_DFAS_CD", "LOA_UIC", "LOA_TRNSN_ID", "LOA_SUB_ACNT_ID", "LOA_BET_CD", "LOA_FND_TY_FG_CD", "LOA_BGT_LN_ITM_ID", "LOA_SCRTY_COOP_IMPL_AGNC_CD", "LOA_SCRTY_COOP_DSGNTR_CD", "LOA_SCRTY_COOP_LN_ITM_ID", "LOA_AGNC_DSBR_CD", "LOA_AGNC_ACNTNG_CD", "LOA_FND_CNTR_ID", "LOA_CST_CNTR_ID", "LOA_PRJ_ID", "LOA_ACTVTY_ID", "LOA_CST_CD", "LOA_WRK_ORD_ID", "LOA_FNCL_AR_ID", "LOA_SCRTY_COOP_CUST_CD", "LOA_END_FY_TX", "LOA_BG_FY_TX", "LOA_BGT_RSTR_CD", "LOA_BGT_SUB_ACT_CD"}

const timeConversionMethod = "2006-01-02 15:04:05"

// Parse the pipe delimited .txt file with the following assumptions:
// 1. The first and last lines are the security classification.
// 2. The second line of the file are the columns that will be a 1:1 match to the LineOfAccountingTrdmFileRecord struct in pipe delimited format.
// 3. There are 57 values per line, excluding the security classification. Again, to know what these values are refer to note #2.
// 4. All values are in pipe delimited format.
// 5. Null values will be present, but it may be desired to filter out LOAs with a null LOA_HS_GDS_CD
func Parse(file io.Reader) ([]models.LineOfAccounting, error) {

	// Init variables
	var codes []models.LineOfAccounting
	scanner := bufio.NewScanner(file)
	var columnHeaders []string

	// Skip first line as it does not hold any necessary data for parsing.
	// Additionally, this will check if it is empty.
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return nil, errors.New("empty file")
	}

	// Read the second line to get the column names from the TRDM file, verify them,
	// and then proceed with parsing the rest of the file.
	if scanner.Scan() {
		columnHeaders = strings.Split(scanner.Text(), "|")
		err := ensureFileStructMatchesColumnNames(columnHeaders)
		if err != nil {
			return nil, err
		}
	}

	// Process the lines of the .txt file into modeled codes
	codes, err := processLines(scanner, columnHeaders, codes)
	if err != nil {
		return nil, err
	}

	return codes, nil
}

// Compare a struct's field names to the columns retrieved from the .txt file
func ensureFileStructMatchesColumnNames(columnNames []string) error {
	if len(columnNames) == 0 {
		return errors.New("column names were not parsed properly from the second line of the loa file")
	}

	if !reflect.DeepEqual(columnNames, expectedColumnNames) {
		return errors.New("column names parsed do not match the expected format of loa file records")
	}
	return nil
}

// Removes all LOAs with an empty HHG code
func PruneEmptyHhgCodes(codes []models.LineOfAccounting) []models.LineOfAccounting {
	var pruned []models.LineOfAccounting

	// If the household goods code is not empty, then it should be appended to the pruned array for return
	for _, code := range codes {
		if *code.LoaHsGdsCd != "" {
			pruned = append(pruned, code)
		}
	}

	return pruned
}

// This function handles the heavy lifting for the main parse function. It processes every line from the .txt file into a proper struct
func processLines(scanner *bufio.Scanner, columnHeaders []string, codes []models.LineOfAccounting) ([]models.LineOfAccounting, error) {
	// Scan every line and parse into the desired Line of Accounting codes
	for scanner.Scan() {
		line := scanner.Text()
		var beginningDate time.Time
		var endingDate time.Time
		var endingFY int
		var beginningFY int
		var err error

		// This check will skip the last line of the file.
		if line == "Unclassified" {
			break
		}

		values := strings.Split(line, "|")
		if len(values) != len(columnHeaders) {
			return nil, errors.New("malformed line in the provided loa file: " + line)
		}

		// Check that the LOA sys id is not empty
		if values[0] == "" {
			return nil, errors.New("malformed line in the provided loa file: " + line)
		}

		loaSysID, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, errors.New("malformed line in the provided loa file: " + line)
		}

		// Check if beginning date and expired date are not blank. If not blank, run time conversion, if not, leave empty.
		if values[27] != "" && values[28] != "" {
			// Parse values[27], this is the beginning date
			parsedDate, parseError := time.Parse(timeConversionMethod, values[27])
			if parseError != nil {
				return nil, fmt.Errorf("malformed effective date in the provided loa file: %s", parseError)
			}
			beginningDate = parsedDate

			// Parse values[28], this is the ending date
			parsedDate, parseError = time.Parse(timeConversionMethod, values[28])
			if parseError != nil {
				return nil, fmt.Errorf("malformed effective date in the provided loa file: %s", parseError)
			}
			endingDate = parsedDate
		}

		// Check if beginning and ending fiscal years are not blank. If not blank, run str to int conversion, if not, leave empty.
		if values[54] != "" && values[55] != "" {
			endingFY, err = strconv.Atoi(values[54])
			if err != nil {
				return nil, fmt.Errorf("malformed ending fiscal year int in the provided loa file: %s", err)
			}

			beginningFY, err = strconv.Atoi(values[55])
			if err != nil {
				return nil, fmt.Errorf("malformed beginning fiscal year int in the provided loa file: %s", err)
			}
		}

		code := models.LineOfAccounting{
			LoaSysID:               &loaSysID,
			LoaDptID:               &values[1],
			LoaTnsfrDptNm:          &values[2],
			LoaBafID:               &values[3],
			LoaTrsySfxTx:           &values[4],
			LoaMajClmNm:            &values[5],
			LoaOpAgncyID:           &values[6],
			LoaAlltSnID:            &values[7],
			LoaPgmElmntID:          &values[8],
			LoaTskBdgtSblnTx:       &values[9],
			LoaDfAgncyAlctnRcpntID: &values[10],
			LoaJbOrdNm:             &values[11],
			LoaSbaltmtRcpntID:      &values[12],
			LoaWkCntrRcpntNm:       &values[13],
			LoaMajRmbsmtSrcID:      &values[14],
			LoaDtlRmbsmtSrcID:      &values[15],
			LoaCustNm:              &values[16],
			LoaObjClsID:            &values[17],
			LoaSrvSrcID:            &values[18],
			LoaSpclIntrID:          &values[19],
			LoaBdgtAcntClsNm:       &values[20],
			LoaDocID:               &values[21],
			LoaClsRefID:            &values[22],
			LoaInstlAcntgActID:     &values[23],
			LoaLclInstlID:          &values[24],
			LoaFmsTrnsactnID:       &values[25],
			LoaDscTx:               &values[26],
			LoaBgnDt:               &beginningDate,
			LoaEndDt:               &endingDate,
			LoaFnctPrsNm:           &values[29],
			LoaStatCd:              &values[30],
			LoaHistStatCd:          &values[31],
			LoaHsGdsCd:             &values[32],
			OrgGrpDfasCd:           &values[33],
			LoaUic:                 &values[34],
			LoaTrnsnID:             &values[35],
			LoaSubAcntID:           &values[36],
			LoaBetCd:               &values[37],
			LoaFndTyFgCd:           &values[38],
			LoaBgtLnItmID:          &values[39],
			LoaScrtyCoopImplAgncCd: &values[40],
			LoaScrtyCoopDsgntrCd:   &values[41],
			LoaScrtyCoopLnItmID:    &values[42],
			LoaAgncDsbrCd:          &values[43],
			LoaAgncAcntngCd:        &values[44],
			LoaFndCntrID:           &values[45],
			LoaCstCntrID:           &values[46],
			LoaPrjID:               &values[47],
			LoaActvtyID:            &values[48],
			LoaCstCd:               &values[49],
			LoaWrkOrdID:            &values[50],
			LoaFnclArID:            &values[51],
			LoaScrtyCoopCustCd:     &values[52],
			LoaEndFyTx:             &endingFY,
			LoaBgFyTx:              &beginningFY,
			LoaBgtRstrCd:           &values[55],
			LoaBgtSubActCd:         &values[56],
		}

		codes = append(codes, code)
	}

	return codes, scanner.Err()
}
