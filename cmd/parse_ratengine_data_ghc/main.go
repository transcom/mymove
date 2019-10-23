package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/gobuffalo/pop"
	"github.com/tealeg/xlsx"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

/*************************************************************************

Parser tool to extract data from the GHC Rate Engine XLSX

For help run: <program> -h

`go run cmd/parse_ratengine_data_ghc/*.go -h`

Rate Engine XLSX sections this tool will be parsing:

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

Rate Engine XLSX sheet tabs:

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

/*************************************************************************

To add new parser functions to this file:

	a.) (optional) Add new verify function for your processing must match signature verifyXlsxSheet
	b.) Add new process function to process XLSX data sheet must match signature processXlsxSheet
	c.) Update initDataSheetInfo() with a.) and b.)
		The index must match the sheet index in the XLSX that you aim to process

You should not have to update the main() or  process() functions. Unless you
intentionally are modifying the pattern of how the processing functions are called.

 *************************************************************************/

const xlsxSheetsCountMax int = 35

type processXlsxSheet func(paramConfig, int, *pop.Connection) error
type verifyXlsxSheet func(paramConfig, int) error

type xlsxDataSheetInfo struct {
	description    *string
	process        *processXlsxSheet
	verify         *verifyXlsxSheet
	outputFilename *string //do not include suffix
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

var xlsxDataSheets []xlsxDataSheetInfo

// initDataSheetInfo: When adding new functions for parsing sheets, must add new xlsxDataSheetInfo
// defining the parse function
//
// The index MUST match the sheet that is being processed. Refer to file comments or XLSX to
// determine the correct index to add.
func initDataSheetInfo() {
	xlsxDataSheets = make([]xlsxDataSheetInfo, xlsxSheetsCountMax, xlsxSheetsCountMax)

	// 4: 	1b) Domestic & International Service Areas
	xlsxDataSheets[4] = xlsxDataSheetInfo{
		description:    stringPointer("1b) Service Areas"),
		outputFilename: stringPointer("1b_service_areas"),
		process:        &parseServiceAreas,
		verify:         &verifyServiceAreas,
	}

	// 6: 	2a) Domestic Linehaul Prices
	xlsxDataSheets[6] = xlsxDataSheetInfo{
		description:    stringPointer("2a) Domestic Linehaul Prices"),
		outputFilename: stringPointer("2a_domestic_linehaul_prices"),
		process:        &parseDomesticLinehaulPrices,
		verify:         &verifyDomesticLinehaulPrices,
	}

	// 7: 	2b) Dom. Service Area Prices
	xlsxDataSheets[7] = xlsxDataSheetInfo{
		description:    stringPointer("2b) Dom. Service Area Prices"),
		outputFilename: stringPointer("2b_domestic_service_area_prices"),
		process:        &parseDomesticServiceAreaPrices,
		verify:         &verifyDomesticServiceAreaPrices,
	}

}

type paramConfig struct {
	processAll   bool
	showOutput   bool
	xlsxFilename *string
	xlsxSheets   []string
	saveToFile   bool
	runTime      time.Time
	xlsxFile     *xlsx.File
}

func xlsxSheetsUsage() string {
	message := "Provide comma separated string of sequential sheet index numbers starting with 0:\n"
	message += "\t e.g. '-xlsxSheets=\"6,7,11\"'\n"
	message += "\t      '-xlsxSheets=\"6\"'\n"
	message += "\n"
	message += "Available sheets for parsing are: \n"

	for i, v := range xlsxDataSheets {
		if v.process != nil {
			description := ""
			if v.description != nil {
				description = *v.description
			}
			message += fmt.Sprintf("%d:  %s\n", i, description)
		}
	}

	return message
}

func main() {
	initDataSheetInfo()
	params := paramConfig{}
	params.runTime = time.Now()

	flag := pflag.CommandLine
	filename := flag.String("filename", "", "Filename including path of the XLSX to parse for Rate Engine GHC import")
	all := flag.Bool("all", true, "Parse entire Rate Engine GHC XLSX")
	sheets := flag.String("xlsxSheets", "", xlsxSheetsUsage())
	display := flag.Bool("display", false, "Display output of parsed info")
	saveToFile := flag.Bool("save", false, "Save output to CSV file")

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false

	flag.Parse(os.Args[1:])

	params.processAll = false
	if all != nil && *all == true {
		params.processAll = true
	}

	// option `xlsxSheets` will override `all` flag
	if sheets != nil && len(*sheets) > 0 {
		// If processes based on `xlsxSheets` indices provided as arguments
		// process those and do not run all
		params.processAll = false
		params.xlsxSheets = strings.Split(*sheets, ",")
	}

	params.xlsxFilename = filename
	if filename != nil {
		log.Printf("Importing file %s\n", *filename)
	} else {
		log.Fatalf("Did not receive an XLSX filename to parse, missing -filename\n")
	}

	xlsxFile, err := xlsx.OpenFile(*params.xlsxFilename)
	params.xlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", *params.xlsxFilename, err)
	}

	params.showOutput = false
	if display != nil && *display == true {
		params.showOutput = true
	}

	params.saveToFile = false
	if saveToFile != nil && *saveToFile == true {
		params.saveToFile = true
	}

	// Connect to the database
	//DB connection
	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = cli.CheckDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	// Create a connection to the DB
	db, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}
	defer db.Close()

	err = db.Transaction(func(connection *pop.Connection) error {
		if params.processAll == true {
			for i, x := range xlsxDataSheets {
				if x.process != nil {
					dbErr := process(params, i, db)
					if dbErr != nil {
						log.Fatalf("Error processing xlsxDataSheets %v\n", dbErr.Error())
						return dbErr
					}
				}
			}
		} else {
			for _, v := range params.xlsxSheets {
				index, dbErr := strconv.Atoi(v)
				if dbErr != nil {
					log.Fatalf("Bad xlsxSheets index provided %v\n", dbErr)
					return dbErr
				}
				dbErr = process(params, index, db)
				if dbErr != nil {
					log.Fatalf("Error processing %v\n", dbErr)
					return dbErr
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Tranaction failed:- %v", err)
	}
}

// process: is the main process function. It will call the
// appropriate verify and process functions based on what is defined
// in the xlsxDataSheets array
//
// Should not need to edit this function when adding new processing functions
//     to add new processing functions update:
//         a.) add new verify function for your processing
//         b.) add new process function for your processing
//         c.) update initDataSheetInfo() with a.) and b.)
func process(params paramConfig, sheetIndex int, db *pop.Connection) error {
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
		err := callFunc(params, sheetIndex)
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
		err := callFunc(params, sheetIndex, db)
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
			f, ferr := strconv.ParseFloat(from, 32)
			if ferr != nil {
				fmt.Printf("ERROR: getInt() ParseFloat error %s\n", ferr.Error())
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

func removeFirstDollarSign(s string) string {
	return strings.Replace(s, "$", "", 1)
}

func splitZip3s(s string) []string {
	if strings.Contains(s, ",") {
		return strings.Split(s, ",")
	}
	i := fmt.Sprintf("%03d", getInt(s))
	return []string{i}
}

func formatServiceAreaNumber(s string) string {
	return fmt.Sprintf("%03d", getInt(s))
}

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

/*************************************************************************/
// GHC Rate Engine XLSX Verification and Process functions
/*************************************************************************/

// verifyDomesticLinehaulPrices: verification for 2a) Domestic Linehaul Prices
var verifyDomesticLinehaulPrices verifyXlsxSheet = func(params paramConfig, sheetIndex int) error {

	if dLhWeightBandNumCells != dLhWeightBandNumCellsExpected {
		return fmt.Errorf("parseDomesticLinehaulPrices(): Exepected %d columns per weight band, found %d defined in golang parser", dLhWeightBandNumCellsExpected, dLhWeightBandNumCells)
	}

	if len(dLhWeightBands) != dLhWeightBandCountExpected {
		return fmt.Errorf("parseDomesticLinehaulPrices(): Exepected %d weight bands, found %d defined in golang parser", dLhWeightBandCountExpected, len(dLhWeightBands))
	}

	log.Println("TODO verifyDomesticLinehaulPrices() not implemented")
	return nil
}

// parseDomesticLinehaulPrices: parser for 2a) Domestic Linehaul Prices
var parseDomesticLinehaulPrices processXlsxSheet = func(params paramConfig, sheetIndex int, db *pop.Connection) error {
	// Create CSV writer to save data to CSV file, returns nil if params.saveToFile=false
	csvWriter := createCsvWriter(params.saveToFile, sheetIndex, params.runTime)
	if csvWriter != nil {
		defer csvWriter.close()

		// Write header to CSV
		dp := domesticLineHaulPrice{}
		csvWriter.write(dp.csvHeader())
	}

	// XLSX Sheet consts
	const xlsxDataSheetNum int = 6  // 2a) Domestic Linehaul Prices
	const feeColIndexStart int = 6  // start at column 6 to get the rates
	const feeRowIndexStart int = 14 // start at row 14 to get the rates
	const serviceAreaNumberColumn int = 2
	const originServiceAreaColumn int = 3
	const serviceScheduleColumn int = 4
	const numEscalationYearsToProcess int = 4

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("parseDomesticLinehaulPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	dataRows := params.xlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart:]
	for _, row := range dataRows {
		colIndex := feeColIndexStart
		// For number of baseline + escalation years
		for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
			// For each rate season
			for _, r := range rateTypes {
				// For each weight band
				for _, w := range dLhWeightBands {
					// For each milage range
					for _, m := range dLhMilesRanges {
						domPrice := domesticLineHaulPrice{
							ServiceAreaNumber: formatServiceAreaNumber(getCell(row.Cells, serviceAreaNumberColumn)),
							OriginServiceArea: getCell(row.Cells, originServiceAreaColumn),
							ServiceSchedule:   getInt(getCell(row.Cells, serviceScheduleColumn)),
							Season:            r,
							WeightBand:        w,
							MilesRange:        m,
							Escalation:        escalation,
							Rate:              getCell(row.Cells, colIndex),
						}
						colIndex++
						if params.showOutput == true {
							log.Println(domPrice.toSlice())
						}
						if csvWriter != nil {
							csvWriter.write(domPrice.toSlice())
						}
					}
				}
				colIndex++ // skip 1 column (empty column) before starting next rate type
			}
		}
	}

	return nil
}

// verifyDomesticServiceAreaPrices: verification 2b) Dom. Service Area Prices
var verifyDomesticServiceAreaPrices verifyXlsxSheet = func(params paramConfig, sheetIndex int) error {
	log.Println("TODO verifyDomesticServiceAreaPrices() not implemented")
	return nil
}

// parseDomesticServiceAreaPrices: parser for: 2b) Dom. Service Area Prices
var parseDomesticServiceAreaPrices processXlsxSheet = func(params paramConfig, sheetIndex int, db *pop.Connection) error {
	// Create CSV writer to save data to CSV file, returns nil if params.saveToFile=false
	csvWriter := createCsvWriter(params.saveToFile, sheetIndex, params.runTime)
	if csvWriter != nil {
		defer csvWriter.close()

		// Write header to CSV
		dp := domesticServiceAreaPrice{}
		csvWriter.write(dp.csvHeader())
	}

	// XLSX Sheet consts
	const xlsxDataSheetNum int = 7  // 2b) Domestic Service Area Prices
	const feeColIndexStart int = 6  // start at column 6 to get the rates
	const feeRowIndexStart int = 10 // start at row 10 to get the rates
	const serviceAreaNumberColumn int = 2
	const originServiceAreaColumn int = 3
	const serviceScheduleColumn int = 4
	const sITPickupDeliveryScheduleColumn int = 5
	const numEscalationYearsToProcess int = 4

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("parseDomesticServiceAreaPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	dataRows := params.xlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart:]
	for _, row := range dataRows {
		colIndex := feeColIndexStart
		// For number of baseline + escalation years
		for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
			// For each rate season
			for _, r := range rateTypes {
				domPrice := domesticServiceAreaPrice{
					ServiceAreaNumber:         formatServiceAreaNumber(getCell(row.Cells, serviceAreaNumberColumn)),
					OriginServiceArea:         getCell(row.Cells, originServiceAreaColumn),
					ServiceSchedule:           getInt(getCell(row.Cells, serviceScheduleColumn)),
					SITPickupDeliverySchedule: getInt(getCell(row.Cells, sITPickupDeliveryScheduleColumn)),
					Season:                    r,
					Escalation:                escalation,
				}

				domPrice.ShorthaulPrice = removeFirstDollarSign(getCell(row.Cells, colIndex))
				colIndex++
				domPrice.OriginDestinationPrice = removeFirstDollarSign(getCell(row.Cells, colIndex))
				colIndex += 3 // skip 2 columns pack and unpack
				domPrice.OriginDestinationSITFirstDayWarehouse = removeFirstDollarSign(getCell(row.Cells, colIndex))
				colIndex++
				domPrice.OriginDestinationSITAddlDays = removeFirstDollarSign(getCell(row.Cells, colIndex))
				colIndex++ // skip column SIT Pickup / Delivery ≤50 miles (per cwt)

				if params.showOutput == true {
					log.Println(domPrice.toSlice())
				}
				if csvWriter != nil {
					csvWriter.write(domPrice.toSlice())
				}

				colIndex += 2 // skip 1 column (empty column) before starting next rate type
			}

		}
	}

	return nil
}

// parseServiceAreas: parser for: 1b) Service Areas
var parseServiceAreas processXlsxSheet = func(params paramConfig, sheetIndex int, db *pop.Connection) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 4          // 1b) Service Areas
	const serviceAreaRowIndexStart int = 10 // start at row 10 to get the rates
	const basePointCityColumn int = 2
	const stateColumn int = 3
	const serviceAreaNumberColumn int = 4
	const zip3sColumn int = 5
	const internationalRateAreaColumn int = 9
	const rateAreaIDColumn int = 10

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("parseServiceAreas expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	log.Println("Parsing Domestic Service Areas")
	// Create CSV writer to save data to CSV file, returns nil if params.saveToFile=false
	csvWriter := createCsvWriter(params.saveToFile, sheetIndex, params.runTime)
	if csvWriter != nil {
		defer csvWriter.close()

		// Write header to CSV
		dsa := domesticServiceArea{}
		csvWriter.write(dsa.csvHeader())
	}

	dataRows := params.xlsxFile.Sheets[xlsxDataSheetNum].Rows[serviceAreaRowIndexStart:]
	for _, row := range dataRows {
		domServArea := domesticServiceArea{
			BasePointCity:     getCell(row.Cells, basePointCityColumn),
			State:             getCell(row.Cells, stateColumn),
			ServiceAreaNumber: formatServiceAreaNumber(getCell(row.Cells, serviceAreaNumberColumn)),
			Zip3s:             splitZip3s(getCell(row.Cells, zip3sColumn)),
		}
		// All the rows are consecutive, if we get to a blank one we're done
		if domServArea.BasePointCity == "" {
			break
		} else if csvWriter != nil {
			csvWriter.write(domServArea.toSlice())
		}
		domServArea.saveToDatabase(db)
	}

	log.Println("Parsing International Service Areas")
	// Create CSV writer to save data to CSV file, returns nil if params.saveToFile=false
	if csvWriter != nil {
		// Write header to CSV
		isa := internationalServiceArea{}
		csvWriter.write(isa.csvHeader())
	}

	for _, row := range dataRows {
		intlServArea := internationalServiceArea{
			RateArea:   getCell(row.Cells, internationalRateAreaColumn),
			RateAreaID: getCell(row.Cells, rateAreaIDColumn),
		}
		// All the rows are consecutive, if we get to a blank one we're done
		if intlServArea.RateArea == "" {
			break
		} else if csvWriter != nil {
			csvWriter.write(intlServArea.toSlice())
		}
	}
	return nil
}

// verifyServiceAreas: verification for: 1b) Service Areas
var verifyServiceAreas verifyXlsxSheet = func(params paramConfig, sheetIndex int) error {
	log.Println("TODO verifyServiceAreas() not implemented")
	return nil
}
