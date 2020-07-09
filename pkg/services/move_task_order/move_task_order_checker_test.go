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

	err := mtoChecker.IsAvailableToPrime(availableMTO.ID)
	suite.NoError(err)

	expectedErr := mtoChecker.IsAvailableToPrime(notAvailableMTO.ID)
	suite.Error(expectedErr)
	suite.IsType(expectedErr, services.InvalidInputError{})
}
