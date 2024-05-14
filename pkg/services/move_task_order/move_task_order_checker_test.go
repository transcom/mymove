package movetaskorder_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	m "github.com/transcom/mymove/pkg/services/move_task_order"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderChecker() {
	mtoChecker := m.NewMoveTaskOrderChecker()

	suite.Run("MTO is available and visible - success", func() {
		availableMTO := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		availableToPrime, err := mtoChecker.MTOAvailableToPrime(suite.AppContextForTest(), availableMTO.ID)
		suite.Equal(availableToPrime, true)
		suite.NoError(err)
	})

	suite.Run("MTO is available but hidden - failure", func() {
		hide := false
		availableHiddenMTO := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: &hide,
				},
			},
		}, nil)

		availableToPrime, err := mtoChecker.MTOAvailableToPrime(suite.AppContextForTest(), availableHiddenMTO.ID)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), availableHiddenMTO.ID.String())
		suite.Equal(availableToPrime, false)
	})

	suite.Run("MTO is not available - no failure, but returns false", func() {
		notAvailableMTO := factory.BuildMove(suite.DB(), nil, nil)

		availableToPrime, err := mtoChecker.MTOAvailableToPrime(suite.AppContextForTest(), notAvailableMTO.ID)
		suite.Equal(availableToPrime, false)
		suite.NoError(err)
	})

	suite.Run("MTO ID is not valid - failure", func() {
		badUUID, _ := uuid.NewV4()
		availableToPrime, err := mtoChecker.MTOAvailableToPrime(suite.AppContextForTest(), badUUID)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), badUUID.String())
		suite.Equal(availableToPrime, false)
	})
}
