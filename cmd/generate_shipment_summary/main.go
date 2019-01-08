package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
)

func noErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	shipmentID := flag.String("shipment", "", "The shipment ID to generate a shipment summary worksheet for")
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

	if *shipmentID == "" {
		log.Fatal("Usage: generate_shipment_summary -shipment <29cb984e-c70d-46f0-926d-cd89e07a6ec3>")
	}

	formLayout := paperwork.ShipmentSummaryPage1Layout

	f, err := os.Open(formLayout.TemplateImagePath)
	noErr(err)
	defer f.Close()

	// Define the data here that you want to populate the form with. Data will only be populated
	// in the form if the field name exist BOTH in the fields map and your data below
	parsedID := uuid.Must(uuid.FromString(*shipmentID))
	data, err := models.FetchShipmentSummaryWorksheetExtractor(db, parsedID)
	noErr(err)

	// Build our form with a template image and field placement
	form, err := paperwork.NewTemplateForm(f, formLayout.FieldsLayout)
	noErr(err)

	// This is very useful for getting field positioning right initially
	if *debug {
		form.Debug()
	}

	// Populate form fields with provided data
	err = form.DrawData(data)
	noErr(err)

	filename := fmt.Sprintf("shipment-summary-worksheet-%s.pdf", time.Now().Format(time.RFC3339))

	output, err := os.Create(filename)
	noErr(err)

	err = form.Output(output)
	noErr(err)

	fmt.Println(filename)
}
