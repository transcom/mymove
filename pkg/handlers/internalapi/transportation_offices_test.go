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
	"github.com/transcom/mymove/pkg/services/address"
	transportationofficeservice "github.com/transcom/mymove/pkg/services/transportation_office"
)

func (suite *HandlerSuite) TestShowDutyLocationTransportationOfficeHandler() {
	location := factory.FetchOrBuildCurrentDutyLocation(suite.DB())

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
	location := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
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
	officeUser := factory.BuildOfficeUser(nil, nil, nil)
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

func (suite *HandlerSuite) TestShowCounselingOfficesHandler() {
	user := factory.BuildDefaultUser(suite.DB())

	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	newAddress := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "CA",
		PostalCode:     "59801",
		County:         "County",
	}
	addressCreator := address.NewAddressCreator()
	createdAddress, err := addressCreator.CreateAddress(suite.AppContextForTest(), &newAddress)
	suite.NoError(err)

	origDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{
				AddressID:                  createdAddress.ID,
				ProvidesServicesCounseling: true,
			},
		},
		{
			Model: models.TransportationOffice{
				Name:             "PPPO Travis AFB - USAF",
				Gbloc:            "KKFA",
				ProvidesCloseout: true,
			},
		},
	}, nil)
	suite.MustSave(&origDutyLocation)

	path := fmt.Sprintf("/transportation_offices/%v/counseling_offices", origDutyLocation.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateUserRequest(req, user)
	params := transportationofficeop.ShowCounselingOfficesParams{
		HTTPRequest:    req,
		DutyLocationID: *handlers.FmtUUID(origDutyLocation.ID),
	}

	handler := ShowCounselingOfficesHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.ShowCounselingOfficesOK{}, response)
	responsePayload := response.(*transportationofficeop.ShowCounselingOfficesOK)

	// Validate outgoing payload
	suite.NoError(responsePayload.Payload.Validate(strfmt.Default))

}
