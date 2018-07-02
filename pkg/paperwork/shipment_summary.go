package paperwork

import (
	"io"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"

	"github.com/transcom/mymove/pkg/models"
)

// Converts a string pointer into either a string or a default value if the pointer
// is nil.
func coalesce(stringPointer *string, defaultValue string) string {
	if stringPointer != nil {
		return *stringPointer
	}
	return defaultValue
}

type formField struct {
	label string
	value string
}

// ShipmentSummary encapsulates the process of drawing a PDF shipment summary form.
type ShipmentSummary struct {
	pdf  *gofpdf.Fpdf
	move *models.Move
}

// NewShipmentSummary creates and returns a new ShipmentSummary.
func NewShipmentSummary(move *models.Move) *ShipmentSummary {
	horizontalMargin := 0.0
	topMargin := 0.0

	pdf := gofpdf.New(PdfOrientation, PdfUnit, PdfPageSize, PdfFontDir)
	pdf.SetMargins(horizontalMargin, topMargin, horizontalMargin)

	return &ShipmentSummary{
		pdf:  pdf,
		move: move,
	}
}

const horizontalMargin = 15.0
const topMargin = 10.0
const bodyWidth = PdfPageWidth - (horizontalMargin * 2)
const fieldHeight = 5.0
const fontFace = "Helvetica"

// DrawForm Writes a Shipment Summary PDF to the provided ReadWriter.
func (s *ShipmentSummary) DrawForm(outputFile io.ReadWriter) error {
	s.pdf.SetMargins(horizontalMargin, topMargin, horizontalMargin)

	sm := s.move.Orders.ServiceMember

	s.pdf.SetHeaderFunc(func() {
		s.pdf.SetFont(fontFace, "B", 17)
		s.pdf.Cell(bodyWidth*0.75, fieldHeight*2, "SHIPMENT SUMMARY WORKSHEET - PPM")
		s.setFieldLabelFont()
		s.pdf.Cell(bodyWidth*0.25, fieldHeight, "Date Prepared (YYYYMMDD)")
		s.setFieldValueFont()
		s.pdf.SetXY(horizontalMargin+bodyWidth*0.75, s.pdf.GetY()+5)
		s.pdf.Cell(bodyWidth*0.25, fieldHeight, time.Now().Format("2006-01-02"))
		s.pdf.SetLineWidth(1.0)
		s.pdf.Line(0, 20, PdfPageWidth, 20)
		s.pdf.Ln(-1)
	})

	s.pdf.AddPage()
	s.addSectionHeader("MEMBER OR EMPLOYEE INFORMATION")

	// Name + phone
	nameFieldWidth := bodyWidth * 0.75
	nameLabelWidth := 40.0
	s.setFieldLabelFont()
	s.pdf.Cell(nameLabelWidth, fieldHeight, "Name")
	s.pdf.Ln(-1)
	s.pdf.SetFont(fontFace, "", 8)
	s.pdf.Cell(nameLabelWidth, fieldHeight, "(Last, First, Middle Initial)")

	s.pdf.SetXY(s.pdf.GetX(), s.pdf.GetY()-fieldHeight)
	s.pdf.SetFont(fontFace, "", 18)

	names := []string{}
	if sm.FirstName != nil && len(*sm.FirstName) > 0 {
		names = append(names, *sm.FirstName)
	}
	if sm.LastName != nil && len(*sm.LastName) > 0 {
		names = append(names, *sm.LastName)
	}
	if sm.MiddleName != nil && len(*sm.MiddleName) > 0 {
		names = append(names, *sm.MiddleName)
	}
	joined := strings.Join(names, ", ")

	s.pdf.Cell(nameFieldWidth-nameLabelWidth, fieldHeight*2, joined)

	s.setFieldLabelFont()
	s.pdf.Cell(bodyWidth-nameFieldWidth, fieldHeight, "Preferred Phone Number")
	s.pdf.SetXY(nameFieldWidth+horizontalMargin, s.pdf.GetY()+fieldHeight)
	s.setFieldValueFont()
	s.pdf.Cell(bodyWidth-nameFieldWidth, fieldHeight, coalesce(sm.Telephone, ""))
	s.pdf.Ln(-1)

	s.drawGrayLineFull(2)

	// More stuff
	row := []formField{
		formField{label: "DoD ID", value: coalesce(sm.Edipi, "")},
		formField{label: "Service Branch/Agency", value: ""},
		formField{label: "Rank/Grade", value: ""},
		formField{label: "Preferred Email", value: ""},
	}
	s.addFormRow(row, bodyWidth)
	s.drawGrayLineFull(2)

	// Address
	s.setFieldLabelFont()
	s.pdf.Cell(bodyWidth*0.3, fieldHeight, "Preferred W2 Mailing Address")
	s.setFieldValueFont()
	s.pdf.Cell(bodyWidth*0.7, fieldHeight, "")
	s.pdf.Ln(-1)

	// Not the right data
	s.addSectionHeader("ORDERS/ACCOUNTING INFORMATION")
	s.addFormRow(row, bodyWidth)

	s.addSectionHeader("ENTITLEMENTS/MOVE SUMMARY")
	y := s.pdf.GetY()
	entitlements := []formField{
		formField{label: "Entitlement", value: "12321 lbs"},
		formField{label: "Pro-Gear", value: "12321 lbs"},
		formField{label: "Spouse Pro-Gear", value: "12321 lbs"},
		formField{label: "Total Weight", value: "12321 lbs"},
	}
	s.addTable("Maximum Weight Entitlement", entitlements, bodyWidth*0.46, fieldHeight)

	middleX := PdfPageWidth * 0.5
	s.pdf.SetXY(middleX, y)
	row = []formField{
		formField{label: "Authorized Origin", value: "Ft. Bragg"},
		formField{label: "Authorized Destination", value: "Pentagon"},
	}
	s.addFormRow(row, bodyWidth*0.5)
	s.drawGrayLine(2, middleX, PdfPageWidth-horizontalMargin)
	s.pdf.SetX(middleX)
	s.addFormRow(row, bodyWidth*0.5)
	s.drawGrayLine(2, middleX, PdfPageWidth-horizontalMargin)

	s.addSectionHeader("FINANCE/PAYMENT")

	return s.pdf.Output(outputFile)
}

func (s *ShipmentSummary) addSectionHeader(title string) {
	s.pdf.Ln(2)
	s.pdf.SetFont(fontFace, "B", 10)
	s.pdf.SetFillColor(221, 231, 240)
	s.pdf.CellFormat(0, 7, title, "", 1, "L", true, 0, "")
	s.pdf.Ln(1)
}

func (s *ShipmentSummary) setFieldLabelFont() {
	s.pdf.SetFont(fontFace, "B", 10)
}

func (s *ShipmentSummary) setFieldValueFont() {
	s.pdf.SetFont(fontFace, "B", 11)
}

func (s *ShipmentSummary) drawGrayLineFull(margin float64) {
	s.drawGrayLine(margin, horizontalMargin, PdfPageWidth-horizontalMargin)
}

func (s *ShipmentSummary) drawGrayLine(margin, x1, x2 float64) {
	s.pdf.SetDrawColor(221, 231, 240)
	s.pdf.SetLineWidth(0.2)
	s.pdf.Ln(margin)
	s.pdf.Line(x1, s.pdf.GetY(), x2, s.pdf.GetY())
	s.pdf.Ln(margin)
}

func (s *ShipmentSummary) addFormRow(fields []formField, width float64) {
	x := s.pdf.GetX()
	fieldWidth := width / float64(len(fields))

	// Add labels
	s.setFieldLabelFont()
	for _, field := range fields {
		s.pdf.Cell(fieldWidth, fieldHeight, field.label)
	}
	s.pdf.Ln(-1)
	s.pdf.SetX(x)

	// Add values
	s.setFieldValueFont()
	for _, field := range fields {
		s.pdf.Cell(fieldWidth, fieldHeight, field.value)
	}
	s.pdf.Ln(-1)
	s.pdf.SetX(x)
}

func (s *ShipmentSummary) addTable(header string, fields []formField, width, cellHeight float64) {
	s.pdf.SetFont(fontFace, "B", 10)
	s.pdf.CellFormat(width, cellHeight, header, "1", 1, "", false, 0, "")

	s.pdf.SetFont(fontFace, "", 10)
	for _, field := range fields {
		s.pdf.CellFormat(width/2, cellHeight, field.label, "LTB", 0, "", false, 0, "")
		s.pdf.CellFormat(width/2, cellHeight, field.value, "TRB", 1, "", false, 0, "")
	}
}
