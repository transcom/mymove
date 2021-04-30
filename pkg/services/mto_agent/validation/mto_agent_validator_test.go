package mtoagentvalidate

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOAgentValidationServiceSuite) TestValidateAgent() {
	// Set up the data needed for AgentValidationData obj
	checker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())
	oldAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
	oldAgentPrime := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		Move: testdatagen.MakeAvailableMove(suite.DB()),
	})

	// Test with bad string key
	suite.T().Run("bad validatorKey - failure", func(t *testing.T) {
		agentData := AgentValidationData{}
		fakeKey := "FakeKey"
		updatedAgent, err := ValidateAgent(&agentData, fakeKey)

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.Contains(err.Error(), fakeKey)
	})

	// Test successful Base validation
	suite.T().Run("BasicAgentValidatorKey - success", func(t *testing.T) {
		newAgent := models.MTOAgent{
			ID:            oldAgent.ID,
			MTOShipmentID: oldAgent.MTOShipmentID,
		}
		agentData := AgentValidationData{
			NewAgent: newAgent,
			OldAgent: &oldAgent,
			Verrs:    validate.NewErrors(),
		}
		updatedAgent, err := ValidateAgent(&agentData, BasicAgentValidatorKey)

		suite.NoError(err)
		suite.NotNil(updatedAgent)
		suite.IsType(models.MTOAgent{}, *updatedAgent)
	})

	// Test unsuccessful Base validation
	suite.T().Run("BasicAgentValidatorKey - failure", func(t *testing.T) {
		newAgent := models.MTOAgent{
			ID:            oldAgent.ID,
			MTOShipmentID: oldAgent.ID, // bad value
		}
		agentData := AgentValidationData{
			NewAgent: newAgent,
			OldAgent: &oldAgent,
			Verrs:    validate.NewErrors(),
		}
		updatedAgent, err := ValidateAgent(&agentData, BasicAgentValidatorKey)

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	// Test successful Prime validation
	suite.T().Run("PrimeAgentValidatorKey - success", func(t *testing.T) {
		newAgentPrime := oldAgentPrime

		// Ensure we have the minimum required contact info
		firstName := "Carol"
		email := "test@example.com"
		newAgentPrime.FirstName = &firstName
		newAgentPrime.Email = &email

		agentData := AgentValidationData{
			NewAgent:            newAgentPrime,
			OldAgent:            &oldAgentPrime,
			Verrs:               validate.NewErrors(),
			AvailabilityChecker: checker,
			Shipment:            &oldAgentPrime.MTOShipment,
		}
		updatedAgent, err := ValidateAgent(&agentData, PrimeAgentValidatorKey)

		suite.NoError(err)
		suite.NotNil(updatedAgent)
		suite.IsType(models.MTOAgent{}, *updatedAgent)
	})

	// Test unsuccessful Prime validation - Not available to Prime
	suite.T().Run("PrimeAgentValidatorKey - not available failure", func(t *testing.T) {
		newAgent := models.MTOAgent{
			ID:            oldAgent.ID,
			MTOShipmentID: oldAgent.MTOShipmentID,
		}
		agentData := AgentValidationData{
			NewAgent:            newAgent,
			OldAgent:            &oldAgent, // this agent should not be Prime-available
			Verrs:               validate.NewErrors(),
			AvailabilityChecker: checker,
			Shipment:            &oldAgent.MTOShipment,
		}
		updatedAgent, err := ValidateAgent(&agentData, PrimeAgentValidatorKey)

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	// Test unsuccessful Prime validation - Invalid input
	suite.T().Run("PrimeAgentValidatorKey - invalid input failure", func(t *testing.T) {
		emptyString := ""
		newAgent := models.MTOAgent{
			ID:            oldAgentPrime.ID,
			MTOShipmentID: oldAgentPrime.ID,
			FirstName:     &emptyString,
			Email:         &emptyString,
			Phone:         &emptyString,
		}
		agentData := AgentValidationData{
			NewAgent:            newAgent,
			OldAgent:            &oldAgentPrime,
			Verrs:               validate.NewErrors(),
			AvailabilityChecker: checker,
			Shipment:            &oldAgentPrime.MTOShipment,
		}
		updatedAgent, err := ValidateAgent(&agentData, PrimeAgentValidatorKey)

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)

		invalidInputError := err.(services.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "mtoShipmentID")
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "firstName")
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "contactInfo")
	})

	// Test with empty string key (successful Base validation)
	suite.T().Run("empty validatorKey - success", func(t *testing.T) {
		newAgent := models.MTOAgent{
			MTOShipmentID: oldAgent.MTOShipmentID,
			ID:            oldAgent.ID,
		}
		agentData := AgentValidationData{
			NewAgent: newAgent,
			OldAgent: &oldAgent,
			Verrs:    validate.NewErrors(),
		}
		updatedAgent, err := ValidateAgent(&agentData, "")

		suite.NoError(err)
		suite.NotNil(updatedAgent)
		suite.IsType(models.MTOAgent{}, *updatedAgent)
	})
}
