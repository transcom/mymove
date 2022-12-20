package internalapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	transportationofficeservice "github.com/transcom/mymove/pkg/services/transportation_office"
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
	showHandler := ShowDutyLocationTransportationOfficeHandler{suite.HandlerConfig()}
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
	showHandler := ShowDutyLocationTransportationOfficeHandler{suite.HandlerConfig()}
	response := showHandler.Handle(params)

	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusNotFound, errResponse.Code)
}

func (suite *HandlerSuite) TestGetTransportationOfficesHandler() {

	user := factory.BuildDefaultUser(suite.DB())

	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices", nil)
	req = suite.AuthenticateUserRequest(req, user)
	params := transportationofficeop.GetTransportationOfficesParams{
		HTTPRequest: req,
		Search:      "test",
	}

	handler := GetTransportationOfficesHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesOK{}, response)
	responsePayload := response.(*transportationofficeop.GetTransportationOfficesOK)

	// Validate outgoing payload
	suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
}

func (suite *HandlerSuite) TestGetTransportationOfficesHandlerUnauthorized() {
	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices", nil)
	params := transportationofficeop.GetTransportationOfficesParams{
		HTTPRequest: req,
		Search:      "test",
	}

	handler := GetTransportationOfficesHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher}

	// Request without authentication
	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesUnauthorized{}, response)
}

func (suite *HandlerSuite) TestGetTransportationOfficesHandlerForbidden() {
	officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices", nil)

	// Auth the request as an office user (forbidden)
	req = suite.AuthenticateOfficeRequest(req, officeUser)
	params := transportationofficeop.GetTransportationOfficesParams{
		HTTPRequest: req,
		Search:      "test",
	}

	handler := GetTransportationOfficesHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesForbidden{}, response)
}
