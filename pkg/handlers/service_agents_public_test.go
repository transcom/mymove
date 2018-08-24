package handlers

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	publicserviceagentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/service_agents"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateServiceAgentHandlerAllValues() {
	numTspUsers := 1
	numShipments := 3
	numShipmentOfferSplit := []int{3}
	status := []models.ShipmentStatus{models.ShipmentStatusDRAFT}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("POST", "/shipments/shipment_id/service_agents", nil)
	req = suite.authenticateTspRequest(req, tspUser)

	pointOfContact := "Pete and Repeat"

	email := "dogs@dogs.bones"
	notes := "This little piggy went to market"
	newServiceAgent := apimessages.ServiceAgent{
		Role:             apimessages.ServiceAgentRole(models.RoleORIGIN),
		PointOfContact:   fmtString(pointOfContact),
		Email:            swag.String(email),
		EmailIsPreferred: fmtBool(false),
		PhoneIsPreferred: fmtBool(true),
		Notes:            swag.String(notes),
	}
	params := publicserviceagentop.CreateServiceAgentParams{
		ServiceAgent: &newServiceAgent,
		ShipmentID:   strfmt.UUID(shipment.ID.String()),
		HTTPRequest:  req,
	}

	handler := PublicCreateServiceAgentHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&publicserviceagentop.CreateServiceAgentOK{}, response)
	okResponse := response.(*publicserviceagentop.CreateServiceAgentOK)

	suite.Equal(newServiceAgent.Role, okResponse.Payload.Role)
	suite.Equal(pointOfContact, *okResponse.Payload.PointOfContact)
	suite.Equal(*newServiceAgent.Email, *okResponse.Payload.Email)
	suite.Equal(*newServiceAgent.EmailIsPreferred, *okResponse.Payload.EmailIsPreferred)
	suite.Equal(*newServiceAgent.PhoneIsPreferred, *okResponse.Payload.PhoneIsPreferred)
	suite.Equal(*newServiceAgent.Notes, *okResponse.Payload.Notes)

	count, err := suite.db.Where("shipment_id=$1", shipment.ID).Count(&models.ServiceAgent{})
	suite.Nil(err, "could not count service agents")
	suite.Equal(1, count)
}
