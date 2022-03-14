package models

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/unit"
)

// FormatValuesShipmentSummaryWorksheet returns the formatted pages for the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheet(shipmentSummaryFormData ShipmentSummaryFormData) (ShipmentSummaryWorksheetPage1Values, ShipmentSummaryWorksheetPage2Values, ShipmentSummaryWorksheetPage3Values, error) {
	page1 := FormatValuesShipmentSummaryWorksheetFormPage1(shipmentSummaryFormData)
	page2, err := FormatValuesShipmentSummaryWorksheetFormPage2(shipmentSummaryFormData)
	page3 := FormatValuesShipmentSummaryWorksheetFormPage3(shipmentSummaryFormData)
	if err != nil {
		return page1, page2, page3, err
	}
	return page1, page2, page3, nil
}

// ShipmentSummaryWorksheetPage1Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage1Values struct {
	FOUOBanner                      string
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
	FOUOBanner      string
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

//FormattedMovingExpenses is an object representing the service member's moving expenses formatted for the SSW
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
	FOUOBanner             string
	PreparationDate        string
	ServiceMemberSignature string
	SignatureDate          string
	FormattedOtherExpenses
}

// ShipmentSummaryFormData is a container for the various objects required for the a Shipment Summary Worksheet
type ShipmentSummaryFormData struct {
	ServiceMember           ServiceMember
	Order                   Order
	CurrentDutyLocation     DutyLocation
	NewDutyLocation         DutyLocation
	WeightAllotment         SSWMaxWeightEntitlement
	PersonallyProcuredMoves PersonallyProcuredMoves
	PreparationDate         time.Time
	Obligations             Obligations
	MovingExpenseDocuments  []MovingExpenseDocument
	PPMRemainingEntitlement unit.Pound
	SignedCertification     SignedCertification
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
func FetchDataShipmentSummaryWorksheetFormData(db *pop.Connection, session *auth.Session, moveID uuid.UUID) (ShipmentSummaryFormData, error) {
	move := Move{}
	dbQErr := db.Q().Eager(
		"Orders",
		"Orders.NewDutyLocation.Address",
		"Orders.ServiceMember",
		"Orders.ServiceMember.DutyLocation.Address",
		"PersonallyProcuredMoves",
	).Find(&move, moveID)

	if dbQErr != nil {
		if errors.Cause(dbQErr).Error() == RecordNotFoundErrorString {
			return ShipmentSummaryFormData{}, ErrFetchNotFound
		}
		return ShipmentSummaryFormData{}, dbQErr
	}

	for i, ppm := range move.PersonallyProcuredMoves {
		ppmDetails, err := FetchPersonallyProcuredMove(db, session, ppm.ID)
		if err != nil {
			return ShipmentSummaryFormData{}, err
		}
		if ppmDetails.Advance != nil {
			status := ppmDetails.Advance.Status
			if status == ReimbursementStatusAPPROVED || status == ReimbursementStatusPAID {
				move.PersonallyProcuredMoves[i].Advance = ppmDetails.Advance
			}
		}
	}

	_, authErr := FetchOrderForUser(db, session, move.OrdersID)
	if authErr != nil {
		return ShipmentSummaryFormData{}, authErr
	}

	serviceMember := move.Orders.ServiceMember
	var rank ServiceMemberRank
	var weightAllotment SSWMaxWeightEntitlement
	if serviceMember.Rank != nil {
		rank = ServiceMemberRank(*serviceMember.Rank)
		weightAllotment = SSWGetEntitlement(rank, move.Orders.HasDependents, move.Orders.SpouseHasProGear)
	}

	movingExpenses, err := FetchMovingExpensesShipmentSummaryWorksheet(move, db, session)
	if err != nil {
		return ShipmentSummaryFormData{}, err
	}

	ppmRemainingEntitlement, err := CalculateRemainingPPMEntitlement(move, weightAllotment.TotalWeight)
	if err != nil {
		return ShipmentSummaryFormData{}, err
	}

	signedCertification, err := FetchSignedCertificationsPPMPayment(db, session, moveID)
	if err != nil {
		return ShipmentSummaryFormData{}, err
	}
	if signedCertification == nil {
		return ShipmentSummaryFormData{},
			errors.New("shipment summary worksheet: signed certification is nil")
	}
	ssd := ShipmentSummaryFormData{
		ServiceMember:           serviceMember,
		Order:                   move.Orders,
		CurrentDutyLocation:     serviceMember.DutyLocation,
		NewDutyLocation:         move.Orders.NewDutyLocation,
		WeightAllotment:         weightAllotment,
		PersonallyProcuredMoves: move.PersonallyProcuredMoves,
		SignedCertification:     *signedCertification,
		PPMRemainingEntitlement: ppmRemainingEntitlement,
		MovingExpenseDocuments:  movingExpenses,
	}
	return ssd, nil
}

//SSWMaxWeightEntitlement weight allotment for the shipment summary worksheet.
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
func SSWGetEntitlement(rank ServiceMemberRank, hasDependents bool, spouseHasProGear bool) SSWMaxWeightEntitlement {
	sswEntitlements := SSWMaxWeightEntitlement{}
	entitlements := GetWeightAllotment(rank)
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
func CalculateRemainingPPMEntitlement(move Move, totalEntitlement unit.Pound) (unit.Pound, error) {
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

// FetchMovingExpensesShipmentSummaryWorksheet fetches moving expenses for the Shipment Summary Worksheet
func FetchMovingExpensesShipmentSummaryWorksheet(move Move, db *pop.Connection, session *auth.Session) ([]MovingExpenseDocument, error) {
	var movingExpenseDocuments []MovingExpenseDocument
	if len(move.PersonallyProcuredMoves) > 0 {
		ppm := move.PersonallyProcuredMoves[0]
		status := MoveDocumentStatusOK
		moveDocuments, err := FetchMoveDocuments(db, session, ppm.ID, &status, MoveDocumentTypeEXPENSE, false)
		if err != nil {
			return movingExpenseDocuments, err
		}
		movingExpenseDocuments = FilterMovingExpenseDocuments(moveDocuments)
	}
	return movingExpenseDocuments, nil
}

const (
	forOfficialUseOnlyText = "UNCLASSIFIED // FOR OFFICIAL USE ONLY"
)

// FormatValuesShipmentSummaryWorksheetFormPage1 formats the data for page 1 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage1(data ShipmentSummaryFormData) ShipmentSummaryWorksheetPage1Values {
	page1 := ShipmentSummaryWorksheetPage1Values{}
	page1.FOUOBanner = forOfficialUseOnlyText
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

	formattedShipments := FormatAllShipments(data.PersonallyProcuredMoves)
	page1.ShipmentNumberAndTypes = formattedShipments.ShipmentNumberAndTypes
	page1.ShipmentPickUpDates = formattedShipments.PickUpDates
	page1.ShipmentCurrentShipmentStatuses = formattedShipments.CurrentShipmentStatuses
	page1.ShipmentWeights = formattedShipments.ShipmentWeights

	formattedSit := FormatAllSITExpenses(data.MovingExpenseDocuments)
	page1.SITNumberAndTypes = formattedSit.NumberAndTypes
	page1.SITEntryDates = formattedSit.EntryDates
	page1.SITEndDates = formattedSit.EndDates
	page1.SITDaysInStorage = formattedSit.DaysInStorage

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
	if len(data.PersonallyProcuredMoves) > 0 && data.PersonallyProcuredMoves[0].Advance != nil {
		advance := data.PersonallyProcuredMoves[0].Advance.RequestedAmount.ToDollarFloatNoRound()
		return FormatDollars(advance)
	}
	return FormatDollars(0)
}

//FormatRank formats the service member's rank for Shipment Summary Worksheet
func FormatRank(rank *ServiceMemberRank) string {
	var rankDisplayValue = map[ServiceMemberRank]string{
		ServiceMemberRankE1:                      "E-1",
		ServiceMemberRankE2:                      "E-2",
		ServiceMemberRankE3:                      "E-3",
		ServiceMemberRankE4:                      "E-4",
		ServiceMemberRankE5:                      "E-5",
		ServiceMemberRankE6:                      "E-6",
		ServiceMemberRankE7:                      "E-7",
		ServiceMemberRankE8:                      "E-8",
		ServiceMemberRankE9:                      "E-9",
		ServiceMemberRankE9SPECIALSENIORENLISTED: "E-9 (Special Senior Enlisted)",
		ServiceMemberRankO1ACADEMYGRADUATE:       "O-1 or Service Academy Graduate",
		ServiceMemberRankO2:                      "O-2",
		ServiceMemberRankO3:                      "O-3",
		ServiceMemberRankO4:                      "O-4",
		ServiceMemberRankO5:                      "O-5",
		ServiceMemberRankO6:                      "O-6",
		ServiceMemberRankO7:                      "O-7",
		ServiceMemberRankO8:                      "O-8",
		ServiceMemberRankO9:                      "O-9",
		ServiceMemberRankO10:                     "O-10",
		ServiceMemberRankW1:                      "W-1",
		ServiceMemberRankW2:                      "W-2",
		ServiceMemberRankW3:                      "W-3",
		ServiceMemberRankW4:                      "W-4",
		ServiceMemberRankW5:                      "W-5",
		ServiceMemberRankAVIATIONCADET:           "Aviation Cadet",
		ServiceMemberRankCIVILIANEMPLOYEE:        "Civilian Employee",
		ServiceMemberRankACADEMYCADET:            "Service Academy Cadet",
		ServiceMemberRankMIDSHIPMAN:              "Midshipman",
	}
	if rank != nil {
		return rankDisplayValue[*rank]
	}
	return ""
}

// FormatValuesShipmentSummaryWorksheetFormPage2 formats the data for page 2 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage2(data ShipmentSummaryFormData) (ShipmentSummaryWorksheetPage2Values, error) {
	var err error
	page2 := ShipmentSummaryWorksheetPage2Values{}
	page2.FOUOBanner = forOfficialUseOnlyText
	page2.TAC = derefStringTypes(data.Order.TAC)
	page2.SAC = derefStringTypes(data.Order.SAC)
	page2.PreparationDate = FormatDate(data.PreparationDate)
	page2.FormattedMovingExpenses, err = FormatMovingExpenses(data.MovingExpenseDocuments)
	if err != nil {
		return page2, err
	}
	page2.TotalMemberPaidRepeated = page2.TotalMemberPaid
	page2.TotalGTCCPaidRepeated = page2.TotalGTCCPaid
	return page2, nil
}

// FormatValuesShipmentSummaryWorksheetFormPage3 formats the data for page 2 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage3(data ShipmentSummaryFormData) ShipmentSummaryWorksheetPage3Values {
	page3 := ShipmentSummaryWorksheetPage3Values{}
	page3.FOUOBanner = forOfficialUseOnlyText
	page3.PreparationDate = FormatDate(data.PreparationDate)
	page3.FormattedOtherExpenses = FormatOtherExpenses(data.MovingExpenseDocuments)
	page3.ServiceMemberSignature = FormatSignature(data.ServiceMember)
	page3.SignatureDate = FormatSignatureDate(data.SignedCertification)
	return page3
}

// FormatOtherExpenses formats other expenses
func FormatOtherExpenses(docs MovingExpenseDocuments) FormattedOtherExpenses {
	var expenseDescriptions []string
	var expenseAmounts []string
	for _, doc := range docs {
		if doc.MovingExpenseType == MovingExpenseTypeOTHER {
			expenseDescriptions = append(expenseDescriptions, doc.MoveDocument.Title)
			expenseAmounts = append(expenseAmounts, FormatDollars(float64(doc.RequestedAmountCents.ToDollarFloatNoRound())))
		}
	}
	return FormattedOtherExpenses{
		Descriptions: strings.Join(expenseDescriptions, "\n\n"),
		AmountsPaid:  strings.Join(expenseAmounts, "\n\n"),
	}
}

//FormatSignature formats a service member's signature for the Shipment Summary Worksheet
func FormatSignature(sm ServiceMember) string {
	first := derefStringTypes(sm.FirstName)
	last := derefStringTypes(sm.LastName)

	return fmt.Sprintf("%s %s electronically signed", first, last)
}

// FormatSignatureDate formats the date the service member electronically signed for the Shipment Summary Worksheet
func FormatSignatureDate(signature SignedCertification) string {
	dateLayout := "02 Jan 2006 at 3:04pm"
	dt := signature.Date.Format(dateLayout)
	return dt
}

//FormatLocation formats AuthorizedOrigin and AuthorizedDestination for Shipment Summary Worksheet
func FormatLocation(dutyLocation DutyLocation) string {
	return fmt.Sprintf("%s, %s %s", dutyLocation.Name, dutyLocation.Address.State, dutyLocation.Address.PostalCode)
}

//FormatServiceMemberFullName formats ServiceMember full name for Shipment Summary Worksheet
func FormatServiceMemberFullName(serviceMember ServiceMember) string {
	lastName := derefStringTypes(serviceMember.LastName)
	suffix := derefStringTypes(serviceMember.Suffix)
	firstName := derefStringTypes(serviceMember.FirstName)
	middleName := derefStringTypes(serviceMember.MiddleName)
	if suffix != "" {
		return fmt.Sprintf("%s %s, %s %s", lastName, suffix, firstName, middleName)
	}
	return strings.TrimSpace(fmt.Sprintf("%s, %s %s", lastName, firstName, middleName))
}

//FormatAllShipments formats Shipment line items for the Shipment Summary Worksheet
func FormatAllShipments(ppms PersonallyProcuredMoves) ShipmentSummaryWorkSheetShipments {
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

//FormatAllSITExpenses formats SIT line items for the Shipment Summary Worksheet
func FormatAllSITExpenses(movingExpenseDocuments MovingExpenseDocuments) ShipmentSummaryWorkSheetSIT {
	formattedShipments := ShipmentSummaryWorkSheetSIT{}
	sitExpenses := FilterSITExpenses(movingExpenseDocuments)
	totalSITExpenses := len(sitExpenses)
	formattedShipmentNumberAndTypes := make([]string, totalSITExpenses)
	formattedEntryDates := make([]string, totalSITExpenses)
	formattedEndDates := make([]string, totalSITExpenses)
	formattedDaysInStorage := make([]string, totalSITExpenses)

	for i, sitExpense := range sitExpenses {
		formattedShipmentNumberAndTypes[i] = FormatPPMNumberAndType(i)
		if sitExpense.StorageStartDate != nil {
			formattedEntryDates[i] = FormatDate(*sitExpense.StorageStartDate)
		}
		if sitExpense.StorageEndDate != nil {
			formattedEndDates[i] = FormatDate(*sitExpense.StorageEndDate)
		}
		days, err := sitExpense.DaysInStorage()
		if err == nil {
			formattedDaysInStorage[i] = fmt.Sprintf("%d", days)
		}
	}

	formattedShipments.NumberAndTypes = strings.Join(formattedShipmentNumberAndTypes, "\n\n")
	formattedShipments.EntryDates = strings.Join(formattedEntryDates, "\n\n")
	formattedShipments.EndDates = strings.Join(formattedEndDates, "\n\n")
	formattedShipments.DaysInStorage = strings.Join(formattedDaysInStorage, "\n\n")

	return formattedShipments
}

//FormatMovingExpenses formats moving expenses for Shipment Summary Worksheet
func FormatMovingExpenses(movingExpenseDocuments MovingExpenseDocuments) (FormattedMovingExpenses, error) {
	return SubTotalsMapToStruct(SubTotalExpenses(movingExpenseDocuments))
}

//SubTotalExpenses groups moving expenses by type and payment method
func SubTotalExpenses(expenseDocuments MovingExpenseDocuments) map[string]float64 {
	var expenseType string
	totals := make(map[string]float64)
	for _, expense := range expenseDocuments {
		expenseType = getExpenseType(expense)
		expenseDollarAmt := expense.RequestedAmountCents.ToDollarFloatNoRound()
		totals[expenseType] += expenseDollarAmt
		addToGrandTotal(totals, expenseType, expenseDollarAmt)
	}
	return totals
}

// SubTotalsMapToStruct takes subtotal map and returns struct
func SubTotalsMapToStruct(subTotals map[string]float64) (FormattedMovingExpenses, error) {
	expenses := FormattedMovingExpenses{}
	err := mapstructure.Decode(subTotals, &expenses)
	if err != nil {
		return FormattedMovingExpenses{}, err
	}
	return expenses, nil
}

func addToGrandTotal(totals map[string]float64, key string, expenseDollarAmt float64) {
	if strings.HasPrefix(key, "Storage") {
		if strings.HasSuffix(key, "GTCCPaid") {
			totals["TotalGTCCPaidSIT"] += expenseDollarAmt
		} else {
			totals["TotalMemberPaidSIT"] += expenseDollarAmt
		}
		totals["TotalPaidSIT"] += expenseDollarAmt
	} else {
		if strings.HasSuffix(key, "GTCCPaid") {
			totals["TotalGTCCPaid"] += expenseDollarAmt
		} else {
			totals["TotalMemberPaid"] += expenseDollarAmt
		}
		totals["TotalPaidNonSIT"] += expenseDollarAmt
	}
}

func getExpenseType(expense MovingExpenseDocument) string {
	expenseType := FormatEnum(string(expense.MovingExpenseType), "")
	if expense.PaymentMethod == "GTCC" {
		return fmt.Sprintf("%s%s", expenseType, "GTCCPaid")
	}
	return fmt.Sprintf("%s%s", expenseType, "MemberPaid")
}

//FormatCurrentPPMStatus formats FormatCurrentPPMStatus for the Shipment Summary Worksheet
func FormatCurrentPPMStatus(ppm PersonallyProcuredMove) string {
	if ppm.Status == "PAYMENT_REQUESTED" {
		return "At destination"
	}
	return FormatEnum(string(ppm.Status), " ")
}

//FormatPPMNumberAndType formats FormatShipmentNumberAndType for the Shipment Summary Worksheet
func FormatPPMNumberAndType(i int) string {
	return fmt.Sprintf("%02d - PPM", i+1)
}

//FormatPPMWeight formats a ppms NetWeight for the Shipment Summary Worksheet
func FormatPPMWeight(ppm PersonallyProcuredMove) string {
	if ppm.NetWeight != nil {
		wtg := FormatWeights(unit.Pound(*ppm.NetWeight))
		return fmt.Sprintf("%s lbs - FINAL", wtg)
	}
	return ""
}

//FormatPPMPickupDate formats a shipments ActualPickupDate for the Shipment Summary Worksheet
func FormatPPMPickupDate(ppm PersonallyProcuredMove) string {
	if ppm.OriginalMoveDate != nil {
		return FormatDate(*ppm.OriginalMoveDate)
	}
	return ""
}

//FormatOrdersTypeAndOrdersNumber formats OrdersTypeAndOrdersNumber for Shipment Summary Worksheet
func FormatOrdersTypeAndOrdersNumber(order Order) string {
	issuingBranch := FormatOrdersType(order)
	ordersNumber := derefStringTypes(order.OrdersNumber)
	return fmt.Sprintf("%s/%s", issuingBranch, ordersNumber)
}

//FormatServiceMemberAffiliation formats ServiceMemberAffiliation in human friendly format
func FormatServiceMemberAffiliation(affiliation *ServiceMemberAffiliation) string {
	if affiliation != nil {
		return FormatEnum(string(*affiliation), " ")
	}
	return ""
}

//FormatOrdersType formats OrdersType for Shipment Summary Worksheet
func FormatOrdersType(order Order) string {
	switch order.OrdersType {
	case internalmessages.OrdersTypePERMANENTCHANGEOFSTATION:
		return "PCS"
	default:
		return ""
	}
}

//FormatDate formats Dates for Shipment Summary Worksheet
func FormatDate(date time.Time) string {
	dateLayout := "02-Jan-2006"
	return date.Format(dateLayout)
}

//FormatEnum titlecases string const types (e.g. THIS_CONSTANT -> This Constant)
//outSep specifies the character to use for rejoining the string
func FormatEnum(s string, outSep string) string {
	words := strings.Replace(strings.ToLower(s), "_", " ", -1)
	return strings.Replace(strings.Title(words), " ", outSep, -1)
}

//FormatWeights formats a unit.Pound using 000s separator
func FormatWeights(wtg unit.Pound) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", wtg)
}

//FormatDollars formats an int using 000s separator
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
