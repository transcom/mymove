package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/gofrs/uuid"
)

type ServiceItemParamType string
type ServiceItemParamOrigin string

func (s ServiceItemParamType) String() string {
	return string(s)
}

func (s ServiceItemParamOrigin) String() string {
	return string(s)
}

const (
	ServiceItemParamTypeString   ServiceItemParamType   = "STRING"
	ServiceItemParamTypeDate     ServiceItemParamType   = "DATE"
	ServiceItemParamTypeInteger  ServiceItemParamType   = "INTEGER"
	ServiceItemParamTypeDecimal  ServiceItemParamType   = "DECIMAL"
	ServiceItemParamOriginPrime  ServiceItemParamOrigin = "PRIME"
	ServiceItemParamOriginSystem ServiceItemParamOrigin = "SYSTEM"
)

var validServiceItemParamType = []string{
	string(ServiceItemParamTypeString),
	string(ServiceItemParamTypeDate),
	string(ServiceItemParamTypeInteger),
	string(ServiceItemParamTypeDecimal),
}

var validServiceItemParamOrigin = []string{
	string(ServiceItemParamOriginPrime),
	string(ServiceItemParamOriginSystem),
}

type ServiceItemParamKey struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	Key         string                 `json:"key" db:"key"`
	Description string                 `json:"description" db:"description"`
	Type        ServiceItemParamType   `json:"type" db:"type"`
	Origin      ServiceItemParamOrigin `json:"origin" db:"origin"`
	CreatedAt   time.Time              `db:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
}

// ServiceItemParamKeys is not required by pop and may be deleted
type ServiceItemParamKeys []ServiceItemParamKey

func (s *ServiceItemParamKey) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(s.Type), Name: "Type"},
		&validators.StringIsPresent{Field: string(s.Origin), Name: "Origin"},
		&validators.StringInclusion{Field: s.Type.String(), Name: "Type", List: validServiceItemParamType},
		&validators.StringInclusion{Field: s.Origin.String(), Name: "Origin", List: validServiceItemParamOrigin},
	), nil
}
