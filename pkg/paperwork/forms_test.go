package paperwork

import (
	"os"

	"github.com/spf13/afero"
)

type fakeModel struct {
	FieldName string
}

// Tests if we can fill a form without blowing up
func (suite *PaperworkSuite) TestFormFillerSmokeTest() {
	templateImagePath := "./testdata/example_template.png"

	f, err := os.Open(templateImagePath)
	suite.FatalNil(err)
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used to close a local server connection to ensure a unit test server is not left running indefinitely
	//RA: Given the functions causing the lint errors are used to close a local server connection for testing purposes, it is not deemed a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	defer f.Close() // nolint:errcheck

	var fields = map[string]FieldPos{
		"FieldName": FormField(28, 11, 79, nil, nil, nil),
	}

	data := fakeModel{
		FieldName: "Data goes here",
	}

	formFiller := NewFormFiller()
	err = formFiller.AppendPage(f, fields, data)
	suite.FatalNil(err)

	testFs := afero.NewMemMapFs()

	output, err := testFs.Create("test-output.pdf")
	suite.FatalNil(err)

	err = formFiller.Output(output)
	suite.FatalNil(err)
}

func (suite *PaperworkSuite) TestFormScaleFont() {
	formFiller := NewFormFiller()
	formFiller.pdf.SetFontSize(10)
	var cellWidth float64 = 60
	value := "Joint Base McGuire-Dix-Lakehurst, NJ  08641"
	stringWidth := formFiller.pdf.GetStringWidth(value)
	tooLong := formFiller.isMoreThanOneLine(stringWidth, cellWidth)

	suite.True(tooLong)
	ptSize, _ := formFiller.pdf.GetFontSize()
	formFiller.ScaleText(value, ptSize, cellWidth)

	stringWidth = formFiller.pdf.GetStringWidth(value)
	tooLong = formFiller.isMoreThanOneLine(stringWidth, cellWidth)
	suite.False(tooLong)

}
