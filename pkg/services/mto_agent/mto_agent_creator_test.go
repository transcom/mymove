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

const agentTypeReceiving = "RECEIVING_AGENT"
const agentTypeReleasing = "RELEASING_AGENT"

func createAgentModel(firstName string, lastName string, agentType string, shipmentID uuid.UUID) *models.MTOAgent {
	// Create valid Receiving Agent for the shipment
	return &models.MTOAgent{
		FirstName:     swag.String(firstName),
		LastName:      swag.String(lastName),
		MTOAgentType:  agentTypeReceiving,
		Email:         swag.String("rileybaker@example.com"),
		Phone:         swag.String("555-555-5555"),
		MTOShipmentID: shipmentID,
	}
}

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

	suite.Run("CreateMTOAgentPrime - Receiving Agent - Success", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use established valid shipment and valid receiving agent
		mtoAgentCreator, shipment := setupTestData()
		receivingAgent := createAgentModel("Jason", "Ash", agentTypeReceiving, shipment.ID)

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
		releasingAgent := createAgentModel("Jason", "Ash", agentTypeReleasing, shipment.ID)

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
		mtoAgentCreator, _ := setupTestData()
		// Expected:	NotFoundError is returned. Agent cannot be created without a valid shipment.
		notFoundUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		releasingAgent := createAgentModel("Jason", "Ash", agentTypeReleasing, notFoundUUID)

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), releasingAgent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID.String())

	})

	suite.Run("Conflict Error", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Add a releasing agent, then attempt to add another agent of same type
		// Expected:	ConflictError is returned. Only one agent of each type is allowed per shipment.
		mtoAgentCreator, shipment := setupTestData()
		releasingAgent := createAgentModel("Jason", "Ash", agentTypeReleasing, shipment.ID)

		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), releasingAgent)
		suite.Nil(err)
		suite.NotNil(createdAgent)

		releasingAgent2 := createAgentModel("Grace", "Griffin", agentTypeReleasing, shipment.ID)

		// Expected:	ConflictError is returned. Only one agent of each type is allowed per shipment.
		createdAgent, err = mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), releasingAgent2)
		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.Run("Conflict Error", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Use same valid receiving agent and mtoShipmentID.
		mtoAgentCreator, shipment := setupTestData()
		receivingAgent := createAgentModel("Jason", "Ash", agentTypeReceiving, shipment.ID)
		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), receivingAgent)
		suite.Nil(err)
		suite.NotNil(createdAgent)

		receivingAgent2 := createAgentModel("Jason", "Ash", agentTypeReceiving, shipment.ID)
		// Expected:	ConflictError is returned. Only one agent of each type is allowed per shipment.
		createdAgent, err = mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), receivingAgent2)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.Run("Not Found Error, unavailable to Prime", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Create a new move and mtoShipment that is unavailable to Prime.
		// Expected:	NotFoundError is returned. Shipment must be available to Prime to add an agent.

		// Creates a shipment, which creates a move that is unavailable to Prime
		unavailableShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{})
		// Add a receiving agent on that shipment
		receivingAgent := createAgentModel("Jason", "Ash", agentTypeReceiving, unavailableShipment.ID)

		// Set up NewMTOAgentCreator
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
		mtoAgentCreator := NewMTOAgentCreator(mtoChecker)
		createdAgent, err := mtoAgentCreator.CreateMTOAgentPrime(suite.AppContextForTest(), receivingAgent)

		suite.Nil(createdAgent)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)

	})

	suite.Run("Validation Error", func() {
		// Under test:	CreateMTOAgentPrime
		// Set up:		Create an agent that is invalid due to missing info
		// Expected:	InvalidInput error is returned. No agent created.

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
