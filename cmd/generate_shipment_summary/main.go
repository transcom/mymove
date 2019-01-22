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

	page1Data, page2Data, err := models.FetchShipmentSummaryWorksheetFormValues(db, parsedID)
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
