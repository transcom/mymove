package paperwork

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/transcom/mymove/pkg/route"

	"github.com/transcom/mymove/pkg/rateengine"

	"github.com/transcom/mymove/pkg/unit"

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

	orders := s.move.Orders
	sm := orders.ServiceMember

	s.pdf.SetHeaderFunc(func() {
		s.pdf.SetFont(fontFace, "B", 17)
		s.pdf.Cell(bodyWidth*0.75, fieldHeight*2, "SHIPMENT SUMMARY WORKSHEET - PPM")
		s.setFieldLabelFont()
		s.pdf.Cell(bodyWidth*0.25, fieldHeight, "Date Prepared (YYYY-MM-DD)")
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

	s.pdf.Cell(nameFieldWidth-nameLabelWidth, fieldHeight*2, sm.ReverseNameLineFormat())

	s.setFieldLabelFont()
	s.pdf.Cell(bodyWidth-nameFieldWidth, fieldHeight, "Preferred Phone Number")
	s.pdf.SetXY(nameFieldWidth+horizontalMargin, s.pdf.GetY()+fieldHeight)
	s.setFieldValueFont()
	s.pdf.Cell(bodyWidth-nameFieldWidth, fieldHeight, coalesce(sm.Telephone, ""))
	s.pdf.Ln(-1)

	s.drawGrayLineFull(2)

	// More stuff
	var affiliation string
	if sm.Affiliation != nil {
		affiliation = string(*sm.Affiliation)
	}

	var rank string
	if sm.Rank != nil {
		rank = string(*sm.Rank)
	}

	row := []formField{
		formField{label: "DoD ID", value: coalesce(sm.Edipi, "")},
		formField{label: "Service Branch/Agency", value: affiliation},
		formField{label: "Rank/Grade", value: rank},
		formField{label: "Preferred Email", value: coalesce(sm.PersonalEmail, "")},
	}
	s.addFormRow(row, bodyWidth)
	s.drawGrayLineFull(2)

	// Address
	var address string
	if sm.ResidentialAddress != nil {
		address = sm.ResidentialAddress.LineFormat()
	}

	s.setFieldLabelFont()
	s.pdf.Cell(bodyWidth*0.3, fieldHeight, "Preferred W2 Mailing Address")
	s.setFieldValueFont()
	s.pdf.Cell(bodyWidth*0.7, fieldHeight, address)
	s.pdf.Ln(-1)

	// Not the right data
	s.addSectionHeader("ORDERS/ACCOUNTING INFORMATION")
	s.addFormRow(row, bodyWidth)

	s.addSectionHeader("ENTITLEMENTS/MOVE SUMMARY")
	y := s.pdf.GetY()

	var allotment models.WeightAllotment
	if sm.Rank != nil {
		allotment = models.GetWeightAllotment(*sm.Rank)
	}

	total := allotment.TotalWeightSelf + allotment.ProGearWeight + allotment.ProGearWeightSpouse
	entitlements := []formField{
		formField{label: "Entitlement", value: formatPounds(allotment.TotalWeightSelf)},
		formField{label: "Pro-Gear", value: formatPounds(allotment.ProGearWeight)},
		formField{label: "Spouse Pro-Gear", value: formatPounds(allotment.ProGearWeightSpouse)},
		formField{label: "Total Weight", value: formatPounds(total)},
	}
	s.addTable("Maximum Weight Entitlement", entitlements, bodyWidth*0.46, fieldHeight)

	middleX := PdfPageWidth * 0.5
	s.pdf.SetXY(middleX, y)
	row = []formField{
		formField{label: "Authorized Origin", value: sm.DutyStation.Name},
		formField{label: "Authorized Destination", value: orders.NewDutyStation.Name},
	}
	s.addFormRow(row, bodyWidth*0.5)
	s.drawGrayLine(2, middleX, PdfPageWidth-horizontalMargin)
	s.pdf.SetX(middleX)
	s.addFormRow(row, bodyWidth*0.5)
	s.drawGrayLine(2, middleX, PdfPageWidth-horizontalMargin)

	s.addSectionHeader("FINANCE/PAYMENT")

	return s.pdf.Output(outputFile)
}

func formatPounds(weight int) string {
	return fmt.Sprintf("%d lbs", weight)
}

func (s *ShipmentSummary) addSectionHeader(title string) {
	s.pdf.Ln(2)
	s.pdf.SetFont(fontFace, "B", 10)
	s.pdf.SetFillColor(221, 231, 240)
	s.pdf.CellFormat(0, 7, title, "", 1, "L", true, 0, "")
	s.pdf.Ln(1)
}

func (s *ShipmentSummary) setFieldLabelFont() {
	s.pdf.SetFont(fontFace, "B", 9)
}

func (s *ShipmentSummary) setFieldValueFont() {
	s.pdf.SetFont(fontFace, "", 10)
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

type ppmComputer interface {
	ComputePPMIncludingLHDiscount(weight unit.Pound, originZip5 string, destinationZip5 string, distanceMiles int, date time.Time, daysInSIT int) (cost rateengine.CostComputation, err error)
}

//SSWPPMComputer a rate engine wrapper with helper functions to simplify ppm cost calculations specific to shipment summary worksheet
type SSWPPMComputer struct {
	ppmComputer
}

//NewSSWPPMComputer creates a SSWPPMComputer
func NewSSWPPMComputer(PPMComputer ppmComputer) *SSWPPMComputer {
	return &SSWPPMComputer{ppmComputer: PPMComputer}
}

//ObligationType type corresponding to obligation sections of shipment summary worksheet
type ObligationType int

//ComputeObligations is helper function for computing the obligations section of the shipment summary worksheet
func (sswPpmComputer *SSWPPMComputer) ComputeObligations(ssfd models.ShipmentSummaryFormData, planner route.Planner) (obligation models.Obligations, err error) {
	firstPPM, err := sswPpmComputer.nilCheckPPM(ssfd)

	if err != nil {
		return models.Obligations{}, err
	}
	distanceMiles, err := planner.Zip5TransitDistance(*firstPPM.PickupPostalCode, *firstPPM.DestinationPostalCode)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating distance")
	}
	maxCost, err := sswPpmComputer.ComputePPMIncludingLHDiscount(
		unit.Pound(ssfd.TotalWeightAllotment),
		*firstPPM.PickupPostalCode,
		*firstPPM.DestinationPostalCode,
		distanceMiles,
		*firstPPM.ActualMoveDate,
		0,
	)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating PPM max obligations")
	}
	actualCost, err := sswPpmComputer.ComputePPMIncludingLHDiscount(
		unit.Pound(ssfd.PPMRemainingEntitlement),
		*firstPPM.PickupPostalCode,
		*firstPPM.DestinationPostalCode,
		distanceMiles,
		*firstPPM.ActualMoveDate,
		0,
	)
	if err != nil {
		return models.Obligations{}, errors.New("error calculating PPM actual obligations")
	}
	var actualSIT unit.Cents
	if firstPPM.TotalSITCost != nil {
		actualSIT = *firstPPM.TotalSITCost
	}
	if actualSIT > maxCost.SITMax {
		actualSIT = maxCost.SITMax
	}
	maxObligation := models.Obligation{Gcc: maxCost.GCC, SIT: maxCost.SITMax}
	actualObligation := models.Obligation{Gcc: actualCost.GCC, SIT: actualSIT}
	obligations := models.Obligations{MaxObligation: maxObligation, ActualObligation: actualObligation}
	return obligations, nil
}

func (sswPpmComputer *SSWPPMComputer) nilCheckPPM(ssfd models.ShipmentSummaryFormData) (models.PersonallyProcuredMove, error) {
	if len(ssfd.PersonallyProcuredMoves) == 0 {
		return models.PersonallyProcuredMove{}, errors.New("missing ppm")
	}
	firstPPM := ssfd.PersonallyProcuredMoves[0]
	if firstPPM.PickupPostalCode == nil || firstPPM.DestinationPostalCode == nil {
		return models.PersonallyProcuredMove{}, errors.New("missing required address parameter")
	}
	if firstPPM.ActualMoveDate == nil {
		return models.PersonallyProcuredMove{}, errors.New("missing required actual move date parameter")
	}
	return firstPPM, nil
}
