package models

import (
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/unit"
)

// MTOShipmentType represents the type of shipments the mto shipment is
type MTOShipmentType string

// using these also in move.go selected move type
const (
	// NTSRaw is the raw string value of the NTS Shipment Type
	NTSRaw = "HHG_INTO_NTS"
	// NTSrRaw is the raw string value of the NTSr Shipment Type
	NTSrRaw = "HHG_OUTOF_NTS"
)

// Market code indicator of international or domestic
type MarketCode string

const (
	MarketCodeDomestic      MarketCode = "d" // domestic
	MarketCodeInternational MarketCode = "i" // international
)

// Add to this list as international service items are implemented
var internationalAccessorialServiceItems = []ReServiceCode{
	ReServiceCodeICRT,
	ReServiceCodeIUCRT,
	ReServiceCodeIOASIT,
	ReServiceCodeIDASIT,
	ReServiceCodeIOFSIT,
	ReServiceCodeIDFSIT,
	ReServiceCodeIOPSIT,
	ReServiceCodeIDDSIT,
	ReServiceCodeIDSHUT,
	ReServiceCodeIOSHUT,
	ReServiceCodeIOSFSC,
	ReServiceCodeIDSFSC,
}

const (
	// MTOShipmentTypeHHG is an HHG Shipment Type default
	MTOShipmentTypeHHG MTOShipmentType = "HHG"
	// MTOShipmentTypeHHGIntoNTS is an HHG Shipment Type for going into NTS
	MTOShipmentTypeHHGIntoNTS MTOShipmentType = NTSRaw
	// MTOShipmentTypeHHGOutOfNTS is an HHG Shipment Type for going out of NTS
	MTOShipmentTypeHHGOutOfNTS MTOShipmentType = NTSrRaw
	// MTOShipmentTypeMobileHome is a Shipment Type for MobileHome
	MTOShipmentTypeMobileHome MTOShipmentType = "MOBILE_HOME"
	// MTOShipmentTypeBoatHaulAway is a Shipment Type for Boat Haul Away
	MTOShipmentTypeBoatHaulAway MTOShipmentType = "BOAT_HAUL_AWAY"
	// MTOShipmentTypeBoatTowAway is a Shipment Type for Boat Tow Away
	MTOShipmentTypeBoatTowAway MTOShipmentType = "BOAT_TOW_AWAY"
	// MTOShipmentTypePPM is a Shipment Type for Personally Procured Move shipments
	MTOShipmentTypePPM MTOShipmentType = "PPM"
	// MTOShipmentTypeUB is a Shipment Type for Unaccompanied Baggage shipments
	MTOShipmentTypeUnaccompaniedBaggage MTOShipmentType = "UNACCOMPANIED_BAGGAGE"
)

// These are meant to be the default number of SIT days that a customer is allowed to have. They should be used when
// creating a shipment and setting the initial value. Other values will likely be added to this once we deal with
// different types of customers.
const (
	// DefaultServiceMemberSITDaysAllowance is the default number of SIT days a service member is allowed
	DefaultServiceMemberSITDaysAllowance = 90
)

// MTOShipmentStatus represents the possible statuses for a mto shipment
type MTOShipmentStatus string

const (
	// MTOShipmentStatusDraft is the draft status type for MTO Shipments
	MTOShipmentStatusDraft MTOShipmentStatus = "DRAFT"
	// MTOShipmentStatusSubmitted is the submitted status type for MTO Shipments
	MTOShipmentStatusSubmitted MTOShipmentStatus = "SUBMITTED"
	// MTOShipmentStatusApproved is the approved status type for MTO Shipments
	MTOShipmentStatusApproved MTOShipmentStatus = "APPROVED"
	// MTOShipmentStatusRejected is the rejected status type for MTO Shipments
	MTOShipmentStatusRejected MTOShipmentStatus = "REJECTED"
	// MTOShipmentStatusCancellationRequested indicates the TOO has requested that the Prime cancel the shipment
	MTOShipmentStatusCancellationRequested MTOShipmentStatus = "CANCELLATION_REQUESTED"
	// MTOShipmentStatusCanceled indicates that a shipment has been canceled by the Prime
	MTOShipmentStatusCanceled MTOShipmentStatus = "CANCELED"
	// MTOShipmentStatusDiversionRequested indicates that the TOO has requested that the Prime divert a shipment
	MTOShipmentStatusDiversionRequested MTOShipmentStatus = "DIVERSION_REQUESTED"
	// MTOShipmentTerminatedForCause indicates that a shipment has been terminated for cause by a COR
	MTOShipmentStatusTerminatedForCause MTOShipmentStatus = "TERMINATED_FOR_CAUSE"
	// MoveStatusAPPROVALSREQUESTED is the approvals requested status type for MTO Shipments
	MTOShipmentStatusApprovalsRequested MTOShipmentStatus = "APPROVALS_REQUESTED"
)

// LOAType represents the possible TAC and SAC types for a mto shipment
type LOAType string

const (
	// LOATypeHHG is the HHG TAC or SAC
	LOATypeHHG LOAType = "HHG"
	// LOATypeNTS is the NTS TAC or SAC
	LOATypeNTS LOAType = "NTS"
)

type DestinationType string

const (
	DestinationTypeHomeOfRecord           DestinationType = "HOME_OF_RECORD"
	DestinationTypeHomeOfSelection        DestinationType = "HOME_OF_SELECTION"
	DestinationTypePlaceEnteredActiveDuty DestinationType = "PLACE_ENTERED_ACTIVE_DUTY"
	DestinationTypeOtherThanAuthorized    DestinationType = "OTHER_THAN_AUTHORIZED"
)

// MTOShipment is an object representing data for a move task order shipment
type MTOShipment struct {
	ID                               uuid.UUID              `json:"id" db:"id"`
	MoveTaskOrder                    Move                   `json:"move_task_order" belongs_to:"moves" fk_id:"move_id"`
	MoveTaskOrderID                  uuid.UUID              `json:"move_task_order_id" db:"move_id"`
	ScheduledPickupDate              *time.Time             `json:"scheduled_pickup_date" db:"scheduled_pickup_date"`
	RequestedPickupDate              *time.Time             `json:"requested_pickup_date" db:"requested_pickup_date"`
	RequestedDeliveryDate            *time.Time             `json:"requested_delivery_date" db:"requested_delivery_date"`
	ApprovedDate                     *time.Time             `json:"approved_date" db:"approved_date"`
	FirstAvailableDeliveryDate       *time.Time             `json:"first_available_delivery_date" db:"first_available_delivery_date"`
	ActualPickupDate                 *time.Time             `json:"actual_pickup_date" db:"actual_pickup_date"`
	RequiredDeliveryDate             *time.Time             `json:"required_delivery_date" db:"required_delivery_date"`
	ScheduledDeliveryDate            *time.Time             `json:"scheduled_delivery_date" db:"scheduled_delivery_date"`
	ActualDeliveryDate               *time.Time             `json:"actual_delivery_date" db:"actual_delivery_date"`
	CustomerRemarks                  *string                `json:"customer_remarks" db:"customer_remarks"`
	CounselorRemarks                 *string                `json:"counselor_remarks" db:"counselor_remarks"`
	PickupAddress                    *Address               `json:"pickup_address" belongs_to:"addresses" fk_id:"pickup_address_id"`
	PickupAddressID                  *uuid.UUID             `json:"pickup_address_id" db:"pickup_address_id"`
	DestinationAddress               *Address               `json:"destination_address" belongs_to:"addresses" fk_id:"destination_address_id"`
	DestinationAddressID             *uuid.UUID             `json:"destination_address_id" db:"destination_address_id"`
	DestinationType                  *DestinationType       `json:"destination_type" db:"destination_address_type"`
	MTOAgents                        MTOAgents              `json:"mto_agents" has_many:"mto_agents" fk_id:"mto_shipment_id"`
	MTOServiceItems                  MTOServiceItems        `json:"mto_service_items" has_many:"mto_service_items" fk_id:"mto_shipment_id"`
	SecondaryPickupAddress           *Address               `json:"secondary_pickup_address" belongs_to:"addresses" fk_id:"secondary_pickup_address_id"`
	SecondaryPickupAddressID         *uuid.UUID             `json:"secondary_pickup_address_id" db:"secondary_pickup_address_id"`
	HasSecondaryPickupAddress        *bool                  `json:"has_secondary_pickup_address" db:"has_secondary_pickup_address"`
	SecondaryDeliveryAddress         *Address               `json:"secondary_delivery_address" belongs_to:"addresses" fk_id:"secondary_delivery_address_id"`
	SecondaryDeliveryAddressID       *uuid.UUID             `json:"secondary_delivery_address_id" db:"secondary_delivery_address_id"`
	HasSecondaryDeliveryAddress      *bool                  `json:"has_secondary_delivery_address" db:"has_secondary_delivery_address"`
	TertiaryPickupAddress            *Address               `json:"tertiary_pickup_address" belongs_to:"addresses" fk_id:"tertiary_pickup_address_id"`
	TertiaryPickupAddressID          *uuid.UUID             `json:"tertiary_pickup_address_id" db:"tertiary_pickup_address_id"`
	HasTertiaryPickupAddress         *bool                  `json:"has_tertiary_pickup_address" db:"has_tertiary_pickup_address"`
	TertiaryDeliveryAddress          *Address               `json:"tertiary_delivery_address" belongs_to:"addresses" fk_id:"tertiary_delivery_address_id"`
	TertiaryDeliveryAddressID        *uuid.UUID             `json:"tertiary_delivery_address_id" db:"tertiary_delivery_address_id"`
	HasTertiaryDeliveryAddress       *bool                  `json:"has_tertiary_delivery_address" db:"has_tertiary_delivery_address"`
	SITDaysAllowance                 *int                   `json:"sit_days_allowance" db:"sit_days_allowance"`
	SITDurationUpdates               SITDurationUpdates     `json:"sit_duration_updates" has_many:"sit_extensions" fk_id:"mto_shipment_id"`
	PrimeEstimatedWeight             *unit.Pound            `json:"prime_estimated_weight" db:"prime_estimated_weight"`
	PrimeEstimatedWeightRecordedDate *time.Time             `json:"prime_estimated_weight_recorded_date" db:"prime_estimated_weight_recorded_date"`
	PrimeActualWeight                *unit.Pound            `json:"prime_actual_weight" db:"prime_actual_weight"`
	BillableWeightCap                *unit.Pound            `json:"billable_weight_cap" db:"billable_weight_cap"`
	BillableWeightJustification      *string                `json:"billable_weight_justification" db:"billable_weight_justification"`
	NTSRecordedWeight                *unit.Pound            `json:"nts_recorded_weight" db:"nts_recorded_weight"`
	ShipmentType                     MTOShipmentType        `json:"shipment_type" db:"shipment_type"`
	Status                           MTOShipmentStatus      `json:"status" db:"status"`
	Diversion                        bool                   `json:"diversion" db:"diversion"`
	DiversionReason                  *string                `json:"diversion_reason" db:"diversion_reason"`
	DivertedFromShipmentID           *uuid.UUID             `json:"diverted_from_shipment_id" db:"diverted_from_shipment_id"`
	ActualProGearWeight              *unit.Pound            `json:"actual_pro_gear_weight" db:"actual_pro_gear_weight"`
	ActualSpouseProGearWeight        *unit.Pound            `json:"actual_spouse_pro_gear_weight" db:"actual_spouse_pro_gear_weight"`
	RejectionReason                  *string                `json:"rejection_reason" db:"rejection_reason"`
	Distance                         *unit.Miles            `json:"distance" db:"distance"`
	Reweigh                          *Reweigh               `json:"reweigh" has_one:"reweighs" fk_id:"shipment_id"`
	UsesExternalVendor               bool                   `json:"uses_external_vendor" db:"uses_external_vendor"`
	StorageFacility                  *StorageFacility       `json:"storage_facility" belongs_to:"storage_facilities" fk:"storage_facility_id"`
	StorageFacilityID                *uuid.UUID             `json:"storage_facility_id" db:"storage_facility_id"`
	ServiceOrderNumber               *string                `json:"service_order_number" db:"service_order_number"`
	TACType                          *LOAType               `json:"tac_type" db:"tac_type"`
	SACType                          *LOAType               `json:"sac_type" db:"sac_type"`
	PPMShipment                      *PPMShipment           `json:"ppm_shipment" has_one:"ppm_shipment" fk_id:"shipment_id"`
	BoatShipment                     *BoatShipment          `json:"boat_shipment" has_one:"boat_shipment" fk_id:"shipment_id"`
	DeliveryAddressUpdate            *ShipmentAddressUpdate `json:"delivery_address_update" has_one:"shipment_address_update" fk_id:"shipment_id"`
	CreatedAt                        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt                        time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt                        *time.Time             `json:"deleted_at" db:"deleted_at"`
	ShipmentLocator                  *string                `json:"shipment_locator" db:"shipment_locator"`
	OriginSITAuthEndDate             *time.Time             `json:"origin_sit_auth_end_date" db:"origin_sit_auth_end_date"`
	DestinationSITAuthEndDate        *time.Time             `json:"destination_sit_auth_end_date" db:"dest_sit_auth_end_date"`
	MobileHome                       *MobileHome            `json:"mobile_home" has_one:"mobile_home" fk_id:"shipment_id"`
	MarketCode                       MarketCode             `json:"market_code" db:"market_code"`
	PrimeAcknowledgedAt              *time.Time             `db:"prime_acknowledged_at"`
	TerminationComments              *string                `json:"termination_comments" db:"termination_comments"`
	TerminatedAt                     *time.Time             `json:"terminated_at" db:"terminated_at"`
}

// TableName overrides the table name used by Pop.
func (m MTOShipment) TableName() string {
	return "mto_shipments"
}

// MTOShipments is a list of mto shipments
type MTOShipments []MTOShipment

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MTOShipment) Validate(db *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	var customVerrs []validate.Errors
	vs = append(vs, &validators.StringInclusion{Field: string(m.Status), Name: "Status", List: []string{
		string(MTOShipmentStatusApproved),
		string(MTOShipmentStatusRejected),
		string(MTOShipmentStatusSubmitted),
		string(MTOShipmentStatusDraft),
		string(MTOShipmentStatusCancellationRequested),
		string(MTOShipmentStatusCanceled),
		string(MTOShipmentStatusDiversionRequested),
		string(MTOShipmentStatusTerminatedForCause),
		string(MTOShipmentStatusApprovalsRequested),
	}})
	// Check if the status of the original shipment is terminated
	if m.ID != uuid.Nil && db != nil {
		var existingShipment MTOShipment
		err := db.Find(&existingShipment, m.ID)
		if err == nil && existingShipment.Status == MTOShipmentStatusTerminatedForCause {
			terminationVerr := validate.NewErrors()
			terminationVerr.Add("status", "Cannot update shipment with status TERMINATED_FOR_CAUSE")
			customVerrs = append(customVerrs, *terminationVerr)
		}
	}
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MoveTaskOrderID, Name: "MoveTaskOrderID"})
	if m.PrimeEstimatedWeight != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: m.PrimeEstimatedWeight.Int(), Compared: 0, Name: "PrimeEstimatedWeight"})
	}
	if m.PrimeActualWeight != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: m.PrimeActualWeight.Int(), Compared: 0, Name: "PrimeActualWeight"})
	}
	if m.NTSRecordedWeight != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: m.NTSRecordedWeight.Int(), Compared: -1, Name: "NTSRecordedWeight"})
	}
	vs = append(vs, &OptionalPoundIsNonNegative{Field: m.BillableWeightCap, Name: "BillableWeightCap"})
	vs = append(vs, &StringIsNilOrNotBlank{Field: m.BillableWeightJustification, Name: "BillableWeightJustification"})
	if m.Status == MTOShipmentStatusRejected {
		var rejectionReason string
		if m.RejectionReason != nil {
			rejectionReason = *m.RejectionReason
		}
		vs = append(vs, &validators.StringIsPresent{Field: rejectionReason, Name: "RejectionReason"})
	}
	if m.SITDaysAllowance != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: *m.SITDaysAllowance, Compared: -1, Name: "SITDaysAllowance"})
	}
	vs = append(vs, &OptionalUUIDIsPresent{Field: m.StorageFacilityID, Name: "StorageFacilityID"})
	vs = append(vs, &StringIsNilOrNotBlank{Field: m.ServiceOrderNumber, Name: "ServiceOrderNumber"})

	var ptrTACType *string
	if m.TACType != nil {
		tacType := string(*m.TACType)
		ptrTACType = &tacType
	}
	vs = append(vs, &OptionalStringInclusion{Field: ptrTACType, Name: "TACType", List: []string{
		string(LOATypeHHG),
		string(LOATypeNTS),
	}})

	var ptrSACType *string
	if m.SACType != nil {
		sacType := string(*m.SACType)
		ptrSACType = &sacType
	}
	vs = append(vs, &OptionalStringInclusion{Field: ptrSACType, Name: "SACType", List: []string{
		string(LOATypeHHG),
		string(LOATypeNTS),
	}})

	var destinationType *string
	if m.DestinationType != nil {
		valDestinationType := string(*m.DestinationType)
		destinationType = &valDestinationType
	}
	vs = append(vs, &OptionalStringInclusion{Field: destinationType, Name: "DestinationType", List: []string{
		string(DestinationTypeHomeOfRecord),
		string(DestinationTypeHomeOfSelection),
		string(DestinationTypePlaceEnteredActiveDuty),
		string(DestinationTypeOtherThanAuthorized),
	}})

	if m.MarketCode != "" {
		vs = append(vs, &validators.StringInclusion{
			Field: string(m.MarketCode),
			Name:  "MarketCode",
			List: []string{
				string(MarketCodeDomestic),
				string(MarketCodeInternational),
			},
		})
	}

	verrs := validate.Validate(vs...)
	// Add our custom verrs the ole manual way because the types
	// didn't want to append together
	for _, e := range customVerrs {
		for field, msgs := range e.Errors {
			for _, msg := range msgs {
				verrs.Add(field, msg)
			}
		}
	}
	return verrs, nil
}

// GetCustomerFromShipment gets the service member given a shipment id
func GetCustomerFromShipment(db *pop.Connection, shipmentID uuid.UUID) (*ServiceMember, error) {
	var serviceMember ServiceMember
	err := db.Q().
		InnerJoin("orders", "orders.service_member_id = service_members.id").
		InnerJoin("moves", "moves.orders_id = orders.id").
		InnerJoin("mto_shipments", "mto_shipments.move_id = moves.id").
		Where("mto_shipments.id = ?", shipmentID).
		First(&serviceMember)
	if err != nil {
		return &serviceMember, fmt.Errorf("error fetching service member for shipment ID: %s with error %w", shipmentID, err)
	}
	return &serviceMember, nil
}

// Helper function to check that an MTO Shipment contains a PPM Shipment
func (m MTOShipment) ContainsAPPMShipment() bool {
	return m.PPMShipment != nil
}

func (m MTOShipment) IsPPMShipment() bool {
	return m.ShipmentType == MTOShipmentTypePPM
}

// determining the market code for a shipment based off of address isOconus value
// this function takes in a shipment and returns the same shipment with the updated MarketCode value
func DetermineShipmentMarketCode(shipment *MTOShipment) *MTOShipment {
	// helper to check if both addresses are CONUS
	isDomestic := func(pickupAddress, destAddress *Address) bool {
		return pickupAddress != nil && destAddress != nil &&
			pickupAddress.IsOconus != nil && destAddress.IsOconus != nil &&
			!*pickupAddress.IsOconus && !*destAddress.IsOconus
	}

	// determine market code based on address and shipment type
	switch shipment.ShipmentType {
	case MTOShipmentTypeHHGIntoNTS:
		if shipment.PickupAddress != nil && shipment.StorageFacility != nil &&
			shipment.PickupAddress.IsOconus != nil && shipment.StorageFacility.Address.IsOconus != nil {
			// If both pickup and storage facility are present, check if both are domestic
			if isDomestic(shipment.PickupAddress, &shipment.StorageFacility.Address) {
				shipment.MarketCode = MarketCodeDomestic
			} else {
				shipment.MarketCode = MarketCodeInternational
			}
		} else if shipment.PickupAddress != nil && shipment.PickupAddress.IsOconus != nil {
			// customers only submit pickup addresses on shipment creation
			if !*shipment.PickupAddress.IsOconus {
				shipment.MarketCode = MarketCodeDomestic
			} else {
				shipment.MarketCode = MarketCodeInternational
			}
		}
	case MTOShipmentTypeHHGOutOfNTS:
		if shipment.StorageFacility != nil && shipment.DestinationAddress != nil &&
			shipment.StorageFacility.Address.IsOconus != nil && shipment.DestinationAddress.IsOconus != nil {
			if isDomestic(&shipment.StorageFacility.Address, shipment.DestinationAddress) {
				shipment.MarketCode = MarketCodeDomestic
			} else {
				shipment.MarketCode = MarketCodeInternational
			}
		} else if shipment.DestinationAddress != nil && shipment.DestinationAddress.IsOconus != nil {
			// customers only submit destination addresses on NTS-release shipments
			if !*shipment.DestinationAddress.IsOconus {
				shipment.MarketCode = MarketCodeDomestic
			} else {
				shipment.MarketCode = MarketCodeInternational
			}
		}
	default:
		if shipment.PickupAddress != nil && shipment.DestinationAddress != nil &&
			shipment.PickupAddress.IsOconus != nil && shipment.DestinationAddress.IsOconus != nil {
			if isDomestic(shipment.PickupAddress, shipment.DestinationAddress) {
				shipment.MarketCode = MarketCodeDomestic
			} else {
				shipment.MarketCode = MarketCodeInternational
			}
		} else {
			// set a default market code for cases where PPM logic needs to be done after shipment creation
			shipment.MarketCode = MarketCodeDomestic
		}
	}
	return shipment
}

func (s MTOShipment) GetDestinationAddress(db *pop.Connection) (*Address, error) {
	if uuid.UUID.IsNil(s.ID) {
		return nil, errors.New("MTOShipment ID is required to fetch destination address.")
	}

	err := db.Load(&s, "DestinationAddress", "PPMShipment.DestinationAddress")
	if err != nil {
		if err.Error() == RecordNotFoundErrorString {
			return nil, errors.WithMessage(ErrSqlRecordNotFound, string(s.ShipmentType)+" ShipmentID: "+s.ID.String())
		}
		return nil, err
	}

	if s.ShipmentType == MTOShipmentTypePPM {
		if s.PPMShipment.DestinationAddress != nil {
			return s.PPMShipment.DestinationAddress, nil
		} else if s.DestinationAddress != nil {
			return s.DestinationAddress, nil
		}
		return nil, errors.WithMessage(ErrMissingDestinationAddress, string(s.ShipmentType))
	}

	if s.DestinationAddress != nil {
		return s.DestinationAddress, nil
	}

	return nil, errors.WithMessage(ErrMissingDestinationAddress, string(s.ShipmentType))
}

// this function takes in two addresses and determines the market code string
func DetermineMarketCode(address1 *Address, address2 *Address) (MarketCode, error) {
	if address1 == nil || address2 == nil {
		return "", fmt.Errorf("both address1 and address2 must be provided")
	}

	// helper to check if both addresses are CONUS
	isDomestic := func(a, b *Address) bool {
		return a != nil && b != nil &&
			a.IsOconus != nil && b.IsOconus != nil &&
			!*a.IsOconus && !*b.IsOconus
	}

	if isDomestic(address1, address2) {
		return MarketCodeDomestic, nil
	} else {
		return MarketCodeInternational, nil
	}
}

// PortLocationInfo holds the ZIP code and port type for a shipment
// this is used in the db function/query below
type PortLocationInfo struct {
	UsprZipID string `db:"uspr_zip_id"`
	PortType  string `db:"port_type"`
}

// GetPortLocationForShipment gets the ZIP and port type associated with the port for the POEFSC/PODFSC service item in a shipment
func GetPortLocationInfoForShipment(db *pop.Connection, shipmentID uuid.UUID) (*string, *string, error) {
	var portLocationInfo PortLocationInfo

	err := db.RawQuery("SELECT * FROM get_port_location_info_for_shipment($1)", shipmentID).
		First(&portLocationInfo)

	if err != nil && err != sql.ErrNoRows {
		return nil, nil, fmt.Errorf("error fetching port location for shipment ID: %s with error %w", shipmentID, err)
	}

	// return the ZIP code and port type, or nil if not found
	if portLocationInfo.UsprZipID != "" && portLocationInfo.PortType != "" {
		return &portLocationInfo.UsprZipID, &portLocationInfo.PortType, nil
	}

	// if nothing was found, return nil - just means we don't have the port info from Prime yet
	return nil, nil, nil
}

func CreateApprovedServiceItemsForShipment(db *pop.Connection, shipment *MTOShipment) error {
	err := db.RawQuery("CALL create_approved_service_items_for_shipment($1)", shipment.ID).Exec()
	if err != nil {
		return fmt.Errorf("error creating approved service items: %w", err)
	}

	return nil
}

func CreateInternationalAccessorialServiceItemsForShipment(db *pop.Connection, shipmentId uuid.UUID, mtoServiceItems MTOServiceItems) ([]string, error) {
	if len(mtoServiceItems) == 0 {
		err := fmt.Errorf("must request service items to create: %s", shipmentId)
		return nil, apperror.NewInvalidInputError(shipmentId, err, nil, err.Error())
	}

	for _, serviceItem := range mtoServiceItems {
		if !slices.Contains(internationalAccessorialServiceItems, serviceItem.ReService.Code) {
			err := fmt.Errorf("cannot create domestic service items for international shipment: %s", shipmentId)
			return nil, apperror.NewInvalidInputError(shipmentId, err, nil, err.Error())
		}
	}

	createdServiceItemIDs := []string{}
	err := db.RawQuery("CALL create_accessorial_service_items_for_shipment($1, $2, $3)", shipmentId, pq.Array(mtoServiceItems), pq.StringArray(createdServiceItemIDs)).All(&createdServiceItemIDs)
	if err != nil {
		return nil, apperror.NewInvalidInputError(shipmentId, err, nil, err.Error())
	}

	return createdServiceItemIDs, nil
}

// a db stored proc that will handle updating the pricing_estimate columns of basic service items for shipment types:
// iHHG
// iUB
func UpdateEstimatedPricingForShipmentBasicServiceItems(db *pop.Connection, shipment *MTOShipment, mileage *int) error {
	err := db.RawQuery("CALL update_service_item_pricing($1, $2)", shipment.ID, mileage).Exec()
	if err != nil {
		return fmt.Errorf("error updating estimated pricing for shipment's service items: %w", err)
	}

	return nil
}

// GetDestinationGblocForShipment gets the GBLOC associated with the shipment's destination address
// there are certain exceptions for OCONUS addresses in Alaska Zone II based on affiliation
func GetDestinationGblocForShipment(db *pop.Connection, shipmentID uuid.UUID) (*string, error) {
	var gbloc *string

	err := db.RawQuery("SELECT * FROM get_destination_gbloc_for_shipment($1)", shipmentID).
		First(&gbloc)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error fetching destination gbloc for shipment ID: %s with error %w", shipmentID, err)
	}

	if gbloc != nil {
		return gbloc, nil
	}

	return nil, nil
}

// Returns a Shipment for a given id
func FetchShipmentByID(db *pop.Connection, shipmentID uuid.UUID) (*MTOShipment, error) {
	var mtoShipment MTOShipment
	err := db.Q().Find(&mtoShipment, shipmentID)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		return nil, err
	}
	return &mtoShipment, nil
}

// filters the returned MtoShipments for each move.
// Ignoring mto shipments that have been deleted, cancelled, or rejected.
func FilterDeletedRejectedCanceledMtoShipments(unfilteredShipments MTOShipments) MTOShipments {
	if len(unfilteredShipments) == 0 {
		return unfilteredShipments
	}

	filteredShipments := MTOShipments{}
	for _, shipment := range unfilteredShipments {
		if shipment.DeletedAt == nil &&
			(shipment.Status != MTOShipmentStatusDraft) &&
			(shipment.Status != MTOShipmentStatusRejected) &&
			(shipment.Status != MTOShipmentStatusCanceled) {
			filteredShipments = append(filteredShipments, shipment)
		}
	}

	return filteredShipments
}

// returns a pointer to a bool indicating whether the shipment is OCONUS
// returns nil if either PickupAddress.IsOconus or DestinationAddress.IsOconus is nil
func IsShipmentOCONUS(shipment MTOShipment) *bool {
	if shipment.PickupAddress == nil || shipment.DestinationAddress == nil {
		return nil
	}
	if shipment.PickupAddress.IsOconus == nil || shipment.DestinationAddress.IsOconus == nil {
		return nil
	}

	isOCONUS := *shipment.PickupAddress.IsOconus || *shipment.DestinationAddress.IsOconus
	return &isOCONUS
}

func (m *MTOShipment) CanSendReweighEmailForShipmentType() bool {
	return m.ShipmentType != MTOShipmentTypePPM
}

func PrimeCanUpdateDeliveryAddress(shipmentType MTOShipmentType) bool {
	isValid := false
	if shipmentType != "" && shipmentType != MTOShipmentTypePPM && shipmentType != MTOShipmentTypeHHGIntoNTS {
		isValid = true
	}

	return isValid
}

func IsShipmentApprovable(dbShipment MTOShipment) bool {
	// check if any service items on current shipment still need to be reviewed
	if dbShipment.MTOServiceItems != nil {
		for _, serviceItem := range dbShipment.MTOServiceItems {
			if serviceItem.Status == MTOServiceItemStatusSubmitted {
				return false
			}
		}
	}
	// check if all SIT Extensions are reviewed
	if dbShipment.SITDurationUpdates != nil {
		for _, sitDurationUpdate := range dbShipment.SITDurationUpdates {
			if sitDurationUpdate.Status == SITExtensionStatusPending {
				return false
			}
		}
	}
	// check if all Delivery Address updates are reviewed
	if dbShipment.DeliveryAddressUpdate != nil && dbShipment.DeliveryAddressUpdate.Status == ShipmentAddressUpdateStatusRequested {
		return false
	}

	return true
}
