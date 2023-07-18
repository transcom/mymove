package tac

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strings"

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
			return nil, errors.New("malformed line in the tac file: " + line)
		}

		code := models.TransportationAccountingCodeDesiredFromTRDM{
			TAC:                      values[2],
			BillingAddressFirstLine:  values[19],
			BillingAddressSecondLine: values[20],
			BillingAddressThirdLine:  values[21],
			BillingAddressFourthLine: values[22],
			Transaction:              values[15],
			EffectiveDate:            values[16],
			ExpirationDate:           values[17],
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
