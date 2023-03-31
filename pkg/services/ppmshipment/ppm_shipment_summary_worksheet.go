package ppmshipment

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// ppmShipmentSummaryWorksheetCreator is the concrete implementation of the services.PPMShipmentSummaryWorksheetCreator interface
type ppmShipmentSummaryWorksheetCreator struct{}

// NewPPMShipmentSummaryWorksheetCreator creates a new struct
func NewPPMShipmentSummaryWorksheetCreator() services.PPMShipmentSummaryWorksheetCreator {
	return &ppmShipmentSummaryWorksheetCreator{}
}

// FormatValuesShipmentSummaryWorksheet returns the formatted pages for the Shipment Summary Worksheet
func (p *ppmShipmentSummaryWorksheetCreator) FormatValuesShipmentSummaryWorksheet(shipmentSummaryFormData ShipmentSummaryFormData) (ShipmentSummaryWorksheetPage1Values, ShipmentSummaryWorksheetPage2Values, ShipmentSummaryWorksheetPage3Values, error) {
	page1 := FormatValuesShipmentSummaryWorksheetFormPage1(shipmentSummaryFormData)
	page2 := FormatValuesShipmentSummaryWorksheetFormPage2(shipmentSummaryFormData)
	page3 := FormatValuesShipmentSummaryWorksheetFormPage3(shipmentSummaryFormData)

	return page1, page2, page3, nil
}

// ShipmentSummaryWorksheetPage1Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage1Values struct {
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

// ShipmentSummaryWorkSheetShipments is an object representing shipment line items on Shipment Summary Worksheet
type ShipmentSummaryWorkSheetShipments struct {
	ShipmentNumberAndTypes  string
	PickUpDates             string
	ShipmentWeights         string
	CurrentShipmentStatuses string
}

// ShipmentSummaryWorkSheetSIT is an object representing SIT on the Shipment Summary Worksheet
type ShipmentSummaryWorkSheetSIT struct {
	NumberAndTypes string
	EntryDates     string
	EndDates       string
	DaysInStorage  string
}

// ShipmentSummaryWorksheetPage2Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage2Values struct {
	CUIBanner       string
	PreparationDate string
	TAC             string
	SAC             string
	FormattedMovingExpenses
}

// Dollar represents a type for dollar monetary unit
type Dollar float64

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

// ShipmentSummaryWorksheetPage3Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage3Values struct {
	CUIBanner              string
	PreparationDate        string
	ServiceMemberSignature string
	SignatureDate          string
	FormattedOtherExpenses
}

// ShipmentSummaryFormData is a container for the various objects required for the a Shipment Summary Worksheet
type ShipmentSummaryFormData struct {
	ServiceMember       models.ServiceMember
	Order               models.Order
	CurrentDutyLocation models.DutyLocation
	NewDutyLocation     models.DutyLocation
	WeightAllotment     models.SSWMaxWeightEntitlement
	PPMShipments        models.PPMShipments
	PreparationDate     time.Time
	Obligations         models.Obligations
	MovingExpenses      models.MovingExpenses
	// PPMRemainingEntitlement unit.Pound
	SignedCertification models.SignedCertification
}

// FetchDataShipmentSummaryWorksheetFormData fetches the pages for the Shipment Summary Worksheet for a given Move ID
func (f *ppmShipmentSummaryWorksheetCreator) FetchDataShipmentSummaryWorksheetFormData(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (ShipmentSummaryFormData, error) {
	ppmShipment, err := FetchPPMShipmentFromMTOShipmentID(appCtx, mtoShipmentID)
	if err != nil {
		return ShipmentSummaryFormData{}, apperror.UnprocessableEntityError{}
	}
	// _, authErr := FetchOrderForUser(db, session, move.OrdersID)
	// if authErr != nil {
	// 	return ShipmentSummaryFormData{}, authErr
	// }

	orders := ppmShipment.Shipment.MoveTaskOrder.Orders
	serviceMember := orders.ServiceMember
	var rank models.ServiceMemberRank
	var weightAllotment models.SSWMaxWeightEntitlement
	if serviceMember.Rank != nil {
		rank = models.ServiceMemberRank(*serviceMember.Rank)
		// QUESTION: would we want to use orders.SpouseHasProgear vs ppmShipment.SpouseProgearWeight
		weightAllotment = SSWGetEntitlement(rank, orders.HasDependents, orders.SpouseHasProGear)
	}

	// Question: where should we now pull this info from? For remaining PPM entitlement
	// ppmRemainingEntitlement, err := CalculateRemainingPPMEntitlement(move, weightAllotment.TotalWeight)
	// if err != nil {
	// 	return ShipmentSummaryFormData{}, err
	// }

	signedCertification := ppmShipment.SignedCertification
	if signedCertification == nil {
		return ShipmentSummaryFormData{},
			apperror.NewConflictError(ppmShipment.ID, "shipment summary worksheet: signed certification is nil")
	}
	ssd := ShipmentSummaryFormData{
		ServiceMember:       serviceMember,
		Order:               orders,
		CurrentDutyLocation: serviceMember.DutyLocation,
		NewDutyLocation:     orders.NewDutyLocation,
		WeightAllotment:     weightAllotment,
		PPMShipments:        models.PPMShipments{*ppmShipment},
		SignedCertification: *signedCertification,
		// PPMRemainingEntitlement: ppmRemainingEntitlement,
	}
	return ssd, nil
}

// SSWGetEntitlement calculates the entitlement for the shipment summary worksheet based on the parameters of
// a move (hasDependents, spouseHasProGear)
func SSWGetEntitlement(rank models.ServiceMemberRank, hasDependents bool, spouseHasProGear bool) models.SSWMaxWeightEntitlement {
	sswEntitlements := models.SSWMaxWeightEntitlement{}
	entitlements := models.GetWeightAllotment(rank)
	sswEntitlements.AddLineItem("ProGear", entitlements.ProGearWeight)
	if !hasDependents {
		sswEntitlements.AddLineItem("Entitlement", entitlements.TotalWeightSelf)
		return sswEntitlements
	}
	sswEntitlements.AddLineItem("Entitlement", entitlements.TotalWeightSelfPlusDependents)
	if spouseHasProGear {
		sswEntitlements.AddLineItem("SpouseProGear", entitlements.ProGearWeightSpouse)
	}
	return sswEntitlements
}

// QUESTION - how does this concept map to our current database structure?
// // CalculateRemainingPPMEntitlement calculates the remaining PPM entitlement for PPM moves
// // a PPMs remaining entitlement weight is equal to total entitlement - hhg weight
// func CalculateRemainingPPMEntitlement(move Move, totalEntitlement unit.Pound) (unit.Pound, error) {
// 	var hhgActualWeight unit.Pound

// 	var ppmActualWeight unit.Pound
// 	if len(move.PersonallyProcuredMoves) > 0 {
// 		if move.PersonallyProcuredMoves[0].NetWeight == nil {
// 			return ppmActualWeight, errors.Errorf("PPM %s does not have NetWeight", move.PersonallyProcuredMoves[0].ID)
// 		}
// 		ppmActualWeight = unit.Pound(*move.PersonallyProcuredMoves[0].NetWeight)
// 	}

// 	switch ppmRemainingEntitlement := totalEntitlement - hhgActualWeight; {
// 	case ppmActualWeight < ppmRemainingEntitlement:
// 		return ppmActualWeight, nil
// 	case ppmRemainingEntitlement < 0:
// 		return 0, nil
// 	default:
// 		return ppmRemainingEntitlement, nil
// 	}
// }

const (
	controlledUnclassifiedInformationText = "CONTROLLED UNCLASSIFIED INFORMATION"
)

// FormatValuesShipmentSummaryWorksheetFormPage1 formats the data for page 1 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage1(data ShipmentSummaryFormData) ShipmentSummaryWorksheetPage1Values {
	page1 := ShipmentSummaryWorksheetPage1Values{}
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
	page1.RankGrade = FormatRank(data.ServiceMember.Rank)

	page1.IssuingBranchOrAgency = FormatServiceMemberAffiliation(sm.Affiliation)
	page1.OrdersIssueDate = FormatDate(data.Order.IssueDate)
	page1.OrdersTypeAndOrdersNumber = FormatOrdersTypeAndOrdersNumber(data.Order)

	page1.AuthorizedOrigin = FormatLocation(data.CurrentDutyLocation)
	page1.AuthorizedDestination = FormatLocation(data.NewDutyLocation)
	page1.NewDutyAssignment = FormatLocation(data.NewDutyLocation)

	page1.WeightAllotment = FormatWeights(data.WeightAllotment.Entitlement)
	page1.WeightAllotmentProgear = FormatWeights(data.WeightAllotment.ProGear)
	page1.WeightAllotmentProgearSpouse = FormatWeights(data.WeightAllotment.SpouseProGear)
	page1.TotalWeightAllotment = FormatWeights(data.WeightAllotment.TotalWeight)

	formattedShipments := FormatAllShipments(data.PPMShipments)
	page1.ShipmentNumberAndTypes = formattedShipments.ShipmentNumberAndTypes
	page1.ShipmentPickUpDates = formattedShipments.PickUpDates
	page1.ShipmentCurrentShipmentStatuses = formattedShipments.CurrentShipmentStatuses
	page1.ShipmentWeights = formattedShipments.ShipmentWeights

	maxObligations := data.Obligations.MaxObligation
	page1.MaxObligationGCC100 = FormatDollars(maxObligations.GCC100())
	page1.TotalWeightAllotmentRepeat = page1.TotalWeightAllotment
	page1.MaxObligationGCC95 = FormatDollars(maxObligations.GCC95())
	page1.MaxObligationSIT = FormatDollars(maxObligations.FormatSIT())
	page1.MaxObligationGCCMaxAdvance = FormatDollars(maxObligations.MaxAdvance())

	actualObligations := data.Obligations.ActualObligation
	page1.ActualObligationGCC100 = FormatDollars(actualObligations.GCC100())
	// page1.PPMRemainingEntitlement = FormatWeights(data.PPMRemainingEntitlement)
	page1.ActualObligationGCC95 = FormatDollars(actualObligations.GCC95())
	page1.ActualObligationSIT = FormatDollars(actualObligations.FormatSIT())
	page1.ActualObligationAdvance = formatActualObligationAdvance(data)
	page1.MileageTotal = actualObligations.Miles.String()
	return page1
}

// UPDATED - to use the new PPMShipment model and the AdvanceAmountRequested field
func formatActualObligationAdvance(data ShipmentSummaryFormData) string {
	if len(data.PPMShipments) > 0 && data.PPMShipments[0].AdvanceAmountRequested != nil {
		advance := data.PPMShipments[0].AdvanceAmountRequested.ToDollarFloatNoRound()
		return FormatDollars(advance)
	}
	return FormatDollars(0)
}

// UPDATED - since this in in a service object now instead of the model package,
// "models." has been prefixed where needed
// FormatRank formats the service member's rank for Shipment Summary Worksheet
func FormatRank(rank *models.ServiceMemberRank) string {
	var rankDisplayValue = map[models.ServiceMemberRank]string{
		models.ServiceMemberRankE1:                      "E-1",
		models.ServiceMemberRankE2:                      "E-2",
		models.ServiceMemberRankE3:                      "E-3",
		models.ServiceMemberRankE4:                      "E-4",
		models.ServiceMemberRankE5:                      "E-5",
		models.ServiceMemberRankE6:                      "E-6",
		models.ServiceMemberRankE7:                      "E-7",
		models.ServiceMemberRankE8:                      "E-8",
		models.ServiceMemberRankE9:                      "E-9",
		models.ServiceMemberRankE9SPECIALSENIORENLISTED: "E-9 (Special Senior Enlisted)",
		models.ServiceMemberRankO1ACADEMYGRADUATE:       "O-1 or Service Academy Graduate",
		models.ServiceMemberRankO2:                      "O-2",
		models.ServiceMemberRankO3:                      "O-3",
		models.ServiceMemberRankO4:                      "O-4",
		models.ServiceMemberRankO5:                      "O-5",
		models.ServiceMemberRankO6:                      "O-6",
		models.ServiceMemberRankO7:                      "O-7",
		models.ServiceMemberRankO8:                      "O-8",
		models.ServiceMemberRankO9:                      "O-9",
		models.ServiceMemberRankO10:                     "O-10",
		models.ServiceMemberRankW1:                      "W-1",
		models.ServiceMemberRankW2:                      "W-2",
		models.ServiceMemberRankW3:                      "W-3",
		models.ServiceMemberRankW4:                      "W-4",
		models.ServiceMemberRankW5:                      "W-5",
		models.ServiceMemberRankAVIATIONCADET:           "Aviation Cadet",
		models.ServiceMemberRankCIVILIANEMPLOYEE:        "Civilian Employee",
		models.ServiceMemberRankACADEMYCADET:            "Service Academy Cadet",
		models.ServiceMemberRankMIDSHIPMAN:              "Midshipman",
	}
	if rank != nil {
		return rankDisplayValue[*rank]
	}
	return ""
}

// FormatValuesShipmentSummaryWorksheetFormPage2 formats the data for page 2 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage2(data ShipmentSummaryFormData) ShipmentSummaryWorksheetPage2Values {
	page2 := ShipmentSummaryWorksheetPage2Values{}
	page2.CUIBanner = controlledUnclassifiedInformationText
	page2.TAC = derefStringTypes(data.Order.TAC)
	page2.SAC = derefStringTypes(data.Order.SAC)
	page2.PreparationDate = FormatDate(data.PreparationDate)
	page2.TotalMemberPaidRepeated = page2.TotalMemberPaid
	page2.TotalGTCCPaidRepeated = page2.TotalGTCCPaid
	return page2
}

// FormatValuesShipmentSummaryWorksheetFormPage3 formats the data for page 2 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage3(data ShipmentSummaryFormData) ShipmentSummaryWorksheetPage3Values {
	page3 := ShipmentSummaryWorksheetPage3Values{}
	page3.CUIBanner = controlledUnclassifiedInformationText
	page3.PreparationDate = FormatDate(data.PreparationDate)
	page3.ServiceMemberSignature = FormatSignature(data.ServiceMember)
	page3.SignatureDate = FormatSignatureDate(data.SignedCertification)
	return page3
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
func FormatAllShipments(ppms models.PPMShipments) ShipmentSummaryWorkSheetShipments {
	totalShipments := len(ppms)
	formattedShipments := ShipmentSummaryWorkSheetShipments{}
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

	formattedShipments.ShipmentNumberAndTypes = strings.Join(formattedNumberAndTypes, "\n\n")
	formattedShipments.PickUpDates = strings.Join(formattedPickUpDates, "\n\n")
	formattedShipments.ShipmentWeights = strings.Join(formattedShipmentWeights, "\n\n")
	formattedShipments.CurrentShipmentStatuses = strings.Join(formattedShipmentStatuses, "\n\n")

	return formattedShipments
}

// // FetchMovingExpensesShipmentSummaryWorksheet fetches moving expenses for the Shipment Summary Worksheet
// // TODO: update to create moving expense summary with the new moving expense model
// func FetchMovingExpensesShipmentSummaryWorksheet(move models.Move, db *pop.Connection, session *auth.Session) (models.MovingExpenses, error) {
// 	var movingExpenseDocuments models.MovingExpenses

// 	return movingExpenseDocuments, nil
// }

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

// UPDATED - to use models.PPMShipment and includes logic for calculating net weight
// NOTE - we no longer have a ppmShipment.NetWeight field, instead it has been calculated
// based on the full - empty weights of approved weight tickets
// we don't yet have a service or reusable function for doing these calculations
// FormatPPMWeight formats a ppms NetWeight for the Shipment Summary Worksheet
func FormatPPMWeight(ppm models.PPMShipment) string {
	weightTickets := ppm.WeightTickets
	var fullWeights unit.Pound
	var emptyWeights unit.Pound
	for _, weightTicket := range weightTickets {
		emptyWeight := weightTicket.EmptyWeight
		fullWeight := weightTicket.FullWeight
		status := weightTicket.Status
		if status != nil && *status == models.PPMDocumentStatusApproved {
			if emptyWeight != nil {
				emptyWeights += *emptyWeight
			}
			if fullWeight != nil {
				fullWeights += *fullWeight
			}
		}
	}
	wtg := FormatWeights(unit.Pound(fullWeights - emptyWeights))
	return fmt.Sprintf("%s lbs - FINAL", wtg)
}

// UPDATED - to use models.PPMShipment
// maps the old personallyProcuredMove.OriginalMoveDate -> ppmShipment.ActualMoveDate
// FormatPPMPickupDate formats a shipments ActualPickupDate for the Shipment Summary Worksheet
func FormatPPMPickupDate(ppm models.PPMShipment) string {
	if ppm.ActualMoveDate != nil {
		return FormatDate(*ppm.ActualMoveDate)
	}
	return ""
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
