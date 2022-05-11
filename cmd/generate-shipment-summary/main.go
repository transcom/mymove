package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
)

// hereRequestTimeout is how long to wait on HERE request before timing out (15 seconds).
const hereRequestTimeout = time.Duration(15) * time.Second

const (
	moveIDFlag string = "move"
	debugFlag  string = "debug"
)

func noErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func checkConfig(v *viper.Viper, logger *zap.Logger) error {

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

	// Scenario config
	flag.String(moveIDFlag, "", "The move ID to generate a shipment summary worksheet for")
	flag.Bool(debugFlag, false, "show field debug output")

	// DB Config
	cli.InitDatabaseFlags(flag)

	// EIA Open Data API
	cli.InitEIAFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func main() {
	flag := pflag.CommandLine
	initFlags(flag)
	parseErr := flag.Parse(os.Args[1:])
	if parseErr != nil {
		log.Fatal("failed to parse flags", zap.Error(parseErr))
	}

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		log.Fatal("failed to bind flags", zap.Error(bindErr))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, _, err := logging.Config(
		logging.WithEnvironment(dbEnv),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	fmt.Println("logger: ", logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// DB connection
	dbConnection, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	moveID := v.GetString(moveIDFlag)
	if moveID == "" {
		log.Fatal("Usage: generate_shipment_summary -move <29cb984e-c70d-46f0-926d-cd89e07a6ec3>")
	}

	// Define the data here that you want to populate the form with. Data will only be populated
	// in the form if the field name exist BOTH in the fields map and your data below
	parsedID := uuid.Must(uuid.FromString(moveID))

	// Build our form with a template image and field placement
	formFiller := paperwork.NewFormFiller()

	debug := v.GetBool(debugFlag)
	// This is very useful for getting field positioning right initially
	if debug {
		formFiller.Debug()
	}

	move, err := models.FetchMoveByMoveID(dbConnection, parsedID)
	if err != nil {
		log.Fatalf("error fetching move: %s", moveIDFlag)
	}

	geocodeEndpoint := os.Getenv("HERE_MAPS_GEOCODE_ENDPOINT")
	routingEndpoint := os.Getenv("HERE_MAPS_ROUTING_ENDPOINT")
	testAppID := os.Getenv("HERE_MAPS_APP_ID")
	testAppCode := os.Getenv("HERE_MAPS_APP_CODE")
	hereClient := &http.Client{Timeout: hereRequestTimeout}
	planner := route.NewHEREPlanner(hereClient, geocodeEndpoint, routingEndpoint, testAppID, testAppCode)
	ppmComputer := paperwork.NewSSWPPMComputer(rateengine.NewRateEngine(move))

	ssfd, err := models.FetchDataShipmentSummaryWorksheetFormData(dbConnection, &auth.Session{}, parsedID)
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "Error fetching shipment summary worksheet data "))
	}
	ssfd.Obligations, err = ppmComputer.ComputeObligations(appCtx, ssfd, planner)
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "Error calculating obligations "))
	}

	page1Data, page2Data, page3Data, err := models.FormatValuesShipmentSummaryWorksheet(ssfd)
	noErr(err)

	// page 1
	page1Layout := paperwork.ShipmentSummaryPage1Layout
	page1Template, err := os.Open(page1Layout.TemplateImagePath)
	noErr(err)

	defer func() {
		if closeErr := page1Template.Close(); closeErr != nil {
			logger.Error("Could not close page1Template file", zap.Error(closeErr))
		}
	}()

	err = formFiller.AppendPage(page1Template, page1Layout.FieldsLayout, page1Data)
	noErr(err)

	// page 2
	page2Layout := paperwork.ShipmentSummaryPage2Layout
	page2Template, err := os.Open(page2Layout.TemplateImagePath)
	noErr(err)

	defer func() {
		if closeErr := page2Template.Close(); closeErr != nil {
			logger.Error("Could not close page2Template file", zap.Error(closeErr))
		}
	}()

	err = formFiller.AppendPage(page2Template, page2Layout.FieldsLayout, page2Data)
	noErr(err)

	// page 3
	page3Layout := paperwork.ShipmentSummaryPage3Layout
	page3Template, err := os.Open(page3Layout.TemplateImagePath)
	noErr(err)

	defer func() {
		if closeErr := page3Template.Close(); closeErr != nil {
			logger.Error("Could not close page3Template file", zap.Error(closeErr))
		}
	}()

	err = formFiller.AppendPage(page3Template, page3Layout.FieldsLayout, page3Data)
	noErr(err)

	filename := fmt.Sprintf("shipment-summary-worksheet-%s.pdf", time.Now().Format(time.RFC3339))

	output, err := os.Create(filename)
	noErr(err)

	defer func() {
		if closeErr := output.Close(); closeErr != nil {
			logger.Error("Could not close output file", zap.Error(closeErr))
		}
	}()

	err = formFiller.Output(output)
	noErr(err)

	fmt.Println(filename)
}
