package ghcrateengine

import (
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/unit"
)

// DefaultContractCode is the default contract code to assume for now
const DefaultContractCode = "TRUSS_TEST"

// minDomesticWeight is the minimum weight used in domestic calculations (weights below this are upgraded to the min)
const minDomesticWeight = unit.Pound(500)

// dateInYear represents a specific date in a year (without caring what year it is)
type dateInYear struct {
	month time.Month
	day   int
}

var (
	// The peak start/end dates
	peakStart = dateInYear{time.May, 15}
	peakEnd   = dateInYear{time.September, 30}
)

// addDate performs the same function as time.Time's AddDate, but ignores the year
func (d dateInYear) addDate(months int, days int) dateInYear {
	// Pick a year so we can use the time.Time functions (just about any year should work)
	fixedDate := time.Date(2019, d.month, d.day, 0, 0, 0, 0, time.UTC)
	newFixedDate := fixedDate.AddDate(0, months, days)
	return dateInYear{newFixedDate.Month(), newFixedDate.Day()}
}

// IsPeakPeriod determines if the given date is in the peak or non-peak part of the year
func IsPeakPeriod(date time.Time) bool {
	dateMonth := date.Month()
	dateDay := date.Day()

	// If the month is between the start/end (exclusive), definitely peak.
	if dateMonth > peakStart.month && dateMonth < peakEnd.month {
		return true
	}

	// If it's in the start month, check to see if it's in the peak part.
	if dateMonth == peakStart.month && dateDay >= peakStart.day {
		return true
	}

	// If it's in the end month, check to see if it's in the peak part.
	if dateMonth == peakEnd.month && dateDay <= peakEnd.day {
		return true
	}

	// Otherwise, it's non-peak
	return false
}

// centPriceAndEscalation is used to hold data returned by the database query
type centPriceAndEscalation struct {
	PriceCents           unit.Cents `db:"price_cents"`
	EscalationCompounded float64    `db:"escalation_compounded"`
}

// MarshalLogObject allows centPriceAndEscalation to be logged by zap
func (p centPriceAndEscalation) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("PriceCents", p.PriceCents.Int())
	encoder.AddFloat64("EscalationCompounded", p.EscalationCompounded)
	return nil
}
