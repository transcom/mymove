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
	"github.com/transcom/mymove/pkg/paperwork"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/route"
	ppmcloseout "github.com/transcom/mymove/pkg/services/ppm_closeout"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	shipmentsummaryworksheet "github.com/transcom/mymove/pkg/services/shipment_summary_worksheet"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

// hereRequestTimeout is how long to wait on HERE request before timing out (15 seconds).
const hereRequestTimeout = time.Duration(15) * time.Second

const (
	PPMShipmentIDFlag string = "ppmshipment"
	debugFlag         string = "debug"
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
	flag.String(PPMShipmentIDFlag, "6d1d9d00-2e5e-4830-a3c1-5c21c951e9c1", "The PPMShipmentID to generate a shipment summary worksheet for")
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
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	moveID := v.GetString(PPMShipmentIDFlag)
	if moveID == "" {
		log.Fatalf("Usage: %s --move <29cb984e-c70d-46f0-926d-cd89e07a6ec3>", os.Args[0])
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

	if err != nil {
		log.Fatalf("error fetching ppmshipment: %s", PPMShipmentIDFlag)
	}

	geocodeEndpoint := os.Getenv("HERE_MAPS_GEOCODE_ENDPOINT")
	routingEndpoint := os.Getenv("HERE_MAPS_ROUTING_ENDPOINT")
	testAppID := os.Getenv("HERE_MAPS_APP_ID")
	testAppCode := os.Getenv("HERE_MAPS_APP_CODE")
	hereClient := &http.Client{Timeout: hereRequestTimeout}

	// TODO: Future cleanup will need to remap to a different planner, but this command should remain for testing purposes
	planner := route.NewHEREPlanner(hereClient, geocodeEndpoint, routingEndpoint, testAppID, testAppCode)

	ppmEstimator := ppmshipment.NewEstimatePPM(planner, &paymentrequesthelper.RequestPaymentHelper{})
	ppmCloseoutFetcher := ppmcloseout.NewPPMCloseoutFetcher(planner, &paymentrequesthelper.RequestPaymentHelper{}, ppmEstimator)

	ppmComputer := shipmentsummaryworksheet.NewSSWPPMComputer(ppmCloseoutFetcher)

	ssfd, err := ppmComputer.FetchDataShipmentSummaryWorksheetFormData(appCtx, &auth.Session{}, parsedID)
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "Error fetching shipment summary worksheet data "))
	}
	ssfd.Obligations, err = ppmComputer.ComputeObligations(appCtx, *ssfd, planner)
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "Error calculating obligations "))
	}

	storer := storage.NewMemory(storage.NewMemoryParams("", ""))
	userUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		log.Fatalf("could not instantiate uploader due to %v", err)
	}
	generator, err := paperwork.NewGenerator(userUploader.Uploader())
	if err != nil {
		log.Fatal(err.Error())
	}

	page1Data, page2Data, page3Data, err := ppmComputer.FormatValuesShipmentSummaryWorksheet(*ssfd, false)
	noErr(err)
	ppmGenerator, err := shipmentsummaryworksheet.NewSSWPPMGenerator(generator)
	noErr(err)
	ssw, info, err := ppmGenerator.FillSSWPDFForm(page1Data, page2Data, page3Data)
	noErr(err)
	fmt.Println(ssw.Name())     // Should always return
	fmt.Println(info.PageCount) // Page count should always be 2
	// This is a testing command, above lines log information on whether PDF was generated successfully.
}
