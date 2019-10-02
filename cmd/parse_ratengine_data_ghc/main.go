package main

import (
	"flag"
	"fmt"
	"strconv"
)
import "go.uber.org/zap"
import "github.com/tealeg/xlsx"


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
19: Domestic  Linehaul Data
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




func help() {

}



func main() {
	logger, err := zap.NewDevelopment()

	/*
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	test := flag.Bool("test", false, "Whether to generate testy mcTest emails")
	*/
	help := flag.Bool("help", false, "Display help/usage info")
	all := flag.Bool("all", false, "True, if parsing entire Rate Engine GHC XLSX")
	display := flag.Bool("display", false, "True, if display output of parsed info")
	filename := flag.String("filename", "", "Filename including path of the XLSX to parse for Rate Engine GHC import")

	flag.Parse()

	//fmt.Printf("File written to %s\n", path)
	fmt.Printf("Importing file %s\n", filename)

	if filename != nil {
		parseDomesticLinehaulPrices(*filename)
	}

}

// A safe way to get a cell from a slice of cells, returning empty string if not found
func getCell(cells []*xlsx.Cell, i int) string {
	if len(cells) > i {
		return cells[i].String()
	}

	return ""
}

func getInt(from string) int  {
	if from == "" {
		return 0
	}
	i, err := strconv.Atoi(from)
	if err != nil {
		return 0
	}
	return i
}



func parseDomesticLinehaulPrices(parseFile string) error {

	/*
	weightBands
	peak and non-peak
	milage bands
	services area -> origin service -> service schedule
	base period year

	available functions:
		ColIndexToLetters
		ColLettersToIndex
	*/

	rateTypes := []string{"NonPeak", "Peak"}

	weightBandNumCellsExpected := 10 //cells per band verify against weightBandNumCells
	weightBandCountExpected := 3 //expected number of weight bands verify against weightBandCount

	type weightBand struct {
		band int
		lowerLbs int
		upperLbs int
		lowerCwt float32
		upperCwt float32
	}

	weightBands := []weightBand{
		{
			band: 1,
			lowerLbs: 500,
			upperLbs: 4999,
			lowerCwt: 5,
			upperCwt: 49.99,
		},
		{
			band: 2,
			lowerLbs: 5000,
			upperLbs: 9999,
			lowerCwt: 50,
			upperCwt: 99.99,
		},
		{
			band: 3,
			lowerLbs: 10000,
			upperLbs: 999999,
			lowerCwt: 100,
			upperCwt: 9999.99,
		},

	}
	weightBandCount := len(weightBands) //number of bands and then repeats

	type milesRange struct {
		rangeNumber int
		lower int
		upper int
	}

	milesRanges := []milesRange {
		{
			rangeNumber: 1,
			lower: 0,
			upper: 250,
		},
		{
			rangeNumber: 2,
			lower: 251,
			upper: 500,
		},
		{
			rangeNumber: 3,
			lower: 501,
			upper: 1000,
		},
		{
			rangeNumber: 4,
			lower: 1001,
			upper: 1500,
		},
		{
			rangeNumber: 5,
			lower: 1501,
			upper: 2000,
		},
		{
			rangeNumber: 6,
			lower: 2001,
			upper: 2500,
		},
		{
			rangeNumber: 7,
			lower: 2501,
			upper: 3000,
		},
		{
			rangeNumber: 8,
			lower: 3001,
			upper: 3500,
		},
		{
			rangeNumber: 9,
			lower: 3501,
			upper: 4000,
		},
		{
			rangeNumber: 10,
			lower: 4001,
			upper: 999999,
		},
	}
	weightBandNumCells := len(milesRanges) //


	type domesticLineHaulPrice struct {
		serviceAreaNumber int
		originServiceArea string
		serviceSchedule int
		season string
		weightBand weightBand
		milesRange milesRange
		optionPeriodYearCount int //the escalation type
		rate string
	}


	if weightBandNumCells != weightBandNumCellsExpected {
		fmt.Errorf("Exepected %d columns per weight band, found %d defined in golang parser\n", weightBandNumCellsExpected, weightBandNumCells)
	}

	if weightBandCount != weightBandCountExpected {
		fmt.Errorf("Exepected %d weight bands, found %d defined in golang parser\n", weightBandCountExpected, weightBandCount)
	}


	xlFile, err := xlsx.OpenFile(parseFile)
	if err != nil {
		return err
	}

	feeColIndexStart := 6 // start at column 6 to get the rates
	colIndex := feeColIndexStart

	dataRows := xlFile.Sheets[6].Rows[14:]
	for _, row := range dataRows {
		// For number of baseline + escalation years
		colIndex = feeColIndexStart
		numEscalationYears := 5
		for escalation := 0; escalation < numEscalationYears; escalation++ {
			// For each rate season
			for _, r := range rateTypes {
				// For each weight band
				for _, w := range weightBands {
					// For each milage range
					for _, m := range milesRanges {
						domesticLineHaulPrice := domesticLineHaulPrice{
							serviceAreaNumber: getInt(getCell(row.Cells, 2)),
							originServiceArea: getCell(row.Cells, 3),
							serviceSchedule: getInt(getCell(row.Cells, 4)),
							season: r,
							weightBand: w,
							milesRange: m,
							optionPeriodYearCount: escalation,
							rate: getCell(row.Cells, colIndex),
						}
						colIndex++
						fmt.Printf("%v ", domesticLineHaulPrice)
					}
				}
				colIndex++ // skip 1 column (empty column) before starting next rate type

			}
		}
	}

	//

	return nil
}


