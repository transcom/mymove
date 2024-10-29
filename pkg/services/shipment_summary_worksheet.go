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
	ActualObligationGCC100          string
	ActualObligationGCC95           string
	ActualObligationAdvance         string
	ActualObligationSIT             string
	MileageTotal                    string
	MailingAddressW2                string
	IsActualExpenseReimbursement    bool
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
	Disbursement                string
	ShipmentPickupDates         string
	TrustedAgentName            string
	ServiceMemberSignature      string
	PPPOPPSORepresentative      string
	SignatureDate               string
	PPMRemainingEntitlement     string
	FormattedMovingExpenses
	FormattedOtherExpenses
}

// Page3Values is an object representing a Shipment Summary Worksheet
type Page3Values struct {
	CUIBanner        string
	PreparationDate3 string
	AddShipments     map[string]string
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
	FormatValuesShipmentSummaryWorksheet(shipmentSummaryFormData models.ShipmentSummaryFormData, isPaymentPacket bool) (Page1Values, Page2Values, Page3Values, error)
	FormatShipment(ppm models.PPMShipment, weightAllotment models.SSWMaxWeightEntitlement, isPaymentPacket bool) models.WorkSheetShipment
	FormatValuesShipmentSummaryWorksheetFormPage1(data models.ShipmentSummaryFormData, isPaymentPacket bool) (Page1Values, error)
	FormatValuesShipmentSummaryWorksheetFormPage2(data models.ShipmentSummaryFormData, isPaymentPacket bool) (Page2Values, error)
	FormatValuesShipmentSummaryWorksheetFormPage3(data models.ShipmentSummaryFormData, isPaymentPacket bool) (Page3Values, error)
}

//go:generate mockery --name SSWPPMGenerator
type SSWPPMGenerator interface {
	FillSSWPDFForm(Page1Values, Page2Values, Page3Values) (afero.File, *pdfcpu.PDFInfo, error)
}
