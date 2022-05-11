package mtoagent

import (
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOAgentServiceSuite) TestMTOAgentUpdater() {

	mtoAgentUpdater := NewMTOAgentUpdater(movetaskorder.NewMoveTaskOrderChecker())

	// Test not found error
	suite.Run("Not Found Error", func() {

		// TESTCASE SCENARIO
		// Under test: UpdateMTOAgentBasic function
		// Set up:     Update an agent that doesn't exist
		// Expected outcome: NotFound Error
		agent := testdatagen.MakeStubbedAgent(suite.DB(), testdatagen.Assertions{})
		agent.ID = uuid.Must(uuid.NewV4())

		eTag := etag.GenerateEtag(agent.UpdatedAt)
		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgentBasic(suite.AppContextForTest(), &agent, eTag) // base validation

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), agent.ID.String())
	})

	// Test validation error
	suite.Run("Validation Error", func() {
		// TESTCASE SCENARIO
		// Under test:  UpdateMTOAgentBasic function
		// Set up:      Create an agent
		//              Update the agent with a shipment that doesn't exist
		// Expected outcome: InvalidInput Error
		originalAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
		originalAgent.MTOShipmentID = uuid.Must(uuid.NewV4())

		eTag := etag.GenerateEtag(originalAgent.UpdatedAt)
		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgentBasic(suite.AppContextForTest(), &originalAgent, eTag) // base validation

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "mtoShipmentID")
	})

	// Test precondition failed (stale eTag)
	suite.Run("Precondition Failed", func() {

		// TESTCASE SCENARIO
		// Under test:  UpdateMTOAgentBasic function
		// Set up:      Create an agent, then update it with a bad etag
		// Expected outcome: PreconditionFailedError
		oldAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgentBasic(suite.AppContextForTest(), &oldAgent, "bloop") // base validation

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	// Test successful update
	suite.Run("Success", func() {
		// TESTCASE SCENARIO
		// Under test:  UpdateMTOAgentBasic function
		// Set up:      Create an agent, then update it successfully
		// Expected outcome: Success and an updated agent
		oldAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
		eTag := etag.GenerateEtag(oldAgent.UpdatedAt)

		firstName := "Carol"
		lastName := "Romilly"
		email := "carol.romilly@example.com"

		newAgent := models.MTOAgent{
			ID:        oldAgent.ID,
			FirstName: swag.String(firstName),
			LastName:  swag.String(lastName),
			Email:     swag.String(email),
			Phone:     nil,
		}

		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgentBasic(suite.AppContextForTest(), &newAgent, eTag) // base validation

		suite.NoError(err)
		suite.NotNil(updatedAgent)
		suite.Equal(updatedAgent.ID, oldAgent.ID)
		suite.Equal(updatedAgent.MTOShipmentID, oldAgent.MTOShipmentID)
		suite.Equal(updatedAgent.FirstName, newAgent.FirstName)
		suite.Equal(updatedAgent.LastName, newAgent.LastName)
		suite.Equal(updatedAgent.Email, newAgent.Email)
		suite.Equal(updatedAgent.Phone, oldAgent.Phone) // should not have been updated
	})
}
