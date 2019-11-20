package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/tealeg/xlsx"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
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
	}

	if len(params.XlsxFilename) > 0 {
		log.Printf("Importing file %s\n", params.XlsxFilename)
	} else {
		log.Fatalf("Did not receive an XLSX filename to parse, missing -filename\n")
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
			Logger: logger,
		}
		err = ghcREImporter.Import(db)
		if err != nil {
			log.Fatalf("GHC Rate Engine import failed due to %v", err)
		}
	}

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

	return message
}
