package mtoagent

import (
	"testing"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOAgentServiceSuite) TestAgentValidationData() {
	// Set up the data needed for AgentValidationData obj
	checker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())
	oldAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())

	// Set up agent models for successful and unsuccessful tests
	successAgent := oldAgent
	errorAgent := oldAgent

	// Test successful check for shipment ID
	suite.T().Run("checkShipmentID - success", func(t *testing.T) {
		agentData := AgentValidationData{
			newAgent:            successAgent, // as-is, should succeed
			oldAgent:            &oldAgent,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := agentData.checkShipmentID()

		suite.NoError(err)
		suite.NoVerrs(agentData.verrs)
	})

	// Test unsuccessful check for shipment ID
	suite.T().Run("checkShipmentID - failure", func(t *testing.T) {
		errorAgent.MTOShipmentID = oldAgent.ID // set an invalid ID value
		agentData := AgentValidationData{
			newAgent:            errorAgent,
			oldAgent:            &oldAgent,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := agentData.checkShipmentID()

		suite.NoError(err)
		suite.True(agentData.verrs.HasAny())
		suite.Contains(agentData.verrs.Keys(), "mtoShipmentID")
	})

	// Test successful check for prime availability
	suite.T().Run("checkPrimeAvailability - success", func(t *testing.T) {
		oldAgentPrime := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()),
		})
		newAgentPrime := oldAgentPrime

		agentData := AgentValidationData{
			newAgent:            newAgentPrime,
			oldAgent:            &oldAgentPrime,
			shipment:            &oldAgentPrime.MTOShipment,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := agentData.checkPrimeAvailability()

		suite.NoError(err)
		suite.NoVerrs(agentData.verrs)
	})

	// Test unsuccessful check for prime availability
	suite.T().Run("checkPrimeAvailability - failure", func(t *testing.T) {
		agentData := AgentValidationData{
			newAgent:            errorAgent, // the default errorAgent should not be Prime-available
			oldAgent:            &oldAgent,
			shipment:            &oldAgent.MTOShipment,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := agentData.checkPrimeAvailability()

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.NoVerrs(agentData.verrs) // this check doesn't add a validation error
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
			newAgent:            successAgent,
			oldAgent:            &oldAgent,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := agentData.checkContactInfo()

		suite.NoError(err)
		suite.NoVerrs(agentData.verrs)
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
			newAgent:            errorAgent,
			oldAgent:            &oldAgent,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := agentData.checkContactInfo()

		suite.NoError(err)
		suite.True(agentData.verrs.HasAny())
		suite.Contains(agentData.verrs.Keys(), "firstName")
		suite.Contains(agentData.verrs.Keys(), "contactInfo")
	})

	// Test getVerrs for successful example
	suite.T().Run("getVerrs - success", func(t *testing.T) {
		agentData := AgentValidationData{
			newAgent:            successAgent, // as-is, should succeed
			oldAgent:            &oldAgent,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		_ = agentData.checkShipmentID() // this test should pass regardless of potential errors here
		_ = agentData.checkContactInfo()
		err := agentData.getVerrs()

		suite.NoError(err)
		suite.NoVerrs(agentData.verrs)
	})

	// Test getVerrs for unsuccessful example
	suite.T().Run("getVerrs - failure", func(t *testing.T) {
		agentData := AgentValidationData{
			newAgent:            errorAgent, // as-is, should fail
			oldAgent:            &oldAgent,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		_ = agentData.checkShipmentID() // this test should pass regardless of potential errors here
		_ = agentData.checkContactInfo()
		err := agentData.getVerrs()

		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		suite.True(agentData.verrs.HasAny())
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
			newAgent:            successAgent,
			oldAgent:            &oldAgent,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		newAgent := agentData.setFullAgent()

		suite.NoVerrs(agentData.verrs)
		suite.Equal(*newAgent.FirstName, *successAgent.FirstName)
		suite.Equal(*newAgent.Email, *successAgent.Email)
		suite.Nil(newAgent.Phone)
		// Checking that the old agent instances weren't changed:
		suite.NotEqual(*newAgent.FirstName, *oldAgent.FirstName)
		suite.NotEqual(*newAgent.FirstName, *agentData.oldAgent.FirstName)
		suite.NotNil(oldAgent.Phone)
		suite.NotNil(agentData.oldAgent.Phone)
		suite.Equal(*oldAgent.Phone, *agentData.oldAgent.Phone)
	})
}
