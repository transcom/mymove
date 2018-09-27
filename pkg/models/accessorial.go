package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/unit"
)

// AccessorialStatus represents the status of an accessorial's lifecycle
type AccessorialStatus string

// AccessorialLocation represents the location of the accessorial item
type AccessorialLocation string

const (
	// AccessorialStatusSUBMITTED captures enum value "SUBMITTED"
	AccessorialStatusSUBMITTED AccessorialStatus = "SUBMITTED"
	// AccessorialStatusAPPROVED captures enum value "APPROVED"
	AccessorialStatusAPPROVED AccessorialStatus = "APPROVED"
	// AccessorialStatusINVOICED captures enum value "INVOICED"
	AccessorialStatusINVOICED AccessorialStatus = "INVOICED"

	// AccessorialLocationORIGIN captures enum value "ORIGIN"
	AccessorialLocationORIGIN AccessorialLocation = "ORIGIN"
	// AccessorialLocationDESTINATION captures enum value "DESTINATION"
	AccessorialLocationDESTINATION AccessorialLocation = "DESTINATION"
	// AccessorialLocationNEITHER captures enum value "NEITHER"
	AccessorialLocationNEITHER AccessorialLocation = "NEITHER"
)

// Accessorial is an object representing an accessorial item in a pre-approval request
type Accessorial struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ShipmentID uuid.UUID `json:"shipment_id" db:"shipment_id"`

	// Code and Item description are linked, how should we store that information? What gets stored in the database?
	// for now, code is a string. It will become a reference to code table once that is designed.
	Code string `json:"code" db:"code"`
	// This is the item description for the code
	Item     string              `json:"item" db:"item"`
	Location AccessorialLocation `json:"location" db:"location"`

	// Enter numbers only, no symbols or units. Examples:
	// Crating: enter "47.4" for crate size of 47.4 cu. ft.
	// 3rd-party service: enter "1299.99" for cost of $1,299.99.
	// Bulky item: enter "1" for a single item.

	Quantity      unit.BaseQuantity `json:"quantity" db:"quantity"`
	Notes         string            `json:"notes" db:"notes"`
	Status        AccessorialStatus `json:"status" db:"status"`
	SubmittedDate time.Time         `json:"submitted_date" db:"submitted_date"`
	ApprovedDate  time.Time         `json:"approved_date" db:"approved_date"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" db:"updated_at"`
}

// FetchAccessorialsByShipmentID returns a list of accessorials by shipment_id
func FetchAccessorialsByShipmentID(dbConnection *pop.Connection, shipmentID *uuid.UUID) ([]Accessorial, error) {
	var err error

	if shipmentID == nil {
		return nil, errors.Wrap(err, "Missing shipmentID")
	}

	accessorials := []Accessorial{}

	query := dbConnection.Where("shipment_id = ?", *shipmentID)

	err = query.All(&accessorials)
	if err != nil {
		return accessorials, errors.Wrap(err, "Accessorials query failed")
	}

	return accessorials, err
}
