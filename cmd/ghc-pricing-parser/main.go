package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
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
	// Set up spreadsheet metadata and parameter configuration
	xlsxDataSheets := pricing.InitDataSheetInfo()
	params := pricing.ParamConfig{}
	params.RunTime = time.Now()

	// Set up parser's command line flags
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

	// Set up DB flags
	cli.InitDatabaseFlags(flag)

	// Parse flags
	flag.SortFlags = false
	err := flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("Could not parse flags: %v\n", err)
	}

	// Bind flags
	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		log.Fatalf("Could not bind flags: %v\n", err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Create logger
	logger, err := logging.Config(v.GetString(cli.DbEnvFlag), v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	// Connect to the database
	err = cli.CheckDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}
	db, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Fatal("Could not close database", zap.Error(closeErr))
		}
	}()

	// Ensure we've been given a spreadsheet to parse
	if params.XlsxFilename == "" {
		logger.Fatal("Did not receive an XLSX filename to parse; missing --filename")
	}

	// Running with a subset of worksheets will turn off ProcessAll flag and the rate engine import
	if len(params.XlsxSheets) > 0 {
		params.ProcessAll = false
		logger.Info("Setting --xlsxSheets disables --re-import so no data will be imported into the rate engine tables. Only stage table data will be updated.")
		params.RunImport = false
	}

	// If we are importing into the rate engine tables, we need a contract code
	if params.RunImport && params.ContractCode == "" {
		logger.Fatal("Did not receive a contract code; missing --contract-code")
	}

	// Open the spreadsheet
	logger.Info("Importing file", zap.String("XlsxFilename", params.XlsxFilename))
	params.XlsxFile, err = xlsx.OpenFile(params.XlsxFilename)
	if err != nil {
		logger.Fatal("Failed to open file", zap.String("XlsxFilename", params.XlsxFilename), zap.Error(err))
	}

	// Now kick off the parsing
	err = pricing.Parse(xlsxDataSheets, params, db, logger)
	if err != nil {
		logger.Fatal("Failed to parse pricing template", zap.Error(err))
	}
	if err = summarizeXlsxStageParsing(db, logger); err != nil {
		logger.Fatal("Failed to summarize XLSX to stage table parsing", zap.Error(err))
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
			logger.Fatal("GHC Rate Engine import failed", zap.Error(err))
		}
		if err := summarizeStageReImport(db, logger, ghcREImporter.ContractID); err != nil {
			logger.Fatal("Failed to summarize stage table to rate engine table import", zap.Error(err))
		}
	}
}

func summarizeXlsxStageParsing(db *pop.Connection, logger logger) error {
	logger.Info("XLSX to Stage Table Parsing Complete")
	logger.Info("Summary:")

	tables := []struct {
		header   string
		ptrSlice interface{}
	}{
		{"1b: Service Areas", &[]models.StageDomesticServiceArea{}},
		{"1b: Service Areas", &[]models.StageInternationalServiceArea{}},
		{"2a: Domestic Linehaul Prices", &[]models.StageDomesticLinehaulPrice{}},
		{"2b: Domestic Service Area Prices", &[]models.StageDomesticServiceAreaPrice{}},
		{"2c: Other Domestic Prices", &[]models.StageDomesticOtherPackPrice{}},
		{"2c: Other Domestic Prices", &[]models.StageDomesticOtherSitPrice{}},
		{"3a: OCONUS to OCONUS Prices", &[]models.StageOconusToOconusPrice{}},
		{"3b: CONUS to OCONUS Prices", &[]models.StageConusToOconusPrice{}},
		{"3c: OCONUS to CONUS Prices", &[]models.StageOconusToConusPrice{}},
		{"3d: Other International Prices", &[]models.StageOtherIntlPrice{}},
		{"3e: Non-Standard Location Prices", &[]models.StageNonStandardLocnPrice{}},
		{"4a: Management, Counseling, and Transition Prices", &[]models.StageShipmentManagementServicesPrice{}},
		{"4a: Management, Counseling, and Transition Prices", &[]models.StageCounselingServicesPrice{}},
		{"4a: Management, Counseling, and Transition Prices", &[]models.StageTransitionPrice{}},
		{"5a: Accessorial and Additional Prices", &[]models.StageDomesticMoveAccessorialPrice{}},
		{"5a: Accessorial and Additional Prices", &[]models.StageInternationalMoveAccessorialPrice{}},
		{"5a: Accessorial and Additional Prices", &[]models.StageDomesticInternationalAdditionalPrice{}},
		{"5b: Price Escalation Discount", &[]models.StagePriceEscalationDiscount{}},
	}

	for _, table := range tables {
		err := summarizeXlsxStageTable(db, logger, table.header, table.ptrSlice)
		if err != nil {
			return err
		}
	}

	return nil
}

func summarizeXlsxStageTable(db *pop.Connection, logger logger, header string, modelSlice interface{}) error {
	modelType := reflect.TypeOf(modelSlice).Elem().Elem()
	modelName := modelType.Name()
	modelInstance := reflect.New(modelType)

	err := db.Limit(2).All(modelSlice)
	if err != nil {
		return err
	}
	length, err := db.Count(modelInstance.Interface())
	if err != nil {
		return err
	}

	modelSliceValue := reflect.ValueOf(modelSlice).Elem()

	headerMsg := fmt.Sprintf("\t%s (%s)", header, modelName)
	logger.Info(headerMsg, zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any(modelName, modelSliceValue.Index(0).Interface()))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any(modelName, modelSliceValue.Index(1).Interface()))
	}
	logger.Info("\t---")

	return nil
}

func summarizeStageReImport(db *pop.Connection, logger logger, contractID uuid.UUID) error {
	logger.Info("Stage Table import into Rate Engine Tables Complete")
	logger.Info(" Summary:")

	// re_contract
	reContracts := []models.ReContract{}
	err := db.Where("id = ?", contractID).Limit(2).All(&reContracts)
	if err != nil {
		return err
	}
	length, err := db.Where("id = ?", contractID).Count(models.ReContract{})
	if err != nil {
		return err
	}

	logger.Info("\tre_contract (ReContract)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReContract", reContracts[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReContract", reContracts[1]))
	}
	logger.Info("\t---")

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

	logger.Info("\tre_contract_years (ReContractYear)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReContractYear", reContractYears[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReContractYear", reContractYears[1]))
	}
	logger.Info("\t---")

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

	logger.Info("\tre_domestic_service_areas (ReDomesticServiceArea)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReDomesticServiceArea", reDomSerAreas[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReDomesticServiceArea", reDomSerAreas[1]))
	}
	logger.Info("\t---")

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
	logger.Info("\tre_rate_areas (ReRateArea)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReRateArea", reRateAreas[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReRateArea", reRateAreas[1]))
	}
	logger.Info("\t---")

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
	logger.Info("\treDomLinePrices (ReDomesticLinehaulPrice)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReDomesticLinehaulPrice", reDomLinePrices[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReRReDomesticLinehaulPriceateArea", reDomLinePrices[1]))
	}
	logger.Info("\t---")

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
	logger.Info("\tre_domestic_service_area_prices (ReDomesticServiceAreaPrice)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReDomesticServiceAreaPrice", reDomSerAreaPrices[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReDomesticServiceAreaPrice", reDomSerAreaPrices[1]))
	}
	logger.Info("\t---")

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
	logger.Info("\tre_domestic_other_prices (ReDomesticOtherPrice)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReDomesticOtherPrice", reDomOtherPrices[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReDomesticOtherPrice", reDomOtherPrices[1]))
	}
	logger.Info("\t---")

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
	logger.Info("\tre_international_prices (ReIntlPrice)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReIntlPrice", reIntlPrices[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReIntlPrice", reIntlPrices[1]))
	}
	logger.Info("\t---")

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
	logger.Info("\tre_international_other_prices (ReIntlOtherPrice)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReIntlOtherPrice", reIntlOtherPrices[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReIntlOtherPrice", reIntlOtherPrices[1]))
	}
	logger.Info("\t---")

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
	logger.Info("\tre_task_order_fees (ReTaskOrderFee)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReTaskOrderFee", reTaskOrderFees[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReTaskOrderFee", reTaskOrderFees[1]))
	}
	logger.Info("\t---")

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
	logger.Info("\tre_domestic_accessorial_prices (ReDomesticAccessorialPrice)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReDomesticAccessorialPrice", reDomAccPrices[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReDomesticAccessorialPrice", reDomAccPrices[1]))
	}
	logger.Info("\t---")

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
	logger.Info("\tre_intl_accessorial_prices (ReIntlAccessorialPrice)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReIntlAccessorialPrice", reIntlAccPrices[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReIntlAccessorialPrice", reIntlAccPrices[1]))
	}
	logger.Info("\t---")

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
	logger.Info("\tre_shipment_type_prices (ReShipmentTypePrice)", zap.Int("length", length))
	if length > 0 {
		logger.Info("\t\tfirst", zap.Any("ReShipmentTypePrice", reShipmentTypePrices[0]))
	}
	if length > 1 {
		logger.Info("\t\tsecond", zap.Any("ReShipmentTypePrice", reShipmentTypePrices[1]))
	}
	logger.Info("\t---")

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
