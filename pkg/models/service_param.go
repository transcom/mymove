package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ServiceParam is a Service Parameter
type ServiceParam struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	ServiceID             uuid.UUID `json:"service_id" db:"service_id"`
	ServiceItemParamKeyID uuid.UUID `json:"service_item_param_key_id" db:"service_item_param_key_id"`
	IsOptional            bool      `json:"is_optional" db:"is_optional"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`

	// Associations
	Service             ReService           `belongs_to:"re_service" fk_id:"service_id"`
	ServiceItemParamKey ServiceItemParamKey `belongs_to:"service_item_param_key" fk_id:"service_item_param_key_id"`
}

// ServiceParams is not required by pop and may be deleted
type ServiceParams []ServiceParam

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *ServiceParam) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.ServiceID, Name: "ServiceID"},
		&validators.UUIDIsPresent{Field: p.ServiceItemParamKeyID, Name: "ServiceItemParamKeyID"},
	), nil
}
