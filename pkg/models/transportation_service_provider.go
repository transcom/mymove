package models

import (
	"encoding/json"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"
)

// TransportationServiceProvider models moving companies used to move
// Shipments.
type TransportationServiceProvider struct {
	ID                       uuid.UUID `json:"id" db:"id"`
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
	StandardCarrierAlphaCode string    `json:"standard_carrier_alpha_code" db:"standard_carrier_alpha_code"`
	Name                     string    `json:"name" db:"name"`
}

// TSPWithBVSAndAwardCount represents a list of TSPs along with their BVS
// and awarded shipment counts.
type TSPWithBVSAndAwardCount struct {
	ID                        uuid.UUID `json:"id" db:"id"`
	Name                      string    `json:"name" db:"name"`
	TrafficDistributionListID uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	BestValueScore            int       `json:"best_value_score" db:"best_value_score"`
	AwardCount                int       `json:"award_count" db:"award_count"`
}

// TSPWithBVSCount represents a list of TSPs along with their BVS counts.
type TSPWithBVSCount struct {
	ID                        uuid.UUID `json:"id" db:"id"`
	Name                      string    `json:"name" db:"name"`
	TrafficDistributionListID uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	BestValueScore            int       `json:"best_value_score" db:"best_value_score"`
}

// String is not required by pop and may be deleted
func (t TransportationServiceProvider) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TransportationServiceProviders is not required by pop and may be deleted
type TransportationServiceProviders []TransportationServiceProvider

// String is not required by pop and may be deleted
func (t TransportationServiceProviders) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (t *TransportationServiceProvider) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.StandardCarrierAlphaCode, Name: "StandardCarrierAlphaCode"},
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
	), nil
}
