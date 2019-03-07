package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

// TransportationServiceProvider models moving companies used to move
// Shipments.
type TransportationServiceProvider struct {
	ID                       uuid.UUID `json:"id" db:"id"`
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
	StandardCarrierAlphaCode string    `json:"standard_carrier_alpha_code" db:"standard_carrier_alpha_code"`
	Enrolled                 bool      `json:"enrolled" db:"enrolled"`
	Name                     *string   `json:"name" db:"name"`
	SupplierID               *string   `json:"supplier_id" db:"supplier_id"`
	PocGeneralName           *string   `json:"poc_general_name" db:"poc_general_name"`
	PocGeneralEmail          *string   `json:"poc_general_email" db:"poc_general_email"`
	PocGeneralPhone          *string   `json:"poc_general_phone" db:"poc_general_phone"`
	PocClaimsName            *string   `json:"poc_claims_name" db:"poc_claims_name"`
	PocClaimsEmail           *string   `json:"poc_claims_email" db:"poc_claims_email"`
	PocClaimsPhone           *string   `json:"poc_claims_phone" db:"poc_claims_phone"`
}

// TSPWithBVSAndOfferCount represents a list of TSPs along with their BVS
// and offered shipment counts.
type TSPWithBVSAndOfferCount struct {
	ID                        uuid.UUID `json:"id" db:"id"`
	TrafficDistributionListID uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	BestValueScore            int       `json:"best_value_score" db:"best_value_score"`
	OfferCount                int       `json:"offer_count" db:"offer_count"`
}

// TSPWithBVSCount represents a list of TSPs along with their BVS counts.
type TSPWithBVSCount struct {
	ID                        uuid.UUID `json:"id" db:"id"`
	TrafficDistributionListID uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	BestValueScore            int       `json:"best_value_score" db:"best_value_score"`
}

// TransportationServiceProviders is not required by pop and may be deleted
type TransportationServiceProviders []TransportationServiceProvider

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (t *TransportationServiceProvider) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.StandardCarrierAlphaCode, Name: "StandardCarrierAlphaCode"},
	), nil
}

// MarshalLogObject is required to control the logging of this
func (t TransportationServiceProvider) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", t.ID.String())
	return nil
}

// FetchTransportationServiceProvider Fetches a TSP model
func FetchTransportationServiceProvider(db *pop.Connection, id uuid.UUID) (*TransportationServiceProvider, error) {
	var transportationServiceProvider TransportationServiceProvider
	err := db.Find(&transportationServiceProvider, id)

	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	return &transportationServiceProvider, nil
}
