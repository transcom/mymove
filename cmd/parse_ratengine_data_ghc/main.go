package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/dbtools"
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
const sharedNumEscalationYearsToProcess int = 1
const xlsxSheetsCountMax int = 35

type processXlsxSheet func(paramConfig, int, services.TableFromSliceCreator, *createCsvHelper) error
type verifyXlsxSheet func(paramConfig, int) error

type xlsxDataSheetInfo struct {
	description    *string
	processMethods []xlsxProcessInfo
	verify         *verifyXlsxSheet
	outputFilename *string //do not include suffix see func generateOutputFilename for details
}

var xlsxDataSheets []xlsxDataSheetInfo

type xlsxProcessInfo struct {
	process    *processXlsxSheet
	adtlSuffix *string
}

// initDataSheetInfo: When adding new functions for parsing sheets, must add new xlsxDataSheetInfo
// defining the parse function
//
// The index MUST match the sheet that is being processed. Refer to file comments or XLSX to
// determine the correct index to add.
func initDataSheetInfo() {
	xlsxDataSheets = make([]xlsxDataSheetInfo, xlsxSheetsCountMax, xlsxSheetsCountMax)

	// 4: 	1b) Domestic & International Service Areas
	xlsxDataSheets[4] = xlsxDataSheetInfo{
		description:    swag.String("1b) Service Areas"),
		outputFilename: swag.String("1b_service_areas"),
		processMethods: []xlsxProcessInfo{
			{
				process:    &parseDomesticServiceAreas,
				adtlSuffix: swag.String("domestic"),
			},
			{
				process:    &parseInternationalServiceAreas,
				adtlSuffix: swag.String("international"),
			},
		},
		verify: &verifyServiceAreas,
	}

	// 6: 	2a) Domestic Linehaul Prices
	xlsxDataSheets[6] = xlsxDataSheetInfo{
		description:    swag.String("2a) Domestic Linehaul Prices"),
		outputFilename: swag.String("2a_domestic_linehaul_prices"),
		processMethods: []xlsxProcessInfo{{
			process: &parseDomesticLinehaulPrices,
		},
		},
		verify: &verifyDomesticLinehaulPrices,
	}

	// 7: 	2b) Dom. Service Area Prices
	xlsxDataSheets[7] = xlsxDataSheetInfo{
		description:    swag.String("2b) Dom. Service Area Prices"),
		outputFilename: swag.String("2b_domestic_service_area_prices"),
		processMethods: []xlsxProcessInfo{{
			process: &parseDomesticServiceAreaPrices,
		},
		},
		verify: &verifyDomesticServiceAreaPrices,
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
	runVerify    bool
}

func xlsxSheetsUsage() string {
	message := "Provide comma separated string of sequential sheet index numbers starting with 0:\n"
	message += "\t e.g. '-xlsxSheets=\"6,7,11\"'\n"
	message += "\t      '-xlsxSheets=\"6\"'\n"
	message += "\n"
	message += "Available sheets for parsing are: \n"

	for i, v := range xlsxDataSheets {
		if len(v.processMethods) > 0 {
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
	runVerify := flag.Bool("verify", true, "Default is true, if false skip sheet format verification")

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false

	err := flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("Could not parse flags: %v\n", err)
	}

	// Process command line params

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

	params.runVerify = false
	if runVerify != nil {
		params.runVerify = *runVerify
	}

	// Connect to the database
	//DB connection
	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		log.Fatalf("Could not bind flags: %v\n", err)
	}
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
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Fatalf("Could not close database: %v", closeErr)
		}
	}()

	tableFromSliceCreator := dbtools.NewTableFromSliceCreator(db, logger, true)

	// Must be after processing config param
	// Run the process function

	err = db.Transaction(func(connection *pop.Connection) error {
		if params.processAll == true {
			for i, x := range xlsxDataSheets {
				for _, p := range x.processMethods {
					if p.process != nil {
						dbErr := process(params, i, tableFromSliceCreator)
						if dbErr != nil {
							log.Printf("Error processing xlsxDataSheets %v\n", dbErr.Error())
							return dbErr
						}
					}
				}
			}
		} else {
			for _, v := range params.xlsxSheets {
				index, dbErr := strconv.Atoi(v)
				if dbErr != nil {
					log.Printf("Bad xlsxSheets index provided %v\n", dbErr)
					return dbErr
				}
				if index < len(xlsxDataSheets) {
					dbErr = process(params, index, tableFromSliceCreator)
					if dbErr != nil {
						log.Printf("Error processing %v\n", dbErr)
						return dbErr
					}
				} else {
					log.Printf("Error processing index %d, not in range of slice xlsxDataSheets\n", index)
					return errors.New("Index out of range of slice xlsxDataSheets")
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Transaction failed:- %v", err)
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
func process(params paramConfig, sheetIndex int, tableFromSliceCreator services.TableFromSliceCreator) error {
	xlsxInfo := xlsxDataSheets[sheetIndex]
	var description string
	if xlsxInfo.description != nil {
		description = *xlsxInfo.description
		log.Printf("Processing sheet index %d with description %s\n", sheetIndex, description)
	} else {
		log.Printf("Processing sheet index %d with missing description\n", sheetIndex)
	}

	// Call verify function
	if params.runVerify == true {
		if xlsxInfo.verify != nil {
			var callFunc verifyXlsxSheet
			callFunc = *xlsxInfo.verify
			err := callFunc(params, sheetIndex)
			if err != nil {
				log.Printf("%s verify error: %v\n", description, err)
				return errors.Wrapf(err, " verify error for sheet index: %d with description: %s", sheetIndex, description)
			}
		} else {
			log.Printf("No verify function for sheet index %d with description %s\n", sheetIndex, description)
		}
	} else {
		log.Print("Skip running the verify functions")
	}

	// Call process function
	if len(xlsxInfo.processMethods) > 0 {
		for _, p := range xlsxInfo.processMethods {
			// Create CSV writer to save data to CSV file, returns nil if params.saveToFile=false
			csvWriter := createCsvWriter(params.saveToFile, sheetIndex, params.runTime, p.adtlSuffix)
			if csvWriter != nil {
				defer csvWriter.close()
			}
			var callFunc processXlsxSheet
			callFunc = *p.process
			err := callFunc(params, sheetIndex, tableFromSliceCreator, csvWriter)
			if err != nil {
				log.Printf("%s process error: %v\n", description, err)
				return errors.Wrapf(err, " process error for sheet index: %d with description: %s", sheetIndex, description)
			}
		}
	} else {
		log.Fatalf("Missing process function for sheet index %d with description %s\n", sheetIndex, description)
	}

	// Verification and Process completed
	log.Printf("Completed processing sheet index %d with description %s\n", sheetIndex, description)
	return nil
}
