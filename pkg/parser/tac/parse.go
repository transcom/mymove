package tac

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

var expectedColumnNames = []string{"TAC_SYS_ID", "LOA_SYS_ID", "TRNSPRTN_ACNT_CD", "TAC_FY_TXT", "TAC_FN_BL_MOD_CD", "ORG_GRP_DFAS_CD", "TAC_MVT_DSG_ID", "TAC_TY_CD", "TAC_USE_CD", "TAC_MAJ_CLMT_ID", "TAC_BILL_ACT_TXT", "TAC_COST_CTR_NM", "BUIC", "TAC_HIST_CD", "TAC_STAT_CD", "TRNSPRTN_ACNT_TX", "TRNSPRTN_ACNT_BGN_DT", "TRNSPRTN_ACNT_END_DT", "DD_ACTVTY_ADRS_ID", "TAC_BLLD_ADD_FRST_LN_TX", "TAC_BLLD_ADD_SCND_LN_TX", "TAC_BLLD_ADD_THRD_LN_TX", "TAC_BLLD_ADD_FRTH_LN_TX", "TAC_FNCT_POC_NM"}

// Parse the pipe delimited .txt file with the following assumptions:
// 1. The first and last lines are the security classification.
// 2. The second line of the file are the columns that will be a 1:1 match to the TransportationAccountingCodeTrdmFileRecord struct in pipe delimited format.
// 3. There are 23 values per line, excluding the security classification. Again, to know what these values are refer to note #2.
// 4. All values are in pipe delimited format.
// 5. Null values will be present, but are not acceptable for TRNSPRTN_ACNT_CD.
func Parse(file io.Reader) ([]models.TransportationAccountingCode, error) {

	// Init variables
	var codes []models.TransportationAccountingCode
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
		return errors.New("column names were not parsed properly from the second line of the tac file")
	}

	if !reflect.DeepEqual(expectedColumnNames, expectedColumnNames) {
		return errors.New("column names parsed do not match the expected format of tac file records")
	}
	return nil
}

// Removes all TACs with an expiration date in the past
func PruneExpiredTACs(codes []models.TransportationAccountingCode) []models.TransportationAccountingCode {
	var pruned []models.TransportationAccountingCode

	// If the expiration date is not before time.Now(), then it is not expired and should be appended to the pruned array
	for _, code := range codes {
		if !code.TrnsprtnAcntEndDt.Before(time.Now()) {
			pruned = append(pruned, code)
		}
	}

	return pruned
}

// Consolidates TACs with the same TAC value. Duplicate "Transaction", aka description & TrnsprtnAcntTx, values are combined with a delimeter of ". Additional description found: "
func ConsolidateDuplicateTACsDesiredFromTRDM(codes []models.TransportationAccountingCode) []models.TransportationAccountingCode {
	consolidatedMap := make(map[string]models.TransportationAccountingCode)

	for _, code := range codes {
		consolidatedMap[code.TAC] = overwriteDuplicateCode(consolidatedMap[code.TAC], code)
	}

	var consolidated []models.TransportationAccountingCode
	for _, value := range consolidatedMap {
		consolidated = append(consolidated, value)
	}

	return consolidated
}

// This function checks two TAC codes: one existing and one new to decide which one to keep based on their ExpirationDates
// If the ExpirationDate is the same, it appends the transactions and maintains the first code found
// If the ExpirationDate is different, it appends the transactions and maintains the ExpirationDate further in the future
func overwriteDuplicateCode(existingCode models.TransportationAccountingCode, newCode models.TransportationAccountingCode) models.TransportationAccountingCode {

	// If the existing code has a nil expiry date but the new code does not, the new code is a better representation and should be return
	if existingCode.TrnsprtnAcntEndDt == nil && newCode.TrnsprtnAcntEndDt != nil {
		return newCode
	}

	// If the new code expires later, append its transaction to the existing one (if not empty), and keep the new code
	if newCode.TrnsprtnAcntEndDt.After(*existingCode.TrnsprtnAcntEndDt) && newCode.TrnsprtnAcntEndDt != nil {
		if existingCode.TrnsprtnAcntTx != nil && *existingCode.TrnsprtnAcntTx != "" {
			*newCode.TrnsprtnAcntTx = *existingCode.TrnsprtnAcntTx + ". Additional description found: " + *newCode.TrnsprtnAcntTx
		}
		return newCode
	}

	// If the new code expires at the same time or earlier compared to the existing code,
	// append its transaction to the existing code (if not empty) because this one expires earlier or is already expired.
	// A separate function handles the pruning of expired codes, not this one
	// Additionally, the new codes end date must not be nil
	if newCode.TrnsprtnAcntEndDt.Before(*existingCode.TrnsprtnAcntEndDt) || newCode.TrnsprtnAcntEndDt.Equal(*existingCode.TrnsprtnAcntEndDt) && newCode.TrnsprtnAcntEndDt != nil {
		if existingCode.TrnsprtnAcntTx != nil && *existingCode.TrnsprtnAcntTx != "" {
			*existingCode.TrnsprtnAcntTx = *existingCode.TrnsprtnAcntTx + ". Additional description found: " + *newCode.TrnsprtnAcntTx
		} else {
			existingCode.TrnsprtnAcntTx = newCode.TrnsprtnAcntTx
		}
	}

	return existingCode
}

// This function handles the heavy lifting for the main parse function. It handles the scanning of every line and conversion into the TransportationAccountingCode model.
func processLines(scanner *bufio.Scanner, columnHeaders []string, codes []models.TransportationAccountingCode) ([]models.TransportationAccountingCode, error) {
	// Scan every line and parse into Transportation Accounting Codes
	for scanner.Scan() {
		line := scanner.Text()
		var tacFyTxt int
		var tacSysID int
		var loaSysID int
		var err error

		// This check will skip the last line of the file.
		if line == "Unclassified" {
			break
		}

		// Gather values from the pipe delimited line
		values := strings.Split(line, "|")
		if len(values) != len(columnHeaders) {
			return nil, errors.New("malformed line in the provided tac file: " + line)
		}

		// Skip the entry if the TAC value is empty
		if values[2] == "" {
			continue
		}

		// If TacSysID is not empty, convert to int
		if values[0] != "" {
			tacSysID, err = strconv.Atoi(values[0])
			if err != nil {
				return nil, errors.New("malformed tac_sys_id in the provided tac file: " + line)
			}
		}

		// If LoaSysId is not empty, convert to int
		if values[1] != "" {
			loaSysID, err = strconv.Atoi(values[1])
			if err != nil {
				return nil, errors.New("malformed loa_sys_id in the provided tac file: " + line)
			}
		}

		// Check if fiscal year text is not empty, convert to int
		if values[3] != "" {
			tacFyTxt, err = strconv.Atoi(values[3])
			if err != nil {
				return nil, fmt.Errorf("malformed tac_fy_txt in the provided tac file: %s", err)
			}
		}

		effectiveDate, err := time.Parse("2006-01-02 15:04:05", values[16])
		if err != nil {
			return nil, fmt.Errorf("malformed effective date in the provided tac file: %s", err)
		}

		expiredDate, err := time.Parse("2006-01-02 15:04:05", values[17])
		if err != nil {
			return nil, fmt.Errorf("malformed expiration date in the provided tac file: %s", err)
		}

		code := models.TransportationAccountingCode{
			TacSysID:           &tacSysID,
			LoaSysID:           &loaSysID,
			TAC:                values[2],
			TacFyTxt:           &tacFyTxt,
			TacFnBlModCd:       &values[4],
			OrgGrpDfasCd:       &values[5],
			TacMvtDsgID:        &values[6],
			TacTyCd:            &values[7],
			TacUseCd:           &values[8],
			TacMajClmtID:       &values[9],
			TacBillActTxt:      &values[10],
			TacCostCtrNm:       &values[11],
			Buic:               &values[12],
			TacHistCd:          &values[13],
			TacStatCd:          &values[14],
			TrnsprtnAcntTx:     &values[15],
			TrnsprtnAcntBgnDt:  &effectiveDate,
			TrnsprtnAcntEndDt:  &expiredDate,
			DdActvtyAdrsID:     &values[18],
			TacBlldAddFrstLnTx: &values[19],
			TacBlldAddScndLnTx: &values[20],
			TacBlldAddThrdLnTx: &values[21],
			TacBlldAddFrthLnTx: &values[22],
			TacFnctPocNm:       &values[23],
		}

		codes = append(codes, code)
	}

	return codes, scanner.Err()
}
