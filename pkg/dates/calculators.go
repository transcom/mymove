package dates

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rickar/cal"
)

// CreateFutureMoveDates generates a list of dates in the future
func CreateFutureMoveDates(startDate time.Time, numDays int, includeWeekendsAndHolidays bool, calendar *cal.Calendar) []time.Time {
	dates := make([]time.Time, 0, numDays)

	daysAdded := 0
	for d := startDate; daysAdded < numDays; d = d.AddDate(0, 0, 1) {
		if includeWeekendsAndHolidays || calendar.IsWorkday(d) {
			dates = append(dates, d)
			daysAdded++
		}
	}

	return dates
}

// CreatePastMoveDates generates a list of dates in the past
func CreatePastMoveDates(startDate time.Time, numDays int, includeWeekendsAndHolidays bool, calendar *cal.Calendar) []time.Time {
	dates := make([]time.Time, numDays)

	daysAdded := 0
	for d := startDate; daysAdded < numDays; d = d.AddDate(0, 0, -1) {
		if includeWeekendsAndHolidays || calendar.IsWorkday(d) {
			// Since we're working backwards, put dates at end of slice.
			dates[numDays-daysAdded-1] = d
			daysAdded++
		}
	}

	return dates
}

// CreateValidDatesBetweenTwoDates returns date range inclusive of startDate, exclusive of endDate (unless endDate is before or equal to startDate and allowEarlierEndDate)
func CreateValidDatesBetweenTwoDates(startDate time.Time, endDate time.Time, includeWeekendsAndHolidays bool, allowEarlierOrSameEndDate bool, calendar *cal.Calendar) ([]time.Time, error) {
	var dates []time.Time

	if startDate.After(endDate) || startDate == endDate {
		if allowEarlierOrSameEndDate == true {
			return dates, nil
		}
		return dates, errors.New("End date cannot be before or equal to start date")
	}

	dateToAdd := startDate

	for dateToAdd.Before(endDate) {
		if includeWeekendsAndHolidays || calendar.IsWorkday(dateToAdd) {
			dates = append(dates, dateToAdd)
		}
		dateToAdd = dateToAdd.AddDate(0, 0, 1)
	}
	return dates, nil
}

// NextValidMoveDate returns next subsequent non-holiday weekday
// This is mostly used for testing purposes
func NextValidMoveDate(d time.Time, calendar *cal.Calendar) time.Time {
	// Add days until we get a non-holiday weekday
	for !calendar.IsWorkday(d) {
		d = d.AddDate(0, 0, 1)
	}
	return d
}
