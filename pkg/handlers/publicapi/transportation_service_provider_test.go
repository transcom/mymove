package publicapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/auth"

	tspop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/transportation_service_provider"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetTransportationServiceProviderHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusSUBMITTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]
	path := fmt.Sprintf("/shipments/%s/transportation_service_provider", shipment.ID.String())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := tspop.GetTransportationServiceProviderParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	// And: get shipment is returned
	handler := GetTransportationServiceProviderHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&tspop.GetTransportationServiceProviderOK{}, response)
	okResponse := response.(*tspop.GetTransportationServiceProviderOK)

	// And: Payload is equivalent to original shipment
	suite.Equal(strfmt.UUID(tspUser.TransportationServiceProviderID.String()), okResponse.Payload.ID)
}

func (suite *HandlerSuite) TestGetTransportationServiceProviderHandlerWhereSessionServiceMemberIDDoesNotEqualShipmentServiceMemberID() {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	path := fmt.Sprintf("/shipments/%s/transportation_service_provider", shipment.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateRequest(req, serviceMember)

	params := tspop.GetTransportationServiceProviderParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	handler := GetTransportationServiceProviderHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.NotEqual(session.ServiceMemberID, shipment.ServiceMemberID)
	suite.Assertions.IsType(&tspop.GetTransportationServiceProviderForbidden{}, response)
}
