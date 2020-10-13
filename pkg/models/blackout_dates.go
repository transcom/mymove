package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// BlackoutDate indicates the range of unavailable times for a TSP and includes its TDL as well.
type BlackoutDate struct {
	ID                              uuid.UUID  `json:"id" db:"id"`
	CreatedAt                       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at" db:"updated_at"`
	TransportationServiceProviderID uuid.UUID  `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	StartBlackoutDate               time.Time  `json:"start_blackout_date" db:"start_blackout_date"`
	EndBlackoutDate                 time.Time  `json:"end_blackout_date" db:"end_blackout_date"`
	TrafficDistributionListID       *uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	Market                          *string    `json:"market" db:"market"`
	SourceGBLOC                     *string    `json:"source_gbloc" db:"source_gbloc"`
	Zip3                            *int       `json:"zip3" db:"zip3"`
	VolumeMove                      *bool      `json:"volume_move" db:"volume_move"`
}

// BlackoutDates is not required by pop and may be deleted
type BlackoutDates []BlackoutDate

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *BlackoutDate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: b.StartBlackoutDate, Name: "StartBlackoutDate"},
		&validators.TimeIsPresent{Field: b.EndBlackoutDate, Name: "EndBlackoutDate"},
		// &validators.StringIsPresent{Field: b.CodeOfService, Name: "CodeOfService"},
		// TODO: write our own validator that can validate pointers; Pop lacks that
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
