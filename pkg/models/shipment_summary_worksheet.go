package models

import (
	"reflect"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/transcom/mymove/pkg/unit"
)

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
	ServiceMember           ServiceMember
	Order                   Order
	CurrentDutyLocation     DutyLocation
	NewDutyLocation         DutyLocation
	WeightAllotment         SSWMaxWeightEntitlement
	PPMShipments            PPMShipments
	PreparationDate         time.Time
	Obligations             Obligations
	MovingExpenses          []MovingExpense
	PPMRemainingEntitlement unit.Pound
	SignedCertification     SignedCertification
}

// SSWMaxWeightEntitlement weight allotment for the shipment summary worksheet.
type SSWMaxWeightEntitlement struct {
	Entitlement   unit.Pound
	ProGear       unit.Pound
	SpouseProGear unit.Pound
	TotalWeight   unit.Pound
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

// QUESTION: do we need this? do we have to use reflect here?
// adds a line item to shipment summary worksheet SSWMaxWeightEntitlement and increments total allotment
func (wa *SSWMaxWeightEntitlement) AddLineItem(field string, value int) {
	r := reflect.ValueOf(wa).Elem()
	f := r.FieldByName(field)
	if f.IsValid() && f.CanSet() {
		f.SetInt(int64(value))
		wa.TotalWeight += unit.Pound(value)
	}
}
