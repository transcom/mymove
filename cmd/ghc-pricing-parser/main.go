package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
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
	if err = summarizeXlsxStageParsing(db); err != nil {
		log.Fatalf("Failed to summarize XLSX to stage table parsing: %v", err)
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
		// Summarize import for verification
		if err := summarizeStageReImport(db, ghcREImporter.ContractID); err != nil {
			log.Fatalf("Failed to summarize stage table to rate engine table import: %v", err)
		}
	}
}

func secondIndex(length int) int {
	if length > 1 {
		return 1
	}
	return 0
}

func summarizeXlsxStageParsing(db *pop.Connection) error {
	log.Println("XLSX to Stage Table Parsing Complete")
	log.Println(" Summary:")

	// 1b Service Areas
	stageDomServiceAreas := []models.StageDomesticServiceArea{}
	db.Limit(2).All(&stageDomServiceAreas)

	length, err := db.Count(models.StageDomesticServiceArea{})
	if err != nil {
		return err
	}

	log.Printf("   1b: Service Areas (StageDomesticServiceArea): %d\n", length)
	log.Printf("     first: %+v\n", stageDomServiceAreas[0])
	log.Printf("      second: %+v\n", stageDomServiceAreas[secondIndex(length)])
	log.Println("   ---")

	stageIntlServiceArea := []models.StageInternationalServiceArea{}
	db.Limit(2).All(&stageIntlServiceArea)

	length, err = db.Count(models.StageInternationalServiceArea{})
	if err != nil {
		return err
	}

	log.Printf("   1b: Service Areas (StageInternationalServiceArea): %d\n", length)
	log.Printf("     first: %+v\n", stageIntlServiceArea[0])
	log.Printf("      second: %+v\n", stageIntlServiceArea[secondIndex(length)])
	log.Println("   ---")

	// 2a Domestic Linehaul Prices
	stageDomLinePrice := []models.StageDomesticLinehaulPrice{}
	db.Limit(2).All(&stageDomLinePrice)

	length, err = db.Count(models.StageDomesticLinehaulPrice{})
	if err != nil {
		return err
	}
	log.Printf("   2a: Domestic Linehaul Prices (StageDomesticLinehaulPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageDomLinePrice[0])
	log.Printf("      second: %+v\n", stageDomLinePrice[secondIndex(length)])
	log.Println("   ---")

	// 2b Domestic Service Area Prices
	stageDomSerAreaPrice := []models.StageDomesticServiceAreaPrice{}
	db.Limit(2).All(&stageDomSerAreaPrice)
	length, err = db.Count(models.StageDomesticServiceAreaPrice{})
	if err != nil {
		return err
	}
	log.Printf("   2b: Domestic Service Area Prices (StageDomesticServiceAreaPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageDomSerAreaPrice[0])
	log.Printf("      second: %+v\n", stageDomSerAreaPrice[secondIndex(length)])
	log.Println("   ---")

	// 2c Other Domestic Prices
	stageDomOtherPackPrice := []models.StageDomesticOtherPackPrice{}
	db.Limit(2).All(&stageDomOtherPackPrice)
	length, err = db.Count(models.StageDomesticOtherPackPrice{})
	if err != nil {
		return err
	}
	log.Printf("   2c: Other Domestic Prices (StageDomesticOtherPackPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageDomOtherPackPrice[0])
	log.Printf("      second: %+v\n", stageDomOtherPackPrice[secondIndex(length)])
	log.Println("   ---")

	stageDomOtherSitPrice := []models.StageDomesticOtherSitPrice{}
	db.Limit(2).All(&stageDomOtherSitPrice)
	// length = len(stageDomOtherSitPrice)
	length, err = db.Count(models.StageDomesticOtherSitPrice{})
	if err != nil {
		return err
	}
	log.Printf("   2c: Other Domestic Prices (StageDomesticOtherSitPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageDomOtherSitPrice[0])
	log.Printf("      second: %+v\n", stageDomOtherSitPrice[secondIndex(length)])
	log.Println("   ---")

	// 3a OCONUS to OCONUS Prices
	stageOconusToOconus := []models.StageOconusToOconusPrice{}
	db.Limit(2).All(&stageOconusToOconus)
	// length = len(stageOconusToOconus)
	length, err = db.Count(models.StageOconusToOconusPrice{})
	if err != nil {
		return err
	}
	log.Printf("   3a: OCONUS to OCONUS Prices (StageOconusToOconusPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageOconusToOconus[0])
	log.Printf("      second: %+v\n", stageOconusToOconus[secondIndex(length)])
	log.Println("   ---")

	// 3b CONUS to OCONUS Prices
	stageConusToOconus := []models.StageConusToOconusPrice{}
	db.Limit(2).All(&stageConusToOconus)
	// length = len(stageConusToOconus)
	length, err = db.Count(models.StageConusToOconusPrice{})
	if err != nil {
		return err
	}
	log.Printf("   3b: CONUS to OCONUS Prices (StageConusToOconusPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageConusToOconus[0])
	log.Printf("      second: %+v\n", stageConusToOconus[secondIndex(length)])
	log.Println("   ---")

	// 3c OCONUS to CONUS Prices
	stageOconusToConus := []models.StageOconusToConusPrice{}
	db.Limit(2).All(&stageOconusToConus)
	// length = len(stageOconusToConus)
	length, err = db.Count(models.StageOconusToConusPrice{})
	if err != nil {
		return err
	}
	log.Printf("   3c: OCONUS to CONUS Prices (StageOconusToConusPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageOconusToConus[0])
	log.Printf("      second: %+v\n", stageOconusToConus[secondIndex(length)])
	log.Println("   ---")

	// 3d Other International Prices
	stageOtherIntlPrices := []models.StageOtherIntlPrice{}
	db.Limit(2).All(&stageOtherIntlPrices)
	length, err = db.Count(models.StageOtherIntlPrice{})
	if err != nil {
		return err
	}

	log.Printf("   3d: Other International Prices (StageOtherIntlPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageOtherIntlPrices[0])
	log.Printf("      second: %+v\n", stageOtherIntlPrices[secondIndex(length)])
	log.Println("   ---")

	// 3e Non-Standard Location Prices
	stageNonStdLocaPrices := []models.StageNonStandardLocnPrice{}
	db.Limit(2).All(&stageNonStdLocaPrices)
	length, err = db.Count(models.StageNonStandardLocnPrice{})
	if err != nil {
		return err
	}
	log.Printf("   3e: Non-Standard Location Prices (StageNonStandardLocnPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageNonStdLocaPrices[0])
	log.Printf("      second: %+v\n", stageNonStdLocaPrices[secondIndex(length)])
	log.Println("   ---")

	// 4a Management, Counseling, and Transition Prices
	stageMgmt := []models.StageShipmentManagementServicesPrice{}
	db.Limit(2).All(&stageMgmt)
	length, err = db.Count(models.StageShipmentManagementServicesPrice{})
	if err != nil {
		return err
	}
	log.Printf("   4a: Management, Counseling, and Transition Prices (StageShipmentManagementServicesPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageMgmt[0])
	log.Printf("      second: %+v\n", stageMgmt[secondIndex(length)])
	log.Println("   ---")

	stageCounsel := []models.StageCounselingServicesPrice{}
	db.Limit(2).All(&stageCounsel)
	length, err = db.Count(models.StageCounselingServicesPrice{})
	if err != nil {
		return err
	}
	log.Printf("   4a: Management, Counseling, and Transition Prices (StageCounselingServicesPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageCounsel[0])
	log.Printf("      second: %+v\n", stageCounsel[secondIndex(length)])
	log.Println("   ---")

	stageTransition := []models.StageTransitionPrice{}
	db.Limit(2).All(&stageTransition)
	length, err = db.Count(models.StageTransitionPrice{})
	if err != nil {
		return err
	}
	log.Printf("   4a: Management, Counseling, and Transition Prices (StageTransitionPrice): %d\n", length)
	log.Printf("     first: %+v\n", stageTransition[0])
	log.Printf("      second: %+v\n", stageTransition[secondIndex(length)])
	log.Println("   ---")

	// 5a Accessorial and Additional Prices
	stageDomMoveAccess := []models.StageDomesticMoveAccessorialPrices{}
	db.Limit(2).All(&stageDomMoveAccess)
	length, err = db.Count(models.StageDomesticMoveAccessorialPrices{})
	if err != nil {
		return err
	}
	log.Printf("   5a Accessorial and Additional Prices (StageDomesticMoveAccessorialPrices): %d\n", length)
	log.Printf("     first: %+v\n", stageDomMoveAccess[0])
	log.Printf("      second: %+v\n", stageDomMoveAccess[secondIndex(length)])
	log.Println("   ---")

	stageIntlMoveAccess := []models.StageInternationalMoveAccessorialPrices{}
	db.Limit(2).All(&stageIntlMoveAccess)
	length, err = db.Count(models.StageInternationalMoveAccessorialPrices{})
	if err != nil {
		return err
	}
	log.Printf("   5a Accessorial and Additional Prices (StageInternationalMoveAccessorialPrices): %d\n", length)
	log.Printf("     first: %+v\n", stageIntlMoveAccess[0])
	log.Printf("      second: %+v\n", stageIntlMoveAccess[secondIndex(length)])
	log.Println("   ---")

	stageDomIntlAdd := []models.StageDomesticInternationalAdditionalPrices{}
	db.Limit(2).All(&stageDomIntlAdd)
	length, err = db.Count(models.StageDomesticInternationalAdditionalPrices{})
	if err != nil {
		return err
	}
	log.Printf("   5a Accessorial and Additional Prices (StageDomesticInternationalAdditionalPrices): %d\n", length)
	log.Printf("     first: %+v\n", stageDomIntlAdd[0])
	log.Printf("      second: %+v\n", stageDomIntlAdd[secondIndex(length)])
	log.Println("   ---")

	// 5b Price Escalation Discount
	stagePriceEsc := []models.StagePriceEscalationDiscount{}
	db.Limit(2).All(&stagePriceEsc)
	length, err = db.Count(models.StagePriceEscalationDiscount{})
	if err != nil {
		return err
	}
	log.Printf("   5b: Price Escalation Discount (StagePriceEscalationDiscount): %d\n", length)
	log.Printf("     first: %+v\n", stagePriceEsc[0])
	log.Printf("      second: %+v\n", stagePriceEsc[secondIndex(length)])
	log.Println("   ---")

	return nil
}

func summarizeStageReImport(db *pop.Connection, contractID uuid.UUID) error {
	log.Println("Stage Table import into Rate Engine Tables Complete")
	log.Println(" Summary:")

	// re_contract
	reContract := []models.ReContract{}
	db.Where("id = ?", contractID).Limit(2).All(&reContract)

	length, err := db.Where("id = ?", contractID).Count(models.ReContract{})
	if err != nil {
		return err
	}

	log.Printf("   re_contract (ReContract): %d\n", length)
	log.Printf("     first: %s\n", reContract[0])
	log.Printf("      second: %s\n", reContract[length-1])
	log.Println("   ---")

	// re_contract_years
	reContractYears := []models.ReContractYear{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reContractYears)

	length, err = db.Where("contract_id = ?", contractID).Count(models.ReContractYear{})
	if err != nil {
		return err
	}
	log.Printf("   re_contract_years (ReContractYear): %d\n", length)
	log.Printf("     first: %+v\n", reContractYears[0])
	log.Printf("      second: %+v\n", reContractYears[secondIndex(length)])
	log.Println("   ---")

	// re_domestic_service_areas
	reDomSerAreas := []models.ReDomesticServiceArea{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reDomSerAreas)

	length, err = db.Where("contract_id = ?", contractID).Count(models.ReDomesticServiceArea{})
	if err != nil {
		return err
	}

	log.Printf("   re_domestic_service_areas (ReDomesticServiceArea): %d\n", length)
	log.Printf("     first: %+v\n", reDomSerAreas[0])
	log.Printf("      second: %+v\n", reDomSerAreas[secondIndex(length)])
	log.Println("   ---")

	// re_rate_areas
	reRateAreas := []models.ReRateArea{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reRateAreas)
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReRateArea{})
	if err != nil {
		return err
	}
	log.Printf("   re_rate_areas (ReRateArea): %d\n", length)
	log.Printf("     first: %+v\n", reRateAreas[0])
	log.Printf("      second: %+v\n", reRateAreas[secondIndex(length)])
	log.Println("   ---")

	// re_domestic_linehaul_prices
	reDomLinePrices := []models.ReDomesticLinehaulPrice{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reDomLinePrices)
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReDomesticLinehaulPrice{})
	if err != nil {
		return err
	}
	log.Printf("   reDomLinePrices (ReDomesticLinehaulPrice): %d\n", length)
	log.Printf("     first: %+v\n", reDomLinePrices[0])
	log.Printf("      second: %+v\n", reDomLinePrices[secondIndex(length)])
	log.Println("   ---")

	// re_domestic_service_area_prices
	reDomSerAreaPrices := []models.ReDomesticServiceAreaPrice{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reDomSerAreaPrices)
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReDomesticServiceAreaPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_domestic_service_area_prices (ReDomesticServiceAreaPrice): %d\n", length)
	log.Printf("     first: %+v\n", reDomSerAreaPrices[0])
	log.Printf("      second: %+v\n", reDomSerAreaPrices[secondIndex(length)])
	log.Println("   ---")

	// re_domestic_other_prices
	reDomOtherPrices := []models.ReDomesticOtherPrice{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reDomOtherPrices)
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReDomesticOtherPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_domestic_other_prices (ReDomesticOtherPrice): %d\n", length)
	log.Printf("     first: %+v\n", reDomOtherPrices[0])
	log.Printf("      second: %+v\n", reDomOtherPrices[secondIndex(length)])
	log.Println("   ---")

	// re_international_prices
	reIntlPrices := []models.ReIntlPrice{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reIntlPrices)
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReIntlPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_international_prices (ReIntlPrice): %d\n", length)
	log.Printf("     first: %+v\n", reIntlPrices[0])
	log.Printf("      second: %+v\n", reIntlPrices[secondIndex(length)])
	log.Println("   ---")

	// re_international_other_prices
	reIntlOtherPrices := []models.ReIntlOtherPrice{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reIntlOtherPrices)
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReIntlOtherPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_international_other_prices (ReIntlOtherPrice): %d\n", length)
	log.Printf("     first: %+v\n", reIntlOtherPrices[0])
	log.Printf("      second: %+v\n", reIntlOtherPrices[secondIndex(length)])
	log.Println("   ---")

	// re_task_order_fees
	reTaskOrderFees := []models.ReTaskOrderFee{}
	db.Limit(2).All(&reTaskOrderFees)
	length, err = db.Count(models.ReTaskOrderFee{})
	if err != nil {
		return err
	}
	log.Printf("   re_task_order_fees (ReTaskOrderFee): %d\n", length)
	log.Printf("     first: %+v\n", reTaskOrderFees[0])
	log.Printf("      second: %+v\n", reTaskOrderFees[secondIndex(length)])
	log.Println("   ---")

	// re_domestic_accessorial_prices
	reDomAccPrices := []models.ReDomesticAccessorialPrice{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reDomAccPrices)
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReDomesticAccessorialPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_domestic_accessorial_prices (ReDomesticAccessorialPrice): %d\n", length)
	log.Printf("     first: %+v\n", reDomAccPrices[0])
	log.Printf("      second: %+v\n", reDomAccPrices[secondIndex(length)])
	log.Println("   ---")

	// re_intl_accessorial_prices
	reIntlAccPrices := []models.ReIntlAccessorialPrice{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reIntlAccPrices)
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReIntlAccessorialPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_intl_accessorial_prices (ReIntlAccessorialPrice): %d\n", length)
	log.Printf("     first: %+v\n", reIntlAccPrices[0])
	log.Printf("      second: %+v\n", reIntlAccPrices[secondIndex(length)])
	log.Println("   ---")

	// re_shipment_type_prices
	reShipmentTypePrices := []models.ReShipmentTypePrice{}
	db.Where("contract_id = ?", contractID).Limit(2).All(&reShipmentTypePrices)
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReShipmentTypePrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_shipment_type_prices (ReShipmentTypePrice): %d\n", length)
	log.Printf("     first: %+v\n", reShipmentTypePrices[0])
	log.Printf("      second: %+v\n", reShipmentTypePrices[secondIndex(length)])
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
