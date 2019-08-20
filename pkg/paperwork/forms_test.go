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
	defer f.Close()

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
