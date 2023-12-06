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

var expectedColumnNames = []string{"LOA_SYS_ID", "LOA_DPT_ID", "LOA_TNSFR_DPT_NM", "LOA_BAF_ID", "LOA_TRSY_SFX_TX", "LOA_MAJ_CLM_NM", "LOA_OP_AGNCY_ID", "LOA_ALLT_SN_ID", "LOA_PGM_ELMNT_ID", "LOA_TSK_BDGT_SBLN_TX", "LOA_DF_AGNCY_ALCTN_RCPNT_ID", "LOA_JB_ORD_NM", "LOA_SBALTMT_RCPNT_ID", "LOA_WK_CNTR_RCPNT_NM", "LOA_MAJ_RMBSMT_SRC_ID", "LOA_DTL_RMBSMT_SRC_ID", "LOA_CUST_NM", "LOA_OBJ_CLS_ID", "LOA_SRV_SRC_ID", "LOA_SPCL_INTR_ID", "LOA_BDGT_ACNT_CLS_NM", "LOA_DOC_ID", "LOA_CLS_REF_ID", "LOA_INSTL_ACNTG_ACT_ID", "LOA_LCL_INSTL_ID", "LOA_FMS_TRNSACTN_ID", "LOA_DSC_TX", "LOA_BGN_DT", "LOA_END_DT", "LOA_FNCT_PRS_NM", "LOA_STAT_CD", "LOA_HIST_STAT_CD", "LOA_HS_GDS_CD", "ORG_GRP_DFAS_CD", "LOA_UIC", "LOA_TRNSN_ID", "LOA_SUB_ACNT_ID", "LOA_BET_CD", "LOA_FND_TY_FG_CD", "LOA_BGT_LN_ITM_ID", "LOA_SCRTY_COOP_IMPL_AGNC_CD", "LOA_SCRTY_COOP_DSGNTR_CD", "LOA_SCRTY_COOP_LN_ITM_ID", "LOA_AGNC_DSBR_CD", "LOA_AGNC_ACNTNG_CD", "LOA_FND_CNTR_ID", "LOA_CST_CNTR_ID", "LOA_PRJ_ID", "LOA_ACTVTY_ID", "LOA_CST_CD", "LOA_WRK_ORD_ID", "LOA_FNCL_AR_ID", "LOA_SCRTY_COOP_CUST_CD", "LOA_END_FY_TX", "LOA_BG_FY_TX", "LOA_BGT_RSTR_CD", "LOA_BGT_SUB_ACT_CD", "ROW_STS_CD"}

const timeConversionMethod = "2006-01-02 15:04:05"

// Parse the pipe delimited .txt file with the following assumptions:
// 1. The first and last lines are the security classification.
// 2. The second line of the file are the columns that will be a 1:1 match to the LineOfAccountingTrdmFileRecord struct in pipe delimited format.
// 3. There are 58 values per line, excluding the security classification. Again, to know what these values are refer to note #2.
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

	columnNameAndLocation := make(map[string]int)

	for i := 0; i < len(columnHeaders); i++ {
		columnNameAndLocation[columnHeaders[i]] = i
	}

	// Process the lines of the .txt file into modeled codes
	codes, err := processLines(scanner, columnNameAndLocation, codes)
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
func processLines(scanner *bufio.Scanner, columnHeaders map[string]int, codes []models.LineOfAccounting) ([]models.LineOfAccounting, error) {
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
		if values[columnHeaders["LOA_SYS_ID"]] == "" {
			return nil, errors.New("malformed line in the provided loa file: " + line)
		}

		// Skip the line if it's deleted
		if values[columnHeaders["ROW_STS_CD"]] == "DLT" {
			continue
		}

		// Check if beginning date and expired date are not blank. If not blank, run time conversion, if not, leave empty.
		if values[columnHeaders["LOA_BGN_DT"]] != "" && values[columnHeaders["LOA_END_DT"]] != "" {
			parsedDate, parseError := time.Parse(timeConversionMethod, values[columnHeaders["LOA_BGN_DT"]])
			if parseError != nil {
				return nil, fmt.Errorf("malformed effective date in the provided loa file: %s", parseError)
			}
			beginningDate = parsedDate

			parsedDate, parseError = time.Parse(timeConversionMethod, values[columnHeaders["LOA_END_DT"]])
			if parseError != nil {
				return nil, fmt.Errorf("malformed effective date in the provided loa file: %s", parseError)
			}
			endingDate = parsedDate
		}

		// Check if beginning and ending fiscal years are not blank. If not blank, run str to int conversion, if not, leave empty.
		if values[columnHeaders["LOA_END_FY_TX"]] != "" && values[columnHeaders["LOA_BG_FY_TX"]] != "" {
			endingFY, err = strconv.Atoi(values[columnHeaders["LOA_END_FY_TX"]])
			if err != nil {
				return nil, fmt.Errorf("malformed ending fiscal year int in the provided loa file: %s", err)
			}

			beginningFY, err = strconv.Atoi(values[columnHeaders["LOA_BG_FY_TX"]])
			if err != nil {
				return nil, fmt.Errorf("malformed beginning fiscal year int in the provided loa file: %s", err)
			}
		}

		code := models.LineOfAccounting{
			LoaSysID:               &values[columnHeaders["LOA_SYS_ID"]],
			LoaDptID:               &values[columnHeaders["LOA_DPT_ID"]],
			LoaTnsfrDptNm:          &values[columnHeaders["LOA_TNSFR_DPT_NM"]],
			LoaBafID:               &values[columnHeaders["LOA_BAF_ID"]],
			LoaTrsySfxTx:           &values[columnHeaders["LOA_TRSY_SFX_TX"]],
			LoaMajClmNm:            &values[columnHeaders["LOA_MAJ_CLM_NM"]],
			LoaOpAgncyID:           &values[columnHeaders["LOA_OP_AGNCY_ID"]],
			LoaAlltSnID:            &values[columnHeaders["LOA_ALLT_SN_ID"]],
			LoaPgmElmntID:          &values[columnHeaders["LOA_PGM_ELMNT_ID"]],
			LoaTskBdgtSblnTx:       &values[columnHeaders["LOA_TSK_BDGT_SBLN_TX"]],
			LoaDfAgncyAlctnRcpntID: &values[columnHeaders["LOA_DF_AGNCY_ALCTN_RCPNT_ID"]],
			LoaJbOrdNm:             &values[columnHeaders["LOA_JB_ORD_NM"]],
			LoaSbaltmtRcpntID:      &values[columnHeaders["LOA_SBALTMT_RCPNT_ID"]],
			LoaWkCntrRcpntNm:       &values[columnHeaders["LOA_WK_CNTR_RCPNT_NM"]],
			LoaMajRmbsmtSrcID:      &values[columnHeaders["LOA_MAJ_RMBSMT_SRC_ID"]],
			LoaDtlRmbsmtSrcID:      &values[columnHeaders["LOA_DTL_RMBSMT_SRC_ID"]],
			LoaCustNm:              &values[columnHeaders["LOA_CUST_NM"]],
			LoaObjClsID:            &values[columnHeaders["LOA_OBJ_CLS_ID"]],
			LoaSrvSrcID:            &values[columnHeaders["LOA_SRV_SRC_ID"]],
			LoaSpclIntrID:          &values[columnHeaders["LOA_SPCL_INTR_ID"]],
			LoaBdgtAcntClsNm:       &values[columnHeaders["LOA_BDGT_ACNT_CLS_NM"]],
			LoaDocID:               &values[columnHeaders["LOA_DOC_ID"]],
			LoaClsRefID:            &values[columnHeaders["LOA_CLS_REF_ID"]],
			LoaInstlAcntgActID:     &values[columnHeaders["LOA_INSTL_ACNTG_ACT_ID"]],
			LoaLclInstlID:          &values[columnHeaders["LOA_LCL_INSTL_ID"]],
			LoaFmsTrnsactnID:       &values[columnHeaders["LOA_FMS_TRNSACTN_ID"]],
			LoaDscTx:               &values[columnHeaders["LOA_DSC_TX"]],
			LoaBgnDt:               &beginningDate,
			LoaEndDt:               &endingDate,
			LoaFnctPrsNm:           &values[columnHeaders["LOA_FNCT_PRS_NM"]],
			LoaStatCd:              &values[columnHeaders["LOA_STAT_CD"]],
			LoaHistStatCd:          &values[columnHeaders["LOA_HIST_STAT_CD"]],
			LoaHsGdsCd:             &values[columnHeaders["LOA_HS_GDS_CD"]],
			OrgGrpDfasCd:           &values[columnHeaders["ORG_GRP_DFAS_CD"]],
			LoaUic:                 &values[columnHeaders["LOA_UIC"]],
			LoaTrnsnID:             &values[columnHeaders["LOA_TRNSN_ID"]],
			LoaSubAcntID:           &values[columnHeaders["LOA_SUB_ACNT_ID"]],
			LoaBetCd:               &values[columnHeaders["LOA_BET_CD"]],
			LoaFndTyFgCd:           &values[columnHeaders["LOA_FND_TY_FG_CD"]],
			LoaBgtLnItmID:          &values[columnHeaders["LOA_BGT_LN_ITM_ID"]],
			LoaScrtyCoopImplAgncCd: &values[columnHeaders["LOA_SCRTY_COOP_IMPL_AGNC_CD"]],
			LoaScrtyCoopDsgntrCd:   &values[columnHeaders["LOA_SCRTY_COOP_DSGNTR_CD"]],
			LoaScrtyCoopLnItmID:    &values[columnHeaders["LOA_SCRTY_COOP_LN_ITM_ID"]],
			LoaAgncDsbrCd:          &values[columnHeaders["LOA_AGNC_DSBR_CD"]],
			LoaAgncAcntngCd:        &values[columnHeaders["LOA_AGNC_ACNTNG_CD"]],
			LoaFndCntrID:           &values[columnHeaders["LOA_FND_CNTR_ID"]],
			LoaCstCntrID:           &values[columnHeaders["LOA_CST_CNTR_ID"]],
			LoaPrjID:               &values[columnHeaders["LOA_PRJ_ID"]],
			LoaActvtyID:            &values[columnHeaders["LOA_ACTVTY_ID"]],
			LoaCstCd:               &values[columnHeaders["LOA_CST_CD"]],
			LoaWrkOrdID:            &values[columnHeaders["LOA_WRK_ORD_ID"]],
			LoaFnclArID:            &values[columnHeaders["LOA_FNCL_AR_ID"]],
			LoaScrtyCoopCustCd:     &values[columnHeaders["LOA_SCRTY_COOP_CUST_CD"]],
			LoaEndFyTx:             &beginningFY,
			LoaBgFyTx:              &endingFY,
			LoaBgtRstrCd:           &values[columnHeaders["LOA_BGT_RSTR_CD"]],
			LoaBgtSubActCd:         &values[columnHeaders["LOA_BGT_SUB_ACT_CD"]],
		}

		codes = append(codes, code)
	}

	return codes, scanner.Err()
}
