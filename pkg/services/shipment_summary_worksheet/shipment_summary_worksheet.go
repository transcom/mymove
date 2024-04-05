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
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
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
func NewSSWPPMGenerator() (services.SSWPPMGenerator, error) {
	// Generator and dependencies must be initiated to handle memory filesystem for AWS
	storer := storage.NewMemory(storage.NewMemoryParams("", ""))
	userUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	generator, err := paperwork.NewGenerator(userUploader.Uploader())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	templateReader, err := createAssetByteReader("paperwork/formtemplates/SSWPDFTemplate.pdf")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &SSWPPMGenerator{
		generator:      *generator,
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

// Page1Values is an object representing a Shipment Summary Worksheet
type Page1Values struct {
	CUIBanner                       string
	ServiceMemberName               string
	MaxSITStorageEntitlement        string
	PreferredPhoneNumber            string
	PreferredEmail                  string
	DODId                           string
	ServiceBranch                   string
	RankGrade                       string
	IssuingBranchOrAgency           string
	OrdersIssueDate                 string
	OrdersTypeAndOrdersNumber       string
	AuthorizedOrigin                string
	AuthorizedDestination           string
	NewDutyAssignment               string
	WeightAllotment                 string
	WeightAllotmentProgear          string
	WeightAllotmentProgearSpouse    string
	TotalWeightAllotment            string
	POVAuthorized                   string
	ShipmentNumberAndTypes          string
	ShipmentPickUpDates             string
	ShipmentWeights                 string
	ShipmentCurrentShipmentStatuses string
	SITNumberAndTypes               string
	SITEntryDates                   string
	SITEndDates                     string
	SITDaysInStorage                string
	PreparationDate                 string
	MaxObligationGCC100             string
	TotalWeightAllotmentRepeat      string
	MaxObligationGCC95              string
	MaxObligationSIT                string
	MaxObligationGCCMaxAdvance      string
	PPMRemainingEntitlement         string
	ActualObligationGCC100          string
	ActualObligationGCC95           string
	ActualObligationAdvance         string
	ActualObligationSIT             string
	MileageTotal                    string
}

// WorkSheetShipments is an object representing shipment line items on Shipment Summary Worksheet
type WorkSheetShipments struct {
	ShipmentNumberAndTypes  string
	PickUpDates             string
	ShipmentWeights         string
	CurrentShipmentStatuses string
}

// WorkSheetSIT is an object representing SIT on the Shipment Summary Worksheet
type WorkSheetSIT struct {
	NumberAndTypes string
	EntryDates     string
	EndDates       string
	DaysInStorage  string
}

// Page2Values is an object representing a Shipment Summary Worksheet
type Page2Values struct {
	CUIBanner       string
	PreparationDate string
	TAC             string
	SAC             string
	FormattedMovingExpenses
}

// Dollar represents a type for dollar monetary unit
type Dollar float64

// String is a string representation of a Dollar
func (d Dollar) String() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("$%.2f", d)
}

// FormattedMovingExpenses is an object representing the service member's moving expenses formatted for the SSW
type FormattedMovingExpenses struct {
	ContractedExpenseMemberPaid Dollar
	ContractedExpenseGTCCPaid   Dollar
	RentalEquipmentMemberPaid   Dollar
	RentalEquipmentGTCCPaid     Dollar
	PackingMaterialsMemberPaid  Dollar
	PackingMaterialsGTCCPaid    Dollar
	WeighingFeesMemberPaid      Dollar
	WeighingFeesGTCCPaid        Dollar
	GasMemberPaid               Dollar
	GasGTCCPaid                 Dollar
	TollsMemberPaid             Dollar
	TollsGTCCPaid               Dollar
	OilMemberPaid               Dollar
	OilGTCCPaid                 Dollar
	OtherMemberPaid             Dollar
	OtherGTCCPaid               Dollar
	TotalMemberPaid             Dollar
	TotalGTCCPaid               Dollar
	TotalMemberPaidRepeated     Dollar
	TotalGTCCPaidRepeated       Dollar
	TotalPaidNonSIT             Dollar
	TotalMemberPaidSIT          Dollar
	TotalGTCCPaidSIT            Dollar
	TotalPaidSIT                Dollar
}

// FormattedOtherExpenses is an object representing the other moving expenses formatted for the SSW
type FormattedOtherExpenses struct {
	Descriptions string
	AmountsPaid  string
}

// Page3Values is an object representing a Shipment Summary Worksheet
type Page3Values struct {
	CUIBanner              string
	PreparationDate        string
	ServiceMemberSignature string
	SignatureDate          string
	FormattedOtherExpenses
}

// ShipmentSummaryFormData is a container for the various objects required for the a Shipment Summary Worksheet
type ShipmentSummaryFormData struct {
	ServiceMember           models.ServiceMember
	Order                   models.Order
	Move                    models.Move
	CurrentDutyLocation     models.DutyLocation
	NewDutyLocation         models.DutyLocation
	WeightAllotment         SSWMaxWeightEntitlement
	PPMShipments            models.PPMShipments
	PreparationDate         time.Time
	Obligations             Obligations
	MovingExpenses          models.MovingExpenses
	PPMRemainingEntitlement unit.Pound
	SignedCertification     models.SignedCertification
}

// Obligations is an object representing the winning and non-winning Max Obligation and Actual Obligation sections of the shipment summary worksheet
type Obligations struct {
	MaxObligation              Obligation
	ActualObligation           Obligation
	NonWinningMaxObligation    Obligation
	NonWinningActualObligation Obligation
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
	if !hasDependents {
		sswEntitlements.addLineItem("Entitlement", entitlements.TotalWeightSelf)
		return services.SSWMaxWeightEntitlement(sswEntitlements)
	}
	sswEntitlements.addLineItem("Entitlement", entitlements.TotalWeightSelfPlusDependents)
	if spouseHasProGear {
		sswEntitlements.addLineItem("SpouseProGear", entitlements.ProGearWeightSpouse)
	}
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
	page1.MaxSITStorageEntitlement = "90 days per each shipment"
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
	page1.WeightAllotmentProgear = FormatWeights(data.WeightAllotment.ProGear)
	page1.WeightAllotmentProgearSpouse = FormatWeights(data.WeightAllotment.SpouseProGear)
	page1.TotalWeightAllotment = FormatWeights(data.WeightAllotment.TotalWeight)

	formattedShipments := FormatAllShipments(data.PPMShipments)
	page1.ShipmentNumberAndTypes = formattedShipments.ShipmentNumberAndTypes
	page1.ShipmentPickUpDates = formattedShipments.PickUpDates
	page1.ShipmentCurrentShipmentStatuses = formattedShipments.CurrentShipmentStatuses
	formattedSIT := FormatAllSITS(data.PPMShipments)

	page1.SITDaysInStorage = formattedSIT.DaysInStorage
	page1.SITEntryDates = formattedSIT.EntryDates
	page1.SITEndDates = formattedSIT.EndDates
	// page1.SITNumberAndTypes
	page1.ShipmentWeights = formattedShipments.ShipmentWeights
	// Obligations cannot be used at this time, require new computer setup.
	page1.TotalWeightAllotmentRepeat = page1.TotalWeightAllotment
	actualObligations := data.Obligations.ActualObligation
	page1.PPMRemainingEntitlement = FormatWeights(data.PPMRemainingEntitlement)
	page1.MileageTotal = actualObligations.Miles.String()
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
	page2 := services.Page2Values{}
	page2.CUIBanner = controlledUnclassifiedInformationText
	page2.TAC = derefStringTypes(data.Order.TAC)
	page2.SAC = derefStringTypes(data.Order.SAC)
	page2.PreparationDate = FormatDate(data.PreparationDate)
	page2.TotalMemberPaidRepeated = page2.TotalMemberPaid
	page2.TotalGTCCPaidRepeated = page2.TotalGTCCPaid
	page2.ServiceMemberSignature = FormatSignature(data.ServiceMember)
	page2.SignatureDate = FormatSignatureDate(data.SignedCertification)
	return page2
}

// FormatSignature formats a service member's signature for the Shipment Summary Worksheet
func FormatSignature(sm models.ServiceMember) string {
	first := derefStringTypes(sm.FirstName)
	last := derefStringTypes(sm.LastName)

	return fmt.Sprintf("%s %s electronically signed", first, last)
}

// FormatSignatureDate formats the date the service member electronically signed for the Shipment Summary Worksheet
func FormatSignatureDate(signature models.SignedCertification) string {
	dateLayout := "02 Jan 2006 at 3:04pm"
	dt := signature.Date.Format(dateLayout)
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
		return "" // Return an empty string if no W2 address
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

// FormatAllShipments formats Shipment line items for the Shipment Summary Worksheet
func FormatAllShipments(ppms models.PPMShipments) WorkSheetShipments {
	totalShipments := len(ppms)
	formattedShipments := WorkSheetShipments{}
	formattedNumberAndTypes := make([]string, totalShipments)
	formattedPickUpDates := make([]string, totalShipments)
	formattedShipmentWeights := make([]string, totalShipments)
	formattedShipmentStatuses := make([]string, totalShipments)
	var shipmentNumber int

	for _, ppm := range ppms {
		formattedNumberAndTypes[shipmentNumber] = FormatPPMNumberAndType(shipmentNumber)
		formattedPickUpDates[shipmentNumber] = FormatPPMPickupDate(ppm)
		formattedShipmentWeights[shipmentNumber] = FormatPPMWeight(ppm)
		formattedShipmentStatuses[shipmentNumber] = FormatCurrentPPMStatus(ppm)
		shipmentNumber++
	}

	formattedShipments.ShipmentNumberAndTypes = strings.Join(formattedNumberAndTypes, newline)
	formattedShipments.PickUpDates = strings.Join(formattedPickUpDates, newline)
	formattedShipments.ShipmentWeights = strings.Join(formattedShipmentWeights, newline)
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

// SubTotalExpenses groups moving expenses by type and payment method
func SubTotalExpenses(expenseDocuments models.MovingExpenses) map[string]float64 {
	var expenseType string
	totals := make(map[string]float64)
	for _, expense := range expenseDocuments {
		expenseType = getExpenseType(expense)
		expenseDollarAmt := expense.Amount.ToDollarFloatNoRound()
		totals[expenseType] += expenseDollarAmt
		// addToGrandTotal(totals, expenseType, expenseDollarAmt)
	}
	return totals
}

func getExpenseType(expense models.MovingExpense) string {
	expenseType := FormatEnum(string(*expense.MovingExpenseType), "")
	paidWithGTCC := expense.PaidWithGTCC
	if paidWithGTCC != nil {
		if *paidWithGTCC {
			return fmt.Sprintf("%s%s", expenseType, "GTCCPaid")
		}
	}

	return fmt.Sprintf("%s%s", expenseType, "MemberPaid")
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
		return fmt.Sprintf("%s lbs - FINAL", wtg)
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
	formattedDifference := fmt.Sprintf("Days: %d\n", int64(difference.Hours()/24))
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
		"Shipment.MoveTaskOrder",
		"Shipment.MoveTaskOrder.Orders",
		"Shipment.MoveTaskOrder.Orders.NewDutyLocation.Address",
		"Shipment.MoveTaskOrder.Orders.OriginDutyLocation.Address",
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

	// Signed Certification needs to be updated
	// signedCertification, err := models.FetchSignedCertificationsPPMPayment(appCtx.DB(), session, ppmShipment.Shipment.MoveTaskOrderID)
	// if err != nil {
	// 	return ShipmentSummaryFormData{}, err
	// }
	// if signedCertification == nil {
	// 	return ShipmentSummaryFormData{},
	// 		errors.New("shipment summary worksheet: signed certification is nil")
	// }

	var ppmShipments []models.PPMShipment

	ppmShipments = append(ppmShipments, ppmShipment)
	if ppmShipment.Shipment.MoveTaskOrder.Orders.OriginDutyLocation == nil {
		return nil, errors.New("order for PPM shipment does not have a origin duty location attached")
	}
	ssd := services.ShipmentSummaryFormData{
		ServiceMember:       serviceMember,
		Order:               ppmShipment.Shipment.MoveTaskOrder.Orders,
		Move:                ppmShipment.Shipment.MoveTaskOrder,
		CurrentDutyLocation: *ppmShipment.Shipment.MoveTaskOrder.Orders.OriginDutyLocation,
		NewDutyLocation:     ppmShipment.Shipment.MoveTaskOrder.Orders.NewDutyLocation,
		WeightAllotment:     weightAllotment,
		PPMShipments:        ppmShipments,
		W2Address:           ppmShipment.W2Address,
		// SignedCertification:     *signedCertification,
		PPMRemainingEntitlement: ppmRemainingEntitlement,
	}
	return &ssd, nil
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
		Creation: "2024-01-22 21:49:12 UTC",
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
	SSWWorksheet, err := SSWPPMGenerator.generator.FillPDFForm(jsonData, SSWPPMGenerator.templateReader)
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
