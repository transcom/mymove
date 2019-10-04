package ghcrateengine

import (
	"github.com/gobuffalo/pop"
)

// RateEngine encapsulates the TSP rate engine process
type GHCRateEngine struct {
	db     *pop.Connection
	logger Logger
}

func NewGHCRateEngine(db *pop.Connection, logger Logger) GHCRateEngine {
	return GHCRateEngine{
		db:     db,
		logger: logger,
	}
}

const MinDomesticPerWeightPounds = 500
const MinInternationalServiceWeightPounds = 500
const MinUBservicesWeightPounds = 300

//const PeakRateCycleStartMonth = time.May
//const PeakRateCycleStartDay = 15
//const PeakRateCycleEndMonth = time.September
//const PeakRateCycleEndDay = 30
//
//func getIsPeakPeriod(moveDate time.Time) bool {
//	isPeakPeriod := false
//	peakRateCycleStart := time.Date(moveDate.Year(), PeakRateCycleStartMonth, PeakRateCycleStartDay, 0, 0, 0, 0, time.UTC)
//	peakRateCycleEnd := time.Date(moveDate.Year(), PeakRateCycleEndMonth, PeakRateCycleEndDay, 0, 0, 0, 0, time.UTC)
//	if moveDate.After(peakRateCycleStart.AddDate(0, 0, -1)) && moveDate.Before(peakRateCycleEnd.AddDate(0, 0, 1)) {
//		isPeakPeriod = true
//	}
//	return isPeakPeriod
//}