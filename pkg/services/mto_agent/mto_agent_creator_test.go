package mtoagent

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOAgentServiceSuite) TestMTOAgentCreator() {
	// Set up NewMTOAgentCreator
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())
	mtoAgentCreator := NewMTOAgentCreator(suite.DB(), mtoChecker)

	// Create new mtoShipment with no agents
	move := testdatagen.MakeAvailableMove(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	const agentTypeReceiving = "RECEIVING_AGENT"
	const agentTypeReleasing = "RELEASING_AGENT"

	// Create valid Receiving Agent for the shipment
	receivingAgent := &models.MTOAgent{

		FirstName:     swag.String("Riley"),
		LastName:      swag.String("Baker"),
		MTOAgentType:  agentTypeReceiving,
		Email:         swag.String("rileybaker@example.com"),
		Phone:         swag.String("555-555-5555"),
		MTOShipmentID: mtoShipment.ID,
	}

	// Create valid Releasing Agent for the shipment
	releasingAgent := &models.MTOAgent{

		FirstName:     swag.String("Jason"),
		LastName:      swag.String("Ash"),
		MTOAgentType:  agentTypeReleasing,
		Email:         swag.String("jasonash@example.com"),
		Phone:         swag.String("555-555-5555"),
		MTOShipmentID: mtoShipment.ID,
	}

	suite.T().Run("CreateMTOAgentPrime - Receiving Agent - Success", func(t *testing.T) {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use established valid shipment and valid receiving agent
		// Expected:	New MTOAgent of type RECEIVING_AGENT is successfully created

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(receivingAgent)

		suite.Nil(err)
		suite.NotNil(createdAgent)
		suite.Equal(receivingAgent.ID, createdAgent.ID)
		suite.Equal(receivingAgent.FirstName, createdAgent.FirstName)
		suite.Equal(receivingAgent.LastName, createdAgent.LastName)
		suite.Equal(receivingAgent.Email, createdAgent.Email)
		suite.Equal(receivingAgent.Phone, createdAgent.Phone)
		suite.Equal(receivingAgent.MTOShipmentID, createdAgent.MTOShipmentID)

	})

	suite.T().Run("CreateMTOAgentPrime - Releasing Agent - Success", func(t *testing.T) {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use established valid shipment and valid releasing agent
		// Expected:	New MTOAgent of type RELEASING_AGENT is successfully created
		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(releasingAgent)

		suite.Nil(err)
		suite.NotNil(createdAgent)
		suite.Equal(releasingAgent.ID, createdAgent.ID)
		suite.Equal(releasingAgent.FirstName, createdAgent.FirstName)
		suite.Equal(releasingAgent.LastName, createdAgent.LastName)
		suite.Equal(releasingAgent.Email, createdAgent.Email)
		suite.Equal(releasingAgent.Phone, createdAgent.Phone)
		suite.Equal(releasingAgent.MTOShipmentID, createdAgent.MTOShipmentID)
	})

	suite.T().Run("Not Found Error", func(t *testing.T) {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use nonexistent mtoShipmentID and valid releasing agent
		// Expected:	NotFoundError is returned. Agent cannot be created without a valid shipment.
		notFoundUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")

		releasingAgent.MTOShipmentID = notFoundUUID

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(releasingAgent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID.String())

	})

	suite.T().Run("Conflict Error", func(t *testing.T) {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use same valid relesing agent and mtoShipmentID.
		// Expected:	ConflictError is returned. Only one agent of each type is allowed per shipment.
		releasingAgent.MTOShipmentID = mtoShipment.ID
		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(releasingAgent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
	})

	suite.T().Run("Conflict Error", func(t *testing.T) {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use same valid receiving agent and mtoShipmentID.
		// Expected:	ConflictError is returned. Only one agent of each type is allowed per shipment.
		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(receivingAgent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
	})

	suite.T().Run("Not Found Error, unavailable to Prime", func(t *testing.T) {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Create a new move and mtoShipment that is unavailable to Prime.
		// Expected:	NotFoundError is returned. Shipment must be available to Prime to add an agent.

		// Creates a shipment, which creates a move that is unavailable to Prime
		unavailableShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{})

		agent := testdatagen.MakeDefaultMTOAgent(suite.DB())
		agent.MTOShipmentID = unavailableShipment.ID

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(&agent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)

	})

	suite.T().Run("Validation Error", func(t *testing.T) {
		invalidAgent := &models.MTOAgent{
			FirstName:     swag.String("Riley"),
			LastName:      swag.String("Baker"),
			MTOShipmentID: mtoShipment.ID,
		}

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(invalidAgent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)

	})

}
