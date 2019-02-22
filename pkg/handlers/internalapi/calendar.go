package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/dates"
	calendarop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/calendar"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
)

// ShowAvailableMoveDatesHandler returns the available move dates starting at a given date.
type ShowAvailableMoveDatesHandler struct {
	handlers.HandlerContext
}

// Handle returns the available move dates.
func (h ShowAvailableMoveDatesHandler) Handle(params calendarop.ShowAvailableMoveDatesParams) middleware.Responder {
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

	return calendarop.NewShowAvailableMoveDatesOK().WithPayload(&availableMoveDatesPayload)
}
