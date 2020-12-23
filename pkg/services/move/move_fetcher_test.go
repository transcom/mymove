package move

import (
	"errors"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveServiceSuite) TestMoveFetcher() {
	moveFetcher := NewMoveFetcher(suite.DB())

	suite.T().Run("successfully returns default draft move", func(t *testing.T) {
		expectedMove := testdatagen.MakeDefaultMove(suite.DB())

		actualMove, err := moveFetcher.FetchMove(expectedMove.Locator)
		suite.FatalNoError(err)

		suite.Equal(expectedMove.ID, actualMove.ID)
		suite.Equal(expectedMove.Locator, actualMove.Locator)
		suite.Equal(expectedMove.CreatedAt.Format(time.RFC3339), actualMove.CreatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.UpdatedAt.Format(time.RFC3339), actualMove.UpdatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.SubmittedAt, actualMove.SubmittedAt)
		suite.Equal(expectedMove.OrdersID, actualMove.OrdersID)
		suite.Equal(expectedMove.Status, actualMove.Status)
		suite.Equal(expectedMove.AvailableToPrimeAt, actualMove.AvailableToPrimeAt)
		suite.Equal(expectedMove.ContractorID, actualMove.ContractorID)
		suite.Equal(expectedMove.ReferenceID, actualMove.ReferenceID)
	})

	suite.T().Run("successfully returns submitted move available to prime", func(t *testing.T) {
		expectedMove := testdatagen.MakeAvailableMove(suite.DB())

		actualMove, err := moveFetcher.FetchMove(expectedMove.Locator)
		suite.FatalNoError(err)

		suite.Equal(expectedMove.ID, actualMove.ID)
		suite.Equal(expectedMove.Locator, actualMove.Locator)
		suite.Equal(expectedMove.CreatedAt.Format(time.RFC3339), actualMove.CreatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.UpdatedAt.Format(time.RFC3339), actualMove.UpdatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.SubmittedAt, actualMove.SubmittedAt)
		suite.Equal(expectedMove.OrdersID, actualMove.OrdersID)
		suite.Equal(expectedMove.Status, actualMove.Status)
		suite.Equal(expectedMove.AvailableToPrimeAt.Format(time.RFC3339), actualMove.AvailableToPrimeAt.Format(time.RFC3339))
		suite.Equal(expectedMove.ContractorID, actualMove.ContractorID)
		suite.Equal(expectedMove.ReferenceID, actualMove.ReferenceID)
	})

	suite.T().Run("returns not found error for unknown locator", func(t *testing.T) {
		_ = testdatagen.MakeAvailableMove(suite.DB())

		_, err := moveFetcher.FetchMove("QX97UY")
		suite.Error(err)
		suite.True(errors.Is(err, services.NotFoundError{}))
	})
}