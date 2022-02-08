package movehistory

import (
	"testing"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveHistoryServiceSuite) TestMoveFetcher() {
	moveHistoryFetcher := NewMoveHistoryFetcher()

	suite.T().Run("successfully returns submitted move history available to prime", func(t *testing.T) {
		expectedMove := testdatagen.MakeAvailableMove(suite.DB())

		moveHistory, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), expectedMove.Locator)
		suite.FatalNoError(err)

		suite.Equal(expectedMove.ID, moveHistory.ID)
		suite.Equal(expectedMove.Locator, moveHistory.Locator)
		suite.Equal(expectedMove.ReferenceID, moveHistory.ReferenceID)
	})

	suite.T().Run("returns not found error for unknown locator", func(t *testing.T) {
		_ = testdatagen.MakeAvailableMove(suite.DB())

		_, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), "QX97UY")
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

}
