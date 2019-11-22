package pricing

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
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
func getCell(cells []*xlsx.Cell, i int) string {
	if len(cells) > i {
		return cells[i].String()
	}

	return ""
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

func removeWhiteSpace(stripString string) string {
	space := regexp.MustCompile(`\s`)
	s := space.ReplaceAllString(stripString, "")

	return s
}

func verifyHeader(row *xlsx.Row, column int, expectedName string) error {
	actual := getCell(row.Cells, column)
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
