package transittime

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
)

/*************************************************************************/
// Shared Helper functions
/*************************************************************************/

// A way to get a cell value from a sheet
// return empty string if not found
func getValueFromSheet(sheet *xlsx.Sheet, row int, col int) (string, error) {
	if sheet == nil {
		return "", fmt.Errorf("sheet is nil")
	}

	if row < 0 || row >= sheet.MaxRow || col < 0 || col >= sheet.MaxCol {
		return "", fmt.Errorf("cell coordinates are out of bounds")
	}

	cell, err := sheet.Cell(row, col)
	if err != nil {
		return "", err
	}

	return cell.String(), nil
}

// A version of getValueFromSheet that panics if it can't read the cell's value
func mustGetValueFromSheet(sheet *xlsx.Sheet, row, col int) string {
	cellString, err := getValueFromSheet(sheet, row, col)
	if err != nil {
		panic(fmt.Sprintf("getValueFromSheet: sheet=\"%s\", row=%d, col=%d: %s", sheet.Name, row, col, err.Error()))
	}

	return cellString
}

// A way to parse domestic header bounds.
// Ex: 0 - 1000 lbs
func getDomesticHeaderBounds(bounds string) ([]string, error) {
	trimmedStr := strings.TrimSpace(bounds)

	var slice []string
	if strings.Contains(trimmedStr, "-") {
		slice = strings.Split(trimmedStr, "-")
	} else {
		// probably >=
		slice = strings.Split(trimmedStr, ">=")
	}

	// header format should be like "0 - 100"
	if len(slice) != 2 {
		return nil, fmt.Errorf("Could not parse lower and upper bounds. Should be of format: %s", "1 - 1000 or >= 8000")
	}

	if strings.Contains(trimmedStr, "-") {
		slice[0] = strings.TrimSpace(slice[0])
		slice[1] = strings.TrimSpace(slice[1])
	} else {
		// flip if >=
		slice = []string{strings.TrimSpace(slice[1]), ""}
	}

	return slice, nil
}

func getInt(from string) int {
	i, err := strconv.Atoi(from)
	if err != nil {
		if strings.HasSuffix(err.Error(), ": invalid syntax") {
			f, ferr := strconv.ParseFloat(from, 32)
			if ferr != nil {
				return 0
			}
			if f != 0.0 {
				return int(f)
			}
		}
		log.Fatalf("ERROR: getInt() Atoi & ParseFloat failed to convert <%s> error %s, returning 0\n", from, err.Error())
	}

	return i
}

func removeFirstDollarSign(s string) string {
	return strings.Replace(s, "$", "", 1)
}

// generateOutputFilename: generates filename using XlsxDataSheetInfo.outputFilename
// with the following format -- <id>_<OutputFilename>_<time.Now().Format("20060102150405")>.csv
// if the adtlSuffix is passed the format is -- <id>_<outputFilename>_<adtlSuffix>_<time.Now().Format("20060102150405")>.csv
func (x *XlsxDataSheetInfo) generateOutputFilename(index int, runTime time.Time, adtlSuffix *string) string {
	var name string
	if x.outputFilename != nil {
		name = *x.outputFilename
	} else {
		name = "transit_time_ghc_parse"
	}

	if adtlSuffix != nil {
		name = name + "_" + *adtlSuffix
	}

	name = strconv.Itoa(index) + "_" + name + "_" + runTime.Format("20060102150405") + ".csv"

	return name
}
