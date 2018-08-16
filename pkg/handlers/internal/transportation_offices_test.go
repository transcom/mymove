package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestShowDutyStationTransportationOfficeHandler() {
	station := testdatagen.MakeDefaultDutyStation(suite.parent.Db)

	path := fmt.Sprintf("/duty_stations/%v/transportation_offices", station.ID.String())
	req := httptest.NewRequest("GET", path, nil)

	params := transportationofficeop.ShowDutyStationTransportationOfficeParams{
		HTTPRequest:   req,
		DutyStationID: *utils.FmtUUID(station.ID),
	}
	showHandler := ShowDutyStationTransportationOfficeHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := showHandler.Handle(params)

	suite.parent.Assertions.IsType(&transportationofficeop.ShowDutyStationTransportationOfficeOK{}, response)
	okResponse := response.(*transportationofficeop.ShowDutyStationTransportationOfficeOK)

	suite.parent.Assertions.Equal(station.TransportationOffice.ID.String(), okResponse.Payload.ID.String())
	suite.parent.Assertions.Equal(station.TransportationOffice.PhoneLines[0].Number, okResponse.Payload.PhoneLines[0])

}

func (suite *HandlerSuite) TestShowDutyStationTransportationOfficeHandlerNoOffice() {
	station := testdatagen.MakeDefaultDutyStation(suite.parent.Db)
	station.TransportationOffice = models.TransportationOffice{}
	station.TransportationOfficeID = nil
	suite.parent.MustSave(&station)

	path := fmt.Sprintf("/duty_stations/%v/transportation_offices", station.ID.String())
	req := httptest.NewRequest("GET", path, nil)

	params := transportationofficeop.ShowDutyStationTransportationOfficeParams{
		HTTPRequest:   req,
		DutyStationID: *utils.FmtUUID(station.ID),
	}
	showHandler := ShowDutyStationTransportationOfficeHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := showHandler.Handle(params)

	suite.parent.Assertions.IsType(&utils.ErrResponse{}, response)
	errResponse := response.(*utils.ErrResponse)

	suite.parent.Assertions.Equal(http.StatusNotFound, errResponse.Code)
}
