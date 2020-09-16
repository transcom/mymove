package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMTOAgentValidation() {
	suite.T().Run("test valid MTOAgent", func(t *testing.T) {
		mtoShipmentID := uuid.Must(uuid.NewV4())
		mtoAgentID := uuid.Must(uuid.NewV4())

		validMTOAgent := models.MTOAgent{
			ID:            mtoAgentID,
			MTOShipmentID: mtoShipmentID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@testagent.agent"),
			Phone:         nil,
			MTOAgentType:  models.MTOAgentReleasing,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOAgent, expErrors)
	})
}

func (suite *ModelSuite) TestMTOAgentTestdataGen() {
	//test with no assertions
	MTOAgent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{})
	suite.Equal(MTOAgent.FirstName, swag.String("Agent First Name"))
	suite.Equal(MTOAgent.LastName, swag.String("Agent Last Name"))
	suite.Equal(MTOAgent.Email, swag.String("agent.email@ghc.gov"))
	suite.Equal(MTOAgent.Phone, swag.String("2025559301"))
	// test with assertions
	someFirstName := "Some other agent name"
	otherMTOAgent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			FirstName: &someFirstName,
		},
	})
	suite.Equal(otherMTOAgent.FirstName, swag.String("Some other agent name"))
}