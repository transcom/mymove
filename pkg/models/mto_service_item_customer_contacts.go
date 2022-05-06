package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// CustomerContactType determines what type of customer contact for a service item.
// For domestic destination 1st day SIT.
type CustomerContactType string

const (
	// CustomerContactTypeFirst describes customer contacts for a FIRST type.
	CustomerContactTypeFirst CustomerContactType = "FIRST"
	// CustomerContactTypeSecond  describes customer contacts for a SECOND type.
	CustomerContactTypeSecond CustomerContactType = "SECOND"
)

// MTOServiceItemCustomerContact is an object representing customer contact for a service item.
type MTOServiceItemCustomerContact struct {
	ID                         uuid.UUID           `db:"id"`
	MTOServiceItem             MTOServiceItem      `belongs_to:"mto_service_items" fk_id:"mto_service_item_id"`
	MTOServiceItemID           uuid.UUID           `db:"mto_service_item_id"`
	Type                       CustomerContactType `db:"type"`
	TimeMilitary               string              `db:"time_military"`
	FirstAvailableDeliveryDate time.Time           `db:"first_available_delivery_date"`
	CreatedAt                  time.Time           `db:"created_at"`
	UpdatedAt                  time.Time           `db:"updated_at"`
}

// MTOServiceItemCustomerContacts is a slice containing MTOServiceItemCustomerContact.
type MTOServiceItemCustomerContacts []MTOServiceItemCustomerContact

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MTOServiceItemCustomerContact) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MTOServiceItemID, Name: "MTOServiceItemID"})
	vs = append(vs, &validators.StringInclusion{Field: string(m.Type), Name: "Type", List: []string{
		string(CustomerContactTypeFirst),
		string(CustomerContactTypeSecond),
	}})
	vs = append(vs, &validators.StringIsPresent{Field: m.TimeMilitary, Name: "TimeMilitary"})
	vs = append(vs, &validators.TimeIsPresent{Field: m.FirstAvailableDeliveryDate, Name: "FirstAvailableDeliveryDate"})

	return validate.Validate(vs...), nil
}

// TableName overrides the table name used by Pop.
func (m MTOServiceItemCustomerContact) TableName() string {
	return "mto_service_item_customer_contacts"
}
