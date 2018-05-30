package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestShowDutyStationTransportationOfficeHandler() {
	dutyStation, _ := testdatagen.MakeDutyStation(suite.db, "Air Station Yuma", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Yuma", State: "Arizona", PostalCode: "85364"})

	path := fmt.Sprintf("/duty_stations/%v/transportation_offices", dutyStation.ID.String())
	req := httptest.NewRequest("GET", path, nil)

	params := transportationofficeop.ShowDutyStationTransportationOfficeParams{
		HTTPRequest:   req,
		DutyStationID: *fmtUUID(dutyStation.ID),
	}
	showHandler := ShowDutyStationTransportationOfficeHandler(NewHandlerContext(suite.db, suite.logger))
	response := showHandler.Handle(params)

	suite.Assertions.IsType(&transportationofficeop.ShowDutyStationTransportationOfficeOK{}, response)
	okResponse := response.(*transportationofficeop.ShowDutyStationTransportationOfficeOK)

	suite.Assertions.Equal(dutyStation.TransportationOffice.ID.String(), okResponse.Payload.ID.String())
	suite.Assertions.Equal(dutyStation.TransportationOffice.PhoneLines[0].Number, okResponse.Payload.PhoneLines[0])

}

func (suite *HandlerSuite) TestShowDutyStationTransportationOfficeHandlerNoOffice() {

	station, _ := testdatagen.MakeDutyStationWithoutTransportationOffice(suite.db, "Air Station Yuma", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Yuma", State: "Arizona", PostalCode: "85364"})

	path := fmt.Sprintf("/duty_stations/%v/transportation_offices", station.ID.String())
	req := httptest.NewRequest("GET", path, nil)

	params := transportationofficeop.ShowDutyStationTransportationOfficeParams{
		HTTPRequest:   req,
		DutyStationID: *fmtUUID(station.ID),
	}
	showHandler := ShowDutyStationTransportationOfficeHandler(NewHandlerContext(suite.db, suite.logger))
	response := showHandler.Handle(params)

	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	suite.Assertions.Equal(http.StatusNotFound, errResponse.code)
}
