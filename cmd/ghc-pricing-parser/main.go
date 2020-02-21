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

func summarizeXlsxStageParsing(db *pop.Connection) error {
	log.Println("XLSX to Stage Table Parsing Complete")
	log.Println(" Summary:")

	// 1b Service Areas
	stageDomServiceAreas := []models.StageDomesticServiceArea{}
	err := db.Limit(2).All(&stageDomServiceAreas)
	if err != nil {
		return err
	}
	length, err := db.Count(models.StageDomesticServiceArea{})
	if err != nil {
		return err
	}

	log.Printf("   1b: Service Areas (StageDomesticServiceArea): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageDomServiceAreas[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageDomServiceAreas[1])
	}
	log.Println("   ---")

	stageIntlServiceArea := []models.StageInternationalServiceArea{}
	err = db.Limit(2).All(&stageIntlServiceArea)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageInternationalServiceArea{})
	if err != nil {
		return err
	}

	log.Printf("   1b: Service Areas (StageInternationalServiceArea): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageIntlServiceArea[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageIntlServiceArea[1])
	}
	log.Println("   ---")

	// 2a Domestic Linehaul Prices
	stageDomLinePrice := []models.StageDomesticLinehaulPrice{}
	err = db.Limit(2).All(&stageDomLinePrice)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageDomesticLinehaulPrice{})
	if err != nil {
		return err
	}

	log.Printf("   2a: Domestic Linehaul Prices (StageDomesticLinehaulPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageDomLinePrice[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageDomLinePrice[1])
	}
	log.Println("   ---")

	// 2b Domestic Service Area Prices
	stageDomSerAreaPrice := []models.StageDomesticServiceAreaPrice{}
	err = db.Limit(2).All(&stageDomSerAreaPrice)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageDomesticServiceAreaPrice{})
	if err != nil {
		return err
	}
	log.Printf("   2b: Domestic Service Area Prices (StageDomesticServiceAreaPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageDomSerAreaPrice[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageDomSerAreaPrice[1])
	}
	log.Println("   ---")

	// 2c Other Domestic Prices
	stageDomOtherPackPrice := []models.StageDomesticOtherPackPrice{}
	err = db.Limit(2).All(&stageDomOtherPackPrice)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageDomesticOtherPackPrice{})
	if err != nil {
		return err
	}
	log.Printf("   2c: Other Domestic Prices (StageDomesticOtherPackPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageDomOtherPackPrice[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageDomOtherPackPrice[1])
	}
	log.Println("   ---")

	stageDomOtherSitPrice := []models.StageDomesticOtherSitPrice{}
	err = db.Limit(2).All(&stageDomOtherSitPrice)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageDomesticOtherSitPrice{})
	if err != nil {
		return err
	}
	log.Printf("   2c: Other Domestic Prices (StageDomesticOtherSitPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageDomOtherSitPrice[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageDomOtherSitPrice[1])
	}
	log.Println("   ---")

	// 3a OCONUS to OCONUS Prices
	stageOconusToOconus := []models.StageOconusToOconusPrice{}
	err = db.Limit(2).All(&stageOconusToOconus)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageOconusToOconusPrice{})
	if err != nil {
		return err
	}
	log.Printf("   3a: OCONUS to OCONUS Prices (StageOconusToOconusPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageOconusToOconus[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageOconusToOconus[1])
	}
	log.Println("   ---")

	// 3b CONUS to OCONUS Prices
	stageConusToOconus := []models.StageConusToOconusPrice{}
	err = db.Limit(2).All(&stageConusToOconus)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageConusToOconusPrice{})
	if err != nil {
		return err
	}
	log.Printf("   3b: CONUS to OCONUS Prices (StageConusToOconusPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageConusToOconus[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageConusToOconus[1])
	}
	log.Println("   ---")

	// 3c OCONUS to CONUS Prices
	stageOconusToConus := []models.StageOconusToConusPrice{}
	err = db.Limit(2).All(&stageOconusToConus)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageOconusToConusPrice{})
	if err != nil {
		return err
	}
	log.Printf("   3c: OCONUS to CONUS Prices (StageOconusToConusPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageOconusToConus[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageOconusToConus[1])
	}
	log.Println("   ---")

	// 3d Other International Prices
	stageOtherIntlPrices := []models.StageOtherIntlPrice{}
	err = db.Limit(2).All(&stageOtherIntlPrices)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageOtherIntlPrice{})
	if err != nil {
		return err
	}

	log.Printf("   3d: Other International Prices (StageOtherIntlPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageOtherIntlPrices[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageOtherIntlPrices[1])
	}
	log.Println("   ---")

	// 3e Non-Standard Location Prices
	stageNonStdLocaPrices := []models.StageNonStandardLocnPrice{}
	err = db.Limit(2).All(&stageNonStdLocaPrices)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageNonStandardLocnPrice{})
	if err != nil {
		return err
	}
	log.Printf("   3e: Non-Standard Location Prices (StageNonStandardLocnPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageNonStdLocaPrices[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageNonStdLocaPrices[1])
	}
	log.Println("   ---")

	// 4a Management, Counseling, and Transition Prices
	stageMgmt := []models.StageShipmentManagementServicesPrice{}
	err = db.Limit(2).All(&stageMgmt)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageShipmentManagementServicesPrice{})
	if err != nil {
		return err
	}
	log.Printf("   4a: Management, Counseling, and Transition Prices (StageShipmentManagementServicesPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageMgmt[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageMgmt[1])
	}
	log.Println("   ---")

	stageCounsel := []models.StageCounselingServicesPrice{}
	err = db.Limit(2).All(&stageCounsel)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageCounselingServicesPrice{})
	if err != nil {
		return err
	}
	log.Printf("   4a: Management, Counseling, and Transition Prices (StageCounselingServicesPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageCounsel[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageCounsel[1])
	}
	log.Println("   ---")

	stageTransition := []models.StageTransitionPrice{}
	err = db.Limit(2).All(&stageTransition)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageTransitionPrice{})
	if err != nil {
		return err
	}
	log.Printf("   4a: Management, Counseling, and Transition Prices (StageTransitionPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageTransition[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageTransition[1])
	}
	log.Println("   ---")

	// 5a Accessorial and Additional Prices
	stageDomMoveAccess := []models.StageDomesticMoveAccessorialPrices{}
	err = db.Limit(2).All(&stageDomMoveAccess)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageDomesticMoveAccessorialPrices{})
	if err != nil {
		return err
	}
	log.Printf("   5a Accessorial and Additional Prices (StageDomesticMoveAccessorialPrices): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageDomMoveAccess[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageDomMoveAccess[1])
	}
	log.Println("   ---")

	stageIntlMoveAccess := []models.StageInternationalMoveAccessorialPrices{}
	err = db.Limit(2).All(&stageIntlMoveAccess)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageInternationalMoveAccessorialPrices{})
	if err != nil {
		return err
	}
	log.Printf("   5a Accessorial and Additional Prices (StageInternationalMoveAccessorialPrices): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageIntlMoveAccess[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageIntlMoveAccess[1])
	}
	log.Println("   ---")

	stageDomIntlAdd := []models.StageDomesticInternationalAdditionalPrices{}
	err = db.Limit(2).All(&stageDomIntlAdd)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StageDomesticInternationalAdditionalPrices{})
	if err != nil {
		return err
	}
	log.Printf("   5a Accessorial and Additional Prices (StageDomesticInternationalAdditionalPrices): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stageDomIntlAdd[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stageDomIntlAdd[1])
	}
	log.Println("   ---")

	// 5b Price Escalation Discount
	stagePriceEsc := []models.StagePriceEscalationDiscount{}
	err = db.Limit(2).All(&stagePriceEsc)
	if err != nil {
		return err
	}
	length, err = db.Count(models.StagePriceEscalationDiscount{})
	if err != nil {
		return err
	}
	log.Printf("   5b: Price Escalation Discount (StagePriceEscalationDiscount): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", stagePriceEsc[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", stagePriceEsc[1])
	}
	log.Println("   ---")

	return nil
}

func summarizeStageReImport(db *pop.Connection, contractID uuid.UUID) error {
	log.Println("Stage Table import into Rate Engine Tables Complete")
	log.Println(" Summary:")

	// re_contract
	reContract := []models.ReContract{}
	err := db.Where("id = ?", contractID).Limit(2).All(&reContract)
	if err != nil {
		return err
	}
	length, err := db.Where("id = ?", contractID).Count(models.ReContract{})
	if err != nil {
		return err
	}

	log.Printf("   re_contract (ReContract): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reContract[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reContract[length-1])
	}
	log.Println("   ---")

	// re_contract_years
	reContractYears := []models.ReContractYear{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reContractYears)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReContractYear{})
	if err != nil {
		return err
	}

	log.Printf("   re_contract_years (ReContractYear): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reContractYears[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reContractYears[1])
	}
	log.Println("   ---")

	// re_domestic_service_areas
	reDomSerAreas := []models.ReDomesticServiceArea{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reDomSerAreas)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReDomesticServiceArea{})
	if err != nil {
		return err
	}

	log.Printf("   re_domestic_service_areas (ReDomesticServiceArea): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reDomSerAreas[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reDomSerAreas[1])
	}
	log.Println("   ---")

	// re_rate_areas
	reRateAreas := []models.ReRateArea{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reRateAreas)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReRateArea{})
	if err != nil {
		return err
	}
	log.Printf("   re_rate_areas (ReRateArea): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reRateAreas[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reRateAreas[1])
	}
	log.Println("   ---")

	// re_domestic_linehaul_prices
	reDomLinePrices := []models.ReDomesticLinehaulPrice{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reDomLinePrices)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReDomesticLinehaulPrice{})
	if err != nil {
		return err
	}
	log.Printf("   reDomLinePrices (ReDomesticLinehaulPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reDomLinePrices[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reDomLinePrices[1])
	}
	log.Println("   ---")

	// re_domestic_service_area_prices
	reDomSerAreaPrices := []models.ReDomesticServiceAreaPrice{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reDomSerAreaPrices)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReDomesticServiceAreaPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_domestic_service_area_prices (ReDomesticServiceAreaPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reDomSerAreaPrices[0])
	}
	log.Printf("\tsecond: %+v\n", reDomSerAreaPrices[1])
	log.Println("   ---")

	// re_domestic_other_prices
	reDomOtherPrices := []models.ReDomesticOtherPrice{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reDomOtherPrices)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReDomesticOtherPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_domestic_other_prices (ReDomesticOtherPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reDomOtherPrices[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reDomOtherPrices[1])
	}
	log.Println("   ---")

	// re_international_prices
	reIntlPrices := []models.ReIntlPrice{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reIntlPrices)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReIntlPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_international_prices (ReIntlPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reIntlPrices[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reIntlPrices[1])
	}
	log.Println("   ---")

	// re_international_other_prices
	reIntlOtherPrices := []models.ReIntlOtherPrice{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reIntlOtherPrices)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReIntlOtherPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_international_other_prices (ReIntlOtherPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reIntlOtherPrices[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reIntlOtherPrices[1])
	}
	log.Println("   ---")

	// re_task_order_fees
	//possibly need a join where contract year id  = contract_year.contract_id
	reTaskOrderFees := []models.ReTaskOrderFee{}
	err = db.Where("contract_id = ?", contractID).Join("re_contract_years", "re_contract_years.id = contract_year_id").Limit(2).All(&reTaskOrderFees)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Join("re_contract_years", "re_contract_years.id = contract_year_id").Count(models.ReTaskOrderFee{})
	if err != nil {
		return err
	}
	log.Printf("   re_task_order_fees (ReTaskOrderFee): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reTaskOrderFees[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reTaskOrderFees[1])
	}
	log.Println("   ---")

	// re_domestic_accessorial_prices
	reDomAccPrices := []models.ReDomesticAccessorialPrice{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reDomAccPrices)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReDomesticAccessorialPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_domestic_accessorial_prices (ReDomesticAccessorialPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reDomAccPrices[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reDomAccPrices[1])
	}
	log.Println("   ---")

	// re_intl_accessorial_prices
	reIntlAccPrices := []models.ReIntlAccessorialPrice{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reIntlAccPrices)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReIntlAccessorialPrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_intl_accessorial_prices (ReIntlAccessorialPrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reIntlAccPrices[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reIntlAccPrices[1])
	}
	log.Println("   ---")

	// re_shipment_type_prices
	reShipmentTypePrices := []models.ReShipmentTypePrice{}
	err = db.Where("contract_id = ?", contractID).Limit(2).All(&reShipmentTypePrices)
	if err != nil {
		return err
	}
	length, err = db.Where("contract_id = ?", contractID).Count(models.ReShipmentTypePrice{})
	if err != nil {
		return err
	}
	log.Printf("   re_shipment_type_prices (ReShipmentTypePrice): %d\n", length)
	if length > 0 {
		log.Printf("\tfirst: %+v\n", reShipmentTypePrices[0])
	}
	if length > 1 {
		log.Printf("\tsecond: %+v\n", reShipmentTypePrices[1])
	}
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
