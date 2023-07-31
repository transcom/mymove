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
func Parse(file io.Reader) ([]models.LineOfAccountingDesiredFromTRDM, error) {

	// Init variables
	var codes []models.LineOfAccountingDesiredFromTRDM
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
func PruneEmptyHhgCodes(codes []models.LineOfAccountingDesiredFromTRDM) []models.LineOfAccountingDesiredFromTRDM {
	var pruned []models.LineOfAccountingDesiredFromTRDM

	// If the household goods code is not empty, then it should be appended to the pruned array for return
	for _, code := range codes {
		if code.HouseholdGoodsCode != "" {
			pruned = append(pruned, code)
		}
	}

	return pruned
}

// This function handles the heavy lifting for the main parse function. It handles the scanning of every line and conversion into the LineOfAccountingDesiredFromTRDM model.
func processLines(scanner *bufio.Scanner, columnHeaders []string, codes []models.LineOfAccountingDesiredFromTRDM) ([]models.LineOfAccountingDesiredFromTRDM, error) {
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

		code := models.LineOfAccountingDesiredFromTRDM{
			LOA:                                values[0],
			DepartmentID:                       values[1],
			TransferDepartmentName:             values[2],
			BasicAppropriationFundID:           values[3],
			TreasurySuffixText:                 values[4],
			MajorClaimantName:                  values[5],
			OperatingAgencyID:                  values[6],
			AllotmentSerialNumberID:            values[7],
			ProgramElementID:                   values[8],
			TaskBudgetSublineText:              values[9],
			DefenseAgencyAllocationRecipientID: values[10],
			JobOrderName:                       values[11],
			SubAllotmentRecipientId:            values[12],
			WorkCenterRecipientName:            values[13],
			MajorReimbursementSourceID:         values[14],
			DetailReimbursementSourceID:        values[15],
			CustomerName:                       values[16],
			ObjectClassID:                      values[17],
			ServiceSourceID:                    values[18],
			SpecialInterestID:                  values[19],
			BudgetAccountClassificationName:    values[20],
			DocumentID:                         values[21],
			ClassReferenceID:                   values[22],
			InstallationAccountingActivityID:   values[23],
			LocalInstallationID:                values[24],
			FMSTransactionID:                   values[25],
			DescriptionText:                    values[26],
			BeginningDate:                      beginningDate,
			EndDate:                            endingDate,
			FunctionalPersonName:               values[29],
			StatusCode:                         values[30],
			HistoryStatusCode:                  values[31],
			HouseholdGoodsCode:                 values[32],
			OrganizationGroupDefenseFinanceAccountingServiceCode: values[33],
			UnitIdentificationCode:                               values[34],
			TransactionID:                                        values[35],
			SubordinateAccountID:                                 values[36],
			BusinessEventTypeCode:                                values[37],
			FundTypeFlagCode:                                     values[38],
			BudgetLineItemID:                                     values[39],
			SecurityCooperationImplementingAgencyCode:            values[40],
			SecurityCooperationDesignatorID:                      values[41],
			SecurityCooperationLineItemID:                        values[42],
			AgencyDisbursingCode:                                 values[43],
			AgencyAccountingCode:                                 values[44],
			FundCenterID:                                         values[45],
			CostCenterID:                                         values[46],
			ProjectTaskID:                                        values[47],
			ActivityID:                                           values[48],
			CostCode:                                             values[49],
			WorkOrderID:                                          values[50],
			FunctionalAreaID:                                     values[51],
			SecurityCooperationCustomerCode:                      values[52],
			EndingFiscalYear:                                     endingFY,
			BeginningFiscalYear:                                  beginningFY,
			BudgetRestrictionCode:                                values[55],
			BudgetSubActivityCode:                                values[56],
		}

		codes = append(codes, code)
	}

	return codes, scanner.Err()
}
