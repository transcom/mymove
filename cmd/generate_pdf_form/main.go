package main

import (
	"fmt"
	"log"
	"os"

	"github.com/transcom/mymove/pkg/paperwork"
)

type fakeModel struct {
	FieldName string
}

func noErr(err error) {
	if err != nil {
		log.Panic("oops ", err)
	}
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	// This is the path to an image you want to use as a form template
	templateImagePath := "./cmd/generate_pdf_form/example_template.png"

	f, err := os.Open(templateImagePath)
	noErr(err)
	defer f.Close()

	// Define your field positions here, it should be a mapping from a struct field name
	// to a FieldPos, which encodes the x and y location, and width of a form field
	var fields = map[string]paperwork.FieldPos{
		"FieldName": paperwork.NewFieldPos(20, 14, 79),
	}

	// Define the data here that you want to populate the form with. Data will only be populated
	// in the form if the field name exist BOTH in the fields map and your data below
	data := fakeModel{
		FieldName: "Data goes here",
	}

	// Build our form with a template image and field placement
	form, err := paperwork.NewTemplateForm(f, fields)
	noErr(err)

	// Uncomment the below line if you want to draw borders around field boxes, very useful
	// for getting field positioning right initially
	form.UseBorders()

	// Populate form fields with provided data
	err = form.DrawData(data)
	noErr(err)

	output, _ := os.Create("./cmd/generate_pdf_form/test-output.pdf")
	err = form.Output(output)
	noErr(err)

	fmt.Println("done!")
}
