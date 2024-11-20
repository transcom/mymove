package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/unit"
)

// MTOServiceItemStatus represents the possible statuses for a mto shipment
type MTOServiceItemStatus string

const (
	// MTOServiceItemStatusSubmitted is the submitted status type for MTO Service Items
	MTOServiceItemStatusSubmitted MTOServiceItemStatus = "SUBMITTED"
	// MTOServiceItemStatusApproved is the approved status type for MTO Service Items
	MTOServiceItemStatusApproved MTOServiceItemStatus = "APPROVED"
	// MTOServiceItemStatusRejected is the rejected status type for MTO Service Items
	MTOServiceItemStatusRejected MTOServiceItemStatus = "REJECTED"
)

// MTOServiceItem is an object representing service items for a move task order.
type MTOServiceItem struct {
	ID                                uuid.UUID                      `db:"id"`
	MoveTaskOrder                     Move                           `belongs_to:"moves" fk_id:"move_id"`
	MoveTaskOrderID                   uuid.UUID                      `db:"move_id"`
	MTOShipment                       MTOShipment                    `belongs_to:"mto_shipments" fk_id:"mto_shipment_id"`
	MTOShipmentID                     *uuid.UUID                     `db:"mto_shipment_id"`
	ReService                         ReService                      `belongs_to:"re_services" fk_id:"re_service_id"`
	ReServiceID                       uuid.UUID                      `db:"re_service_id"`
	Reason                            *string                        `db:"reason"`
	RejectionReason                   *string                        `db:"rejection_reason"`
	Status                            MTOServiceItemStatus           `db:"status"`
	PickupPostalCode                  *string                        `db:"pickup_postal_code"`
	SITPostalCode                     *string                        `db:"sit_postal_code"`
	SITEntryDate                      *time.Time                     `db:"sit_entry_date"`
	SITDepartureDate                  *time.Time                     `db:"sit_departure_date"`
	SITCustomerContacted              *time.Time                     `db:"sit_customer_contacted"`
	SITRequestedDelivery              *time.Time                     `db:"sit_requested_delivery"`
	SITOriginHHGOriginalAddress       *Address                       `belongs_to:"addresses" fk_id:"sit_origin_hhg_original_address_id"`
	SITOriginHHGOriginalAddressID     *uuid.UUID                     `db:"sit_origin_hhg_original_address_id"`
	SITOriginHHGActualAddress         *Address                       `belongs_to:"addresses" fk_id:"sit_origin_hhg_actual_address_id"`
	SITOriginHHGActualAddressID       *uuid.UUID                     `db:"sit_origin_hhg_actual_address_id"`
	SITDestinationOriginalAddress     *Address                       `belongs_to:"addresses" fk_id:"sit_destination_original_address_id"`
	SITDestinationOriginalAddressID   *uuid.UUID                     `db:"sit_destination_original_address_id"`
	SITDestinationFinalAddress        *Address                       `belongs_to:"addresses" fk_id:"sit_destination_final_address_id"`
	SITDestinationFinalAddressID      *uuid.UUID                     `db:"sit_destination_final_address_id"`
	Description                       *string                        `db:"description"`
	EstimatedWeight                   *unit.Pound                    `db:"estimated_weight"`
	ActualWeight                      *unit.Pound                    `db:"actual_weight"`
	Dimensions                        MTOServiceItemDimensions       `has_many:"mto_service_item_dimensions" fk_id:"mto_service_item_id"`
	CustomerContacts                  MTOServiceItemCustomerContacts `many_to_many:"service_items_customer_contacts"`
	ServiceRequestDocuments           ServiceRequestDocuments        `has_many:"service_request_document" fk_id:"mto_service_item_id"`
	CreatedAt                         time.Time                      `db:"created_at"`
	UpdatedAt                         time.Time                      `db:"updated_at"`
	ApprovedAt                        *time.Time                     `db:"approved_at"`
	RejectedAt                        *time.Time                     `db:"rejected_at"`
	RequestedApprovalsRequestedStatus *bool                          `db:"requested_approvals_requested_status"`
	CustomerExpense                   bool                           `db:"customer_expense"`
	CustomerExpenseReason             *string                        `db:"customer_expense_reason"`
	SITDeliveryMiles                  *int                           `db:"sit_delivery_miles"`
	PricingEstimate                   *unit.Cents                    `db:"pricing_estimate"`
	StandaloneCrate                   *bool                          `db:"standalone_crate"`
	LockedPriceCents                  *unit.Cents                    `db:"locked_price_cents"`
	POELocation                       *PortLocation                  `belongs_to:"port_locations" fk_id:"poe_location_id"`
	POELocationID                     *uuid.UUID                     `db:"poe_location_id"`
	PODLocation                       *PortLocation                  `belongs_to:"port_locations" fk_id:"pod_location_id"`
	PODLocationID                     *uuid.UUID                     `db:"pod_location_id"`
	ServiceLocation                   *ServiceLocationType           `db:"service_location"`
}

// MTOServiceItemSingle is an object representing a single column in the service items table
type MTOServiceItemSingle struct {
	ID                              uuid.UUID            `db:"id"`
	MoveTaskOrderID                 uuid.UUID            `db:"move_id"`
	MTOShipmentID                   *uuid.UUID           `db:"mto_shipment_id"`
	ReServiceID                     uuid.UUID            `db:"re_service_id"`
	CreatedAt                       time.Time            `db:"created_at"`
	UpdatedAt                       time.Time            `db:"updated_at"`
	Reason                          *string              `db:"reason"`
	PickupPostalCode                *string              `db:"pickup_postal_code"`
	Description                     *string              `db:"description"`
	Status                          MTOServiceItemStatus `db:"status"`
	RejectionReason                 *string              `db:"rejection_reason"`
	ApprovedAt                      *time.Time           `db:"approved_at"`
	RejectedAt                      *time.Time           `db:"rejected_at"`
	SITPostalCode                   *string              `db:"sit_postal_code"`
	SITEntryDate                    *time.Time           `db:"sit_entry_date"`
	SITDepartureDate                *time.Time           `db:"sit_departure_date"`
	SITDestinationFinalAddressID    *uuid.UUID           `db:"sit_destination_final_address_id"`
	SITOriginHHGOriginalAddressID   *uuid.UUID           `db:"sit_origin_hhg_original_address_id"`
	SITOriginHHGActualAddressID     *uuid.UUID           `db:"sit_origin_hhg_actual_address_id"`
	EstimatedWeight                 *unit.Pound          `db:"estimated_weight"`
	ActualWeight                    *unit.Pound          `db:"actual_weight"`
	SITDestinationOriginalAddressID *uuid.UUID           `db:"sit_destination_original_address_id"`
	SITCustomerContacted            *time.Time           `db:"sit_customer_contacted"`
	SITRequestedDelivery            *time.Time           `db:"sit_requested_delivery"`
	CustomerExpense                 bool                 `db:"customer_expense"`
	CustomerExpenseReason           *string              `db:"customer_expense_reason"`
	SITDeliveryMiles                *unit.Miles          `db:"sit_delivery_miles"`
	PricingEstimate                 *unit.Cents          `db:"pricing_estimate"`
	POELocationID                   *uuid.UUID           `db:"poe_location_id"`
	PODLocationID                   *uuid.UUID           `db:"pod_location_id"`
}

// TableName overrides the table name used by Pop.
func (m MTOServiceItem) TableName() string {
	return "mto_service_items"
}

// MTOServiceItems is a slice containing MTOServiceItems
type MTOServiceItems []MTOServiceItem

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MTOServiceItem) Validate(_ *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.StringInclusion{Field: string(m.Status), Name: "Status", List: []string{
		string(MTOServiceItemStatusSubmitted),
		string(MTOServiceItemStatusApproved),
		string(MTOServiceItemStatusRejected),
	}})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MoveTaskOrderID, Name: "MoveTaskOrderID"})
	vs = append(vs, &OptionalUUIDIsPresent{Field: m.MTOShipmentID, Name: "MTOShipmentID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.ReServiceID, Name: "ReServiceID"})
	vs = append(vs, &StringIsNilOrNotBlank{Field: m.Reason, Name: "Reason"})
	vs = append(vs, &StringIsNilOrNotBlank{Field: m.PickupPostalCode, Name: "PickupPostalCode"})
	vs = append(vs, &StringIsNilOrNotBlank{Field: m.Description, Name: "Description"})

	return validate.Validate(vs...), nil
}

// FetchRelatedDestinationSITServiceItems returns all service items with destination SIT ReService codes
// that are associated with the same shipment as the provided service item.
func FetchRelatedDestinationSITServiceItems(tx *pop.Connection, mtoServiceItemID uuid.UUID) (MTOServiceItems, error) {
	var relatedDestinationSITServiceItems MTOServiceItems
	err := tx.RawQuery(
		`SELECT msi.id
			FROM mto_service_items msi
			INNER JOIN re_services res ON msi.re_service_id = res.id
			WHERE res.code IN (?, ?, ?, ?) AND mto_shipment_id IN (
				SELECT mto_shipment_id FROM mto_service_items WHERE id = ?)`, ReServiceCodeDDFSIT, ReServiceCodeDDASIT, ReServiceCodeDDDSIT, ReServiceCodeDDSFSC, mtoServiceItemID).
		All(&relatedDestinationSITServiceItems)
	return relatedDestinationSITServiceItems, err
}

func FetchServiceItem(db *pop.Connection, serviceItemID uuid.UUID) (MTOServiceItem, error) {
	var serviceItem MTOServiceItem
	err := db.Eager("SITDestinationOriginalAddress",
		"SITDestinationFinalAddress",
		"ReService",
		"CustomerContacts").Where("id = ?", serviceItemID).First(&serviceItem)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return MTOServiceItem{}, ErrFetchNotFound
		}
		return MTOServiceItem{}, err
	}

	return serviceItem, nil
}

func FetchRelatedDestinationSITFuelCharge(tx *pop.Connection, mtoServiceItemID uuid.UUID) (MTOServiceItem, error) {
	var serviceItem MTOServiceItem
	err := tx.RawQuery(
		`SELECT msi.id
            FROM mto_service_items msi
            INNER JOIN re_services res ON msi.re_service_id = res.id
            WHERE res.code IN (?) AND msi.mto_shipment_id IN (
                SELECT mto_shipment_id FROM mto_service_items WHERE id = ?)`, ReServiceCodeDDSFSC, mtoServiceItemID).First(&serviceItem)
	return serviceItem, err
}
