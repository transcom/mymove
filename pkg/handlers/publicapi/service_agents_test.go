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

func (suite *HandlerSuite) TestPatchServiceAgentHandler() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusACCEPTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]
	serviceAgents, _ := models.FetchServiceAgentsByTSP(suite.TestDB(), tspUser.TransportationServiceProviderID, shipment.ID)
	serviceAgent := serviceAgents[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/shipments/shipment_id/service_agents/service_agents_id", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	UpdatePayload := apimessages.ServiceAgent{
		PointOfContact: models.StringPointer("Not Jenny at ACME"),
		Email:          models.StringPointer("notjenny@example.com"),
		PhoneNumber:    models.StringPointer("3039035768"),
		Notes:          models.StringPointer("Some notes"),
	}

	params := serviceagentop.PatchServiceAgentParams{
		HTTPRequest:    req,
		ShipmentID:     strfmt.UUID(shipment.ID.String()),
		ServiceAgentID: strfmt.UUID(serviceAgent.ID.String()),
		Update:         &UpdatePayload,
	}

	// And: patch service agent is returned
	handler := PatchServiceAgentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&serviceagentop.PatchServiceAgentOK{}, response)
	okResponse := response.(*serviceagentop.PatchServiceAgentOK)

	// And: Payload has new values
	suite.Equal(strfmt.UUID(serviceAgent.ID.String()), okResponse.Payload.ID)
	suite.Equal(*UpdatePayload.PointOfContact, *okResponse.Payload.PointOfContact)
	suite.Equal(*UpdatePayload.Email, *okResponse.Payload.Email)
	suite.Equal(*UpdatePayload.PhoneNumber, *okResponse.Payload.PhoneNumber)
	suite.Equal(UpdatePayload.Notes, okResponse.Payload.Notes)
}

func (suite *HandlerSuite) TestPatchServiceAgentHandlerOnlyPOC() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusACCEPTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]
	serviceAgents, _ := models.FetchServiceAgentsByTSP(suite.TestDB(), tspUser.TransportationServiceProviderID, shipment.ID)
	serviceAgent := serviceAgents[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/shipments/shipment_id/service_agents/service_agents_id", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	UpdatePayload := apimessages.ServiceAgent{
		PointOfContact: models.StringPointer("Not Jenny at ACME"),
	}

	params := serviceagentop.PatchServiceAgentParams{
		HTTPRequest:    req,
		ShipmentID:     strfmt.UUID(shipment.ID.String()),
		ServiceAgentID: strfmt.UUID(serviceAgent.ID.String()),
		Update:         &UpdatePayload,
	}

	// And: patch service agent is returned
	handler := PatchServiceAgentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&serviceagentop.PatchServiceAgentOK{}, response)
	okResponse := response.(*serviceagentop.PatchServiceAgentOK)

	// And: Payload has new values
	suite.Equal(strfmt.UUID(serviceAgent.ID.String()), okResponse.Payload.ID)
	suite.Equal(*UpdatePayload.PointOfContact, *okResponse.Payload.PointOfContact)
	suite.Equal("jenny_acme@example.com", *okResponse.Payload.Email)
	suite.Equal("303-867-5309", *okResponse.Payload.PhoneNumber)
	suite.Nil(okResponse.Payload.Notes)
}

func (suite *HandlerSuite) TestPatchServiceAgentHandlerOnlyEmail() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusACCEPTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]
	serviceAgents, _ := models.FetchServiceAgentsByTSP(suite.TestDB(), tspUser.TransportationServiceProviderID, shipment.ID)
	serviceAgent := serviceAgents[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/shipments/shipment_id/service_agents/service_agents_id", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	UpdatePayload := apimessages.ServiceAgent{
		Email: models.StringPointer("notjenny@example.com"),
	}

	params := serviceagentop.PatchServiceAgentParams{
		HTTPRequest:    req,
		ShipmentID:     strfmt.UUID(shipment.ID.String()),
		ServiceAgentID: strfmt.UUID(serviceAgent.ID.String()),
		Update:         &UpdatePayload,
	}

	// And: patch service agent is returned
	handler := PatchServiceAgentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&serviceagentop.PatchServiceAgentOK{}, response)
	okResponse := response.(*serviceagentop.PatchServiceAgentOK)

	// And: Payload has new values
	suite.Equal(strfmt.UUID(serviceAgent.ID.String()), okResponse.Payload.ID)
	suite.Equal("Jenny at ACME Movers", *okResponse.Payload.PointOfContact)
	suite.Equal(*UpdatePayload.Email, *okResponse.Payload.Email)
	suite.Equal("303-867-5309", *okResponse.Payload.PhoneNumber)
	suite.Nil(okResponse.Payload.Notes)
}

func (suite *HandlerSuite) TestPatchServiceAgentHandlerOnlyPhoneNumber() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusACCEPTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]
	serviceAgents, _ := models.FetchServiceAgentsByTSP(suite.TestDB(), tspUser.TransportationServiceProviderID, shipment.ID)
	serviceAgent := serviceAgents[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/shipments/shipment_id/service_agents/service_agents_id", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	UpdatePayload := apimessages.ServiceAgent{
		Notes: models.StringPointer("Some notes"),
	}

	params := serviceagentop.PatchServiceAgentParams{
		HTTPRequest:    req,
		ShipmentID:     strfmt.UUID(shipment.ID.String()),
		ServiceAgentID: strfmt.UUID(serviceAgent.ID.String()),
		Update:         &UpdatePayload,
	}

	// And: patch service agent is returned
	handler := PatchServiceAgentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&serviceagentop.PatchServiceAgentOK{}, response)
	okResponse := response.(*serviceagentop.PatchServiceAgentOK)

	// And: Payload has new values
	suite.Equal(strfmt.UUID(serviceAgent.ID.String()), okResponse.Payload.ID)
	suite.Equal("Jenny at ACME Movers", *okResponse.Payload.PointOfContact)
	suite.Equal("jenny_acme@example.com", *okResponse.Payload.Email)
	suite.Equal("303-867-5309", *okResponse.Payload.PhoneNumber)
	suite.Equal(*UpdatePayload.Notes, *okResponse.Payload.Notes)
}

func (suite *HandlerSuite) TestPatchServiceAgentHandlerOnlyNotes() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusACCEPTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]
	serviceAgents, _ := models.FetchServiceAgentsByTSP(suite.TestDB(), tspUser.TransportationServiceProviderID, shipment.ID)
	serviceAgent := serviceAgents[0]

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/shipments/shipment_id/service_agents/service_agents_id", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	UpdatePayload := apimessages.ServiceAgent{
		PhoneNumber: models.StringPointer("3039035768"),
	}

	params := serviceagentop.PatchServiceAgentParams{
		HTTPRequest:    req,
		ShipmentID:     strfmt.UUID(shipment.ID.String()),
		ServiceAgentID: strfmt.UUID(serviceAgent.ID.String()),
		Update:         &UpdatePayload,
	}

	// And: patch service agent is returned
	handler := PatchServiceAgentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&serviceagentop.PatchServiceAgentOK{}, response)
	okResponse := response.(*serviceagentop.PatchServiceAgentOK)

	// And: Payload has new values
	suite.Equal(strfmt.UUID(serviceAgent.ID.String()), okResponse.Payload.ID)
	suite.Equal("Jenny at ACME Movers", *okResponse.Payload.PointOfContact)
	suite.Equal("jenny_acme@example.com", *okResponse.Payload.Email)
	suite.Equal(*UpdatePayload.PhoneNumber, *okResponse.Payload.PhoneNumber)
	suite.Nil(okResponse.Payload.Notes)
}

func (suite *HandlerSuite) TestPatchServiceAgentHandlerWrongTSP() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusACCEPTED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.TestDB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]
	serviceAgents, _ := models.FetchServiceAgentsByTSP(suite.TestDB(), tspUser.TransportationServiceProviderID, shipment.ID)
	serviceAgent := serviceAgents[0]

	otherTspUser := testdatagen.MakeDefaultTspUser(suite.TestDB())

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/shipments/shipment_id/service_agents/service_agents_id", nil)
	req = suite.AuthenticateTspRequest(req, otherTspUser)

	UpdatePayload := apimessages.ServiceAgent{
		PointOfContact: models.StringPointer("Not Jenny at ACME"),
		Email:          models.StringPointer("notjenny@example.com"),
		PhoneNumber:    models.StringPointer("3039035768"),
		Notes:          models.StringPointer("Some notes"),
	}

	params := serviceagentop.PatchServiceAgentParams{
		HTTPRequest:    req,
		ShipmentID:     strfmt.UUID(shipment.ID.String()),
		ServiceAgentID: strfmt.UUID(serviceAgent.ID.String()),
		Update:         &UpdatePayload,
	}

	// And: patch service agent is returned
	handler := PatchServiceAgentHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 400 status code
	suite.Assertions.IsType(&serviceagentop.PatchServiceAgentBadRequest{}, response)
}
