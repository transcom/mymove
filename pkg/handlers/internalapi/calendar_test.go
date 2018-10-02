package internalapi

import (
	"github.com/go-openapi/strfmt"
	calendarop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/calendar"
	"github.com/transcom/mymove/pkg/handlers"
	"net/http/httptest"
	"time"
)

func (suite *HandlerSuite) TestShowUnavailableMoveDatesHandler() {
	req := httptest.NewRequest("GET", "/calendar/unavailable_move_dates", nil)

	params := calendarop.ShowUnavailableMoveDatesParams{
		HTTPRequest: req,
		StartDate:   strfmt.Date(time.Date(2018, 9, 26, 0, 0, 0, 0, time.UTC)),
	}

	unavailableDates := []strfmt.Date{
		strfmt.Date(time.Date(2018, 9, 26, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 9, 27, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 9, 28, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 9, 29, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 9, 30, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 1, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 2, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 3, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 6, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 7, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 8, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 13, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 14, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 20, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 21, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 27, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 28, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 3, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 4, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 10, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 11, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 12, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 17, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 18, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 22, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 24, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 25, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 1, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 2, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 8, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 9, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 22, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 23, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 25, 0, 0, 0, 0, time.UTC)),
	}

	showHandler := ShowUnavailableMoveDatesHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := showHandler.Handle(params)

	suite.IsType(&calendarop.ShowUnavailableMoveDatesOK{}, response)
	okResponse := response.(*calendarop.ShowUnavailableMoveDatesOK)

	suite.Equal(unavailableDates, okResponse.Payload)
}
