package movetaskorder_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"

	. "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderChecker() {
	availableMTO := testdatagen.MakeAvailableMove(suite.DB())
	notAvailableMTO := testdatagen.MakeDefaultMove(suite.DB())
	mtoChecker := NewMoveTaskOrderChecker()

	suite.RunWithPreloadedData("MTO is available and visible - success", func() {
		availableToPrime, err := mtoChecker.MTOAvailableToPrime(suite.AppContextForTest(), availableMTO.ID)
		suite.Equal(availableToPrime, true)
		suite.NoError(err)
	})

	suite.RunWithPreloadedData("MTO is available but hidden - failure", func() {
		now := time.Now()
		hide := false
		availableHiddenMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
				Show:               &hide,
			},
		})

		availableToPrime, err := mtoChecker.MTOAvailableToPrime(suite.AppContextForTest(), availableHiddenMTO.ID)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), availableHiddenMTO.ID.String())
		suite.Equal(availableToPrime, false)
	})

	suite.RunWithPreloadedData("MTO is not available - no failure, but returns false", func() {
		availableToPrime, err := mtoChecker.MTOAvailableToPrime(suite.AppContextForTest(), notAvailableMTO.ID)
		suite.Equal(availableToPrime, false)
		suite.NoError(err)
	})

	suite.RunWithPreloadedData("MTO ID is not valid - failure", func() {
		badUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		availableToPrime, err := mtoChecker.MTOAvailableToPrime(suite.AppContextForTest(), badUUID)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), badUUID.String())
		suite.Equal(availableToPrime, false)
	})
}
