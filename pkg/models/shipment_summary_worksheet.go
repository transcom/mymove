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
	ShipmentNumberAndTypes          string
	ShipmentPickUpDates             string
	ShipmentWeights                 string
	ShipmentCurrentShipmentStatuses string
	PreparationDate                 string
	MaxObligationGCC100             string
	TotalWeightAllotmentRepeat      string
	MaxObligationGCC95              string
	MaxObligationGCCMaxAdvance      string
	ActualWeight                    string
	ActualObligationGCC100          string
	ActualObligationGCC95           string
	ActualObligationAdvance         string
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
}

// ShipmentSummaryFormData is a container for the various objects required for the a Shipment Summary Worksheet
type ShipmentSummaryFormData struct {
	ServiceMember           ServiceMember
	Order                   Order
	CurrentDutyStation      DutyStation
	NewDutyStation          DutyStation
	WeightAllotment         WeightAllotment
	TotalWeightAllotment    int
	Shipments               Shipments
	PersonallyProcuredMoves PersonallyProcuredMoves
	PreparationDate         time.Time
	MaxObligation           Obligation
	ActualObligation        Obligation
	MovingExpenseDocuments  []MovingExpenseDocument
}

//Obligation an object representing the obligations section on the shipment summary worksheet
type Obligation struct {
	Gcc unit.Cents
}

//GCC100 calculates the 95% GCC on shipment summary worksheet
func (obligation Obligation) GCC100() float64 {
	return obligation.Gcc.ToDollarFloat()
}

//GCC95 calculates the 100% GCC on shipment summary worksheet
func (obligation Obligation) GCC95() float64 {
	return obligation.Gcc.MultiplyFloat64(.95).ToDollarFloat()
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
	var totalEntitlement int
	if serviceMember.Rank != nil {
		rank = ServiceMemberRank(*serviceMember.Rank)
		weightAllotment = GetWeightAllotment(rank)
		totalEntitlement, err = GetEntitlement(rank, move.Orders.HasDependents, move.Orders.SpouseHasProGear)
		if err != nil {
			return ShipmentSummaryFormData{}, err
		}
	}

	movingExpenses, err := FetchMovingExpensesShipmentSummaryWorksheet(move, db, session)
	if err != nil {
		return ShipmentSummaryFormData{}, err
	}

	ssd := ShipmentSummaryFormData{
		ServiceMember:           serviceMember,
		Order:                   move.Orders,
		CurrentDutyStation:      serviceMember.DutyStation,
		NewDutyStation:          move.Orders.NewDutyStation,
		WeightAllotment:         weightAllotment,
		TotalWeightAllotment:    totalEntitlement,
		Shipments:               move.Shipments,
		PersonallyProcuredMoves: move.PersonallyProcuredMoves,
		MovingExpenseDocuments:  movingExpenses,
	}
	return ssd, nil
}

//FetchMovingExpensesShipmentSummaryWorksheet fetches moving expenses for the Shipment Summary Worksheet
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

	page1.AuthorizedOrigin = FormatAuthorizedLocation(data.CurrentDutyStation)
	page1.AuthorizedDestination = FormatAuthorizedLocation(data.NewDutyStation)
	page1.NewDutyAssignment = FormatDutyStation(data.NewDutyStation)

	page1.WeightAllotment = FormatWeightAllotment(data)
	page1.WeightAllotmentProgear = FormatWeights(data.WeightAllotment.ProGearWeight)
	page1.WeightAllotmentProgearSpouse = FormatWeights(data.WeightAllotment.ProGearWeightSpouse)
	page1.TotalWeightAllotment = FormatWeights(data.TotalWeightAllotment)

	formattedShipments := FormatAllShipments(data.PersonallyProcuredMoves, data.Shipments)
	page1.ShipmentNumberAndTypes = formattedShipments.ShipmentNumberAndTypes
	page1.ShipmentPickUpDates = formattedShipments.PickUpDates
	page1.ShipmentCurrentShipmentStatuses = formattedShipments.CurrentShipmentStatuses
	page1.ShipmentWeights = formattedShipments.ShipmentWeights

	page1.MaxObligationGCC100 = FormatDollars(data.MaxObligation.GCC100())
	page1.TotalWeightAllotmentRepeat = page1.TotalWeightAllotment
	page1.MaxObligationGCC95 = FormatDollars(data.MaxObligation.GCC95())
	page1.MaxObligationGCCMaxAdvance = FormatDollars(data.MaxObligation.MaxAdvance())

	page1.ActualObligationGCC100 = FormatDollars(data.ActualObligation.GCC100())
	page1.ActualWeight = formatActualWeight(data)
	page1.ActualObligationAdvance = formatActualObligationAdvance(data)

	page1.ActualObligationGCC95 = FormatDollars(data.ActualObligation.GCC95())
	return page1
}

func formatActualObligationAdvance(data ShipmentSummaryFormData) string {
	if len(data.PersonallyProcuredMoves) > 0 && data.PersonallyProcuredMoves[0].Advance != nil {
		advance := data.PersonallyProcuredMoves[0].Advance.RequestedAmount.ToDollarFloat()
		return FormatDollars(advance)
	}
	return FormatDollars(0)
}

func formatActualWeight(data ShipmentSummaryFormData) string {
	if len(data.PersonallyProcuredMoves) > 0 && data.PersonallyProcuredMoves[0].NetWeight != nil {
		return FormatWeights(int(*(data.PersonallyProcuredMoves[0].NetWeight)))
	}
	return FormatWeights(0)
}

//FormatRank formats the service member's rank for Shipment Summary Worksheet
func FormatRank(rank *ServiceMemberRank) string {
	var rankDisplayValue = map[ServiceMemberRank]string{
		ServiceMemberRankE1:                     "E-1",
		ServiceMemberRankE2:                     "E-2",
		ServiceMemberRankE3:                     "E-3",
		ServiceMemberRankE4:                     "E-4",
		ServiceMemberRankE5:                     "E-5",
		ServiceMemberRankE6:                     "E-6",
		ServiceMemberRankE7:                     "E-7",
		ServiceMemberRankE8:                     "E-8",
		ServiceMemberRankE9:                     "E-9",
		ServiceMemberRankO1W1ACADEMYGRADUATE:    "O-1/W-1/Service Academy Graduate",
		ServiceMemberRankO2W2:                   "O-2/W-2",
		ServiceMemberRankO3W3:                   "O-3/W-3",
		ServiceMemberRankO4W4:                   "O-4/W-4",
		ServiceMemberRankO5W5:                   "O-5/W-5",
		ServiceMemberRankO6:                     "O-6",
		ServiceMemberRankO7:                     "O-7",
		ServiceMemberRankO8:                     "O-8",
		ServiceMemberRankO9:                     "O-9",
		ServiceMemberRankO10:                    "O-10",
		ServiceMemberRankAVIATIONCADET:          "Aviation Cadet",
		ServiceMemberRankCIVILIANEMPLOYEE:       "Civilian Employee",
		ServiceMemberRankACADEMYCADETMIDSHIPMAN: "Service Academy Cadet/Midshipman",
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
	if err != nil {
		return page2, err
	}
	return page2, nil
}

//FormatWeightAllotment formats the weight allotment for Shipment Summary Worksheet
func FormatWeightAllotment(data ShipmentSummaryFormData) string {
	if data.Order.HasDependents {
		return FormatWeights(data.WeightAllotment.TotalWeightSelfPlusDependents)
	}
	return FormatWeights(data.WeightAllotment.TotalWeightSelf)
}

//FormatAuthorizedLocation formats AuthorizedOrigin and AuthorizedDestination for Shipment Summary Worksheet
func FormatAuthorizedLocation(dutyStation DutyStation) string {
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
		wtg := FormatWeights(int(*shipment.NetWeight))
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
		wtg := FormatWeights(int(*ppm.NetWeight))
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

//FormatDutyStation formats DutyStation for Shipment Summary Worksheet
func FormatDutyStation(dutyStation DutyStation) string {
	return fmt.Sprintf("%s, %s", dutyStation.Name, dutyStation.Address.State)
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

//FormatWeights formats an int using 000s separator
func FormatWeights(wtg int) string {
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
