package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/unit"
)

// ShipmentLineItemStatus represents the status of a line item's lifecycle
type ShipmentLineItemStatus string

// ShipmentLineItemLocation represents the location of the line item
type ShipmentLineItemLocation string

const (
	// ShipmentLineItemStatusSUBMITTED captures enum value "SUBMITTED"
	ShipmentLineItemStatusSUBMITTED ShipmentLineItemStatus = "SUBMITTED"
	// ShipmentLineItemStatusAPPROVED captures enum value "APPROVED"
	ShipmentLineItemStatusAPPROVED ShipmentLineItemStatus = "APPROVED"
	// ShipmentLineItemStatusINVOICED captures enum value "INVOICED"
	ShipmentLineItemStatusINVOICED ShipmentLineItemStatus = "INVOICED"

	// ShipmentLineItemLocationORIGIN captures enum value "ORIGIN"
	ShipmentLineItemLocationORIGIN ShipmentLineItemLocation = "ORIGIN"
	// ShipmentLineItemLocationDESTINATION captures enum value "DESTINATION"
	ShipmentLineItemLocationDESTINATION ShipmentLineItemLocation = "DESTINATION"
	// ShipmentLineItemLocationNEITHER captures enum value "NEITHER"
	ShipmentLineItemLocationNEITHER ShipmentLineItemLocation = "NEITHER"
)

// ShipmentLineItem is an object representing a line item in a pre-approval request
type ShipmentLineItem struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ShipmentID uuid.UUID `json:"shipment_id" db:"shipment_id"`

	Tariff400ngItemID uuid.UUID                `json:"tariff400ng_item_id" db:"tariff400ng_item_id"`
	Tariff400ngItem   Tariff400ngItem          `belongs_to:"tariff400ng_items"`
	Location          ShipmentLineItemLocation `json:"location" db:"location"`

	// Enter numbers only, no symbols or units. Examples:
	// Crating: enter "47.4" for crate size of 47.4 cu. ft.
	// 3rd-party service: enter "1299.99" for cost of $1,299.99.
	// Bulky item: enter "1" for a single item.
	Quantity1     unit.BaseQuantity      `json:"quantity_1" db:"quantity_1"`
	Quantity2     unit.BaseQuantity      `json:"quantity_2" db:"quantity_2"`
	Notes         string                 `json:"notes" db:"notes"`
	Status        ShipmentLineItemStatus `json:"status" db:"status"`
	SubmittedDate time.Time              `json:"submitted_date" db:"submitted_date"`
	ApprovedDate  time.Time              `json:"approved_date" db:"approved_date"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
}

// FetchLineItemsByShipmentID returns a list of line items by shipment_id
func FetchLineItemsByShipmentID(dbConnection *pop.Connection, shipmentID *uuid.UUID) ([]ShipmentLineItem, error) {
	var err error

	if shipmentID == nil {
		return nil, errors.Wrap(err, "Missing shipmentID")
	}

	shipmentLineItems := []ShipmentLineItem{}

	query := dbConnection.Where("shipment_id = ?", *shipmentID)

	err = query.Eager().All(&shipmentLineItems)
	if err != nil {
		return shipmentLineItems, errors.Wrap(err, "Fetch line items query failed")
	}

	return shipmentLineItems, err
}

// FetchShipmentLineItemByID returns a shipment line item by id
func FetchShipmentLineItemByID(dbConnection *pop.Connection, shipmentLineItemID *uuid.UUID) (ShipmentLineItem, error) {
	var err error

	if shipmentLineItemID == nil {
		return ShipmentLineItem{}, errors.Wrap(err, "Missing shipmentLineItemID")
	}

	shipmentLineItem := ShipmentLineItem{}

	err = dbConnection.Eager().Find(&shipmentLineItem, shipmentLineItemID)
	if err != nil {
		return shipmentLineItem, errors.Wrap(err, "Shipment line items query failed")
	}

	return shipmentLineItem, err
}

// Approve marks the ShipmentLineItem request as Approved. Must be in a submitted state.
func (s *ShipmentLineItem) Approve() error {
	if s.Status != ShipmentLineItemStatusSUBMITTED {
		return errors.Wrap(ErrInvalidTransition, "Approve")
	}
	s.Status = ShipmentLineItemStatusAPPROVED
	return nil
}
