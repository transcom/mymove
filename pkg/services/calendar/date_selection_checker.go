package calendar

import (
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/services"
)

type dateSelectionChecker struct{}

func NewDateSelectionChecker() services.DateSelectionChecker {
	return &dateSelectionChecker{}
}

func (g *dateSelectionChecker) IsDateWeekendHoliday(appCtx appcontext.AppContext, countryCode string, date time.Time) (*services.IsDateWeekendHolidayInfo, error) {

	calendar, country, err := dates.NewCalendar(appCtx, countryCode)
	if err != nil {
		return nil, err
	}

	isHoliday, _, _ := calendar.IsHoliday(date)
	isWeekend := country.Weekends.IsWeekend(date)

	var isDateWeekendHolidayInfo = services.IsDateWeekendHolidayInfo{}
	isDateWeekendHolidayInfo.CountryCode = countryCode
	isDateWeekendHolidayInfo.CountryName = cases.Title(language.English).String(strings.ToLower(country.CountryName))
	isDateWeekendHolidayInfo.Date = date
	isDateWeekendHolidayInfo.IsWeekend = isWeekend
	isDateWeekendHolidayInfo.IsHoliday = isHoliday
	return &isDateWeekendHolidayInfo, nil
}
