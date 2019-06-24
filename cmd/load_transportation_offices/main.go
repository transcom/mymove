package main

import (
	"bufio"
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

// Command: go run github.com/transcom/mymove/cmd/save_fuel_price_data
func main() {

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

	fileBytes := transportationoffices.ReadXMLFile("To_Cntct_info_201906070930.xml")

	o := transportationoffices.UnmarshalXML(fileBytes)

	offices := o.LISTGCNSLORGID.GCNSLORGID

	fmt.Printf("Name: %s\n", offices[0].LISTGCNSLINFO.GCNSLINFO.CNSLNAME)
	fmt.Printf("# total offices: %d\n", len(offices))

	usOfficesFilter := func(o transportationoffices.Office) bool {
		return o.LISTGCNSLINFO.GCNSLINFO.CNSLCOUNTRY == "US"
	}
	usOffices := transportationoffices.Filter(offices, usOfficesFilter)
	fmt.Printf("# us only offices: %d\n", len(usOffices))

	conusOfficesFilter := func(o transportationoffices.Office) bool {
		return o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE != "AK" &&
			o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE != "HI"
	}
	conusOffices := transportationoffices.Filter(usOffices, conusOfficesFilter)
	fmt.Printf("# conus only offices: %d\n", len(conusOffices))

	// f := transportationoffices.OpenFile()

	f, err := os.Create("/Users/lynzt/Downloads/transportationoffices.txt")
	defer f.Close()
	w := bufio.NewWriter(f)

	// for _, o := range conusOffices[1:5] {
	for _, o := range conusOffices {
		dbOffices := transportationoffices.CheckDbForConusOffices(dbConnection, o)
		transportationoffices.OutputResults(o, dbOffices, w)
	}
	w.Flush()

}
