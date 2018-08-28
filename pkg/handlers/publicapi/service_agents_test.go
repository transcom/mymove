package publicapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	serviceagentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/service_agents"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexServiceAgentsHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusACCEPTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	path := fmt.Sprintf("/shipments/%s/service_agents", shipment.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params := serviceagentop.IndexServiceAgentsParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	handler := IndexServiceAgentsHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&serviceagentop.IndexServiceAgentsOK{}, response)
	okResponse := response.(*serviceagentop.IndexServiceAgentsOK)

	suite.Equal(2, len(okResponse.Payload))
}

func (suite *HandlerSuite) TestCreateServiceAgentHandlerAllValues() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []models.ShipmentStatus{models.ShipmentStatusAWARDED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/shipment_id/service_agents", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	pointOfContact := "Pete and Repeat"

	email := "dogs@dogs.bones"
	notes := "This little piggy went to market"
	newServiceAgent := apimessages.ServiceAgent{
		Role:             apimessages.ServiceAgentRole(models.RoleORIGIN),
		PointOfContact:   handlers.FmtString(pointOfContact),
		Email:            swag.String(email),
		EmailIsPreferred: handlers.FmtBool(false),
		PhoneIsPreferred: handlers.FmtBool(true),
		Notes:            swag.String(notes),
	}
	params := serviceagentop.CreateServiceAgentParams{
		ServiceAgent: &newServiceAgent,
		ShipmentID:   strfmt.UUID(shipment.ID.String()),
		HTTPRequest:  req,
	}

	handler := CreateServiceAgentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&serviceagentop.CreateServiceAgentOK{}, response)
	okResponse := response.(*serviceagentop.CreateServiceAgentOK)

	suite.Equal(newServiceAgent.Role, okResponse.Payload.Role)
	suite.Equal(pointOfContact, *okResponse.Payload.PointOfContact)
	suite.Equal(*newServiceAgent.Email, *okResponse.Payload.Email)
	suite.Equal(*newServiceAgent.EmailIsPreferred, *okResponse.Payload.EmailIsPreferred)
	suite.Equal(*newServiceAgent.PhoneIsPreferred, *okResponse.Payload.PhoneIsPreferred)
	suite.Equal(*newServiceAgent.Notes, *okResponse.Payload.Notes)

	count, err := suite.TestDB().Where("shipment_id=$1", shipment.ID).Count(&models.ServiceAgent{})
	suite.Nil(err, "could not count service agents")
	suite.Equal(1, count)
}
