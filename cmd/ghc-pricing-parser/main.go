package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/tealeg/xlsx"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/parser/pricing"
	"github.com/transcom/mymove/pkg/services/ghcimport"
)

/*************************************************************************

Parser tool to extract data from the GHC Rate Engine XLSX

For help run: <program> -h

`go run ./cmd/ghc-pricing-parser/ -h`

 *************************************************************************/

func main() {
	xlsxDataSheets := pricing.InitDataSheetInfo()

	params := pricing.ParamConfig{}
	params.RunTime = time.Now()

	flag := pflag.CommandLine
	flag.StringVar(&params.XlsxFilename, "filename", "", "Filename including path of the XLSX to parse for Rate Engine GHC import")
	flag.BoolVar(&params.ProcessAll, "all", true, "Parse entire Rate Engine GHC XLSX")
	flag.StringSliceVar(&params.XlsxSheets, "xlsxSheets", []string{}, xlsxSheetsUsage(xlsxDataSheets))
	flag.BoolVar(&params.ShowOutput, "display", false, "Display output of parsed info")
	flag.BoolVar(&params.SaveToFile, "save-csv", false, "Save output to CSV file")
	flag.BoolVar(&params.RunVerify, "verify", true, "Default is true, if false skip sheet format verification")
	flag.BoolVar(&params.RunImport, "re-import", true, "Run GHC Rate Engine Import")
	flag.BoolVar(&params.UseTempTables, "use-temp-tables", true, "Default is true, if false stage tables are NOT temp tables")
	flag.BoolVar(&params.DropIfExists, "drop", false, "Default is false, if true stage tables will be dropped if they exist")
	flag.StringVar(&params.ContractCode, "contract-code", "", "Contract code to use for this import")
	flag.StringVar(&params.ContractName, "contract-name", "", "Contract name to use for this import")

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false

	err := flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("Could not parse flags: %v\n", err)
	}

	// option `xlsxSheets` will override `all` flag
	if len(params.XlsxSheets) > 0 {
		params.ProcessAll = false
		log.Println("Setting --xlsxSheets disables --re-import so no data will be imported into the rate engine tables. Only stage table data will be updated.")
		params.RunImport = false
	}

	if params.XlsxFilename == "" {
		log.Fatalf("Did not receive an XLSX filename to parse; missing --filename\n")
	}
	log.Printf("Importing file %s\n", params.XlsxFilename)

	if params.RunImport && params.ContractCode == "" {
		log.Fatalf("Did not receive a contract code; missing --contract-code\n")
	}

	params.XlsxFile, err = xlsx.OpenFile(params.XlsxFilename)
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", params.XlsxFilename, err)
	}

	// Connect to the database
	// DB connection
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

	// Now kick off the parsing
	err = pricing.Parse(xlsxDataSheets, params, db, logger)
	if err != nil {
		log.Fatalf("Failed to parse pricing template due to %v", err)
	}

	// If the parsing was successful, run GHC Rate Engine importer
	if params.RunImport {
		ghcREImporter := ghcimport.GHCRateEngineImporter{
			Logger:       logger,
			ContractCode: params.ContractCode,
			ContractName: params.ContractName,
		}
		err = ghcREImporter.Import(db)
		if err != nil {
			log.Fatalf("GHC Rate Engine import failed due to %v", err)
		}
	}

	// Summarize import for verification
	if err := summarizeXlsxStageParsing(db); err != nil {
		log.Fatalf("Failed to summarize XLSX to stage table parsing: %v", err)
	}

	if params.RunImport {
		if err := summarizeStageReImport(db); err != nil {
			log.Fatalf("Failed to summarize stage table to rate engine table import: %v", err)
		}
	}

}

func summarizeXlsxStageParsing(db *pop.Connection) error {
	log.Println("XLSX to Stage Table Parsing Complete")
	log.Println(" Summary:")

	// 1b Service Areas
	stageDomServiceAreas := []models.StageDomesticServiceArea{}
	if err := db.All(&stageDomServiceAreas); err != nil {
		return err
	}
	length := len(stageDomServiceAreas)
	log.Printf("   1b: Service Areas (StageDomesticServiceArea): %d\n", length)
	log.Printf("     first: %+v\n", stageDomServiceAreas[0])
	log.Printf("      last: %+v\n", stageDomServiceAreas[length-1])
	log.Println("   ---")

	stageIntlServiceArea := []models.StageInternationalServiceArea{}
	if err := db.All(&stageIntlServiceArea); err != nil {
		return err
	}
	length = len(stageIntlServiceArea)
	log.Printf("   1b: Service Areas (StageInternationalServiceArea): %d\n", length)
	log.Printf("     first: %+v\n", stageIntlServiceArea[0])
	log.Printf("      last: %+v\n", stageIntlServiceArea[length-1])
	log.Println("   ---")

	// 2a Domestic Linehaul Prices
	stageDomLinePrice := []models.StageDomesticLinehaulPrice{}
	if err := db.All(&stageDomLinePrice); err != nil {
		return err
	}
	length = len(stageDomLinePrice)
	log.Printf("   2a: Domestic Linehaul Prices (StageDomesticLinehaulPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageDomLinePrice[0])
	log.Printf("      last: %+v\n", stageDomLinePrice[length-1])
	log.Println("   ---")

	// 2b Domestic Service Area Prices
	stageDomSerAreaPrice := []models.StageDomesticServiceAreaPrice{}
	if err := db.All(&stageDomSerAreaPrice); err != nil {
		return err
	}
	length = len(stageDomSerAreaPrice)
	log.Printf("   2b: Domestic Service Area Prices (StageDomesticServiceAreaPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageDomSerAreaPrice[0])
	log.Printf("      last: %+v\n", stageDomSerAreaPrice[length-1])
	log.Println("   ---")

	// 2c Other Domestic Prices
	stageDomOtherPackPrice := []models.StageDomesticOtherPackPrice{}
	if err := db.All(&stageDomOtherPackPrice); err != nil {
		return err
	}
	length = len(stageDomOtherPackPrice)
	log.Printf("   2c: Other Domestic Prices (StageDomesticOtherPackPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageDomOtherPackPrice[0])
	log.Printf("      last: %+v\n", stageDomOtherPackPrice[length-1])
	log.Println("   ---")

	stageDomOtherSitPrice := []models.StageDomesticOtherSitPrice{}
	if err := db.All(&stageDomOtherSitPrice); err != nil {
		return err
	}
	length = len(stageDomOtherSitPrice)
	log.Printf("   2c: Other Domestic Prices (StageDomesticOtherSitPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageDomOtherSitPrice[0])
	log.Printf("      last: %+v\n", stageDomOtherSitPrice[length-1])
	log.Println("   ---")

	// 3a OCONUS to OCONUS Prices
	stageOconusToOconus := []models.StageOconusToOconusPrice{}
	if err := db.All(&stageOconusToOconus); err != nil {
		return err
	}
	length = len(stageOconusToOconus)
	log.Printf("   3a: OCONUS to OCONUS Prices (StageOconusToOconusPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageOconusToOconus[0])
	log.Printf("      last: %+v\n", stageOconusToOconus[length-1])
	log.Println("   ---")

	// 3b CONUS to OCONUS Prices
	stageConusToOconus := []models.StageConusToOconusPrice{}
	if err := db.All(&stageConusToOconus); err != nil {
		return err
	}
	length = len(stageConusToOconus)
	log.Printf("   3b: CONUS to OCONUS Prices (StageConusToOconusPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageConusToOconus[0])
	log.Printf("      last: %+v\n", stageConusToOconus[length-1])
	log.Println("   ---")

	// 3c OCONUS to CONUS Prices
	stageOconusToConus := []models.StageOconusToConusPrice{}
	if err := db.All(&stageOconusToConus); err != nil {
		return err
	}
	length = len(stageOconusToConus)
	log.Printf("   3c: OCONUS to CONUS Prices (StageOconusToConusPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageOconusToConus[0])
	log.Printf("      last: %+v\n", stageOconusToConus[length-1])
	log.Println("   ---")

	// 3d Other International Prices
	stageOtherIntlPrices := []models.StageOtherIntlPrice{}
	if err := db.All(&stageOtherIntlPrices); err != nil {
		return err
	}
	length = len(stageOtherIntlPrices)
	log.Printf("   3d: Other International Prices (StageOtherIntlPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageOtherIntlPrices[0])
	log.Printf("      last: %+v\n", stageOtherIntlPrices[length-1])
	log.Println("   ---")

	// 3e Non-Standard Location Prices
	stageNonStdLocaPrices := []models.StageNonStandardLocnPrice{}
	if err := db.All(&stageNonStdLocaPrices); err != nil {
		return err
	}
	length = len(stageNonStdLocaPrices)
	log.Printf("   3e: Non-Standard Location Prices (StageNonStandardLocnPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageNonStdLocaPrices[0])
	log.Printf("      last: %+v\n", stageNonStdLocaPrices[length-1])
	log.Println("   ---")

	// 4a Management, Counseling, and Transition Prices
	stageMgmt := []models.StageShipmentManagementServicesPrice{}
	if err := db.All(&stageMgmt); err != nil {
		return err
	}
	length = len(stageMgmt)
	log.Printf("   4a: Management, Counseling, and Transition Prices (StageShipmentManagementServicesPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageMgmt[0])
	log.Printf("      last: %+v\n", stageMgmt[length-1])
	log.Println("   ---")

	stageCounsel := []models.StageCounselingServicesPrice{}
	if err := db.All(&stageCounsel); err != nil {
		return err
	}
	length = len(stageCounsel)
	log.Printf("   4a: Management, Counseling, and Transition Prices (StageCounselingServicesPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageCounsel[0])
	log.Printf("      last: %+v\n", stageCounsel[length-1])
	log.Println("   ---")

	stageTransition := []models.StageTransitionPrice{}
	if err := db.All(&stageTransition); err != nil {
		return err
	}
	length = len(stageTransition)
	log.Printf("   4a: Management, Counseling, and Transition Prices (StageTransitionPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageTransition[0])
	log.Printf("      last: %+v\n", stageTransition[length-1])
	log.Println("   ---")

	// 5a Accessorial and Additional Prices
	stageDomMoveAccess := []models.StageDomesticMoveAccessorialPrices{}
	if err := db.All(&stageDomMoveAccess); err != nil {
		return err
	}
	length = len(stageDomMoveAccess)
	log.Printf("   5a Accessorial and Additional Prices (StageDomesticMoveAccessorialPrices): %d\n", length)
	log.Printf("     first: %+v\n", stageDomMoveAccess[0])
	log.Printf("      last: %+v\n", stageDomMoveAccess[length-1])
	log.Println("   ---")

	stageIntlMoveAccess := []models.StageInternationalMoveAccessorialPrices{}
	if err := db.All(&stageIntlMoveAccess); err != nil {
		return err
	}
	length = len(stageIntlMoveAccess)
	log.Printf("   5a Accessorial and Additional Prices (StageInternationalMoveAccessorialPrices): %d\n", length)
	log.Printf("     first: %+v\n", stageIntlMoveAccess[0])
	log.Printf("      last: %+v\n", stageIntlMoveAccess[length-1])
	log.Println("   ---")

	stageDomIntlAdd := []models.StageDomesticInternationalAdditionalPrices{}
	if err := db.All(&stageDomIntlAdd); err != nil {
		return err
	}
	length = len(stageDomIntlAdd)
	log.Printf("   5a Accessorial and Additional Prices (StageDomesticInternationalAdditionalPrices): %d\n", length)
	log.Printf("     first: %+v\n", stageDomIntlAdd[0])
	log.Printf("      last: %+v\n", stageDomIntlAdd[length-1])
	log.Println("   ---")

	// 5b Price Escalation Discount
	stagePriceEsc := []models.StagePriceEscalationDiscount{}
	if err := db.All(&stagePriceEsc); err != nil {
		return err
	}
	length = len(stagePriceEsc)
	log.Printf("   5b: Price Escalation Discount (StagePriceEscalationDiscount): %d\n", length)
	log.Printf("     first: %+v\n", stagePriceEsc[0])
	log.Printf("      last: %+v\n", stagePriceEsc[length-1])
	log.Println("   ---")

	return nil
}

func summarizeStageReImport(db *pop.Connection) error {
	log.Println("Stage Table import into Rate Engine Tables Complete")
	log.Println(" Summary:")

	// re_contract
	reContract := []models.ReContract{}
	if err := db.All(&reContract); err != nil {
		return err
	}
	length := len(reContract)
	log.Printf("   re_contract (ReContract): %d\n", length)
	log.Printf("     first: %s\n", reContract[0])
	log.Printf("      last: %s\n", reContract[length-1])
	log.Println("   ---")

	// re_contract_years
	reContractYears := []models.ReContractYear{}
	if err := db.All(&reContractYears); err != nil {
		return err
	}
	length = len(reContractYears)
	log.Printf("   re_contract_years (ReContractYear): %d\n", length)
	log.Printf("     first: %+v\n", reContractYears[0])
	log.Printf("      last: %+v\n", reContractYears[length-1])
	log.Println("   ---")

	// re_domestic_service_areas
	reDomSerAreas := []models.ReDomesticServiceArea{}
	if err := db.All(&reDomSerAreas); err != nil {
		return err
	}
	length = len(reDomSerAreas)
	log.Printf("   re_domestic_service_areas (ReDomesticServiceArea): %d\n", length)
	log.Printf("     first: %+v\n", reDomSerAreas[0])
	log.Printf("      last: %+v\n", reDomSerAreas[length-1])
	log.Println("   ---")

	// re_rate_areas
	reRateAreas := []models.ReRateArea{}
	if err := db.All(&reRateAreas); err != nil {
		return err
	}
	length = len(reRateAreas)
	log.Printf("   re_rate_areas (ReRateArea): %d\n", length)
	log.Printf("     first: %+v\n", reRateAreas[0])
	log.Printf("      last: %+v\n", reRateAreas[length-1])
	log.Println("   ---")

	// re_domestic_linehaul_prices
	reDomLinePrices := []models.ReDomesticLinehaulPrice{}
	if err := db.All(&reDomLinePrices); err != nil {
		return err
	}
	length = len(reDomLinePrices)
	log.Printf("   reDomLinePrices (ReDomesticLinehaulPrice): %d\n", length)
	log.Printf("     first: %+v\n", reDomLinePrices[0])
	log.Printf("      last: %+v\n", reDomLinePrices[length-1])
	log.Println("   ---")

	// re_domestic_service_area_prices
	reDomSerAreaPrices := []models.ReDomesticServiceAreaPrice{}
	if err := db.All(&reDomSerAreaPrices); err != nil {
		return err
	}
	length = len(reDomSerAreaPrices)
	log.Printf("   re_domestic_service_area_prices (ReDomesticServiceAreaPrice): %d\n", length)
	log.Printf("     first: %+v\n", reDomSerAreaPrices[0])
	log.Printf("      last: %+v\n", reDomSerAreaPrices[length-1])
	log.Println("   ---")

	// re_domestic_other_prices
	reDomOtherPrices := []models.ReDomesticOtherPrice{}
	if err := db.All(&reDomOtherPrices); err != nil {
		return err
	}
	length = len(reDomOtherPrices)
	log.Printf("   re_domestic_other_prices (ReDomesticOtherPrice): %d\n", length)
	log.Printf("     first: %+v\n", reDomOtherPrices[0])
	log.Printf("      last: %+v\n", reDomOtherPrices[length-1])
	log.Println("   ---")

	// re_international_prices
	reIntlPrices := []models.ReIntlPrice{}
	if err := db.All(&reIntlPrices); err != nil {
		return err
	}
	length = len(reIntlPrices)
	log.Printf("   re_international_prices (ReIntlPrice): %d\n", length)
	log.Printf("     first: %+v\n", reIntlPrices[0])
	log.Printf("      last: %+v\n", reIntlPrices[length-1])
	log.Println("   ---")

	// re_international_other_prices
	reIntlOtherPrices := []models.ReIntlOtherPrice{}
	if err := db.All(&reIntlOtherPrices); err != nil {
		return err
	}
	length = len(reIntlOtherPrices)
	log.Printf("   re_international_other_prices (ReIntlOtherPrice): %d\n", length)
	log.Printf("     first: %+v\n", reIntlOtherPrices[0])
	log.Printf("      last: %+v\n", reIntlOtherPrices[length-1])
	log.Println("   ---")

	// re_task_order_fees
	reTaskOrderFees := []models.ReTaskOrderFee{}
	if err := db.All(&reTaskOrderFees); err != nil {
		return err
	}
	length = len(reTaskOrderFees)
	log.Printf("   re_task_order_fees (ReTaskOrderFee): %d\n", length)
	log.Printf("     first: %+v\n", reTaskOrderFees[0])
	log.Printf("      last: %+v\n", reTaskOrderFees[length-1])
	log.Println("   ---")

	// re_domestic_accessorial_prices
	reDomAccPrices := []models.ReDomesticAccessorialPrice{}
	if err := db.All(&reDomAccPrices); err != nil {
		return err
	}
	length = len(reDomAccPrices)
	log.Printf("   re_domestic_accessorial_prices (ReDomesticAccessorialPrice): %d\n", length)
	log.Printf("     first: %+v\n", reDomAccPrices[0])
	log.Printf("      last: %+v\n", reDomAccPrices[length-1])
	log.Println("   ---")

	// re_intl_accessorial_prices
	reIntlAccPrices := []models.ReIntlAccessorialPrice{}
	if err := db.All(&reIntlAccPrices); err != nil {
		return err
	}
	length = len(reIntlAccPrices)
	log.Printf("   re_intl_accessorial_prices (ReIntlAccessorialPrice): %d\n", length)
	log.Printf("     first: %+v\n", reIntlAccPrices[0])
	log.Printf("      last: %+v\n", reIntlAccPrices[length-1])
	log.Println("   ---")

	// re_shipment_type_prices
	reShipmentTypePrices := []models.ReShipmentTypePrice{}
	if err := db.All(&reShipmentTypePrices); err != nil {
		return err
	}
	length = len(reShipmentTypePrices)
	log.Printf("   re_shipment_type_prices (ReShipmentTypePrice): %d\n", length)
	log.Printf("     first: %+v\n", reShipmentTypePrices[0])
	log.Printf("      last: %+v\n", reShipmentTypePrices[length-1])
	log.Println("   ---")

	return nil
}

func xlsxSheetsUsage(xlsxDataSheets []pricing.XlsxDataSheetInfo) string {
	message := "Provide comma separated string of sequential sheet index numbers starting with 0:\n"
	message += "\t e.g. '-xlsxSheets=\"6,7,11\"'\n"
	message += "\t      '-xlsxSheets=\"6\"'\n"
	message += "\n"
	message += "Available sheets for parsing are: \n"

	for i, v := range xlsxDataSheets {
		if len(v.ProcessMethods) > 0 {
			description := ""
			if v.Description != nil {
				description = *v.Description
			}
			message += fmt.Sprintf("%d:  %s\n", i, description)
		}
	}

	message += "\n"
	message += "NOTE: This option disables the Rate Engine table import by disabling the --re-import flag\n"

	return message
}
