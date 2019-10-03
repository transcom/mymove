package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

/*************************************************************************

1) 1b) Service Areas

2) Domestic Price Tabs
        2a) Domestic Linehaul Prices
	    2b) Domestic Service Area Prices
	    2c) Other Domestic Prices

3) International Price Tabs
        3a) OCONUS to OCONUS Prices
	    3b) CONUS to OCONUS Prices
	    3c) OCONUS to CONUS Prices
	    3d) Other International Prices
	    3e) Non-Standard Loc'n Prices

4) Mgmt., Coun., Trans. Prices Tab
        4a) Mgmt., Coun., Trans. Prices

5) Other Prices Tabs
        5a) Access. and Add. Prices
	    5b) Price Escalation Discount


--------------------------------------------------------------------------

0: 	Guide to Pricing Rate Table
1: 	Total Evaluated Price
2: 	Submission Checklist
3: 	1a) Directions
4: 	1b) Service Areas
5: 	Domestic Price Tabs >>>
6: 	2a) Domestic Linehaul Prices
7: 	2b) Dom. Service Area Prices
8: 	2c) Other Domestic Prices
9: 	International Prices Tables >>>
10: 3a) OCONUS to OCONUS Prices
11: 3b) CONUS to OCONUS Prices
12: 3c) OCONUS to CONUS Prices
13: 3d) Other International Prices
14: 3e) Non-Standard Loc'n Prices
15:	Other Price Tables
16: 4a) Mgmt., Coun., Trans. Prices
17: 5a) Access. and Add. Prices
18: 5b) Price Escalation Discount
19: Domestic Linehaul Data
20: Domestic Move Count
21: Domestic Avg Weight
22: Domestic Avg Milage
23: Domestic Price Calculation >>>
24: Domestic Linehaul Calculation
25: Domestic SA Price Calculation
26: NTS Packing Calculation
27: Int'l Price Calculation >>>
28: OCONUS to OCONUS Calculation
29: CONUS to OCONUS Calculation
30: OCONUS to CONUS Calculation
31: Other Int'l Prices Calculation
32: Non-Standard Loc'n Calculation
33: Other Calculations >>>
34: Mgmt., Coun., Trans., Calc
35: Access. and Add. Calculation


 *************************************************************************/
const xlsxSheetsCountMax int = 35
type processXlsxSheet func(config, int) error
type verifyXlsxSheet func(config, int) error

type xlsxDataSheetInfo struct {
	description *string
	process *processXlsxSheet
	verify *verifyXlsxSheet
	outputFilename *string //<id>_<outputFilename>_<time.Now().Format("20060102150405")>.csv
}

func (x *xlsxDataSheetInfo) generateOutputFilename(index int, runTime time.Time) string {
	var name string
	if x.outputFilename != nil {
		name = *x.outputFilename
	}

	name = string(index) + "_" + name + "_" + runTime.Format("20060102150405") + ".csv"

	return name
}

var xlsxDataSheets []xlsxDataSheetInfo

func initDataSheetInfo()  {
	xlsxDataSheets := make([]xlsxDataSheetInfo, xlsxSheetsCountMax, xlsxSheetsCountMax)

	// 6: 	2a) Domestic Linehaul Prices
	xlsxDataSheets[6] = xlsxDataSheetInfo {
		description: stringPointer("2a) Domestic Linehaul Prices"),
		outputFilename: stringPointer("2b_domestic_linehaul_prices"),
		process: &parseDomesticLinehaulPrices,
		verify: &verifyDomesticLinehaulPrices,
	}
}

type config struct {
	displayHelp bool
	processAll bool
	showOutput bool
	xlsxFilename *string
	xlsxSheet *int
	saveToFile bool
	runTime time.Time
}

func usage() {
	fmt.Printf(`%s: -f=<CSV Input File> -o=<XLSX Output File> -d=<Delimiter>
`,
		os.Args[0])
}


func main() {

	initDataSheetInfo()
	config := config{}
	config.runTime = time.Now()

	help := flag.Bool("help", false, "Display help/usage info")
	all := flag.Bool("all", false, "true, if parsing entire Rate Engine GHC XLSX")
	display := flag.Bool("display", false, "true, if displaying output of parsed info")
	filename := flag.String("filename", "", "Filename including path of the XLSX to parse for Rate Engine GHC import")
	saveToFile := flag.Bool("save", false, "true, if saving output to file")
	//TODO change xlsxSheet to a string and name to xlsxSheets
	config.xlsxSheet = flag.Int("xlsxSheet", 99, "Sequential sheet index number starting with 0")

	/**
	TODO: - implement xlsxSheets (string of indexes to process)
	TODO: - implement help!!!
	TODO: - implement print out availalble indices for processing xlsxDataSheets
	TODO: - implement verification perSheet of expected filled in cells (add verify function xlsxDataSheetInfo
	 */

	flag.Parse()

	config.displayHelp = false
	if help != nil && *help == true {
		config.displayHelp = true
		usage()
		return
	}

	config.processAll = false
	if all != nil && *all == true {
		config.processAll = true
		// TODO parse everything
	}

	config.xlsxFilename = filename
	if filename != nil {
		log.Printf("Importing file %s\n", *filename)
	}

	config.showOutput = false
	if display != nil && *display == true {
		config.showOutput = true
	}

	config.saveToFile = false
	if saveToFile != nil && *saveToFile == true {
		config.saveToFile = true
	}

	// TODO process config.xlsxSheets

	// TODO properly call the process function
	err := process(config, 6)
	if err != nil {
		log.Fatalf("Error processing %v\n", err)
	}

}

func process(config config, sheetIndex int) error {
	xlsxInfo := xlsxDataSheets[sheetIndex]
	var description string
	if xlsxInfo.description != nil {
		description = *xlsxInfo.description
		log.Printf("Processing sheet index %d with description %s\n", sheetIndex, description)
	} else {
		log.Printf("Processing sheet index %d with missing description\n", sheetIndex)
	}

	// Call verify function
	if xlsxInfo.verify != nil {
		var callFunc verifyXlsxSheet
		callFunc = *xlsxInfo.verify
		err := callFunc(config, sheetIndex)
		if err != nil {
			log.Printf("%s verify error: %v\n", description, err)
		}
	} else {
		log.Printf("No verify function for sheet index %d with description %s\n", sheetIndex, description)
	}

	// Call process function
	if xlsxInfo.process != nil {
		var callFunc processXlsxSheet
		callFunc = *xlsxInfo.process
		err := callFunc(config, sheetIndex)
		if err != nil {
			log.Printf("%s process error: %v\n", description, err)
		}
	} else {
		log.Fatalf("Missing process function for sheet index %d with description %s\n", sheetIndex, description)
	}

	// Verification and Process completed
	log.Printf("Completed processing sheet index %d with description %s\n", sheetIndex, description)
	return nil
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
			fmt.Printf("WARNING: getInt() invalid int syntax checking string <%s> for float string\n", from)
			f, err := strconv.ParseFloat(from, 32)
			if err != nil {
				fmt.Printf("ERROR: getInt() ParseFloat error %s\n", err.Error())
				return 0
			}
			if f != 0.0 {
				fmt.Printf("SUCCESS: getInt() converting string <%s> from float to int <%d>\n", from, int(f))
				return int(f)
			}
		}
		fmt.Printf("ERROR: getInt() Atoi error %s\n", err.Error())
		return 0
	}

	return i
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

var verifyDomesticLinehaulPrices verifyXlsxSheet = func (params config, sheetIndex int) error {
	log.Println("TODO verifyDomesticLinehaulPrices() not implemented")
	return nil
}

var parseDomesticLinehaulPrices processXlsxSheet = func (params config, sheetIndex int) error {

	/*
		peak and non-peak
		weightBands
		milage bands
		services area -> origin service -> service schedule
		base period year

		available functions:
			ColIndexToLetters
			ColLettersToIndex
	*/

	if weightBandNumCells != weightBandNumCellsExpected {
		return fmt.Errorf("parseDomesticLinehaulPrices(): Exepected %d columns per weight band, found %d defined in golang parser\n", weightBandNumCellsExpected, weightBandNumCells)
	}

	if len(weightBands) != weightBandCountExpected {
		return fmt.Errorf("parseDomesticLinehaulPrices(): Exepected %d weight bands, found %d defined in golang parser\n", weightBandCountExpected, len(weightBands))
	}

	// TODO can this be a function? I think yes
	// TODO return point, write to file if point is not nil
	var createCsv createCsvHelper
	if params.saveToFile == true {
		err := createCsv.createCsvWriter(xlsxDataSheets[sheetIndex].generateOutputFilename(sheetIndex, params.runTime))
		checkError("Failed to create CSV writer", err)
		defer createCsv.close()
	}

	if params.xlsxFilename == nil {
		return fmt.Errorf("parseDomesticLinehaulPrices(): did not receive an XLSX filename to parse")
	}
	xlFile, err := xlsx.OpenFile(*params.xlsxFilename)
	if err != nil {
		return err
	}

	const xlsxDataSheetNum int = 6 // 2a) Domestic Linehaul Prices
	const feeColIndexStart int = 6 // start at column 6 to get the rates
	colIndex := feeColIndexStart
	const feeRowIndexStart int = 14 // start at row 14 to get the rates
	const serviceAreaNumberColumn int = 2
	const originServiceAreaColumn int = 3
	const serviceScheduleColumn int = 4

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("parseDomesticLinehaulPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	dataRows := xlFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart:]
	for _, row := range dataRows {
		// For number of baseline + escalation years
		colIndex = feeColIndexStart
		numEscalationYears := 1
		for escalation := 0; escalation < numEscalationYears; escalation++ {
			// For each rate season
			for _, r := range rateTypes {
				// For each weight band
				for _, w := range weightBands {
					// For each milage range
					for _, m := range milesRanges {
						domesticLineHaulPrice := domesticLineHaulPrice{
							serviceAreaNumber:     getInt(getCell(row.Cells, serviceAreaNumberColumn)),
							originServiceArea:     getCell(row.Cells, originServiceAreaColumn),
							serviceSchedule:       getInt(getCell(row.Cells, serviceScheduleColumn)),
							season:                r,
							weightBand:            w,
							milesRange:            m,
							optionPeriodYearCount: escalation,
							rate:                  getCell(row.Cells, colIndex),
						}
						colIndex++
						if params.showOutput == true {
							log.Println(domesticLineHaulPrice.toSlice())
						}
						if params.saveToFile == true {
							createCsv.write(domesticLineHaulPrice.toSlice())
						}
					}
					//TODO DEBUG REMOVE return
					return nil
				}
				colIndex++ // skip 1 column (empty column) before starting next rate type
			}
		}
	}

	return nil
}
