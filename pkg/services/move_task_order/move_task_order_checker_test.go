package movetaskorder_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"
	. "github.com/transcom/mymove/pkg/services/move_task_order"

	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderChecker() {
	now := time.Now()
	availableMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
	})
	notAvailableMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
	mtoChecker := NewMoveTaskOrderChecker(suite.DB())

	availableToPrime, err := mtoChecker.MTOAvailableToPrime(availableMTO.ID)
	suite.Equal(availableToPrime, true)
	suite.NoError(err)

	availableToPrime2, err2 := mtoChecker.MTOAvailableToPrime(notAvailableMTO.ID)
	suite.Equal(availableToPrime2, false)
	suite.NoError(err2)

	availableToPrime3, err3 := mtoChecker.MTOAvailableToPrime(uuid.Nil)
	suite.Error(err3)
	suite.IsType(err3, services.NotFoundError{})
	suite.Equal(availableToPrime3, false)
}
