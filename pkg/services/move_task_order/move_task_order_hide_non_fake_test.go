package movetaskorder

import (
	"fmt"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_Hide() {
	mtoHider := NewMoveTaskOrderHider(suite.DB())
	suite.T().Run("Clear path", func(t *testing.T) {
		// Create move with service member using fake data
		result, err := mtoHider.Hide()
		suite.NoError(err)
		suite.Equal(0, len(result))
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelMTOAgent() {
	suite.T().Run("valid fake data", func(t *testing.T) {
		agent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{})
		result, err := isValidFakeModelMTOAgent(agent)
		suite.NoError(err)
		suite.Equal(true, result)
	})

	badFakeData := []testdatagen.Assertions{
		{MTOAgent: models.MTOAgent{FirstName: swag.String("Billy")}},
		{MTOAgent: models.MTOAgent{LastName: swag.String("Smith")}},
		{MTOAgent: models.MTOAgent{Phone: swag.String("111-111-1111")}},
		{MTOAgent: models.MTOAgent{Email: swag.String("billy@move.mil")}},
	}
	for idx, badData := range badFakeData {
		suite.T().Run(fmt.Sprintf("invalid fake data %d", idx), func(t *testing.T) {
			agent := testdatagen.MakeMTOAgent(suite.DB(), badData)
			result, err := isValidFakeModelMTOAgent(agent)
			suite.NoError(err)
			suite.Equal(false, result)
		})
	}
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelMTOShipment() {
	suite.T().Run("Clear path", func(t *testing.T) {
		result, err := isValidFakeModelMTOShipment(models.MTOShipment{})
		suite.NoError(err)
		suite.Equal(true, result)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelMTOShipments() {
	suite.T().Run("Clear path", func(t *testing.T) {
		result, err := isValidFakeModelMTOShipments(models.MTOShipments{})
		suite.NoError(err)
		suite.Equal(true, result)
	})
}
