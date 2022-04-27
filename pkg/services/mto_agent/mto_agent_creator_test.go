package mtoagent

import (
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOAgentServiceSuite) TestMTOAgentCreator() {

	setupTestData := func() (services.MTOAgentCreator, models.MTOShipment) {
		// Set up NewMTOAgentCreator
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
		mtoAgentCreator := NewMTOAgentCreator(mtoChecker)

		// Create new mtoShipment with no agents
		move := testdatagen.MakeAvailableMove(suite.DB())
		mtoShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		return mtoAgentCreator, mtoShipment
	}

	const agentTypeReceiving = "RECEIVING_AGENT"
	const agentTypeReleasing = "RELEASING_AGENT"

	// Create valid Receiving Agent for the shipment
	receivingAgent := &models.MTOAgent{

		FirstName:    swag.String("Riley"),
		LastName:     swag.String("Baker"),
		MTOAgentType: agentTypeReceiving,
		Email:        swag.String("rileybaker@example.com"),
		Phone:        swag.String("555-555-5555"),
	}

	// Create valid Releasing Agent for the shipment
	releasingAgent := &models.MTOAgent{

		FirstName:    swag.String("Jason"),
		LastName:     swag.String("Ash"),
		MTOAgentType: agentTypeReleasing,
		Email:        swag.String("jasonash@example.com"),
		Phone:        swag.String("555-555-5555"),
	}

	suite.Run("CreateMTOAgentPrime - Receiving Agent - Success", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use established valid shipment and valid receiving agent
		mtoAgentCreator, shipment := setupTestData()
		receivingAgent.MTOShipmentID = shipment.ID

		// Expected:	New MTOAgent of type RECEIVING_AGENT is successfully created

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), receivingAgent)

		suite.Nil(err)
		suite.NotNil(createdAgent)
		suite.Equal(receivingAgent.FirstName, createdAgent.FirstName)
		suite.Equal(receivingAgent.LastName, createdAgent.LastName)
		suite.Equal(receivingAgent.Email, createdAgent.Email)
		suite.Equal(receivingAgent.Phone, createdAgent.Phone)
		suite.Equal(receivingAgent.MTOShipmentID, createdAgent.MTOShipmentID)

	})

	suite.Run("CreateMTOAgentPrime - Releasing Agent - Success", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use established valid shipment and valid releasing agent
		mtoAgentCreator, shipment := setupTestData()
		releasingAgent.MTOShipmentID = shipment.ID

		// Expected:	New MTOAgent of type RELEASING_AGENT is successfully created
		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), releasingAgent)

		suite.Nil(err)
		suite.NotNil(createdAgent)
		suite.Equal(releasingAgent.FirstName, createdAgent.FirstName)
		suite.Equal(releasingAgent.LastName, createdAgent.LastName)
		suite.Equal(releasingAgent.Email, createdAgent.Email)
		suite.Equal(releasingAgent.Phone, createdAgent.Phone)
		suite.Equal(releasingAgent.MTOShipmentID, createdAgent.MTOShipmentID)
	})

	suite.Run("Not Found Error", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use nonexistent mtoShipmentID and valid releasing agent
		mtoAgentCreator, shipment := setupTestData()
		releasingAgent.MTOShipmentID = shipment.ID

		// Expected:	NotFoundError is returned. Agent cannot be created without a valid shipment.
		notFoundUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")

		releasingAgent.MTOShipmentID = notFoundUUID

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), releasingAgent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID.String())

	})

	suite.Run("Conflict Error", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use same valid relesing agent and mtoShipmentID.
		// Expected:	ConflictError is returned. Only one agent of each type is allowed per shipment.
		mtoAgentCreator, shipment := setupTestData()
		releasingAgent.MTOShipmentID = shipment.ID

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), releasingAgent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.Run("Conflict Error", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use same valid receiving agent and mtoShipmentID.
		mtoAgentCreator, shipment := setupTestData()
		receivingAgent.MTOShipmentID = shipment.ID

		// Expected:	ConflictError is returned. Only one agent of each type is allowed per shipment.
		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), receivingAgent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.Run("Not Found Error, unavailable to Prime", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Create a new move and mtoShipment that is unavailable to Prime.
		mtoAgentCreator, _ := setupTestData()

		// Expected:	NotFoundError is returned. Shipment must be available to Prime to add an agent.

		// Creates a shipment, which creates a move that is unavailable to Prime
		unavailableShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{})

		// Create an agent associated with the shipment
		agent := testdatagen.MakeDefaultMTOAgent(suite.DB())
		agent.MTOShipmentID = unavailableShipment.ID

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), &agent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)

	})

	suite.Run("Validation Error", func() {
		mtoAgentCreator, shipment := setupTestData()

		invalidAgent := &models.MTOAgent{
			FirstName:     swag.String("Riley"),
			LastName:      swag.String("Baker"),
			MTOShipmentID: shipment.ID,
		}

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), invalidAgent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

	})

}
