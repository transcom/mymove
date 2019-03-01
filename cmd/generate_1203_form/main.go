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

type fakeModel struct {
	FieldName string
}

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
	shipmentID := flag.String("shipment", "", "The shipment ID to generate 1203 form for")
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
		log.Fatal("Usage: generate_1203_form -shipment <29cb984e-c70d-46f0-926d-cd89e07a6ec3>")
	}

	formLayout := paperwork.Form1203Layout

	templateImage, err := os.Open(formLayout.TemplateImagePath)
	noErr(err)
	defer templateImage.Close()

	// Define the data here that you want to populate the form with. Data will only be populated
	// in the form if the field name exist BOTH in the fields map and your data below
	parsedID := uuid.Must(uuid.FromString(*shipmentID))

	// Build our form with a template image and field placement
	formFiller := paperwork.NewFormFiller()

	gbl, err := models.FetchGovBillOfLadingFormValues(db, parsedID)
	noErr(err)

	// This is very useful for getting field positioning right initially
	if *debug {
		formFiller.Debug()
	}

	// Populate form fields with provided data
	err = formFiller.AppendPage(templateImage, formLayout.FieldsLayout, gbl)
	noErr(err)

	filename := fmt.Sprintf("form-1203-%s.pdf", time.Now().Format(time.RFC3339))

	output, err := os.Create(filename)
	noErr(err)

	err = formFiller.Output(output)
	noErr(err)

	fmt.Println(filename)
}
