package models

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/unit"
)

var qualityBands = []int{1, 2, 3, 4}

// OffersPerQualityBand is a map of the number of shipments to be offered per round to each quality band
// TODO: change these back to [5, 3, 2, 1] after the B&M pilot
var OffersPerQualityBand = map[int]int{
	1: 1,
	2: 1,
	3: 1,
	4: 1,
}

// TransportationServiceProviderPerformance is a combination of all TSP
// performance metrics (BVS, Quality Band) for a performance period.
type TransportationServiceProviderPerformance struct {
	ID                              uuid.UUID                     `json:"id" db:"id"`
	CreatedAt                       time.Time                     `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time                     `json:"updated_at" db:"updated_at"`
	PerformancePeriodStart          time.Time                     `json:"performance_period_start" db:"performance_period_start"`
	PerformancePeriodEnd            time.Time                     `json:"performance_period_end" db:"performance_period_end"`
	RateCycleStart                  time.Time                     `json:"rate_cycle_start" db:"rate_cycle_start"`
	RateCycleEnd                    time.Time                     `json:"rate_cycle_end" db:"rate_cycle_end"`
	TrafficDistributionListID       uuid.UUID                     `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	TrafficDistributionList         TrafficDistributionList       `belongs_to:"traffic_distribution_list" fk_id:"traffic_distribution_list_id"`
	TransportationServiceProviderID uuid.UUID                     `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	TransportationServiceProvider   TransportationServiceProvider `belongs_to:"transportation_service_provider" fk_id:"transportation_service_provider_id"`
	QualityBand                     *int                          `json:"quality_band" db:"quality_band"`
	BestValueScore                  float64                       `json:"best_value_score" db:"best_value_score"`
	LinehaulRate                    unit.DiscountRate             `json:"linehaul_rate" db:"linehaul_rate"`
	SITRate                         unit.DiscountRate             `json:"sit_rate" db:"sit_rate"`
	OfferCount                      int                           `json:"offer_count" db:"offer_count"`
}

// TransportationServiceProviderPerformances is a handy type for multiple TransportationServiceProviderPerformance structs
type TransportationServiceProviderPerformances []TransportationServiceProviderPerformance

// TSPPerformanceGroup contains the fields required to uniquely identify a TransportationServiceProviderPerformances
// grouping for quality band assignment (currently done in the award queue).
type TSPPerformanceGroup struct {
	TrafficDistributionListID uuid.UUID
	PerformancePeriodStart    time.Time
	PerformancePeriodEnd      time.Time
	RateCycleStart            time.Time
	RateCycleEnd              time.Time
}

// TSPPerformanceGroups is a handy type for multiple TSPPerformanceGroup structs
type TSPPerformanceGroups []TSPPerformanceGroup

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

		// Best Value Scores can range from 0 - 100, with up to four decimal places, as defined
		// in DTR403. See page 7 of https://www.ustranscom.mil/dtr/part-iv/dtr-part-4-403.pdf
		&validators.IntIsGreaterThan{Field: int(t.BestValueScore), Name: "BestValueScore", Compared: -1},
		&validators.IntIsLessThan{Field: int(t.BestValueScore), Name: "BestValueScore", Compared: 101},

		&DiscountRateIsValid{Field: t.LinehaulRate, Name: "LinehaulRate"},
		&DiscountRateIsValid{Field: t.SITRate, Name: "SITRate"},
	), nil
}

// MarshalLogObject is required to be able to zap.Object log a TSPP
func (t *TransportationServiceProviderPerformance) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("ID", t.ID.String())
	encoder.AddTime("PerformancePeriodStart", t.PerformancePeriodStart)
	encoder.AddTime("PerformancePeriodEnd", t.PerformancePeriodEnd)
	encoder.AddTime("RateCycleStart", t.RateCycleStart)
	encoder.AddTime("RateCycleEnd", t.RateCycleEnd)
	encoder.AddString("TrafficDistributionListID", t.TrafficDistributionListID.String())
	encoder.AddString("TransportationServiceProviderID", t.TransportationServiceProviderID.String())
	if t.QualityBand != nil {
		encoder.AddInt("QualityBand", *t.QualityBand)
	}
	encoder.AddFloat64("BestValueScore", t.BestValueScore)
	encoder.AddFloat64("LinehaulRate", t.LinehaulRate.Float64())
	encoder.AddFloat64("SITRate", t.SITRate.Float64())
	encoder.AddInt("OfferCount", t.OfferCount)

	return nil
}

// NextTSPPerformanceInQualityBand returns the TSP performance record in a given TDL
// and Quality Band that will next be offered a shipment.
func NextTSPPerformanceInQualityBand(tx *pop.Connection, tdlID uuid.UUID,
	qualityBand int, bookDate time.Time, requestedPickupDate time.Time) (
	TransportationServiceProviderPerformance, error) {

	sql := `SELECT
			tspp.*
		FROM
			transportation_service_provider_performances AS tspp
		LEFT JOIN
			transportation_service_providers AS tsp ON
				tspp.transportation_service_provider_id = tsp.id
		WHERE
			tspp.traffic_distribution_list_id = $1
			AND
			tspp.quality_band = $2
			AND
			$3 BETWEEN tspp.performance_period_start AND tspp.performance_period_end
			AND
			$4 BETWEEN tspp.rate_cycle_start AND tspp.rate_cycle_end
			AND
			tsp.enrolled = true
		ORDER BY
			offer_count ASC,
			best_value_score DESC
		`
	// Note: For PPM estimates, we ensure we have a tiebreaker that always returns the same TSPP
	// record in case multiple records match in the query above.  We may want to adjust the award
	// queue for consistency if we start doing HHGs again.  For more information, see:
	// https://docs.google.com/document/d/1T-KYb7BGNWpybkz-LrLGFRfWyKXhAD2w4fwJOBHko5A/edit#

	tspp := TransportationServiceProviderPerformance{}
	err := tx.RawQuery(sql, tdlID, qualityBand, bookDate, requestedPickupDate).First(&tspp)

	return tspp, err
}

// GatherNextEligibleTSPPerformances returns a map of QualityBands to their next eligible TSPPerformance.
func GatherNextEligibleTSPPerformances(tx *pop.Connection, tdlID uuid.UUID, bookDate time.Time, requestedPickupDate time.Time) (map[int]TransportationServiceProviderPerformance, error) {
	tspPerformances := make(map[int]TransportationServiceProviderPerformance)
	qualityBandsWithoutTSPs := 0

	for _, qualityBand := range qualityBands {
		tspPerformance, err := NextTSPPerformanceInQualityBand(tx, tdlID, qualityBand, bookDate, requestedPickupDate)
		if err != nil {
			if err.Error() == RecordNotFoundErrorString {
				// Some quality bands might not have TSPs, and that's OK. We
				// just need to make sure SOME quality bands have TSPs.
				qualityBandsWithoutTSPs++
			} else {
				return tspPerformances, err
			}
		} else {
			tspPerformances[qualityBand] = tspPerformance
		}
	}
	if qualityBandsWithoutTSPs >= len(qualityBands) {
		return tspPerformances, fmt.Errorf("Could not find any TSPs to fill quality bands in TDL: %s", tdlID)
	}
	return tspPerformances, nil
}

// NextEligibleTSPPerformance wraps GatherNextEligibleTSPPerformances and DetermineNextTSPPerformance.
func NextEligibleTSPPerformance(db *pop.Connection, tdlID uuid.UUID, bookDate time.Time, requestedPickupDate time.Time) (TransportationServiceProviderPerformance, error) {
	var tspPerformance TransportationServiceProviderPerformance
	tspPerformances, err := GatherNextEligibleTSPPerformances(db, tdlID, bookDate, requestedPickupDate)
	if err == nil {
		return SelectNextTSPPerformance(tspPerformances), nil
	}
	return tspPerformance, err
}

// SelectNextTSPPerformance returns the tspPerformance that is next to receive a shipment.
func SelectNextTSPPerformance(tspPerformances map[int]TransportationServiceProviderPerformance) TransportationServiceProviderPerformance {
	bands := sortedMapIntKeys(tspPerformances)
	// First time through, no rounds have yet occurred so rounds is set to the maximum rounds that have already occurred.
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

// FetchTSPPerformancesForQualityBandAssignment returns TSPPs in the given TSPP grouping in the order
// that they should be assigned quality bands.
func FetchTSPPerformancesForQualityBandAssignment(tx *pop.Connection, perfGroup TSPPerformanceGroup, mps float64) (TransportationServiceProviderPerformances, error) {
	var perfs TransportationServiceProviderPerformances
	err := tx.
		Select("transportation_service_provider_performances.*").
		Join("transportation_service_providers AS tsp", "tsp.id = transportation_service_provider_performances.transportation_service_provider_id").
		Where("traffic_distribution_list_id = ?", perfGroup.TrafficDistributionListID).
		Where("performance_period_start = ?", perfGroup.PerformancePeriodStart).
		Where("performance_period_end = ?", perfGroup.PerformancePeriodEnd).
		Where("rate_cycle_start = ?", perfGroup.RateCycleStart).
		Where("rate_cycle_end = ?", perfGroup.RateCycleEnd).
		Where("best_value_score > ?", mps).
		Where("enrolled = true").
		Order("best_value_score DESC").
		All(&perfs)
	// Note: For PPM estimates, we ensure we have a tiebreaker that always returns the same TSPP
	// record in case multiple records match in the query above.  We may want to adjust the award
	// queue for consistency if we start doing HHGs again.  For more information, see:
	// https://docs.google.com/document/d/1T-KYb7BGNWpybkz-LrLGFRfWyKXhAD2w4fwJOBHko5A/edit#

	return perfs, err
}

// FetchUnbandedTSPPerformanceGroups gets all groupings of TSPPs that have at least one entry with
// an unassigned quality band.
func FetchUnbandedTSPPerformanceGroups(db *pop.Connection) (TSPPerformanceGroups, error) {
	var perfs TransportationServiceProviderPerformances
	err := db.
		Select("traffic_distribution_list_id", "performance_period_start", "performance_period_end", "rate_cycle_start", "rate_cycle_end").
		Join("transportation_service_providers AS tsp", "tsp.id = transportation_service_provider_performances.transportation_service_provider_id").
		Where("quality_band IS NULL").
		Where("enrolled = true").
		GroupBy("traffic_distribution_list_id", "performance_period_start", "performance_period_end", "rate_cycle_start", "rate_cycle_end").
		Order("traffic_distribution_list_id, performance_period_start, rate_cycle_start").
		All(&perfs)

	perfGroups := make(TSPPerformanceGroups, len(perfs))
	for i, perf := range perfs {
		perfGroups[i] = TSPPerformanceGroup{
			TrafficDistributionListID: perf.TrafficDistributionListID,
			PerformancePeriodStart:    perf.PerformancePeriodStart,
			PerformancePeriodEnd:      perf.PerformancePeriodEnd,
			RateCycleStart:            perf.RateCycleStart,
			RateCycleEnd:              perf.RateCycleEnd,
		}
	}

	return perfGroups, err
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
// It returns the updated TSPPerformance record.
func IncrementTSPPerformanceOfferCount(db *pop.Connection, tspPerformanceID uuid.UUID) (TransportationServiceProviderPerformance, error) {
	var tspPerformance TransportationServiceProviderPerformance
	if err := db.Find(&tspPerformance, tspPerformanceID); err != nil {
		return tspPerformance, err
	}
	tspPerformance.OfferCount++
	validationErr, databaseErr := db.ValidateAndSave(&tspPerformance)
	if databaseErr != nil {
		return tspPerformance, databaseErr
	} else if validationErr.HasAny() {
		return tspPerformance, fmt.Errorf("Validation failure: %s", validationErr)
	}
	return tspPerformance, nil
}

// GetRateCycle returns the start date and end dates for a rate cycle of the
// given year and season (peak/non-peak), inclusive.
func GetRateCycle(year int, peak bool) (start time.Time, end time.Time) {
	if peak {
		start = time.Date(year, time.May, 15, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.September, 30, 0, 0, 0, 0, time.UTC)
	} else {
		start = time.Date(year, time.October, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(year+1, time.May, 14, 0, 0, 0, 0, time.UTC)
	}

	return start, end
}

// FetchDiscountRates returns the discount linehaul and SIT rates for the TSP with the highest
// BVS during the specified date, limited to those TSPs in the channel defined by the
// originZip and destinationZip.  In case of more than one TSP having the same highest BVS score,
// we return the one whose TSPP ID comes first alphabetically.
func FetchDiscountRates(db *pop.Connection, originZip string, destinationZip string, cos string, date time.Time) (linehaulDiscount unit.DiscountRate, sitDiscount unit.DiscountRate, err error) {
	// the discount rates are based on BVS rate data which we no longer insert into our system.
	// this rate data is going to change in outcome 5 so we don't really care what these values are.
	// hard code discount rates so ppm pricing doesn't fail w/ no data
	return .505, .5, nil
}
