package dates

import (
	"time"

	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/us"
)

// NewUSCalendar returns a new Calendar object initialized with standard US Federal Holidays.
func NewUSCalendar() *cal.BusinessCalendar {
	// NOTE: For now, we are returning a new calendar object for each call.  Could consider
	// caching this in the future.
	usCalendar := cal.NewBusinessCalendar()
	usCalendar.AddHoliday(us.Holidays...)
	return usCalendar
}

// NextWorkday returns the next workday after the given date, using the given calendar
func NextWorkday(cal cal.BusinessCalendar, date time.Time) time.Time {
	for {
		date = date.AddDate(0, 0, 1)
		if cal.IsWorkday(date) {
			return date
		}
	}
}

// NextNonWorkday returns the next weekend or holiday after the given date, using the given calendar
func NextNonWorkday(cal cal.BusinessCalendar, date time.Time) time.Time {
	for {
		date = date.AddDate(0, 0, 1)
		if !cal.IsWorkday(date) {
			return date
		}
	}
}
