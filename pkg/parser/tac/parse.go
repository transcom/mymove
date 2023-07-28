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

	// If the expiration date is not before time.Now(), then it is not expired and should be appended to the pruned array
	for _, code := range codes {
		if !code.ExpirationDate.Before(time.Now()) {
			pruned = append(pruned, code)
		}
	}

	return pruned
}

// Consolidates TACs with the same TAC value. Duplicate "Transaction", aka description, values are combined with a delimeter of ". Additional description found: "
func ConsolidateDuplicateTACsDesiredFromTRDM(codes []models.TransportationAccountingCodeDesiredFromTRDM) []models.TransportationAccountingCodeDesiredFromTRDM {
	consolidatedMap := make(map[string]models.TransportationAccountingCodeDesiredFromTRDM)

	for _, code := range codes {
		consolidatedMap[code.TAC] = overwriteDuplicateCode(consolidatedMap[code.TAC], code)
	}

	var consolidated []models.TransportationAccountingCodeDesiredFromTRDM
	for _, value := range consolidatedMap {
		consolidated = append(consolidated, value)
	}

	return consolidated
}

// This function checks two TAC codes: one existing and one new to decide which one to keep based on their ExpirationDates
// If the ExpirationDate is the same, it appends the transactions and maintains the first code found
// If the ExpirationDate is different, it appends the transactions and maintains the ExpirationDate further in the future
func overwriteDuplicateCode(existingCode models.TransportationAccountingCodeDesiredFromTRDM, newCode models.TransportationAccountingCodeDesiredFromTRDM) models.TransportationAccountingCodeDesiredFromTRDM {

	// If the new code expires later, append its transaction to the existing one (if not empty), and keep the new code
	if newCode.ExpirationDate.After(existingCode.ExpirationDate) {
		if existingCode.Transaction != "" {
			newCode.Transaction = existingCode.Transaction + ". Additional description found: " + newCode.Transaction
		}
		return newCode
	}

	// If the new code expires at the same time or earlier compared to the existing code,
	// append its transaction to the existing code (if not empty) because this one expires earlier or is already expired.
	// A separate function handles the pruning of expired codes, not this one
	if newCode.ExpirationDate.Before(existingCode.ExpirationDate) || newCode.ExpirationDate.Equal(existingCode.ExpirationDate) {
		if existingCode.Transaction != "" {
			existingCode.Transaction = existingCode.Transaction + ". Additional description found: " + newCode.Transaction
		} else {
			existingCode.Transaction = newCode.Transaction
		}
	}

	return existingCode
}

// This function handles the heavy lifting for the main parse function. It handles the scanning of every line and conversion into the TransportationAccountingCodeDesiredFromTRDM model.
func processLines(scanner *bufio.Scanner, columnHeaders []string, codes []models.TransportationAccountingCodeDesiredFromTRDM) ([]models.TransportationAccountingCodeDesiredFromTRDM, error) {
	// Scan every line and parse into Transportation Accounting Codes
	for scanner.Scan() {
		line := scanner.Text()

		// This check will skip the last line of the file.
		if line == "Unclassified" {
			break
		}

		values := strings.Split(line, "|")
		if len(values) != len(columnHeaders) {
			return nil, errors.New("malformed line in the provided tac file: " + line)
		}

		// Skip the entry if the TAC value is empty
		if values[2] == "" {
			continue
		}

		effectiveDate, err := time.Parse("2006-01-02 15:04:05", values[16])
		if err != nil {
			return nil, fmt.Errorf("malformed effective date in the provided tac file: %s", err)
		}

		expiredDate, err := time.Parse("2006-01-02 15:04:05", values[17])
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

	return codes, scanner.Err()
}
