package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"

	"github.com/transcom/mymove/pkg/auth"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
)

// hereRequestTimeout is how long to wait on HERE request before timing out (15 seconds).
const hereRequestTimeout = time.Duration(15) * time.Second

func noErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	moveID := flag.String("move", "", "The move ID to generate a shipment summary worksheet for")
	debug := flag.Bool("debug", false, "show field debug output")
	flag.Parse()

	// DB connection
	err := pop.AddLookupPaths(*config)
	if err != nil {
		log.Fatal(err)
	}
	db, err := pop.Connect(*env)
	if err != nil {
		log.Fatal(err)
	}

	if *moveID == "" {
		log.Fatal("Usage: generate_shipment_summary -move <29cb984e-c70d-46f0-926d-cd89e07a6ec3>")
	}

	// Define the data here that you want to populate the form with. Data will only be populated
	// in the form if the field name exist BOTH in the fields map and your data below
	parsedID := uuid.Must(uuid.FromString(*moveID))

	// Build our form with a template image and field placement
	formFiller := paperwork.NewFormFiller()

	// This is very useful for getting field positioning right initially
	if *debug {
		formFiller.Debug()
	}

	logger, err := logging.Config(*env, true)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	geocodeEndpoint := os.Getenv("HERE_MAPS_GEOCODE_ENDPOINT")
	routingEndpoint := os.Getenv("HERE_MAPS_ROUTING_ENDPOINT")
	testAppID := os.Getenv("HERE_MAPS_APP_ID")
	testAppCode := os.Getenv("HERE_MAPS_APP_CODE")
	hereClient := &http.Client{Timeout: hereRequestTimeout}
	planner := route.NewHEREPlanner(logger, hereClient, geocodeEndpoint, routingEndpoint, testAppID, testAppCode)
	ppmComputer := paperwork.NewSSWPPMComputer(rateengine.NewRateEngine(db, logger))

	ssfd, err := models.FetchDataShipmentSummaryWorksheetFormData(db, &auth.Session{}, parsedID)
	ssfd.Obligations, err = ppmComputer.ComputeObligations(ssfd, planner)
	if err != nil {
		log.Println("Error calculating obligations ")
	}

	page1Data, page2Data, err := models.FormatValuesShipmentSummaryWorksheet(ssfd)
	noErr(err)

	// page 1
	page1Layout := paperwork.ShipmentSummaryPage1Layout
	page1Template, err := os.Open(page1Layout.TemplateImagePath)
	noErr(err)
	defer page1Template.Close()

	err = formFiller.AppendPage(page1Template, page1Layout.FieldsLayout, page1Data)
	noErr(err)

	// page 2
	page2Layout := paperwork.ShipmentSummaryPage2Layout
	page2Template, err := os.Open(page2Layout.TemplateImagePath)
	noErr(err)
	defer page2Template.Close()

	err = formFiller.AppendPage(page2Template, page2Layout.FieldsLayout, page2Data)
	noErr(err)

	filename := fmt.Sprintf("shipment-summary-worksheet-%s.pdf", time.Now().Format(time.RFC3339))

	output, err := os.Create(filename)
	noErr(err)
	defer output.Close()

	err = formFiller.Output(output)
	noErr(err)

	fmt.Println(filename)
}
