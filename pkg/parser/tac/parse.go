package tac

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/models"
)

// Parse the pipe delimited .txt file with the following assumptions:
// 1. The first and last lines are the security classification.
// 2. The second line of the file are the columns that will be a 1:1 match to the TransportationAccountingCodeTrdmFileRecord struct in pipe delimited format.
// 3. There are 23 values per line, excluding the security classification. Again, to know what these values are refer to note #2.
// 4. All values are in pipe delimited format.
// 5. Null values will be present, but are not acceptable for TRNSPRTN_ACNT_CD.
func Parse(file io.Reader) ([]models.TransportationAccountingCodeDesiredFromTRDM, error) {

	// Init variables
	var codes []models.TransportationAccountingCodeDesiredFromTRDM
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
		ensureFileStructMatchesColumnNames(columnHeaders)
	}

	// Scan every line and parse into Transportation Accounting Codes
	for scanner.Scan() {
		line := scanner.Text()

		// This check will skip the last line of the file.
		// bufio does not appear to have a method of skipping the last line
		// nor check the total amount of lines present in a file without
		// scanning every single line to memory first.
		if line == "Unclassified" {
			break
		}

		values := strings.Split(line, "|")
		if len(values) != len(columnHeaders) {
			return nil, errors.New("malformed line in the provided tac file: " + line)
		}

		effectiveDate, err := time.Parse(time.RFC3339, values[16])
		if err != nil {
			return nil, fmt.Errorf("malformed effective date in the provided tac file: %s", err)
		}

		expiredDate, err := time.Parse(time.RFC3339, values[17])
		if err != nil {
			return nil, fmt.Errorf("malformed expiration date in the provided tac file: %s", err)
		}

		code := models.TransportationAccountingCodeDesiredFromTRDM{
			TAC:                      values[2],
			BillingAddressFirstLine:  values[19],
			BillingAddressSecondLine: values[20],
			BillingAddressThirdLine:  values[21],
			BillingAddressFourthLine: values[22],
			Transaction:              values[15],
			EffectiveDate:            effectiveDate,
			ExpirationDate:           expiredDate,
			FiscalYear:               values[3],
		}

		codes = append(codes, code)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return codes, nil
}

// Compare a struct's field names to the columns retrieved from the .txt file
func ensureFileStructMatchesColumnNames(columnNames []string) error {
	if len(columnNames) == 0 {
		return errors.New("column names were not parsed properly from the second line of the tac file")
	}
	expectedColumnNames := getFieldNames(models.TransportationAccountingCodeTrdmFileRecord{})
	if !reflect.DeepEqual(columnNames, expectedColumnNames) {
		return errors.New("column names parsed do not match the expected format of tac file records")
	}
	return nil
}

// This function gathers the struct field names for comparison to
// line 2 of the .txt file - The columns
func getFieldNames(obj interface{}) []string {
	var fieldNames []string

	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldNames = append(fieldNames, field.Name)
	}

	return fieldNames
}

// Removes all TACs with an expiration date in the past
func PruneExpiredTACsDesiredFromTRDM(codes []models.TransportationAccountingCodeDesiredFromTRDM) []models.TransportationAccountingCodeDesiredFromTRDM {
	var pruned []models.TransportationAccountingCodeDesiredFromTRDM

	for _, code := range codes {
		if code.ExpirationDate.Before(time.Now()) {
			pruned = append(pruned, code)
		}
	}

	return pruned
}

// Consoliddates TACs with the same TAC value. Duplicate "Transaction", aka description, calues are combined with a delimeter of ". Additional description found: "
func ConsolidateDuplicateTACsDesiredFromTRDM(codes []models.TransportationAccountingCodeDesiredFromTRDM) []models.TransportationAccountingCodeDesiredFromTRDM {
	consolidatedMap := make(map[string]models.TransportationAccountingCodeDesiredFromTRDM)

	for _, code := range codes {
		existingCode, exists := consolidatedMap[code.TAC]
		if exists && existingCode.Transaction != code.Transaction {
			existingCode.Transaction = existingCode.Transaction + ". Additional description found: " + code.Transaction
		} else {
			consolidatedMap[code.TAC] = code
		}
	}

	var consolidated []models.TransportationAccountingCodeDesiredFromTRDM
	for _, value := range consolidatedMap {
		consolidated = append(consolidated, value)
	}

	return consolidated
}
