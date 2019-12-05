package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type PaymentServiceItemParam struct {
	ID               uuid.UUID                `json:"id" db:"id"`
	PaymentServiceItemID uuid.UUID                `json:"payment_service_item_id" db:"payment_service_item_id"`
	ServiceItemParamKeyID    uuid.UUID                `json:"service_item_param_key_id" db:"service_item_param_key_id"`
	Value string `json:"value" db:"value"`
	CreatedAt        time.Time                `db:"created_at"`
	UpdatedAt        time.Time                `db:"updated_at"`

	//Associations
	PaymentServiceItem PaymentServiceItem `belongs_to:"payment_service_items"`
	ServiceItemParamKey ServiceItemParamKey `belongs_to:"service_item_param_keys"`
}

// PaymentServiceItemParams is not required by pop and may be deleted
type PaymentServiceItemParams []PaymentServiceItemParam

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PaymentServiceItemParam) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.PaymentServiceItemID, Name: "PaymentServiceItemID"},
		&validators.UUIDIsPresent{Field: p.ServiceItemParamKeyID, Name: "ServiceItemParamKeyID"},
	), nil
}
