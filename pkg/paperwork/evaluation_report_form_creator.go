package paperwork

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/models"
)

const (
	regularFontPath = "pkg/paperwork/formtemplates/PublicSans-Regular.ttf"
	boldFontPath    = "pkg/paperwork/formtemplates/PublicSans-Bold.ttf"
	arrowImagePath  = "pkg/paperwork/formtemplates/arrowright.png"
	// We can figure this out at runtime, but it adds complexity
	arrowImageFormat          = "png"
	arrowImageName            = "arrowright"
	mmPerPixel                = (letterWidthMm / 1204)
	pageHeightMm              = 279.4
	pageBottomMarginMm        = 10.0
	pageTopMarginMm           = 14.0
	pageSideMarginMm          = 48.0 * mmPerPixel
	reportHeadingFontSize     = 40.0 * mmPerPixel
	sectionHeadingFontSize    = 28.0 * mmPerPixel
	subsectionHeadingFontSize = 22.0 * mmPerPixel
	textFontSize              = 15.0 * mmPerPixel
	textSmallFontSize         = 13.0 * mmPerPixel
)

// TODO do we want to keep this or inline it?
func pxToMM(px float64) float64 {
	// 1204px is the px width of the designs
	return (letterWidthMm / 1204) * px
}

type TableRow struct {
	LeftFieldName  string
	LeftLabel      string
	RightFieldName string
	RightLabel     string
}

func getField(fieldName string, data interface{}) (string, error) {
	r := reflect.ValueOf(data)
	val := reflect.Indirect(r).FieldByName(fieldName).Interface()

	switch v := val.(type) {
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("only string type is supported")
	}
}

type EvaluationReportFormFiller struct {
	pdf      *gofpdf.Fpdf
	reportID string
}

func loadFont(pdf *gofpdf.Fpdf, family string, style string, path string) error {
	font, err := assets.Asset(path)
	if err != nil {
		return err
	}
	pdf.AddUTF8FontFromBytes(family, style, font)

	return pdf.Error()
}

func NewEvaluationReportFormFiller() (*EvaluationReportFormFiller, error) {
	pdf := gofpdf.New(pageOrientation, distanceUnit, pageSize, fontDir)
	pdf.SetMargins(pageSideMarginMm, pageTopMarginMm, pageSideMarginMm)
	pdf.SetAutoPageBreak(false, pageBottomMarginMm)
	pdf.AliasNbPages("")

	err := loadFont(pdf, "PublicSans", "", regularFontPath)
	if err != nil {
		return nil, err
	}
	err = loadFont(pdf, "PublicSans", "B", boldFontPath)
	if err != nil {
		return nil, err
	}
	pdf.SetFont("PublicSans", fontStyle, fontSize)

	return &EvaluationReportFormFiller{
		pdf: pdf,
	}, nil
}

func (f *EvaluationReportFormFiller) CreateShipmentReport(report models.EvaluationReport, violations models.PWSViolations, shipment models.MTOShipment, customer models.ServiceMember) error {
	err := f.loadArrowImage()
	if err != nil {
		return err
	}
	f.reportID = fmt.Sprintf("QA-%s", strings.ToUpper(report.ID.String()[:5]))

	f.pdf.AddPage()
	f.reportHeading("Shipment report", f.reportID, report.Move.Locator, *report.Move.ReferenceID)
	f.contactInformation(customer, report.OfficeUser)

	err = f.shipmentCard(shipment)
	if err != nil {
		return fmt.Errorf("draw shipment card error %w", err)
	}

	err = f.inspectionInformationSection(report, violations)
	if err != nil {
		return err
	}

	f.pdf.AddPage()
	f.sectionHeading("Violations", pxToMM(56.0))

	err = f.violationsSection(violations)
	if err != nil {
		return err
	}

	return f.pdf.Error()
}

func (f *EvaluationReportFormFiller) CreateCounselingReport(report models.EvaluationReport, violations models.PWSViolations, shipments models.MTOShipments, customer models.ServiceMember) error {
	err := f.loadArrowImage()
	if err != nil {
		return err
	}

	f.reportID = fmt.Sprintf("QA-%s", strings.ToUpper(report.ID.String()[:5]))
	f.pdf.AddPage()

	f.reportHeading("Counseling report", f.reportID, report.Move.Locator, *report.Move.ReferenceID)
	f.sectionHeading("Move information", pxToMM(18.0))
	f.contactInformation(customer, report.OfficeUser)

	for _, shipment := range shipments {
		err = f.shipmentCard(shipment)
		if err != nil {
			return fmt.Errorf("draw shipment card error %w", err)
		}
	}

	f.pdf.AddPage()
	f.sectionHeading("Evaluation report", pxToMM(56.0))
	err = f.inspectionInformationSection(report, violations)
	if err != nil {
		return err
	}

	f.pdf.AddPage()
	f.sectionHeading("Violations", pxToMM(56.0))

	err = f.violationsSection(violations)
	if err != nil {
		return err
	}

	return f.pdf.Error()
}

// Output outputs the form to the provided file
func (f *EvaluationReportFormFiller) Output(output io.Writer) error {
	// Loop through all pages and add headings. This must be done right before output
	// because we need to be able to calculate the number of pages to show "Page X of Y"
	numPages := f.pdf.PageCount()
	for i := 1; i <= numPages; i++ {
		f.pdf.SetPage(i)
		f.reportPageHeader()
	}
	return f.pdf.Output(output)
}
func (f *EvaluationReportFormFiller) violationsSection(violations models.PWSViolations) error {
	f.subsectionHeading(fmt.Sprintf("Violations observed (%d)", len(violations)))

	kpis := map[string]bool{}
	for _, violation := range violations {
		if violation.IsKpi {
			if violation.AdditionalDataElem == "observedPickupSpreadDates" {
				kpis["ObservedPickupSpreadStartDate"] = true
				kpis["ObservedPickupSpreadEndDate"] = true
			} else {
				elementName := violation.AdditionalDataElem
				kpis[strings.ToUpper(elementName[0:1])+elementName[1:]] = true
			}
		}
		f.violation(violation)
		f.addVerticalSpace(pxToMM(16.0))
	}

	if len(kpis) > 0 {
		allKPIs := []string{}
		for kpi, present := range kpis {
			if present {
				allKPIs = append(allKPIs, kpi)
			}
		}
		err := f.subsection("Additional data for KPIs", allKPIs, KPIFieldLabels, AdditionalKPIData{
			ObservedPickupSpreadStartDate: "?",
			ObservedPickupSpreadEndDate:   "?",
			ObservedClaimDate:             "?",
			ObservedPickupDate:            "?",
			ObservedDeliveryDate:          "?",
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *EvaluationReportFormFiller) inspectionInformationSection(report models.EvaluationReport, violations models.PWSViolations) error {
	inspectionInfo := FormatValuesInspectionInformation(report, violations)

	err := f.subsection("Inspection information", InspectionInformationFields, InspectionInformationFieldLabels, inspectionInfo)
	if err != nil {
		return err
	}

	err = f.subsection("Violations", ViolationsFields, ViolationsFieldLabels, inspectionInfo)
	if err != nil {
		return err
	}

	err = f.subsection("QAE remarks", QAERemarksFields, QAERemarksFieldLabels, inspectionInfo)
	if err != nil {
		return err
	}
	return nil
}

func (f *EvaluationReportFormFiller) addVerticalSpace(dy float64) {
	f.pdf.SetY(f.pdf.GetY() + dy)
}

func (f *EvaluationReportFormFiller) reportPageHeader() {
	stripeHeight := pxToMM(18.0)
	textHeight := pxToMM(34.0)

	f.pdf.MoveTo(0.0, 0.0)
	f.pdf.SetTextColor(255, 255, 255)
	f.pdf.SetFillColor(0, 0, 0)
	f.pdf.SetFontUnitSize(textSmallFontSize) // 28px
	f.pdf.CellFormat(letterWidthMm, stripeHeight, controlledUnclassifiedInformationText, "", 1, "CM", true, 0, "")
	f.setTextColorBaseDarker()
	f.pdf.SetFontStyle("B")
	f.pdf.CellFormat(-pageSideMarginMm+letterWidthMm/2.0, textHeight, fmt.Sprintf("Report #%s", f.reportID), "", 0, "LM", false, 0, "")
	f.pdf.CellFormat(0.0, textHeight, fmt.Sprintf("Page %d of %d", f.pdf.PageNo(), f.pdf.PageCount()), "", 1, "RM", false, 0, "")
	f.pdf.SetFontStyle("")
}

func (f *EvaluationReportFormFiller) reportHeading(text string, reportID string, moveCode string, mtoReferenceID string) {
	//f.pdf.SetFontSize(30.0) // 40px
	f.pdf.SetFontUnitSize(reportHeadingFontSize)
	f.setTextColorBaseDarkest()
	headingY := f.pdf.GetY()
	headingWidth := 100.0
	height := pxToMM(70.0)
	bottomMargin := pxToMM(43.0)
	rightMargin := pageSideMarginMm
	idsWidth := letterWidthMm - pageSideMarginMm - headingWidth - rightMargin

	// Heading (left aligned)
	f.pdf.MoveTo(pageSideMarginMm, headingY)
	f.pdf.CellFormat(headingWidth, height, text, "", 0, "LM", false, 0, "")

	// Report ID/Move Code/MTO reference ID (right aligned)
	f.pdf.SetFontUnitSize(textSmallFontSize)
	f.setTextColorBaseDark()
	f.pdf.CellFormat(idsWidth, height/3.0, fmt.Sprintf("REPORT ID #%s", reportID), "", 1, "RM", false, 0, "")
	f.pdf.SetX(pageSideMarginMm + headingWidth)
	f.pdf.CellFormat(idsWidth, height/3.0, fmt.Sprintf("MOVE CODE #%s", moveCode), "", 1, "RM", false, 0, "")
	f.pdf.SetX(pageSideMarginMm + headingWidth)
	f.pdf.CellFormat(idsWidth, height/3.0, fmt.Sprintf("MTO REFERENCE ID #%s", mtoReferenceID), "", 1, "RM", false, 0, "")
	f.pdf.MoveTo(pageSideMarginMm, headingY+height)
	f.addVerticalSpace(bottomMargin)
}
func (f *EvaluationReportFormFiller) setTextColorBaseDark() {
	f.pdf.SetTextColor(86, 92, 101)
}
func (f *EvaluationReportFormFiller) setTextColorBaseDarker() {
	f.pdf.SetTextColor(61, 69, 81)
}
func (f *EvaluationReportFormFiller) setTextColorBaseDarkest() {
	f.pdf.SetTextColor(23, 23, 23)
}
func (f *EvaluationReportFormFiller) setFillColorBaseLight() {
	f.pdf.SetFillColor(169, 174, 177)
}

func (f *EvaluationReportFormFiller) setBorderColor() {
	borderR := 220
	borderG := 222
	borderB := 224

	f.pdf.SetDrawColor(borderR, borderG, borderB)
}

func (f *EvaluationReportFormFiller) sectionHeading(text string, bottomMargin float64) {
	f.pdf.SetFontStyle("B")
	f.setTextColorBaseDarkest()
	f.pdf.SetFontUnitSize(sectionHeadingFontSize)

	f.pdf.SetX(pageSideMarginMm)
	f.pdf.CellFormat(0.0, 10.0, text, "", 1, "LT", false, 0, "")
	f.pdf.SetFontStyle("")
	f.pdf.SetFontSize(fontSize)

	f.addVerticalSpace(bottomMargin)
}

func (f *EvaluationReportFormFiller) subsection(heading string, fieldOrder []string, fieldLabels map[string]string, data interface{}) error {
	bottomMargin := pxToMM(40.0)
	f.subsectionHeading(heading)
	for _, field := range fieldOrder {
		labelText, found := fieldLabels[field]
		if !found {
			return fmt.Errorf("not found %s", field)
		}
		fieldValue, err := getField(field, data)
		if err != nil {
			return err
		}
		if fieldValue != "" {
			f.subsectionRow(labelText, fieldValue)
		}
	}
	f.addVerticalSpace(bottomMargin)

	return nil
}

func (f *EvaluationReportFormFiller) subsectionHeading(heading string) {
	topMargin := pxToMM(16.0)
	bottomMargin := pxToMM(24.0)
	f.pdf.SetFontStyle("B")
	f.setTextColorBaseDarkest()
	f.pdf.SetFontUnitSize(subsectionHeadingFontSize)
	f.addVerticalSpace(topMargin)
	f.pdf.SetX(pageSideMarginMm)
	f.pdf.CellFormat(0.0, 10.0, heading, "", 1, "LT", false, 0, "")
	f.addVerticalSpace(bottomMargin)

	// Reset font
	f.pdf.SetFontStyle("")
	f.pdf.SetFontSize(fontSize)
}

// Assumptions: we wont have long enough labels to want auto page break
func (f *EvaluationReportFormFiller) subsectionRow(key string, value string) {
	f.pdf.SetX(pageSideMarginMm)
	f.setTextColorBaseDarkest()
	f.pdf.SetFontStyle("B")
	f.pdf.SetCellMargin(pxToMM(8.0))
	f.setBorderColor()
	labelWidth := pxToMM(200.0)
	valueWidth := letterWidthMm - 2.0*pageSideMarginMm - labelWidth
	textLineHeight := pxToMM(18.0)
	minFieldHeight := pxToMM(40.0)

	// If the text is long, or contains line breaks, we will want to display across multiple lines
	needToLineWrapValue := f.pdf.GetStringWidth(value) > valueWidth-2*f.pdf.GetCellMargin() || strings.Contains(value, "\n")
	// I'm assuming that we will not have line breaks in labels
	needToLineWrapLabel := f.pdf.GetStringWidth(key) > labelWidth-2*f.pdf.GetCellMargin()
	estimatedHeight := minFieldHeight
	// TODO the estimated height calculation is not quite right, it diverges for really long text.
	// using AutoPageBreak prevents this from being an issue, but it is weird.
	if needToLineWrapValue {
		// Auto page break doesnt work super well for us in other places in the document because we have lines that
		// should be kept together, but here, for a potentially large block of paragraphy text, it works great.
		f.pdf.SetAutoPageBreak(true, pageBottomMarginMm)
		estimatedHeight = math.Ceil(f.pdf.GetStringWidth(value)/(valueWidth-2*f.pdf.GetCellMargin())) * textLineHeight
	}
	if needToLineWrapLabel {
		estimatedHeight = math.Max(estimatedHeight, math.Ceil(f.pdf.GetStringWidth(key)/(labelWidth-2*f.pdf.GetCellMargin()))*textLineHeight)
	}
	if f.pdf.GetY()+estimatedHeight > pageHeightMm-pageBottomMarginMm {
		f.pdf.AddPage()
	}
	f.pdf.SetFontUnitSize(textFontSize)
	y := f.pdf.GetY()

	if needToLineWrapLabel {
		f.pdf.MultiCell(labelWidth, textLineHeight, key, "T", "LM", false)
	} else {
		f.pdf.CellFormat(labelWidth, minFieldHeight, key, "T", 0, "LM", false, 0, "")
	}

	labelY := f.pdf.GetY()
	f.pdf.SetFontStyle("")
	f.pdf.MoveTo(pageSideMarginMm+labelWidth, y)
	if needToLineWrapValue {
		f.pdf.MultiCell(valueWidth, textLineHeight, value, "T", "LM", false)
	} else {
		f.pdf.CellFormat(valueWidth, minFieldHeight, value, "T", 1, "LM", false, 0, "")
	}
	valueY := f.pdf.GetY()
	endY := math.Max(math.Max(labelY, valueY), y+minFieldHeight)
	f.pdf.SetY(endY)
	f.pdf.SetAutoPageBreak(false, pageBottomMarginMm)
}

func (f *EvaluationReportFormFiller) violation(violation models.PWSViolation) {
	// - 1.2.3 Violation Title
	//   Requirement summary
	height := pxToMM(18.0)
	bulletWidth := pxToMM(22.0)
	f.pdf.SetX(pageSideMarginMm)
	f.pdf.SetFontUnitSize(textSmallFontSize)
	f.pdf.SetFontStyle("B")

	totalHeight := 2 * height
	if f.pdf.GetY()+totalHeight > pageHeightMm-pageBottomMarginMm {
		f.pdf.AddPage()
	}
	f.pdf.CellFormat(bulletWidth, height, "•", "", 0, "RM", false, 0, "")
	f.pdf.CellFormat(letterWidthMm-2.0*pageSideMarginMm-bulletWidth, height, violation.ParagraphNumber+" "+violation.Title, "", 1, "LM", false, 0, "")
	f.pdf.SetX(pageSideMarginMm + bulletWidth)
	f.pdf.SetFontStyle("")
	f.pdf.CellFormat(letterWidthMm-2.0*pageSideMarginMm, height, violation.RequirementSummary, "", 1, "LM", false, 0, "")
}

// contactInformation displays side by side contact info for customer and QAE users
func (f *EvaluationReportFormFiller) contactInformation(customer models.ServiceMember, officeUser models.OfficeUser) {
	contactInfo := FormatContactInformationValues(customer, officeUser)

	gap := pxToMM(16.0)
	columnWidth := -pageSideMarginMm + (letterWidthMm-gap)/2.0
	textHeight := pxToMM(21.0)
	customerContactText := strings.Join([]string{contactInfo.CustomerFullName, contactInfo.CustomerPhone, contactInfo.CustomerRank, contactInfo.CustomerAffiliation}, "\n")
	qaeContactText := strings.Join([]string{contactInfo.QAEFullName, contactInfo.QAEPhone, contactInfo.QAEEmail}, "\n")

	f.pdf.SetFontStyle("B")
	f.pdf.SetFontUnitSize(textFontSize)
	f.setTextColorBaseDarkest()
	f.setBorderColor()
	f.pdf.SetX(pageSideMarginMm)
	f.pdf.CellFormat(columnWidth, textHeight, "Customer information", "B", 0, "LM", false, 0, "")
	f.pdf.SetX(pageSideMarginMm + columnWidth + gap)
	f.pdf.CellFormat(columnWidth, textHeight, "QAE", "B", 1, "LM", false, 0, "")
	f.pdf.SetFontStyle("")
	contentY := f.pdf.GetY()
	f.pdf.MultiCell(columnWidth, textHeight, customerContactText, "", "LM", false)
	endY := f.pdf.GetY()
	f.pdf.MoveTo(pageSideMarginMm+columnWidth+gap, contentY)
	f.pdf.MultiCell(columnWidth, textHeight, qaeContactText, "", "LM", false)

	bottomMargin := pxToMM(36.0)
	f.pdf.MoveTo(pageSideMarginMm, endY+bottomMargin)
}

/*
==============================================================================
| HHG                                                     SHIPMENT ID #12345 |
|                                                                            |
| 123 Main St                        ->  456 Freedom Rd                      |
| Scheduled pickup date 20 May 2022      Scheduled delivery date 29 May 2022 |
| ---------------------------------      ----------------------------------- |
| Requested pickup date 20 May 2022      Requested delivery date 29 May 2022 |
| ---------------------------------      ----------------------------------- |
| .....                                                                      |
'.__________________________________________________________________________.'
*/
func (f *EvaluationReportFormFiller) shipmentCard(shipment models.MTOShipment) error {
	layout := PickShipmentCardLayout(shipment.ShipmentType)
	vals := FormatValuesShipment(shipment)
	headingMargin := 2.0
	headingHeight := 5.0
	headingBottomMargin := pxToMM(18.0)
	stripeHeight := pxToMM(9.0)
	addressHeight := 10.0
	tableRowHeight := 10.0
	estimatedHeight := stripeHeight + headingMargin + headingHeight + headingBottomMargin + addressHeight + tableRowHeight + float64(len(layout))
	if f.pdf.GetY()+estimatedHeight > pageHeightMm-pageBottomMarginMm {
		f.pdf.AddPage()
	}

	cardWidth := letterWidthMm - 2*pageSideMarginMm
	f.setHHGStripeColor(shipment.ShipmentType)
	f.pdf.Rect(pageSideMarginMm, f.pdf.GetY(), cardWidth, stripeHeight, "DF")
	startY := f.pdf.GetY() + stripeHeight

	headingX := pageSideMarginMm + pxToMM(8.0)
	headingY := startY + headingMargin
	shipmentTypeX := headingX

	f.pdf.MoveTo(shipmentTypeX, headingY)
	f.pdf.SetFontStyle("B")
	shipmentTypeText := f.formatShipmentType(shipment.ShipmentType)

	f.pdf.SetFontUnitSize(textFontSize)
	f.pdf.CellFormat(f.pdf.GetStringWidth(shipmentTypeText)+2*f.pdf.GetCellMargin(), headingHeight, shipmentTypeText, "", 0, "LM", false, 0, "")
	f.pdf.SetFontStyle("")
	if shipment.UsesExternalVendor {
		f.setFillColorBaseLight()
		const externalVendorText = "EXTERNAL VENDOR"
		vendorTagWidth := f.pdf.GetStringWidth(externalVendorText) + 2*f.pdf.GetCellMargin()
		f.pdf.CellFormat(vendorTagWidth, headingHeight, externalVendorText, "", 0, "LM", true, 0, "")

	}
	// heading - shipment ID
	f.setTextColorBaseDark()
	// pagewidth - x - margin
	shipmentIDWidth := ((pageSideMarginMm + cardWidth) - f.pdf.GetX()) - pxToMM(8.0)
	f.pdf.SetFontUnitSize(textSmallFontSize)
	f.pdf.CellFormat(shipmentIDWidth, headingHeight, "Shipment ID: "+vals.ShipmentID, "", 0, "RM", false, 0, "")
	f.addVerticalSpace(headingHeight + headingBottomMargin)

	tableHMargin := pxToMM(12.0)
	tableWidth := cardWidth - 2.0*tableHMargin
	tableX := pageSideMarginMm + tableHMargin
	if shipment.ShipmentType != models.MTOShipmentTypePPM {
		gap := 2.0
		labelWidth := 0.3 * ((tableWidth - gap) / 2.0)
		valueWidth := 0.7 * ((tableWidth - gap) / 2.0)
		rightX := tableX + labelWidth + valueWidth + gap
		leftAddressLabel := ""
		rightAddressLabel := ""
		if shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
			leftAddressLabel = vals.StorageFacilityName
			rightAddressLabel = "DELIVERY ADDRESS"
		} else if shipment.ShipmentType == models.MTOShipmentTypeHHGIntoNTSDom {
			leftAddressLabel = "PICKUP ADDRESS"
			rightAddressLabel = vals.StorageFacilityName
		}
		f.sideBySideAddress(gap, tableX, vals.PickupAddress, leftAddressLabel, rightX, vals.DeliveryAddress, rightAddressLabel)
		f.addVerticalSpace(pxToMM(12.0))
	}
	err := f.twoColumnTable(tableX, f.pdf.GetY(), tableWidth, layout, vals)
	if err != nil {
		return fmt.Errorf("TwoColumnTable %w", err)
	}
	f.pdf.RoundedRect(pageSideMarginMm, startY, cardWidth, f.pdf.GetY()-startY, 1.0, "34", "D")
	shipmentCardBottomMargin := pxToMM(16.0)
	f.addVerticalSpace(shipmentCardBottomMargin)
	return nil
}

func (f *EvaluationReportFormFiller) setHHGStripeColor(shipmentType models.MTOShipmentType) {
	r := 0
	g := 150
	b := 244

	if strings.Contains(string(shipmentType), "PPM") {
		r = 230
		g = 199
		b = 76
	} else if strings.Contains(string(shipmentType), "NTS") {
		r = 129
		g = 104
		b = 179
	} else if strings.Contains(string(shipmentType), "HHG") {
		r = 0
		g = 150
		b = 244
	}
	f.pdf.SetDrawColor(r, g, b)
	f.pdf.SetFillColor(r, g, b)
}

func (f *EvaluationReportFormFiller) sideBySideAddress(gap float64, leftAddressX float64, leftAddress string, leftAddressLabel string, rightAddressX float64, rightAddress string, rightAddressLabel string) {
	f.pdf.SetFontUnitSize(textFontSize)
	addressY := f.pdf.GetY()
	startY := f.pdf.GetY()
	f.pdf.SetX(leftAddressX)
	f.setTextColorBaseDark()
	addressWidth := rightAddressX - leftAddressX - gap
	if leftAddressLabel != "" {
		f.pdf.CellFormat(addressWidth, pxToMM(18.0), leftAddressLabel, "", 1, "LT", false, 0, "")
		addressY = f.pdf.GetY() + pxToMM(8.0)
	}
	if rightAddressLabel != "" {
		f.pdf.MoveTo(rightAddressX, startY)
		f.pdf.CellFormat(addressWidth, pxToMM(18.0), rightAddressLabel, "", 1, "LT", false, 0, "")
		addressY = math.Max(addressY, f.pdf.GetY()+pxToMM(8.0))
	}
	f.pdf.MoveTo(leftAddressX, addressY)
	f.setTextColorBaseDarkest()
	f.pdf.CellFormat(addressWidth, pxToMM(18.0), leftAddress, "", 1, "LT", false, 0, "")
	leftY := f.pdf.GetY()
	f.pdf.MoveTo(rightAddressX-pxToMM(20.0), addressY)
	f.drawArrow()

	f.pdf.MoveTo(rightAddressX, addressY)
	f.pdf.CellFormat(addressWidth, pxToMM(18.0), rightAddress, "", 1, "LT", false, 0, "")
	addressY = math.Max(leftY, f.pdf.GetY())
	f.pdf.SetY(addressY)
}

func (f *EvaluationReportFormFiller) formatShipmentType(shipmentType models.MTOShipmentType) string {
	if shipmentType == models.MTOShipmentTypePPM {
		return "PPM"
	} else if shipmentType == models.MTOShipmentTypeHHGIntoNTSDom {
		return "NTS"
	} else if shipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
		return "NTS-R"
	} else if strings.Contains(string(shipmentType), "HHG") {
		return "HHG"
	}
	return string(shipmentType)
}

func (f *EvaluationReportFormFiller) twoColumnTable(x float64, y float64, w float64, layout []TableRow, data interface{}) error {
	gap := pxToMM(28.0)
	columnWidth := (w - gap) / 2.0
	labelWidth := 0.3 * columnWidth
	valueWidth := 0.7 * columnWidth
	f.pdf.SetY(y)
	f.pdf.SetFontUnitSize(textSmallFontSize)

	for i, row := range layout {
		err := f.twoColumnTableRow(x, gap, labelWidth, valueWidth, row, data)
		if err != nil {
			return err
		}
		if i < len(layout)-1 {
			f.setBorderColor()
			f.pdf.Line(x, f.pdf.GetY(), x+labelWidth+valueWidth, f.pdf.GetY())
			f.pdf.Line(x+labelWidth+valueWidth+gap, f.pdf.GetY(), x+gap+2.0*(labelWidth+valueWidth), f.pdf.GetY())
			f.addVerticalSpace(1.0)
		}
	}
	f.addVerticalSpace(2.0)
	return nil
}

func (f *EvaluationReportFormFiller) twoColumnTableRow(x float64, gap float64, labelWidth float64, valueWidth float64, row TableRow, data interface{}) error {
	rowStartY := f.pdf.GetY()

	leftVal, err := getField(row.LeftFieldName, data)
	if err != nil {
		return err
	}
	f.tableColumn(x, labelWidth, valueWidth, row.LeftLabel, leftVal)
	leftValY := f.pdf.GetY()
	if row.RightFieldName == "" {
		return nil
	}
	f.pdf.SetY(rowStartY)
	rightVal, err := getField(row.RightFieldName, data)
	if err != nil {
		return err
	}
	f.tableColumn(x+labelWidth+valueWidth+gap, labelWidth, valueWidth, row.RightLabel, rightVal)
	rightValY := f.pdf.GetY()
	f.pdf.SetY(math.Max(leftValY, rightValY))
	return nil
}

// tableColumn draws one side of a two-column table row
func (f *EvaluationReportFormFiller) tableColumn(x float64, labelWidth float64, valueWidth float64, label string, value string) {
	lineHeight := 5.0 // TODO this shadows a global
	f.pdf.SetX(x)
	f.pdf.SetFontStyle("B")
	f.setTextColorBaseDarker()
	f.pdf.CellFormat(labelWidth, 10.0, label, "", 0, "LT", false, 0, "")
	f.pdf.SetFontStyle("")
	f.setTextColorBaseDarkest()
	f.pdf.MultiCell(valueWidth, lineHeight, value, "", "LT", false)
}

// drawArrow draws an image of an arrow. loadArrowImage MUST be called before this.
func (f *EvaluationReportFormFiller) drawArrow() {
	f.pdf.Image(arrowImageName, f.pdf.GetX(), f.pdf.GetY(), pxToMM(20.0), 0.0, flow, arrowImageFormat, imageLink, imageLinkURL)
}

func (f *EvaluationReportFormFiller) loadArrowImage() error {
	// load image from assets
	arrow, err := assets.Asset(arrowImagePath)
	if err != nil {
		return errors.Wrap(err, "could not load image asset")
	}
	arrowImage := bytes.NewReader(arrow)

	opt := gofpdf.ImageOptions{
		ImageType: arrowImageFormat,
		ReadDpi:   true,
	}

	// After the image is registered, we can use its name to draw it
	f.pdf.RegisterImageOptionsReader(arrowImageName, opt, arrowImage)
	return f.pdf.Error()
}
