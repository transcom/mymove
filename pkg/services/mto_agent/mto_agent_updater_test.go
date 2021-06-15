package mtoagent

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"

	"github.com/transcom/mymove/pkg/services"
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

		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgentBasic(&notFoundAgent, eTag) // base validation

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID)
	})

	// Test validation error
	suite.T().Run("Validation Error", func(t *testing.T) {
		invalidAgent := newAgent
		invalidAgent.MTOShipmentID = newAgent.ID

		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgentBasic(&invalidAgent, eTag) // base validation

		suite.Nil(updatedAgent)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)

		invalidInputError := err.(services.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "mtoShipmentID")
	})

	// Test precondition failed (stale eTag)
	suite.T().Run("Precondition Failed", func(t *testing.T) {
		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgentBasic(&newAgent, "bloop") // base validation

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

		updatedAgent, err := mtoAgentUpdater.UpdateMTOAgentBasic(&newAgent, eTag) // base validation

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
