package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"
)

// BlackoutDate indicates the unavailable times for a TSP and includes
// its TDL as well.
type BlackoutDate struct {
	ID                              uuid.UUID `json:"id" db:"id"`
	CreatedAt                       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at" db:"updated_at"`
	TransportationServiceProviderId uuid.UUID `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	BlackoutDate                    time.Time `json:"blackout_date" db:"blackout_date"`
	TrafficDistributionListId       uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	CodeOfService                   string    `json:"code_of_service" db:"code_of_service"`
}

// String is not required by pop and may be deleted
func (b BlackoutDate) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// BlackoutDates is not required by pop and may be deleted
type BlackoutDates []BlackoutDate

// String is not required by pop and may be deleted
func (b BlackoutDates) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *BlackoutDate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: b.BlackoutDate, Name: "BlackoutDate"},
		&validators.StringIsPresent{Field: b.CodeOfService, Name: "CodeOfService"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *BlackoutDate) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *BlackoutDate) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
