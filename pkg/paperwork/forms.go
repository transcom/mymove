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

// These are hardcoded for now
const (
	pageOrientation string  = "P"
	distanceUnit    string  = "mm"
	pageSize        string  = "letter"
	fontFamily      string  = "Helvetica"
	fontStyle       string  = ""
	fontSize        float64 = 8
	fontDir         string  = ""
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

// FieldPos encapsulates the starting position and width of a form field
type FieldPos struct {
	xPos  float64
	yPos  float64
	width float64
}

// NewFieldPos returns a new field position
func NewFieldPos(xPos, yPos, width float64) FieldPos {
	return FieldPos{
		xPos:  xPos,
		yPos:  yPos,
		width: width,
	}
}

// Form is a fillable pdf form
type Form struct {
	pdf       *gofpdf.Fpdf
	fields    map[string]FieldPos
	useBorder bool
}

// NewTemplateForm turns a template image and fields mapping into a Form instance
func NewTemplateForm(templateImage io.ReadSeeker, fields map[string]FieldPos) (Form, error) {
	// Determine image type
	_, format, err := image.DecodeConfig(templateImage)
	if err != nil {
		return Form{}, err
	}
	templateImage.Seek(0, io.SeekStart)

	pdf := gofpdf.New(pageOrientation, distanceUnit, pageSize, fontDir)
	pdf.SetMargins(0, 0, 0)
	pdf.AddPage()

	// Use provided image as document background
	var opt gofpdf.ImageOptions
	opt.ImageType = format
	pdf.RegisterImageReader("form_template", format, templateImage)
	pdf.ImageOptions("form_template", imageXPos, imageYPos, letterWidthMm, letterHeightMm, flow, opt, imageLink, imageLinkURL)

	pdf.SetFont(fontFamily, fontStyle, fontSize)

	newForm := Form{
		pdf:    pdf,
		fields: fields,
	}

	return newForm, pdf.Error()
}

// UseBorders draws boxes around each form field
func (f *Form) UseBorders() {
	f.useBorder = true
}

// DrawData draws the provided data set onto the form using the fields mapping
func (f *Form) DrawData(data interface{}) error {
	borderStr := ""
	if f.useBorder {
		borderStr = "1"
	}

	r := reflect.ValueOf(data)
	for k := range f.fields {
		fieldVal := reflect.Indirect(r).FieldByName(k)
		val := fieldVal.Interface()

		layout := f.fields[k]
		f.pdf.MoveTo(layout.xPos, layout.yPos)

		// Turn value into a display string depending on type, will need
		// an explicit case for each type we're accommodating
		var displayValue string
		switch v := val.(type) {
		case string:
			displayValue = v
		case int64:
			displayValue = strconv.FormatInt(v, 10)
		case time.Time:
			displayValue = v.Format("20060102")
		case internalmessages.ServiceMemberRank:
			displayValue = string(v)
		case *internalmessages.ServiceMemberRank:
			if v != nil {
				displayValue = string(*v)
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

		f.pdf.MultiCell(layout.width, 4, displayValue, borderStr, "", false)
	}

	return f.pdf.Error()
}

// Output outputs the form to the provided file
func (f *Form) Output(output io.WriteCloser) error {
	return f.pdf.OutputAndClose(output)
}
