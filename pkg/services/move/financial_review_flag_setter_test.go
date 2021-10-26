package move

import (
	"errors"
	"testing"

	"github.com/transcom/mymove/pkg/etag"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveServiceSuite) TestFinancialReviewFlagSetter() {
	flagCreator := NewFinancialReviewFlagSetter()
	defaultFlagReason := "destination address is far from duty location"

	suite.T().Run("flag can be set", func(t *testing.T) {
		move := testdatagen.MakeDefaultMove(suite.DB())
		eTag := etag.GenerateEtag(move.UpdatedAt)

		suite.Require().Equal(false, move.FinancialReviewFlag)
		suite.Require().Nil(move.FinancialReviewFlagSetAt)
		suite.Require().Nil(move.FinancialReviewRemarks)
		m, err := flagCreator.SetFinancialReviewFlag(suite.TestAppContext(), move.ID, eTag, defaultFlagReason)
		suite.NoError(suite.DB().Reload(&move))
		suite.Require().NotNil(m)
		suite.Require().NoError(err)
		suite.Require().Equal(true, move.FinancialReviewFlag)
		suite.Require().NotNil(move.FinancialReviewFlagSetAt)
		suite.Require().Equal(defaultFlagReason, *move.FinancialReviewRemarks)
	})

	suite.T().Run("Wrong moveID should result in error", func(t *testing.T) {
		wrongUUID := uuid.Must(uuid.NewV4())

		_, err := flagCreator.SetFinancialReviewFlag(suite.TestAppContext(), wrongUUID, "", defaultFlagReason)
		suite.Error(err)
		suite.Require().True(errors.As(err, &apperror.NotFoundError{}))
	})

	suite.T().Run("Empty remarks param should result in error", func(t *testing.T) {
		move := testdatagen.MakeDefaultMove(suite.DB())
		eTag := etag.GenerateEtag(move.UpdatedAt)

		_, err := flagCreator.SetFinancialReviewFlag(suite.TestAppContext(), move.ID, eTag, "")
		suite.Error(err)
		suite.Require().True(errors.As(err, &apperror.InvalidInputError{}))
	})

	suite.T().Run("setting flag after it has already been set should have no effect", func(t *testing.T) {
		move := testdatagen.MakeDefaultMove(suite.DB())
		eTag := etag.GenerateEtag(move.UpdatedAt)
		// Make sure move starts out as we expect it to
		suite.Require().False(move.FinancialReviewFlag)
		suite.Require().Nil(move.FinancialReviewFlagSetAt)
		suite.Require().Nil(move.FinancialReviewRemarks)

		// Set the flag once
		_, err := flagCreator.SetFinancialReviewFlag(suite.TestAppContext(), move.ID, eTag, defaultFlagReason)
		suite.Require().NoError(err)
		suite.Require().NoError(suite.DB().Reload(&move))
		suite.Require().True(move.FinancialReviewFlag)
		suite.Require().NotNil(move.FinancialReviewFlagSetAt)
		suite.Require().Equal(defaultFlagReason, *move.FinancialReviewRemarks)
		originalFlagTime := move.FinancialReviewFlagSetAt

		suite.Require().NoError(suite.DB().Reload(&move))
		eTag = etag.GenerateEtag(move.UpdatedAt)

		// Attempt to set it again, and check to make sure nothing has changed
		_, err = flagCreator.SetFinancialReviewFlag(suite.TestAppContext(), move.ID, eTag, "new reason")
		suite.Require().NoError(err)
		suite.Require().NoError(suite.DB().Reload(&move))
		suite.Require().True(move.FinancialReviewFlag)
		suite.Require().Equal(defaultFlagReason, *move.FinancialReviewRemarks)
		suite.Require().Equal(originalFlagTime, move.FinancialReviewFlagSetAt)
	})
}
