package testdatagen

import (
	"time"
)

// PerformancePeriodStart is the first day of the first performance period in 2019.
var PerformancePeriodStart = time.Date(2019, time.May, 15, 0, 0, 0, 0, time.UTC)

// PerformancePeriodEnd is the last day of the first performance period in 2019.
var PerformancePeriodEnd = time.Date(2019, time.July, 31, 0, 0, 0, 0, time.UTC)

// DateInsidePerformancePeriod is within the performance period defined by
// PerformancePeriodStart and PerformancePeriodEnd.
var DateInsidePerformancePeriod = time.Date(2019, time.May, 16, 0, 0, 0, 0, time.UTC)

// DateOutsidePerformancePeriod is after the performance period defined by
// PerformancePeriodStart and PerformancePeriodEnd.
var DateOutsidePerformancePeriod = time.Date(2019, time.August, 1, 0, 0, 0, 0, time.UTC)
