package movetaskorder

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderChecker() {
	availableMTO := testdatagen.MakeAvailableMove(suite.DB())
	notAvailableMTO := testdatagen.MakeDefaultMove(suite.DB())
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
