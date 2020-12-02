package mtoagent

import (
	"testing"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOAgentServiceSuite) TestMTOAgentUpdater() {
	// Set up the updater
	mtoAgentUpdater := NewMTOAgentUpdater(suite.DB())
	oldAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
	eTag := etag.GenerateEtag(oldAgent.UpdatedAt)

	newAgent := oldAgent

	// Test not found error
	suite.T().Run("Not Found Error", func(t *testing.T) {
		notFoundUUID := "00000000-0000-0000-0000-000000000001"
		notFoundAgent := newAgent
		notFoundAgent.ID = uuid.FromStringOrNil(notFoundUUID)

		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgent(&notFoundAgent, eTag, "") // base validation

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID)
	})

	// Test validation error
	suite.T().Run("Validation Error", func(t *testing.T) {
		invalidAgent := newAgent
		invalidAgent.MTOShipmentID = newAgent.ID

		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgent(&invalidAgent, eTag, "") // base validation

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)

		invalidInputError := err.(services.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "mtoShipmentID")
	})

	// Test precondition failed (stale eTag)
	suite.T().Run("Precondition Failed", func(t *testing.T) {
		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgent(&newAgent, "bloop", "") // base validation

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	// Test successful update
	suite.T().Run("Success", func(t *testing.T) {
		firstName := "Carol"
		lastName := "Romilly"
		email := "carol.romilly@example.com"

		newAgent.FirstName = &firstName
		newAgent.LastName = &lastName
		newAgent.Email = &email
		newAgent.Phone = nil // should keep the phone number from oldAgent

		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgent(&newAgent, eTag, "") // base validation

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

func (suite *MTOAgentServiceSuite) TestValidateUpdateMTOAgent() {
	// Set up the data needed for updateMTOAgentData obj
	checker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())
	oldAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
	oldAgentPrime := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		Move: testdatagen.MakeAvailableMove(suite.DB()),
	})

	// Test with bad string key
	suite.T().Run("bad validatorKey - failure", func(t *testing.T) {
		agentData := updateMTOAgentData{}
		fakeKey := "FakeKey"
		updatedAgent, err := ValidateUpdateMTOAgent(&agentData, fakeKey)

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.Contains(err.Error(), fakeKey)
	})

	// Test successful Base validation
	suite.T().Run("UpdateMTOAgentBaseValidator - success", func(t *testing.T) {
		newAgent := models.MTOAgent{
			ID:            oldAgent.ID,
			MTOShipmentID: oldAgent.MTOShipmentID,
		}
		agentData := updateMTOAgentData{
			updatedAgent: newAgent,
			oldAgent:     oldAgent,
			verrs:        validate.NewErrors(),
		}
		updatedAgent, err := ValidateUpdateMTOAgent(&agentData, UpdateMTOAgentBaseValidator)

		suite.NoError(err)
		suite.NotNil(updatedAgent)
		suite.IsType(models.MTOAgent{}, *updatedAgent)
	})

	// Test unsuccessful Base validation
	suite.T().Run("UpdateMTOAgentBaseValidator - failure", func(t *testing.T) {
		newAgent := models.MTOAgent{
			ID:            oldAgent.ID,
			MTOShipmentID: oldAgent.ID, // bad value
		}
		agentData := updateMTOAgentData{
			updatedAgent: newAgent,
			oldAgent:     oldAgent,
			verrs:        validate.NewErrors(),
		}
		updatedAgent, err := ValidateUpdateMTOAgent(&agentData, UpdateMTOAgentBaseValidator)

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	// Test successful Prime validation
	suite.T().Run("UpdateMTOAgentPrimeValidator - success", func(t *testing.T) {
		newAgentPrime := oldAgentPrime

		// Ensure we have the minimum required contact info
		firstName := "Carol"
		email := "test@example.com"
		newAgentPrime.FirstName = &firstName
		newAgentPrime.Email = &email

		agentData := updateMTOAgentData{
			updatedAgent:        newAgentPrime,
			oldAgent:            oldAgentPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedAgent, err := ValidateUpdateMTOAgent(&agentData, UpdateMTOAgentPrimeValidator)

		suite.NoError(err)
		suite.NotNil(updatedAgent)
		suite.IsType(models.MTOAgent{}, *updatedAgent)
	})

	// Test unsuccessful Prime validation - Not available to Prime
	suite.T().Run("UpdateMTOAgentPrimeValidator - not available failure", func(t *testing.T) {
		newAgent := models.MTOAgent{
			ID:            oldAgent.ID,
			MTOShipmentID: oldAgent.MTOShipmentID,
		}
		agentData := updateMTOAgentData{
			updatedAgent:        newAgent,
			oldAgent:            oldAgent, // this agent should not be Prime-available
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedAgent, err := ValidateUpdateMTOAgent(&agentData, UpdateMTOAgentPrimeValidator)

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	// Test unsuccessful Prime validation - Invalid input
	suite.T().Run("UpdateMTOAgentPrimeValidator - invalid input failure", func(t *testing.T) {
		emptyString := ""
		newAgent := models.MTOAgent{
			ID:            oldAgentPrime.ID,
			MTOShipmentID: oldAgentPrime.ID,
			FirstName:     &emptyString,
			Email:         &emptyString,
			Phone:         &emptyString,
		}
		agentData := updateMTOAgentData{
			updatedAgent:        newAgent,
			oldAgent:            oldAgentPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedAgent, err := ValidateUpdateMTOAgent(&agentData, UpdateMTOAgentPrimeValidator)

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
		}
		agentData := updateMTOAgentData{
			updatedAgent: newAgent,
			oldAgent:     oldAgent,
			verrs:        validate.NewErrors(),
		}
		updatedAgent, err := ValidateUpdateMTOAgent(&agentData, "")

		suite.NoError(err)
		suite.NotNil(updatedAgent)
		suite.IsType(models.MTOAgent{}, *updatedAgent)
	})
}
