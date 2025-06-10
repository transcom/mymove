package dates

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/us"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// NewUSCalendar returns a new Calendar object initialized with standard US Federal Holidays.
func NewUSCalendar() *cal.BusinessCalendar {
	// NOTE: For now, we are returning a new calendar object for each call.  Could consider
	// caching this in the future.
	usCalendar := cal.NewBusinessCalendar()
	usCalendar.AddHoliday(us.Holidays...)
	return usCalendar
}

// NewCalendar returns a new calendar object based on the provided country code.
func NewCalendar(appCtx appcontext.AppContext, countryCode string) (*cal.BusinessCalendar, *models.Country, error) {

	if appCtx == nil {
		return nil, nil, errors.New("app context is nil")
	}

	db := appCtx.DB()
	if db == nil {
		return nil, nil, errors.New("database connection is nil")
	}

	if len(countryCode) != 2 {
		return nil, nil, errors.New("countryCode should be precisely two characters")
	}

	// Grab the country with its holidays and weekends
	var country models.Country
	err := db.Where("country = ?", countryCode).EagerPreload("Holidays", "Weekends").First(&country)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, nil, errors.Wrap(models.ErrFetchNotFound, "the country code provided was not found")
		}
		return nil, nil, errors.Wrap(err, "failed to fetch country data")
	}

	calendar := cal.NewBusinessCalendar()

	// Populate holidays
	setHolidays(calendar, &country.Holidays)

	// Set workdays
	setWorkdays(calendar, &country.Weekends)

	return calendar, &country, nil
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

func setHolidays(calendar *cal.BusinessCalendar, holidays *models.CountryHolidays) {
	if calendar != nil && holidays != nil {
		for _, holiday := range *holidays {
			year, month, day := holiday.ObservationDate.Date()
			calendar.AddHoliday(&cal.Holiday{
				Name:      holiday.HolidayName,
				Month:     month,
				Day:       day,
				StartYear: year,
				EndYear:   year,
				Func:      cal.CalcDayOfMonth,
			})
		}
	}
}

func setWorkdays(calendar *cal.BusinessCalendar, weekends *models.CountryWeekend) {

	if calendar != nil && weekends != nil {

		// The default workdays are Monday - Friday, so we'll check for deviation from that.
		if weekends.IsMondayWeekend {
			calendar.SetWorkday(time.Monday, false)
		}
		if weekends.IsTuesdayWeekend {
			calendar.SetWorkday(time.Tuesday, false)
		}
		if weekends.IsWednesdayWeekend {
			calendar.SetWorkday(time.Wednesday, false)
		}
		if weekends.IsThursdayWeekend {
			calendar.SetWorkday(time.Thursday, false)
		}
		if weekends.IsFridayWeekend {
			calendar.SetWorkday(time.Friday, false)
		}

		// Saturday and Sunday are default weekend days, so we'll check for deviation from that.
		if !weekends.IsSaturdayWeekend {
			calendar.SetWorkday(time.Saturday, true)
		}
		if !weekends.IsSundayWeekend {
			calendar.SetWorkday(time.Sunday, true)
		}
	}
}
