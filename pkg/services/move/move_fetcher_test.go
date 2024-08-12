package move

import (
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *MoveServiceSuite) TestMoveFetcher() {
	moveFetcher := NewMoveFetcher()
	defaultSearchParams := services.MoveFetcherParams{}

	suite.Run("successfully returns default draft move", func() {
		expectedMove := factory.BuildMove(suite.DB(), nil, nil)

		actualMove, err := moveFetcher.FetchMove(suite.AppContextForTest(), expectedMove.Locator, &defaultSearchParams)
		suite.FatalNoError(err)

		suite.Equal(expectedMove.ID, actualMove.ID)
		suite.Equal(expectedMove.Locator, actualMove.Locator)
		suite.Equal(expectedMove.CreatedAt.Format(time.RFC3339), actualMove.CreatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.UpdatedAt.Format(time.RFC3339), actualMove.UpdatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.SubmittedAt, actualMove.SubmittedAt)
		suite.Equal(expectedMove.OrdersID, actualMove.OrdersID)
		suite.Equal(expectedMove.Status, actualMove.Status)
		suite.Equal(expectedMove.AvailableToPrimeAt, actualMove.AvailableToPrimeAt)
		suite.Equal(expectedMove.ApprovedAt, actualMove.ApprovedAt)
		suite.Equal(expectedMove.ContractorID, actualMove.ContractorID)
		suite.Equal(expectedMove.Contractor.ContractNumber, actualMove.Contractor.ContractNumber)
		suite.Equal(expectedMove.ReferenceID, actualMove.ReferenceID)
	})

	suite.Run("successfully returns submitted move available to prime", func() {
		expectedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		actualMove, err := moveFetcher.FetchMove(suite.AppContextForTest(), expectedMove.Locator, &defaultSearchParams)
		suite.FatalNoError(err)

		suite.Equal(expectedMove.ID, actualMove.ID)
		suite.Equal(expectedMove.Locator, actualMove.Locator)
		suite.Equal(expectedMove.CreatedAt.Format(time.RFC3339), actualMove.CreatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.UpdatedAt.Format(time.RFC3339), actualMove.UpdatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.SubmittedAt, actualMove.SubmittedAt)
		suite.Equal(expectedMove.OrdersID, actualMove.OrdersID)
		suite.Equal(expectedMove.Status, actualMove.Status)
		suite.Equal(expectedMove.AvailableToPrimeAt.Format(time.RFC3339), actualMove.AvailableToPrimeAt.Format(time.RFC3339))
		suite.Equal(expectedMove.ApprovedAt.Format(time.RFC3339), actualMove.ApprovedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.ContractorID, actualMove.ContractorID)
		suite.Equal(expectedMove.Contractor.Name, actualMove.Contractor.Name)
		suite.Equal(expectedMove.ReferenceID, actualMove.ReferenceID)
	})

	suite.Run("returns not found error for unknown locator", func() {
		_ = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		_, err := moveFetcher.FetchMove(suite.AppContextForTest(), "QX97UY", &defaultSearchParams)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns not found for a move that is marked hidden in the db", func() {
		hide := false
		hiddenMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: &hide,
				},
			},
		}, nil)
		locator := hiddenMove.Locator
		searchParams := services.MoveFetcherParams{
			IncludeHidden: false,
		}

		_, err := moveFetcher.FetchMove(suite.AppContextForTest(), locator, &searchParams)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns hidden move if explicit param is passed in", func() {
		hide := false
		actualMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: &hide,
				},
			},
		}, nil)
		locator := actualMove.Locator
		searchParams := services.MoveFetcherParams{
			IncludeHidden: true,
		}

		expectedMove, err := moveFetcher.FetchMove(suite.AppContextForTest(), locator, &searchParams)

		suite.FatalNoError(err)
		suite.Equal(expectedMove.ID, actualMove.ID)
		suite.Equal(expectedMove.Locator, actualMove.Locator)
	})
}
