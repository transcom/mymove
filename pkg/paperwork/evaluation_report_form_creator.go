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

type DynamicFormFiller struct {
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

func NewDynamicFormFiller() (*DynamicFormFiller, error) {
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

	return &DynamicFormFiller{
		pdf: pdf,
	}, nil
}

func (d *DynamicFormFiller) loadArrowImage() error {
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
	d.pdf.RegisterImageOptionsReader(arrowImageName, opt, arrowImage)
	return d.pdf.Error()
}

// Output outputs the form to the provided file
func (d *DynamicFormFiller) Output(output io.Writer) error {
	d.addPageHeaders()
	return d.pdf.Output(output)
}
func (d *DynamicFormFiller) ViolationsSection(violations models.PWSViolations) error {
	d.subsectionHeading(fmt.Sprintf("Violations observed (%d)", len(violations)))

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
		d.violation(violation)
		d.addVerticalSpace(pxToMM(16.0))
	}

	if len(kpis) > 0 {
		allKPIs := []string{}
		for kpi, present := range kpis {
			if present {
				allKPIs = append(allKPIs, kpi)
			}
		}
		err := d.subsection("Additional data for KPIs", allKPIs, KPIFieldLabels, AdditionalKPIData{
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

func (d *DynamicFormFiller) InspectionInformationSection(report models.EvaluationReport, violations models.PWSViolations) error {
	inspectionInfo := FormatValuesInspectionInformation(report, violations)

	err := d.subsection("Inspection information", InspectionInformationFields, InspectionInformationFieldLabels, inspectionInfo)
	if err != nil {
		return err
	}

	err = d.subsection("Violations", ViolationsFields, ViolationsFieldLabels, inspectionInfo)
	if err != nil {
		return err
	}

	err = d.subsection("QAE remarks", QAERemarksFields, QAERemarksFieldLabels, inspectionInfo)
	if err != nil {
		return err
	}
	return nil
}
func (d *DynamicFormFiller) CreateShipmentReport(report models.EvaluationReport, violations models.PWSViolations, shipment models.MTOShipment, customer models.ServiceMember) error {
	err := d.loadArrowImage()
	if err != nil {
		return err
	}
	d.reportID = fmt.Sprintf("QA-%s", strings.ToUpper(report.ID.String()[:5]))

	d.pdf.AddPage()
	d.reportHeading("Shipment report", d.reportID, report.Move.Locator, *report.Move.ReferenceID)
	d.contactInformation(customer, report.OfficeUser)

	err = d.shipmentCard(shipment)
	if err != nil {
		return fmt.Errorf("draw shipment card error %w", err)
	}

	err = d.InspectionInformationSection(report, violations)
	if err != nil {
		return err
	}

	d.pdf.AddPage()
	d.sectionHeading("Violations", pxToMM(56.0))

	err = d.ViolationsSection(violations)
	if err != nil {
		return err
	}

	return d.pdf.Error()
}
func (d *DynamicFormFiller) CreateCounselingReport(report models.EvaluationReport, violations models.PWSViolations, shipments models.MTOShipments, customer models.ServiceMember) error {
	err := d.loadArrowImage()
	if err != nil {
		return err
	}

	d.reportID = fmt.Sprintf("QA-%s", strings.ToUpper(report.ID.String()[:5]))
	d.pdf.AddPage()

	d.reportHeading("Counseling report", d.reportID, report.Move.Locator, *report.Move.ReferenceID)
	d.sectionHeading("Move information", pxToMM(18.0))
	d.contactInformation(customer, report.OfficeUser)

	for _, shipment := range shipments {
		err = d.shipmentCard(shipment)
		if err != nil {
			return fmt.Errorf("draw shipment card error %w", err)
		}
	}

	d.pdf.AddPage()
	d.sectionHeading("Evaluation report", pxToMM(56.0))
	err = d.InspectionInformationSection(report, violations)
	if err != nil {
		return err
	}

	d.pdf.AddPage()
	d.sectionHeading("Violations", pxToMM(56.0))

	err = d.ViolationsSection(violations)
	if err != nil {
		return err
	}

	return d.pdf.Error()
}

func (d *DynamicFormFiller) addVerticalSpace(dy float64) {
	d.pdf.SetY(d.pdf.GetY() + dy)
}

func (d *DynamicFormFiller) reportPageHeader() {
	stripeHeight := pxToMM(18.0)
	textHeight := pxToMM(34.0)

	d.pdf.MoveTo(0.0, 0.0)
	d.pdf.SetTextColor(255, 255, 255)
	d.pdf.SetFillColor(0, 0, 0)
	d.pdf.SetFontUnitSize(textSmallFontSize) // 28px
	d.pdf.CellFormat(letterWidthMm, stripeHeight, controlledUnclassifiedInformationText, "", 1, "CM", true, 0, "")
	d.setTextColorBaseDarker()
	d.pdf.SetFontStyle("B")
	d.pdf.CellFormat(-pageSideMarginMm+letterWidthMm/2.0, textHeight, fmt.Sprintf("Report #%s", d.reportID), "", 0, "LM", false, 0, "")
	d.pdf.CellFormat(0.0, textHeight, fmt.Sprintf("Page %d of %d", d.pdf.PageNo(), d.pdf.PageCount()), "", 1, "RM", false, 0, "")
	d.pdf.SetFontStyle("")
}

func (d *DynamicFormFiller) reportHeading(text string, reportID string, moveCode string, mtoReferenceID string) {
	//d.pdf.SetFontSize(30.0) // 40px
	d.pdf.SetFontUnitSize(reportHeadingFontSize)
	d.setTextColorBaseDarkest()
	headingX := pageSideMarginMm
	headingY := d.pdf.GetY()
	headingWidth := 100.0
	height := pxToMM(70.0)
	bottomMargin := pxToMM(43.0)
	rightMargin := pageSideMarginMm
	idsWidth := letterWidthMm - headingX - headingWidth - rightMargin

	// Heading (left aligned)
	d.pdf.MoveTo(headingX, headingY)
	d.pdf.CellFormat(headingWidth, height, text, "", 0, "LM", false, 0, "")

	// Report ID/Move Code/MTO reference ID (right aligned)
	d.pdf.SetFontUnitSize(textSmallFontSize)
	d.setTextColorBaseDark()
	d.pdf.CellFormat(idsWidth, height/3.0, fmt.Sprintf("REPORT ID #%s", reportID), "", 1, "RM", false, 0, "")
	d.pdf.SetX(headingX + headingWidth)
	d.pdf.CellFormat(idsWidth, height/3.0, fmt.Sprintf("MOVE CODE #%s", moveCode), "", 1, "RM", false, 0, "")
	d.pdf.SetX(headingX + headingWidth)
	d.pdf.CellFormat(idsWidth, height/3.0, fmt.Sprintf("MTO REFERENCE ID #%s", mtoReferenceID), "", 1, "RM", false, 0, "")
	d.pdf.MoveTo(headingX, headingY+height)
	d.addVerticalSpace(bottomMargin)
}
func (d *DynamicFormFiller) setTextColorBaseDark() {
	d.pdf.SetTextColor(86, 92, 101)
}
func (d *DynamicFormFiller) setTextColorBaseDarker() {
	d.pdf.SetTextColor(61, 69, 81)
}
func (d *DynamicFormFiller) setTextColorBaseDarkest() {
	d.pdf.SetTextColor(23, 23, 23)
}
func (d *DynamicFormFiller) setFillColorBaseLight() {
	d.pdf.SetFillColor(169, 174, 177)
}

func (d *DynamicFormFiller) setBorderColor() {
	borderR := 220
	borderG := 222
	borderB := 224

	d.pdf.SetDrawColor(borderR, borderG, borderB)
}

func (d *DynamicFormFiller) sectionHeading(text string, bottomMargin float64) {
	d.pdf.SetFontStyle("B")
	d.setTextColorBaseDarkest()
	d.pdf.SetFontUnitSize(sectionHeadingFontSize)

	d.pdf.SetX(pageSideMarginMm)
	d.pdf.CellFormat(0.0, 10.0, text, "", 1, "LT", false, 0, "")
	d.pdf.SetFontStyle("")
	d.pdf.SetFontSize(fontSize)

	d.addVerticalSpace(bottomMargin)
}

// TODO would love a better name for this
// TODO also how do i send the data in?
// TODO - [{key,value},...]
// TODO - map key -> value (but then how do i order?)
// TODO - list of keys, map of key to key text, struct with keys that match

func (d *DynamicFormFiller) subsection(heading string, fieldOrder []string, fieldLabels map[string]string, data interface{}) error {
	bottomMargin := pxToMM(40.0)
	d.subsectionHeading(heading)
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
			d.subsectionRow(labelText, fieldValue)
		}
	}
	d.addVerticalSpace(bottomMargin)

	return nil
}

func (d *DynamicFormFiller) subsectionHeading(heading string) {
	topMargin := pxToMM(16.0)
	bottomMargin := pxToMM(24.0)
	d.pdf.SetFontStyle("B")
	d.setTextColorBaseDarkest()
	d.pdf.SetFontUnitSize(subsectionHeadingFontSize)
	d.addVerticalSpace(topMargin)
	d.pdf.SetX(pageSideMarginMm)
	d.pdf.CellFormat(0.0, 10.0, heading, "", 1, "LT", false, 0, "")
	d.addVerticalSpace(bottomMargin)

	// Reset font
	d.pdf.SetFontStyle("")
	d.pdf.SetFontSize(fontSize)
}

// Assumptions: we wont have long enough labels to want auto page break
func (d *DynamicFormFiller) subsectionRow(key string, value string) {
	d.pdf.SetX(pageSideMarginMm)
	d.setTextColorBaseDarkest()
	d.pdf.SetFontStyle("B")
	d.pdf.SetCellMargin(pxToMM(8.0))
	d.setBorderColor()
	labelWidth := pxToMM(200.0)
	valueWidth := letterWidthMm - 2.0*pageSideMarginMm - labelWidth
	textLineHeight := pxToMM(18.0)
	minFieldHeight := pxToMM(40.0)

	// If the text is long, or contains line breaks, we will want to display across multiple lines
	needToLineWrapValue := d.pdf.GetStringWidth(value) > valueWidth-2*d.pdf.GetCellMargin() || strings.Contains(value, "\n")
	// I'm assuming that we will not have line breaks in labels
	needToLineWrapLabel := d.pdf.GetStringWidth(key) > labelWidth-2*d.pdf.GetCellMargin()
	estimatedHeight := minFieldHeight
	// TODO the estimated height calculation is not quite right, it diverges for really long text.
	// using AutoPageBreak prevents this from being an issue, but it is weird.
	if needToLineWrapValue {
		// Auto page break doesnt work super well for us in other places in the document because we have lines that
		// should be kept together, but here, for a potentially large block of paragraphy text, it works great.
		d.pdf.SetAutoPageBreak(true, pageBottomMarginMm)
		estimatedHeight = math.Ceil(d.pdf.GetStringWidth(value)/(valueWidth-2*d.pdf.GetCellMargin())) * textLineHeight
	}
	if needToLineWrapLabel {
		estimatedHeight = math.Max(estimatedHeight, math.Ceil(d.pdf.GetStringWidth(key)/(labelWidth-2*d.pdf.GetCellMargin()))*textLineHeight)
	}
	if d.pdf.GetY()+estimatedHeight > pageHeightMm-pageBottomMarginMm {
		d.pdf.AddPage()
	}
	d.pdf.SetFontUnitSize(textFontSize)
	y := d.pdf.GetY()

	if needToLineWrapLabel {
		d.pdf.MultiCell(labelWidth, textLineHeight, key, "T", "LM", false)
	} else {
		d.pdf.CellFormat(labelWidth, minFieldHeight, key, "T", 0, "LM", false, 0, "")
	}

	labelY := d.pdf.GetY()
	d.pdf.SetFontStyle("")
	d.pdf.MoveTo(pageSideMarginMm+labelWidth, y)
	if needToLineWrapValue {
		d.pdf.MultiCell(valueWidth, textLineHeight, value, "T", "LM", false)
	} else {
		d.pdf.CellFormat(valueWidth, minFieldHeight, value, "T", 1, "LM", false, 0, "")
	}
	valueY := d.pdf.GetY()
	endY := math.Max(math.Max(labelY, valueY), y+minFieldHeight)
	d.pdf.SetY(endY)
	d.pdf.SetAutoPageBreak(false, pageBottomMarginMm)
}

func (d *DynamicFormFiller) violation(violation models.PWSViolation) {
	// - 1.2.3 Violation Title
	//   Requirement summary
	height := pxToMM(18.0)
	bulletWidth := pxToMM(22.0)
	d.pdf.SetX(pageSideMarginMm)
	d.pdf.SetFontUnitSize(textSmallFontSize)
	d.pdf.SetFontStyle("B")

	totalHeight := 2 * height
	if d.pdf.GetY()+totalHeight > pageHeightMm-pageBottomMarginMm {
		d.pdf.AddPage()
	}
	d.pdf.CellFormat(bulletWidth, height, "•", "", 0, "RM", false, 0, "")
	d.pdf.CellFormat(letterWidthMm-2.0*pageSideMarginMm-bulletWidth, height, violation.ParagraphNumber+" "+violation.Title, "", 1, "LM", false, 0, "")
	d.pdf.SetX(pageSideMarginMm + bulletWidth)
	d.pdf.SetFontStyle("")
	d.pdf.CellFormat(letterWidthMm-2.0*pageSideMarginMm, height, violation.RequirementSummary, "", 1, "LM", false, 0, "")
}

func (d *DynamicFormFiller) contactInformation(customer models.ServiceMember, officeUser models.OfficeUser) {
	contactInfo := FormatContactInformationValues(customer, officeUser)

	gap := pxToMM(16.0)
	columnWidth := -pageSideMarginMm + (letterWidthMm-gap)/2.0
	textHeight := pxToMM(21.0)
	customerContactText := strings.Join([]string{contactInfo.CustomerFullName, contactInfo.CustomerPhone, contactInfo.CustomerRank, contactInfo.CustomerAffiliation}, "\n")
	qaeContactText := strings.Join([]string{contactInfo.QAEFullName, contactInfo.QAEPhone, contactInfo.QAEEmail}, "\n")

	d.pdf.SetFontStyle("B")
	d.pdf.SetFontUnitSize(textFontSize)
	d.setTextColorBaseDarkest()
	d.setBorderColor()
	d.pdf.SetX(pageSideMarginMm)
	d.pdf.CellFormat(columnWidth, textHeight, "Customer information", "B", 0, "LM", false, 0, "")
	d.pdf.SetX(pageSideMarginMm + columnWidth + gap)
	d.pdf.CellFormat(columnWidth, textHeight, "QAE", "B", 1, "LM", false, 0, "")
	d.pdf.SetFontStyle("")
	contentY := d.pdf.GetY()
	d.pdf.MultiCell(columnWidth, textHeight, customerContactText, "", "LM", false)
	endY := d.pdf.GetY()
	d.pdf.MoveTo(pageSideMarginMm+columnWidth+gap, contentY)
	d.pdf.MultiCell(columnWidth, textHeight, qaeContactText, "", "LM", false)

	bottomMargin := pxToMM(36.0)
	d.pdf.MoveTo(pageSideMarginMm, endY+bottomMargin)
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
func (d *DynamicFormFiller) shipmentCard(shipment models.MTOShipment) error {
	layout := PickShipmentCardLayout(shipment.ShipmentType)
	vals := FormatValuesShipment(shipment)
	stripeHeight := pxToMM(9.0)
	addressHeight := 10.0
	tableRowHeight := 10.0
	//estimatedHeight := stripeHeight + headingMargin + headingHeight + headingBottomMargin + addressHeight + tableRowHeight + len(layout)
	estimatedHeight := stripeHeight + 5.0 + 5.0 + 5.0 + addressHeight + tableRowHeight + float64(len(layout))
	if d.pdf.GetY()+estimatedHeight > pageHeightMm-pageBottomMarginMm {
		d.pdf.AddPage()
	}

	cardWidth := letterWidthMm - 2*pageSideMarginMm
	d.setHHGStripeColor(shipment.ShipmentType)
	d.pdf.Rect(pageSideMarginMm, d.pdf.GetY(), cardWidth, stripeHeight, "DF")
	startY := d.pdf.GetY() + stripeHeight

	headingMargin := 2.0
	headingX := pageSideMarginMm + pxToMM(8.0)
	headingY := startY + headingMargin
	headingHeight := 5.0
	shipmentTypeX := headingX
	headingBottomMargin := pxToMM(18.0)

	// in/72 px * 25.4mm/in * 18px (25.4/72) = 6.35
	// 1204px / 8.5in = 141.65px/in * in/25.4mm = 5.576px/mm
	d.pdf.MoveTo(shipmentTypeX, headingY)
	d.pdf.SetFontStyle("B")
	shipmentTypeText := d.formatShipmentType(shipment.ShipmentType)

	d.pdf.SetFontUnitSize(textFontSize)
	d.pdf.CellFormat(d.pdf.GetStringWidth(shipmentTypeText)+2*d.pdf.GetCellMargin(), headingHeight, shipmentTypeText, "", 0, "LM", false, 0, "")
	d.pdf.SetFontStyle("")
	if shipment.UsesExternalVendor {
		d.setFillColorBaseLight()
		const externalVendorText = "EXTERNAL VENDOR"
		vendorTagWidth := d.pdf.GetStringWidth(externalVendorText) + 2*d.pdf.GetCellMargin()
		d.pdf.CellFormat(vendorTagWidth, headingHeight, externalVendorText, "", 0, "LM", true, 0, "")

	}
	// heading - shipment ID
	d.setTextColorBaseDark()
	// pagewidth - x - margin
	shipmentIDWidth := ((pageSideMarginMm + cardWidth) - d.pdf.GetX()) - pxToMM(8.0)
	d.pdf.SetFontUnitSize(textSmallFontSize)
	d.pdf.CellFormat(shipmentIDWidth, headingHeight, "Shipment ID: "+vals.ShipmentID, "", 0, "RM", false, 0, "")
	d.addVerticalSpace(headingHeight + headingBottomMargin)

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
		d.sideBySideAddress(gap, tableX, vals.PickupAddress, leftAddressLabel, rightX, vals.DeliveryAddress, rightAddressLabel)
		d.addVerticalSpace(pxToMM(12.0))
	}
	err := d.twoColumnTable(tableX, d.pdf.GetY(), tableWidth, layout, vals)
	if err != nil {
		return fmt.Errorf("TwoColumnTable %w", err)
	}
	d.pdf.RoundedRect(pageSideMarginMm, startY, cardWidth, d.pdf.GetY()-startY, 1.0, "34", "D")
	shipmentCardBottomMargin := pxToMM(16.0)
	d.addVerticalSpace(shipmentCardBottomMargin)
	return nil
}

func (d *DynamicFormFiller) setHHGStripeColor(shipmentType models.MTOShipmentType) {
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
	d.pdf.SetDrawColor(r, g, b)
	d.pdf.SetFillColor(r, g, b)
}

func (d *DynamicFormFiller) sideBySideAddress(gap float64, leftAddressX float64, leftAddress string, leftAddressLabel string, rightAddressX float64, rightAddress string, rightAddressLabel string) {
	d.pdf.SetFontUnitSize(textFontSize)
	addressY := d.pdf.GetY()
	startY := d.pdf.GetY()
	d.pdf.SetX(leftAddressX)
	d.setTextColorBaseDark()
	addressWidth := rightAddressX - leftAddressX - gap
	if leftAddressLabel != "" {
		d.pdf.CellFormat(addressWidth, pxToMM(18.0), leftAddressLabel, "", 1, "LT", false, 0, "")
		addressY = d.pdf.GetY() + pxToMM(8.0)
	}
	if rightAddressLabel != "" {
		d.pdf.MoveTo(rightAddressX, startY)
		d.pdf.CellFormat(addressWidth, pxToMM(18.0), rightAddressLabel, "", 1, "LT", false, 0, "")
		addressY = math.Max(addressY, d.pdf.GetY()+pxToMM(8.0))
	}
	d.pdf.MoveTo(leftAddressX, addressY)
	d.setTextColorBaseDarkest()
	d.pdf.CellFormat(addressWidth, pxToMM(18.0), leftAddress, "", 1, "LT", false, 0, "")
	leftY := d.pdf.GetY()
	d.pdf.MoveTo(rightAddressX-pxToMM(20.0), addressY)
	d.drawArrow()

	d.pdf.MoveTo(rightAddressX, addressY)
	d.pdf.CellFormat(addressWidth, pxToMM(18.0), rightAddress, "", 1, "LT", false, 0, "")
	addressY = math.Max(leftY, d.pdf.GetY())
	d.pdf.SetY(addressY)
}

func (d *DynamicFormFiller) formatShipmentType(shipmentType models.MTOShipmentType) string {
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

func (d *DynamicFormFiller) twoColumnTable(x float64, y float64, w float64, layout []TableRow, data interface{}) error {
	gap := pxToMM(28.0)
	columnWidth := (w - gap) / 2.0
	labelWidth := 0.3 * columnWidth
	valueWidth := 0.7 * columnWidth
	d.pdf.SetY(y)
	d.pdf.SetFontUnitSize(textSmallFontSize)

	for i, row := range layout {
		err := d.twoColumnTableRow(x, gap, labelWidth, valueWidth, row, data)
		if err != nil {
			return err
		}
		if i < len(layout)-1 {
			d.setBorderColor()
			d.pdf.Line(x, d.pdf.GetY(), x+labelWidth+valueWidth, d.pdf.GetY())
			d.pdf.Line(x+labelWidth+valueWidth+gap, d.pdf.GetY(), x+gap+2.0*(labelWidth+valueWidth), d.pdf.GetY())
			d.addVerticalSpace(1.0)
		}
	}
	d.addVerticalSpace(2.0)
	return nil
}

func (d *DynamicFormFiller) twoColumnTableRow(x float64, gap float64, labelWidth float64, valueWidth float64, row TableRow, data interface{}) error {
	rowStartY := d.pdf.GetY()

	leftVal, err := getField(row.LeftFieldName, data)
	if err != nil {
		return err
	}
	d.tableColumn(x, labelWidth, valueWidth, row.LeftLabel, leftVal)
	leftValY := d.pdf.GetY()
	if row.RightFieldName == "" {
		return nil
	}
	d.pdf.SetY(rowStartY)
	rightVal, err := getField(row.RightFieldName, data)
	if err != nil {
		return err
	}
	d.tableColumn(x+labelWidth+valueWidth+gap, labelWidth, valueWidth, row.RightLabel, rightVal)
	rightValY := d.pdf.GetY()
	d.pdf.SetY(math.Max(leftValY, rightValY))
	return nil
}

// tableColumn draws one side of a two-column table row
func (d *DynamicFormFiller) tableColumn(x float64, labelWidth float64, valueWidth float64, label string, value string) {
	lineHeight := 5.0 // TODO this shadows a global
	d.pdf.SetX(x)
	d.pdf.SetFontStyle("B")
	d.setTextColorBaseDarker()
	d.pdf.CellFormat(labelWidth, 10.0, label, "", 0, "LT", false, 0, "")
	d.pdf.SetFontStyle("")
	d.setTextColorBaseDarkest()
	d.pdf.MultiCell(valueWidth, lineHeight, value, "", "LT", false)
}

func (d *DynamicFormFiller) drawArrow() {
	d.pdf.Image(arrowImageName, d.pdf.GetX(), d.pdf.GetY(), pxToMM(20.0), 0.0, flow, arrowImageFormat, imageLink, imageLinkURL)
}

// Loop through all pages and add headings. This must be done at the end because it uses the number of pages
func (d *DynamicFormFiller) addPageHeaders() {
	numPages := d.pdf.PageCount()
	for i := 1; i <= numPages; i++ {
		d.pdf.SetPage(i)
		d.reportPageHeader()
	}
}
