package movehistory

import (
	"fmt"
	"testing"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveHistoryServiceSuite) TestMoveFetcher() {
	moveHistoryFetcher := NewMoveHistoryFetcher()

	suite.T().Run("successfully returns default draft move history", func(t *testing.T) {
		expectedMove := testdatagen.MakeDefaultMove(suite.DB())

		moveHistory, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), expectedMove.Locator)
		suite.FatalNoError(err)

		fmt.Printf("moveHistory : +%v\n\n", moveHistory)

		suite.Equal(expectedMove.ID, moveHistory.ID)
		suite.Equal(expectedMove.Locator, moveHistory.Locator)
		suite.Equal(expectedMove.ReferenceID, moveHistory.ReferenceID)
		//suite.Equal(expectedMove.CreatedAt.Format(time.RFC3339), actualMove.CreatedAt.Format(time.RFC3339))
		//suite.Equal(expectedMove.UpdatedAt.Format(time.RFC3339), actualMove.UpdatedAt.Format(time.RFC3339))
		//suite.Equal(expectedMove.SubmittedAt, actualMove.SubmittedAt)
		//suite.Equal(expectedMove.OrdersID, actualMove.OrdersID)
		//suite.Equal(expectedMove.Status, actualMove.Status)
		//suite.Equal(expectedMove.AvailableToPrimeAt, actualMove.AvailableToPrimeAt)
		//suite.Equal(expectedMove.ContractorID, actualMove.ContractorID)

	})

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
