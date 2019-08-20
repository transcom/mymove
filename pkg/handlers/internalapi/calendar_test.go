package internalapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"

	calendarop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/calendar"
	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *HandlerSuite) TestShowAvailableMoveDatesHandler() {
	req := httptest.NewRequest("GET", "/calendar/available_move_dates", nil)

	startDate := strfmt.Date(time.Date(2018, 9, 27, 0, 0, 0, 0, time.UTC))
	params := calendarop.ShowAvailableMoveDatesParams{
		HTTPRequest: req,
		StartDate:   startDate,
	}

	availableDates := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 5, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 9, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 10, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 11, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 12, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 16, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 17, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 18, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 19, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 22, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 23, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 24, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 25, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 26, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 29, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 30, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 31, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 1, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 2, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 5, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 6, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 7, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 8, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 9, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 13, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 14, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 15, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 16, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 19, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 20, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 21, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 23, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 26, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 27, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 28, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 29, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 30, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 3, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 4, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 5, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 18, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 19, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 20, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 21, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 24, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 26, 0, 0, 0, 0, time.UTC)),
	}

	showHandler := ShowAvailableMoveDatesHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := showHandler.Handle(params)

	suite.IsType(&calendarop.ShowAvailableMoveDatesOK{}, response)
	okResponse := response.(*calendarop.ShowAvailableMoveDatesOK)

	suite.Equal(startDate, *okResponse.Payload.StartDate)
	suite.Equal(availableDates, okResponse.Payload.Available)
}
