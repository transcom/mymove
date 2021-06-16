package pricing

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
)

/*************************************************************************************************************/
// COMMON Types
/*************************************************************************************************************/

const sharedNumEscalationYearsToProcess int = 1

var rateSeasons = []string{"NonPeak", "Peak"}

type headerInfo struct {
	headerName string
	column     int
}

/*************************************************************************/
// Shared Helper functions
/*************************************************************************/

// A safe way to get a cell from a slice of cells, returning empty string if not found
func getCell(sheet *xlsx.Sheet, rowIndex, colIndex int) string {
	if rowIndex >= sheet.MaxRow || colIndex >= sheet.MaxCol {
		return ""
	}

	cell, err := sheet.Cell(rowIndex, colIndex)
	if err != nil {
		// TODO: Is this panic OK? For now, just trying to avoid having to add in error-checking on every
		//   getCell call.  It's a CLI, so what would we do to react to it other than end the process?
		panic(err)
	}
	return cell.String()
}

func getInt(from string) (int, error) {
	i, err := strconv.Atoi(from)
	if err != nil {
		if strings.HasSuffix(err.Error(), ": invalid syntax") {
			f, ferr := strconv.ParseFloat(from, 32)
			if ferr != nil {
				return 0, ferr
			}
			if f != 0.0 {
				return int(f), nil
			}
		}

		return 0, err
	}

	return i, nil
}

func removeFirstDollarSign(s string) string {
	return strings.Replace(s, "$", "", 1)
}

func removeWhiteSpace(stripString string) string {
	space := regexp.MustCompile(`\s`)
	s := space.ReplaceAllString(stripString, "")

	return s
}

func verifyHeader(sheet *xlsx.Sheet, rowIndex, colIndex int, expectedName string) error {
	actual := getCell(sheet, rowIndex, colIndex)
	if removeWhiteSpace(expectedName) != removeWhiteSpace(actual) {
		return fmt.Errorf("format error: Header <%s> is missing; got <%s> instead", expectedName, actual)
	}

	return nil
}

// generateOutputFilename: generates filename using XlsxDataSheetInfo.outputFilename
// with the following format -- <id>_<OutputFilename>_<time.Now().Format("20060102150405")>.csv
// if the adtlSuffix is passed the format is -- <id>_<outputFilename>_<adtlSuffix>_<time.Now().Format("20060102150405")>.csv
func (x *XlsxDataSheetInfo) generateOutputFilename(index int, runTime time.Time, adtlSuffix *string) string {
	var name string
	if x.outputFilename != nil {
		name = *x.outputFilename
	} else {
		name = "rate_engine_ghc_parse"
	}

	if adtlSuffix != nil {
		name = name + "_" + *adtlSuffix
	}

	name = strconv.Itoa(index) + "_" + name + "_" + runTime.Format("20060102150405") + ".csv"

	return name
}
