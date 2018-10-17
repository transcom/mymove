package paperwork

import (
	"fmt"
	"image"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

var rankDisplayValue = map[internalmessages.ServiceMemberRank]string{
	internalmessages.ServiceMemberRankE1:                     "E-1",
	internalmessages.ServiceMemberRankE2:                     "E-2",
	internalmessages.ServiceMemberRankE3:                     "E-3",
	internalmessages.ServiceMemberRankE4:                     "E-4",
	internalmessages.ServiceMemberRankE5:                     "E-5",
	internalmessages.ServiceMemberRankE6:                     "E-6",
	internalmessages.ServiceMemberRankE7:                     "E-7",
	internalmessages.ServiceMemberRankE8:                     "E-8",
	internalmessages.ServiceMemberRankE9:                     "E-9",
	internalmessages.ServiceMemberRankO1W1ACADEMYGRADUATE:    "O-1/W-1/Service Academy Graduate",
	internalmessages.ServiceMemberRankO2W2:                   "O-2/W-2",
	internalmessages.ServiceMemberRankO3W3:                   "O-3/W-3",
	internalmessages.ServiceMemberRankO4W4:                   "O-4/W-4",
	internalmessages.ServiceMemberRankO5W5:                   "O-5/W-5",
	internalmessages.ServiceMemberRankO6:                     "O-6",
	internalmessages.ServiceMemberRankO7:                     "O-7",
	internalmessages.ServiceMemberRankO8:                     "O-8",
	internalmessages.ServiceMemberRankO9:                     "O-9",
	internalmessages.ServiceMemberRankO10:                    "O-10",
	internalmessages.ServiceMemberRankAVIATIONCADET:          "Aviation Cadet",
	internalmessages.ServiceMemberRankCIVILIANEMPLOYEE:       "Civilian Employee",
	internalmessages.ServiceMemberRankACADEMYCADETMIDSHIPMAN: "Service Academy Cadet/Midshipman",
}

// These are hardcoded for now
const (
	pageOrientation string  = "P"
	distanceUnit    string  = "mm"
	pageSize        string  = "letter"
	fontFamily      string  = "Helvetica"
	fontStyle       string  = ""
	fontSize        float64 = 7
	fontDir         string  = ""
	lineHeight      float64 = 3
	templateName    string  = "form_template"
	imageXPos       float64 = 0
	imageYPos       float64 = 0
	// 0-value will be auto-calculated from aspect ratio
	letterWidthMm  float64 = 215.9
	letterHeightMm float64 = 0
	// Whether the cursor should be advanced after placing image
	flow bool = false
	// Ties image to an existing link, either by link ID or URL
	imageLink    int    = 0
	imageLinkURL string = ""
)

func floatPtr(f float64) *float64 {
	return &f
}

// FormLayout houses both a background image form template and the layout of individual fields
type FormLayout struct {
	TemplateImagePath string
	FieldsLayout      map[string]FieldPos
}

// FieldPos encapsulates the starting position and width of a form field
type FieldPos struct {
	xPos       float64
	yPos       float64
	width      float64
	fontSize   *float64
	lineHeight *float64
}

// FormField returns a new field position
func FormField(xPos, yPos, width float64, fontSize, lineHeight *float64) FieldPos {
	return FieldPos{
		xPos:       xPos,
		yPos:       yPos,
		width:      width,
		fontSize:   fontSize,
		lineHeight: lineHeight,
	}
}

// FormFiller is a fillable pdf form
type FormFiller struct {
	pdf       *gofpdf.Fpdf
	fields    map[string]FieldPos
	useBorder bool
}

// NewTemplateForm turns a template image and fields mapping into a FormFiller instance
func NewTemplateForm(templateImage io.ReadSeeker, fields map[string]FieldPos) (FormFiller, error) {
	// Determine image type
	_, format, err := image.DecodeConfig(templateImage)
	if err != nil {
		return FormFiller{}, err
	}
	templateImage.Seek(0, io.SeekStart)

	pdf := gofpdf.New(pageOrientation, distanceUnit, pageSize, fontDir)
	pdf.SetMargins(0, 0, 0)
	pdf.AddPage()

	// Use provided image as document background
	opt := gofpdf.ImageOptions{
		ImageType: format,
		ReadDpi:   true,
	}
	pdf.RegisterImageOptionsReader("form_template", opt, templateImage)
	pdf.Image("form_template", imageXPos, imageYPos, letterWidthMm, letterHeightMm, flow, format, imageLink, imageLinkURL)

	pdf.SetFont(fontFamily, fontStyle, fontSize)

	newForm := FormFiller{
		pdf:    pdf,
		fields: fields,
	}

	return newForm, pdf.Error()
}

// UseBorders draws boxes around each form field
func (f *FormFiller) UseBorders() {
	f.useBorder = true
}

// DrawData draws the provided data set onto the form using the fields mapping
func (f *FormFiller) DrawData(data interface{}) error {
	borderStr := ""
	if f.useBorder {
		borderStr = "1"
	}

	r := reflect.ValueOf(data)
	for k := range f.fields {
		fieldVal := reflect.Indirect(r).FieldByName(k)
		val := fieldVal.Interface()

		formField := f.fields[k]
		f.pdf.MoveTo(formField.xPos, formField.yPos)

		// Turn value into a display string depending on type, will need
		// an explicit case for each type we're accommodating
		var displayValue string
		switch v := val.(type) {
		case string:
			displayValue = v
		case int64:
			displayValue = strconv.FormatInt(v, 10)
		case time.Time:
			displayValue = v.Format("02-Jan-2006")
		case internalmessages.ServiceMemberRank:
			displayValue = rankDisplayValue[v]
		case *internalmessages.ServiceMemberRank:
			if v != nil {
				displayValue = rankDisplayValue[*v]
			}
		case internalmessages.Affiliation:
			displayValue = string(v)
		case *internalmessages.Affiliation:
			if v != nil {
				displayValue = string(*v)
			}
		case models.Address:
			displayValue = v.Format()
		case *models.Address:
			if v != nil {
				displayValue = v.Format()
			}
		default:
			fmt.Println(v)
		}

		// Apply custom formatting options
		if formField.fontSize != nil {
			f.pdf.SetFontSize(*formField.fontSize)
		} else {
			f.pdf.SetFontSize(fontSize)
		}

		tempLineHeight := lineHeight
		if formField.lineHeight != nil {
			tempLineHeight = *formField.lineHeight
		}

		f.pdf.MultiCell(formField.width, tempLineHeight, displayValue, borderStr, "", false)
	}

	return f.pdf.Error()
}

// Output outputs the form to the provided file
func (f *FormFiller) Output(output io.Writer) error {
	return f.pdf.Output(output)
}
