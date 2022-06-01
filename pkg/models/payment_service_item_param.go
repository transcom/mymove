package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// PaymentServiceItemParam represents a parameter of the Payment Service Item
type PaymentServiceItemParam struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	PaymentServiceItemID  uuid.UUID `json:"payment_service_item_id" db:"payment_service_item_id"`
	ServiceItemParamKeyID uuid.UUID `json:"service_item_param_key_id" db:"service_item_param_key_id"`
	Value                 string    `json:"value" db:"value"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`

	// Associations
	PaymentServiceItem  PaymentServiceItem  `belongs_to:"payment_service_item" fk_id:"payment_service_item_id"`
	ServiceItemParamKey ServiceItemParamKey `belongs_to:"service_item_param_key" fk_id:"service_item_param_key_id"`

	// Used to lookup the appropriate ServiceItemParamKeyID when creating a PaymentServiceItemParam
	IncomingKey string `db:"-"`
}

// PaymentServiceItemParams is not required by pop and may be deleted
type PaymentServiceItemParams []PaymentServiceItemParam

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PaymentServiceItemParam) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.PaymentServiceItemID, Name: "PaymentServiceItemID"},
		&validators.UUIDIsPresent{Field: p.ServiceItemParamKeyID, Name: "ServiceItemParamKeyID"},
		&validators.StringIsPresent{Field: p.Value, Name: "Value"},
	), nil
}
