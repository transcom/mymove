package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	transportationofficeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/transportation_office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/address"
	transportationofficeservice "github.com/transcom/mymove/pkg/services/transportation_office"
)

func (suite *HandlerSuite) TestGetTransportationOfficesHandler() {
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "LRC Fort Knox",
				ProvidesCloseout: true,
			},
		},
	}, nil)
	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices", nil)
	params := transportationofficeop.GetTransportationOfficesParams{
		HTTPRequest: req,
		Search:      "LRC Fort Knox",
	}

	handler := GetTransportationOfficesHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesOK{}, response)
	responsePayload := response.(*transportationofficeop.GetTransportationOfficesOK)

	suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
	suite.Equal(transportationOffice.Name, *responsePayload.Payload[0].Name)
	suite.Equal(transportationOffice.Address.ID.String(), responsePayload.Payload[0].Address.ID.String())
	suite.Equal(transportationOffice.Gbloc, responsePayload.Payload[0].Gbloc)

}

func (suite *HandlerSuite) TestNoTransportationOfficesHandler() {

	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices", nil)
	params := transportationofficeop.GetTransportationOfficesParams{
		HTTPRequest: req,
		Search:      "LRC Fort Knox",
	}

	handler := GetTransportationOfficesHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesOK{}, response)
	responsePayload, ok := response.(*transportationofficeop.GetTransportationOfficesOK)

	suite.True(ok)
	suite.NotNil(responsePayload, "Response should not be nil")

}

func (suite *HandlerSuite) TestGetTransportationOfficesOpenHandler() {
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "NSF Dahlgren - East",
				ProvidesCloseout: true,
			},
		},
	}, nil)

	transportationOffice2 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "NSF Dahlgren - West",
				ProvidesCloseout: false,
			},
		},
	}, nil)

	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/open/transportation_offices", nil)
	params := transportationofficeop.GetTransportationOfficesOpenParams{
		HTTPRequest: req,
		Search:      "NSF Dahlgren - E",
	}

	handler := GetTransportationOfficesOpenHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesOpenOK{}, response)
	responsePayload := response.(*transportationofficeop.GetTransportationOfficesOpenOK)

	suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
	suite.Equal(transportationOffice.Name, *responsePayload.Payload[0].Name)
	suite.Equal(transportationOffice2.Name, *responsePayload.Payload[1].Name)
}

func (suite *HandlerSuite) TestGetTransportationOfficesGBLOCsHandler() {
	transportationOffice1 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "LRC Fort Knox",
				ProvidesCloseout: true,
				Gbloc:            "AGFM",
			},
		},
	}, nil)
	transportationOffice2 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "NSF End of Alphabet",
				ProvidesCloseout: true,
				Gbloc:            "WXYZ",
			},
		},
	}, nil)
	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices/gblocs", nil)
	params := transportationofficeop.GetTransportationOfficesGBLOCsParams{
		HTTPRequest: req,
	}

	handler := GetTransportationOfficesGBLOCsHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesGBLOCsOK{}, response)
	responsePayload := response.(*transportationofficeop.GetTransportationOfficesGBLOCsOK)

	suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
	suite.Equal(transportationOffice1.Gbloc, responsePayload.Payload[0])
	suite.Equal(transportationOffice2.Gbloc, responsePayload.Payload[1])
}

func (suite *HandlerSuite) TestShowCounselingOfficesHandler() {
	user := factory.BuildDefaultUser(suite.DB())

	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	newAddress := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "CA",
		PostalCode:     "59801",
		County:         models.StringPointer("County"),
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
