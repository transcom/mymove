package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type ServiceParam struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	ServiceID             uuid.UUID `json:"service_id" db:"service_id"`
	ServiceItemParamKeyID uuid.UUID `json:"service_item_param_key_id" db:"service_item_param_key_id"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`

	//Associations
	Service             ReServices          `belongs_to:"re_services"`
	ServiceItemParamKey ServiceItemParamKey `belongs_to:"service_item_param_keys"`
}

// ServiceParams is not required by pop and may be deleted
type ServiceParams []ServiceParams

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *ServiceParam) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.ServiceID, Name: "ServiceID"},
		&validators.UUIDIsPresent{Field: p.ServiceItemParamKeyID, Name: "ServiceItemParamKeyID"},
	), nil
}
