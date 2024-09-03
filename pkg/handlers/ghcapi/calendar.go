package ghcapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	calendarop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/calendar"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

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
			var isDateWeekendHolidayInfo ghcmessages.IsDateWeekendHolidayInfo
			isDateWeekendHolidayInfo.CountryCode = &info.CountryCode
			isDateWeekendHolidayInfo.CountryName = &info.CountryName
			isDateWeekendHolidayInfo.Date = handlers.FmtDate(info.Date)
			isDateWeekendHolidayInfo.IsWeekend = &info.IsWeekend
			isDateWeekendHolidayInfo.IsHoliday = &info.IsHoliday
			return calendarop.NewIsDateWeekendHolidayOK().WithPayload(&isDateWeekendHolidayInfo), nil
		})
}
