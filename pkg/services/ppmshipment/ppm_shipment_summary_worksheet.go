package ppmshipment

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/utilities"
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
	Obligations         Obligations
	MovingExpenses      models.MovingExpenses
	// PPMRemainingEntitlement unit.Pound
	SignedCertification models.SignedCertification
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

// FetchDataShipmentSummaryWorksheetFormData fetches the pages for the Shipment Summary Worksheet for a given Move ID
func (f *ppmShipmentSummaryWorksheetCreator) FetchDataShipmentSummaryWorksheetFormData(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (ShipmentSummaryFormData, error) {
	ppmShipment, err := FetchPPMShipmentFromMTOShipmentID(appCtx, mtoShipmentID)

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
		weightAllotment = models.SSWGetEntitlement(rank, orders.HasDependents, orders.SpouseHasProGear)
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

// SSWMaxWeightEntitlement weight allotment for the shipment summary worksheet.
type SSWMaxWeightEntitlement struct {
	Entitlement   unit.Pound
	ProGear       unit.Pound
	SpouseProGear unit.Pound
	TotalWeight   unit.Pound
}

// QUESTION: do we need this? do we have to use reflect here?
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
func SSWGetEntitlement(rank models.ServiceMemberRank, hasDependents bool, spouseHasProGear bool) SSWMaxWeightEntitlement {
	sswEntitlements := SSWMaxWeightEntitlement{}
	entitlements := models.GetWeightAllotment(rank)
	sswEntitlements.addLineItem("ProGear", entitlements.ProGearWeight)
	if !hasDependents {
		sswEntitlements.addLineItem("Entitlement", entitlements.TotalWeightSelf)
		return sswEntitlements
	}
	sswEntitlements.addLineItem("Entitlement", entitlements.TotalWeightSelfPlusDependents)
	if spouseHasProGear {
		sswEntitlements.addLineItem("SpouseProGear", entitlements.ProGearWeightSpouse)
	}
	return sswEntitlements
}

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
	page1.PPMRemainingEntitlement = FormatWeights(data.PPMRemainingEntitlement)
	page1.ActualObligationGCC95 = FormatDollars(actualObligations.GCC95())
	page1.ActualObligationSIT = FormatDollars(actualObligations.FormatSIT())
	page1.ActualObligationAdvance = formatActualObligationAdvance(data)
	page1.MileageTotal = actualObligations.Miles.String()
	return page1
}

func formatActualObligationAdvance(data ShipmentSummaryFormData) string {
	if len(data.PPMShipments) > 0 && data.PPMShipments[0].AdvanceAmountRequested != nil {
		advance := data.PPMShipments[0].AdvanceAmountRequested.ToDollarFloatNoRound()
		return FormatDollars(advance)
	}
	return FormatDollars(0)
}

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
func FormatAllShipments(ppms models.PersonallyProcuredMoves) ShipmentSummaryWorkSheetShipments {
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

// FetchMovingExpensesShipmentSummaryWorksheet fetches moving expenses for the Shipment Summary Worksheet
// TODO: update to create moving expense summary with the new moving expense model
func FetchMovingExpensesShipmentSummaryWorksheet(move models.Move, db *pop.Connection, session *auth.Session) (models.MovingExpenses, error) {
	var movingExpenseDocuments models.MovingExpenses

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
func FormatCurrentPPMStatus(ppm models.PersonallyProcuredMove) string {
	if ppm.Status == "PAYMENT_REQUESTED" {
		return "At destination"
	}
	return FormatEnum(string(ppm.Status), " ")
}

// FormatPPMNumberAndType formats FormatShipmentNumberAndType for the Shipment Summary Worksheet
func FormatPPMNumberAndType(i int) string {
	return fmt.Sprintf("%02d - PPM", i+1)
}

// FormatPPMWeight formats a ppms NetWeight for the Shipment Summary Worksheet
func FormatPPMWeight(ppm models.PersonallyProcuredMove) string {
	if ppm.NetWeight != nil {
		wtg := FormatWeights(unit.Pound(*ppm.NetWeight))
		return fmt.Sprintf("%s lbs - FINAL", wtg)
	}
	return ""
}

// FormatPPMPickupDate formats a shipments ActualPickupDate for the Shipment Summary Worksheet
func FormatPPMPickupDate(ppm models.PersonallyProcuredMove) string {
	if ppm.OriginalMoveDate != nil {
		return FormatDate(*ppm.OriginalMoveDate)
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

// GetPPMDocuments returns all documents associated with a PPM shipment.
func (f *ppmShipmentSummaryWorksheetCreator) GetPPMDocuments(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (*models.PPMDocuments, error) {
	var documents models.PPMDocuments

	err := appCtx.DB().
		Scope(utilities.ExcludeDeletedScope(models.WeightTicket{})).
		EagerPreload(
			"EmptyDocument.UserUploads.Upload",
			"FullDocument.UserUploads.Upload",
			"ProofOfTrailerOwnershipDocument.UserUploads.Upload",
		).
		InnerJoin("ppm_shipments ppm", "ppm.id = weight_tickets.ppm_shipment_id").
		Where("ppm.shipment_id = ? AND ppm.deleted_at IS NULL", mtoShipmentID).
		All(&documents.WeightTickets)

	if err != nil {
		return nil, apperror.NewQueryError("WeightTicket", err, "unable to search for WeightTickets")
	}

	for i := range documents.WeightTickets {
		documents.WeightTickets[i].EmptyDocument.UserUploads = documents.WeightTickets[i].EmptyDocument.UserUploads.FilterDeleted()
		documents.WeightTickets[i].FullDocument.UserUploads = documents.WeightTickets[i].FullDocument.UserUploads.FilterDeleted()
		documents.WeightTickets[i].ProofOfTrailerOwnershipDocument.UserUploads = documents.WeightTickets[i].ProofOfTrailerOwnershipDocument.UserUploads.FilterDeleted()
	}

	err = appCtx.DB().
		Scope(utilities.ExcludeDeletedScope(models.ProgearWeightTicket{})).
		EagerPreload(
			"Document.UserUploads.Upload",
		).
		InnerJoin("ppm_shipments ppm", "ppm.id = progear_weight_tickets.ppm_shipment_id").
		Where("ppm.shipment_id = ? AND ppm.deleted_at IS NULL", mtoShipmentID).
		All(&documents.ProgearWeightTickets)

	if err != nil {
		return nil, apperror.NewQueryError("ProgearWeightTicket", err, "unable to search for ProgearWeightTickets")
	}

	for i := range documents.ProgearWeightTickets {
		documents.ProgearWeightTickets[i].Document.UserUploads = documents.ProgearWeightTickets[i].Document.UserUploads.FilterDeleted()
	}

	err = appCtx.DB().
		Scope(utilities.ExcludeDeletedScope(models.MovingExpense{})).
		EagerPreload(
			"Document.UserUploads.Upload",
		).
		InnerJoin("ppm_shipments ppm", "ppm.id = moving_expenses.ppm_shipment_id").
		Where("ppm.shipment_id = ? AND ppm.deleted_at IS NULL", mtoShipmentID).
		All(&documents.MovingExpenses)

	if err != nil {
		return nil, apperror.NewQueryError("MovingExpense", err, "unable to search for MovingExpenses")
	}

	for i := range documents.MovingExpenses {
		documents.MovingExpenses[i].Document.UserUploads = documents.MovingExpenses[i].Document.UserUploads.FilterDeleted()
	}

	return &documents, nil
}
