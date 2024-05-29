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
}

// NewSSWPPMComputer creates a SSWPPMComputer
func NewSSWPPMComputer() services.SSWPPMComputer {
	return &SSWPPMComputer{}
}

// SSWPPMGenerator is the concrete struct implementing the services.shipmentsummaryworksheet interface
type SSWPPMGenerator struct {
	generator      paperwork.Generator
	templateReader *bytes.Reader
}

// NewSSWPPMGenerator creates a SSWPPMGenerator
func NewSSWPPMGenerator(pdfGenerator *paperwork.Generator) (services.SSWPPMGenerator, error) {
	templateReader, err := createAssetByteReader("paperwork/formtemplates/SSWPDFTemplate.pdf")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &SSWPPMGenerator{
		generator:      *pdfGenerator,
		templateReader: templateReader,
	}, nil
}

// FormatValuesShipmentSummaryWorksheet returns the formatted pages for the Shipment Summary Worksheet
func (SSWPPMComputer *SSWPPMComputer) FormatValuesShipmentSummaryWorksheet(shipmentSummaryFormData services.ShipmentSummaryFormData) (services.Page1Values, services.Page2Values) {
	page1 := FormatValuesShipmentSummaryWorksheetFormPage1(shipmentSummaryFormData)
	page2 := FormatValuesShipmentSummaryWorksheetFormPage2(shipmentSummaryFormData)

	return page1, page2
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

var newline = "\n\n"

// WorkSheetShipments is an object representing shipment line items on Shipment Summary Worksheet
type WorkSheetShipments struct {
	ShipmentNumberAndTypes      string
	PickUpDates                 string
	ShipmentWeights             string
	ShipmentWeightForObligation string
	CurrentShipmentStatuses     string
}

// WorkSheetShipment is an object representing specific shipment items on Shipment Summary Worksheet
type WorkSheetShipment struct {
	EstimatedIncentive    string
	MaxAdvance            string
	FinalIncentive        string
	AdvanceAmountReceived string
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

type Agent struct {
	Name  string
	Email string
	Date  string
	Phone string
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
func SSWGetEntitlement(grade internalmessages.OrderPayGrade, hasDependents bool, spouseHasProGear bool) services.SSWMaxWeightEntitlement {
	sswEntitlements := SSWMaxWeightEntitlement{}
	entitlements := models.GetWeightAllotment(grade)
	sswEntitlements.addLineItem("ProGear", entitlements.ProGearWeight)
	sswEntitlements.addLineItem("SpouseProGear", entitlements.ProGearWeightSpouse)
	if !hasDependents {
		sswEntitlements.addLineItem("Entitlement", entitlements.TotalWeightSelf)
		return services.SSWMaxWeightEntitlement(sswEntitlements)
	}
	sswEntitlements.addLineItem("Entitlement", entitlements.TotalWeightSelfPlusDependents)
	return services.SSWMaxWeightEntitlement(sswEntitlements)
}

// CalculateRemainingPPMEntitlement calculates the remaining PPM entitlement for PPM moves
// a PPMs remaining entitlement weight is equal to total entitlement - hhg weight
func CalculateRemainingPPMEntitlement(move models.Move, totalEntitlement unit.Pound) (unit.Pound, error) {

	var hhgActualWeight unit.Pound

	ppmActualWeight := models.GetTotalNetWeightForMove(move)

	switch ppmRemainingEntitlement := totalEntitlement - hhgActualWeight; {
	case ppmActualWeight < ppmRemainingEntitlement:
		return ppmActualWeight, nil
	case ppmRemainingEntitlement < 0:
		return 0, nil
	default:
		return ppmRemainingEntitlement, nil
	}
}

const (
	controlledUnclassifiedInformationText = "CONTROLLED UNCLASSIFIED INFORMATION"
)

// FormatValuesShipmentSummaryWorksheetFormPage1 formats the data for page 1 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage1(data services.ShipmentSummaryFormData) services.Page1Values {
	page1 := services.Page1Values{}
	page1.CUIBanner = controlledUnclassifiedInformationText
	page1.MaxSITStorageEntitlement = fmt.Sprintf("%02d Days in SIT", data.MaxSITStorageEntitlement)
	// We don't currently know what allows POV to be authorized, so we are hardcoding it to "No" to start
	page1.POVAuthorized = "No"
	page1.PreparationDate = FormatDate(data.PreparationDate)

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

	page1.AuthorizedOrigin = FormatLocation(data.CurrentDutyLocation)
	page1.AuthorizedDestination = data.NewDutyLocation.Name
	page1.NewDutyAssignment = data.NewDutyLocation.Name

	page1.WeightAllotment = FormatWeights(data.WeightAllotment.Entitlement)
	page1.WeightAllotmentProGear = FormatWeights(data.WeightAllotment.ProGear)
	page1.WeightAllotmentProgearSpouse = FormatWeights(data.WeightAllotment.SpouseProGear)
	page1.TotalWeightAllotment = FormatWeights(data.WeightAllotment.TotalWeight)

	formattedShipments := FormatAllShipments(data.PPMShipments)
	page1.ShipmentNumberAndTypes = formattedShipments.ShipmentNumberAndTypes
	page1.ShipmentPickUpDates = formattedShipments.PickUpDates
	page1.ShipmentCurrentShipmentStatuses = formattedShipments.CurrentShipmentStatuses
	formattedSIT := FormatAllSITS(data.PPMShipments)
	formattedShipment := FormatCurrentShipment(data.PPMShipment)
	page1.SITDaysInStorage = formattedSIT.DaysInStorage
	page1.SITEntryDates = formattedSIT.EntryDates
	page1.SITEndDates = formattedSIT.EndDates
	page1.SITNumberAndTypes = formattedShipments.ShipmentNumberAndTypes
	page1.ShipmentWeights = formattedShipments.ShipmentWeights
	page1.MaxObligationGCC100 = FormatWeights(data.WeightAllotment.TotalWeight) + " lbs; " + formattedShipment.EstimatedIncentive
	page1.ActualObligationGCC100 = formattedShipments.ShipmentWeightForObligation + " lbs; " + formattedShipment.FinalIncentive
	page1.MaxObligationGCCMaxAdvance = formattedShipment.MaxAdvance
	page1.ActualObligationAdvance = formattedShipment.AdvanceAmountReceived
	page1.MaxObligationSIT = fmt.Sprintf("%02d Days in SIT", data.MaxSITStorageEntitlement)
	page1.ActualObligationSIT = formattedSIT.DaysInStorage
	page1.TotalWeightAllotmentRepeat = page1.TotalWeightAllotment
	page1.PPMRemainingEntitlement = FormatWeights(data.PPMRemainingEntitlement)
	return page1
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

// FormatValuesShipmentSummaryWorksheetFormPage2 formats the data for page 2 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage2(data services.ShipmentSummaryFormData) services.Page2Values {

	expensesMap := SubTotalExpenses(data.MovingExpenses)
	agentInfo := FormatAgentInfo(data.MTOAgents)
	formattedShipments := FormatAllShipments(data.PPMShipments)

	page2 := services.Page2Values{}
	page2.CUIBanner = controlledUnclassifiedInformationText
	page2.TAC = derefStringTypes(data.Order.TAC)
	page2.SAC = derefStringTypes(data.Order.SAC)
	page2.PreparationDate = FormatDate(data.PreparationDate)
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
	page2.TrustedAgentName = agentInfo.Name
	page2.TrustedAgentDate = agentInfo.Date
	page2.TrustedAgentEmail = agentInfo.Email
	page2.TrustedAgentPhone = agentInfo.Phone
	page2.ServiceMemberSignature = FormatSignature(data.ServiceMember)
	page2.SignatureDate = FormatSignatureDate(data.SignedCertification.UpdatedAt)
	return page2
}

func formatMaxAdvance(estimatedIncentive *unit.Cents) string {
	if estimatedIncentive != nil {
		maxAdvance := float64(*estimatedIncentive) * 0.6
		return FormatDollars(maxAdvance / 100)
	}
	maxAdvanceString := "No Incentive Found"
	return maxAdvanceString

}

func FormatAgentInfo(agentArray []models.MTOAgent) Agent {
	agentObject := Agent{}
	if len(agentArray) == 0 {
		agentObject.Name = "No agent specified"
		agentObject.Email = "No agent specified"
		agentObject.Date = "No agent specified"
		agentObject.Phone = "No agent specified"
		return agentObject
	}

	agent := agentArray[0]

	switch {
	case agent.FirstName != nil && agent.LastName != nil:
		agentObject.Name = fmt.Sprintf("%s, %s", *agent.LastName, *agent.FirstName)
	case agent.FirstName == nil && agent.LastName == nil:
		agentObject.Name = "No name specified"
	case agent.FirstName == nil:
		agentObject.Name = fmt.Sprintf("No first name provided, Last Name: %s", *agent.LastName)
	case agent.LastName == nil:
		agentObject.Name = fmt.Sprintf("First Name: %s, No last name provided", *agent.FirstName)
	}

	agentObject.Email = getOrDefault(agent.Email, "No Email Specified")
	agentObject.Phone = getOrDefault(agent.Phone, "No Phone Specified")
	agentObject.Date = agent.UpdatedAt.Format("20060102")

	return agentObject
}

func getOrDefault(value *string, defaultValue string) string {
	if value != nil {
		return *value
	}
	return defaultValue
}

// FormatSignature formats a service member's signature for the Shipment Summary Worksheet
func FormatSignature(sm models.ServiceMember) string {
	first := derefStringTypes(sm.FirstName)
	last := derefStringTypes(sm.LastName)

	return fmt.Sprintf("%s %s electronically signed", first, last)
}

// FormatSignatureDate formats the date the service member electronically signed for the Shipment Summary Worksheet
func FormatSignatureDate(signature time.Time) string {
	dateLayout := "02 Jan 2006 at 3:04pm"
	dt := signature.Format(dateLayout)
	return dt
}

// FormatLocation formats AuthorizedOrigin and AuthorizedDestination for Shipment Summary Worksheet
func FormatLocation(dutyLocation models.DutyLocation) string {
	return fmt.Sprintf("%s, %s %s", dutyLocation.Name, dutyLocation.Address.State, dutyLocation.Address.PostalCode)
}

// FormatAddress retrieves a PPMShipment W2Address and formats it for the SSW Document
func FormatAddress(w2Address *models.Address) string {
	var addressString string

	if w2Address != nil {
		addressString = fmt.Sprintf("%s, %s %s%s %s %s%s",
			w2Address.StreetAddress1,
			nilOrValue(w2Address.StreetAddress2),
			nilOrValue(w2Address.StreetAddress3),
			w2Address.City,
			w2Address.State,
			nilOrValue(w2Address.Country),
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

func FormatCurrentShipment(ppm models.PPMShipment) WorkSheetShipment {
	formattedShipment := WorkSheetShipment{}

	if ppm.FinalIncentive != nil {
		formattedShipment.FinalIncentive = ppm.FinalIncentive.ToDollarString()
	} else {
		formattedShipment.FinalIncentive = "No final incentive."
	}
	if ppm.EstimatedIncentive != nil {
		formattedShipment.MaxAdvance = formatMaxAdvance(ppm.EstimatedIncentive)
		formattedShipment.EstimatedIncentive = ppm.EstimatedIncentive.ToDollarString()
	} else {
		formattedShipment.MaxAdvance = "Advance not available."
		formattedShipment.EstimatedIncentive = "No estimated incentive."
	}
	if ppm.AdvanceAmountReceived != nil {
		formattedShipment.AdvanceAmountReceived = ppm.AdvanceAmountReceived.ToDollarString()
	} else {
		formattedShipment.AdvanceAmountReceived = "No advance received."
	}

	return formattedShipment
}

// FormatAllShipments formats Shipment line items for the Shipment Summary Worksheet
func FormatAllShipments(ppms models.PPMShipments) WorkSheetShipments {
	totalShipments := len(ppms)
	formattedShipments := WorkSheetShipments{}
	formattedNumberAndTypes := make([]string, totalShipments)
	formattedPickUpDates := make([]string, totalShipments)
	formattedShipmentWeights := make([]string, totalShipments)
	formattedShipmentStatuses := make([]string, totalShipments)
	formattedShipmentTotalWeights := unit.Pound(0)
	var shipmentNumber int

	for _, ppm := range ppms {
		formattedNumberAndTypes[shipmentNumber] = FormatPPMNumberAndType(shipmentNumber)
		formattedPickUpDates[shipmentNumber] = FormatPPMPickupDate(ppm)
		formattedShipmentWeights[shipmentNumber] = FormatPPMWeight(ppm)
		formattedShipmentStatuses[shipmentNumber] = FormatCurrentPPMStatus(ppm)
		if ppm.EstimatedWeight != nil {
			formattedShipmentTotalWeights += *ppm.EstimatedWeight
		}
		shipmentNumber++
	}

	formattedShipments.ShipmentNumberAndTypes = strings.Join(formattedNumberAndTypes, newline)
	formattedShipments.PickUpDates = strings.Join(formattedPickUpDates, newline)
	formattedShipments.ShipmentWeights = strings.Join(formattedShipmentWeights, newline)
	formattedShipments.ShipmentWeightForObligation = FormatWeights(formattedShipmentTotalWeights)
	formattedShipments.CurrentShipmentStatuses = strings.Join(formattedShipmentStatuses, newline)
	return formattedShipments
}

// FormatAllSITs formats SIT line items for the Shipment Summary Worksheet
func FormatAllSITS(ppms models.PPMShipments) WorkSheetSIT {
	totalSITS := len(ppms)
	formattedSIT := WorkSheetSIT{}
	formattedSITNumberAndTypes := make([]string, totalSITS)
	formattedSITEntryDates := make([]string, totalSITS)
	formattedSITEndDates := make([]string, totalSITS)
	formattedSITDaysInStorage := make([]string, totalSITS)
	var sitNumber int

	for _, ppm := range ppms {
		// formattedSITNumberAndTypes[sitNumber] = FormatPPMNumberAndType(sitNumber)
		formattedSITEntryDates[sitNumber] = FormatSITEntryDate(ppm)
		formattedSITEndDates[sitNumber] = FormatSITEndDate(ppm)
		formattedSITDaysInStorage[sitNumber] = FormatSITDaysInStorage(ppm)

		sitNumber++
	}
	formattedSIT.NumberAndTypes = strings.Join(formattedSITNumberAndTypes, newline)
	formattedSIT.EntryDates = strings.Join(formattedSITEntryDates, newline)
	formattedSIT.EndDates = strings.Join(formattedSITEndDates, newline)
	formattedSIT.DaysInStorage = strings.Join(formattedSITDaysInStorage, newline)

	return formattedSIT
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
		expenseType, addToTotal := getExpenseType(expense)
		expenseDollarAmt := expense.Amount.ToDollarFloatNoRound()

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

// FormatPPMNumberAndType formats FormatShipmentNumberAndType for the Shipment Summary Worksheet
func FormatPPMNumberAndType(i int) string {
	return fmt.Sprintf("%02d - PPM", i+1)
}

// FormatSITNumberAndType formats FormatSITNumberAndType for the Shipment Summary Worksheet
func FormatSITNumberAndType(i int) string {
	return fmt.Sprintf("%02d - SIT", i+1)
}

// FormatPPMWeight formats a ppms NetWeight for the Shipment Summary Worksheet
func FormatPPMWeight(ppm models.PPMShipment) string {
	if ppm.EstimatedWeight != nil {
		wtg := FormatWeights(unit.Pound(*ppm.EstimatedWeight))
		return fmt.Sprintf("%s lbs - Estimated", wtg)
	}
	return ""
}

// FormatPPMPickupDate formats a shipments ActualPickupDate for the Shipment Summary Worksheet
func FormatPPMPickupDate(ppm models.PPMShipment) string {
	return FormatDate(ppm.ExpectedDepartureDate)
}

// FormatSITEntryDate formats a SIT EstimatedEntryDate for the Shipment Summary Worksheet
func FormatSITEntryDate(ppm models.PPMShipment) string {
	if ppm.SITEstimatedEntryDate == nil {
		return "No Entry Data" // Return string if no SIT attached
	}
	return FormatDate(*ppm.SITEstimatedEntryDate)
}

// FormatSITEndDate formats a SIT EstimatedPickupDate for the Shipment Summary Worksheet
func FormatSITEndDate(ppm models.PPMShipment) string {
	if ppm.SITEstimatedDepartureDate == nil {
		return "No Departure Data" // Return string if no SIT attached
	}
	return FormatDate(*ppm.SITEstimatedDepartureDate)
}

// FormatSITDaysInStorage formats a SIT DaysInStorage for the Shipment Summary Worksheet
func FormatSITDaysInStorage(ppm models.PPMShipment) string {
	if ppm.SITEstimatedEntryDate == nil || ppm.SITEstimatedDepartureDate == nil {
		return "No Entry/Departure Data" // Return string if no SIT attached
	}
	firstDate := ppm.SITEstimatedDepartureDate
	secondDate := *ppm.SITEstimatedEntryDate
	difference := firstDate.Sub(secondDate)
	formattedDifference := fmt.Sprintf("Days: %d\n", int64(difference.Hours()/24)+1)
	return formattedDifference
}

// FormatOrdersTypeAndOrdersNumber formats OrdersTypeAndOrdersNumber for Shipment Summary Worksheet
func FormatOrdersTypeAndOrdersNumber(order models.Order) string {
	issuingBranch := FormatOrdersType(order)
	ordersNumber := derefStringTypes(order.OrdersNumber)
	return fmt.Sprintf("%s/%s", issuingBranch, ordersNumber)
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
func (SSWPPMComputer *SSWPPMComputer) ComputeObligations(_ appcontext.AppContext, _ services.ShipmentSummaryFormData, _ route.Planner) (obligation services.Obligations, err error) {
	// Obligations must remain test data until new computer system is finished
	obligations := services.Obligations{
		ActualObligation:           services.Obligation{Gcc: 123, SIT: 123, Miles: unit.Miles(123456)},
		MaxObligation:              services.Obligation{Gcc: 456, SIT: 456, Miles: unit.Miles(123456)},
		NonWinningActualObligation: services.Obligation{Gcc: 789, SIT: 789, Miles: unit.Miles(12345)},
		NonWinningMaxObligation:    services.Obligation{Gcc: 1000, SIT: 1000, Miles: unit.Miles(12345)},
	}
	return obligations, nil
}

// FetchDataShipmentSummaryWorksheetFormData fetches the pages for the Shipment Summary Worksheet for a given Move ID
func (SSWPPMComputer *SSWPPMComputer) FetchDataShipmentSummaryWorksheetFormData(appCtx appcontext.AppContext, _ *auth.Session, ppmShipmentID uuid.UUID) (*services.ShipmentSummaryFormData, error) {

	ppmShipment := models.PPMShipment{}
	dbQErr := appCtx.DB().Q().Eager(
		"Shipment.MoveTaskOrder.Orders.ServiceMember",
		"Shipment.MoveTaskOrder.Orders.NewDutyLocation.Address",
		"Shipment.MoveTaskOrder.Orders.OriginDutyLocation.Address",
		"Shipment.MTOAgents",
		"W2Address",
		"SignedCertification",
		"MovingExpenses",
	).Find(&ppmShipment, ppmShipmentID)

	if dbQErr != nil {
		if errors.Cause(dbQErr).Error() == models.RecordNotFoundErrorString {
			return nil, models.ErrFetchNotFound
		}
		return nil, dbQErr
	}

	serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
	if ppmShipment.Shipment.MoveTaskOrder.Orders.Grade == nil {
		return nil, errors.New("order for requested shipment summary worksheet data does not have a pay grade attached")
	}

	weightAllotment := SSWGetEntitlement(*ppmShipment.Shipment.MoveTaskOrder.Orders.Grade, ppmShipment.Shipment.MoveTaskOrder.Orders.HasDependents, ppmShipment.Shipment.MoveTaskOrder.Orders.SpouseHasProGear)
	ppmRemainingEntitlement, err := CalculateRemainingPPMEntitlement(ppmShipment.Shipment.MoveTaskOrder, weightAllotment.TotalWeight)
	if err != nil {
		return nil, err
	}

	maxSit, err := CalculateShipmentSITAllowance(appCtx, ppmShipment.Shipment)
	if err != nil {
		return nil, err
	}

	// DOES NOT INCLUDE PPPO/PPSO SIGNATURE
	signedCertification := ppmShipment.SignedCertification

	var ppmShipments []models.PPMShipment

	ppmShipments = append(ppmShipments, ppmShipment)
	if ppmShipment.Shipment.MoveTaskOrder.Orders.OriginDutyLocation == nil {
		return nil, errors.New("order for PPM shipment does not have a origin duty location attached")
	}
	ssd := services.ShipmentSummaryFormData{
		ServiceMember:            serviceMember,
		Order:                    ppmShipment.Shipment.MoveTaskOrder.Orders,
		Move:                     ppmShipment.Shipment.MoveTaskOrder,
		CurrentDutyLocation:      *ppmShipment.Shipment.MoveTaskOrder.Orders.OriginDutyLocation,
		NewDutyLocation:          ppmShipment.Shipment.MoveTaskOrder.Orders.NewDutyLocation,
		WeightAllotment:          weightAllotment,
		PPMShipment:              ppmShipment,
		PPMShipments:             ppmShipments,
		W2Address:                ppmShipment.W2Address,
		MovingExpenses:           ppmShipment.MovingExpenses,
		MTOAgents:                ppmShipment.Shipment.MTOAgents,
		SignedCertification:      *signedCertification,
		PPMRemainingEntitlement:  ppmRemainingEntitlement,
		MaxSITStorageEntitlement: maxSit,
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
func (SSWPPMGenerator *SSWPPMGenerator) FillSSWPDFForm(Page1Values services.Page1Values, Page2Values services.Page2Values) (sswfile afero.File, pdfInfo *pdfcpu.PDFInfo, err error) {

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
		Source:   "SSWPDFTemplate.pdf",
		Version:  "pdfcpu v0.6.0 dev",
		Creation: "2024-03-08 17:36:47 UTC",
		Producer: "macOS Version 13.5 (Build 22G74) Quartz PDFContext, AppendMode 1.1",
	}

	var sswCheckbox = []checkbox{
		{
			Pages:   []int{2},
			ID:      "797",
			Name:    "EDOther",
			Value:   true,
			Default: false,
			Locked:  false,
		},
	}

	formData := pdFData{ // This is unique to each PDF template, must be found for new templates using PDFCPU's export function used on the template (can be done through CLI)
		Header: sswHeader,
		Forms: []form{
			{ // Dynamically loops, creates, and aggregates json for text fields, merges page 1 and 2
				TextField: mergeTextFields(createTextFields(Page1Values, 1), createTextFields(Page2Values, 2)),
			},
			// The following is the structure for using a Checkbox field
			{
				Checkbox: sswCheckbox,
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

	// pdfInfo.PageCount is a great way to tell whether returned PDF is corrupted
	pdfInfoResult, err := SSWPPMGenerator.generator.GetPdfFileInfo(SSWWorksheet.Name())
	if err != nil || pdfInfoResult.PageCount != 2 {
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
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i).Interface()

		textFieldEntry := textField{
			Pages:     pages,
			ID:        fmt.Sprintf("%d", len(textFields)+1),
			Name:      field.Name,
			Value:     fmt.Sprintf("%v", value),
			Multiline: false,
			Locked:    false,
		}

		textFields = append(textFields, textFieldEntry)
	}

	return textFields
}

// MergeTextFields merges page 1 and page 2 data
func mergeTextFields(fields1, fields2 []textField) []textField {
	return append(fields1, fields2...)
}

// createAssetByteReader creates a new byte reader based on the TemplateImagePath of the formLayout
func createAssetByteReader(path string) (*bytes.Reader, error) {
	asset, err := assets.Asset(path)
	if err != nil {
		return nil, errors.Wrap(err, "error creating asset from path; check image path : "+path)
	}

	return bytes.NewReader(asset), nil
}
