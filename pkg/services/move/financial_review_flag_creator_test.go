package move

import (
	"errors"
	"testing"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveServiceSuite) TestFinancialReviewFlagCreator() {
	flagCreator := NewFinancialReviewFlagCreator()
	defaultFlagReason := "destination address is far from duty location"

	suite.T().Run("flag can be set", func(t *testing.T) {
		move := testdatagen.MakeDefaultMove(suite.DB())
		suite.Require().Equal(false, move.FinancialReviewRequested)
		suite.Require().Nil(move.FinancialReviewRequestedAt)
		suite.Require().Nil(move.FinancialReviewRemarks)
		m, err := flagCreator.CreateFinancialReviewFlag(suite.TestAppContext(), move.ID, defaultFlagReason)
		suite.NoError(suite.DB().Reload(&move))
		suite.Require().NotNil(m)
		suite.Require().NoError(err)
		suite.Require().Equal(true, move.FinancialReviewRequested)
		suite.Require().NotNil(move.FinancialReviewRequestedAt)
		suite.Require().Equal(defaultFlagReason, *move.FinancialReviewRemarks)
	})

	suite.T().Run("Wrong moveID should result in error", func(t *testing.T) {
		wrongUUID := uuid.Must(uuid.NewV4())

		_, err := flagCreator.CreateFinancialReviewFlag(suite.TestAppContext(), wrongUUID, defaultFlagReason)
		suite.Error(err)
		suite.Require().True(errors.As(err, &apperror.NotFoundError{}))
	})

	suite.T().Run("Empty remarks param should result in error", func(t *testing.T) {
		move := testdatagen.MakeDefaultMove(suite.DB())

		_, err := flagCreator.CreateFinancialReviewFlag(suite.TestAppContext(), move.ID, "")
		suite.Error(err)
		suite.Require().True(errors.As(err, &apperror.InvalidInputError{}))
	})

	suite.T().Run("setting flag after it has already been set should have no effect", func(t *testing.T) {
		move := testdatagen.MakeDefaultMove(suite.DB())
		// Make sure move starts out as we expect it to
		suite.Require().False(move.FinancialReviewRequested)
		suite.Require().Nil(move.FinancialReviewRequestedAt)
		suite.Require().Nil(move.FinancialReviewRemarks)

		// Set the flag once
		_, err := flagCreator.CreateFinancialReviewFlag(suite.TestAppContext(), move.ID, defaultFlagReason)
		suite.Require().NoError(err)
		suite.Require().NoError(suite.DB().Reload(&move))
		suite.Require().True(move.FinancialReviewRequested)
		suite.Require().NotNil(move.FinancialReviewRequestedAt)
		suite.Require().Equal(defaultFlagReason, *move.FinancialReviewRemarks)
		originalFlagTime := move.FinancialReviewRequestedAt

		// Attempt to set it again, and check to make sure nothing has changed
		_, err = flagCreator.CreateFinancialReviewFlag(suite.TestAppContext(), move.ID, "new reason")
		suite.Require().NoError(err)
		suite.Require().NoError(suite.DB().Reload(&move))
		suite.Require().True(move.FinancialReviewRequested)
		suite.Require().Equal(defaultFlagReason, *move.FinancialReviewRemarks)
		suite.Require().Equal(originalFlagTime, move.FinancialReviewRequestedAt)
	})
}
