package pricing

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"
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

// A safe way to get a cell's value (as a string) from a sheet
func getCell(sheet *xlsx.Sheet, rowIndex, colIndex int) (string, error) {
	if sheet == nil {
		return "", fmt.Errorf("sheet is nil")
	}

	if rowIndex < 0 || rowIndex >= sheet.MaxRow || colIndex < 0 || colIndex >= sheet.MaxCol {
		return "", fmt.Errorf("cell coordinates are out of bounds")
	}

	cell, err := sheet.Cell(rowIndex, colIndex)
	if err != nil {
		return "", err
	}

	return cell.String(), nil
}

// A version of getCell that panics if it can't read the cell's value
func mustGetCell(sheet *xlsx.Sheet, rowIndex, colIndex int) string {
	cellString, err := getCell(sheet, rowIndex, colIndex)
	if err != nil {
		panic(fmt.Sprintf("getCell: sheet=\"%s\", row=%d, col=%d: %s", sheet.Name, rowIndex, colIndex, err.Error()))
	}

	return cellString
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
	actual := mustGetCell(sheet, rowIndex, colIndex)
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

	if adtlSuffix != nil && *adtlSuffix != "" {
		name = name + "_" + *adtlSuffix
	}

	name = strconv.Itoa(index) + "_" + name + "_" + runTime.Format("20060102150405") + ".csv"

	return name
}

// newDebugPrefix creates a debug-based PrefixPrinter with the specified prefix text.
func newDebugPrefix(prefixText string) *pterm.PrefixPrinter {
	return pterm.Debug.WithPrefix(pterm.Prefix{
		Text:  prefixText,
		Style: &pterm.ThemeDefault.DebugPrefixStyle,
	})
}
