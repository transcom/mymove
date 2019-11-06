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
	filename := flag.String("filename", "", "Filename including path of the XLSX to parse for Rate Engine GHC import")
	all := flag.Bool("all", true, "Parse entire Rate Engine GHC XLSX")
	sheets := flag.String("xlsxSheets", "", xlsxSheetsUsage(xlsxDataSheets))
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

	params.ProcessAll = false
	if all != nil && *all == true {
		params.ProcessAll = true
	}

	// option `xlsxSheets` will override `all` flag
	if sheets != nil && len(*sheets) > 0 {
		// If processes based on `xlsxSheets` indices provided as arguments
		// process those and do not run all
		params.ProcessAll = false
		params.XlsxSheets = strings.Split(*sheets, ",")
	}

	if filename != nil {
		log.Printf("Importing file %s\n", *filename)
	} else {
		log.Fatalf("Did not receive an XLSX filename to parse, missing -filename\n")
	}
	params.XlsxFilename = *filename

	xlsxFile, err := xlsx.OpenFile(params.XlsxFilename)
	params.XlsxFile = xlsxFile
	if err != nil {
		log.Fatalf("Failed to open file %s with error %v\n", params.XlsxFilename, err)
	}

	params.ShowOutput = false
	if display != nil && *display == true {
		params.ShowOutput = true
	}

	params.SaveToFile = false
	if saveToFile != nil && *saveToFile == true {
		params.SaveToFile = true
	}

	params.RunVerify = false
	if runVerify != nil {
		params.RunVerify = *runVerify
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

	// Now kick off the parsing
	err = pricing.Parse(xlsxDataSheets, params, db, logger)
	if err != nil {
		log.Fatalf("Failed to parse pricing template due to %v", err)
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
