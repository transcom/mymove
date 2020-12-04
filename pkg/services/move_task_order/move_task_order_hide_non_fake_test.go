package movetaskorder

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
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
	suite.T().Run("Clear path", func(t *testing.T) {
		result, err := isValidFakeModelMTOAgent(models.MTOAgent{})
		suite.NoError(err)
		suite.Equal(true, result)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelMTOShipment() {
	suite.T().Run("Clear path", func(t *testing.T) {
		result, err := isValidFakeModelMTOShipment(models.MTOShipment{})
		suite.NoError(err)
		suite.Equal(true, result)
	})
}
