package testdatagen

import (
	"time"
)

var TestYear = TestYear

// Peak Rate Cycle test dates.
var PeakRateCycleStart = time.Date(TestYear, time.May, 15, 0, 0, 0, 0, time.UTC)
var PeakRateCycleEnd = time.Date(TestYear, time.October, 1, 0, 0, 0, 0, time.UTC)
var DateInsidePeakRateCycle = time.Date(TestYear, time.May, 16, 0, 0, 0, 0, time.UTC)
var DateOutsidePeakRateCycle = time.Date(TestYear, time.October, 1, 0, 0, 0, 0, time.UTC)

// Non-Peak Rate Cycle test dates.
var NonPeakRateCycleStart = time.Date(TestYear, time.October, 1, 0, 0, 0, 0, time.UTC)
var NonPeakRateCycleEnd = time.Date(TestYear+1, time.May, 15, 0, 0, 0, 0, time.UTC)
var DateInsideNonPeakRateCycle = time.Date(TestYear, time.October, 2, 0, 0, 0, 0, time.UTC)
var DateOutsideNonPeakRateCycle = time.Date(TestYear+1, time.May, 16, 0, 0, 0, 0, time.UTC)

// PerformancePeriodStart is the first day of the first performance period
var PerformancePeriodStart = time.Date(TestYear, time.May, 15, 0, 0, 0, 0, time.UTC)

// PerformancePeriodEnd is the last day of the first performance period
var PerformancePeriodEnd = time.Date(TestYear, time.July, 31, 0, 0, 0, 0, time.UTC)

// DateInsidePerformancePeriod is within the performance period defined by
// PerformancePeriodStart and PerformancePeriodEnd.
var DateInsidePerformancePeriod = time.Date(TestYear, time.May, 16, 0, 0, 0, 0, time.UTC)

// DateOutsidePerformancePeriod is after the performance period defined by
// PerformancePeriodStart and PerformancePeriodEnd.
var DateOutsidePerformancePeriod = time.Date(TestYear, time.August, 1, 0, 0, 0, 0, time.UTC)
