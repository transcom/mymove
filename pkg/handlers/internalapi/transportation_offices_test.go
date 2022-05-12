package internalapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestShowDutyLocationTransportationOfficeHandler() {
	location := testdatagen.FetchOrMakeDefaultCurrentDutyLocation(suite.DB())

	path := fmt.Sprintf("/duty_locations/%v/transportation_offices", location.ID.String())
	req := httptest.NewRequest("GET", path, nil)

	params := transportationofficeop.ShowDutyLocationTransportationOfficeParams{
		HTTPRequest:    req,
		DutyLocationID: *handlers.FmtUUID(location.ID),
	}
	showHandler := ShowDutyLocationTransportationOfficeHandler{handlers.NewHandlerConfig(suite.DB(), suite.Logger())}
	response := showHandler.Handle(params)

	suite.Assertions.IsType(&transportationofficeop.ShowDutyLocationTransportationOfficeOK{}, response)
	okResponse := response.(*transportationofficeop.ShowDutyLocationTransportationOfficeOK)

	suite.Assertions.Equal(location.TransportationOffice.ID.String(), okResponse.Payload.ID.String())
	suite.Assertions.Equal(location.TransportationOffice.PhoneLines[0].Number, okResponse.Payload.PhoneLines[0])

}

func (suite *HandlerSuite) TestShowDutyLocationTransportationOfficeHandlerNoOffice() {
	location := testdatagen.FetchOrMakeDefaultCurrentDutyLocation(suite.DB())
	location.TransportationOffice = models.TransportationOffice{}
	location.TransportationOfficeID = nil
	suite.MustSave(&location)

	path := fmt.Sprintf("/duty_locations/%v/transportation_offices", location.ID.String())
	req := httptest.NewRequest("GET", path, nil)

	params := transportationofficeop.ShowDutyLocationTransportationOfficeParams{
		HTTPRequest:    req,
		DutyLocationID: *handlers.FmtUUID(location.ID),
	}
	showHandler := ShowDutyLocationTransportationOfficeHandler{handlers.NewHandlerConfig(suite.DB(), suite.Logger())}
	response := showHandler.Handle(params)

	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusNotFound, errResponse.Code)
}
