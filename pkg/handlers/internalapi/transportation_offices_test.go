package internalapi

import (
	"fmt"
	"net/http"

	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestShowDutyStationTransportationOfficeHandler() {
	station := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())

	path := fmt.Sprintf("/duty_stations/%v/transportation_offices", station.ID.String())
	req := suite.NewRequestWithContext("GET", path, nil)

	params := transportationofficeop.ShowDutyStationTransportationOfficeParams{
		HTTPRequest:   req,
		DutyStationID: *handlers.FmtUUID(station.ID),
	}
	showHandler := ShowDutyStationTransportationOfficeHandler{handlers.NewHandlerConfig(suite.DB(), suite.TestLogger())}
	response := showHandler.Handle(params)

	suite.Assertions.IsType(&transportationofficeop.ShowDutyStationTransportationOfficeOK{}, response)
	okResponse := response.(*transportationofficeop.ShowDutyStationTransportationOfficeOK)

	suite.Assertions.Equal(station.TransportationOffice.ID.String(), okResponse.Payload.ID.String())
	suite.Assertions.Equal(station.TransportationOffice.PhoneLines[0].Number, okResponse.Payload.PhoneLines[0])

}

func (suite *HandlerSuite) TestShowDutyStationTransportationOfficeHandlerNoOffice() {
	station := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	station.TransportationOffice = models.TransportationOffice{}
	station.TransportationOfficeID = nil
	suite.MustSave(&station)

	path := fmt.Sprintf("/duty_stations/%v/transportation_offices", station.ID.String())
	req := suite.NewRequestWithContext("GET", path, nil)

	params := transportationofficeop.ShowDutyStationTransportationOfficeParams{
		HTTPRequest:   req,
		DutyStationID: *handlers.FmtUUID(station.ID),
	}
	showHandler := ShowDutyStationTransportationOfficeHandler{handlers.NewHandlerConfig(suite.DB(), suite.TestLogger())}
	response := showHandler.Handle(params)

	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusNotFound, errResponse.Code)
}
