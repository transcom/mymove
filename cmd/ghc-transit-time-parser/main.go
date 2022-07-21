package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx/v3"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/parser/transittime"
)

/*************************************************************************

Parser tool to extract data from the GHC HHG Transit Times XLSX

For help run: <program> -h

`go run ./cmd/ghc-hhg-transit-times-parser/ -h`

 *************************************************************************/

func main() {
	xlsxDataSheets := transittime.InitDataSheetInfo()

	params := transittime.ParamConfig{}
	params.RunTime = time.Now()

	flag := pflag.CommandLine
	flag.StringVar(&params.XlsxFilename, "filename", "", "Filename including path of the XLSX to parse for Transit Times GHC import")
	flag.BoolVar(&params.ProcessAll, "all", true, "Parse entire Transit Times GHC XLSX")
	flag.StringSliceVar(&params.XlsxSheets, "xlsxSheets", []string{}, xlsxSheetsUsage(xlsxDataSheets))
	flag.BoolVar(&params.ShowOutput, "display", false, "Display output of parsed info")
	flag.BoolVar(&params.SaveToFile, "save-csv", true, "Save output to CSV file")
	flag.BoolVar(&params.RunVerify, "verify", true, "Default is true, if false skip sheet format verification")
	flag.BoolVar(&params.RunImport, "re-import", true, "Run GHC Transit Times Import")

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
		log.Println("Setting --xlsxSheets disables --re-import so no data will be imported into the Transit Times tables. Only stage table data will be updated.")
	}

	if params.XlsxFilename == "" {
		log.Fatalf("Did not receive an XLSX filename to parse; missing --filename\n")
	}
	log.Printf("Importing file %s\n", params.XlsxFilename)

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

	logger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
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

	appCtx := appcontext.NewAppContext(db, logger, nil)

	// Now kick off the parsing
	err = transittime.Parse(appCtx, xlsxDataSheets, params)
	if err != nil {
		log.Fatalf("Failed to parse transit times template due to %v", err)
	}
}

func xlsxSheetsUsage(xlsxDataSheets []transittime.XlsxDataSheetInfo) string {
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
	message += "NOTE: This option disables the Transit Times table import by disabling the --re-import flag\n"

	return message
}
