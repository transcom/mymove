package movetaskorder_test

import (
	"github.com/transcom/mymove/pkg/services"
	. "github.com/transcom/mymove/pkg/services/move_task_order"

	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderChecker() {
	now := time.Now()
	availableMTO := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			AvailableToPrimeAt: &now,
		},
	})
	notAvailableMTO := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
	mtoChecker := NewMoveTaskOrderChecker(suite.DB())

	availableToPrime, err := mtoChecker.MTOAvailableToPrime(availableMTO.ID)
	suite.Equal(availableToPrime, true)
	suite.NoError(err)

	availableToPrime2, expectedErr := mtoChecker.MTOAvailableToPrime(notAvailableMTO.ID)
	suite.Error(expectedErr)
	suite.IsType(expectedErr, services.InvalidInputError{})
	suite.Equal(availableToPrime2, false)
}
