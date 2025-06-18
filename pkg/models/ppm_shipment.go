package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/unit"
)

type HaulType string

const (
	LINEHAUL  HaulType = "Linehaul"
	SHORTHAUL HaulType = "Shorthaul"
)

type PPMCloseout struct {
	ID                    *uuid.UUID
	PlannedMoveDate       *time.Time
	ActualMoveDate        *time.Time
	Miles                 *int
	EstimatedWeight       *unit.Pound
	ActualWeight          *unit.Pound
	ProGearWeightCustomer *unit.Pound
	ProGearWeightSpouse   *unit.Pound
	GrossIncentive        *unit.Cents
	GCC                   *unit.Cents
	AOA                   *unit.Cents
	RemainingIncentive    *unit.Cents
	HaulPrice             *unit.Cents
	HaulFSC               *unit.Cents
	HaulType              *HaulType
	DOP                   *unit.Cents
	DDP                   *unit.Cents
	PackPrice             *unit.Cents
	UnpackPrice           *unit.Cents
	IntlPackPrice         *unit.Cents
	IntlUnpackPrice       *unit.Cents
	IntlLinehaulPrice     *unit.Cents
	SITReimbursement      *unit.Cents
	GCCMultiplier         *float64
}

type PPMActualWeight struct {
	ActualWeight *unit.Pound
}

type PPMSITEstimatedCost struct {
	PPMSITEstimatedCost *unit.Pound
}

type PPMSITEstimatedCostParams struct {
	ContractYearName       string
	PriceRateOrFactor      string
	IsPeak                 string
	EscalationCompounded   string
	ServiceAreaOrigin      string
	ServiceAreaDestination string
	NumberDaysSIT          string
}

type PPMSITEstimatedCostInfo struct {
	EstimatedSITCost       *unit.Cents
	PriceFirstDaySIT       *unit.Cents
	PriceAdditionalDaySIT  *unit.Cents
	ParamsFirstDaySIT      PPMSITEstimatedCostParams
	ParamsAdditionalDaySIT PPMSITEstimatedCostParams
}

// PPMShipmentStatus represents the status of an order record's lifecycle
type PPMShipmentStatus string

const (
	// PPMShipmentStatusCanceled captures enum value "CANCELED"
	PPMShipmentStatusCanceled PPMShipmentStatus = "CANCELED"
	// PPMShipmentStatusDraft captures enum value "DRAFT"
	PPMShipmentStatusDraft PPMShipmentStatus = "DRAFT"
	// PPMShipmentStatusSubmitted captures enum value "SUBMITTED"
	PPMShipmentStatusSubmitted PPMShipmentStatus = "SUBMITTED"
	// PPMShipmentStatusWaitingOnCustomer captures enum value "WAITING_ON_CUSTOMER"
	PPMShipmentStatusWaitingOnCustomer PPMShipmentStatus = "WAITING_ON_CUSTOMER"
	// PPMShipmentStatusNeedsAdvanceApproval captures enum value "NEEDS_ADVANCE_APPROVAL"
	PPMShipmentStatusNeedsAdvanceApproval PPMShipmentStatus = "NEEDS_ADVANCE_APPROVAL"
	// PPMShipmentStatusNeedsCloseout captures enum value "NEEDS_CLOSEOUT"
	PPMShipmentStatusNeedsCloseout PPMShipmentStatus = "NEEDS_CLOSEOUT"
	// PPMShipmentStatusCloseoutComplete captures enum value "CLOSEOUT_COMPLETE"
	PPMShipmentStatusCloseoutComplete PPMShipmentStatus = "CLOSEOUT_COMPLETE"
	// PPMStatusCOMPLETED captures enum value "COMPLETED"
	PPMShipmentStatusComplete PPMShipmentStatus = "COMPLETED"
)

// AllowedPPMShipmentStatuses is a list of all the allowed values for the Status of a PPMShipment as strings. Needed for
// validation.
var AllowedPPMShipmentStatuses = []string{
	string(PPMShipmentStatusCanceled),
	string(PPMShipmentStatusDraft),
	string(PPMShipmentStatusSubmitted),
	string(PPMShipmentStatusWaitingOnCustomer),
	string(PPMShipmentStatusNeedsAdvanceApproval),
	string(PPMShipmentStatusNeedsCloseout),
	string(PPMShipmentStatusCloseoutComplete),
}

// PPMAdvanceStatus represents the status of an advance that can be approved, edited or rejected by a SC
type PPMAdvanceStatus string

const (
	// PPMAdvanceStatusApproved captures enum value "APPROVED"
	PPMAdvanceStatusApproved PPMAdvanceStatus = "APPROVED"
	// PPMAdvanceStatusRejected captures enum value "REJECTED"
	PPMAdvanceStatusRejected PPMAdvanceStatus = "REJECTED"
	// PPMAdvanceStatusReceived captures enum value "RECEIVED"
	PPMAdvanceStatusReceived PPMAdvanceStatus = "RECEIVED"
	// PPMAdvanceStatusNotReceived captures enum value "NOT RECEIVED"
	PPMAdvanceStatusNotReceived PPMAdvanceStatus = "NOT_RECEIVED"
)

// AllowedPPMAdvanceStatuses is a list of all the allowed values for AdvanceStatus on a PPMShipment, as strings. Needed
// for validation.
var AllowedPPMAdvanceStatuses = []string{
	string(PPMAdvanceStatusApproved),
	string(PPMAdvanceStatusRejected),
	string(PPMAdvanceStatusReceived),
	string(PPMAdvanceStatusNotReceived),
}

// SITLocationType represents whether the SIT at the origin or destination
type SITLocationType string

const (
	// SITLocationTypeOrigin captures enum value "ORIGIN"
	SITLocationTypeOrigin SITLocationType = "ORIGIN"
	// SITLocationTypeDestination captures enum value "DESTINATION"
	SITLocationTypeDestination SITLocationType = "DESTINATION"
)

// AllowedSITLocationTypes is a list of all the allowed values for SITLocationType on a PPMShipment, as strings. Needed
// for validation.
var AllowedSITLocationTypes = []string{
	string(SITLocationTypeOrigin),
	string(SITLocationTypeDestination),
}

// PPMDocumentStatus represents the status of a PPMShipment's documents. Lives here since we have multiple PPM document
// models.
type PPMDocumentStatus string

const (
	// PPMDocumentStatusApproved captures enum value "DRAFT"
	PPMDocumentStatusDRAFT PPMDocumentStatus = "DRAFT"
	// PPMDocumentStatusApproved captures enum value "APPROVED"
	PPMDocumentStatusApproved PPMDocumentStatus = "APPROVED"
	// PPMDocumentStatusExcluded captures enum value "EXCLUDED"
	PPMDocumentStatusExcluded PPMDocumentStatus = "EXCLUDED"
	// PPMDocumentStatusRejected captures enum value "REJECTED"
	PPMDocumentStatusRejected PPMDocumentStatus = "REJECTED"
)

// AllowedPPMDocumentStatuses is a list of all the allowed values for the Status of a PPMShipment's documents, as
// strings. Needed for validation.
var AllowedPPMDocumentStatuses = []string{
	string(PPMDocumentStatusApproved),
	string(PPMDocumentStatusExcluded),
	string(PPMDocumentStatusRejected),
}

// PPMDocuments is a collection of the different PPMShipment documents. This type exists mainly to make it easier to
// work with the group of documents as a whole when we don't actually retrieve the PPM Shipment itself.
type PPMDocuments struct {
	WeightTickets
	MovingExpenses
	ProgearWeightTickets
}

// PPMType represents the type of a PPM shipment
type PPMType string

const (
	// PPMTypeIncentiveBased captures enum value "INCENTIVE_BASED"
	PPMTypeIncentiveBased PPMType = "INCENTIVE_BASED"
	// PPMTypeActualExpense captures enum value "ACTUAL_EXPENSE"
	PPMTypeActualExpense PPMType = "ACTUAL_EXPENSE"
	// PPMTypeSmallPackage captures enum value "SMALL_PACKAGE"
	PPMTypeSmallPackage PPMType = "SMALL_PACKAGE"
)

// AllowedPPMTypes is a list of all the allowed values for PPM types
var AllowedPPMTypes = []string{
	string(PPMTypeIncentiveBased),
	string(PPMTypeActualExpense),
	string(PPMTypeSmallPackage),
}

// PPMShipment is the portion of a move that a service member performs themselves
type PPMShipment struct {
	ID                             uuid.UUID            `json:"id" db:"id"`
	PPMType                        PPMType              `json:"ppm_type" db:"ppm_type"`
	ShipmentID                     uuid.UUID            `json:"shipment_id" db:"shipment_id"`
	Shipment                       MTOShipment          `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	CreatedAt                      time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt                      time.Time            `json:"updated_at" db:"updated_at"`
	DeletedAt                      *time.Time           `json:"deleted_at" db:"deleted_at"`
	Status                         PPMShipmentStatus    `json:"status" db:"status"`
	ExpectedDepartureDate          time.Time            `json:"expected_departure_date" db:"expected_departure_date"`
	ActualMoveDate                 *time.Time           `json:"actual_move_date" db:"actual_move_date"`
	SubmittedAt                    *time.Time           `json:"submitted_at" db:"submitted_at"`
	ReviewedAt                     *time.Time           `json:"reviewed_at" db:"reviewed_at"`
	ApprovedAt                     *time.Time           `json:"approved_at" db:"approved_at"`
	W2Address                      *Address             `belongs_to:"addresses" fk_id:"w2_address_id"`
	W2AddressID                    *uuid.UUID           `db:"w2_address_id"`
	PickupAddress                  *Address             `belongs_to:"addresses" fk_id:"pickup_postal_address_id"`
	PickupAddressID                *uuid.UUID           `db:"pickup_postal_address_id"`
	SecondaryPickupAddress         *Address             `belongs_to:"addresses" fk_id:"secondary_pickup_postal_address_id"`
	SecondaryPickupAddressID       *uuid.UUID           `db:"secondary_pickup_postal_address_id"`
	HasSecondaryPickupAddress      *bool                `db:"has_secondary_pickup_address"`
	TertiaryPickupAddress          *Address             `belongs_to:"addresses" fk_id:"tertiary_pickup_postal_address_id"`
	TertiaryPickupAddressID        *uuid.UUID           `db:"tertiary_pickup_postal_address_id"`
	HasTertiaryPickupAddress       *bool                `db:"has_tertiary_pickup_address"`
	DestinationAddress             *Address             `belongs_to:"addresses" fk_id:"destination_postal_address_id"`
	DestinationAddressID           *uuid.UUID           `db:"destination_postal_address_id"`
	SecondaryDestinationAddress    *Address             `belongs_to:"addresses" fk_id:"secondary_destination_postal_address_id"`
	SecondaryDestinationAddressID  *uuid.UUID           `db:"secondary_destination_postal_address_id"`
	HasSecondaryDestinationAddress *bool                `db:"has_secondary_destination_address"`
	TertiaryDestinationAddress     *Address             `belongs_to:"addresses" fk_id:"tertiary_destination_postal_address_id"`
	TertiaryDestinationAddressID   *uuid.UUID           `db:"tertiary_destination_postal_address_id"`
	HasTertiaryDestinationAddress  *bool                `db:"has_tertiary_destination_address"`
	EstimatedWeight                *unit.Pound          `json:"estimated_weight" db:"estimated_weight"`
	AllowableWeight                *unit.Pound          `json:"allowable_weight" db:"allowable_weight"`
	HasProGear                     *bool                `json:"has_pro_gear" db:"has_pro_gear"`
	ProGearWeight                  *unit.Pound          `json:"pro_gear_weight" db:"pro_gear_weight"`
	SpouseProGearWeight            *unit.Pound          `json:"spouse_pro_gear_weight" db:"spouse_pro_gear_weight"`
	EstimatedIncentive             *unit.Cents          `json:"estimated_incentive" db:"estimated_incentive"`
	MaxIncentive                   *unit.Cents          `json:"max_incentive" db:"max_incentive"`
	FinalIncentive                 *unit.Cents          `json:"final_incentive" db:"final_incentive"`
	HasRequestedAdvance            *bool                `json:"has_requested_advance" db:"has_requested_advance"`
	AdvanceAmountRequested         *unit.Cents          `json:"advance_amount_requested" db:"advance_amount_requested"`
	HasReceivedAdvance             *bool                `json:"has_received_advance" db:"has_received_advance"`
	AdvanceStatus                  *PPMAdvanceStatus    `json:"advance_status" db:"advance_status"`
	AdvanceAmountReceived          *unit.Cents          `json:"advance_amount_received" db:"advance_amount_received"`
	SITExpected                    *bool                `json:"sit_expected" db:"sit_expected"`
	SITLocation                    *SITLocationType     `json:"sit_location" db:"sit_location"`
	SITEstimatedWeight             *unit.Pound          `json:"sit_estimated_weight" db:"sit_estimated_weight"`
	SITEstimatedEntryDate          *time.Time           `json:"sit_estimated_entry_date" db:"sit_estimated_entry_date"`
	SITEstimatedDepartureDate      *time.Time           `json:"sit_estimated_departure_date" db:"sit_estimated_departure_date"`
	SITEstimatedCost               *unit.Cents          `json:"sit_estimated_cost" db:"sit_estimated_cost"`
	WeightTickets                  WeightTickets        `has_many:"weight_tickets" fk_id:"ppm_shipment_id" order_by:"created_at asc"`
	MovingExpenses                 MovingExpenses       `has_many:"moving_expenses" fk_id:"ppm_shipment_id" order_by:"created_at asc"`
	ProgearWeightTickets           ProgearWeightTickets `has_many:"progear_weight_tickets" fk_id:"ppm_shipment_id" order_by:"created_at asc"`
	SignedCertification            *SignedCertification `has_one:"signed_certification" fk_id:"ppm_id"`
	AOAPacketID                    *uuid.UUID           `json:"aoa_packet_id" db:"aoa_packet_id"`
	AOAPacket                      *Document            `belongs_to:"documents" fk_id:"aoa_packet_id"`
	PaymentPacketID                *uuid.UUID           `json:"payment_packet_id" db:"payment_packet_id"`
	PaymentPacket                  *Document            `belongs_to:"documents" fk_id:"payment_packet_id"`
	IsActualExpenseReimbursement   *bool                `json:"is_actual_expense_reimbursement" db:"is_actual_expense_reimbursement"`
	HasGunSafe                     *bool                `json:"has_gun_safe" db:"has_gun_safe"`
	GunSafeWeight                  *unit.Pound          `json:"gun_safe_weight" db:"gun_safe_weight"`
	GCCMultiplierID                *uuid.UUID           `json:"gcc_multiplier_id" db:"gcc_multiplier_id"`
	GCCMultiplier                  *GCCMultiplier       `belongs_to:"gcc_multipliers" fk_id:"gcc_multiplier_id"`
}

// TableName overrides the table name used by Pop.
func (p PPMShipment) TableName() string {
	return "ppm_shipments"
}

// Cancel marks the PPM as Canceled
func (p *PPMShipment) CancelShipment() error {
	p.Status = PPMShipmentStatusCanceled
	return nil
}

// PPMShipments is a list of PPMs
type PPMShipments []PPMShipment

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (p PPMShipment) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "ShipmentID", Field: p.ShipmentID},
		&validators.StringInclusion{Name: "PPMType", Field: string(p.PPMType), List: AllowedPPMTypes},
		&OptionalTimeIsPresent{Name: "DeletedAt", Field: p.DeletedAt},
		&validators.TimeIsPresent{Name: "ExpectedDepartureDate", Field: p.ExpectedDepartureDate},
		&validators.StringInclusion{Name: "Status", Field: string(p.Status), List: AllowedPPMShipmentStatuses},
		&OptionalTimeIsPresent{Name: "ActualMoveDate", Field: p.ActualMoveDate},
		&OptionalTimeIsPresent{Name: "SubmittedAt", Field: p.SubmittedAt},
		&OptionalTimeIsPresent{Name: "ReviewedAt", Field: p.ReviewedAt},
		&OptionalTimeIsPresent{Name: "ApprovedAt", Field: p.ApprovedAt},
		&OptionalUUIDIsPresent{Name: "W2AddressID", Field: p.W2AddressID},
		&OptionalUUIDIsPresent{Name: "PickupAddressID", Field: p.PickupAddressID},
		&OptionalUUIDIsPresent{Name: "SecondaryPickupAddressID", Field: p.SecondaryPickupAddressID},
		&OptionalUUIDIsPresent{Name: "DestinationAddressID", Field: p.DestinationAddressID},
		&OptionalUUIDIsPresent{Name: "SecondaryDestinationAddressID", Field: p.SecondaryDestinationAddressID},
		&OptionalPoundIsNonNegative{Name: "EstimatedWeight", Field: p.EstimatedWeight},
		&OptionalPoundIsNonNegative{Name: "AllowableWeight", Field: p.AllowableWeight},
		&OptionalPoundIsNonNegative{Name: "ProGearWeight", Field: p.ProGearWeight},
		&OptionalPoundIsNonNegative{Name: "SpouseProGearWeight", Field: p.SpouseProGearWeight},
		&OptionalPoundIsNonNegative{Name: "GunSafeWeight", Field: p.GunSafeWeight},
		&OptionalPoundIsMax{Name: "GunSafeWeight", Field: p.GunSafeWeight, Max: 500},
		&OptionalCentIsNotNegative{Name: "EstimatedIncentive", Field: p.EstimatedIncentive},
		&OptionalCentIsNotNegative{Name: "MaxIncentive", Field: p.MaxIncentive},
		&OptionalCentIsPositive{Name: "FinalIncentive", Field: p.FinalIncentive},
		&OptionalCentIsNotNegative{Name: "AdvanceAmountRequested", Field: p.AdvanceAmountRequested},
		&OptionalStringInclusion{Name: "AdvanceStatus", List: AllowedPPMAdvanceStatuses, Field: (*string)(p.AdvanceStatus)},
		&OptionalCentIsPositive{Name: "AdvanceAmountReceived", Field: p.AdvanceAmountReceived},
		&OptionalStringInclusion{Name: "SITLocation", List: AllowedSITLocationTypes, Field: (*string)(p.SITLocation)},
		&OptionalPoundIsNonNegative{Name: "SITEstimatedWeight", Field: p.SITEstimatedWeight},
		&OptionalTimeIsPresent{Name: "SITEstimatedEntryDate", Field: p.SITEstimatedEntryDate},
		&OptionalTimeIsPresent{Name: "SITEstimatedDepartureDate", Field: p.SITEstimatedDepartureDate},
		&OptionalCentIsPositive{Name: "SITEstimatedCost", Field: p.SITEstimatedCost},
		&OptionalUUIDIsPresent{Name: "AOAPacketID", Field: p.AOAPacketID},
		&OptionalUUIDIsPresent{Name: "PaymentPacketID", Field: p.PaymentPacketID},
		&OptionalUUIDIsPresent{Name: "GCCMultiplierID", Field: p.GCCMultiplierID},
	), nil

}
func GetPPMNetWeight(ppm PPMShipment) unit.Pound {
	totalNetWeight := unit.Pound(0)

	// small package PPMs do not have weight tickets so we will add up approved moving expenses
	if ppm.PPMType == PPMTypeSmallPackage {
		if len(ppm.MovingExpenses) >= 1 {
			for _, movingExpense := range ppm.MovingExpenses {
				if movingExpense.WeightShipped != nil && movingExpense.Status != nil && *movingExpense.Status != PPMDocumentStatusRejected {
					totalNetWeight += *movingExpense.WeightShipped
				}
			}
			return totalNetWeight
		} else {
			return unit.Pound(0)
		}
	}

	// incentive-based and actual expense PPMs have weight tickets
	weightTickets := ppm.WeightTickets.FilterRejected()
	for _, weightTicket := range weightTickets {
		if weightTicket.AdjustedNetWeight != nil && *weightTicket.AdjustedNetWeight > 0 {
			totalNetWeight += *weightTicket.AdjustedNetWeight
		} else {
			totalNetWeight += GetWeightTicketNetWeight(weightTicket)
		}
	}
	return totalNetWeight
}

// FetchPPMShipmentByPPMShipmentID returns a PPM Shipment for a given id
func FetchPPMShipmentByPPMShipmentID(db *pop.Connection, ppmShipmentID uuid.UUID) (*PPMShipment, error) {
	var ppmShipment PPMShipment
	err := db.Q().Find(&ppmShipment, ppmShipmentID)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		return nil, err
	}
	return &ppmShipment, nil
}

type PPMIncentiveOCONUS struct {
	TotalIncentive int `db:"total_incentive"`
	PriceISLH      int `db:"price_islh"`
	PriceIHPK      int `db:"price_ihpk"`
	PriceIHUPK     int `db:"price_ihupk"`
	PriceFSC       int `db:"price_fsc"`
}

// a db function that will handle updating the estimated_incentive value
// this simulates pricing of a basic iHHG shipment with ISLH, IHPK, IHUPK, and the CONUS portion for a FSC
func CalculatePPMIncentive(db *pop.Connection, ppmID uuid.UUID, pickupAddressID uuid.UUID, destAddressID uuid.UUID, moveDate time.Time, mileage int, weight int, isEstimated bool, isActual bool, isMax bool) (*PPMIncentiveOCONUS, error) {
	var incentive PPMIncentiveOCONUS

	err := db.RawQuery("SELECT * FROM calculate_ppm_incentive($1, $2, $3, $4, $5, $6, $7, $8, $9)", ppmID, pickupAddressID, destAddressID, moveDate, mileage, weight, isEstimated, isActual, isMax).
		First(&incentive)
	if err != nil {
		return nil, fmt.Errorf("error calculating PPM incentive for PPM ID %s: %w", ppmID, err)
	}

	return &incentive, nil
}

type PPMSITCosts struct {
	TotalSITCost     int `db:"total_cost"`
	PriceFirstDaySIT int `db:"price_first_day"`
	PriceAddlDaySIT  int `db:"price_addl_day"`
}

// a db function that will handle calculating and returning the SIT costs related to a PPM shipment
func CalculatePPMSITCost(db *pop.Connection, ppmID uuid.UUID, addressID uuid.UUID, isOrigin bool, moveDate time.Time, weight int, sitDays int) (*PPMSITCosts, error) {
	var costs PPMSITCosts

	err := db.RawQuery("SELECT * FROM calculate_ppm_SIT_cost($1, $2, $3, $4, $5, $6)", ppmID, addressID, isOrigin, moveDate, weight, sitDays).
		First(&costs)
	if err != nil {
		return nil, fmt.Errorf("error calculating PPM SIT costs for PPM ID %s: %w", ppmID, err)
	}

	return &costs, nil
}
