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

// FetchTransportationServiceProvidersInTDL returns TSPs in a given TDL in the
// order that they should be awarded new shipments.
func FetchTransportationServiceProvidersInTDL(tx *pop.Connection, tdlID uuid.UUID) ([]TSPWithBVSAndAwardCount, error) {
	// We need to get TSPs, along with their Best Value Scores and total
	// awarded shipments, hence the two joins. Some notes on the query:
	// - We min() the id and scores, because we need an aggregate function given
	//   that it's a GROUP BY
	// - the UUID is CAST() to text to work inside the MIN(), it doesn't accept UUIDs
	// - We might be able to replace this with Pop's join syntax for easier reading:
	//   https://github.com/markbates/pop#join-query
	sql := `SELECT
			MIN(CAST(transportation_service_providers.id AS text)) as id,
			MIN(transportation_service_providers.name) as name,
			MIN(CAST(best_value_scores.traffic_distribution_list_id AS text)) as traffic_distribution_list_id,
			MIN(best_value_scores.score) as best_value_score,
			COUNT(shipment_awards.id) as award_count
		FROM
			transportation_service_providers
		JOIN best_value_scores ON
			transportation_service_providers.id = best_value_scores.transportation_service_provider_id
		LEFT JOIN shipment_awards ON
			transportation_service_providers.id = shipment_awards.transportation_service_provider_id
		WHERE
			best_value_scores.traffic_distribution_list_id = ?
		GROUP BY best_value_scores.id
		ORDER BY award_count ASC, best_value_score DESC
		`

	tsps := []TSPWithBVSAndAwardCount{}
	err := tx.RawQuery(sql, tdlID).All(&tsps)

	return tsps, err
}
