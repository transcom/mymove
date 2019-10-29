package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

/*************************************************************************************************************/
// COMMON Types
/*************************************************************************************************************/

var rateSeasons = []string{"NonPeak", "Peak"}

type createCsvHelper struct {
	csvFilename string
	csvFile     *os.File
	csvWriter   *csv.Writer
}

func (cCH *createCsvHelper) createCsvWriter(filename string) error {

	cCH.csvFilename = filename
	file, err := os.Create(cCH.csvFilename)

	if err != nil {
		log.Fatal(err.Error())
	}
	cCH.csvFile = file
	cCH.csvWriter = csv.NewWriter(cCH.csvFile)

	return nil
}

func (cCH *createCsvHelper) write(record []string) {
	if cCH.csvWriter == nil {
		log.Fatalln("createCsvHelper.createCsvWriter() was not called to initialize cCH.csvWriter")
	}
	err := cCH.csvWriter.Write(record)
	if err != nil {
		log.Fatal(err.Error())
	}
	cCH.csvWriter.Flush()
}

func (cCH *createCsvHelper) close() {
	cCH.csvFile.Close()
	cCH.csvWriter.Flush()
}

/*************************************************************************/
// Shared Helper functions
/*************************************************************************/

func createCsvWriter(create bool, sheetIndex int, runTime time.Time) *createCsvHelper {
	var createCsv createCsvHelper

	if create == true {
		err := createCsv.createCsvWriter(xlsxDataSheets[sheetIndex].generateOutputFilename(sheetIndex, runTime))
		checkError("Failed to create CSV writer", err)
	} else {
		return nil
	}
	return &createCsv
}

// A safe way to get a cell from a slice of cells, returning empty string if not found
func getCell(cells []*xlsx.Cell, i int) string {
	if len(cells) > i {
		return cells[i].String()
	}

	return ""
}

// Gotta have a stringPointer function. Returns nil if empty string
func stringPointer(s string) *string {
	if s == "" {
		return nil
	}

	return &s
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

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func removeFirstDollarSign(s string) string {
	return strings.Replace(s, "$", "", 1)
}

func removeWhiteSpace(stripString string) string {
	space := regexp.MustCompile(`\s`)
	s := space.ReplaceAllString(stripString, "")

	return s
}

func formatServiceAreaNumber(s string) string {
	return fmt.Sprintf("%03d", getInt(s))
}

// generateOutputFilename: generates filename using xlsxDataSheetInfo.outputFilename
// with the folling fomat -- <id>_<outputFilename>_<time.Now().Format("20060102150405")>.csv
func (x *xlsxDataSheetInfo) generateOutputFilename(index int, runTime time.Time) string {
	var name string
	if x.outputFilename != nil {
		name = *x.outputFilename
	} else {
		name = "rate_engine_ghc_parse"
	}

	name = strconv.Itoa(index) + "_" + name + "_" + runTime.Format("20060102150405") + ".csv"

	return name
}
