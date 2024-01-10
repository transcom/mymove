package shipmentsummaryworksheet

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

// FormatValuesShipmentSummaryWorksheet returns the formatted pages for the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheet(shipmentSummaryFormData ShipmentSummaryFormData) (Page1Values, Page2Values, Page3Values, error) {
	page1 := FormatValuesShipmentSummaryWorksheetFormPage1(shipmentSummaryFormData)
	page2 := FormatValuesShipmentSummaryWorksheetFormPage2(shipmentSummaryFormData)
	page3 := FormatValuesShipmentSummaryWorksheetFormPage3(shipmentSummaryFormData)

	return page1, page2, page3, nil
}

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

// FetchDataShipmentSummaryWorksheetFormData fetches the pages for the Shipment Summary Worksheet for a given Move ID
func FetchDataShipmentSummaryWorksheetFormData(appCtx appcontext.AppContext, session *auth.Session, ppmShipmentID uuid.UUID) (ShipmentSummaryFormData, error) {
	ppmShipment := models.PPMShipment{}
	dbQErr := appCtx.DB().Q().Eager(
		"Shipment.MoveTaskOrder.Orders.ServiceMember",
		"Shipment.MoveTaskOrder",
		"Shipment.MoveTaskOrder.Orders",
		"Shipment.MoveTaskOrder.Orders.NewDutyLocation.Address",
		"Shipment.MoveTaskOrder.Orders.ServiceMember.DutyLocation.Address",
	).Find(&ppmShipment, ppmShipmentID)

	if dbQErr != nil {
		if errors.Cause(dbQErr).Error() == models.RecordNotFoundErrorString {
			return ShipmentSummaryFormData{}, models.ErrFetchNotFound
		}
		return ShipmentSummaryFormData{}, dbQErr
	}

	// for i, ppm := range move.PersonallyProcuredMoves {
	// 	ppmDetails, err := models.FetchPersonallyProcuredMove(appCtx.DB(), session, ppm.ID)
	// 	if err != nil {
	// 		return ShipmentSummaryFormData{}, err
	// 	}
	// 	if ppmDetails.Advance != nil {
	// 		status := ppmDetails.Advance.Status
	// 		if status == models.ReimbursementStatusAPPROVED || status == models.ReimbursementStatusPAID {
	// 			move.PersonallyProcuredMoves[i].Advance = ppmDetails.Advance
	// 		}
	// 	}
	// }

	// _, authErr := models.FetchOrderForUser(appCtx.DB(), session, move.OrdersID)
	// if authErr != nil {
	// 	return ShipmentSummaryFormData{}, authErr
	// }

	serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
	var rank models.ServiceMemberRank
	var weightAllotment SSWMaxWeightEntitlement
	if serviceMember.Rank != nil {
		rank = models.ServiceMemberRank(*serviceMember.Rank)
		weightAllotment = SSWGetEntitlement(rank, ppmShipment.Shipment.MoveTaskOrder.Orders.HasDependents, ppmShipment.Shipment.MoveTaskOrder.Orders.SpouseHasProGear)
	}

	ppmRemainingEntitlement, err := CalculateRemainingPPMEntitlement(ppmShipment.Shipment.MoveTaskOrder, weightAllotment.TotalWeight)
	if err != nil {
		return ShipmentSummaryFormData{}, err
	}

	signedCertification, err := models.FetchSignedCertificationsPPMPayment(appCtx.DB(), session, ppmShipment.Shipment.MoveTaskOrderID)
	if err != nil {
		return ShipmentSummaryFormData{}, err
	}
	if signedCertification == nil {
		return ShipmentSummaryFormData{},
			errors.New("shipment summary worksheet: signed certification is nil")
	}

	moveHolder := models.Move{}
	var ppmShipments []models.PPMShipment

	// MTOShipments is inherently plural
	for _, mtoShipment := range moveHolder.MTOShipments {
		if mtoShipment.PPMShipment != nil {
			// We have a PPM shipment present, append it
			ppmShipments = append(ppmShipments, *mtoShipment.PPMShipment)
		}
	}

	ssd := ShipmentSummaryFormData{
		ServiceMember:           serviceMember,
		Order:                   ppmShipment.Shipment.MoveTaskOrder.Orders,
		Move:                    ppmShipment.Shipment.MoveTaskOrder,
		CurrentDutyLocation:     serviceMember.DutyLocation,
		NewDutyLocation:         ppmShipment.Shipment.MoveTaskOrder.Orders.NewDutyLocation,
		WeightAllotment:         weightAllotment,
		PPMShipments:            ppmShipments,
		SignedCertification:     *signedCertification,
		PPMRemainingEntitlement: ppmRemainingEntitlement,
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

// CalculateRemainingPPMEntitlement calculates the remaining PPM entitlement for PPM moves
// a PPMs remaining entitlement weight is equal to total entitlement - hhg weight
func CalculateRemainingPPMEntitlement(move models.Move, totalEntitlement unit.Pound) (unit.Pound, error) {
	var hhgActualWeight unit.Pound

	var ppmActualWeight unit.Pound
	if len(move.PersonallyProcuredMoves) > 0 {
		if move.PersonallyProcuredMoves[0].NetWeight == nil {
			return ppmActualWeight, errors.Errorf("PPM %s does not have NetWeight", move.PersonallyProcuredMoves[0].ID)
		}
		ppmActualWeight = unit.Pound(*move.PersonallyProcuredMoves[0].NetWeight)
	}

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
func FormatValuesShipmentSummaryWorksheetFormPage1(data ShipmentSummaryFormData) Page1Values {
	page1 := Page1Values{}
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
	if len(data.PPMShipments) > 0 && data.PPMShipments[0].HasRequestedAdvance != nil {
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
func FormatValuesShipmentSummaryWorksheetFormPage2(data ShipmentSummaryFormData) Page2Values {
	page2 := Page2Values{}
	page2.CUIBanner = controlledUnclassifiedInformationText
	page2.TAC = derefStringTypes(data.Order.TAC)
	page2.SAC = derefStringTypes(data.Order.SAC)
	page2.PreparationDate = FormatDate(data.PreparationDate)
	page2.TotalMemberPaidRepeated = page2.TotalMemberPaid
	page2.TotalGTCCPaidRepeated = page2.TotalGTCCPaid
	return page2
}

// FormatValuesShipmentSummaryWorksheetFormPage3 formats the data for page 2 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage3(data ShipmentSummaryFormData) Page3Values {
	page3 := Page3Values{}
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

	formattedShipments.ShipmentNumberAndTypes = strings.Join(formattedNumberAndTypes, "\n\n")
	formattedShipments.PickUpDates = strings.Join(formattedPickUpDates, "\n\n")
	formattedShipments.ShipmentWeights = strings.Join(formattedShipmentWeights, "\n\n")
	formattedShipments.CurrentShipmentStatuses = strings.Join(formattedShipmentStatuses, "\n\n")

	return formattedShipments
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
	return FormatDate(*&ppm.ExpectedDepartureDate)
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

// type ppmComputer interface {
// 	ComputePPMMoveCosts(appCtx appcontext.AppContext, weight unit.Pound, originPickupZip5 string, originDutyLocationZip5 string, destinationZip5 string, distanceMilesFromOriginPickupZip int, distanceMilesFromOriginDutyLocationZip int, date time.Time, daysInSit int) (cost rateengine.CostDetails, err error)
// }

// SSWPPMComputer a rate engine wrapper with helper functions to simplify ppm cost calculations specific to shipment summary worksheet
type SSWPPMComputer struct {
	ppmComputer
}

// NewSSWPPMComputer creates a SSWPPMComputer
func NewSSWPPMComputer(PPMShipment models.PPMShipment) *SSWPPMComputer {
	return &SSWPPMComputer{ppmComputer: PPMComputer}
}

// ObligationType type corresponding to obligation sections of shipment summary worksheet
type ObligationType int

// ComputeObligations is helper function for computing the obligations section of the shipment summary worksheet
func (sswPpmComputer *SSWPPMComputer) ComputeObligations(appCtx appcontext.AppContext, ssfd ShipmentSummaryFormData, planner route.Planner) (obligation Obligations, err error) {
	firstPPM, err := sswPpmComputer.nilCheckPPM(ssfd)
	if err != nil {
		return Obligations{}, err
	}

	originDutyLocationZip := ssfd.CurrentDutyLocation.Address.PostalCode
	destDutyLocationZip := ssfd.Order.NewDutyLocation.Address.PostalCode

	distanceMilesFromPickupZip, err := planner.ZipTransitDistance(appCtx, firstPPM.PickupPostalCode, destDutyLocationZip)
	if err != nil {
		return Obligations{}, errors.New("error calculating distance")
	}

	distanceMilesFromDutyLocationZip, err := planner.ZipTransitDistance(appCtx, originDutyLocationZip, destDutyLocationZip)
	if err != nil {
		return Obligations{}, errors.New("error calculating distance")
	}

	actualCosts, err := sswPpmComputer.ComputePPMMoveCosts(
		appCtx,
		ssfd.PPMRemainingEntitlement,
		firstPPM.PickupPostalCode,
		originDutyLocationZip,
		destDutyLocationZip,
		distanceMilesFromPickupZip,
		distanceMilesFromDutyLocationZip,
		firstPPM.ExpectedDepartureDate,
		0,
	)
	if err != nil {
		return Obligations{}, errors.New("error calculating PPM actual obligations")
	}

	maxCosts, err := sswPpmComputer.ComputePPMMoveCosts(
		appCtx,
		ssfd.WeightAllotment.TotalWeight,
		firstPPM.PickupPostalCode,
		originDutyLocationZip,
		destDutyLocationZip,
		distanceMilesFromPickupZip,
		distanceMilesFromDutyLocationZip,
		firstPPM.ExpectedDepartureDate,
		0,
	)
	if err != nil {
		return Obligations{}, errors.New("error calculating PPM max obligations")
	}

	actualCost := rateengine.GetWinningCostMove(actualCosts)
	maxCost := rateengine.GetWinningCostMove(maxCosts)
	nonWinningActualCost := rateengine.GetNonWinningCostMove(actualCosts)
	nonWinningMaxCost := rateengine.GetNonWinningCostMove(maxCosts)

	var actualSIT unit.Cents
	if firstPPM.SITEstimatedCost != nil {
		actualSIT = *firstPPM.SITEstimatedCost
	}

	if actualSIT > maxCost.SITMax {
		actualSIT = maxCost.SITMax
	}

	obligations := Obligations{
		ActualObligation:           Obligation{Gcc: actualCost.GCC, SIT: actualSIT, Miles: unit.Miles(actualCost.Mileage)},
		MaxObligation:              Obligation{Gcc: maxCost.GCC, SIT: actualSIT, Miles: unit.Miles(actualCost.Mileage)},
		NonWinningActualObligation: Obligation{Gcc: nonWinningActualCost.GCC, SIT: actualSIT, Miles: unit.Miles(nonWinningActualCost.Mileage)},
		NonWinningMaxObligation:    Obligation{Gcc: nonWinningMaxCost.GCC, SIT: actualSIT, Miles: unit.Miles(nonWinningActualCost.Mileage)},
	}
	return obligations, nil
}

func (sswPpmComputer *SSWPPMComputer) nilCheckPPM(ssfd ShipmentSummaryFormData) (models.PPMShipment, error) {
	if len(ssfd.PPMShipments) == 0 {
		return models.PPMShipment{}, errors.New("missing ppm")
	}
	firstPPM := ssfd.PPMShipments[0]
	if firstPPM.PickupPostalCode == "" || firstPPM.DestinationPostalCode == "" {
		return models.PPMShipment{}, errors.New("missing required address parameter")
	}
	// This test is being removed as they are checking whether required values exist
	// if firstPPM.ExpectedDepartureDate == nil {
	// 	return models.PersonallyProcuredMove{}, errors.New("missing required original move date parameter")
	// }
	return firstPPM, nil
}
