package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/gofrs/uuid"
)

// ServiceItemParamType is a type of service item parameter
type ServiceItemParamType string

// String is a string representation of a ServiceItemParamType
func (s ServiceItemParamType) String() string {
	return string(s)
}

const (
	// ServiceItemParamTypeString is a string
	ServiceItemParamTypeString ServiceItemParamType = "STRING"
	// ServiceItemParamTypeDate is a date
	ServiceItemParamTypeDate ServiceItemParamType = "DATE"
	// ServiceItemParamTypeInteger is an integer
	ServiceItemParamTypeInteger ServiceItemParamType = "INTEGER"
	// ServiceItemParamTypeDecimal is a decimal
	ServiceItemParamTypeDecimal ServiceItemParamType = "DECIMAL"
)

// ServiceItemParamOrigin is a type of service item parameter origin
type ServiceItemParamOrigin string

// String is a string representation of a ServiceItemParamOrigin
func (s ServiceItemParamOrigin) String() string {
	return string(s)
}

const (
	// ServiceItemParamOriginPrime is the Prime origin
	ServiceItemParamOriginPrime ServiceItemParamOrigin = "PRIME"
	// ServiceItemParamOriginSystem is the System origin
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

// ServiceItemParamKey is a key for a Service Item Param
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

// Validate validates a ServiceItemParamKey
func (s *ServiceItemParamKey) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Key, Name: "Key"},
		&validators.StringIsPresent{Field: s.Description, Name: "Description"},
		&validators.StringIsPresent{Field: string(s.Type), Name: "Type"},
		&validators.StringIsPresent{Field: string(s.Origin), Name: "Origin"},
		&validators.StringInclusion{Field: s.Type.String(), Name: "Type", List: validServiceItemParamType},
		&validators.StringInclusion{Field: s.Origin.String(), Name: "Origin", List: validServiceItemParamOrigin},
	), nil
}
