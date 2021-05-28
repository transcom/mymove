package mtoagentvalidate

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOAgentValidateServiceSuite) TestAgentValidationData() {
	// Set up the data needed for AgentValidationData obj
	checker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())
	oldAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())

	// Set up agent models for successful and unsuccessful tests
	successAgent := oldAgent
	errorAgent := oldAgent

	// Test successful check for Shipment ID
	suite.T().Run("checkShipmentID - success", func(t *testing.T) {
		agentData := AgentValidationData{
			NewAgent:            successAgent, // as-is, should succeed
			OldAgent:            &oldAgent,
			AvailabilityChecker: checker,
			Verrs:               validate.NewErrors(),
		}
		err := agentData.checkShipmentID()

		suite.NoError(err)
		suite.NoVerrs(agentData.Verrs)
	})

	// Test unsuccessful check for Shipment ID
	suite.T().Run("checkShipmentID - failure", func(t *testing.T) {
		errorAgent.MTOShipmentID = oldAgent.ID // set an invalid ID value
		agentData := AgentValidationData{
			NewAgent:            errorAgent,
			OldAgent:            &oldAgent,
			AvailabilityChecker: checker,
			Verrs:               validate.NewErrors(),
		}
		err := agentData.checkShipmentID()

		suite.NoError(err)
		suite.True(agentData.Verrs.HasAny())
		suite.Contains(agentData.Verrs.Keys(), "mtoShipmentID")
	})

	// Test successful check for prime availability
	suite.T().Run("checkPrimeAvailability - success", func(t *testing.T) {
		oldAgentPrime := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()),
		})
		newAgentPrime := oldAgentPrime

		agentData := AgentValidationData{
			NewAgent:            newAgentPrime,
			OldAgent:            &oldAgentPrime,
			Shipment:            &oldAgentPrime.MTOShipment,
			AvailabilityChecker: checker,
			Verrs:               validate.NewErrors(),
		}
		err := agentData.checkPrimeAvailability()

		suite.NoError(err)
		suite.NoVerrs(agentData.Verrs)
	})

	// Test unsuccessful check for prime availability
	suite.T().Run("checkPrimeAvailability - failure", func(t *testing.T) {
		agentData := AgentValidationData{
			NewAgent:            errorAgent, // the default errorAgent should not be Prime-available
			OldAgent:            &oldAgent,
			Shipment:            &oldAgent.MTOShipment,
			AvailabilityChecker: checker,
			Verrs:               validate.NewErrors(),
		}
		err := agentData.checkPrimeAvailability()

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.NoVerrs(agentData.Verrs) // this check doesn't add a validation error
	})

	// Test successful check for contact info
	suite.T().Run("checkContactInfo - success", func(t *testing.T) {
		firstName := "Carol"
		lastName := ""
		email := ""
		phone := "234-555-4567"

		successAgent.FirstName = &firstName
		successAgent.LastName = &lastName
		successAgent.Email = &email
		successAgent.Phone = &phone

		agentData := AgentValidationData{
			NewAgent:            successAgent,
			OldAgent:            &oldAgent,
			AvailabilityChecker: checker,
			Verrs:               validate.NewErrors(),
		}
		err := agentData.checkContactInfo()

		suite.NoError(err)
		suite.NoVerrs(agentData.Verrs)
	})

	// Test unsuccessful check for contact info
	suite.T().Run("checkContactInfo - failure", func(t *testing.T) {
		firstName := ""
		email := ""
		phone := ""

		errorAgent.FirstName = &firstName
		errorAgent.Email = &email
		errorAgent.Phone = &phone

		agentData := AgentValidationData{
			NewAgent:            errorAgent,
			OldAgent:            &oldAgent,
			AvailabilityChecker: checker,
			Verrs:               validate.NewErrors(),
		}
		err := agentData.checkContactInfo()

		suite.NoError(err)
		suite.True(agentData.Verrs.HasAny())
		suite.Contains(agentData.Verrs.Keys(), "firstName")
		suite.Contains(agentData.Verrs.Keys(), "contactInfo")
	})

	// Test getVerrs for successful example
	suite.T().Run("getVerrs - success", func(t *testing.T) {
		agentData := AgentValidationData{
			NewAgent:            successAgent, // as-is, should succeed
			OldAgent:            &oldAgent,
			AvailabilityChecker: checker,
			Verrs:               validate.NewErrors(),
		}
		// These checks should not fail with an error
		err := agentData.checkShipmentID()
		suite.FatalNoError(err)

		err = agentData.checkContactInfo()
		suite.FatalNoError(err)

		err = agentData.getVerrs()
		suite.NoError(err)
		suite.NoVerrs(agentData.Verrs)
	})

	// Test getVerrs for unsuccessful example
	suite.T().Run("getVerrs - failure", func(t *testing.T) {
		agentData := AgentValidationData{
			NewAgent:            errorAgent, // as-is, should fail
			OldAgent:            &oldAgent,
			AvailabilityChecker: checker,
			Verrs:               validate.NewErrors(),
		}
		// These checks will find validation errors, but should not return other errors
		err := agentData.checkShipmentID()
		suite.FatalNoError(err)

		err = agentData.checkContactInfo()
		suite.FatalNoError(err)

		err = agentData.getVerrs()
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		suite.True(agentData.Verrs.HasAny())
	})

	// Test setFullAgent for successful example
	suite.T().Run("setFullAgent - success", func(t *testing.T) {
		firstName := "First"
		email := "email@email.email"
		phone := ""

		successAgent.FirstName = &firstName
		successAgent.Email = &email
		successAgent.Phone = &phone

		agentData := AgentValidationData{
			NewAgent:            successAgent,
			OldAgent:            &oldAgent,
			AvailabilityChecker: checker,
			Verrs:               validate.NewErrors(),
		}
		newAgent := agentData.setFullAgent()

		suite.NoVerrs(agentData.Verrs)
		suite.Equal(*newAgent.FirstName, *successAgent.FirstName)
		suite.Equal(*newAgent.Email, *successAgent.Email)
		suite.Nil(newAgent.Phone)
		// Checking that the old agent instances weren't changed:
		suite.NotEqual(*newAgent.FirstName, *oldAgent.FirstName)
		suite.NotEqual(*newAgent.FirstName, *agentData.OldAgent.FirstName)
		suite.NotNil(oldAgent.Phone)
		suite.NotNil(agentData.OldAgent.Phone)
		suite.Equal(*oldAgent.Phone, *agentData.OldAgent.Phone)
	})
}

func (suite *MTOAgentValidateServiceSuite) TestAgentValidationData_checkAgentID() {
	suite.T().Run("SUCCESS - When creating a new agent, ID should be nil", func(t *testing.T) {
		agentData := AgentValidationData{
			NewAgent: models.MTOAgent{
				// No ID because of create
			},
			Verrs: validate.NewErrors(),
		}

		err := agentData.checkAgentID()
		suite.NoError(err)
		suite.NoVerrs(agentData.Verrs)
	})

	suite.T().Run("FAIL - ID is set when creating a new agent", func(t *testing.T) {
		randomUUID, _ := uuid.NewV4()
		agentData := AgentValidationData{
			NewAgent: models.MTOAgent{
				ID: randomUUID,
			},
			Verrs: validate.NewErrors(),
		}

		err := agentData.checkAgentID()
		suite.NoError(err)
		suite.NotEmpty(agentData.Verrs)

		err = agentData.getVerrs()
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	suite.T().Run("SUCCESS - When updating an agent, old and new IDs should match", func(t *testing.T) {
		randomUUID, _ := uuid.NewV4()
		agentData := AgentValidationData{
			NewAgent: models.MTOAgent{
				ID: randomUUID,
			},
			OldAgent: &models.MTOAgent{
				ID: randomUUID,
			},
			Verrs: validate.NewErrors(),
		}

		err := agentData.checkAgentID()
		suite.NoError(err)
		suite.NoVerrs(agentData.Verrs)
	})

	suite.T().Run("FAIL - Old and new IDs do not match for update", func(t *testing.T) {
		agentData := AgentValidationData{
			NewAgent: models.MTOAgent{
				ID: uuid.Nil,
			},
			OldAgent: &models.MTOAgent{
				ID: uuid.Must(uuid.NewV4()),
			},
			Verrs: validate.NewErrors(),
		}

		err := agentData.checkAgentID()
		suite.Error(err)
		suite.IsType(services.ImplementationError{}, err)
		suite.NoVerrs(agentData.Verrs)
	})
}

func (suite *MTOAgentValidateServiceSuite) TestAgentValidationData_checkAgentType() {
	// Set up - no need to create real DB records, we just need models:
	randomUUID, _ := uuid.NewV4()
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{Stub: true})
	oldAgent := models.MTOAgent{
		ID:            randomUUID,
		MTOShipmentID: shipment.ID,
		MTOShipment:   shipment,
		MTOAgentType:  models.MTOAgentReceiving,
	}
	shipment.MTOAgents = models.MTOAgents{
		oldAgent,
	}

	suite.T().Run("SUCCESS - MTOAgentType was not changed", func(t *testing.T) {
		lastName := "Baker"
		agentData := AgentValidationData{
			NewAgent: models.MTOAgent{
				LastName: &lastName,
			},
			Shipment: &shipment,
			Verrs:    validate.NewErrors(),
		}

		err := agentData.checkAgentType()
		suite.NoError(err, "Unexpected error from checkAgentType with unchanged MTOAgentType")

		err = agentData.getVerrs()
		suite.NoError(err, "Unexpected validation error with unchanged MTOAgentType")
	})

	suite.T().Run("SUCCESS - MTOAgentType was changed to valid type", func(t *testing.T) {
		agentData := AgentValidationData{
			NewAgent: models.MTOAgent{
				ID:           oldAgent.ID,
				MTOAgentType: models.MTOAgentReleasing, // oldAgent is RECEIVING, so we're switching types
			},
			OldAgent: &oldAgent,
			Shipment: &shipment,
			Verrs:    validate.NewErrors(),
		}

		err := agentData.checkAgentType()
		suite.NoError(err, "Unexpected error from checkAgentType with new MTOAgentType")

		err = agentData.getVerrs()
		suite.NoError(err, "Unexpected validation error with new MTOAgentType")
	})

	suite.T().Run("FAIL - Shipment already has another agent with the same type", func(t *testing.T) {
		agentData := AgentValidationData{
			NewAgent: models.MTOAgent{
				// No ID because we're simulating a create
				MTOAgentType: models.MTOAgentReceiving, // oldAgent is RECEIVING, so this is the same type
			},
			// No old agent because this is for creating a new agent
			Shipment: &shipment,
			Verrs:    validate.NewErrors(),
		}

		err := agentData.checkAgentType()
		suite.Error(err, "Unexpectedly no error from checkAgentType with duplicated MTOAgentType")
		suite.IsType(services.ConflictError{}, err)
		suite.NoVerrs(agentData.Verrs)

		if err != nil {
			suite.Contains(err.Error(), models.MTOAgentReceiving)
		}
	})

	suite.T().Run("FAIL - Shipment already has the max number of agents", func(t *testing.T) {
		maxedShipment := shipment // Copy the shipment so we don't affect other tests
		secondAgent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			Stub: true,
			MTOAgent: models.MTOAgent{
				MTOShipmentID: maxedShipment.ID,
				MTOShipment:   maxedShipment,
				MTOAgentType:  models.MTOAgentReleasing,
			},
		})
		maxedShipment.MTOAgents = append(maxedShipment.MTOAgents, secondAgent)
		suite.Len(maxedShipment.MTOAgents, 2)

		agentData := AgentValidationData{
			NewAgent: models.MTOAgent{
				// No ID because we're simulating a create
				MTOAgentType: models.MTOAgentReceiving, // value doesn't matter, but we need one to validate
			},
			// No old agent because this is for creating a new agent
			Shipment: &maxedShipment,
			Verrs:    validate.NewErrors(),
		}

		err := agentData.checkAgentType()
		suite.Error(err, "Unexpectedly no error from checkAgentType with max number of agents")
		suite.IsType(services.ConflictError{}, err)
		suite.NoVerrs(agentData.Verrs)

		if err != nil {
			suite.Contains(err.Error(), "This shipment already has 2 agents - no more can be added")
		}
	})
}
