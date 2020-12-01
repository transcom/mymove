package mtoagent

import (
	"testing"

	"github.com/getlantern/deepcopy"
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOAgentServiceSuite) TestUpdateMTOAgentData() {
	// Set up the data needed for updateMTOAgentData obj
	checker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())
	oldAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())

	// Set up an agent model for successful tests
	var successAgent models.MTOAgent
	err := deepcopy.Copy(&successAgent, &oldAgent)
	suite.FatalNoError(err, "error while copying the old MTO Agent model to success model")

	// Set up an agent model for error tests
	var errorAgent models.MTOAgent
	err = deepcopy.Copy(&errorAgent, &oldAgent)
	suite.FatalNoError(err, "error while copying the old MTO Agent model to error model")

	// Test successful check for shipment ID
	suite.T().Run("checkShipmentID - success", func(t *testing.T) {
		agentData := updateMTOAgentData{
			updatedAgent:        successAgent, // as-is, should succeed
			oldAgent:            oldAgent,
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
		agentData := updateMTOAgentData{
			updatedAgent:        errorAgent,
			oldAgent:            oldAgent,
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
		var newAgentPrime models.MTOAgent
		err := deepcopy.Copy(&newAgentPrime, &oldAgentPrime)
		suite.FatalNoError(err, "error while copying Prime-available agent models")

		agentData := updateMTOAgentData{
			updatedAgent:        newAgentPrime,
			oldAgent:            oldAgentPrime,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err = agentData.checkPrimeAvailability()

		suite.NoError(err)
		suite.NoVerrs(agentData.verrs)
	})

	// Test unsuccessful check for prime availability
	suite.T().Run("checkPrimeAvailability - failure", func(t *testing.T) {
		agentData := updateMTOAgentData{
			updatedAgent:        errorAgent, // the default errorAgent should not be Prime-available
			oldAgent:            oldAgent,
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

		agentData := updateMTOAgentData{
			updatedAgent:        successAgent,
			oldAgent:            oldAgent,
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

		agentData := updateMTOAgentData{
			updatedAgent:        errorAgent,
			oldAgent:            oldAgent,
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
		agentData := updateMTOAgentData{
			updatedAgent:        successAgent, // as-is, should succeed
			oldAgent:            oldAgent,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := agentData.getVerrs()

		suite.NoError(err)
		suite.NoVerrs(agentData.verrs)
	})

	// Test getVerrs for unsuccessful example
	suite.T().Run("getVerrs - failure", func(t *testing.T) {
		agentData := updateMTOAgentData{
			updatedAgent:        errorAgent, // as-is, should fail
			oldAgent:            oldAgent,
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

	// Test setNewAgent for successful example
	suite.T().Run("setNewAgent - success", func(t *testing.T) {
		agentData := updateMTOAgentData{
			updatedAgent:        successAgent, // as-is, should succeed
			oldAgent:            oldAgent,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := agentData.getVerrs()

		suite.NoError(err)
		suite.NoVerrs(agentData.verrs)
	})
}
