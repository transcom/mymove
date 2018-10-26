package models

import "github.com/rickar/cal"

// NewUSCalendar returns a new Calendar object initialized with standard US Federal Holidays.
func NewUSCalendar() *cal.Calendar {
	// NOTE: For now, we are returning a new calendar object for each call.  Could consider
	// caching this in the future.
	usCalendar := cal.NewCalendar()
	cal.AddUsHolidays(usCalendar)
	return usCalendar
}
