package ghcrateengine

import (
	"github.com/gobuffalo/pop"
)

const MinDomesticPerWeightPounds = 500
const MinInternationalServicesWeightPounds = 500
const MinUBServicesWeightPounds = 300

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

// Use this function to apply the escalation rate
// when all sub-costs have been calculated
// func lookupContractYearEscalation(db *pop.Connection, moveDate time.Time, contractCode string) float64 {
// 	// TODO: look up contract using contractCode and move Date
// 	// select escalation from re_contracts innerjoin re_contract_years on re_contract_years.id = re_contract_years.contract_id
// 	// where contract_code = contractCode and moveDate is between contract_years.start_date and contract_years.end_date
// 	stubEscalation := 1.02
// 	return stubEscalation
// }

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