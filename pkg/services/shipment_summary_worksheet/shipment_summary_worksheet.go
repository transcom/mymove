package shipmentsummaryworksheet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// SSWPPMComputer is the concrete struct implementing the services.shipmentsummaryworksheet interface
type SSWPPMComputer struct {
	services.PPMCloseoutFetcher
}

// NewSSWPPMComputer creates a SSWPPMComputer
func NewSSWPPMComputer(ppmCloseoutFetcher services.PPMCloseoutFetcher) services.SSWPPMComputer {
	return &SSWPPMComputer{
		ppmCloseoutFetcher,
	}
}

// SSWPPMGenerator is the concrete struct implementing the services.shipmentsummaryworksheet interface
type SSWPPMGenerator struct {
	generator      paperwork.Generator
	templateReader *bytes.Reader
}

// NewSSWPPMGenerator creates a SSWPPMGenerator
func NewSSWPPMGenerator(pdfGenerator *paperwork.Generator) (services.SSWPPMGenerator, error) {
	templateReader, err := createAssetByteReader("paperwork/formtemplates/ShipmentSummaryWorksheet.pdf")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &SSWPPMGenerator{
		generator:      *pdfGenerator,
		templateReader: templateReader,
	}, nil
}

// FormatValuesShipmentSummaryWorksheet returns the formatted pages for the Shipment Summary Worksheet
func (SSWPPMComputer *SSWPPMComputer) FormatValuesShipmentSummaryWorksheet(shipmentSummaryFormData models.ShipmentSummaryFormData, isPaymentPacket bool) (services.Page1Values, services.Page2Values, services.Page3Values, error) {
	page1, err := SSWPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage1(shipmentSummaryFormData, isPaymentPacket)
	if err != nil {
		return page1, services.Page2Values{}, services.Page3Values{}, errors.WithStack(err)
	}
	page2, err := SSWPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage2(shipmentSummaryFormData, isPaymentPacket)
	if err != nil {
		return page1, page2, services.Page3Values{}, errors.WithStack(err)
	}
	page3, err := SSWPPMComputer.FormatValuesShipmentSummaryWorksheetFormPage3(shipmentSummaryFormData, isPaymentPacket)
	if err != nil {
		return page1, page2, services.Page3Values{}, errors.WithStack(err)
	}
	return page1, page2, page3, nil
}

// textField represents a text field within a form.
type textField struct {
	Pages     []int  `json:"pages"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	Multiline bool   `json:"multiline"`
	Locked    bool   `json:"locked"`
}

// WorkSheetSIT is an object representing SIT on the Shipment Summary Worksheet
type WorkSheetSIT struct {
	NumberAndTypes string
	EntryDates     string
	EndDates       string
	DaysInStorage  string
}

// Dollar represents a type for dollar monetary unit
type Dollar float64

// String is a string representation of a Dollar
func (d Dollar) String() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("$%.2f", d)
}

// Obligation an object representing the obligations section on the shipment summary worksheet
type Obligation struct {
	Gcc   unit.Cents
	SIT   unit.Cents
	Miles unit.Miles
}

// GCC100 calculates the 100% GCC on shipment summary worksheet
func (obligation Obligation) GCC100() float64 {
	return obligation.Gcc.ToDollarFloatNoRound()
}

// GCC95 calculates the 95% GCC on shipment summary worksheet
func (obligation Obligation) GCC95() float64 {
	return obligation.Gcc.MultiplyFloat64(.95).ToDollarFloatNoRound()
}

// FormatSIT formats the SIT Cost into a dollar float for the shipment summary worksheet
func (obligation Obligation) FormatSIT() float64 {
	return obligation.SIT.ToDollarFloatNoRound()
}

// MaxAdvance calculates the Max Advance on the shipment summary worksheet
func (obligation Obligation) MaxAdvance() float64 {
	return obligation.Gcc.MultiplyFloat64(.60).ToDollarFloatNoRound()
}

// SSWMaxWeightEntitlement weight allotment for the shipment summary worksheet.
type SSWMaxWeightEntitlement struct {
	Entitlement   unit.Pound
	ProGear       unit.Pound
	SpouseProGear unit.Pound
	TotalWeight   unit.Pound
}

type Certifications struct {
	CustomerField string
	OfficeField   string
	DateField     string
}

// adds a line item to shipment summary worksheet SSWMaxWeightEntitlement and increments total allotment
func (wa *SSWMaxWeightEntitlement) addLineItem(field string, value int) {
	r := reflect.ValueOf(wa).Elem()
	f := r.FieldByName(field)
	if f.IsValid() && f.CanSet() {
		f.SetInt(int64(value))
		wa.TotalWeight += unit.Pound(value)
	}
}

// SSWGetEntitlement calculates the entitlement for the shipment summary worksheet based on the parameters of
// a move (hasDependents, spouseHasProGear)
func SSWGetEntitlement(grade internalmessages.OrderPayGrade, hasDependents bool, spouseHasProGear bool) models.SSWMaxWeightEntitlement {
	sswEntitlements := SSWMaxWeightEntitlement{}
	entitlements := models.GetWeightAllotment(grade)
	sswEntitlements.addLineItem("ProGear", entitlements.ProGearWeight)
	sswEntitlements.addLineItem("SpouseProGear", entitlements.ProGearWeightSpouse)
	if !hasDependents {
		sswEntitlements.addLineItem("Entitlement", entitlements.TotalWeightSelf)
		return models.SSWMaxWeightEntitlement(sswEntitlements)
	}
	sswEntitlements.addLineItem("Entitlement", entitlements.TotalWeightSelfPlusDependents)
	return models.SSWMaxWeightEntitlement(sswEntitlements)
}

// Calculates cost for the Remaining PPM Incentive (pre-tax) field on page 2 of SSW form.
func CalculateRemainingPPMEntitlement(finalIncentive *unit.Cents, sitMemberPaid float64, sitGTCCPaid float64, aoa *unit.Cents) float64 {
	// FinalIncentive
	var finalIncentiveFloat float64 = 0
	if finalIncentive != nil {
		finalIncentiveFloat = float64(*finalIncentive) / 100.0
	}

	var aoaFloat float64 = 0
	if aoa != nil {
		aoaFloat = float64(*aoa) / 100.0
	}

	// This costing is computed by taking the Actual Obligations 100% GCC plus the
	// SIT cost calculated (if SIT was approved and accepted) minus any Advance
	// Operating Allowance (AOA) the customer identified as receiving in the Document upload process
	return (finalIncentiveFloat + sitMemberPaid + sitGTCCPaid) - aoaFloat
}

const (
	controlledUnclassifiedInformationText = "CONTROLLED UNCLASSIFIED INFORMATION"
)

const (
	trustedAgentText = "Trusted Agent Requires POA \nor Letter of Authorization"
)

// FormatValuesShipmentSummaryWorksheetFormPage1 formats the data for page 1 of the Shipment Summary Worksheet
func (s SSWPPMComputer) FormatValuesShipmentSummaryWorksheetFormPage1(data models.ShipmentSummaryFormData, isPaymentPacket bool) (services.Page1Values, error) {
	var err error
	page1 := services.Page1Values{}
	page1.CUIBanner = controlledUnclassifiedInformationText
	page1.MaxSITStorageEntitlement = fmt.Sprintf("%02d Days in SIT", data.MaxSITStorageEntitlement)
	// We don't currently know what allows POV to be authorized, so we are hardcoding it to "No" to start
	page1.POVAuthorized = "No"

	sm := data.ServiceMember
	page1.ServiceMemberName = FormatServiceMemberFullName(sm)
	page1.PreferredPhoneNumber = derefStringTypes(sm.Telephone)
	page1.ServiceBranch = FormatServiceMemberAffiliation(sm.Affiliation)
	page1.PreferredEmail = derefStringTypes(sm.PersonalEmail)
	page1.DODId = derefStringTypes(sm.Edipi)
	page1.MailingAddressW2 = FormatAddress(data.W2Address)
	page1.RankGrade = FormatGrade(data.Order.Grade)

	page1.IssuingBranchOrAgency = FormatServiceMemberAffiliation(sm.Affiliation)
	page1.OrdersIssueDate = FormatDate(data.Order.IssueDate)
	page1.OrdersTypeAndOrdersNumber = FormatOrdersTypeAndOrdersNumber(data.Order)

	page1.AuthorizedOrigin = data.CurrentDutyLocation.Name
	page1.AuthorizedDestination = data.NewDutyLocation.Name

	page1.NewDutyAssignment = data.NewDutyLocation.Name

	page1.WeightAllotment = FormatWeights(data.WeightAllotment.Entitlement)
	page1.WeightAllotmentProGear = FormatWeights(data.WeightAllotment.ProGear)
	page1.WeightAllotmentProgearSpouse = FormatWeights(data.WeightAllotment.SpouseProGear)
	page1.TotalWeightAllotment = FormatWeights(data.WeightAllotment.TotalWeight)

	formattedSIT := WorkSheetSIT{}

	formattedShipment := s.FormatShipment(data.PPMShipment, data.WeightAllotment, isPaymentPacket)
	page1.ShipmentNumberAndTypes = formattedShipment.ShipmentNumberAndTypes
	page1.ShipmentPickUpDates = formattedShipment.PickUpDates
	page1.ShipmentCurrentShipmentStatuses = formattedShipment.CurrentShipmentStatuses

	// Shipment weights for Payment Packet are actual, for AOA Packet are estimated.
	if isPaymentPacket {
		formattedSIT = FormatAllSITSForPaymentPacket(data.MovingExpenses)

		finalPPMWeight := FormatPPMWeightFinal(data.PPMShipmentFinalWeight)
		page1.ShipmentWeights = finalPPMWeight
		page1.ActualObligationGCC100 = finalPPMWeight + "; " + formattedShipment.FinalIncentive
		page1.PreparationDate1, err = formatSSWDate(data.SignedCertifications, data.PPMShipment.ID)
		if err != nil {
			return page1, err
		}
	} else {
		formattedSIT = FormatAllSITSForAOAPacket(data.PPMShipment)

		page1.ShipmentWeights = formattedShipment.ShipmentWeights
		page1.ActualObligationGCC100 = formattedShipment.ShipmentWeightForObligation + " - Actual lbs; "

		page1.PreparationDate1 = formatAOADate(data.SignedCertifications, data.PPMShipment.ID)
	}

	// Fill out form fields related to Actual Expense Reimbursement status
	if data.PPMShipment.IsActualExpenseReimbursement != nil {
		page1.IsActualExpenseReimbursement = *data.PPMShipment.IsActualExpenseReimbursement
	}

	page1.SITDaysInStorage = formattedSIT.DaysInStorage
	page1.SITEntryDates = formattedSIT.EntryDates
	page1.SITEndDates = formattedSIT.EndDates
	page1.SITNumberAndTypes = formattedShipment.ShipmentNumberAndTypes

	page1.MaxObligationGCC100 = FormatWeights(data.WeightAllotment.Entitlement) + " lbs; " + formattedShipment.EstimatedIncentive
	page1.MaxObligationGCCMaxAdvance = formattedShipment.MaxAdvance
	page1.ActualObligationAdvance = formattedShipment.AdvanceAmountReceived
	page1.MaxObligationSIT = fmt.Sprintf("%02d Days in SIT", data.MaxSITStorageEntitlement)
	page1.ActualObligationSIT = formattedSIT.DaysInStorage
	page1.TotalWeightAllotmentRepeat = page1.TotalWeightAllotment
	return page1, nil
}

// FormatValuesShipmentSummaryWorksheetFormPage2 formats the data for page 2 of the Shipment Summary Worksheet
func (s *SSWPPMComputer) FormatValuesShipmentSummaryWorksheetFormPage2(data models.ShipmentSummaryFormData, isPaymentPacket bool) (services.Page2Values, error) {
	var err error
	expensesMap := SubTotalExpenses(data.MovingExpenses)
	certificationInfo := formatSignedCertifications(data.SignedCertifications, data.PPMShipment.ID, isPaymentPacket)
	formattedShipments := s.FormatShipment(data.PPMShipment, data.WeightAllotment, isPaymentPacket)

	page2 := services.Page2Values{}
	page2.CUIBanner = controlledUnclassifiedInformationText
	page2.TAC = derefStringTypes(data.Order.TAC)
	page2.SAC = derefStringTypes(data.Order.SAC)
	if isPaymentPacket {
		data.PPMRemainingEntitlement = CalculateRemainingPPMEntitlement(data.PPMShipment.FinalIncentive, expensesMap["StorageMemberPaid"], expensesMap["StorageGTCCPaid"], data.PPMShipment.AdvanceAmountReceived)
		page2.PPMRemainingEntitlement = FormatDollars(data.PPMRemainingEntitlement)
		page2.PreparationDate2, err = formatSSWDate(data.SignedCertifications, data.PPMShipment.ID)
		if err != nil {
			return page2, err
		}
		page2.Disbursement = formatDisbursement(expensesMap, data.PPMRemainingEntitlement)
	} else {
		page2.PreparationDate2 = formatAOADate(data.SignedCertifications, data.PPMShipment.ID)
		page2.Disbursement = "N/A"
		page2.PPMRemainingEntitlement = "N/A"
	}
	page2.ContractedExpenseMemberPaid = FormatDollars(expensesMap["ContractedExpenseMemberPaid"])
	page2.ContractedExpenseGTCCPaid = FormatDollars(expensesMap["ContractedExpenseGTCCPaid"])
	page2.PackingMaterialsMemberPaid = FormatDollars(expensesMap["PackingMaterialsMemberPaid"])
	page2.PackingMaterialsGTCCPaid = FormatDollars(expensesMap["PackingMaterialsGTCCPaid"])
	page2.WeighingFeesMemberPaid = FormatDollars(expensesMap["WeighingFeeMemberPaid"])
	page2.WeighingFeesGTCCPaid = FormatDollars(expensesMap["WeighingFeeGTCCPaid"])
	page2.RentalEquipmentMemberPaid = FormatDollars(expensesMap["RentalEquipmentMemberPaid"])
	page2.RentalEquipmentGTCCPaid = FormatDollars(expensesMap["RentalEquipmentGTCCPaid"])
	page2.TollsMemberPaid = FormatDollars(expensesMap["TollsMemberPaid"])
	page2.TollsGTCCPaid = FormatDollars(expensesMap["TollsGTCCPaid"])
	page2.OilMemberPaid = FormatDollars(expensesMap["OilMemberPaid"])
	page2.OilGTCCPaid = FormatDollars(expensesMap["OilGTCCPaid"])
	page2.OtherMemberPaid = FormatDollars(expensesMap["OtherMemberPaid"])
	page2.OtherGTCCPaid = FormatDollars(expensesMap["OtherGTCCPaid"])
	page2.TotalMemberPaid = FormatDollars(expensesMap["TotalMemberPaid"])
	page2.TotalGTCCPaid = FormatDollars(expensesMap["TotalGTCCPaid"])
	page2.TotalMemberPaidRepeated = FormatDollars(expensesMap["TotalMemberPaid"])
	page2.TotalGTCCPaidRepeated = FormatDollars(expensesMap["TotalGTCCPaid"])
	page2.TotalMemberPaidSIT = FormatDollars(expensesMap["StorageMemberPaid"])
	page2.TotalGTCCPaidSIT = FormatDollars(expensesMap["StorageGTCCPaid"])
	page2.TotalMemberPaidRepeated = page2.TotalMemberPaid
	page2.TotalGTCCPaidRepeated = page2.TotalGTCCPaid
	page2.ShipmentPickupDates = formattedShipments.PickUpDates
	page2.TrustedAgentName = trustedAgentText
	page2.ServiceMemberSignature = certificationInfo.CustomerField
	page2.PPPOPPSORepresentative = certificationInfo.OfficeField
	page2.SignatureDate = certificationInfo.DateField

	return page2, nil
}

// FormatValuesShipmentSummaryWorksheetFormPage3 formats the data for page 3 of the Shipment Summary Worksheet
func (s *SSWPPMComputer) FormatValuesShipmentSummaryWorksheetFormPage3(data models.ShipmentSummaryFormData, isPaymentPacket bool) (services.Page3Values, error) {
	var err error
	var page3 services.Page3Values

	if isPaymentPacket {
		page3.PreparationDate3, err = formatSSWDate(data.SignedCertifications, data.PPMShipment.ID)
		if err != nil {
			return page3, err
		}
	} else {
		page3.PreparationDate3 = formatAOADate(data.SignedCertifications, data.PPMShipment.ID)
	}

	page3Map, err := formatAdditionalShipments(data)
	if err != nil {
		return page3, err
	}
	page3.AddShipments = page3Map
	return page3, nil
}

func formatAdditionalShipments(ssfd models.ShipmentSummaryFormData) (map[string]string, error) {
	page3Map := make(map[string]string)
	hasCurrentPPM := false
	const rows = 16
	for i, shipment := range ssfd.AllShipments {

		// If this is the shipment the SSW is being generated for, skip it.
		if shipment.PPMShipment.ID == ssfd.PPMShipment.ID {
			hasCurrentPPM = true
			continue
		}

		// This ensures that skipping the current PPM does not cause any row skips due to db fetch order
		if !hasCurrentPPM {
			i = i + 1
		}

		// If after skipping the current PPM, i is more than the amount of rows we have, throw an error.
		if i > rows {
			err := errors.New("PDF is being generated for a move with more than 17 shipments, SSW cannot display them all")
			return nil, err
		}

		// Default values will be configured for HHG, shipment-specific values configured below in switch case
		// This helps us to prevent redundant and confusing code for each shipment type
		page3Map, err := formatAdditionalHHG(page3Map, i, shipment)
		if err != nil {
			return nil, err
		}

		// Switch handles unique values by shipment type
		switch {
		case shipment.ShipmentType == models.MTOShipmentTypePPM:
			// Weights
			totalWeight := models.GetPPMNetWeight(*shipment.PPMShipment)
			if totalWeight != 0 {
				page3Map[fmt.Sprintf("AddShipmentWeights%d", i)] = FormatPPMWeightFinal(totalWeight) // Comment happens in formatter
			} else if shipment.PPMShipment.EstimatedWeight != nil {
				page3Map[fmt.Sprintf("AddShipmentWeights%d", i)] = FormatPPMWeightEstimated(*shipment.PPMShipment) // Comment happens in formatter
			} else {
				page3Map[fmt.Sprintf("AddShipmentWeights%d", i)] = " - "
			}
			// Dates
			if shipment.PPMShipment.ActualMoveDate != nil {
				page3Map[fmt.Sprintf("AddShipmentPickUpDates%d", i)] = FormatDate(*shipment.PPMShipment.ActualMoveDate) + " Actual"

			} else {
				page3Map[fmt.Sprintf("AddShipmentPickUpDates%d", i)] = FormatDate(shipment.PPMShipment.ExpectedDepartureDate) + " Expected"

			}
			// PPM Status instead of shipment status
			page3Map[fmt.Sprintf("AddShipmentStatus%d", i)] = FormatCurrentPPMStatus(*shipment.PPMShipment)
		case shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom:
			page3Map[fmt.Sprintf("AddShipmentNumberAndTypes%d", i)] = *shipment.ShipmentLocator + " NTS Release"
		case shipment.ShipmentType == models.MTOShipmentTypeHHGIntoNTSDom:
			page3Map[fmt.Sprintf("AddShipmentNumberAndTypes%d", i)] = *shipment.ShipmentLocator + " NTS"
		case shipment.ShipmentType == models.MTOShipmentTypeMobileHome:
			page3Map[fmt.Sprintf("AddShipmentNumberAndTypes%d", i)] = *shipment.ShipmentLocator + " Mobile Home"
		case shipment.ShipmentType == models.MTOShipmentTypeBoatHaulAway:
			page3Map[fmt.Sprintf("AddShipmentNumberAndTypes%d", i)] = *shipment.ShipmentLocator + " Boat Haul"
		case shipment.ShipmentType == models.MTOShipmentTypeBoatTowAway:
			page3Map[fmt.Sprintf("AddShipmentNumberAndTypes%d", i)] = *shipment.ShipmentLocator + " Boat Tow"
		case shipment.ShipmentType == models.MTOShipmentTypeUnaccompaniedBaggage:
			page3Map[fmt.Sprintf("AddShipmentNumberAndTypes%d", i)] = *shipment.ShipmentLocator + " UB"
		}

	}
	return page3Map, nil
}

func formatAdditionalHHG(page3Map map[string]string, i int, shipment models.MTOShipment) (map[string]string, error) {
	// If no shipment locator, throw error because something is wrong
	if shipment.ShipmentLocator != nil {
		page3Map[fmt.Sprintf("AddShipmentNumberAndTypes%d", i)] = *shipment.ShipmentLocator + " " + string(shipment.ShipmentType)
	} else {
		err := errors.New("PDF is being generated for a move without a locator")
		return nil, err
	}

	// If we're missing pickup dates or weights, we return " - " instead of error. Also it may be a PPM
	// For dates, we prefer actual -> scheduled -> requested -> -
	if shipment.ActualPickupDate != nil {
		page3Map[fmt.Sprintf("AddShipmentPickUpDates%d", i)] = FormatDate(*shipment.ActualPickupDate) + " Actual"
	} else if shipment.ScheduledPickupDate != nil {
		page3Map[fmt.Sprintf("AddShipmentPickUpDates%d", i)] = FormatDate(*shipment.ScheduledPickupDate) + " Scheduled"

	} else if shipment.RequestedPickupDate != nil {
		page3Map[fmt.Sprintf("AddShipmentPickUpDates%d", i)] = FormatDate(*shipment.RequestedPickupDate) + " Requested"

	} else {
		page3Map[fmt.Sprintf("AddShipmentPickUpDates%d", i)] = " - "
	}

	// For weights, we prefer actual -> estimated -> -
	if shipment.PrimeActualWeight != nil {
		page3Map[fmt.Sprintf("AddShipmentWeights%d", i)] = FormatWeights(*shipment.PrimeActualWeight) + " Actual"
	} else if shipment.PrimeEstimatedWeight != nil {
		page3Map[fmt.Sprintf("AddShipmentWeights%d", i)] = FormatWeights(*shipment.PrimeEstimatedWeight) + " Estimated"
	} else {
		page3Map[fmt.Sprintf("AddShipmentWeights%d", i)] = " - "
	}

	// Status is always available
	page3Map[fmt.Sprintf("AddShipmentStatus%d", i)] = FormatEnum(string(shipment.Status), "")

	return page3Map, nil
}

// FormatGrade formats the service member's rank for Shipment Summary Worksheet
func FormatGrade(grade *internalmessages.OrderPayGrade) string {
	var gradeDisplayValue = map[internalmessages.OrderPayGrade]string{
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
		models.ServiceMemberGradeO1ACADEMYGRADUATE:       "O-1 or Service Academy Graduate",
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
	if grade != nil {
		return gradeDisplayValue[*grade]
	}
	return ""
}

func formatEmplid(serviceMember models.ServiceMember) (*string, error) {
	const prefix = "EMPLID:"
	const separator = " "
	if *serviceMember.Affiliation == models.AffiliationCOASTGUARD && serviceMember.Emplid != nil {
		slice := []string{prefix, *serviceMember.Emplid}
		formattedReturn := strings.Join(slice, separator)
		return &formattedReturn, nil
	} else {
		return serviceMember.Edipi, nil
	}
}

func formatMaxAdvance(estimatedIncentive *unit.Cents) string {
	if estimatedIncentive != nil {
		maxAdvance := float64(*estimatedIncentive) * 0.6
		return FormatDollars(maxAdvance / 100)
	}
	maxAdvanceString := "No Incentive Found"
	return maxAdvanceString

}

func formatSignedCertifications(signedCertifications []*models.SignedCertification, ppmid uuid.UUID, isPaymentPacket bool) Certifications {
	certifications := Certifications{}
	// Strings used to build return values
	var customerSignature string
	var aoaSignature string
	var sswSignature string
	var aoaDate string
	var sswDate string

	// This loop evaluates all certs, move-level customer signature doesn't have a ppm id, it's collected first, then office signatures with ppmids
	for _, cert := range signedCertifications {
		if cert.PpmID == nil { // Original move signature required, doesn't have ppmid. All others of that type do
			if *cert.CertificationType == models.SignedCertificationTypeSHIPMENT {
				customerSignature = cert.Signature
			}
		} else if *cert.PpmID == ppmid { // PPM ID needs to be checked to prevent signatures from other PPMs on the same move from populating
			switch {
			case *cert.CertificationType == models.SignedCertificationTypePreCloseoutReviewedPPMPAYMENT:
				aoaSignature = cert.Signature
				aoaDate = FormatDate(cert.UpdatedAt) // We use updatedat to get the most recent signature dates
			case *cert.CertificationType == models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT:
				sswSignature = cert.Signature
				sswDate = FormatDate(cert.UpdatedAt) // We use updatedat to get the most recent signature dates
			}
		}
	}
	certifications.CustomerField = customerSignature
	certifications.OfficeField = "AOA: " + aoaSignature
	certifications.DateField = "AOA: " + aoaDate

	if isPaymentPacket {
		certifications.OfficeField += "\nSSW: " + sswSignature
		certifications.DateField += "\nSSW: " + sswDate
	}

	return certifications
}

// The following formats the preparation date, as the preparation date for AOAs is the date the service counselor certifies the advance.
func formatAOADate(signedCertifications []*models.SignedCertification, ppmid uuid.UUID) string {
	// This loop evaluates certs to find Office AOA Signature date
	for _, cert := range signedCertifications {
		if cert.PpmID != nil { // Required to avoid error, service members signatures have nil ppm ids
			if *cert.PpmID == ppmid { // PPM ID needs to be checked to prevent signatures from other PPMs on the same move from populating
				if *cert.CertificationType == models.SignedCertificationTypePreCloseoutReviewedPPMPAYMENT {
					aoaDate := FormatDate(cert.UpdatedAt) // We use updatedat to get the most recent signature dates
					return aoaDate
				}
			}
		}
	}
	return FormatDate(time.Now())
}

// The following formats the preparation date, as the preparation date for SSWs is the date the closeout counselor certifies the closeout.
func formatSSWDate(signedCertifications []*models.SignedCertification, ppmid uuid.UUID) (string, error) {
	// This loop evaluates certs to find Office SSW Signature date
	for _, cert := range signedCertifications {
		if cert.PpmID != nil { // Required to avoid error, service members signatures have nil ppm ids
			if *cert.PpmID == ppmid { // PPM ID needs to be checked to prevent signatures from other PPMs on the same move from populating
				if *cert.CertificationType == models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT {
					sswDate := FormatDate(cert.UpdatedAt) // We use updatedat to get the most recent signature dates
					return sswDate, nil
				}
			}
		}
	}
	return "", errors.New("Payment Packet is not certified")
}

// FormatAddress retrieves a PPMShipment W2Address and formats it for the SSW Document
func FormatAddress(w2Address *models.Address) string {
	var addressString string

	var country string
	if w2Address != nil && w2Address.Country != nil && w2Address.Country.Country != "" {
		country = w2Address.Country.Country
	} else {
		country = ""
	}

	if w2Address != nil {
		addressString = fmt.Sprintf("%s, %s %s%s %s %s%s",
			w2Address.StreetAddress1,
			nilOrValue(w2Address.StreetAddress2),
			nilOrValue(w2Address.StreetAddress3),
			w2Address.City,
			w2Address.State,
			nilOrValue(&country),
			w2Address.PostalCode,
		)
	} else {
		return "W2 Address not found"
	}

	return addressString
}

// nilOrValue returns the dereferenced value if the pointer is not nil, otherwise an empty string.
func nilOrValue(str *string) string {
	if str != nil {
		return *str
	}
	return ""
}

// FormatServiceMemberFullName formats ServiceMember full name for Shipment Summary Worksheet
func FormatServiceMemberFullName(serviceMember models.ServiceMember) string {
	lastName := derefStringTypes(serviceMember.LastName)
	suffix := derefStringTypes(serviceMember.Suffix)
	firstName := derefStringTypes(serviceMember.FirstName)
	middleName := derefStringTypes(serviceMember.MiddleName)
	if suffix != "" {
		return fmt.Sprintf("%s %s, %s %s", lastName, suffix, firstName, middleName)
	}
	return strings.TrimSpace(fmt.Sprintf("%s, %s %s", lastName, firstName, middleName))
}

func (s SSWPPMComputer) FormatShipment(ppm models.PPMShipment, weightAllotment models.SSWMaxWeightEntitlement, isPaymentPacket bool) models.WorkSheetShipment {
	formattedShipment := models.WorkSheetShipment{}

	if ppm.FinalIncentive != nil {
		formattedShipment.FinalIncentive = FormatDollarFromCents(*ppm.FinalIncentive)
	} else {
		formattedShipment.FinalIncentive = "No final incentive."
	}
	if ppm.EstimatedIncentive != nil {
		formattedShipment.MaxAdvance = formatMaxAdvance(ppm.EstimatedIncentive)
		formattedShipment.EstimatedIncentive = FormatDollarFromCents(*ppm.EstimatedIncentive)
	} else {
		formattedShipment.MaxAdvance = "Advance not available."
		formattedShipment.EstimatedIncentive = "No estimated incentive."
	}
	formattedShipmentTotalWeights := unit.Pound(0)
	formattedNumberAndTypes := *ppm.Shipment.ShipmentLocator + " PPM"
	formattedShipmentWeights := FormatPPMWeightEstimated(ppm)
	formattedShipmentStatuses := FormatCurrentPPMStatus(ppm)
	if ppm.EstimatedWeight != nil {
		formattedShipmentTotalWeights += s.calculateShipmentTotalWeight(ppm, weightAllotment)
	}
	formattedPickUpDates := FormatDate(ppm.ExpectedDepartureDate)
	// If advance isn't configured or received, it's false
	var hasRequestedAdvance bool
	if ppm.HasRequestedAdvance == nil {
		hasRequestedAdvance = false
	} else {
		hasRequestedAdvance = *ppm.HasRequestedAdvance
	}
	if isPaymentPacket {
		formattedPickUpDates = "N/A"
		if ppm.ActualMoveDate != nil {
			formattedPickUpDates = FormatDate(*ppm.ActualMoveDate)
		}
		// If it's received, reflect that
		if ppm.AdvanceAmountReceived != nil {
			formattedShipment.AdvanceAmountReceived = FormatDollarFromCents(*ppm.AdvanceAmountReceived) + "Customer received"
		} else if hasRequestedAdvance {
			// If it's requested, give amount and status
			if *ppm.AdvanceStatus != models.PPMAdvanceStatusReceived {
				formattedShipment.AdvanceAmountReceived = FormatDollarFromCents(*ppm.AdvanceAmountRequested) + " Requested, " + FormatEnum(string(*ppm.AdvanceStatus), "")
			} else {
				// If it's received, give received amount and status
				formattedShipment.AdvanceAmountReceived = FormatDollarFromCents(*ppm.AdvanceAmountReceived) + " Requested, " + FormatEnum(string(*ppm.AdvanceStatus), "")
			}
		} else {
			formattedShipment.AdvanceAmountReceived = "No Advance Requested."
		}
	} else {
		// No customer received amount in AOA packet
		if hasRequestedAdvance {
			if ppm.AdvanceStatus != nil {
				if *ppm.AdvanceStatus != models.PPMAdvanceStatusReceived {
					formattedShipment.AdvanceAmountReceived = FormatDollarFromCents(*ppm.AdvanceAmountRequested) + " Requested, " + FormatEnum(string(*ppm.AdvanceStatus), "")
				} else {
					// If it's received, give received amount and status
					formattedShipment.AdvanceAmountReceived = FormatDollarFromCents(*ppm.AdvanceAmountReceived) + " Requested, " + FormatEnum(string(*ppm.AdvanceStatus), "")
				}
			}
			// If it's requested, give amount and status
		} else {
			formattedShipment.AdvanceAmountReceived = "No Advance Requested."
		}
	}

	// Last resort in case any dates are stored incorrectly
	if formattedPickUpDates == "01-Jan-0001" {
		formattedPickUpDates = "N/A"
	}

	formattedShipment.ShipmentNumberAndTypes = formattedNumberAndTypes
	formattedShipment.PickUpDates = formattedPickUpDates
	formattedShipment.ShipmentWeights = formattedShipmentWeights
	formattedShipment.ShipmentWeightForObligation = FormatWeights(formattedShipmentTotalWeights)
	formattedShipment.CurrentShipmentStatuses = formattedShipmentStatuses

	return formattedShipment
}

// FormatAllSITs formats SIT line items for the Shipment Summary Worksheet Payment Packet
func FormatAllSITSForPaymentPacket(expenseDocuments models.MovingExpenses) WorkSheetSIT {
	formattedSIT := WorkSheetSIT{}

	for _, expense := range expenseDocuments {
		if *expense.MovingExpenseType == models.MovingExpenseReceiptTypeStorage {
			formattedSIT.EntryDates = FormatSITDate(expense.SITStartDate)
			formattedSIT.EndDates = FormatSITDate(expense.SubmittedSITEndDate)
			formattedSIT.DaysInStorage = FormatSITDaysInStorage(expense.SITStartDate, expense.SubmittedSITEndDate)
			return formattedSIT
		}
	}

	return formattedSIT
}

// FormatAllSITs formats SIT line items for the Shipment Summary Worksheet AOA Packet
func FormatAllSITSForAOAPacket(ppm models.PPMShipment) WorkSheetSIT {
	formattedSIT := WorkSheetSIT{}

	if ppm.SITEstimatedEntryDate != nil && ppm.SITEstimatedDepartureDate != nil {
		formattedSIT.EntryDates = FormatSITDate(ppm.SITEstimatedEntryDate)
		formattedSIT.EndDates = FormatSITDate(ppm.SITEstimatedDepartureDate)
		formattedSIT.DaysInStorage = FormatSITDaysInStorage(ppm.SITEstimatedEntryDate, ppm.SITEstimatedDepartureDate)
	}

	return formattedSIT
}

func (s SSWPPMComputer) calculateShipmentTotalWeight(ppmShipment models.PPMShipment, weightAllotment models.SSWMaxWeightEntitlement) unit.Pound {

	var err error
	var ppmActualWeight unit.Pound
	var maxLimit unit.Pound

	// Set maxLimit equal to the maximum weight entitlement or the allowable weight, whichever is lower
	if weightAllotment.TotalWeight < weightAllotment.Entitlement {
		maxLimit = weightAllotment.TotalWeight
	} else {
		maxLimit = weightAllotment.Entitlement
	}

	// Get the actual weight of the ppmShipment
	if len(ppmShipment.WeightTickets) > 0 {
		ppmActualWeight, err = s.PPMCloseoutFetcher.GetActualWeight(&ppmShipment)
		if err != nil {
			return 0
		}
	}

	// If actual weight is less than the lessor of maximum weight entitlement or the allowable weight, then use ppmActualWeight
	if ppmActualWeight < maxLimit {
		return ppmActualWeight
	} else {
		return maxLimit
	}
}

// FetchMovingExpensesShipmentSummaryWorksheet fetches moving expenses for the Shipment Summary Worksheet
// TODO: update to create moving expense summary with the new moving expense model
func FetchMovingExpensesShipmentSummaryWorksheet(PPMShipment models.PPMShipment, _ appcontext.AppContext, _ *auth.Session) (models.MovingExpenses, error) {
	var movingExpenseDocuments = PPMShipment.MovingExpenses

	return movingExpenseDocuments, nil
}

func SubTotalExpenses(expenseDocuments models.MovingExpenses) map[string]float64 {
	totals := make(map[string]float64)

	for _, expense := range expenseDocuments {
		if expense.MovingExpenseType == nil || expense.Amount == nil {
			continue
		} // Added quick nil check to ensure SSW returns while moving expenses are being added still
		var nilPPMDocumentStatus *models.PPMDocumentStatus
		if expense.Status != nilPPMDocumentStatus && (*expense.Status == models.PPMDocumentStatusRejected || *expense.Status == models.PPMDocumentStatusExcluded) {
			continue
		}
		expenseType, addToTotal := getExpenseType(expense)
		expenseDollarAmt := expense.Amount.ToDollarFloatNoRound()

		if expenseType == "StorageMemberPaid" {
			expenseDollarAmt = expense.SITReimburseableAmount.ToDollarFloatNoRound()
		}

		totals[expenseType] += expenseDollarAmt

		if addToTotal && expenseType != "Storage" {
			if paidWithGTCC := expense.PaidWithGTCC; paidWithGTCC != nil && *paidWithGTCC {
				totals["TotalGTCCPaid"] += expenseDollarAmt
			} else {
				totals["TotalMemberPaid"] += expenseDollarAmt
			}
		}
	}

	return totals
}

func getExpenseType(expense models.MovingExpense) (string, bool) {
	expenseType := FormatEnum(string(*expense.MovingExpenseType), "")
	addToTotal := expenseType != "Storage"

	if paidWithGTCC := expense.PaidWithGTCC; paidWithGTCC != nil && *paidWithGTCC {
		return fmt.Sprintf("%s%s", expenseType, "GTCCPaid"), addToTotal
	}

	return fmt.Sprintf("%s%s", expenseType, "MemberPaid"), addToTotal
}

// FormatCurrentPPMStatus formats FormatCurrentPPMStatus for the Shipment Summary Worksheet
func FormatCurrentPPMStatus(ppm models.PPMShipment) string {
	if ppm.Status == "PAYMENT_REQUESTED" {
		return "At destination"
	}
	return FormatEnum(string(ppm.Status), " ")
}

// FormatSITNumberAndType formats FormatSITNumberAndType for the Shipment Summary Worksheet
func FormatSITNumberAndType(i int) string {
	return fmt.Sprintf("%02d - SIT", i+1)
}

// FormatPPMWeight formats a ppms EstimatedNetWeight for the Shipment Summary Worksheet
func FormatPPMWeightEstimated(ppm models.PPMShipment) string {
	if ppm.EstimatedWeight != nil {
		wtg := FormatWeights(unit.Pound(*ppm.EstimatedWeight))
		return fmt.Sprintf("%s lbs - Estimated", wtg)
	}
	return ""
}

// FormatPPMWeight formats a ppms final NetWeight for the Shipment Summary Worksheet
func FormatPPMWeightFinal(weight unit.Pound) string {
	wtg := FormatWeights(unit.Pound(weight))
	return fmt.Sprintf("%s lbs - Actual", wtg)
}

// FormatSITDate formats a SIT Date for the Shipment Summary Worksheet
func FormatSITDate(sitDate *time.Time) string {
	if sitDate == nil {
		return "No SIT date" // Return string if no date found
	}
	return FormatDate(*sitDate)
}

// FormatSITDaysInStorage formats a SIT DaysInStorage for the Shipment Summary Worksheet
func FormatSITDaysInStorage(entryDate *time.Time, departureDate *time.Time) string {
	if entryDate == nil || departureDate == nil {
		return "No Entry/Departure Data" // Return string if no SIT attached
	}
	firstDate := *departureDate
	secondDate := *entryDate
	difference := firstDate.Sub(secondDate)
	formattedDifference := fmt.Sprintf("Days: %d\n", int64(difference.Hours()/24)+1)
	return formattedDifference
}

func formatDisbursement(expensesMap map[string]float64, ppmRemainingEntitlement float64) string {
	disbursementGTCC := expensesMap["TotalGTCCPaid"] + expensesMap["StorageGTCCPaid"]
	disbursementGTCCB := ppmRemainingEntitlement + expensesMap["StorageMemberPaid"]
	var disbursementMember float64
	// Disbursement GTCC is the lowest value of the above 2 calculations
	if disbursementGTCCB < disbursementGTCC {
		disbursementGTCC = disbursementGTCCB
	}
	if disbursementGTCC < 0 {
		// The only way this can happen is if the member overdrafted on their advance, resulting in negative GTCCB. In this case, the
		// disbursement member value will be liable for the negative difference, meaning they owe this money to the govt.
		disbursementMember = disbursementGTCC
		disbursementGTCC = 0
	} else {
		// Disbursement Member is remaining entitlement plus member SIT minus GTCC Disbursement, not less than 0.
		disbursementMember = ppmRemainingEntitlement + expensesMap["StorageMemberPaid"] - disbursementGTCC
	}

	// Return formatted values in string
	disbursementString := "GTCC: " + FormatDollars(disbursementGTCC) + "\nMember: " + FormatDollars(disbursementMember)
	return disbursementString
}

// FormatOrdersTypeAndOrdersNumber formats OrdersTypeAndOrdersNumber for Shipment Summary Worksheet
func FormatOrdersTypeAndOrdersNumber(order models.Order) string {
	orderType := FormatOrdersType(order)
	ordersNumber := derefStringTypes(order.OrdersNumber)
	return fmt.Sprintf("%s/%s", orderType, ordersNumber)
}

// FormatServiceMemberAffiliation formats ServiceMemberAffiliation in human friendly format
func FormatServiceMemberAffiliation(affiliation *models.ServiceMemberAffiliation) string {
	if affiliation != nil {
		return FormatEnum(string(*affiliation), " ")
	}
	return ""
}

// FormatOrdersType formats OrdersType for Shipment Summary Worksheet
func FormatOrdersType(order models.Order) string {
	switch order.OrdersType {
	case internalmessages.OrdersTypePERMANENTCHANGEOFSTATION:
		return "PCS"
	default:
		return ""
	}
}

// FormatDate formats Dates for Shipment Summary Worksheet
func FormatDate(date time.Time) string {
	dateLayout := "02-Jan-2006"
	return date.Format(dateLayout)
}

// FormatEnum titlecases string const types (e.g. THIS_CONSTANT -> This Constant)
// outSep specifies the character to use for rejoining the string
func FormatEnum(s string, outSep string) string {
	words := strings.Replace(strings.ToLower(s), "_", " ", -1)
	return strings.Replace(cases.Title(language.English).String(words), " ", outSep, -1)
}

// FormatWeights formats a unit.Pound using 000s separator
func FormatWeights(wtg unit.Pound) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", wtg)
}

// FormatDollars formats an int using 000s separator
func FormatDollars(dollars float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("$%.2f", dollars)
}

// FormatDollars formats cents using 000s separator
func FormatDollarFromCents(cents unit.Cents) string {
	d := float64(cents) / 100.0
	p := message.NewPrinter(language.English)
	return p.Sprintf("$%.2f", d)
}

func derefStringTypes(st interface{}) string {
	switch v := st.(type) {
	case *string:
		if v != nil {
			return *v
		}
	case string:
		return v
	}
	return ""
}

// ObligationType type corresponding to obligation sections of shipment summary worksheet
type ObligationType int

// ComputeObligations is helper function for computing the obligations section of the shipment summary worksheet
// Obligations must remain as static test data until new computer system is finished
func (SSWPPMComputer *SSWPPMComputer) ComputeObligations(_ appcontext.AppContext, _ models.ShipmentSummaryFormData, _ route.Planner) (obligation models.Obligations, err error) {
	// Obligations must remain test data until new computer system is finished
	obligations := models.Obligations{
		ActualObligation:           models.Obligation{Gcc: 123, SIT: 123, Miles: unit.Miles(123456)},
		MaxObligation:              models.Obligation{Gcc: 456, SIT: 456, Miles: unit.Miles(123456)},
		NonWinningActualObligation: models.Obligation{Gcc: 789, SIT: 789, Miles: unit.Miles(12345)},
		NonWinningMaxObligation:    models.Obligation{Gcc: 1000, SIT: 1000, Miles: unit.Miles(12345)},
	}
	return obligations, nil
}

// FetchDataShipmentSummaryWorksheetFormData fetches the pages for the Shipment Summary Worksheet for a given Move ID
func (SSWPPMComputer *SSWPPMComputer) FetchDataShipmentSummaryWorksheetFormData(appCtx appcontext.AppContext, session *auth.Session, ppmShipmentID uuid.UUID) (*models.ShipmentSummaryFormData, error) {

	ppmShipment := models.PPMShipment{}
	dbQErr := appCtx.DB().Q().Eager(
		"Shipment.MoveTaskOrder.Orders.ServiceMember",
		"Shipment.MoveTaskOrder.Orders.NewDutyLocation.Address",
		"Shipment.MoveTaskOrder.Orders.OriginDutyLocation.Address",
		"Shipment.MoveTaskOrder.MTOShipments.PPMShipment",
		"Shipment.MoveTaskOrder.MTOShipments.BoatShipment",
		"W2Address",
		"WeightTickets",
		"MovingExpenses",
	).Find(&ppmShipment, ppmShipmentID)

	if dbQErr != nil {
		if errors.Cause(dbQErr).Error() == models.RecordNotFoundErrorString {
			return nil, models.ErrFetchNotFound
		}
		return nil, dbQErr
	}

	// Final actual weight is a calculated value we don't store. This needs to be fetched independently
	// Requires WeightTickets eager preload
	ppmShipmentFinalWeight := models.GetPPMNetWeight(ppmShipment)

	serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
	if ppmShipment.Shipment.MoveTaskOrder.Orders.Grade == nil {
		return nil, errors.New("order for requested shipment summary worksheet data does not have a pay grade attached")
	}

	weightAllotment := SSWGetEntitlement(*ppmShipment.Shipment.MoveTaskOrder.Orders.Grade, ppmShipment.Shipment.MoveTaskOrder.Orders.HasDependents, ppmShipment.Shipment.MoveTaskOrder.Orders.SpouseHasProGear)

	maxSit, err := CalculateShipmentSITAllowance(appCtx, ppmShipment.Shipment)
	if err != nil {
		return nil, err
	}

	serviceMember.Edipi, err = formatEmplid(serviceMember)
	if err != nil {
		return nil, err
	}

	// Fetches all signed certifications for a move to be filtered in this file by ppmid and type
	signedCertifications, err := models.FetchSignedCertifications(appCtx.DB(), session, ppmShipment.Shipment.MoveTaskOrderID)
	if err != nil {
		return nil, err
	}

	var ppmShipments []models.PPMShipment

	ppmShipments = append(ppmShipments, ppmShipment)
	if ppmShipment.Shipment.MoveTaskOrder.Orders.OriginDutyLocation == nil {
		return nil, errors.New("order for PPM shipment does not have a origin duty location attached")
	}

	isActualExpenseReimbursement := false
	if ppmShipment.IsActualExpenseReimbursement != nil {
		isActualExpenseReimbursement = true
	}

	ssd := models.ShipmentSummaryFormData{
		AllShipments:                 ppmShipment.Shipment.MoveTaskOrder.MTOShipments,
		ServiceMember:                serviceMember,
		Order:                        ppmShipment.Shipment.MoveTaskOrder.Orders,
		Move:                         ppmShipment.Shipment.MoveTaskOrder,
		CurrentDutyLocation:          *ppmShipment.Shipment.MoveTaskOrder.Orders.OriginDutyLocation,
		NewDutyLocation:              ppmShipment.Shipment.MoveTaskOrder.Orders.NewDutyLocation,
		WeightAllotment:              weightAllotment,
		PPMShipment:                  ppmShipment,
		PPMShipments:                 ppmShipments,
		PPMShipmentFinalWeight:       ppmShipmentFinalWeight,
		W2Address:                    ppmShipment.W2Address,
		MovingExpenses:               ppmShipment.MovingExpenses,
		SignedCertifications:         signedCertifications,
		MaxSITStorageEntitlement:     maxSit,
		IsActualExpenseReimbursement: isActualExpenseReimbursement,
	}
	return &ssd, nil
}

// CalculateShipmentSITAllowance finds the number of days allowed in SIT for a shipment based on its entitlement and any approved SIT extensions
func CalculateShipmentSITAllowance(appCtx appcontext.AppContext, shipment models.MTOShipment) (int, error) {
	entitlement, err := fetchEntitlement(appCtx, shipment)
	if err != nil {
		return 0, err
	}

	totalSITAllowance := 0
	if entitlement.StorageInTransit != nil {
		totalSITAllowance = *entitlement.StorageInTransit
	}
	for _, ext := range shipment.SITDurationUpdates {
		if ext.ApprovedDays != nil {
			totalSITAllowance += *ext.ApprovedDays
		}
	}
	return totalSITAllowance, nil
}

func fetchEntitlement(appCtx appcontext.AppContext, mtoShipment models.MTOShipment) (*models.Entitlement, error) {
	var move models.Move
	err := appCtx.DB().Q().EagerPreload("Orders.Entitlement").Find(&move, mtoShipment.MoveTaskOrderID)

	if err != nil {
		return nil, err
	}

	return move.Orders.Entitlement, nil
}

// FillSSWPDFForm takes form data and fills an existing PDF form template with said data
func (SSWPPMGenerator *SSWPPMGenerator) FillSSWPDFForm(Page1Values services.Page1Values, Page2Values services.Page2Values, Page3Values services.Page3Values) (sswfile afero.File, pdfInfo *pdfcpu.PDFInfo, err error) {

	// header represents the header section of the JSON.
	type header struct {
		Source   string `json:"source"`
		Version  string `json:"version"`
		Creation string `json:"creation"`
		Producer string `json:"producer"`
	}

	// checkbox represents a checkbox within a form.
	type checkbox struct {
		Pages   []int  `json:"pages"`
		ID      string `json:"id"`
		Name    string `json:"name"`
		Default bool   `json:"value"`
		Value   bool   `json:"multiline"`
		Locked  bool   `json:"locked"`
	}

	// forms represents a form containing text fields.
	type form struct {
		TextField []textField `json:"textfield"`
		Checkbox  []checkbox  `json:"checkbox"`
	}

	// pdFData represents the entire JSON structure.
	type pdFData struct {
		Header header `json:"header"`
		Forms  []form `json:"forms"`
	}

	var sswHeader = header{
		Source:   "ShipmentSummaryWorksheet.pdf",
		Version:  "pdfcpu v0.9.1 dev",
		Creation: "2024-10-29 20:00:40 UTC",
		Producer: "macOS Version 13.5 (Build 22G74) Quartz PDFContext, AppendMode 1.1",
	}

	isActualExpenseReimbursement := false
	if Page1Values.IsActualExpenseReimbursement {
		isActualExpenseReimbursement = true
		Page1Values.GCCIsActualExpenseReimbursement = "Actual Expense Reimbursement"
		Page2Values.IncentiveIsActualExpenseReimbursement = "Actual Expense Reimbursement"
		Page2Values.HeaderIsActualExpenseReimbursement = `This PPM is being processed at actual expense reimbursement for valid expenses not to exceed the
		government constructed cost (GCC).`
	}

	var sswCheckbox = []checkbox{
		{
			Pages:   []int{2},
			ID:      "198",
			Name:    "EDOther",
			Value:   true,
			Default: false,
			Locked:  false,
		},
		{
			Pages:   []int{1},
			ID:      "444",
			Name:    "IsActualExpenseReimbursement",
			Value:   true,
			Default: isActualExpenseReimbursement,
			Locked:  false,
		},
	}

	formData := pdFData{ // This is unique to each PDF template, must be found for new templates using PDFCPU's export function used on the template (can be done through CLI)
		Header: sswHeader,
		Forms: []form{
			{ // Dynamically loops, creates, and aggregates json for text fields, merges page 1 and 2
				TextField: mergeTextFields(createTextFields(Page1Values, 1), createTextFields(Page2Values, 2), createTextFields(Page3Values, 3)),
				Checkbox:  sswCheckbox,
			},
		},
	}

	// Marshal the FormData struct into a JSON-encoded byte slice
	jsonData, err := json.MarshalIndent(formData, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	SSWWorksheet, err := SSWPPMGenerator.generator.FillPDFForm(jsonData, SSWPPMGenerator.templateReader, "")
	if err != nil {
		return nil, nil, err
	}

	// pdfInfo.PageCount is a great way to tell whether returned PDF is corrupted. Pages is expected pages
	const pages = 3
	pdfInfoResult, err := SSWPPMGenerator.generator.GetPdfFileInfo(SSWWorksheet.Name())
	if err != nil || pdfInfoResult.PageCount != pages {
		return nil, nil, errors.Wrap(err, "SSWGenerator output a corrupted or incorretly altered PDF")
	}
	// Return PDFInfo for additional testing in other functions
	pdfInfo = pdfInfoResult
	return SSWWorksheet, pdfInfo, err
}

// CreateTextFields formats the SSW Page data to match PDF-accepted JSON
func createTextFields(data interface{}, pages ...int) []textField {
	var textFields []textField
	val := reflect.ValueOf(data)

	// Process top-level struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// Handle map for additional shipments on page 3
		if fieldValue.Kind() == reflect.Map {
			for _, key := range fieldValue.MapKeys() {
				mapValue := fieldValue.MapIndex(key)

				textFieldEntry := textField{
					Pages:     pages,
					ID:        fmt.Sprintf("%d", len(textFields)+1),
					Name:      fmt.Sprintf("%v", key),
					Value:     fmt.Sprintf("%v", mapValue.Interface()),
					Multiline: true,
					Locked:    false,
				}
				textFields = append(textFields, textFieldEntry)
			}
		} else {
			// handle primitive fields
			textFieldEntry := textField{
				Pages:     pages,
				ID:        fmt.Sprintf("%d", len(textFields)+1),
				Name:      field.Name,
				Value:     fmt.Sprintf("%v", fieldValue.Interface()),
				Multiline: true,
				Locked:    false,
			}
			textFields = append(textFields, textFieldEntry)
		}
	}
	return textFields
}

// MergeTextFields merges page 1, page 2, and page 3 data
func mergeTextFields(fields1, fields2, fields3 []textField) []textField {
	totalFields := append(fields1, fields2...)
	return append(totalFields, fields3...)
}

// createAssetByteReader creates a new byte reader based on the TemplateImagePath of the formLayout
func createAssetByteReader(path string) (*bytes.Reader, error) {
	asset, err := assets.Asset(path)
	if err != nil {
		return nil, errors.Wrap(err, "error creating asset from path; check image path : "+path)
	}

	return bytes.NewReader(asset), nil
}
