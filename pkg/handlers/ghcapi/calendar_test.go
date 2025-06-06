package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	calendarop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/calendar"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *HandlerSuite) TestIsDateSelectionWeekendHolidayHandler() {
	expectedUrl := fmt.Sprintf("/calendar/%s/is-weekend-holiday/%s", "US", "2023-01-01")
	req := httptest.NewRequest("GET", expectedUrl, nil)

	params := calendarop.IsDateWeekendHolidayParams{
		HTTPRequest: req,
	}

	expectedCountryCode := "US"
	expectedCountryName := "United States"
	expectedDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockDateSelectionChecker := mocks.DateSelectionChecker{}
	info := services.IsDateWeekendHolidayInfo{}
	info.CountryCode = expectedCountryCode
	info.CountryName = expectedCountryName
	info.Date = expectedDate
	info.IsHoliday = true
	info.IsWeekend = true

	mockDateSelectionChecker.On("IsDateWeekendHoliday",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Time"),
	).Return(&info, nil)

	showHandler := IsDateWeekendHolidayHandler{suite.NewHandlerConfig(), &mockDateSelectionChecker}
	response := showHandler.Handle(params)

	suite.IsType(&calendarop.IsDateWeekendHolidayOK{}, response)
	okResponse := response.(*calendarop.IsDateWeekendHolidayOK)

	suite.Equal(expectedCountryCode, *okResponse.Payload.CountryCode)
	suite.Equal(expectedCountryName, *okResponse.Payload.CountryName)
	suite.Equal(strfmt.Date(expectedDate), *okResponse.Payload.Date)
	suite.True(*okResponse.Payload.IsHoliday)
	suite.True(*okResponse.Payload.IsWeekend)
}
