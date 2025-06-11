package calendar

import (
	"errors"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/utils"
)

type dateSelectionChecker struct{}

func NewDateSelectionChecker() services.DateSelectionChecker {
	return &dateSelectionChecker{}
}

func (g *dateSelectionChecker) IsDateWeekendHoliday(appCtx appcontext.AppContext, countryCode string, date time.Time) (*services.IsDateWeekendHolidayInfo, error) {

	if appCtx == nil {
		return nil, errors.New("app context is nil")
	}

	if len(countryCode) != 2 {
		return nil, errors.New("countryCode should be precisely two characters")
	}

	if date.IsZero() {
		return nil, errors.New("date value is zero")
	}

	calendar, country, err := dates.NewCalendar(appCtx, countryCode)
	if err != nil {
		return nil, err
	}

	isHoliday, _, _ := calendar.IsHoliday(date)
	isWeekend := country.Weekends.IsWeekend(date)

	var isDateWeekendHolidayInfo = services.IsDateWeekendHolidayInfo{}
	isDateWeekendHolidayInfo.CountryCode = countryCode
	isDateWeekendHolidayInfo.CountryName = utils.ToTitleCase(country.CountryName)
	isDateWeekendHolidayInfo.Date = date
	isDateWeekendHolidayInfo.IsWeekend = isWeekend
	isDateWeekendHolidayInfo.IsHoliday = isHoliday
	return &isDateWeekendHolidayInfo, nil
}
