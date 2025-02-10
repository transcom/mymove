package models

import (
	"database/sql/driver"
	"fmt"
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
	ExternalCrate                     *bool                          `db:"external_crate"`
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
	if db != nil {
		var serviceItem MTOServiceItem
		err := db.Eager("SITDestinationOriginalAddress",
			"SITDestinationFinalAddress",
			"ReService",
			"CustomerContacts",
			"MTOShipment.PickupAddress",
			"MTOShipment.DestinationAddress",
			"Dimensions").Where("id = ?", serviceItemID).First(&serviceItem)

		if err != nil {
			if errors.Cause(err).Error() == RecordNotFoundErrorString {
				return MTOServiceItem{}, ErrFetchNotFound
			}
			return MTOServiceItem{}, err
		}
		return serviceItem, nil
	} else {
		return MTOServiceItem{}, errors.New("db connection is nil; unable to fetch service item")
	}
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

type MTOServiceItemType struct {
	ID                                *uuid.UUID           `json:"id"`
	MoveID                            *uuid.UUID           `json:"move_Id"`
	ReServiceID                       *uuid.UUID           `json:"re_service_id"`
	CreatedAt                         *time.Time           `json:"created_at"`
	UpdatedAt                         *time.Time           `json:"updated_at"`
	Reason                            *string              `json:"reason"`
	PickupPostalCode                  *string              `json:"pickup_postal_code"`
	Description                       *string              `json:"description"`
	Status                            *string              `json:"status"`
	RejectionReason                   *string              `json:"rejected_reason"`
	ApprovedAt                        *time.Time           `json:"approved_at"`
	SITPostalCode                     *string              `json:"sit_postal_code"`
	SITEntryDate                      *time.Time           `json:"sit_entry_date"`
	SITDepartureDate                  *time.Time           `json:"sit_departure_date"`
	SITDestinationFinalAddressID      *uuid.UUID           `json:"sit_destination_final_address_id"`
	SITOriginHHGOriginalAddressID     *uuid.UUID           `json:"sit_origin_hhg_original_address_id"`
	SITOriginHHGActualAddressID       *uuid.UUID           `json:"sit_origin_hhg_actual_address_id"`
	EstimatedWeight                   *unit.Pound          `json:"estimated_weight"`
	ActualWeight                      *unit.Pound          `json:"actual_weight"`
	SITDestinationOriginalAddressID   *uuid.UUID           `json:"sit_destination_original_address_id"`
	SITCustomerContacted              *time.Time           `json:"sit_customer_contacted"`
	SITRequestedDelivery              *time.Time           `json:"sit_requested_delivery"`
	RequestedApprovalsRequestedStatus *bool                `json:"requested_approvals_requested_status"`
	CustomerExpense                   *bool                `json:"customer_expense"`
	CustomerExpenseReason             *string              `json:"customer_expense_reason"`
	SITDeliveryMiles                  *int                 `json:"sit_delivery_miles"`
	PricingEstimate                   *unit.Cents          `json:"pricing_estimate"`
	StandaloneCrate                   *bool                `json:"standalone_crate"`
	LockedPriceCents                  *unit.Cents          `json:"locked_price_cents"`
	ServiceLocation                   *ServiceLocationType `json:"service_location"`
	POELocationID                     *uuid.UUID           `json:"poe_location_id"`
	PODLocationID                     *uuid.UUID           `json:"pod_location_id"`
	ExternalCrate                     *bool                `json:"external_crate"`
}

func (m MTOServiceItem) GetMTOServiceItemTypeFromServiceItem() MTOServiceItemType {
	return MTOServiceItemType{
		ID:                                &m.ID,
		MoveID:                            &m.MoveTaskOrderID,
		ReServiceID:                       &m.ReServiceID,
		CreatedAt:                         &m.CreatedAt,
		UpdatedAt:                         &m.UpdatedAt,
		Reason:                            m.Reason,
		PickupPostalCode:                  m.PickupPostalCode,
		Description:                       m.Description,
		Status:                            (*string)(&m.Status),
		RejectionReason:                   m.RejectionReason,
		ApprovedAt:                        m.ApprovedAt,
		SITPostalCode:                     m.SITPostalCode,
		SITEntryDate:                      m.SITEntryDate,
		SITDepartureDate:                  m.SITDepartureDate,
		SITDestinationFinalAddressID:      m.SITDestinationFinalAddressID,
		SITOriginHHGOriginalAddressID:     m.SITOriginHHGOriginalAddressID,
		SITOriginHHGActualAddressID:       m.SITOriginHHGActualAddressID,
		EstimatedWeight:                   m.EstimatedWeight,
		ActualWeight:                      m.ActualWeight,
		SITDestinationOriginalAddressID:   m.SITDestinationOriginalAddressID,
		SITCustomerContacted:              m.SITCustomerContacted,
		SITRequestedDelivery:              m.SITRequestedDelivery,
		RequestedApprovalsRequestedStatus: m.RequestedApprovalsRequestedStatus,
		CustomerExpense:                   &m.CustomerExpense,
		CustomerExpenseReason:             m.CustomerExpenseReason,
		SITDeliveryMiles:                  m.SITDeliveryMiles,
		PricingEstimate:                   m.PricingEstimate,
		StandaloneCrate:                   m.StandaloneCrate,
		LockedPriceCents:                  m.LockedPriceCents,
		ServiceLocation:                   m.ServiceLocation,
		POELocationID:                     m.POELocationID,
		PODLocationID:                     m.PODLocationID,
		ExternalCrate:                     m.ExternalCrate,
	}
}

func (m MTOServiceItem) Value() (driver.Value, error) {
	var id string
	var moveTaskOrderID string
	var mtoShipmentID string
	var reason string
	var pickupPostalCode string
	var description string
	var rejectionReason string
	var approvedAt string
	var rejectedAt string
	var sitPostalCode string
	var sitRequestedDelivery string
	var requestedApprovalsRequestedStatus bool
	var serviceLocation string
	var sitEntryDate string
	var sitDepartureDate string
	var sitCustomerContacted string
	var poeLocationID string
	var podLocationID string
	var sitDestinationFinalAddressID string
	var sitOriginHHGOriginalAddressID string
	var sitOriginHHGActualAddressID string
	var sitDestinationOriginalAddressID string
	var standaloneCrate bool
	var lockedPriceCents int64
	var sitDeliveryMiles int
	var customerExpenseReason string
	var estimatedWeight int64
	var actualWeight int64
	var pricingEstimate int64
	var externalCrate bool

	if m.ID != uuid.Nil {
		id = m.ID.String()
	}

	if m.MoveTaskOrderID != uuid.Nil {
		moveTaskOrderID = m.MoveTaskOrderID.String()
	}

	if *m.MTOShipmentID != uuid.Nil {
		mtoShipmentID = m.MTOShipmentID.String()
	}

	if m.Reason != nil {
		reason = *m.Reason
	}

	if m.PickupPostalCode != nil {
		pickupPostalCode = *m.PickupPostalCode
	}

	if m.Description != nil {
		description = *m.Description
	}

	if m.RejectionReason != nil {
		rejectionReason = *m.RejectionReason
	}

	if m.ApprovedAt != nil {
		approvedAt = m.ApprovedAt.Format("2006-01-02 15:04:05")
	}

	if m.RejectedAt != nil {
		rejectedAt = m.RejectedAt.Format("2006-01-02 15:04:05")
	}

	if m.SITPostalCode != nil {
		sitPostalCode = *m.SITPostalCode
	}

	if m.SITRequestedDelivery != nil {
		sitRequestedDelivery = m.SITRequestedDelivery.Format("2006-01-02 15:04:05")
	}

	if m.RequestedApprovalsRequestedStatus != nil {
		requestedApprovalsRequestedStatus = *m.RequestedApprovalsRequestedStatus
	}

	if m.ServiceLocation != nil {
		serviceLocation = string(*m.ServiceLocation)
	}

	if m.SITEntryDate != nil {
		sitEntryDate = m.SITEntryDate.Format("2006-01-02 15:04:05")
	}

	if m.SITDepartureDate != nil {
		sitDepartureDate = m.SITDepartureDate.Format("2006-01-02 15:04:05")
	}

	if m.SITCustomerContacted != nil {
		sitCustomerContacted = m.SITCustomerContacted.Format("2006-01-02 15:04:05")
	}

	if m.POELocationID != nil {
		poeLocationID = m.POELocationID.String()
	}

	if m.PODLocationID != nil {
		podLocationID = m.PODLocationID.String()
	}

	if m.SITDestinationFinalAddressID != nil {
		sitDestinationFinalAddressID = m.SITDestinationFinalAddressID.String()
	}

	if m.SITOriginHHGActualAddressID != nil {
		sitOriginHHGActualAddressID = m.SITOriginHHGActualAddressID.String()
	}

	if m.SITDestinationOriginalAddressID != nil {
		sitDestinationOriginalAddressID = m.SITDestinationOriginalAddressID.String()
	}

	if m.SITOriginHHGOriginalAddressID != nil {
		sitOriginHHGOriginalAddressID = m.SITOriginHHGOriginalAddressID.String()
	}

	if m.StandaloneCrate != nil {
		standaloneCrate = *m.StandaloneCrate
	}

	if m.ExternalCrate != nil {
		externalCrate = *m.ExternalCrate
	}

	if m.LockedPriceCents != nil {
		lockedPriceCents = m.LockedPriceCents.Int64()
	}

	if m.SITDeliveryMiles != nil {
		sitDeliveryMiles = *m.SITDeliveryMiles
	}

	if m.CustomerExpenseReason != nil {
		customerExpenseReason = *m.CustomerExpenseReason
	}

	if m.EstimatedWeight != nil {
		estimatedWeight = m.EstimatedWeight.Int64()
	}

	if m.ActualWeight != nil {
		actualWeight = m.ActualWeight.Int64()
	}

	if m.PricingEstimate != nil {
		pricingEstimate = m.PricingEstimate.Int64()
	}

	s := fmt.Sprintf("(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%d,%d,%s,%s,%s,%t,%t,%s,%d,%d,%t,%d,%s,%s,%s,%s,%t)",
		id,
		moveTaskOrderID,
		mtoShipmentID,
		m.CreatedAt.Format("2006-01-02 15:04:05"),
		m.UpdatedAt.Format("2006-01-02 15:04:05"),
		reason,
		pickupPostalCode,
		description,
		m.Status,
		rejectionReason,
		approvedAt,
		rejectedAt,
		sitPostalCode,
		sitEntryDate,
		sitDepartureDate,
		sitDestinationFinalAddressID,
		sitOriginHHGOriginalAddressID,
		sitOriginHHGActualAddressID,
		estimatedWeight,
		actualWeight,
		sitDestinationOriginalAddressID,
		sitCustomerContacted,
		sitRequestedDelivery,
		requestedApprovalsRequestedStatus,
		m.CustomerExpense,
		customerExpenseReason,
		sitDeliveryMiles,
		pricingEstimate,
		standaloneCrate,
		lockedPriceCents,
		serviceLocation,
		poeLocationID,
		podLocationID,
		m.ReService.Code.String(),
		externalCrate,
	)
	return []byte(s), nil
}
