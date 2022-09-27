package transittime

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx/v3"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/dbtools"
)

/*************************************************************************

Parser tool to extract data from the GHC Transit Times XLSX

For help run: <program> -h

`go run ./cmd/ghc-hhg-transit-times-parser/ -h`

Transit Times XLSX sections this tool will be parsing:

1) Domestic HHG Transit Times Tab

2) International HHG and Unaccompanied Baggage (UB) - NOT IMPLEMENTED YET

--------------------------------------------------------------------------

Transit Times XLSX sheet tabs:

0: 	Instructions for Transit Times
1: 	Domestic Transit Times
2: 	International Transit Times
3:  Domestic Transit Times - Blank
4:  International Transit Times - Blank

 *************************************************************************/

/*************************************************************************

To add new parser functions to this file:

	a.) (optional) Add new verify function for your processing must match signature verifyXlsxSheet
	b.) Add new process function to process XLSX data sheet must match signature processXlsxSheet
	c.) Update InitDataSheetInfo() with a.) and b.)
		The index must match the sheet index in the XLSX that you aim to process

You should not have to update the Parse() or process() functions unless you
intentionally are modifying the pattern of how the processing functions are called.

 *************************************************************************/

const xlsxSheetsCountMax int = 35

type processXlsxSheet func(ParamConfig, int, *zap.Logger) (interface{}, error)
type verifyXlsxSheet func(ParamConfig, int, *zap.Logger) error

// XlsxDataSheetInfo describes the excel sheet info
type XlsxDataSheetInfo struct {
	Description    *string
	ProcessMethods []xlsxProcessInfo
	verify         *verifyXlsxSheet
	outputFilename *string //do not include suffix see func generateOutputFilename for details
}

type xlsxProcessInfo struct {
	process    *processXlsxSheet
	adtlSuffix *string
}

// ParamConfig config for the transit time parser program
type ParamConfig struct {
	ProcessAll    bool
	ShowOutput    bool
	XlsxFilename  string
	XlsxSheets    []string
	SaveToFile    bool
	RunTime       time.Time
	XlsxFile      *xlsx.File
	RunVerify     bool
	RunImport     bool
	UseTempTables bool
	DropIfExists  bool
}

// InitDataSheetInfo : When adding new functions for parsing sheets, must add new XlsxDataSheetInfo
// defining the parse function
//
// The index MUST match the sheet that is being processed. Refer to file comments or XLSX to
// determine the correct index to add.
func InitDataSheetInfo() []XlsxDataSheetInfo {
	xlsxDataSheets := make([]XlsxDataSheetInfo, xlsxSheetsCountMax)

	// Tab Index
	// 1: 	Domestic Transit Times
	xlsxDataSheets[1] = XlsxDataSheetInfo{
		Description:    swag.String("HHG Domestic Transit Times"),
		outputFilename: swag.String("hhg_domestic_transit_times"),
		ProcessMethods: []xlsxProcessInfo{
			{
				process:    &parseDomesticTransitTime,
				adtlSuffix: swag.String("domestic"),
			},
		},
		verify: &verifyTransitTime,
	}

	return xlsxDataSheets
}

// Parse parses the XLSX file
func Parse(appCtx appcontext.AppContext, xlsxDataSheets []XlsxDataSheetInfo, params ParamConfig) error {
	// Must be after processing config param
	// Run the process function

	err := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		tableFromSliceCreator := dbtools.NewTableFromSliceCreator(params.UseTempTables, params.DropIfExists)

		if params.ProcessAll {
			for i, x := range xlsxDataSheets {
				if len(x.ProcessMethods) >= 1 {
					dbErr := process(txnAppCtx, xlsxDataSheets, params, i, tableFromSliceCreator)
					if dbErr != nil {
						log.Printf("Error processing xlsxDataSheets %v\n", dbErr.Error())
						return dbErr
					}
				}
			}
		} else {
			for _, v := range params.XlsxSheets {
				index, dbErr := strconv.Atoi(v)
				if dbErr != nil {
					log.Printf("Bad XlsxSheets index provided %v\n", dbErr)
					return dbErr
				}
				if index < len(xlsxDataSheets) {
					dbErr = process(appCtx, xlsxDataSheets, params, index, tableFromSliceCreator)
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
		return errors.Wrap(err, "Transaction failed")
	}

	return nil
}

// process: is the main process function. It will call the
// appropriate verify and process functions based on what is defined
// in the xlsxDataSheets array
//
// Should not need to edit this function when adding new processing functions
//
//	to add new processing functions update:
//	    a.) add new verify function for your processing
//	    b.) add new process function for your processing
//	    c.) update InitDataSheetInfo() with a.) and b.)
func process(appCtx appcontext.AppContext, xlsxDataSheets []XlsxDataSheetInfo, params ParamConfig, sheetIndex int, tableFromSliceCreator services.TableFromSliceCreator) error {
	xlsxInfo := xlsxDataSheets[sheetIndex]
	var description string
	if xlsxInfo.Description != nil {
		description = *xlsxInfo.Description
		log.Printf("Processing sheet index %d with Description %s\n", sheetIndex, description)
	} else {
		log.Printf("Processing sheet index %d with missing Description\n", sheetIndex)
	}

	// Call verify function
	if params.RunVerify {
		if xlsxInfo.verify != nil {
			callFunc := *xlsxInfo.verify
			err := callFunc(params, sheetIndex, appCtx.Logger())
			if err != nil {
				log.Printf("%s Verify error: %v\n", description, err)
				return errors.Wrapf(err, " Verify error for sheet index: %d with Description: %s", sheetIndex, description)
			}
		} else {
			log.Printf("No verify function for sheet index %d with Description %s\n", sheetIndex, description)
		}
	} else {
		log.Print("Skip running the verify functions")
	}

	// Call process function
	if len(xlsxInfo.ProcessMethods) > 0 {
		for methodIndex, processMethods := range xlsxInfo.ProcessMethods {
			if processMethods.process != nil {
				callFunc := *processMethods.process
				slice, err := callFunc(params, sheetIndex, appCtx.Logger())
				if err != nil {
					log.Printf("%s process error: %v\n", description, err)
					return errors.Wrapf(err, " process error for sheet index: %d with Description: %s", sheetIndex, description)
				}

				if params.SaveToFile {
					filename := xlsxDataSheets[sheetIndex].generateOutputFilename(sheetIndex, params.RunTime, processMethods.adtlSuffix)
					if err := createCSV(filename, slice); err != nil {
						return errors.Wrapf(err, "Could not create CSV for sheet index: %d with Description: %s", sheetIndex, description)
					}
					log.Println("File created:")
					log.Println(filename)
				}
				// ToDo: needs extra work
				//if err := tableFromSliceCreator.CreateTableFromSlice(slice); err != nil {
				//	return errors.Wrapf(err, "Could not create table for sheet index: %d with Description: %s", sheetIndex, description)
				//}
			} else {
				log.Printf("No process function for sheet index %d with Description %s method index: %d\n", sheetIndex, description, methodIndex)
			}
		}
	} else {
		log.Fatalf("Missing process function for sheet index %d with Description %s\n", sheetIndex, description)
	}

	// Verification and Process completed
	log.Printf("Completed processing sheet index %d with Description %s\n", sheetIndex, description)
	return nil
}

func createCSV(filename string, slice interface{}) error {
	// Create file for writing the CSV
	csvFile, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "Could not create CSV file")
	}
	defer func() {
		if closeErr := csvFile.Close(); closeErr != nil {
			log.Fatalf("Could not close CSV file: %v", closeErr)
		}
	}()

	// Write the CSV
	if err := gocsv.MarshalFile(slice, csvFile); err != nil {
		return errors.Wrapf(err, "Could not marshal CSV file")
	}

	return nil
}
