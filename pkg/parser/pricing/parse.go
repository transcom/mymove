package pricing

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/tealeg/xlsx/v3"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/dbtools"
)

/*************************************************************************

Parser tool to extract data from the GHC Rate Engine XLSX

For help run: <program> -h

`go run ./cmd/ghc-pricing-parser/ -h`

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
	c.) Update InitDataSheetInfo() with a.) and b.)
		The index must match the sheet index in the XLSX that you aim to process

You should not have to update the Parse() or process() functions unless you
intentionally are modifying the pattern of how the processing functions are called.

 *************************************************************************/

const xlsxSheetsCountMax int = 35

type processXlsxSheet func(appcontext.AppContext, ParamConfig, int) (interface{}, error)
type verifyXlsxSheet func(ParamConfig, int) error

// XlsxDataSheetInfo is the xlsx data sheet info
type XlsxDataSheetInfo struct {
	Description    *string
	ProcessMethods []xlsxProcessInfo
	verify         *verifyXlsxSheet
	outputFilename *string //do not include suffix see func generateOutputFilename for details
}

type xlsxProcessInfo struct {
	process     *processXlsxSheet
	description *string
	adtlSuffix  *string
}

// ParamConfig is the parameter conifguration
type ParamConfig struct {
	ProcessAll        bool
	ShowOutput        bool
	XlsxFilename      string
	XlsxSheets        []string
	SaveToFile        bool
	RunTime           time.Time
	XlsxFile          *xlsx.File
	RunVerify         bool
	RunImport         bool
	UseTempTables     bool
	DropIfExists      bool
	ContractCode      string
	ContractName      string
	ContractStartDate string
}

// InitDataSheetInfo - When adding new functions for parsing sheets, must add new XlsxDataSheetInfo
// defining the parse function
//
// The index MUST match the sheet that is being processed. Refer to file comments or XLSX to
// determine the correct index to add.
func InitDataSheetInfo() []XlsxDataSheetInfo {
	xlsxDataSheets := make([]XlsxDataSheetInfo, xlsxSheetsCountMax)

	// 4: 	1b) Domestic & International Service Areas
	xlsxDataSheets[4] = XlsxDataSheetInfo{
		Description:    swag.String("1b) Service Areas"),
		outputFilename: swag.String("1b_service_areas"),
		ProcessMethods: []xlsxProcessInfo{
			{
				process:     &parseDomesticServiceAreas,
				description: swag.String("domestic service areas"),
				adtlSuffix:  swag.String("domestic"),
			},
			{
				process:     &parseInternationalServiceAreas,
				description: swag.String("international service areas"),
				adtlSuffix:  swag.String("international"),
			},
		},
		verify: &verifyServiceAreas,
	}

	// 6: 	2a) Domestic Linehaul Prices
	xlsxDataSheets[6] = XlsxDataSheetInfo{
		Description:    swag.String("2a) Domestic Linehaul Prices"),
		outputFilename: swag.String("2a_domestic_linehaul_prices"),
		ProcessMethods: []xlsxProcessInfo{{
			process:     &parseDomesticLinehaulPrices,
			description: swag.String("domestic linehaul prices"),
		},
		},
		verify: &verifyDomesticLinehaulPrices,
	}

	// 7: 	2b) Dom. Service Area Prices
	xlsxDataSheets[7] = XlsxDataSheetInfo{
		Description:    swag.String("2b) Dom. Service Area Prices"),
		outputFilename: swag.String("2b_domestic_service_area_prices"),
		ProcessMethods: []xlsxProcessInfo{{
			process:     &parseDomesticServiceAreaPrices,
			description: swag.String("domestic service area prices"),
		},
		},
		verify: &verifyDomesticServiceAreaPrices,
	}

	// 8: 	2c) Dom. Other Prices
	xlsxDataSheets[8] = XlsxDataSheetInfo{
		Description:    swag.String("2c) Dom. Other Prices"),
		outputFilename: swag.String("2c_domestic_other_prices"),
		ProcessMethods: []xlsxProcessInfo{
			{
				process:     &parseDomesticOtherPricesPack,
				description: swag.String("domestic other (pack/unpack) prices"),
				adtlSuffix:  swag.String("pack"),
			},
			{
				process:     &parseDomesticOtherPricesSit,
				description: swag.String("domestic other (SIT pickup/delivery) prices"),
				adtlSuffix:  swag.String("sit"),
			},
		},
		verify: &verifyDomesticOtherPrices,
	}

	// 10: 	3a) OCONUS TO OCONUS Prices
	xlsxDataSheets[10] = XlsxDataSheetInfo{
		Description:    swag.String("3a) OCONUS to OCONUS Prices"),
		outputFilename: swag.String("3a_oconus_to_oconus_prices"),
		ProcessMethods: []xlsxProcessInfo{{
			process:     &parseOconusToOconusPrices,
			description: swag.String("OCONUS to OCONUS prices"),
		},
		},
		verify: &verifyIntlOconusToOconusPrices,
	}

	// 11: 	3b) CONUS TO OCONUS Prices
	xlsxDataSheets[11] = XlsxDataSheetInfo{
		Description:    swag.String("3b) CONUS to OCONUS Prices"),
		outputFilename: swag.String("3b_conus_to_oconus_prices"),
		ProcessMethods: []xlsxProcessInfo{{
			process:     &parseConusToOconusPrices,
			description: swag.String("CONUS to OCONUS prices"),
		},
		},
		verify: &verifyIntlConusToOconusPrices,
	}

	// 12: 	3c) OCONUS TO CONUS Prices
	xlsxDataSheets[12] = XlsxDataSheetInfo{
		Description:    swag.String("3c) OCONUS to CONUS Prices"),
		outputFilename: swag.String("3c_oconus_to_conus_prices"),
		ProcessMethods: []xlsxProcessInfo{{
			process:     &parseOconusToConusPrices,
			description: swag.String("OCONUS to CONUS prices"),
		},
		},
		verify: &verifyIntlOconusToConusPrices,
	}

	// 14: 	3e) Non-Standard Locn Prices
	xlsxDataSheets[14] = XlsxDataSheetInfo{
		Description:    swag.String("3e) Non-Standard Loc'n Prices"),
		outputFilename: swag.String("3e_non_standard_locn_prices"),
		ProcessMethods: []xlsxProcessInfo{{
			process:     &parseNonStandardLocnPrices,
			description: swag.String("non-standard location prices"),
		},
		},
		verify: &verifyNonStandardLocnPrices,
	}

	// 18:	5b) Price Escalation Discount
	xlsxDataSheets[18] = XlsxDataSheetInfo{
		Description:    swag.String("5b) Price Escalation Discount"),
		outputFilename: swag.String("5b_price_escalation_discount"),
		ProcessMethods: []xlsxProcessInfo{{
			process:     &parsePriceEscalationDiscount,
			description: swag.String("price escalation discount"),
		},
		},
		verify: &verifyPriceEscalationDiscount,
	}

	// 13: 	5a) Other International Prices
	xlsxDataSheets[13] = XlsxDataSheetInfo{
		Description:    swag.String("3d) Other International Prices"),
		outputFilename: swag.String("3d_other_international_prices"),
		ProcessMethods: []xlsxProcessInfo{{
			process:     &parseOtherIntlPrices,
			description: swag.String("other international prices"),
		},
		},
		verify: &verifyOtherIntlPrices,
	}

	// 16: 	4a) Mgmt., Coun., Trans. Prices
	xlsxDataSheets[16] = XlsxDataSheetInfo{
		Description:    swag.String("4a) Mgmt., Coun., Trans. Prices"),
		outputFilename: swag.String("4a_mgmt_coun_trans_prices"),
		ProcessMethods: []xlsxProcessInfo{
			{
				process:     &parseShipmentManagementServicesPrices,
				description: swag.String("shipment management services prices"),
				adtlSuffix:  swag.String("management"),
			},
			{
				process:     &parseCounselingServicesPrices,
				description: swag.String("counseling services prices"),
				adtlSuffix:  swag.String("counsel"),
			},
			{
				process:     &parseTransitionPrices,
				description: swag.String("transition prices"),
				adtlSuffix:  swag.String("transition"),
			},
		},
		verify: &verifyManagementCounselTransitionPrices,
	}

	// 17: 	5a) Access. and Add. Prices
	xlsxDataSheets[17] = XlsxDataSheetInfo{
		Description:    swag.String("5a) Access. and Add. Prices"),
		outputFilename: swag.String("5a_access_and_add_prices"),
		ProcessMethods: []xlsxProcessInfo{
			{
				process:     &parseDomesticMoveAccessorialPrices,
				description: swag.String("domestic move accessorial prices"),
				adtlSuffix:  swag.String("domestic"),
			},
			{
				process:     &parseInternationalMoveAccessorialPrices,
				description: swag.String("international move accessorial prices"),
				adtlSuffix:  swag.String("international"),
			},
			{
				process:     &parseDomesticInternationalAdditionalPrices,
				description: swag.String("domestic/international additional prices"),
				adtlSuffix:  swag.String("additional"),
			},
		},
		verify: &verifyAccessAndAddPrices,
	}

	return xlsxDataSheets
}

// Parse will parsh xlsx data sheet info
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
						txnAppCtx.Logger().Error("Error processing xlsxDataSheets", zap.Error(dbErr))
						return dbErr
					}
				}
			}
		} else {
			for _, v := range params.XlsxSheets {
				index, dbErr := strconv.Atoi(v)
				if dbErr != nil {
					txnAppCtx.Logger().Error("Bad XlsxSheets index provided", zap.Error(dbErr))
					return dbErr
				}
				if index < len(xlsxDataSheets) {
					dbErr = process(txnAppCtx, xlsxDataSheets, params, index, tableFromSliceCreator)
					if dbErr != nil {
						txnAppCtx.Logger().Error("Error processing", zap.Error(dbErr))
						return dbErr
					}
				} else {
					txnAppCtx.Logger().Error("Error processing index not in range of slice xlsxDataSheets", zap.Int("index", index))
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

	description := "(no description)"
	if xlsxInfo.Description != nil {
		description = *xlsxInfo.Description
	}

	pterm.Println(pterm.BgGray.Sprint(fmt.Sprintf("Processing sheet index %d: %s", sheetIndex, description)))

	// Call verify function
	if params.RunVerify {
		if xlsxInfo.verify != nil {
			callFunc := *xlsxInfo.verify
			err := callFunc(params, sheetIndex)
			if err != nil {
				appCtx.Logger().Error("Verify error", zap.String("description", description), zap.Error(err))
				return errors.Wrapf(err, " Verify error for sheet index: %d with description: %s", sheetIndex, description)
			}
		} else {
			appCtx.Logger().Info("No verify function", zap.Int("sheet index", sheetIndex), zap.String("description", description))
		}
	} else {
		appCtx.Logger().Info("Skip running the verify functions")
	}

	// Call process function
	if len(xlsxInfo.ProcessMethods) > 0 {
		for methodIndex, p := range xlsxInfo.ProcessMethods {
			if p.process != nil {
				processDescription := "(no description)"
				if p.description != nil {
					processDescription = *p.description
				}

				spinner, err := pterm.DefaultSpinner.Start(fmt.Sprintf("Processing section: %s", processDescription))
				if err != nil {
					appCtx.Logger().Fatal("Failed to create pterm spinner", zap.Error(err))
				}

				callFunc := *p.process
				slice, err := callFunc(appCtx, params, sheetIndex)
				if err != nil {
					spinner.Fail()
					appCtx.Logger().Error("process error", zap.String("description", description), zap.Error(err))
					return errors.Wrapf(err, "process error for sheet index: %d with description: %s", sheetIndex, description)
				}

				if params.SaveToFile {
					filename := xlsxDataSheets[sheetIndex].generateOutputFilename(sheetIndex, params.RunTime, p.adtlSuffix)
					if err := createCSV(appCtx, filename, slice); err != nil {
						spinner.Fail()
						return errors.Wrapf(err, "Could not create CSV for sheet index: %d with description: %s", sheetIndex, description)
					}
				}
				if err := tableFromSliceCreator.CreateTableFromSlice(appCtx, slice); err != nil {
					spinner.Fail()
					return errors.Wrapf(err, "Could not create table for sheet index: %d with description: %s", sheetIndex, description)
				}

				spinner.Success()
			} else {
				appCtx.Logger().Info("No process function", zap.Int("sheet index", sheetIndex), zap.String("description", description), zap.Int("method index", methodIndex))
			}
		}
	} else {
		appCtx.Logger().Fatal("Missing process function", zap.Int("sheet index", sheetIndex), zap.String("description", description))
	}

	return nil
}

func createCSV(appCtx appcontext.AppContext, filename string, slice interface{}) error {
	// Create file for writing the CSV
	csvFile, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "Could not create CSV file")
	}
	defer func() {
		if closeErr := csvFile.Close(); closeErr != nil {
			appCtx.Logger().Fatal("Could not close CSV file", zap.Error(closeErr))
		}
	}()

	// Write the CSV
	if err := gocsv.MarshalFile(slice, csvFile); err != nil {
		return errors.Wrapf(err, "Could not marshal CSV file")
	}

	return nil
}
