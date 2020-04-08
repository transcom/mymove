package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

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

// SswPageData holds all ssw page data formatted for inserting into the ssw pdf
type SswPageData struct {
	page1 models.ShipmentSummaryWorksheetPage1Values
	page2 models.ShipmentSummaryWorksheetPage2Values
	page3 models.ShipmentSummaryWorksheetPage3Values
}

func noErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func checkCloseoutDataConfig(v *viper.Viper, logger logger) error {

	logger.Debug("checking config")

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

func initFetchCloseoutDataFlags(flag *pflag.FlagSet) {

	// Scenario config
	flag.String(moveIDFlag, "", "The move ID to generate a shipment summary worksheet for")
	flag.Bool(debugFlag, false, "show field debug output")

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

// Command: go run ./cmd/milmove-tasks fetch-closeout-data
func fetchCloseoutData(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "Could not parse args")
	}
	flags := cmd.Flags()
	v := viper.New()
	err = v.BindPFlags(flags)
	if err != nil {
		return errors.Wrap(err, "Could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkCloseoutDataConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	var session *awssession.Session
	if v.GetBool(cli.DbIamFlag) || (v.GetString(cli.EmailBackendFlag) == "ses") {
		c, errorConfig := cli.GetAWSConfig(v, v.GetBool(cli.VerboseFlag))
		if errorConfig != nil {
			logger.Fatal(errors.Wrap(errorConfig, "error creating aws config").Error())
		}
		s, errorSession := awssession.NewSession(c)
		if errorSession != nil {
			logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
		}
		session = s
	}

	var dbCreds *credentials.Credentials
	if v.GetBool(cli.DbIamFlag) {
		if session != nil {
			// We want to get the credentials from the logged in AWS session rather than create directly,
			// because the session conflates the environment, shared, and container metdata config
			// within NewSession.  With stscreds, we use the Secure Token Service,
			// to assume the given role (that has rds db connect permissions).
			dbIamRole := v.GetString(cli.DbIamRoleFlag)
			logger.Info(fmt.Sprintf("assuming AWS role %q for db connection", dbIamRole))
			dbCreds = stscreds.NewCredentials(session, dbIamRole)
		}
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, dbCreds, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	// closeout, err := closeout.NewCloseoutData(dbConnection, logger)
	// if err != nil {
	// 	logger.Fatal("initializing CloseoutData", zap.Error(err))
	// }
	// fmt.Println(closeout)
	// moveIDs := []string{"77CXF9"}
	// _ = closeout.FetchCloseoutDetails(moveIDs)

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
	planner := route.NewHEREPlanner(logger, hereClient, geocodeEndpoint, routingEndpoint, testAppID, testAppCode)
	ppmComputer := paperwork.NewSSWPPMComputer(rateengine.NewRateEngine(dbConnection, logger, move))

	ssfd, err := models.FetchDataShipmentSummaryWorksheetFormData(dbConnection, &auth.Session{}, parsedID)
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "Error fetching shipment summary worksheet data "))
	}
	ssfd.Obligations, err = ppmComputer.ComputeObligations(ssfd, planner)
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "Error calculating obligations "))
	}

	page1Data, page2Data, page3Data, err := models.FormatValuesShipmentSummaryWorksheet(ssfd)
	noErr(err)

	fmt.Println("------------------------------")
	fmt.Println("------------------------------")

	sswPageData := SswPageData{page1Data, page2Data, page3Data}
	fmt.Printf("%+v\n", sswPageData)

	// this didn't want to work - idk...
	// file, _ := json.MarshalIndent(sswPageData, "", " ")
	// fmt.Println(file)

	file, _ := json.MarshalIndent(ssfd, "", " ")
	ioutil.WriteFile("pptas.json", file, 0644)

	return nil
}
