package services

import (
	"github.com/gofrs/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
)

// Dollar represents a type for dollar monetary unit
type Dollar float64

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
	WeightAllotmentProGear          string
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
	PreparationDate1                string
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
	MailingAddressW2                string
}

// Page2Values is an object representing a Shipment Summary Worksheet
type Page2Values struct {
	CUIBanner                   string
	PreparationDate2            string
	TAC                         string
	SAC                         string
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
	TotalPaidNonSIT             string
	TotalMemberPaidSIT          string
	TotalGTCCPaidSIT            string
	TotalPaidSIT                string
	ShipmentPickupDates         string
	TrustedAgentName            string
	FormattedMovingExpenses
	ServiceMemberSignature string
	PPPOPPSORepresentative string
	SignatureDate          string
	FormattedOtherExpenses
}

// FormattedOtherExpenses is an object representing the other moving expenses formatted for the SSW
type FormattedOtherExpenses struct {
	Descriptions string
	AmountsPaid  string
}

// FormattedMovingExpenses is an object representing the service member's moving expenses formatted for the SSW
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
	TotalPaidNonSIT             string
	TotalMemberPaidSIT          string
	TotalGTCCPaidSIT            string
	TotalPaidSIT                string
}

//go:generate mockery --name SSWPPMComputer
type SSWPPMComputer interface {
	FetchDataShipmentSummaryWorksheetFormData(appCtx appcontext.AppContext, _ *auth.Session, ppmShipmentID uuid.UUID) (*models.ShipmentSummaryFormData, error)
	ComputeObligations(_ appcontext.AppContext, _ models.ShipmentSummaryFormData, _ route.Planner) (models.Obligations, error)
	FormatValuesShipmentSummaryWorksheet(shipmentSummaryFormData models.ShipmentSummaryFormData, isPaymentPacket bool) (Page1Values, Page2Values, error)
	FormatAllShipments(data models.ShipmentSummaryFormData) models.WorkSheetShipments
	FormatValuesShipmentSummaryWorksheetFormPage1(data models.ShipmentSummaryFormData, isPaymentPacket bool) (Page1Values, error)
	FormatValuesShipmentSummaryWorksheetFormPage2(data models.ShipmentSummaryFormData, isPaymentPacket bool) (Page2Values, error)
}

//go:generate mockery --name SSWPPMGenerator
type SSWPPMGenerator interface {
	FillSSWPDFForm(Page1Values, Page2Values) (afero.File, *pdfcpu.PDFInfo, error)
}
