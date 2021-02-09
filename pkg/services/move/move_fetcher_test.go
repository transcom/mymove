package move

import (
	"errors"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveServiceSuite) TestMoveFetcher() {
	moveFetcher := NewMoveFetcher(suite.DB())
	defaultSearchParams := services.MoveFetcherParams{}

	suite.T().Run("successfully returns default draft move", func(t *testing.T) {
		expectedMove := testdatagen.MakeDefaultMove(suite.DB())

		actualMove, err := moveFetcher.FetchMove(expectedMove.Locator, &defaultSearchParams)
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

		actualMove, err := moveFetcher.FetchMove(expectedMove.Locator, &defaultSearchParams)
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

		_, err := moveFetcher.FetchMove("QX97UY", &defaultSearchParams)
		suite.Error(err)
		suite.True(errors.Is(err, services.NotFoundError{}))
	})

	suite.T().Run("Returns not found for a move that is marked hidden in the db", func(t *testing.T) {
		hide := false
		hiddenMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Show: &hide,
			},
		})
		locator := hiddenMove.Locator
		searchParams := services.MoveFetcherParams{
			IncludeHidden: false,
		}

		_, err := moveFetcher.FetchMove(locator, &searchParams)

		suite.Error(err)
		suite.True(errors.Is(err, services.NotFoundError{}))
	})

	suite.T().Run("Returns hidden move if explicit param is passed in", func(t *testing.T) {
		hide := false
		actualMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Show: &hide,
			},
		})
		locator := actualMove.Locator
		searchParams := services.MoveFetcherParams{
			IncludeHidden: true,
		}

		expectedMove, err := moveFetcher.FetchMove(locator, &searchParams)

		suite.FatalNoError(err)
		suite.Equal(expectedMove.ID, actualMove.ID)
		suite.Equal(expectedMove.Locator, actualMove.Locator)
	})
}
