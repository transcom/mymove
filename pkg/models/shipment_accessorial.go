package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/unit"
)

// ShipmentAccessorialStatus represents the status of an accessorial's lifecycle
type ShipmentAccessorialStatus string

// ShipmentAccessorialLocation represents the location of the accessorial item
type ShipmentAccessorialLocation string

const (
	// ShipmentAccessorialStatusSUBMITTED captures enum value "SUBMITTED"
	ShipmentAccessorialStatusSUBMITTED ShipmentAccessorialStatus = "SUBMITTED"
	// ShipmentAccessorialStatusAPPROVED captures enum value "APPROVED"
	ShipmentAccessorialStatusAPPROVED ShipmentAccessorialStatus = "APPROVED"
	// ShipmentAccessorialStatusINVOICED captures enum value "INVOICED"
	ShipmentAccessorialStatusINVOICED ShipmentAccessorialStatus = "INVOICED"

	// ShipmentAccessorialLocationORIGIN captures enum value "ORIGIN"
	ShipmentAccessorialLocationORIGIN ShipmentAccessorialLocation = "ORIGIN"
	// ShipmentAccessorialLocationDESTINATION captures enum value "DESTINATION"
	ShipmentAccessorialLocationDESTINATION ShipmentAccessorialLocation = "DESTINATION"
	// ShipmentAccessorialLocationNEITHER captures enum value "NEITHER"
	ShipmentAccessorialLocationNEITHER ShipmentAccessorialLocation = "NEITHER"
)

// ShipmentAccessorial is an object representing an accessorial item in a pre-approval request
type ShipmentAccessorial struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ShipmentID uuid.UUID `json:"shipment_id" db:"shipment_id"`

	AccessorialID uuid.UUID                   `json:"accessorial_id" db:"accessorial_id"`
	Accessorial   Accessorial                 `belongs_to:"accessorials"`
	Location      ShipmentAccessorialLocation `json:"location" db:"location"`

	// Enter numbers only, no symbols or units. Examples:
	// Crating: enter "47.4" for crate size of 47.4 cu. ft.
	// 3rd-party service: enter "1299.99" for cost of $1,299.99.
	// Bulky item: enter "1" for a single item.
	Quantity1     unit.BaseQuantity         `json:"quantity_1" db:"quantity_1"`
	Quantity2     unit.BaseQuantity         `json:"quantity_2" db:"quantity_2"`
	Notes         string                    `json:"notes" db:"notes"`
	Status        ShipmentAccessorialStatus `json:"status" db:"status"`
	SubmittedDate time.Time                 `json:"submitted_date" db:"submitted_date"`
	ApprovedDate  time.Time                 `json:"approved_date" db:"approved_date"`
	CreatedAt     time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at" db:"updated_at"`
}

// FetchAccessorialsByShipmentID returns a list of accessorials by shipment_id
func FetchAccessorialsByShipmentID(dbConnection *pop.Connection, shipmentID *uuid.UUID) ([]ShipmentAccessorial, error) {
	var err error

	if shipmentID == nil {
		return nil, errors.Wrap(err, "Missing shipmentID")
	}

	accessorials := []ShipmentAccessorial{}

	query := dbConnection.Where("shipment_id = ?", *shipmentID)

	err = query.Eager().All(&accessorials)
	if err != nil {
		return accessorials, errors.Wrap(err, "Accessorials query failed")
	}

	return accessorials, err
}

// FetchShipmentAccessorialByID returns a shipment accessorial by id
func FetchShipmentAccessorialByID(dbConnection *pop.Connection, shipmentAccessorialID *uuid.UUID) (ShipmentAccessorial, error) {
	var err error

	if shipmentAccessorialID == nil {
		return ShipmentAccessorial{}, errors.Wrap(err, "Missing shipmentAccessorialID")
	}

	shipmentAccessorial := ShipmentAccessorial{}

	err = dbConnection.Eager().Find(&shipmentAccessorial, shipmentAccessorialID)
	if err != nil {
		return shipmentAccessorial, errors.Wrap(err, "Shipment accessorials query failed")
	}

	return shipmentAccessorial, err
}
