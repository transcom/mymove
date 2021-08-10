package movetaskorder_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/services"
	. "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderChecker() {
	availableMTO := testdatagen.MakeAvailableMove(suite.DB())
	notAvailableMTO := testdatagen.MakeDefaultMove(suite.DB())
	mtoChecker := NewMoveTaskOrderChecker()

	suite.RunWithRollback("MTO is available and visible - success", func() {
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		availableToPrime, err := mtoChecker.MTOAvailableToPrime(appCfg, availableMTO.ID)
		suite.Equal(availableToPrime, true)
		suite.NoError(err)
	})

	suite.RunWithRollback("MTO is available but hidden - failure", func() {
		now := time.Now()
		hide := false
		availableHiddenMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
				Show:               &hide,
			},
		})

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		availableToPrime, err := mtoChecker.MTOAvailableToPrime(appCfg, availableHiddenMTO.ID)
		suite.Error(err)
		suite.IsType(err, services.NotFoundError{})
		suite.Contains(err.Error(), availableHiddenMTO.ID.String())
		suite.Equal(availableToPrime, false)
	})

	suite.RunWithRollback("MTO is not available - no failure, but returns false", func() {
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		availableToPrime, err := mtoChecker.MTOAvailableToPrime(appCfg, notAvailableMTO.ID)
		suite.Equal(availableToPrime, false)
		suite.NoError(err)
	})

	suite.RunWithRollback("MTO ID is not valid - failure", func() {
		badUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		availableToPrime, err := mtoChecker.MTOAvailableToPrime(appCfg, badUUID)
		suite.Error(err)
		suite.IsType(err, services.NotFoundError{})
		suite.Contains(err.Error(), badUUID.String())
		suite.Equal(availableToPrime, false)
	})
}
