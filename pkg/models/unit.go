package models

import "github.com/gofrs/uuid"

// Unit represents a Department of Defense Entity, uniquely identified by a Unit Identification Code (UIC)
type Unit struct {
	ID                     uuid.UUID            `json:"id" db:"id"`
	UIC                    *string              `json:"uic" db:"uic"`
	Name                   *string              `json:"name" db:"name"`
	City                   *string              `json:"city" db:"city"`
	Locality               *string              `json:"locality" db:"locality"`
	Country                *string              `json:"country" db:"country"`
	PostalCode             *string              `json:"postal_code" db:"postal_code"`
	TransportationOfficeID uuid.UUID            `json:"transportation_office_id" db:"transportation_office_id"`
	TransportationOffice   TransportationOffice `belongs_to:"transportation_offices"`
}
