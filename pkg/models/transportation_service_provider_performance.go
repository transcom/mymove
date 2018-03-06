package models

import (
	"encoding/json"
	"fmt"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"
)

// awardCounts struct contains the number of shipments awarded to each tsp according to quality band
var awardCounts = map[int]int{
	1: 5,
	2: 3,
	3: 2,
	4: 1,
}

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

// NormalizedTransportationServiceProviderPerformance is a combination of all TSP
// performance metrics (BVS, Quality Band) for a performance period.
type NormalizedTransportationServiceProviderPerformance struct {
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
	NormalizedAwardCount            float64   `json:"normalized_award_count" db:"normalized_award_count"`
}

// String is not required by pop and may be deleted
func (t TransportationServiceProviderPerformance) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TransportationServiceProviderPerformances is a handy type for multiple TransportationServiceProviderPerformance structs
type TransportationServiceProviderPerformances []TransportationServiceProviderPerformance

// NormalizedTransportationServiceProviderPerformances is a handy type for multiple NormalizedTransportationServiceProviderPerformance structs
type NormalizedTransportationServiceProviderPerformances []NormalizedTransportationServiceProviderPerformance

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
	NormalizedTransportationServiceProviderPerformances, error) {

	sql := `SELECT *,
			CASE WHEN quality_band = 1 THEN award_count / 5
			WHEN quality_band = 2 THEN award_count / 3
			WHEN quality_band = 3 THEN award_count / 2
			WHEN quality_band = 4 THEN award_count / 1
            END AS normalized_award_count
		FROM
			transportation_service_provider_performances
		WHERE
			traffic_distribution_list_id = $1
			AND
			best_value_score > $2
		ORDER BY
			normalized_award_count ASC,
			best_value_score DESC
		`

	normalizedtsps := NormalizedTransportationServiceProviderPerformances{}
	err := tx.RawQuery(sql, tdlID, mps).All(&normalizedtsps)
	fmt.Println(normalizedtsps)
	return normalizedtsps, err
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

// type NormalizedTSPPerformances TransportationServiceProviderPerformance
func normalizeAwardCounts(tspPerfs TransportationServiceProviderPerformances) (tspNormPerfs NormalizedTransportationServiceProviderPerformances) {
	for _, tspPerf := range tspPerfs {
		tspPerf["NormalizedAwardCount"] = tspPerf.AwardCount / awardCounts[tspPerf.QualityBand]
	}
}

// ByNormalizedAwardCount implements sort.Interface for []tspPerformance based on
// the QualityBand field then the BestValueScore field then Award Count field.
// type ByNormalizedAwardCount TransportationServiceProviderPerformances

// func (q ByNormalizedAwardCount) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
// func (q ByNormalizedAwardCount) Len() int      { return len(q) }
// func (q ByNormalizedAwardCount) Less(i, j int) bool {
// 	if *q[i].NormalizedAwardCount != *q[j].NormalizedAwardCount {
// 		return *q[i].NormalizedAwardCount < *q[j].NormalizedAwardCount
// 	}
// 	if q[i].QualityBand != q[j].QualityBand {
// 		return q[i].QualityBand < q[j].QualityBand
// 	}
// 	return q[i].BestValueScore < q[j].BestValueScore
// }
