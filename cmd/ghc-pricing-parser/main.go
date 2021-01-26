package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
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
	flag.StringVar(&params.XlsxFilename, "filename", "", "Filename (including path) of the XLSX to parse for the GHC rate engine data import")
	flag.StringVar(&params.ContractCode, "contract-code", "", "Contract code to use for this import")
	flag.StringVar(&params.ContractName, "contract-name", "", "Contract name to use for this import; if not provided, the contract-code value will be used")
	flag.StringVar(&params.ContractStartDate, "contract-start-date", "2021-02-01", "Beginning base date for contracts periods, in format: YYYY-MM-DD; if not provided, 2021-02-01 will be used")
	flag.BoolVar(&params.ProcessAll, "all", true, "Parse entire GHC Rate Engine XLSX")
	flag.StringSliceVar(&params.XlsxSheets, "xlsxSheets", []string{}, xlsxSheetsUsage(xlsxDataSheets))
	flag.BoolVar(&params.ShowOutput, "display", false, "Display output of parsed info")
	flag.BoolVar(&params.SaveToFile, "save-csv", false, "Save output of XLSX sheets to CSV file")
	flag.BoolVar(&params.RunVerify, "verify", true, "Perform sheet format verification -- but does not validate data")
	flag.BoolVar(&params.RunImport, "re-import", true, "Perform the import from staging tables to GHC rate engine tables")
	flag.BoolVar(&params.UseTempTables, "use-temp-tables", true, "Make the staging tables be temp tables that don't persist after import")
	flag.BoolVar(&params.DropIfExists, "drop", false, "Drop any existing staging tables prior to creating them; useful when turning `--use-temp-tables` off")

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
	logger, err := logging.Config(logging.WithEnvironment(v.GetString(cli.DbEnvFlag)), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
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

	// Before parsing spreadsheet, ensure there's a valid contract start date
	basePeriodStartDateForPrimeContract1, err := time.Parse("2006-01-02", params.ContractStartDate)
	if err != nil {
		logger.Fatal("could not parse the given contract start date", zap.Error(err))
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
			Logger:            logger,
			ContractCode:      params.ContractCode,
			ContractName:      params.ContractName,
			ContractStartDate: basePeriodStartDateForPrimeContract1,
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
	logger.Info("XLSX to stage table parsing complete. Summary follows:")
	logger.Info("====")

	models := []struct {
		header        string
		modelInstance interface{}
	}{
		{"1b: Service Areas", models.StageDomesticServiceArea{}},
		{"1b: Service Areas", models.StageInternationalServiceArea{}},
		{"2a: Domestic Linehaul Prices", models.StageDomesticLinehaulPrice{}},
		{"2b: Domestic Service Area Prices", models.StageDomesticServiceAreaPrice{}},
		{"2c: Other Domestic Prices", models.StageDomesticOtherPackPrice{}},
		{"2c: Other Domestic Prices", models.StageDomesticOtherSitPrice{}},
		{"3a: OCONUS to OCONUS Prices", models.StageOconusToOconusPrice{}},
		{"3b: CONUS to OCONUS Prices", models.StageConusToOconusPrice{}},
		{"3c: OCONUS to CONUS Prices", models.StageOconusToConusPrice{}},
		{"3d: Other International Prices", models.StageOtherIntlPrice{}},
		{"3e: Non-Standard Location Prices", models.StageNonStandardLocnPrice{}},
		{"4a: Management, Counseling, and Transition Prices", models.StageShipmentManagementServicesPrice{}},
		{"4a: Management, Counseling, and Transition Prices", models.StageCounselingServicesPrice{}},
		{"4a: Management, Counseling, and Transition Prices", models.StageTransitionPrice{}},
		{"5a: Accessorial and Additional Prices", models.StageDomesticMoveAccessorialPrice{}},
		{"5a: Accessorial and Additional Prices", models.StageInternationalMoveAccessorialPrice{}},
		{"5a: Accessorial and Additional Prices", models.StageDomesticInternationalAdditionalPrice{}},
		{"5b: Price Escalation Discount", models.StagePriceEscalationDiscount{}},
	}

	for index, model := range models {
		if index != 0 {
			logger.Info("----")
		}
		err := summarizeModel(db, logger, model.header, model.modelInstance, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func summarizeStageReImport(db *pop.Connection, logger logger, contractID uuid.UUID) error {
	logger.Info("Stage table import into rate engine tables complete. Summary follows:")
	logger.Info("====")

	models := []struct {
		header        string
		modelInstance interface{}
		filter        *pop.Query
	}{
		{
			"re_contract",
			models.ReContract{},
			db.Where("id = ?", contractID),
		},
		{
			"re_contract_years",
			models.ReContractYear{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_domestic_service_areas",
			models.ReDomesticServiceArea{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_zip3s",
			models.ReZip3{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_rate_areas",
			models.ReRateArea{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_domestic_linehaul_prices",
			models.ReDomesticLinehaulPrice{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_domestic_service_area_prices",
			models.ReDomesticServiceAreaPrice{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_domestic_other_prices",
			models.ReDomesticOtherPrice{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_intl_prices",
			models.ReIntlPrice{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_intl_other_prices",
			models.ReIntlOtherPrice{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_task_order_fees",
			models.ReTaskOrderFee{},
			db.Where("contract_id = ?", contractID).Join("re_contract_years", "re_contract_years.id = contract_year_id"),
		},
		{
			"re_domestic_accessorial_prices",
			models.ReDomesticAccessorialPrice{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_intl_accessorial_prices",
			models.ReIntlAccessorialPrice{},
			db.Where("contract_id = ?", contractID),
		},
		{
			"re_shipment_type_prices",
			models.ReShipmentTypePrice{},
			db.Where("contract_id = ?", contractID),
		},
	}

	for index, model := range models {
		if index != 0 {
			logger.Info("----")
		}
		err := summarizeModel(db, logger, model.header, model.modelInstance, model.filter)
		if err != nil {
			return err
		}
	}

	return nil
}

func summarizeModel(db *pop.Connection, logger logger, header string, modelInstance interface{}, filter *pop.Query) error {
	// Inspired by https://stackoverflow.com/a/25386460
	modelType := reflect.TypeOf(modelInstance)
	if modelType.Kind() != reflect.Struct {
		return fmt.Errorf("model type under header [%s] should be a struct, but got %s instead", header, modelType.Kind())
	}

	modelName := modelType.Name()
	modelSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 2)
	modelPtrSlice := reflect.New(modelSlice.Type())
	modelPtrSlice.Elem().Set(modelSlice)

	if filter == nil {
		filter = db.Q()
	}

	err := filter.Limit(2).All(modelPtrSlice.Interface())
	if err != nil {
		return err
	}
	length, err := filter.Count(modelInstance)
	if err != nil {
		return err
	}

	modelSlice = modelPtrSlice.Elem()

	headerMsg := fmt.Sprintf("%s (%s)", header, modelName)
	logger.Info(headerMsg, zap.Int("row count", length))
	if length > 0 {
		logger.Info("first:", zap.Any(modelName, modelSlice.Index(0).Interface()))
	}
	if length > 1 {
		logger.Info("second:", zap.Any(modelName, modelSlice.Index(1).Interface()))
	}

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
