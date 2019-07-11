package testdatagen

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// GetDateInRateCycle returns a date that is guaranteed to be in the requested
// year and Peak/NonPeak season.
func GetDateInRateCycle(year int, peak bool) time.Time {
	start, end := models.GetRateCycle(year, peak)

	center := end.Sub(start) / 2

	return start.Add(center)
}

// MakeTSPPerformance makes a single transportation service provider record.
func MakeTSPPerformance(db *pop.Connection, assertions Assertions) (models.TransportationServiceProviderPerformance, error) {

	var tsp models.TransportationServiceProvider
	id := assertions.TransportationServiceProviderPerformance.TransportationServiceProviderID
	qualityBand := assertions.TransportationServiceProviderPerformance.QualityBand
	score := assertions.TransportationServiceProviderPerformance.BestValueScore
	offerCount := assertions.TransportationServiceProviderPerformance.OfferCount
	linehaulRate := assertions.TransportationServiceProviderPerformance.LinehaulRate
	sitRate := assertions.TransportationServiceProviderPerformance.SITRate
	quartile := assertions.TransportationServiceProviderPerformance.Quartile
	rank := assertions.TransportationServiceProviderPerformance.Rank
	surveyScore := assertions.TransportationServiceProviderPerformance.SurveyScore
	rateScore := assertions.TransportationServiceProviderPerformance.RateScore

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

	if score == 0 {
		score = 88
	}

	if linehaulRate == 0 {
		linehaulRate = 0.34
	}

	if sitRate == 0 {
		sitRate = 0.45
	}

	if quartile == 0 {
		quartile = 1
	}

	if rank == 0 {
		rank = 1
	}

	if surveyScore == 0 {
		surveyScore = 63
	}

	if rateScore == 0 {
		rateScore = 25
	}

	tspp := models.TransportationServiceProviderPerformance{
		TransportationServiceProvider:   tsp,
		TransportationServiceProviderID: tsp.ID,

		PerformancePeriodStart:    PerformancePeriodStart,
		PerformancePeriodEnd:      PerformancePeriodEnd,
		RateCycleStart:            PeakRateCycleStart,
		RateCycleEnd:              PeakRateCycleEnd,
		TrafficDistributionListID: tdl.ID,
		QualityBand:               qualityBand,
		BestValueScore:            score,
		OfferCount:                offerCount,
		LinehaulRate:              linehaulRate,
		SITRate:                   sitRate,
		Quartile:                  quartile,
		Rank:                      rank,
		SurveyScore:               surveyScore,
		RateScore:                 rateScore,
	}

	mergeModels(&tspp, assertions.TransportationServiceProviderPerformance)

	verrs, err := db.ValidateAndCreate(&tspp)
	if verrs.HasAny() {
		err = fmt.Errorf("TSPP validation errors: %v", verrs)
	}
	if err != nil {
		log.Panic(err)
	}

	return tspp, err
}

// MakeDefaultTSPPerformance makes a TransportationServiceProviderPerformance with default values
func MakeDefaultTSPPerformance(db *pop.Connection) (models.TransportationServiceProviderPerformance, error) {
	return MakeTSPPerformance(db, Assertions{})
}
