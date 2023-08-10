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
			return nil, errors.New("file column headers do not match")
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
	expectedColumnNames := getFieldNames(models.LineOfAccountingTrdmFileRecord{})
	if !reflect.DeepEqual(columnNames, expectedColumnNames) {
		return errors.New("column names parsed do not match the expected format of loa file records")
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

		loaSysId, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, errors.New("malformed line in the provided loa file: " + line)
		}

		// Check if beginning date and expired date are not blank. If not blank, run time conversion, if not, leave empty.
		if values[27] != "" && values[28] != "" {
			// Parse values[27], this is the beginning date
			parsedDate, err := time.Parse(timeConversionMethod, values[27])
			if err != nil {
				return nil, fmt.Errorf("malformed effective date in the provided loa file: %s", err)
			}
			beginningDate = parsedDate

			// Parse values[28], this is the ending date
			parsedDate, err = time.Parse(timeConversionMethod, values[28])
			if err != nil {
				return nil, fmt.Errorf("malformed effective date in the provided loa file: %s", err)
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
			LoaSysID:               &loaSysId,
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
