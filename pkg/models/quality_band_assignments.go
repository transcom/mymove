package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"
)

// QualityBandAssignment connects a Transportation Service Provider to a Traffic
// Distribution List, assigns a quality band number and a performance period ID,
// and indicates how many shipments are made for each quality band.
type QualityBandAssignment struct {
	ID                              uuid.UUID `json:"id" db:"id"`
	CreatedAt                       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at" db:"updated_at"`
	TransportationServiceProviderID uuid.UUID `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	TrafficDistributionListID       uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	BandNumber                      int       `json:"band_number" db:"band_number"`
	PerformancePeriodID             uuid.UUID `json:"performance_period_id" db:"performance_period_id"`
	ShipmentsPerBand                int       `json:"shipments_per_band" db:"shipments_per_band"`
}

// String is not required by pop and may be deleted
func (q QualityBandAssignment) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// QualityBandAssignments is not required by pop and may be deleted
type QualityBandAssignments []QualityBandAssignment

// String is not required by pop and may be deleted
func (q QualityBandAssignments) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (q *QualityBandAssignment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: q.BandNumber, Name: "BandNumber"},
		&validators.IntIsPresent{Field: q.ShipmentsPerBand, Name: "ShipmentsPerBand"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (q *QualityBandAssignment) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (q *QualityBandAssignment) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
