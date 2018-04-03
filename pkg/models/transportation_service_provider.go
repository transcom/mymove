package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"go.uber.org/zap/zapcore"
)

// TransportationServiceProvider models moving companies used to move
// Shipments.
type TransportationServiceProvider struct {
	ID                       uuid.UUID `db:"id"`
	CreatedAt                time.Time `db:"created_at"`
	UpdatedAt                time.Time `db:"updated_at"`
	StandardCarrierAlphaCode string    `db:"standard_carrier_alpha_code"`
	Name                     string    `db:"name"`
}

// TSPWithBVSAndOfferCount represents a list of TSPs along with their BVS
// and offered shipment counts.
type TSPWithBVSAndOfferCount struct {
	ID                        uuid.UUID `db:"id"`
	Name                      string    `db:"name"`
	TrafficDistributionListID uuid.UUID `db:"traffic_distribution_list_id"`
	BestValueScore            int       `db:"best_value_score"`
	OfferCount                int       `db:"offer_count"`
}

// TSPWithBVSCount represents a list of TSPs along with their BVS counts.
type TSPWithBVSCount struct {
	ID                        uuid.UUID `db:"id"`
	Name                      string    `db:"name"`
	TrafficDistributionListID uuid.UUID `db:"traffic_distribution_list_id"`
	BestValueScore            int       `db:"best_value_score"`
}

// TransportationServiceProviders is not required by pop and may be deleted
type TransportationServiceProviders []TransportationServiceProvider

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (t *TransportationServiceProvider) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.StandardCarrierAlphaCode, Name: "StandardCarrierAlphaCode"},
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
	), nil
}

// MarshalLogObject is required to control the logging of this
func (t TransportationServiceProvider) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", t.ID.String())
	encoder.AddString("name", t.Name)
	return nil
}
