package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/unit"
)

// FormatValuesShipmentSummaryWorksheet returns the formatted pages for the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheet(shipmentSummaryFormData ShipmentSummaryFormData) (ShipmentSummaryWorksheetPage1Values, ShipmentSummaryWorksheetPage2Values, error) {
	page1 := FormatValuesShipmentSummaryWorksheetFormPage1(shipmentSummaryFormData)
	page2, err := FormatValuesShipmentSummaryWorksheetFormPage2(shipmentSummaryFormData)
	if err != nil {
		return page1, page2, err
	}
	return page1, page2, nil
}

// ShipmentSummaryWorksheetPage1Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage1Values struct {
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
	TAC                             string
	SAC                             string
	ShipmentNumberAndTypes          string
	ShipmentPickUpDates             string
	ShipmentWeights                 string
	ShipmentCurrentShipmentStatuses string
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
}

//ShipmentSummaryWorkSheetShipments is and object representing shipment line items on Shipment Summary Worksheet
type ShipmentSummaryWorkSheetShipments struct {
	ShipmentNumberAndTypes  string
	PickUpDates             string
	ShipmentWeights         string
	CurrentShipmentStatuses string
}

// ShipmentSummaryWorksheetPage2Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage2Values struct {
	PreparationDate string
	FormattedMovingExpenses
}

//FormattedMovingExpenses is an object representing the service member's moving expenses formatted for the SSW
type FormattedMovingExpenses struct {
	ContractedExpenseMemberPaid string
	ContractedExpenseGTCCPaid   string
	RentalEquipmentMemberPaid   string
	RentalEquipmentGTCCPaid     string
	PackingMaterialsMemberPaid  string
	PackingMaterialsGTCCPaid    string
	WeighingFeesMemberPaid      string
	WeighingFeesGTCCPaid        string
	GasMemberPaid               string
	GasGTCCPaid                 string
	TollsMemberPaid             string
	TollsGTCCPaid               string
	OilMemberPaid               string
	OilGTCCPaid                 string
	OtherMemberPaid             string
	OtherGTCCPaid               string
	TotalMemberPaid             string
	TotalGTCCPaid               string
	TotalMemberPaidRepeated     string
	TotalGTCCPaidRepeated       string
}

// ShipmentSummaryFormData is a container for the various objects required for the a Shipment Summary Worksheet
type ShipmentSummaryFormData struct {
	ServiceMember           ServiceMember
	Order                   Order
	CurrentDutyStation      DutyStation
	NewDutyStation          DutyStation
	WeightAllotment         WeightAllotment
	TotalWeightAllotment    unit.Pound
	Shipments               Shipments
	PersonallyProcuredMoves PersonallyProcuredMoves
	PreparationDate         time.Time
	Obligations             Obligations
	MovingExpenseDocuments  []MovingExpenseDocument
	PPMRemainingEntitlement unit.Pound
}

//Obligations an object representing the Max Obligation and Actual Obligation sections of the shipment summary worksheet
type Obligations struct {
	MaxObligation    Obligation
	ActualObligation Obligation
}

//Obligation an object representing the obligations section on the shipment summary worksheet
type Obligation struct {
	Gcc unit.Cents
	SIT unit.Cents
}

//GCC100 calculates the 100% GCC on shipment summary worksheet
func (obligation Obligation) GCC100() float64 {
	return obligation.Gcc.ToDollarFloat()
}

//GCC95 calculates the 95% GCC on shipment summary worksheet
func (obligation Obligation) GCC95() float64 {
	return obligation.Gcc.MultiplyFloat64(.95).ToDollarFloat()
}

// FormatSIT formats the SIT Cost into a dollar float for the shipment summary worksheet
func (obligation Obligation) FormatSIT() float64 {
	return obligation.SIT.ToDollarFloat()
}

//MaxAdvance calculates the Max Advance on the shipment summary worksheet
func (obligation Obligation) MaxAdvance() float64 {
	return obligation.Gcc.MultiplyFloat64(.60).ToDollarFloat()
}

// FetchDataShipmentSummaryWorksheetFormData fetches the pages for the Shipment Summary Worksheet for a given Move ID
func FetchDataShipmentSummaryWorksheetFormData(db *pop.Connection, session *auth.Session, moveID uuid.UUID) (data ShipmentSummaryFormData, err error) {
	move := Move{}
	err = db.Q().Eager(
		"Orders",
		"Orders.NewDutyStation.Address",
		"Orders.ServiceMember",
		"Orders.ServiceMember.DutyStation.Address",
		"Shipments",
		"PersonallyProcuredMoves",
	).Find(&move, moveID)

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

	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return ShipmentSummaryFormData{}, ErrFetchNotFound
		}
		return ShipmentSummaryFormData{}, err
	}

	_, authErr := FetchOrderForUser(db, session, move.OrdersID)
	if authErr != nil {
		return ShipmentSummaryFormData{}, authErr
	}

	serviceMember := move.Orders.ServiceMember
	var rank ServiceMemberRank
	var weightAllotment WeightAllotment
	var totalEntitlement unit.Pound
	if serviceMember.Rank != nil {
		rank = ServiceMemberRank(*serviceMember.Rank)
		weightAllotment = GetWeightAllotment(rank)
		te, err := GetEntitlement(rank, move.Orders.HasDependents, move.Orders.SpouseHasProGear)
		totalEntitlement = unit.Pound(te)
		if err != nil {
			return ShipmentSummaryFormData{}, err
		}
	}

	movingExpenses, err := FetchMovingExpensesShipmentSummaryWorksheet(move, db, session)
	if err != nil {
		return ShipmentSummaryFormData{}, err
	}

	ppmRemainingEntitlement := CalculateRemainingPPMEntitlement(move, totalEntitlement)

	ssd := ShipmentSummaryFormData{
		ServiceMember:           serviceMember,
		Order:                   move.Orders,
		CurrentDutyStation:      serviceMember.DutyStation,
		NewDutyStation:          move.Orders.NewDutyStation,
		WeightAllotment:         weightAllotment,
		TotalWeightAllotment:    totalEntitlement,
		Shipments:               move.Shipments,
		PersonallyProcuredMoves: move.PersonallyProcuredMoves,
		PPMRemainingEntitlement: ppmRemainingEntitlement,
		MovingExpenseDocuments:  movingExpenses,
	}
	return ssd, nil
}

// CalculateRemainingPPMEntitlement calculates the remaining PPM entitlement for PPM moves
// a PPMs remaining entitlement weight is equal to total entitlement - hhg weight
func CalculateRemainingPPMEntitlement(move Move, totalEntitlement unit.Pound) unit.Pound {
	var hhgActualWeight unit.Pound
	if len(move.Shipments) > 0 && move.Shipments[0].NetWeight != nil {
		hhgActualWeight = *move.Shipments[0].NetWeight
	}
	var ppmActualWeight unit.Pound
	if len(move.PersonallyProcuredMoves) > 0 {
		ppmActualWeight = unit.Pound(*move.PersonallyProcuredMoves[0].NetWeight)
	}
	switch ppmRemainingEntitlement := totalEntitlement - hhgActualWeight; {
	case ppmActualWeight < ppmRemainingEntitlement:
		return ppmActualWeight
	case ppmRemainingEntitlement < 0:
		return 0
	default:
		return ppmRemainingEntitlement
	}
}

// FetchMovingExpensesShipmentSummaryWorksheet fetches moving expenses for the Shipment Summary Worksheet
func FetchMovingExpensesShipmentSummaryWorksheet(move Move, db *pop.Connection, session *auth.Session) ([]MovingExpenseDocument, error) {
	var movingExpenses []MovingExpenseDocument
	if len(move.PersonallyProcuredMoves) > 0 {
		ppm := move.PersonallyProcuredMoves[0]
		moveDocuments, err := FetchApprovedMovingExpenseDocuments(db, session, ppm.ID)
		if err != nil {
			return movingExpenses, err
		}
		for _, moveDocument := range moveDocuments {
			if moveDocument.MovingExpenseDocument != nil {
				movingExpenses = append(movingExpenses, *moveDocument.MovingExpenseDocument)
			}
		}
	}
	return movingExpenses, nil
}

// FormatValuesShipmentSummaryWorksheetFormPage1 formats the data for page 1 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage1(data ShipmentSummaryFormData) ShipmentSummaryWorksheetPage1Values {
	page1 := ShipmentSummaryWorksheetPage1Values{}
	page1.MaxSITStorageEntitlement = "90 days per each shipment"
	// We don't currently know what allows POV to be authorized, so we are hardcoding it to "No" to start
	page1.POVAuthorized = "NO"
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
	page1.TAC = derefStringTypes(data.Order.TAC)
	page1.SAC = derefStringTypes(data.Order.SAC)

	page1.AuthorizedOrigin = FormatLocation(data.CurrentDutyStation)
	page1.AuthorizedDestination = FormatLocation(data.NewDutyStation)
	page1.NewDutyAssignment = FormatLocation(data.NewDutyStation)

	page1.WeightAllotment = FormatWeightAllotment(data)
	page1.WeightAllotmentProgear = FormatWeights(unit.Pound(data.WeightAllotment.ProGearWeight))
	page1.WeightAllotmentProgearSpouse = FormatWeights(unit.Pound(data.WeightAllotment.ProGearWeightSpouse))
	page1.TotalWeightAllotment = FormatWeights(data.TotalWeightAllotment)

	formattedShipments := FormatAllShipments(data.PersonallyProcuredMoves, data.Shipments)
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
	return page1
}

func formatActualObligationAdvance(data ShipmentSummaryFormData) string {
	if len(data.PersonallyProcuredMoves) > 0 && data.PersonallyProcuredMoves[0].Advance != nil {
		advance := data.PersonallyProcuredMoves[0].Advance.RequestedAmount.ToDollarFloat()
		return FormatDollars(advance)
	}
	return FormatDollars(0)
}

//FormatRank formats the service member's rank for Shipment Summary Worksheet
func FormatRank(rank *ServiceMemberRank) string {
	var rankDisplayValue = map[ServiceMemberRank]string{
		ServiceMemberRankE1:                "E-1",
		ServiceMemberRankE2:                "E-2",
		ServiceMemberRankE3:                "E-3",
		ServiceMemberRankE4:                "E-4",
		ServiceMemberRankE5:                "E-5",
		ServiceMemberRankE6:                "E-6",
		ServiceMemberRankE7:                "E-7",
		ServiceMemberRankE8:                "E-8",
		ServiceMemberRankE9:                "E-9",
		ServiceMemberRankO1ACADEMYGRADUATE: "O-1/Service Academy Graduate",
		ServiceMemberRankO2:                "O-2",
		ServiceMemberRankO3:                "O-3",
		ServiceMemberRankO4:                "O-4",
		ServiceMemberRankO5:                "O-5",
		ServiceMemberRankO6:                "O-6",
		ServiceMemberRankO7:                "O-7",
		ServiceMemberRankO8:                "O-8",
		ServiceMemberRankO9:                "O-9",
		ServiceMemberRankO10:               "O-10",
		ServiceMemberRankW1:                "W-1",
		ServiceMemberRankW2:                "W-2",
		ServiceMemberRankW3:                "W-3",
		ServiceMemberRankW4:                "W-4",
		ServiceMemberRankW5:                "W-5",
		ServiceMemberRankAVIATIONCADET:     "Aviation Cadet",
		ServiceMemberRankCIVILIANEMPLOYEE:  "Civilian Employee",
		ServiceMemberRankACADEMYCADET:      "Service Academy Cadet",
		ServiceMemberRankMIDSHIPMAN:        "Midshipman",
	}
	if rank != nil {
		return rankDisplayValue[*rank]
	}
	return ""
}

//FormatValuesShipmentSummaryWorksheetFormPage2 formats the data for page 2 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage2(data ShipmentSummaryFormData) (ShipmentSummaryWorksheetPage2Values, error) {
	var err error
	page2 := ShipmentSummaryWorksheetPage2Values{}
	page2.PreparationDate = FormatDate(data.PreparationDate)
	page2.FormattedMovingExpenses, err = FormatMovingExpenses(data.MovingExpenseDocuments)
	page2.TotalMemberPaidRepeated = page2.TotalMemberPaid
	page2.TotalGTCCPaidRepeated = page2.TotalGTCCPaid
	if err != nil {
		return page2, err
	}
	return page2, nil
}

//FormatWeightAllotment formats the weight allotment for Shipment Summary Worksheet
func FormatWeightAllotment(data ShipmentSummaryFormData) string {
	if data.Order.HasDependents {
		return FormatWeights(unit.Pound(data.WeightAllotment.TotalWeightSelfPlusDependents))
	}
	return FormatWeights(unit.Pound(data.WeightAllotment.TotalWeightSelf))
}

//FormatLocation formats AuthorizedOrigin and AuthorizedDestination for Shipment Summary Worksheet
func FormatLocation(dutyStation DutyStation) string {
	return fmt.Sprintf("%s, %s %s", dutyStation.Name, dutyStation.Address.State, dutyStation.Address.PostalCode)
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
func FormatAllShipments(ppms PersonallyProcuredMoves, shipments Shipments) ShipmentSummaryWorkSheetShipments {
	totalShipments := len(shipments) + len(ppms)
	formattedShipments := ShipmentSummaryWorkSheetShipments{}
	formattedNumberAndTypes := make([]string, totalShipments)
	formattedPickUpDates := make([]string, totalShipments)
	formattedShipmentWeights := make([]string, totalShipments)
	formattedShipmentStatuses := make([]string, totalShipments)
	var shipmentNumber int

	for _, shipment := range shipments {
		formattedNumberAndTypes[shipmentNumber] = FormatShipmentNumberAndType(shipmentNumber)
		formattedPickUpDates[shipmentNumber] = FormatShipmentPickupDate(shipment)
		formattedShipmentWeights[shipmentNumber] = FormatShipmentWeight(shipment)
		formattedShipmentStatuses[shipmentNumber] = FormatCurrentShipmentStatus(shipment)
		shipmentNumber++
	}
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

//FormatMovingExpenses formats moving expenses for Shipment Summary Worksheet
func FormatMovingExpenses(movingExpenseDocuments MovingExpenseDocuments) (FormattedMovingExpenses, error) {
	expenses := FormattedMovingExpenses{
		ContractedExpenseMemberPaid: FormatDollars(0),
		ContractedExpenseGTCCPaid:   FormatDollars(0),
		RentalEquipmentMemberPaid:   FormatDollars(0),
		RentalEquipmentGTCCPaid:     FormatDollars(0),
		PackingMaterialsMemberPaid:  FormatDollars(0),
		PackingMaterialsGTCCPaid:    FormatDollars(0),
		WeighingFeesMemberPaid:      FormatDollars(0),
		WeighingFeesGTCCPaid:        FormatDollars(0),
		GasMemberPaid:               FormatDollars(0),
		GasGTCCPaid:                 FormatDollars(0),
		TollsMemberPaid:             FormatDollars(0),
		TollsGTCCPaid:               FormatDollars(0),
		OilMemberPaid:               FormatDollars(0),
		OilGTCCPaid:                 FormatDollars(0),
		OtherMemberPaid:             FormatDollars(0),
		OtherGTCCPaid:               FormatDollars(0),
		TotalMemberPaid:             FormatDollars(0),
		TotalGTCCPaid:               FormatDollars(0),
	}
	subTotals := SubTotalExpenses(movingExpenseDocuments)
	formattedExpenses := make(map[string]string)
	for key, value := range subTotals {
		formattedExpenses[key] = FormatDollars(value)
	}

	err := mapstructure.Decode(formattedExpenses, &expenses)
	if err != nil {
		return expenses, err
	}
	return expenses, nil
}

//SubTotalExpenses groups moving expenses by type and payment method
func SubTotalExpenses(expenseDocuments MovingExpenseDocuments) map[string]float64 {
	var expenseType string
	totals := make(map[string]float64)
	for _, expense := range expenseDocuments {
		expenseType = getExpenseType(expense)
		expenseDollarAmt := expense.RequestedAmountCents.ToDollarFloat()
		totals[expenseType] += expenseDollarAmt
		addToGrandTotal(totals, expenseType, expenseDollarAmt)
	}
	return totals
}

func addToGrandTotal(totals map[string]float64, key string, expenseDollarAmt float64) {
	if strings.HasSuffix(key, "GTCCPaid") {
		totals["TotalGTCCPaid"] += expenseDollarAmt
	} else {
		totals["TotalMemberPaid"] += expenseDollarAmt
	}
}

func getExpenseType(expense MovingExpenseDocument) string {
	expenseType := FormatEnum(string(expense.MovingExpenseType), "")
	if expense.PaymentMethod == "GTCC" {
		return fmt.Sprintf("%s%s", expenseType, "GTCCPaid")
	}
	return fmt.Sprintf("%s%s", expenseType, "MemberPaid")
}

//FormatCurrentShipmentStatus formats FormatCurrentShipmentStatus for the Shipment Summary Worksheet
func FormatCurrentShipmentStatus(shipment Shipment) string {
	return FormatEnum(string(shipment.Status), " ")
}

//FormatCurrentPPMStatus formats FormatCurrentPPMStatus for the Shipment Summary Worksheet
func FormatCurrentPPMStatus(ppm PersonallyProcuredMove) string {
	if ppm.Status == "PAYMENT_REQUESTED" {
		return "At destination"
	}
	return FormatEnum(string(ppm.Status), " ")
}

//FormatShipmentNumberAndType formats FormatShipmentNumberAndType for the Shipment Summary Worksheet
func FormatShipmentNumberAndType(i int) string {
	return fmt.Sprintf("%02d - HHG (GBL)", i+1)
}

//FormatPPMNumberAndType formats FormatShipmentNumberAndType for the Shipment Summary Worksheet
func FormatPPMNumberAndType(i int) string {
	return fmt.Sprintf("%02d - PPM", i+1)
}

//FormatShipmentWeight formats a shipments ShipmentWeight for the Shipment Summary Worksheet
func FormatShipmentWeight(shipment Shipment) string {
	if shipment.NetWeight != nil {
		wtg := FormatWeights(*shipment.NetWeight)
		return fmt.Sprintf("%s lbs - FINAL", wtg)
	}
	return ""
}

//FormatShipmentPickupDate formats a shipments ActualPickupDate for the Shipment Summary Worksheet
func FormatShipmentPickupDate(shipment Shipment) string {
	if shipment.ActualPickupDate != nil {
		return FormatDate(*shipment.ActualPickupDate)
	}
	return ""
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
