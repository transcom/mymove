package primeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoagent "github.com/transcom/mymove/pkg/services/mto_agent"

	"github.com/transcom/mymove/pkg/handlers"
)

type updateMTOAgentSubtestData struct {
	agent    models.MTOAgent
	newAgent models.MTOAgent
	req      *http.Request
	handler  UpdateMTOAgentHandler
	eTag     string
}

func (suite *HandlerSuite) makeUpdateMTOAgentSubtestData() (subtestData *updateMTOAgentSubtestData) {
	subtestData = &updateMTOAgentSubtestData{}
	// Set up db objects
	subtestData.agent = testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		Move: testdatagen.MakeAvailableMove(suite.DB()),
	})

	firstName := "Carol"
	lastName := "Romilly"
	email := "carol.romilly@example.com"
	phone := "456-555-7890"

	subtestData.newAgent = models.MTOAgent{
		FirstName: &firstName,
		LastName:  &lastName,
		Email:     &email,
		Phone:     &phone,
	}

	// Create handler and request
	subtestData.handler = UpdateMTOAgentHandler{
		suite.HandlerConfig(),
		mtoagent.NewMTOAgentUpdater(movetaskorder.NewMoveTaskOrderChecker()),
	}
	subtestData.req = httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/agents/%s", subtestData.agent.MTOShipmentID.String(), subtestData.agent.ID.String()), nil)

	subtestData.eTag = etag.GenerateEtag(subtestData.agent.UpdatedAt)

	return subtestData
}

func (suite *HandlerSuite) TestUpdateMTOAgentHandler() {

	// Test a successful request + update
	suite.Run("200 - OK response", func() {
		subtestData := suite.makeUpdateMTOAgentSubtestData()
		payload := payloads.MTOAgent(&subtestData.newAgent)
		params := mtoshipmentops.UpdateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			AgentID:       *handlers.FmtUUID(subtestData.agent.ID),
			MtoShipmentID: *handlers.FmtUUID(subtestData.agent.MTOShipmentID),
			Body:          payload,
			IfMatch:       subtestData.eTag,
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOAgentOK{}, response)

		// Check values
		agentOK := response.(*mtoshipmentops.UpdateMTOAgentOK)
		suite.Equal(agentOK.Payload.ID.String(), subtestData.agent.ID.String())
		suite.Equal(agentOK.Payload.MtoShipmentID.String(), subtestData.agent.MTOShipmentID.String())
		suite.Equal(string(agentOK.Payload.AgentType), string(subtestData.agent.MTOAgentType)) // wasn't updated, should be original value
		suite.Equal(agentOK.Payload.FirstName, subtestData.newAgent.FirstName)
		suite.Equal(agentOK.Payload.LastName, subtestData.newAgent.LastName)
		suite.Equal(agentOK.Payload.Email, subtestData.newAgent.Email)
		suite.Equal(agentOK.Payload.Phone, subtestData.newAgent.Phone)
	})

	// Test stale eTag
	suite.Run("412 - Precondition failed response", func() {
		subtestData := suite.makeUpdateMTOAgentSubtestData()
		// Let's test with the same valid values, but with a bad eTag
		payload := payloads.MTOAgent(&subtestData.newAgent)
		badETag := etag.GenerateEtag(subtestData.agent.UpdatedAt.Add(time.Duration(-10)))
		params := mtoshipmentops.UpdateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			AgentID:       *handlers.FmtUUID(subtestData.agent.ID),
			MtoShipmentID: *handlers.FmtUUID(subtestData.agent.MTOShipmentID),
			Body:          payload,
			IfMatch:       badETag, // stale
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOAgentPreconditionFailed{}, response)
	})

	// Test invalid IDs in the body vs. path values
	suite.Run("422 - Unprocessable response for bad ID values", func() {
		subtestData := suite.makeUpdateMTOAgentSubtestData()
		fakeUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")

		badAgent := subtestData.newAgent
		badAgent.ID = fakeUUID
		badAgent.MTOShipmentID = fakeUUID

		payload := payloads.MTOAgent(&badAgent)
		params := mtoshipmentops.UpdateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			AgentID:       *handlers.FmtUUID(subtestData.agent.ID),
			MtoShipmentID: *handlers.FmtUUID(subtestData.agent.MTOShipmentID),
			Body:          payload,
			IfMatch:       subtestData.eTag,
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOAgentUnprocessableEntity{}, response)

		// Check error message for the invalid fields
		agentUnprocessable := response.(*mtoshipmentops.UpdateMTOAgentUnprocessableEntity)
		_, okID := agentUnprocessable.Payload.InvalidFields["id"]
		_, okMTOShipmentID := agentUnprocessable.Payload.InvalidFields["mtoShipmentID"]
		suite.True(okID)
		suite.True(okMTOShipmentID)
	})

	// Test invalid input
	suite.Run("422 - Unprocessable response for invalid input", func() {
		subtestData := suite.makeUpdateMTOAgentSubtestData()
		empty := ""

		payload := payloads.MTOAgent(&subtestData.newAgent)
		payload.FirstName = &empty
		payload.Email = &empty
		payload.Phone = &empty

		params := mtoshipmentops.UpdateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			AgentID:       *handlers.FmtUUID(subtestData.agent.ID),
			MtoShipmentID: *handlers.FmtUUID(subtestData.agent.MTOShipmentID),
			Body:          payload,
			IfMatch:       subtestData.eTag,
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOAgentUnprocessableEntity{}, response)

		// Check error message for the invalid fields
		agentUnprocessable := response.(*mtoshipmentops.UpdateMTOAgentUnprocessableEntity)
		_, okFirstName := agentUnprocessable.Payload.InvalidFields["firstName"]
		_, okContactInfo := agentUnprocessable.Payload.InvalidFields["contactInfo"]
		suite.True(okFirstName)
		suite.True(okContactInfo)
	})

	// Test not found response
	suite.Run("404 - Not found response", func() {
		subtestData := suite.makeUpdateMTOAgentSubtestData()
		payload := payloads.MTOAgent(&subtestData.newAgent)
		params := mtoshipmentops.UpdateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			AgentID:       *handlers.FmtUUID(subtestData.agent.MTOShipmentID), // instead of agent.ID
			MtoShipmentID: *handlers.FmtUUID(subtestData.agent.MTOShipmentID),
			Body:          payload,
			IfMatch:       subtestData.eTag,
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOAgentNotFound{}, response)

		// Check error message for the incorrect ID
		agentNotFound := response.(*mtoshipmentops.UpdateMTOAgentNotFound)
		suite.Contains(*agentNotFound.Payload.Detail, subtestData.agent.MTOShipmentID.String())
	})

	// Test not Prime-available (not found response)
	suite.Run("404 - Not available response", func() {
		subtestData := suite.makeUpdateMTOAgentSubtestData()
		unavailableAgent := testdatagen.MakeDefaultMTOAgent(suite.DB()) // default is not available to Prime

		payload := payloads.MTOAgent(&unavailableAgent)
		params := mtoshipmentops.UpdateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			AgentID:       *handlers.FmtUUID(unavailableAgent.ID),
			MtoShipmentID: *handlers.FmtUUID(unavailableAgent.MTOShipmentID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(unavailableAgent.UpdatedAt),
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOAgentNotFound{}, response)

		// Check error message for the unavailable ID
		agentNotFound := response.(*mtoshipmentops.UpdateMTOAgentNotFound)
		suite.Contains(*agentNotFound.Payload.Detail, unavailableAgent.ID.String())
	})
}

type createMTOAgentSubtestData struct {
	move           models.Move
	mtoShipment    models.MTOShipment
	receivingAgent *primemessages.MTOAgent
	releasingAgent *primemessages.MTOAgent
	handler        CreateMTOAgentHandler
	req            *http.Request
}

func (suite *HandlerSuite) makeCreateMTOAgentSubtestData() (subtestData *createMTOAgentSubtestData) {
	subtestData = &createMTOAgentSubtestData{}

	// Create new mtoShipment with no agents
	subtestData.move = testdatagen.MakeAvailableMove(suite.DB())
	subtestData.mtoShipment = testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
		Move: subtestData.move,
	})

	const agentTypeReceiving = "RECEIVING_AGENT"
	const agentTypeReleasing = "RELEASING_AGENT"

	// Create valid Receiving Agent payload for the shipment
	subtestData.receivingAgent = &primemessages.MTOAgent{

		FirstName:     swag.String("Riley"),
		LastName:      swag.String("Baker"),
		AgentType:     agentTypeReceiving,
		Email:         swag.String("rileybaker@example.com"),
		Phone:         swag.String("555-555-5555"),
		MtoShipmentID: strfmt.UUID(subtestData.mtoShipment.ID.String()),
	}

	// Create valid Releasing Agent payload for the shipment
	subtestData.releasingAgent = &primemessages.MTOAgent{

		FirstName:     swag.String("Jason"),
		LastName:      swag.String("Ash"),
		AgentType:     agentTypeReleasing,
		Email:         swag.String("jasonash@example.com"),
		Phone:         swag.String("555-555-5555"),
		MtoShipmentID: strfmt.UUID(subtestData.mtoShipment.ID.String()),
	}

	// Create Handler
	subtestData.handler = CreateMTOAgentHandler{
		suite.HandlerConfig(),
		mtoagent.NewMTOAgentCreator(movetaskorder.NewMoveTaskOrderChecker()),
	}
	subtestData.req = httptest.NewRequest("POST", fmt.Sprintf("/mto-shipments/%s/agents", subtestData.mtoShipment.ID), nil)

	return subtestData
}

func (suite *HandlerSuite) TestCreateMTOAgentHandler() {
	suite.Run("200 - OK response Receiving Agent", func() {
		// Under test: 	CreateMTOAgentHandler, MTOAgentCreator
		// Set up: 		Pass in valid payload for a receiving agent.
		// Expected:	Handler returns 200 response with payload of new agent.

		subtestData := suite.makeCreateMTOAgentSubtestData()
		payload := subtestData.receivingAgent
		params := mtoshipmentops.CreateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			MtoShipmentID: *handlers.FmtUUID(subtestData.mtoShipment.ID),
			Body:          payload,
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOAgentOK{}, response)

		// Check Values
		agentOK := response.(*mtoshipmentops.CreateMTOAgentOK)
		suite.Equal(subtestData.receivingAgent.MtoShipmentID.String(), agentOK.Payload.MtoShipmentID.String())
		suite.Equal(string(subtestData.receivingAgent.AgentType), string(agentOK.Payload.AgentType)) // wasn't updated, should be original value
		suite.Equal(subtestData.receivingAgent.FirstName, agentOK.Payload.FirstName)
		suite.Equal(subtestData.receivingAgent.LastName, agentOK.Payload.LastName)
		suite.Equal(subtestData.receivingAgent.Email, agentOK.Payload.Email)
		suite.Equal(subtestData.receivingAgent.Phone, agentOK.Payload.Phone)

	})

	suite.Run("200 - OK response Releasing Agent", func() {
		// Under test: 	CreateMTOAgentHandler, MTOAgentCreator
		// Set up: 		Pass in valid payload for a releasing agent.
		// Expected:	Handler returns 200 response with payload of new agent.

		subtestData := suite.makeCreateMTOAgentSubtestData()
		payload := subtestData.releasingAgent
		params := mtoshipmentops.CreateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			MtoShipmentID: *handlers.FmtUUID(subtestData.mtoShipment.ID),
			Body:          payload,
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOAgentOK{}, response)

		// Check Values
		agentOK := response.(*mtoshipmentops.CreateMTOAgentOK)
		suite.Equal(subtestData.releasingAgent.MtoShipmentID.String(), agentOK.Payload.MtoShipmentID.String())
		suite.Equal(string(subtestData.releasingAgent.AgentType), string(agentOK.Payload.AgentType)) // wasn't updated, should be original value
		suite.Equal(subtestData.releasingAgent.FirstName, agentOK.Payload.FirstName)
		suite.Equal(subtestData.releasingAgent.LastName, agentOK.Payload.LastName)
		suite.Equal(subtestData.releasingAgent.Email, agentOK.Payload.Email)
		suite.Equal(subtestData.releasingAgent.Phone, agentOK.Payload.Phone)

	})

	suite.Run("404 - Not Found response", func() {
		subtestData := suite.makeCreateMTOAgentSubtestData()
		subtestData.releasingAgent.MtoShipmentID = "00000000-0000-0000-0000-000000000001"
		payload := subtestData.releasingAgent
		params := mtoshipmentops.CreateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			MtoShipmentID: subtestData.releasingAgent.MtoShipmentID,
			Body:          payload,
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOAgentNotFound{}, response)

	})

	suite.Run("409 - Conflict response", func() {
		// Under test: 	CreateMTOAgentHandler, MTOAgentCreator
		// Set up: 		Pass in valid payload for a receiving agent, and
		//				a shipment that already has an existing receiving agent.
		// Expected:	Handler returns 409 Conflict Error.

		subtestData := suite.makeCreateMTOAgentSubtestData()
		payload := subtestData.receivingAgent
		// set up the shipment and agent as already associated with
		// each other
		_, err := subtestData.handler.MTOAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), payloads.MTOAgentModel(payload))
		suite.NoError(err)
		params := mtoshipmentops.CreateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			MtoShipmentID: *handlers.FmtUUID(subtestData.mtoShipment.ID),
			Body:          payload,
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOAgentConflict{}, response)
	})

	suite.Run("422 - Unprocessable response for invalid input", func() {
		// Under test: 	CreateMTOAgentHandler, MTOAgentCreator
		// Set up: 		Pass an invalid payload for a releasing agent.
		// Expected:	Handler returns 422 Unprocessable Entity Error.
		subtestData := suite.makeCreateMTOAgentSubtestData()
		newMTOShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.move,
		})
		subtestData.releasingAgent.MtoShipmentID = strfmt.UUID(newMTOShipment.ID.String())
		empty := ""

		payload := subtestData.releasingAgent
		payload.FirstName = &empty
		payload.Email = &empty
		payload.Phone = &empty

		params := mtoshipmentops.CreateMTOAgentParams{
			HTTPRequest:   subtestData.req,
			MtoShipmentID: *handlers.FmtUUID(newMTOShipment.ID),
			Body:          payload,
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := subtestData.handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOAgentUnprocessableEntity{}, response)
	})
}
