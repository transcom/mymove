package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/dates"
	calendarop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/calendar"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// ShowAvailableMoveDatesHandler returns the available move dates starting at a given date.
type ShowAvailableMoveDatesHandler struct {
	handlers.HandlerConfig
}

// Handle returns the available move dates.
func (h ShowAvailableMoveDatesHandler) Handle(params calendarop.ShowAvailableMoveDatesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(_ appcontext.AppContext) (middleware.Responder, error) {

			startDate := time.Time(params.StartDate)

			var availableMoveDatesPayload internalmessages.AvailableMoveDates
			availableMoveDatesPayload.StartDate = handlers.FmtDate(startDate)

			var datesPayload []strfmt.Date
			const daysToCheckAfterStartDate = 90
			const shortFuseTotalDays = 5
			daysChecked := 0
			shortFuseDaysFound := 0

			usCalendar := dates.NewUSCalendar()
			firstPossibleDate := startDate.AddDate(0, 0, 1) // We never include the start date.
			for d := firstPossibleDate; daysChecked < daysToCheckAfterStartDate; d = d.AddDate(0, 0, 1) {
				if usCalendar.IsWorkday(d) {
					if shortFuseDaysFound < shortFuseTotalDays {
						shortFuseDaysFound++
					} else {
						datesPayload = append(datesPayload, strfmt.Date(d))
					}
				}
				daysChecked++
			}
			availableMoveDatesPayload.Available = datesPayload

			return calendarop.NewShowAvailableMoveDatesOK().WithPayload(&availableMoveDatesPayload), nil
		})
}

type IsDateWeekendHolidayHandler struct {
	handlers.HandlerConfig
	services.DateSelectionChecker
}

func (h IsDateWeekendHolidayHandler) Handle(params calendarop.IsDateWeekendHolidayParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			date := time.Time(params.Date)
			info, err := h.DateSelectionChecker.IsDateWeekendHoliday(appCtx, params.CountryCode, date)
			if err != nil {
				return calendarop.NewIsDateWeekendHolidayInternalServerError().WithPayload(nil), nil
			}
			var isDateWeekendHolidayInfo internalmessages.IsDateWeekendHolidayInfo
			isDateWeekendHolidayInfo.CountryCode = &info.CountryCode
			isDateWeekendHolidayInfo.CountryName = &info.CountryName
			isDateWeekendHolidayInfo.Date = handlers.FmtDate(info.Date)
			isDateWeekendHolidayInfo.IsWeekend = &info.IsWeekend
			isDateWeekendHolidayInfo.IsHoliday = &info.IsHoliday
			return calendarop.NewIsDateWeekendHolidayOK().WithPayload(&isDateWeekendHolidayInfo), nil
		})
}
