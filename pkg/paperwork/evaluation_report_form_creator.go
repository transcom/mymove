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
	// We can figure the image format out at runtime, but it makes the code more complicated.
	// If this form ever uses more images, it would probably make more sense to do that.
	arrowImageFormat = "png"
	arrowImageName   = "arrowright"

	// The designs for this report are 1204px wide. This lets us convert pixel sizes to millimeter sizes
	mmPerPixel                = (letterWidthMm / 1204)
	pageHeightMm              = 279.4
	pageBottomMarginMm        = 52 * mmPerPixel
	pageTopMarginMm           = 90 * mmPerPixel
	pageSideMarginMm          = 48.0 * mmPerPixel
	reportHeadingFontSize     = 40.0 * mmPerPixel
	sectionHeadingFontSize    = 28.0 * mmPerPixel
	subsectionHeadingFontSize = 22.0 * mmPerPixel
	textFontSize              = 15.0 * mmPerPixel
	textSmallFontSize         = 13.0 * mmPerPixel

	// fpdf interprets a width of zero as "fill all the space up to the right margin"
	widthFill = 0.0
	moveRight = 0
	moveDown  = 1
)

// pxToMM converts pixels (design units) to millimeters (pdf units)
func pxToMM(px float64) float64 {
	return mmPerPixel * px
}

// getField gets a value from a struct based on the string value of the field name
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

func loadFont(pdf *gofpdf.Fpdf, family string, style string, path string) error {
	font, err := assets.Asset(path)
	if err != nil {
		return err
	}
	pdf.AddUTF8FontFromBytes(family, style, font)

	return pdf.Error()
}

// EvaluationReportFormFiller is used to create counseling and shipment evaluation reports
type EvaluationReportFormFiller struct {
	pdf      *gofpdf.Fpdf
	reportID string
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

func (f *EvaluationReportFormFiller) CreateShipmentReport(report models.EvaluationReport, violations models.ReportViolations, shipment models.MTOShipment, customer models.ServiceMember) error {
	err := f.loadArrowImage()
	if err != nil {
		return err
	}
	f.reportID = fmt.Sprintf("QA-%s", strings.ToUpper(report.ID.String()[:5]))

	f.pdf.AddPage()
	seriousIncident := false
	if report.SeriousIncident != nil {
		seriousIncident = *report.SeriousIncident
	}
	f.reportHeading("Shipment report", seriousIncident, f.reportID, report.Move.Locator, *report.Move.ReferenceID)
	f.contactInformation(customer, report.OfficeUser)

	err = f.shipmentCard(shipment)
	if err != nil {
		return fmt.Errorf("draw shipment card error %w", err)
	}

	err = f.inspectionInformationSection(report)
	if err != nil {
		return err
	}

	if len(violations) != 0 {
		f.pdf.AddPage()
		f.sectionHeading("Violations", pxToMM(56.0))

		additionalKPIData := FormatAdditionalKPIValues(report)
		err = f.violationsSection(violations, additionalKPIData)
		if err != nil {
			return err
		}
	}

	return f.pdf.Error()
}

func (f *EvaluationReportFormFiller) CreateCounselingReport(report models.EvaluationReport, violations models.ReportViolations, shipments models.MTOShipments, customer models.ServiceMember) error {
	err := f.loadArrowImage()
	if err != nil {
		return err
	}

	f.reportID = fmt.Sprintf("QA-%s", strings.ToUpper(report.ID.String()[:5]))
	f.pdf.AddPage()

	seriousIncident := false
	if report.SeriousIncident != nil {
		seriousIncident = *report.SeriousIncident
	}
	f.reportHeading("Counseling report", seriousIncident, f.reportID, report.Move.Locator, *report.Move.ReferenceID)
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
	err = f.inspectionInformationSection(report)
	if err != nil {
		return err
	}

	if len(violations) != 0 {
		f.pdf.AddPage()
		f.sectionHeading("Violations", pxToMM(56.0))

		additionalKPIData := FormatAdditionalKPIValues(report)
		err = f.violationsSection(violations, additionalKPIData)
		if err != nil {
			return err
		}
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

// reportPageHeader draws the header at the top of every page
// It looks a bit like this:
// ####### CONTROLLED UNCLASSIFIED INFORMATION #######
// Report #QA-12345                        Page 1 of 3
func (f *EvaluationReportFormFiller) reportPageHeader() {
	stripeHeight := pxToMM(18.0)
	textHeight := pxToMM(34.0)

	f.pdf.MoveTo(0.0, 0.0)
	f.pdf.SetTextColor(162, 214, 61)
	f.pdf.SetFillColor(0, 0, 0)
	f.pdf.SetFontUnitSize(textSmallFontSize)
	f.pdf.CellFormat(letterWidthMm, stripeHeight, controlledUnclassifiedInformationText, "", moveDown, "CM", true, 0, "")
	f.setTextColorBaseDarker()
	f.pdf.SetFontStyle("B")
	f.pdf.CellFormat(-pageSideMarginMm+letterWidthMm/2.0, textHeight, fmt.Sprintf("Report #%s", f.reportID), "", moveRight, "LM", false, 0, "")
	f.pdf.CellFormat(widthFill, textHeight, fmt.Sprintf("Page %d of %d", f.pdf.PageNo(), f.pdf.PageCount()), "", moveDown, "RM", false, 0, "")
	f.pdf.SetFontStyle("")
}

// reportHeading draws the heading at the beginning of a report
// It looks a bit like this:
//
//	                                         REPORT ID #QA-12345
//	Counseling report                          MOVE CODE #ABC123
//	                                 MTO REFERENCE ID #1234-5678
func (f *EvaluationReportFormFiller) reportHeading(text string, seriousIncident bool, reportID string, moveCode string, mtoReferenceID string) {
	headingY := f.pdf.GetY()
	if seriousIncident {
		f.seriousIncidentFlag()
	}
	f.pdf.SetFontUnitSize(reportHeadingFontSize)
	f.setTextColorBaseDarkest()
	headingWidth := pxToMM(900)
	height := pxToMM(70.0)
	bottomMargin := pxToMM(43.0)

	// Heading (left aligned)
	f.pdf.SetX(pageSideMarginMm)
	if seriousIncident {
		f.pdf.CellFormat(headingWidth, reportHeadingFontSize+pxToMM(16.0), text, "", moveRight, "LM", false, 0, "")
		f.pdf.MoveTo(f.pdf.GetX(), headingY)
	} else {
		f.pdf.CellFormat(headingWidth, height, text, "", moveRight, "LM", false, 0, "")
	}

	// Report ID/Move Code/MTO reference ID (right aligned)
	f.pdf.SetFontUnitSize(textSmallFontSize)
	f.setTextColorBaseDark()
	f.pdf.CellFormat(widthFill, height/3.0, fmt.Sprintf("REPORT ID #%s", reportID), "", moveDown, "RM", false, 0, "")
	f.pdf.SetX(pageSideMarginMm + headingWidth)
	f.pdf.CellFormat(widthFill, height/3.0, fmt.Sprintf("MOVE CODE #%s", moveCode), "", moveDown, "RM", false, 0, "")
	f.pdf.SetX(pageSideMarginMm + headingWidth)
	f.pdf.CellFormat(widthFill, height/3.0, fmt.Sprintf("MTO REFERENCE ID #%s", mtoReferenceID), "", moveDown, "RM", false, 0, "")
	f.pdf.MoveTo(pageSideMarginMm, headingY+height)
	f.addVerticalSpace(bottomMargin)
}

func (f *EvaluationReportFormFiller) seriousIncidentFlag() {
	f.setErrorFlagColors()
	// bump the tag over so the left edge of the bubble lines up with the left edge
	// of text below it
	f.pdf.SetX(f.pdf.GetX() + f.pdf.GetCellMargin())
	f.drawTag("SERIOUS INCIDENT", moveDown)
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
	f.pdf.CellFormat(columnWidth, textHeight+pxToMM(8.0), "Customer information", "B", moveRight, "LM", false, 0, "")
	f.pdf.SetX(pageSideMarginMm + columnWidth + gap)
	f.pdf.CellFormat(columnWidth, textHeight+pxToMM(8.0), "QAE", "B", moveDown, "LM", false, 0, "")
	f.pdf.SetFontStyle("")
	f.addVerticalSpace(pxToMM(4.0))
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
	headingMargin := pxToMM(12.0)
	headingHeight := pxToMM(21.0)
	headingBottomMargin := pxToMM(18.0)
	stripeHeight := pxToMM(9.0)
	// Rough overestimates of subcomponent heights used to guess whether we should page break
	addressHeight := pxToMM(8.0) + 2.0*pxToMM(18.0)
	tableRowHeight := pxToMM(2.0*12.0 + 42.0)
	estimatedHeight := stripeHeight + headingMargin + headingHeight + headingBottomMargin + addressHeight + tableRowHeight*float64(len(layout))
	if f.pdf.GetY()+estimatedHeight > pageHeightMm-pageBottomMarginMm {
		f.pdf.AddPage()
	}

	// Colored stripe at top of card
	cardWidth := letterWidthMm - 2*pageSideMarginMm
	f.setHHGStripeColor(shipment.ShipmentType)
	f.pdf.Rect(pageSideMarginMm, f.pdf.GetY(), cardWidth, stripeHeight, "DF")
	startY := f.pdf.GetY() + stripeHeight

	// Shipment type (HHG, PPM, NTS, ...)
	headingX := pageSideMarginMm + pxToMM(8.0)
	headingY := startY + headingMargin
	f.pdf.MoveTo(headingX, headingY)
	f.pdf.SetFontStyle("B")
	shipmentTypeText := f.formatShipmentType(shipment.ShipmentType)

	f.pdf.SetFontUnitSize(textFontSize)
	f.pdf.CellFormat(f.pdf.GetStringWidth(shipmentTypeText)+2*f.pdf.GetCellMargin(), headingHeight, shipmentTypeText, "", moveRight, "LM", false, 0, "")
	f.pdf.SetFontStyle("")
	if shipment.UsesExternalVendor {
		f.setFillColorBaseLight()
		f.drawTag("EXTERNAL VENDOR", moveRight)
	}
	// heading - shipment ID (right aligned)
	f.setTextColorBaseDark()
	shipmentIDWidth := ((pageSideMarginMm + cardWidth) - f.pdf.GetX()) - pxToMM(8.0)
	f.pdf.SetFontUnitSize(textSmallFontSize)
	f.pdf.CellFormat(shipmentIDWidth, headingHeight, "Shipment ID: "+vals.ShipmentID, "", moveRight, "RM", false, 0, "")
	f.addVerticalSpace(headingHeight + headingBottomMargin)

	tableHMargin := pxToMM(12.0)
	tableWidth := cardWidth - 2.0*tableHMargin
	tableX := pageSideMarginMm + tableHMargin
	// Display pickup and destination addresses for non-PPM shipments
	if shipment.ShipmentType != models.MTOShipmentTypePPM {
		gap := pxToMM(48.0)
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
		return err
	}

	// Draw a border around the entire shipment card
	f.pdf.RoundedRect(pageSideMarginMm, startY, cardWidth, f.pdf.GetY()-startY, 1.0, "34", "D")
	shipmentCardBottomMargin := pxToMM(16.0)
	f.addVerticalSpace(shipmentCardBottomMargin)
	return f.pdf.Error()
}

// inspectionInformationSection draws the Inspection Information section of the report
func (f *EvaluationReportFormFiller) inspectionInformationSection(report models.EvaluationReport) error {
	inspectionInfo := FormatValuesInspectionInformation(report)

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
	return f.pdf.Error()
}

// violationsSection draws the violations section of the report, which lists all PWS violations and
// associated KPIs
func (f *EvaluationReportFormFiller) violationsSection(violations models.ReportViolations, additionalKPIData AdditionalKPIData) error {
	f.subsectionHeading(fmt.Sprintf("Violations observed (%d)", len(violations)))

	kpis := map[string]bool{}
	for _, reportViolation := range violations {
		violation := reportViolation.Violation
		if violation.IsKpi {
			// Save all the KPI fields that we'll need to display after the violations
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
		err := f.subsection("Additional data for KPIs", allKPIs, KPIFieldLabels, additionalKPIData)
		if err != nil {
			return err
		}
	}
	return f.pdf.Error()
}

func (f *EvaluationReportFormFiller) sectionHeading(text string, bottomMargin float64) {
	f.pdf.SetFontStyle("B")
	f.setTextColorBaseDarkest()
	f.pdf.SetFontUnitSize(sectionHeadingFontSize)

	f.pdf.SetX(pageSideMarginMm)
	f.pdf.CellFormat(widthFill, pxToMM(34.0), text, "", moveDown, "LT", false, 0, "")
	f.pdf.SetFontStyle("")
	f.pdf.SetFontSize(fontSize)

	f.addVerticalSpace(bottomMargin)
}

// subsection draws a heading and a series of fields
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

	return f.pdf.Error()
}

func (f *EvaluationReportFormFiller) subsectionHeading(heading string) {
	topMargin := pxToMM(16.0)
	bottomMargin := pxToMM(24.0)
	f.pdf.SetFontStyle("B")
	f.setTextColorBaseDarkest()
	f.pdf.SetFontUnitSize(subsectionHeadingFontSize)
	f.addVerticalSpace(topMargin)
	safetyMargin := 2.0 * subsectionHeadingFontSize
	if f.pdf.GetY()+subsectionHeadingFontSize+bottomMargin+safetyMargin > pageHeightMm-pageBottomMarginMm {
		// Bump heading to next page if we're too close to the bottom
		f.pdf.AddPage()
	}
	f.pdf.SetX(pageSideMarginMm)
	f.pdf.CellFormat(widthFill, pxToMM(26.0), heading, "", moveDown, "LT", false, 0, "")
	f.addVerticalSpace(bottomMargin)

	// Reset font
	f.pdf.SetFontStyle("")
	f.pdf.SetFontSize(fontSize)
}

// subsectionRow draws one field label and value
func (f *EvaluationReportFormFiller) subsectionRow(key string, value string) {
	f.pdf.SetX(pageSideMarginMm)
	f.setTextColorBaseDarkest()
	f.pdf.SetFontStyle("B")
	f.pdf.SetFontUnitSize(textFontSize)
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
	y := f.pdf.GetY()

	// border line
	f.pdf.Line(f.pdf.GetX(), y, f.pdf.GetX()+letterWidthMm-2.0*pageSideMarginMm, y)

	fieldInternalPadding := pxToMM(12.0)
	if needToLineWrapLabel {
		f.addVerticalSpace(fieldInternalPadding)
		f.pdf.MultiCell(labelWidth, textLineHeight, key, "", "LM", false)
		f.addVerticalSpace(fieldInternalPadding)
	} else {
		f.pdf.CellFormat(labelWidth, minFieldHeight, key, "", moveRight, "LM", false, 0, "")
	}

	labelY := f.pdf.GetY()
	f.pdf.SetFontStyle("")
	f.pdf.MoveTo(pageSideMarginMm+labelWidth, y)
	if needToLineWrapValue {
		f.addVerticalSpace(fieldInternalPadding)
		f.pdf.MultiCell(widthFill, textLineHeight, value, "", "LM", false)
		f.addVerticalSpace(fieldInternalPadding)
	} else {
		f.pdf.CellFormat(widthFill, minFieldHeight, value, "", moveDown, "LM", false, 0, "")
	}
	valueY := f.pdf.GetY()
	// Figure out where our tallest text field stopped, so we can make sure the next
	// element gets drawn below it
	endY := math.Max(math.Max(labelY, valueY), y+minFieldHeight)
	f.pdf.SetY(endY)
	f.pdf.SetAutoPageBreak(false, pageBottomMarginMm)
}

// violation displays a PWS requirement that was violated
func (f *EvaluationReportFormFiller) violation(violation models.PWSViolation) {
	height := pxToMM(18.0)
	bulletWidth := pxToMM(22.0)
	f.pdf.SetX(pageSideMarginMm)
	f.pdf.SetFontUnitSize(textSmallFontSize)
	f.pdf.SetFontStyle("B")

	totalHeight := 2 * height
	if f.pdf.GetY()+totalHeight > pageHeightMm-pageBottomMarginMm {
		f.pdf.AddPage()
	}
	// bullet point
	f.pdf.CellFormat(bulletWidth, height, "•", "", moveRight, "RM", false, 0, "")
	// paragraph number and title
	f.pdf.CellFormat(widthFill, height, violation.ParagraphNumber+" "+violation.Title, "", moveDown, "LM", false, 0, "")

	// requirement summary
	f.pdf.SetX(pageSideMarginMm + bulletWidth)
	f.pdf.SetFontStyle("")
	f.pdf.CellFormat(widthFill, height, violation.RequirementSummary, "", moveDown, "LM", false, 0, "")
}

// sideBySideAddress draws a pickup address and a delivery address in one line for a shipment card
func (f *EvaluationReportFormFiller) sideBySideAddress(gap float64, leftAddressX float64, leftAddress string, leftAddressLabel string, rightAddressX float64, rightAddress string, rightAddressLabel string) {
	f.pdf.SetFontUnitSize(textFontSize)
	addressHeight := pxToMM(18.0)
	addressY := f.pdf.GetY()
	startY := f.pdf.GetY()
	f.pdf.SetX(leftAddressX)
	f.setTextColorBaseDark()
	addressWidth := rightAddressX - leftAddressX - gap
	if leftAddressLabel != "" {
		f.pdf.CellFormat(addressWidth, addressHeight, leftAddressLabel, "", moveDown, "LT", false, 0, "")
		addressY = f.pdf.GetY() + pxToMM(8.0)
	}
	if rightAddressLabel != "" {
		f.pdf.MoveTo(rightAddressX, startY)
		f.pdf.CellFormat(addressWidth, addressHeight, rightAddressLabel, "", moveDown, "LT", false, 0, "")
		addressY = math.Max(addressY, f.pdf.GetY()+pxToMM(8.0))
	}
	f.pdf.MoveTo(leftAddressX, addressY)
	f.setTextColorBaseDarkest()
	f.pdf.MultiCell(addressWidth, addressHeight, leftAddress, "", "LT", false)
	leftY := f.pdf.GetY()
	f.pdf.MoveTo(leftAddressX+addressWidth, addressY)
	f.drawArrow(gap)

	f.pdf.MoveTo(rightAddressX, addressY)
	f.pdf.MultiCell(addressWidth, addressHeight, rightAddress, "", "LT", false)
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

// twoColumnTable is used to display fields within a shipment card
// typically, the left-hand fields are related to pickup and the right-hand fields are related to delivery
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
	return f.pdf.Error()
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
		// Skip drawing right field if it doesn't exist
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
	return f.pdf.Error()
}

// tableColumn draws one side of a two-column table row
func (f *EvaluationReportFormFiller) tableColumn(x float64, labelWidth float64, valueWidth float64, label string, value string) {
	textVerticalMargin := pxToMM(12.0)
	f.pdf.MoveTo(x, f.pdf.GetY()+textVerticalMargin)
	f.pdf.SetFontStyle("B")
	f.setTextColorBaseDarker()
	f.pdf.CellFormat(labelWidth, pxToMM(42.0), label, "", moveRight, "LT", false, 0, "")
	f.pdf.SetFontStyle("")
	f.setTextColorBaseDarkest()
	if value == "" {
		value = "-"
	}
	f.pdf.MultiCell(valueWidth, pxToMM(18.0), value, "", "LT", false)
	f.addVerticalSpace(textVerticalMargin)
}

func (f *EvaluationReportFormFiller) addVerticalSpace(dy float64) {
	f.pdf.MoveTo(f.pdf.GetX(), f.pdf.GetY()+dy)
}

// drawArrow draws an image of an arrow. loadArrowImage MUST be called before this.
func (f *EvaluationReportFormFiller) drawArrow(width float64) {
	arrowWidth := pxToMM(13.33)
	centerX := f.pdf.GetX() + (width-arrowWidth)/2.0
	f.pdf.Image(arrowImageName, centerX, f.pdf.GetY(), arrowWidth, 0.0, flow, arrowImageFormat, imageLink, imageLinkURL)
}

// drawTag draws text on top of a rounded rectangle
// text color and fill color should be set before calling this function
// If ln is 0, the pdf cursor will move to the right after drawing
// If ln is 1, the pdf cursor will move to the next line
func (f *EvaluationReportFormFiller) drawTag(text string, ln int) {
	tagFontSize := pxToMM(16.0)
	f.pdf.SetFontUnitSize(tagFontSize)
	f.pdf.SetFontStyle("")

	sidePadding := pxToMM(8.0)
	topPadding := pxToMM(5.0)
	bottomPadding := pxToMM(3.0)
	textWidth := f.pdf.GetStringWidth(text) + 2*sidePadding
	height := tagFontSize + bottomPadding + topPadding
	startX, startY := f.pdf.GetXY()
	// moving the rectangle up by 3px and the text down 2px gives us a top padding of 5px and
	// makes the tag line up well if it's drawn on the same line as other text.
	f.pdf.RoundedRect(startX, startY-pxToMM(3.0), textWidth, height, pxToMM(2.0), "1234", "F")
	f.pdf.MoveTo(startX, startY+pxToMM(2.0))
	f.pdf.CellFormat(textWidth, tagFontSize, text, "", ln, "CM", false, 0, "")
	if ln == moveRight {
		// if we're moving right, we have to reset Y so other text will line up with stuff that came before this tag
		f.pdf.SetXY(f.pdf.GetX(), startY)
	}
}

// loadArrowImage loads a specific image into the PDF. This needs to be called once
// before we try to draw the image. If these reports ever require more than this one
// image, we should make this more generic.
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
	f.pdf.SetDrawColor(220, 222, 224)
}

func (f *EvaluationReportFormFiller) setErrorFlagColors() {
	f.pdf.SetTextColor(255, 255, 255)
	f.pdf.SetFillColor(181, 9, 9)
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
