package paperwork

import (
	"fmt"
	"image"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/pkg/errors"

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

var affiliationDisplayValue = map[internalmessages.Affiliation]string{
	internalmessages.AffiliationARMY:       "Army",
	internalmessages.AffiliationNAVY:       "Navy",
	internalmessages.AffiliationMARINES:    "Marines",
	internalmessages.AffiliationAIRFORCE:   "Air Force",
	internalmessages.AffiliationCOASTGUARD: "Coast Guard",
}

var deptIndDisplayValue = map[internalmessages.DeptIndicator]string{
	internalmessages.DeptIndicatorAIRFORCE: "Air Force",
	internalmessages.DeptIndicatorMARINES:  "Marines",
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
	pdf   *gofpdf.Fpdf
	debug bool
	pages int
}

// NewFormFiller turns a template image and fields mapping into a FormFiller instance
func NewFormFiller() *FormFiller {
	pdf := gofpdf.New(pageOrientation, distanceUnit, pageSize, fontDir)
	pdf.SetMargins(0, 0, 0)
	pdf.SetFont(fontFamily, fontStyle, fontSize)

	return &FormFiller{
		debug: false,
		pdf:   pdf,
	}
}

// AppendPage adds a page to a PDF
func (f *FormFiller) AppendPage(templateImage io.ReadSeeker, fields map[string]FieldPos, data interface{}) error {

	f.pdf.AddPage()
	f.pages++

	// Determine image type
	_, format, err := image.DecodeConfig(templateImage)
	if err != nil {
		return errors.Wrap(err, "could not decode image config")
	}
	templateImage.Seek(0, io.SeekStart)

	// Use provided image as document background
	opt := gofpdf.ImageOptions{
		ImageType: format,
		ReadDpi:   true,
	}

	formTemplate := fmt.Sprintf("form_template_%d", f.pages)
	f.pdf.RegisterImageOptionsReader(formTemplate, opt, templateImage)
	f.pdf.Image(formTemplate, imageXPos, imageYPos, letterWidthMm, letterHeightMm, flow, format, imageLink, imageLinkURL)

	err = f.drawData(fields, data)
	if err != nil {
		return err
	}

	if f.pdf.Error() != nil {
		return errors.Wrap(f.pdf.Error(), "error creating PDF")
	}

	return nil
}

// Debug draws boxes around each form field and overlays the
// fields name on the right.
func (f *FormFiller) Debug() {
	f.debug = true
}

// DrawDebugOverlay draws a red bordered box and overlays the field's name on the right.
func (f *FormFiller) drawDebugOverlay(xPos, yPos, width, lineHeight float64, label string) {
	dr, dg, db := f.pdf.GetDrawColor()
	f.pdf.SetDrawColor(255, 0, 0)

	tr, tg, tb := f.pdf.GetTextColor()
	f.pdf.SetTextColor(255, 0, 0)

	fs, _ := f.pdf.GetFontSize()
	f.pdf.SetFontSize(4)

	f.pdf.MoveTo(xPos, yPos)
	f.pdf.CellFormat(width, lineHeight, label, "1", 0, "R", false, 0, "")

	// Restore settings
	f.pdf.SetTextColor(tr, tg, tb)
	f.pdf.SetDrawColor(dr, dg, db)
	f.pdf.SetFontSize(fs)
}

// DrawData draws the provided data set onto the form using the fields mapping
func (f *FormFiller) drawData(fields map[string]FieldPos, data interface{}) error {
	r := reflect.ValueOf(data)
	for k := range fields {
		fieldVal := reflect.Indirect(r).FieldByName(k)
		val := fieldVal.Interface()

		formField := fields[k]
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
			displayValue = affiliationDisplayValue[v]
		case *internalmessages.Affiliation:
			if v != nil {
				displayValue = affiliationDisplayValue[*v]
			}
		case internalmessages.DeptIndicator:
			displayValue = deptIndDisplayValue[v]
		case *internalmessages.DeptIndicator:
			if v != nil {
				displayValue = "DI: " + deptIndDisplayValue[*v]
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

		f.pdf.MultiCell(formField.width, tempLineHeight, displayValue, "", "", false)

		// Draw a red-bordered box with the display value's key to the right
		if f.debug {
			f.drawDebugOverlay(formField.xPos, formField.yPos, formField.width, tempLineHeight, k)
		}
	}

	return f.pdf.Error()
}

// Output outputs the form to the provided file
func (f *FormFiller) Output(output io.Writer) error {
	return f.pdf.Output(output)
}
