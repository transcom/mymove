package testdatagen

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// GetDateInRateCycle returns a date that is guaranteed to be in the requested
// year and Peak/NonPeak season.
func GetDateInRateCycle(year int, peak bool) time.Time {
	start, end := models.GetRateCycle(year, peak)

	center := end.Sub(start) / 2

	return start.Add(center)
}

// MakeTSPPerformance makes a single transportation service provider record.
func MakeTSPPerformance(db *pop.Connection, assertions Assertions) models.TransportationServiceProviderPerformance {

	var tsp models.TransportationServiceProvider
	id := assertions.TransportationServiceProviderPerformance.TransportationServiceProviderID
	if id == uuid.Nil {
		tsp = MakeDefaultTSP(db)
	} else {
		tsp = assertions.TransportationServiceProviderPerformance.TransportationServiceProvider
	}

	var tdl models.TrafficDistributionList
	id = assertions.TransportationServiceProviderPerformance.TrafficDistributionListID
	if id == uuid.Nil {
		tdl = MakeDefaultTDL(db)
	} else {
		tdl = assertions.TransportationServiceProviderPerformance.TrafficDistributionList
	}

	qualityBand := 1
	tspp := models.TransportationServiceProviderPerformance{
		TransportationServiceProvider:   tsp,
		TransportationServiceProviderID: tsp.ID,

		PerformancePeriodStart:    PerformancePeriodStart,
		PerformancePeriodEnd:      PerformancePeriodEnd,
		RateCycleStart:            PeakRateCycleStart,
		RateCycleEnd:              PeakRateCycleEnd,
		TrafficDistributionListID: tdl.ID,
		QualityBand:               &qualityBand,
		BestValueScore:            0.88,
		OfferCount:                0,
		LinehaulRate:              0.34,
		SITRate:                   0.45,
	}

	mergeModels(&tspp, assertions.TransportationServiceProviderPerformance)

	verrs, err := db.ValidateAndCreate(&tspp)
	if verrs.HasAny() {
		err = fmt.Errorf("TSPP validation errors: %v", verrs)
	}
	if err != nil {
		log.Panic(err)
	}

	return tspp
}

// MakeDefaultTSPPerformance makes a TransportationServiceProviderPerformance with default values
func MakeDefaultTSPPerformance(db *pop.Connection) models.TransportationServiceProviderPerformance {
	return MakeTSPPerformance(db, Assertions{})
}

// MakeTSPPerformanceDeprecated makes a single best_value_score record
//
// Deprecated: Use MakeTSPPErformance or MakeDEfaultTSPPERformance instead.
func MakeTSPPerformanceDeprecated(db *pop.Connection,
	tsp models.TransportationServiceProvider,
	tdl models.TrafficDistributionList,
	qualityBand *int,
	score float64,
	offerCount int,
	linehaulDiscountRate unit.DiscountRate,
	SITDiscountRate unit.DiscountRate) (models.TransportationServiceProviderPerformance, error) {

	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          PerformancePeriodStart,
		PerformancePeriodEnd:            PerformancePeriodEnd,
		RateCycleStart:                  PeakRateCycleStart,
		RateCycleEnd:                    PeakRateCycleEnd,
		TransportationServiceProviderID: tsp.ID,
		TrafficDistributionListID:       tdl.ID,
		QualityBand:                     qualityBand,
		BestValueScore:                  score,
		OfferCount:                      offerCount,
		LinehaulRate:                    linehaulDiscountRate,
		SITRate:                         SITDiscountRate,
	}

	verrs, err := db.ValidateAndCreate(&tspPerformance)
	if verrs.HasAny() {
		err = fmt.Errorf("TSP Performance validation errors: %v", verrs)
	}
	if err != nil {
		log.Panic(err)
	}

	return tspPerformance, err
}
