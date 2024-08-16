package calendar

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/services"
)

type dateSelectionChecker struct{}

func NewDateSelectionChecker() services.DateSelectionChecker {
	return &dateSelectionChecker{}
}

func (g *dateSelectionChecker) IsDateWeekendHoliday(appCtx appcontext.AppContext, countryCode string, date time.Time) (*services.IsDateWeekendHolidayInfo, error) {
	// Assume for now invocation is US based only
	// TODO: query TRDM data to determine if date is weekend/holiday for particular country for international moves
	var calendar = dates.NewUSCalendar()
	isHoliday, _, _ := calendar.IsHoliday(date)
	var isDateWeekendHolidayInfo = services.IsDateWeekendHolidayInfo{}
	isDateWeekendHolidayInfo.CountryCode = countryCode
	// TODO - look up country name. For now return US.
	isDateWeekendHolidayInfo.CountryName = "United States"
	isDateWeekendHolidayInfo.Date = date
	isDateWeekendHolidayInfo.IsWeekend = date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
	isDateWeekendHolidayInfo.IsHoliday = isHoliday
	return &isDateWeekendHolidayInfo, nil
}
