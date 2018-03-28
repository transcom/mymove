package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	//"go.uber.org/zap"
)

var qualityBands = []int{1, 2, 3, 4}

// OffersPerQualityBand is a map of the number of shipments to be offered per round to each quality band
var OffersPerQualityBand = map[int]int{
	1: 5,
	2: 3,
	3: 2,
	4: 1,
}

// TransportationServiceProviderPerformance is a combination of all TSP
// performance metrics (BVS, Quality Band) for a performance period.
type TransportationServiceProviderPerformance struct {
	ID                              uuid.UUID `db:"id"`
	CreatedAt                       time.Time `db:"created_at"`
	UpdatedAt                       time.Time `db:"updated_at"`
	PerformancePeriodStart          time.Time `db:"performance_period_start"`
	PerformancePeriodEnd            time.Time `db:"performance_period_end"`
	RateCycleStart                  time.Time `db:"rate_cycle_start"`
	RateCycleEnd                    time.Time `db:"rate_cycle_end"`
	TrafficDistributionListID       uuid.UUID `db:"traffic_distribution_list_id"`
	TransportationServiceProviderID uuid.UUID `db:"transportation_service_provider_id"`
	QualityBand                     *int      `db:"quality_band"`
	BestValueScore                  int       `db:"best_value_score"`
	OfferCount                      int       `db:"offer_count"`
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
		// Start times should be before End times
		&validators.TimeIsBeforeTime{FirstTime: t.PerformancePeriodStart, FirstName: "PerformancePeriodStart",
			SecondTime: t.PerformancePeriodEnd, SecondName: "PerformancePeriodEnd"},
		&validators.TimeIsBeforeTime{FirstTime: t.RateCycleStart, FirstName: "RateCycleStart",
			SecondTime: t.RateCycleEnd, SecondName: "RateCycleEnd"},

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

// NextTSPPerformanceInQualityBand returns the TSP performance record in a given TDL
// and Quality Band that will next be offered a shipment.
func NextTSPPerformanceInQualityBand(tx *pop.Connection, tdlID uuid.UUID, qualityBand int, shipment Shipment) (
	TransportationServiceProviderPerformance, error) {

	sql := `SELECT
			*
		FROM
			transportation_service_provider_performances
		WHERE
			traffic_distribution_list_id = $1
			AND
			quality_band = $2
			AND
			$3 BETWEEN performance_period_start AND performance_period_end
			AND
			$4 BETWEEN rate_cycle_start AND rate_cycle_end
		ORDER BY
			offer_count ASC,
			best_value_score DESC
		`

	tspp := TransportationServiceProviderPerformance{}
	err := tx.RawQuery(sql, tdlID, qualityBand, shipment.BookDate, shipment.RequestedPickupDate).First(&tspp)

	return tspp, err
}

// GatherNextEligibleTSPPerformances returns a map of QualityBands to their next eligible TSPPerformance.
func GatherNextEligibleTSPPerformances(tx *pop.Connection, tdlID uuid.UUID, shipment Shipment) (map[int]TransportationServiceProviderPerformance, error) {
	tspPerformances := make(map[int]TransportationServiceProviderPerformance)
	for _, qualityBand := range qualityBands {
		tspPerformance, err := NextTSPPerformanceInQualityBand(tx, tdlID, qualityBand, shipment)
		if err != nil {
			// We don't want the program to error out if Quality Bands don't have a TSPPerformance.
			//zap.S().Errorf("\tNo TSP returned for Quality Band: %d\n; See error: %s", qualityBand, err)
		} else {
			tspPerformances[qualityBand] = tspPerformance
		}
	}
	if len(tspPerformances) == 0 {
		return tspPerformances, fmt.Errorf("\tNo TSPPerformances found for TDL %s", tdlID)
	}
	return tspPerformances, nil
}

// NextEligibleTSPPerformance wraps GatherNextEligibleTSPPerformances and DetermineNextTSPPerformance.
func NextEligibleTSPPerformance(db *pop.Connection, tdlID uuid.UUID, shipment Shipment) (TransportationServiceProviderPerformance, error) {
	var tspPerformance TransportationServiceProviderPerformance
	tspPerformances, err := GatherNextEligibleTSPPerformances(db, tdlID, shipment)
	if err == nil {
		return SelectNextTSPPerformance(tspPerformances), nil
	}
	return tspPerformance, err
}

// SelectNextTSPPerformance returns the tspPerformance that is next to receive a shipment.
func SelectNextTSPPerformance(tspPerformances map[int]TransportationServiceProviderPerformance) TransportationServiceProviderPerformance {
	bands := sortedMapIntKeys(tspPerformances)
	// First time through, no rounds have yet occurred so rounds is set to the maximum rounds that have already occured.
	// Since the TSPs in quality band 1 will always have been offered the greatest number of shipments, we use that to calculate max.
	maxRounds := float64(tspPerformances[bands[0]].OfferCount) / float64(OffersPerQualityBand[bands[0]])
	previousRounds := math.Ceil(maxRounds)

	for _, band := range bands {
		tspPerformance := tspPerformances[band]
		rounds := float64(tspPerformance.OfferCount) / float64(OffersPerQualityBand[band])

		if rounds < previousRounds {
			return tspPerformance
		}
		previousRounds = rounds
	}

	// If we get all the way through, it means all of the TSPPerformances have had the
	// same number of offers and we should wrap around and assign the next offer to
	// the first quality band.
	return tspPerformances[bands[0]]
}

func sortedMapIntKeys(mapWithIntKeys map[int]TransportationServiceProviderPerformance) []int {
	keys := []int{}
	for key := range mapWithIntKeys {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return keys
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

// AssignQualityBandToTSPPerformance sets the QualityBand value for a TransportationServiceProviderPerformance.
func AssignQualityBandToTSPPerformance(db *pop.Connection, band int, id uuid.UUID) error {
	performance := TransportationServiceProviderPerformance{}
	if err := db.Find(&performance, id); err != nil {
		return err
	}
	performance.QualityBand = &band
	verrs, err := db.ValidateAndUpdate(&performance)
	if err != nil {
		return err
	} else if verrs.Count() > 0 {
		return errors.New("could not update quality band")
	}
	return nil
}

// IncrementTSPPerformanceOfferCount increments the offer_count column by 1 and validates.
func IncrementTSPPerformanceOfferCount(db *pop.Connection, tspPerformanceID uuid.UUID) error {
	var tspPerformance TransportationServiceProviderPerformance
	if err := db.Find(&tspPerformance, tspPerformanceID); err != nil {
		return err
	}
	tspPerformance.OfferCount++
	validationErr, databaseErr := db.ValidateAndSave(&tspPerformance)
	if databaseErr != nil {
		return databaseErr
	} else if validationErr.HasAny() {
		return fmt.Errorf("Validation failure: %s", validationErr)
	}
	return nil
}

// GetRateCycle returns the start date and end dates for a rate cycle of the
// given year and season (peak/non-peak).
func GetRateCycle(year int, peak bool) (start time.Time, end time.Time) {
	if peak {
		start = time.Date(year, time.May, 15, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.October, 1, 0, 0, 0, 0, time.UTC)
	} else {
		start = time.Date(year, time.October, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(year+1, time.May, 15, 0, 0, 0, 0, time.UTC)
	}

	return start, end
}

// DateIsPeakRateCycle determines if a given date is within the a peak rate
// cycle or not. This is the authoritative source on rate cycles.
// Peak rate cycles are: May 15th - September 30th, inclusive.
func DateIsPeakRateCycle(t time.Time) bool {
	year, _, _ := t.Date()

	peakStartInc, peakEndExcl := GetRateCycle(year, true)

	if (t.After(peakStartInc) || t.Equal(peakStartInc)) && t.Before(peakEndExcl) {
		fmt.Println(" true")
		return true
	}

	fmt.Println(" false")
	return false
}
