package services

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
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
	MailingAddressW2                string
}

// Page2Values is an object representing a Shipment Summary Worksheet
type Page2Values struct {
	CUIBanner                   string
	PreparationDate             string
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

// ShipmentSummaryFormData is a container for the various objects required for the a Shipment Summary Worksheet
type ShipmentSummaryFormData struct {
	ServiceMember            models.ServiceMember
	Order                    models.Order
	Move                     models.Move
	CurrentDutyLocation      models.DutyLocation
	NewDutyLocation          models.DutyLocation
	WeightAllotment          SSWMaxWeightEntitlement
	PPMShipment              models.PPMShipment
	PPMShipments             models.PPMShipments
	W2Address                *models.Address
	PreparationDate          time.Time
	Obligations              Obligations
	MovingExpenses           models.MovingExpenses
	PPMRemainingEntitlement  unit.Pound
	SignedCertifications     []*models.SignedCertification
	MaxSITStorageEntitlement int
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

// SSWMaxWeightEntitlement weight allotment for the shipment summary worksheet.
type SSWMaxWeightEntitlement struct {
	Entitlement   unit.Pound
	ProGear       unit.Pound
	SpouseProGear unit.Pound
	TotalWeight   unit.Pound
}

//go:generate mockery --name SSWPPMComputer
type SSWPPMComputer interface {
	FetchDataShipmentSummaryWorksheetFormData(appCtx appcontext.AppContext, _ *auth.Session, ppmShipmentID uuid.UUID) (*ShipmentSummaryFormData, error)
	ComputeObligations(_ appcontext.AppContext, _ ShipmentSummaryFormData, _ route.Planner) (Obligations, error)
	FormatValuesShipmentSummaryWorksheet(shipmentSummaryFormData ShipmentSummaryFormData) (Page1Values, Page2Values)
}

//go:generate mockery --name SSWPPMGenerator
type SSWPPMGenerator interface {
	FillSSWPDFForm(Page1Values, Page2Values) (afero.File, *pdfcpu.PDFInfo, error)
}
