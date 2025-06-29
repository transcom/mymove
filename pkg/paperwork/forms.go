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

var rankDisplayValue = map[internalmessages.OrderPayGrade]string{
	models.ServiceMemberGradeE1:                      "E-1",
	models.ServiceMemberGradeE2:                      "E-2",
	models.ServiceMemberGradeE3:                      "E-3",
	models.ServiceMemberGradeE4:                      "E-4",
	models.ServiceMemberGradeE5:                      "E-5",
	models.ServiceMemberGradeE6:                      "E-6",
	models.ServiceMemberGradeE7:                      "E-7",
	models.ServiceMemberGradeE8:                      "E-8",
	models.ServiceMemberGradeE9:                      "E-9",
	models.ServiceMemberGradeE9SPECIALSENIORENLISTED: "E-9 (Special Senior Enlisted)",
	models.ServiceMemberGradeO1:                      "O-1 or Service Academy Graduate",
	models.ServiceMemberGradeO2:                      "O-2",
	models.ServiceMemberGradeO3:                      "O-3",
	models.ServiceMemberGradeO4:                      "O-4",
	models.ServiceMemberGradeO5:                      "O-5",
	models.ServiceMemberGradeO6:                      "O-6",
	models.ServiceMemberGradeO7:                      "O-7",
	models.ServiceMemberGradeO8:                      "O-8",
	models.ServiceMemberGradeO9:                      "O-9",
	models.ServiceMemberGradeO10:                     "O-10",
	models.ServiceMemberGradeW1:                      "W-1",
	models.ServiceMemberGradeW2:                      "W-2",
	models.ServiceMemberGradeW3:                      "W-3",
	models.ServiceMemberGradeW4:                      "W-4",
	models.ServiceMemberGradeW5:                      "W-5",
	models.ServiceMemberGradeAVIATIONCADET:           "Aviation Cadet",
	models.ServiceMemberGradeCIVILIANEMPLOYEE:        "Civilian Employee",
	models.ServiceMemberGradeACADEMYCADET:            "Service Academy Cadet",
	models.ServiceMemberGradeMIDSHIPMAN:              "Midshipman",
}

var affiliationDisplayValue = map[internalmessages.Affiliation]string{
	internalmessages.AffiliationARMY:       "Army",
	internalmessages.AffiliationNAVY:       "Navy",
	internalmessages.AffiliationMARINES:    "Marines",
	internalmessages.AffiliationAIRFORCE:   "Air Force",
	internalmessages.AffiliationCOASTGUARD: "Coast Guard",
	internalmessages.AffiliationSPACEFORCE: "Space Force",
}

var serviceMemberAffiliationDisplayValue = map[models.ServiceMemberAffiliation]string{
	models.AffiliationARMY:       "Army",
	models.AffiliationNAVY:       "Navy",
	models.AffiliationMARINES:    "Marines",
	models.AffiliationAIRFORCE:   "Air Force",
	models.AffiliationCOASTGUARD: "Coast Guard",
	models.AffiliationSPACEFORCE: "Space Force",
}

var deptIndDisplayValue = map[internalmessages.DeptIndicator]string{
	internalmessages.DeptIndicatorAIRANDSPACEFORCE: "Air Force and Space Force",
	internalmessages.DeptIndicatorNAVYANDMARINES:   "Navy and Marine Corps",
}

// These are hardcoded for now
const (
	pageOrientation string  = "P"
	distanceUnit    string  = "mm"
	pageSize        string  = "letter"
	fontFamily      string  = "Helvetica"
	fontStyle       string  = ""
	fontSize        float64 = 7
	alignment       string  = "LM"
	// Horizontal alignment is controlled by including "L", "C" or "R" (left, center, right) in alignStr.
	// Vertical alignment is controlled by including "T", "M", "B" or "A" (top, middle, bottom, baseline) in alignStr.
	// The default alignment is left middle.
	fontDir    string  = ""
	lineHeight float64 = 3
	imageXPos  float64 = 0
	imageYPos  float64 = 0
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

func stringPtr(s string) *string {
	return &s
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
	alignStr   *string
}

// FormField returns a new field position
func FormField(xPos, yPos, width float64, fontSize, lineHeight *float64, alignStr *string) FieldPos {
	return FieldPos{
		xPos:       xPos,
		yPos:       yPos,
		width:      width,
		fontSize:   fontSize,
		lineHeight: lineHeight,
		alignStr:   alignStr,
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
	pdf.SetAutoPageBreak(false, 0)

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
	_, err = templateImage.Seek(0, io.SeekStart)
	if err != nil {
		return errors.Wrap(err, "could not read image data")
	}

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
	f.pdf.CellFormat(width, lineHeight, "", "1", 0, "R", false, 0, "")

	f.pdf.MoveTo(xPos+1.2, yPos-2.9)
	f.pdf.CellFormat(width, 4, label, "0", 0, "R", false, 0, "")

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
		case internalmessages.OrderPayGrade:
			displayValue = rankDisplayValue[v]
		case *internalmessages.OrderPayGrade:
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
		case fmt.Stringer:
			displayValue = v.String()
		default:
			// TODO: error out?
			// fmt.Println(v)
		}

		// Apply custom formatting options
		if formField.fontSize != nil {
			f.pdf.SetFontSize(*formField.fontSize)
		} else {
			f.pdf.SetFontSize(fontSize)
		}

		fs, _ := f.pdf.GetFontSize()
		f.ScaleText(displayValue, fs, formField.width)

		tempLineHeight := lineHeight
		if formField.lineHeight != nil {
			tempLineHeight = *formField.lineHeight
		}

		fieldAlignment := alignment
		if formField.alignStr != nil {
			fieldAlignment = *formField.alignStr
		}

		f.pdf.MultiCell(formField.width, tempLineHeight, displayValue, "", fieldAlignment, false)

		// Draw a red-bordered box with the display value's key to the right
		if f.debug {
			f.drawDebugOverlay(formField.xPos, formField.yPos, formField.width, tempLineHeight, k)
		}
	}

	return f.pdf.Error()
}

// ScaleText scales the text down to fit the cell if it is too long to fit in one line
func (f *FormFiller) ScaleText(displayValue string, _ float64, width float64) {
	stringWidth := f.pdf.GetStringWidth(displayValue)
	for f.isMoreThanOneLine(stringWidth, width) {
		ptSize, _ := f.pdf.GetFontSize()
		newFontSize := ptSize * .95
		if newFontSize <= 5 {
			break
		}
		f.pdf.SetFontSize(newFontSize)
		stringWidth = f.pdf.GetStringWidth(displayValue)
	}
}

func (f *FormFiller) isMoreThanOneLine(stringWidth float64, width float64) bool {
	return stringWidth > (width - 2*f.pdf.GetCellMargin())
}

// Output outputs the form to the provided file
func (f *FormFiller) Output(output io.Writer) error {
	return f.pdf.Output(output)
}
