package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"
)

// TransportationServiceProviderPerformance is a combination of all TSP
// performance metrics (BVS, Quality Band) for a performance period.
type TransportationServiceProviderPerformance struct {
	ID                              uuid.UUID `json:"id" db:"id"`
	CreatedAt                       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at" db:"updated_at"`
	PerformancePeriodStart          time.Time `json:"performance_period_start" db:"performance_period_start"`
	PerformancePeriodEnd            time.Time `json:"performance_period_end" db:"performance_period_end"`
	TrafficDistributionListID       uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	TransportationServiceProviderID uuid.UUID `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	QualityBand                     *int      `json:"quality_band" db:"quality_band"`
	BestValueScore                  int       `json:"best_value_score" db:"best_value_score"`
	AwardCount                      int       `json:"award_count" db:"award_count"`
}

// String is not required by pop and may be deleted
func (t TransportationServiceProviderPerformance) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TransportationServiceProviderPerformances is a handy type for multiple TransportationServiceProviderPerformance structs
type TransportationServiceProviderPerformances []TransportationServiceProviderPerformance

// String is not required by pop and may be deleted
func (t TransportationServiceProviderPerformances) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (t *TransportationServiceProviderPerformance) Validate(tx *pop.Connection) (*validate.Errors, error) {
	// Pop can't validate pointers to ints, so turn the pointer into an integer.
	// Our valid values are [nil, 1, 2, 3, 4]
	qualityBand := 1
	if t.QualityBand != nil {
		qualityBand = *t.QualityBand
	}

	return validate.Validate(
		// Quality Bands can have a range from 1 - 4 as defined in DTR 402. See page 67 of
		// https://www.ustranscom.mil/dtr/part-iv/dtr-part-4-402.pdf
		&validators.IntIsGreaterThan{Field: qualityBand, Name: "QualityBand", Compared: 0},
		&validators.IntIsLessThan{Field: qualityBand, Name: "QualityBand", Compared: 5},

		// Best Value Scores can range from 0 - 100, as defined in DTR403. See page 7
		// of https://www.ustranscom.mil/dtr/part-iv/dtr-part-4-403.pdf
		&validators.IntIsGreaterThan{Field: t.BestValueScore, Name: "BestValueScore", Compared: -1},
		&validators.IntIsLessThan{Field: t.BestValueScore, Name: "BestValueScore", Compared: 101},
	), nil
}

// FetchTSPPerformanceForAwardQueue returns TSP performance records in a given TDL
// in the order that they should be awarded new shipments.
func FetchTSPPerformanceForAwardQueue(tx *pop.Connection, tdlID uuid.UUID, mps int) (
	TransportationServiceProviderPerformances, error) {

	sql := `SELECT
			*
		FROM
			transportation_service_provider_performances
		WHERE
			traffic_distribution_list_id = $1
			AND
			best_value_score > $2
		ORDER BY
			award_count ASC,
			quality_band DESC,
			best_value_score DESC
		`

	tsps := TransportationServiceProviderPerformances{}
	err := tx.RawQuery(sql, tdlID, mps).All(&tsps)

	return tsps, err
}

// FetchTSPPerformanceForQualityBandAssignment returns TSPs in a given TDL in the
// order that they should be assigned quality bands.
func FetchTSPPerformanceForQualityBandAssignment(tx *pop.Connection, tdlID uuid.UUID, mps int) (TransportationServiceProviderPerformances, error) {

	sql := `SELECT
			*
		FROM
			transportation_service_provider_performances
		WHERE
			traffic_distribution_list_id = $1
			AND
			best_value_score > $2
		ORDER BY
			best_value_score DESC
		`

	tsps := TransportationServiceProviderPerformances{}
	err := tx.RawQuery(sql, tdlID, mps).All(&tsps)

	return tsps, err
}
