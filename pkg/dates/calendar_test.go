package dates

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type CalendarSuite struct {
	*testingsuite.PopTestSuite
}

const defaultCountryCode = "US"

func TestCalendarSuite(t *testing.T) {

	hs := &CalendarSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func TestNextWorkday(t *testing.T) {
	var dateTests = []struct {
		name string
		in   time.Time
		out  time.Time
	}{
		{
			"No weekend or holiday",
			time.Date(2019, 1, 24, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 1, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			"Holiday",
			time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 12, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			"Weekend",
			time.Date(2019, 1, 25, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 1, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			"Weekend and holiday",
			time.Date(2019, 1, 18, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 1, 22, 0, 0, 0, 0, time.UTC),
		},
	}
	cal := NewUSCalendar()
	for _, dt := range dateTests {
		t.Run(dt.name, func(t *testing.T) {
			nextDate := NextWorkday(*cal, dt.in)
			if nextDate != dt.out {
				t.Fatalf("Actual date: %v is not equal to expected date: %v", nextDate, dt.out)
			}
		})
	}
}

func TestNextNonWorkday(t *testing.T) {
	var dateTests = []struct {
		name string
		in   time.Time
		out  time.Time
	}{
		{
			"Saturday after weekday",
			time.Date(2019, 1, 24, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 1, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			"Saturday after Sunday",
			time.Date(2019, 1, 27, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 2, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			"Holiday after weekday",
			time.Date(2019, 12, 23, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC),
		},
	}
	cal := NewUSCalendar()
	for _, dt := range dateTests {
		t.Run(dt.name, func(t *testing.T) {
			nextDate := NextNonWorkday(*cal, dt.in)
			if nextDate != dt.out {
				t.Fatalf("Actual date: %v is not equal to expected date: %v", nextDate, dt.out)
			}
		})
	}
}

func (suite *DatesSuite) TestNewCalendar() {

	suite.Run("calendar creation fails with nil app context", func() {
		var appCtx appcontext.AppContext = nil
		calendar, country, err := NewCalendar(appCtx, defaultCountryCode)
		suite.Nil(calendar, "calendar should be nil when app context is nil")
		suite.Nil(country, "country should be nil when app context is nil")
		suite.Error(err, "Expected an error when app context is nil")
		suite.Contains(err.Error(), "app context is nil", "Error should mention nil app context")
	})

	suite.Run("calendar creation fails with nil db", func() {
		var appCtx appcontext.AppContext = appcontext.NewAppContext(nil, nil, nil, nil)
		calendar, country, err := NewCalendar(appCtx, defaultCountryCode)
		suite.Nil(calendar, "calendar should be nil when db is nil")
		suite.Nil(country, "country should be nil when db is nil")
		suite.Error(err, "Expected an error when db is nil")
		suite.Contains(err.Error(), "database connection is nil", "Error should mention nil db")
	})

	failureTestCases := []struct {
		name        string
		countryCode string
		expectError string
	}{
		{
			name:        "empty country code",
			countryCode: "",
		},
		{
			name:        "single character country code",
			countryCode: "a",
		},
		{
			name:        "3 character country code",
			countryCode: "123",
		},
	}

	for _, tc := range failureTestCases {
		suite.Run("calendar creation fails for "+tc.name, func() {
			calendar, country, err := NewCalendar(suite.AppContextForTest(), tc.countryCode)

			suite.Nil(calendar, "calendar should be nil when countryCode is invalid")
			suite.Nil(country, "country should be nil when countryCode is invalid")
			suite.Error(err, "Expected an error when countryCode is invalid")
			suite.Contains(err.Error(), "countryCode should be precisely two characters", "Error should mention invalid countryCode")
		})
	}

	suite.Run("calendar creation fails when country code is not found", func() {
		calendar, country, err := NewCalendar(suite.AppContextForTest(), "XX")
		suite.Nil(calendar, "calendar should be nil when the country code is not found")
		suite.Nil(country, "country should be nil when the country code is not found")
		suite.Error(err, "Expected an error when the country code is not found")
		suite.Contains(err.Error(), "the country code provided was not found", "Error should mention that the country code provided was not found")
	})

	successTestCases := []struct {
		name        string
		countryCode string
	}{
		{
			countryCode: "US",
		},
		{
			countryCode: "IN",
		},
		{
			countryCode: "KW",
		},
	}

	for _, tc := range successTestCases {
		suite.Run("calendar creation succeeds for valid country code "+tc.countryCode, func() {
			calendar, country, err := NewCalendar(suite.AppContextForTest(), tc.countryCode)

			suite.NotNil(calendar, "calendar should not be nil")
			suite.NotNil(country, "country should not be nil")
			suite.NoError(err, "should not error when passed a valid country code")

			// Verify that the calendar reflects all the holidays
			for _, holiday := range country.Holidays {
				isHoliday, _, _ := calendar.IsHoliday(holiday.ObservationDate)
				suite.True(isHoliday, "calendar should correctly identify holiday")
			}

			// loop through every day from 6 months ago until 6 months from now to get some coverage
			start := time.Now().AddDate(0, -6, 0)
			end := time.Now().AddDate(0, 6, 0)
			for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {

				/*
				 First check if the day is a holiday. If if is a holiday, it should not be considered a workday.
				 If it's not a holiday, check if it falls on a weekend. Any day that does not fall on a holiday or weekend should be considered a workday.
				*/
				isHoliday, _, _ := calendar.IsHoliday(day)

				if isHoliday {
					suite.Equal(false, calendar.IsWorkday(day), "holiday should not be considered a workday")
				} else {
					var isWeekend bool
					isWorkday := calendar.IsWorkday(day)
					switch day.Weekday() {
					case time.Monday:
						isWeekend = country.Weekends.IsMondayWeekend
					case time.Tuesday:
						isWeekend = country.Weekends.IsTuesdayWeekend
					case time.Wednesday:
						isWeekend = country.Weekends.IsWednesdayWeekend
					case time.Thursday:
						isWeekend = country.Weekends.IsThursdayWeekend
					case time.Friday:
						isWeekend = country.Weekends.IsFridayWeekend
					case time.Saturday:
						isWeekend = country.Weekends.IsSaturdayWeekend
					case time.Sunday:
						isWeekend = country.Weekends.IsSundayWeekend
					}

					// Since we know the current day is not a holiday, then it should only be considered a workday if it doesn't fall on a weekend
					suite.NotEqual(isWorkday, isWeekend, "only non-holiday weekdays should be considered workdays")
				}
			}
		})
	}
}
