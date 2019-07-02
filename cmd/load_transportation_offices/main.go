package main

import (
	// "bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	transportationoffices "github.com/transcom/mymove/pkg/services/transportation_offices"
)

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Debug("checking config")

	err := cli.CheckEIA(v)
	if err != nil {
		return err
	}

	err = cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {

	// DB Config
	cli.InitDatabaseFlags(flag)

	// EIA Open Data API
	cli.InitEIAFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

// Command: go run github.com/transcom/mymove/cmd/load_transportation_offices
func main() {
	inputFile := "./cmd/load_transportation_offices/data/To_Cntct_info_201906070930.xml"
	// officesPath := "./testdata/transportation_offices.xml"
	// inputFile := "To_Cntct_info_201906070930.xml"
	// outputFile := "/Users/lynzt/Downloads/transportationoffices.txt"

	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	fmt.Println("hi...")
	builder := transportationoffices.NewMigrationBuilder(dbConnection, logger)
	builder.Build(inputFile)

	// fileBytes := transportationoffices.ReadXMLFile(inputFile)
	// o := transportationoffices.UnmarshalXML(fileBytes)

	// offices := o.LISTGCNSLORGID.GCNSLORGID

	// fmt.Printf("# total offices: %d\n", len(offices))

	// usOfficesFilter := func(o transportationoffices.Office) bool {
	// 	return o.LISTGCNSLINFO.GCNSLINFO.CNSLCOUNTRY == "US"
	// }
	// usOffices := transportationoffices.FilterOffice(offices, usOfficesFilter)
	// fmt.Printf("# us only offices: %d\n", len(usOffices))

	// conusOfficesFilter := func(o transportationoffices.Office) bool {
	// 	return o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE != "AK" &&
	// 		o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE != "HI"
	// }
	// conusOffices := transportationoffices.FilterOffice(usOffices, conusOfficesFilter)
	// fmt.Printf("# conus only offices: %d\n", len(conusOffices))

	// f, err := os.Create(outputFile)
	// defer f.Close()
	// w := bufio.NewWriter(f)

	// counter := 0
	// for _, o := range conusOffices {
	// 	transportationoffices.WriteXMLLine(o, w)
	// 	dbOffices := transportationoffices.FindConusOffices(dbConnection, o, w)
	// 	dbPPSOs := transportationoffices.FindPPSOs(dbConnection, o)
	// 	res := transportationoffices.WriteDbRecs("office", dbOffices, w)
	// 	transportationoffices.WriteDbRecs("JPPSO", dbPPSOs, w)
	// 	counter += res
	// }
	// w.Flush()
	// fmt.Println(counter)

}
